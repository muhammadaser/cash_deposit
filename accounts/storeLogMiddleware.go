package accounts

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

func (mw *storeLogMiddleware) GetListAccounts() (accounts []Account, err error) {
	defer func() {
		mw.logger.Log("method", "GetListAccounts", "err", err)
	}()

	return mw.next.GetListAccounts()
}
func (mw *storeLogMiddleware) GetAccount(accountID string) (account Account, err error) {
	defer func() {
		mw.logger.Log("method", "GetAccount", "err", err)
	}()

	return mw.next.GetAccount(accountID)
}
func (mw *storeLogMiddleware) PostAccount(account Account) (err error) {
	defer func() {
		mw.logger.Log("method", "PostAccount", "err", err)
	}()

	return mw.next.PostAccount(account)
}
