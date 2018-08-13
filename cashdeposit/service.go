package cashdeposit

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
	ListDeposits() ([]CashDeposit, error)
	TotalBalance(accountID string) (TotalBalance, error)
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

func (s *setService) ListDeposits() ([]CashDeposit, error) {
	deposits, err := s.store.GetListDeposits()
	if err != nil {
		return deposits, ErrDatabase
	}
	return deposits, nil
}
func (s *setService) TotalBalance(accountID string) (TotalBalance, error) {
	balance, err := s.store.GetTotalBalance(accountID)
	if err != nil {
		return balance, err
	}
	return balance, nil
}
