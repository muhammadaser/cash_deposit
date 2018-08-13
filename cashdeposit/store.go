package cashdeposit

import (
	"time"

	"github.com/go-pg/pg"
)

// Store of products
type Store interface {
	GetListDeposits() ([]CashDeposit, error)
	GetTotalBalance(accountID string) (int64, error)
	// PostAccount(account Account) error
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

func (s *setStore) GetListDeposits() ([]CashDeposit, error) {
	cashDeposits := []CashDeposit{}
	_, err := s.pgDB.Query(&cashDeposits, "select * from public.cash_deposit")
	return nil, err
}
func (s *setStore) GetTotalBalance(accountID string) (int64, error) {
	tb := struct {
		balance int64
	}{}
	_, err := s.pgDB.Query(&tb, "select sum(deposit_amount) as balance from public.cash_deposit")
	return tb.balance, err
}

// func (s *setStore) GetAccount(accountID string) (Account, error) {
// 	account := Account{}
// 	_, err := s.pgDB.QueryOne(&account, "select * from public.account where account_id=?", accountID)
// 	return account, err
// }
// func (s *setStore) PostAccount(account Account) error {
// 	return s.pgDB.Insert(&account)
// }
