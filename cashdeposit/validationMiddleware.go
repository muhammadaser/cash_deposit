package cashdeposit

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

func (mw *validMiddleware) ListDeposits() ([]CashDeposit, error) {
	return mw.next.ListDeposits()
}
func (mw *validMiddleware) TotalBalance(accountID string) (balacne TotalBalance, err error) {
	if accountID == "" {
		return balacne, ErrBadRouting
	}
	return mw.next.TotalBalance(accountID)
}
func (mw *validMiddleware) NewDeposits(deposit CashDeposit) (err error) {
	_, err = govalidator.ValidateStruct(deposit)
	if err != nil {
		return err
	}
	return mw.next.NewDeposits(deposit)
}
