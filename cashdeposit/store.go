package cashdeposit

import (
	"time"

	"github.com/go-pg/pg"
)

// Store of products
type Store interface {
	GetListDeposits() ([]CashDeposit, error)
	GetListDepositsByAccount(accountID string) ([]CashDeposit, error)
	GetTotalBalance(accountID string) (TotalBalance, error)
	PostDeposit(deposit CashDeposit) error
}

// NewStore return struct that implement store interface
func NewStore(pgDB *pg.DB) Store {
	return &setStore{pgDB}
}

type setStore struct {
	pgDB *pg.DB
}

// CashDeposit merupakan tabel menyimpan data cash deposit nasabah
type CashDeposit struct {
	tableName struct{} `sql:"cash_deposit"`

	DepositID     string    `json:"deposit_id" sql:",pk" valid:"required"`
	AccountID     string    `json:"account_id" valid:"required"`
	DepositDate   time.Time `json:"deposit_date"`
	DepositAmount int64     `json:"deposit_amount" valid:"required"`
}

// TotalBalance nasabah berdasarkan account id
type TotalBalance struct {
	Balance int64 `json:"balance"`
}

func (s *setStore) GetListDeposits() ([]CashDeposit, error) {
	cashDeposits := []CashDeposit{}
	_, err := s.pgDB.Query(&cashDeposits, "select * from public.cash_deposit")
	return cashDeposits, err
}
func (s *setStore) GetListDepositsByAccount(accountID string) ([]CashDeposit, error) {
	cashDeposits := []CashDeposit{}
	_, err := s.pgDB.Query(
		&cashDeposits,
		"select * from public.cash_deposit where account_id = ?0",
		accountID,
	)
	return cashDeposits, err
}
func (s *setStore) GetTotalBalance(accountID string) (TotalBalance, error) {
	tb := TotalBalance{}
	_, err := s.pgDB.Query(
		&tb,
		"select sum(deposit_amount) as balance from public.cash_deposit where account_id=?0",
		accountID,
	)
	return tb, err
}
func (s *setStore) PostDeposit(deposit CashDeposit) error {
	return s.pgDB.Insert(&deposit)
}
