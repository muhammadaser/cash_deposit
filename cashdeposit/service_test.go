package cashdeposit_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jonboulle/clockwork"

	"github.com/go-kit/kit/log"
	"github.com/muhammadaser/cash_deposit/cashdeposit"
	mockStore "github.com/muhammadaser/cash_deposit/mocks/cashdeposit"
	"github.com/sebdah/goldie"
	"github.com/stretchr/testify/assert"
)

var clock = clockwork.NewFakeClock()
var baseUrl = "/v1/cash-deposits"
var listDepositsEmpty = []cashdeposit.CashDeposit{}
var listDeposits = []cashdeposit.CashDeposit{
	{
		DepositID:     "201808013849",
		AccountID:     "03212546",
		DepositDate:   clock.Now(),
		DepositAmount: 100000,
	},
	{
		DepositID:     "201808013849",
		AccountID:     "03253548",
		DepositDate:   clock.Now(),
		DepositAmount: 100000,
	},
}
var singleDeposit = listDeposits[0]
var singleDepositEmpty = cashdeposit.CashDeposit{}

func TestListDeposits(t *testing.T) {
	tests := map[string]struct {
		output     []cashdeposit.CashDeposit
		err        error
		goldenFile string
	}{
		"Success": {
			output:     listDeposits,
			err:        nil,
			goldenFile: "testdata/list-deposits/success",
		},
		"Success while list empty": {
			output:     listDepositsEmpty,
			err:        nil,
			goldenFile: "testdata/list-deposits/success-whlie-empty",
		},
		"Failure While db error": {
			output:     listDepositsEmpty,
			err:        cashdeposit.ErrDatabase,
			goldenFile: "testdata/list-deposits/failure-db-error",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			res := httptest.NewRecorder()
			req, err := http.NewRequest("GET", baseUrl, nil)
			assert.Nil(t, err)

			s := new(mockStore.Store)
			s.On("GetListDeposits").
				Return(test.output, test.err).
				Maybe()

			svc := cashdeposit.NewService(s, clock)
			svc = cashdeposit.NewValidationMiddleware()(svc)

			endpoint := cashdeposit.NewEndpoint(svc, log.NewNopLogger())
			handler := cashdeposit.NewHTTPHandler(endpoint, log.NewNopLogger())

			handler.ServeHTTP(res, req)
			goldie.Assert(t, test.goldenFile, res.Body.Bytes())

			s.AssertExpectations(t)
		})
	}
}

func TestTotalBalance(t *testing.T) {
	totalBalance := cashdeposit.TotalBalance{Balance: 10000}
	totalBalanceZero := totalBalance
	totalBalanceZero.Balance = 0

	tests := map[string]struct {
		output     cashdeposit.TotalBalance
		input      string
		err        error
		goldenFile string
	}{
		"Success": {
			output:     totalBalance,
			input:      singleDeposit.AccountID,
			err:        nil,
			goldenFile: "testdata/total-balance/success",
		},
		"Success while balance zero": {
			output:     totalBalanceZero,
			input:      singleDeposit.AccountID,
			err:        nil,
			goldenFile: "testdata/total-balance/success-whlie-zero",
		},
		"Failure While db error": {
			output:     totalBalanceZero,
			input:      singleDeposit.AccountID,
			err:        cashdeposit.ErrDatabase,
			goldenFile: "testdata/total-balance/failure-db-error",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			res := httptest.NewRecorder()
			url := baseUrl + "/account/" + test.input + "/balance"
			req, err := http.NewRequest("GET", url, nil)
			assert.Nil(t, err)

			s := new(mockStore.Store)
			s.On("GetTotalBalance", test.input).
				Return(test.output, test.err).
				Maybe()

			svc := cashdeposit.NewService(s, clock)
			svc = cashdeposit.NewValidationMiddleware()(svc)

			endpoint := cashdeposit.NewEndpoint(svc, log.NewNopLogger())
			handler := cashdeposit.NewHTTPHandler(endpoint, log.NewNopLogger())

			handler.ServeHTTP(res, req)
			goldie.Assert(t, test.goldenFile, res.Body.Bytes())

			s.AssertExpectations(t)
		})
	}
}
