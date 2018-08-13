package cashdeposit

import "github.com/go-kit/kit/log"

// MailLogMiddleware for log
type MailLogMiddleware func(next Mail) Mail

// NewMailLogMiddleware return loggerMiddleware
func NewMailLogMiddleware(logger log.Logger) MailLogMiddleware {
	return func(next Mail) Mail {
		return &mailLogMiddleware{next, logger}
	}
}

type mailLogMiddleware struct {
	next   Mail
	logger log.Logger
}

func (mw *mailLogMiddleware) SendReceiptNotif(deposit CashDeposit) (err error) {
	defer func() {
		mw.logger.Log("method", "SendReceiptNotif", "err", err)
	}()

	return mw.next.SendReceiptNotif(deposit)
}
