package cashdeposit

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

func (mw *validMiddleware) ListDeposits() ([]CashDeposit, error) {
	return mw.next.ListDeposits()
}
func (mw *validMiddleware) TotalBalance(accountID string) (balacne TotalBalance, err error) {
	if accountID == "" {
		return balacne, ErrBadRouting
	}
	return mw.next.TotalBalance(accountID)
}

// func (mw *validMiddleware) Account(accountID string) (account Account, err error) {
// 	if accountID == "" {
// 		return account, ErrAccountNotFound
// 	}
// 	return mw.next.Account(accountID)
// }
// func (mw *validMiddleware) NewAccount(account Account) (err error) {
// 	_, err = govalidator.ValidateStruct(account)
// 	if err != nil {
// 		return err
// 	}
// 	return mw.next.NewAccount(account)
// }
