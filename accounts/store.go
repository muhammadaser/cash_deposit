package accounts

import "github.com/go-pg/pg"

// Store of products
type Store interface {
	GetListAccounts() (accounts []Account, err error)
	GetAccount(accountID string) (account Account, err error)
	PostAccount(account Account) error
}

// NewStore return struct that implement store interface
func NewStore(pgDB *pg.DB) Store {
	return &setStore{pgDB}
}

type setStore struct {
	pgDB *pg.DB
}

// Account nasabah yang di input oleh bank officer.
//  untuk saat ini data pada table account hanya seadanya,
//  perlu di tambah dan di sempurnakan.
//  data ini hanya untuk keperluan test.
type Account struct {
	tableName struct{} `sql:"account"`

	AccountID string `json:"account_id" sql:",pk" valid:"required"`
	FirstName string `json:"first_name" valid:"required"`
	LastName  string `json:"last_name" valid:"required"`
	Email     string `json:"email" valid:"email"`
	PhoneNo   string `json:"phone_no" valid:"required"`
	Address   string `json:"address" valid:"required"`
}

func (s *setStore) GetListAccounts() (accounts []Account, err error) {
	_, err = s.pgDB.Query(&accounts, "select * from public.account")
	return
}
func (s *setStore) GetAccount(accountID string) (Account, error) {
	account := Account{}
	_, err := s.pgDB.QueryOne(&account, "select * from public.account where account_id=?", accountID)
	return account, err
}
func (s *setStore) PostAccount(account Account) error {
	return s.pgDB.Insert(&account)
}
