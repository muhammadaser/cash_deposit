package cashdeposit

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler")
)

// NewHTTPHandler mounts all of the service endpoints into an http.Handler.
func NewHTTPHandler(e Endpoints, logger log.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorLogger(logger),
	}

	listDepositsHandler := httptransport.NewServer(
		e.ListDepositsEndpoint,
		decodeHTTPListDeposits,
		encodeHTTPGenericResponse,
		options...,
	)
	totalBalanceHandler := httptransport.NewServer(
		e.TotalBalanceEndpoint,
		decodeHTTPTotalBalance,
		encodeHTTPGenericResponse,
		options...,
	)
	newDepositHandler := httptransport.NewServer(
		e.NewDepositEndpoint,
		decodeHTTPNewDeposit,
		encodeHTTPGenericResponse,
		options...,
	)

	r := mux.NewRouter()

	r.Handle("/v1/cash-deposits", listDepositsHandler).Methods("GET")
	r.Handle("/v1/cash-deposits", newDepositHandler).Methods("POST")
	r.Handle("/v1/cash-deposits/account/{accountID}/balance", totalBalanceHandler).Methods("GET")

	return r
}

func decodeHTTPListDeposits(_ context.Context, r *http.Request) (request interface{}, err error) {
	return listDepositsRequest{}, nil
}

func decodeHTTPTotalBalance(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	accountID, ok := vars["accountID"]
	if !ok {
		return nil, ErrBadRouting
	}
	return totalBalanceRequest{
		AccountID: accountID,
	}, nil
}

func decodeHTTPNewDeposit(_ context.Context, r *http.Request) (request interface{}, err error) {
	req := newDepositsRequest{}

	if e := json.NewDecoder(r.Body).Decode(&req.Dp); e != nil {
		return nil, e
	}

	return req, nil
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	code := err2code(err)
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error(), Code: code})
}

func err2code(err error) int {
	switch err {
	// case ErrAlreadyExists, ErrInconsistentIDs:
	// 	return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

type errorWrapper struct {
	Code  int    `json:"code,omitempty"`
	Error string `json:"messsage"`
}

// encodeHTTPGenericResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if f, ok := response.(Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}
