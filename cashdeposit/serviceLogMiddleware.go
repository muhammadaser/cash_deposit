package cashdeposit

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

// ServiceLogMiddleware for log
type ServiceLogMiddleware func(next Service) Service

// NewServiceLogMiddleware return loggerMiddleware
func NewServiceLogMiddleware(logger log.Logger) ServiceLogMiddleware {
	return func(next Service) Service {
		return &serviceLogMiddleware{next, logger}
	}
}

type serviceLogMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw *serviceLogMiddleware) ListDeposits() (cd []CashDeposit, err error) {
	defer func() {
		mw.logger.Log("method", "ListDeposits", "err", err)
	}()

	return mw.next.ListDeposits()
}
func (mw *serviceLogMiddleware) TotalBalance(accountID string) (balacne TotalBalance, err error) {
	defer func() {
		mw.logger.Log("method", "TotalBalance", "accountID", accountID, "totalBalance", balacne.Balance, "err", err)
	}()

	return mw.next.TotalBalance(accountID)
}

func (mw *serviceLogMiddleware) NewDeposits(deposit CashDeposit) (err error) {
	defer func() {
		jsonString, _ := json.Marshal(deposit)
		mw.logger.Log("method", "NewDeposits", "deposit", jsonString, "err", err)
	}()

	return mw.next.NewDeposits(deposit)
}

// LoggingEndpointMiddleware returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func LoggingEndpointMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			defer func(begin time.Time) {
				logger.Log("transport_error", err, "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)

		}
	}
}
