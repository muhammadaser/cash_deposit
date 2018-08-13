package cashdeposit

import "github.com/go-kit/kit/log"

// StoreLogMiddleware for log
type StoreLogMiddleware func(next Store) Store

// NewStoreLogMiddleware return loggerMiddleware
func NewStoreLogMiddleware(logger log.Logger) StoreLogMiddleware {
	return func(next Store) Store {
		return &storeLogMiddleware{next, logger}
	}
}

type storeLogMiddleware struct {
	next   Store
	logger log.Logger
}

func (mw *storeLogMiddleware) GetListDeposits() (cd []CashDeposit, err error) {
	defer func() {
		mw.logger.Log("method", "GetListDeposits", "err", err)
	}()

	return mw.next.GetListDeposits()
}
func (mw *storeLogMiddleware) GetListDepositsByAccount(accountID string) (cd []CashDeposit, err error) {
	defer func() {
		mw.logger.Log("method", "GetListDepositsByAccount", "accountID", accountID, "err", err)
	}()

	return mw.next.GetListDepositsByAccount(accountID)
}
func (mw *storeLogMiddleware) GetTotalBalance(accountID string) (balance TotalBalance, err error) {
	defer func() {
		mw.logger.Log("method", "GetTotalBalance", "accountID", accountID, "balance", balance.Balance, "err", err)
	}()

	return mw.next.GetTotalBalance(accountID)
}
func (mw *storeLogMiddleware) PostDeposit(deposit CashDeposit) (err error) {
	defer func() {
		mw.logger.Log("method", "PostDeposit", "err", err)
	}()

	return mw.next.PostDeposit(deposit)
}
