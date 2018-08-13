package accounts

import "github.com/asaskevich/govalidator"

// ValidationMiddleware type for markup service
type ValidationMiddleware func(s Service) Service

// NewValidationMiddleware for validation service
func NewValidationMiddleware() ValidationMiddleware {
	return func(next Service) Service {
		return &validMiddleware{next}
	}
}

type validMiddleware struct {
	next Service
}

func (mw *validMiddleware) ListAccounts() ([]Account, error) {
	return mw.next.ListAccounts()
}
func (mw *validMiddleware) Account(accountID string) (account Account, err error) {
	if accountID == "" {
		return account, ErrAccountNotFound
	}
	return mw.next.Account(accountID)
}
func (mw *validMiddleware) NewAccount(account Account) (err error) {
	_, err = govalidator.ValidateStruct(account)
	if err != nil {
		return err
	}
	return mw.next.NewAccount(account)
}
