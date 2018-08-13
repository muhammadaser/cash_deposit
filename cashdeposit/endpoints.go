package cashdeposit

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

// Endpoints list of Service endpoint
type Endpoints struct {
	ListDepositsEndpoint          endpoint.Endpoint
	ListDepositsByAccountEndpoint endpoint.Endpoint
	TotalBalanceEndpoint          endpoint.Endpoint
	NewDepositEndpoint            endpoint.Endpoint
}

// NewEndpoint return Endpoints
func NewEndpoint(svc Service, logger log.Logger) Endpoints {
	var listDepositsEndpoint endpoint.Endpoint
	{
		listDepositsEndpoint = MakeListDepositsEndpoint(svc)
		listDepositsEndpoint = LoggingEndpointMiddleware(logger)(listDepositsEndpoint)
	}
	var listDepositsByAccountEndpoint endpoint.Endpoint
	{
		listDepositsByAccountEndpoint = MakeListDepositsByAccountEndpoint(svc)
		listDepositsByAccountEndpoint = LoggingEndpointMiddleware(logger)(listDepositsByAccountEndpoint)
	}
	var totalBalanceEndpoint endpoint.Endpoint
	{
		totalBalanceEndpoint = MakeTotalBalanceEndpoint(svc)
		totalBalanceEndpoint = LoggingEndpointMiddleware(logger)(totalBalanceEndpoint)
	}
	var newDepositsEndpoint endpoint.Endpoint
	{
		newDepositsEndpoint = MakeNewDepositsEndpoint(svc)
		newDepositsEndpoint = LoggingEndpointMiddleware(logger)(newDepositsEndpoint)
	}

	return Endpoints{
		ListDepositsEndpoint:          listDepositsEndpoint,
		ListDepositsByAccountEndpoint: listDepositsByAccountEndpoint,
		TotalBalanceEndpoint:          totalBalanceEndpoint,
		NewDepositEndpoint:            newDepositsEndpoint,
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

// MakeListDepositsByAccountEndpoint for list cash by account deposit
func MakeListDepositsByAccountEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listDepositsByAccountRequest)
		d, err := svc.ListDepositsByAcount(req.AccountID)
		res := listDepositsByAccountResponse{CashDeposit: d, Err: err}
		if err == nil {
			res.Code = http.StatusOK
			res.Message = "success"
		} else {
			res.Message = err.Error()
		}
		return res, nil
	}
}

// MakeTotalBalanceEndpoint for account endpoint
func MakeTotalBalanceEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(totalBalanceRequest)
		balance, err := svc.TotalBalance(req.AccountID)
		res := totalBalanceResponse{Balance: balance, Err: err}
		if err == nil {
			res.Code = http.StatusOK
			res.Message = "success"
		} else {
			res.Message = err.Error()
		}
		return res, nil
	}
}

// MakeNewDepositsEndpoint for New account
func MakeNewDepositsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(newDepositsRequest)
		err = svc.NewDeposits(req.Dp)
		res := newDepositsResponse{Err: err}
		if err == nil {
			res.Code = http.StatusOK
			res.Message = "success"
		} else {
			res.Message = err.Error()
		}
		return res, nil
	}
}

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

//====== List Account ======
type listDepositsByAccountRequest struct {
	AccountID string
}

type listDepositsByAccountResponse struct {
	Code        int           `json:"code,omitempty"`
	Message     string        `json:"message,omitempty"`
	CashDeposit []CashDeposit `json:"data"`
	Err         error         `json:"-"`
}

func (r listDepositsByAccountResponse) Failed() error { return r.Err }

//====== Account ======
type totalBalanceRequest struct {
	AccountID string
}

type totalBalanceResponse struct {
	Code    int          `json:"code,omitempty"`
	Message string       `json:"message,omitempty"`
	Balance TotalBalance `json:"data"`
	Err     error        `json:"-"`
}

func (r totalBalanceResponse) Failed() error { return r.Err }

//====== New Account ======
type newDepositsRequest struct {
	Dp CashDeposit
}

type newDepositsResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Err     error  `json:"-"`
}

func (r newDepositsResponse) Failed() error { return r.Err }
