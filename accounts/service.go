package accounts

import (
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/go-pg/pg"
	"github.com/jonboulle/clockwork"
)

var (
	// ErrDatabase from database connection
	ErrDatabase = errors.New("Database Error")
	// ErrAccountNotFound for not exits account
	ErrAccountNotFound = errors.New("Account Not Found")
)

// Service of accounts
type Service interface {
	ListAccounts() (accounts []Account, err error)
	Account(accountID string) (account Account, err error)
	NewAccount(account Account) error
}

// New return Service
func New(logger log.Logger, pgDB *pg.DB) Service {
	var s Store
	{
		s = NewStore(pgDB)
		s = NewStoreLogMiddleware(logger)(s)
	}
	var svc Service
	{
		svc = NewService(s, clockwork.NewRealClock())
		svc = NewValidationMiddleware()(svc)
		svc = NewServiceLogMiddleware(logger)(svc)
	}
	return svc
}

type setService struct {
	store Store
	t     clockwork.Clock
}

// NewService return the struct that implements Service
func NewService(store Store, t clockwork.Clock) Service {
	return &setService{store, t}
}

func (s *setService) ListAccounts() (accounts []Account, err error) {
	accounts, err = s.store.GetListAccounts()
	if err != nil {
		return accounts, ErrDatabase
	}
	return accounts, nil
}
func (s *setService) Account(accountID string) (Account, error) {
	account, err := s.store.GetAccount(accountID)
	if err == pg.ErrNoRows {
		return account, ErrAccountNotFound
	}
	if err != nil {
		return account, ErrDatabase
	}
	return account, nil
}
func (s *setService) NewAccount(account Account) error {
	err := s.store.PostAccount(account)
	if err != nil {
		return ErrDatabase
	}
	return nil
}
