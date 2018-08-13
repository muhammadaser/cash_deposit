package cashdeposit

import (
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/go-pg/pg"
	"github.com/jonboulle/clockwork"

	"github.com/muhammadaser/cash_deposit/accounts"
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
	ListDepositsByAcount(accountID string) ([]CashDeposit, error)
	TotalBalance(accountID string) (TotalBalance, error)
	NewDeposits(deposit CashDeposit) error
}

// New return Service
func New(logger log.Logger, pgDB *pg.DB) Service {
	var s Store
	{
		s = NewStore(pgDB)
		s = NewStoreLogMiddleware(logger)(s)
	}
	var as accounts.Store
	{
		as = accounts.NewStore(pgDB)
		as = accounts.NewStoreLogMiddleware(logger)(as)
	}
	var mail Mail
	{
		mail = NewMail(as)
		mail = NewMailLogMiddleware(logger)(mail)
	}
	var svc Service
	{
		svc = NewService(s, clockwork.NewRealClock(), mail)
		svc = NewValidationMiddleware()(svc)
		svc = NewServiceLogMiddleware(logger)(svc)
	}
	return svc
}

type setService struct {
	store Store
	t     clockwork.Clock
	mail  Mail
}

// NewService return the struct that implements Service
func NewService(store Store, t clockwork.Clock, mail Mail) Service {
	return &setService{store, t, mail}
}

func (s *setService) ListDeposits() ([]CashDeposit, error) {
	deposits, err := s.store.GetListDeposits()
	if err != nil {
		return deposits, ErrDatabase
	}
	return deposits, nil
}
func (s *setService) ListDepositsByAcount(accountID string) ([]CashDeposit, error) {
	deposits, err := s.store.GetListDepositsByAccount(accountID)
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
func (s *setService) NewDeposits(deposit CashDeposit) error {
	deposit.DepositDate = s.t.Now()
	err := s.store.PostDeposit(deposit)
	if err != nil {
		return ErrDatabase
	}
	go s.mail.SendReceiptNotif(deposit)
	return nil
}
