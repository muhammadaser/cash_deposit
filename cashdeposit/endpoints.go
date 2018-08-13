package cashdeposit

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

// Endpoints list of Service endpoint
type Endpoints struct {
	ListDepositsEndpoint endpoint.Endpoint
	// AccountEndpoint      endpoint.Endpoint
	// NewAccountEndpoint   endpoint.Endpoint
}

// NewEndpoint return Endpoints
func NewEndpoint(svc Service, logger log.Logger) Endpoints {
	var listDepositsEndpoint endpoint.Endpoint
	{
		listDepositsEndpoint = MakeListDepositsEndpoint(svc)
		listDepositsEndpoint = LoggingEndpointMiddleware(logger)(listDepositsEndpoint)
	}
	// var accountEndpoint endpoint.Endpoint
	// {
	// 	accountEndpoint = MakeAccountEndpoint(svc)
	// 	accountEndpoint = LoggingEndpointMiddleware(logger)(accountEndpoint)
	// }
	// var newAccountEndpoint endpoint.Endpoint
	// {
	// 	newAccountEndpoint = MakeNewAccountEndpoint(svc)
	// 	newAccountEndpoint = LoggingEndpointMiddleware(logger)(newAccountEndpoint)
	// }

	return Endpoints{
		ListDepositsEndpoint: listDepositsEndpoint,
		// AccountEndpoint:      accountEndpoint,
		// NewAccountEndpoint:   newAccountEndpoint,
	}
}

// MakeListDepositsEndpoint for list cash deposit
func MakeListDepositsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// req := request.(listDepositsRequest)
		d, err := svc.ListDeposits()
		res := listDepositsResponse{CashDeposit: d, Err: err}
		if err == nil {
			res.Code = http.StatusOK
			res.Message = "success"
		} else {
			res.Message = err.Error()
		}
		return res, nil
	}
}

// // MakeAccountEndpoint for account endpoint
// func MakeAccountEndpoint(svc Service) endpoint.Endpoint {
// 	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
// 		req := request.(accountRequest)
// 		account, err := svc.Account(req.AccountID)
// 		res := accountResponse{Account: account, Err: err}
// 		if err == nil {
// 			res.Code = http.StatusOK
// 			res.Message = "success"
// 		} else if err == ErrAccountNotFound {
// 			res.Code = http.StatusOK
// 			res.Message = err.Error()
// 			res.Err = nil
// 		} else {
// 			res.Message = err.Error()
// 		}
// 		return res, nil
// 	}
// }

// // MakeNewAccountEndpoint for New account
// func MakeNewAccountEndpoint(svc Service) endpoint.Endpoint {
// 	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
// 		req := request.(newAccountRequest)
// 		err = svc.NewAccount(req.Ac)
// 		res := newAccountResponse{Err: err}
// 		if err == nil {
// 			res.Code = http.StatusOK
// 			res.Message = "success"
// 		} else {
// 			res.Message = err.Error()
// 		}
// 		return res, nil
// 	}
// }

// Failer is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so if they've
// failed, and if so encode them using a separate write path based on the error.
type Failer interface {
	Failed() error
}

//====== List Account ======
type listDepositsRequest struct {
}

type listDepositsResponse struct {
	Code        int           `json:"code,omitempty"`
	Message     string        `json:"message,omitempty"`
	CashDeposit []CashDeposit `json:"data"`
	Err         error         `json:"-"`
}

func (r listDepositsResponse) Failed() error { return r.Err }

// //====== Account ======
// type accountRequest struct {
// 	AccountID string
// }

// type accountResponse struct {
// 	Code    int     `json:"code,omitempty"`
// 	Message string  `json:"message,omitempty"`
// 	Account Account `json:"data"`
// 	Err     error   `json:"-"`
// }

// func (r accountResponse) Failed() error { return r.Err }

// //====== New Account ======
// type newAccountRequest struct {
// 	Ac Account
// }

// type newAccountResponse struct {
// 	Code    int    `json:"code,omitempty"`
// 	Message string `json:"message,omitempty"`
// 	Err     error  `json:"-"`
// }

// func (r newAccountResponse) Failed() error { return r.Err }
