package cashdeposit_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jonboulle/clockwork"

	"github.com/go-kit/kit/log"
	"github.com/muhammadaser/cash_deposit/accounts"
	"github.com/muhammadaser/cash_deposit/cashdeposit"
	mockStore "github.com/muhammadaser/cash_deposit/mocks/cashdeposit"
	"github.com/sebdah/goldie"
	"github.com/stretchr/testify/assert"
)

var clock = clockwork.NewFakeClock()
var baseUrl = "/cash-deposit/v1/deposits"
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
var listDepositsByAccount = []cashdeposit.CashDeposit{
	{
		DepositID:     "201808013849",
		AccountID:     "03212546",
		DepositDate:   clock.Now(),
		DepositAmount: 300000,
	},
	{
		DepositID:     "201808013850",
		AccountID:     "03212546",
		DepositDate:   clock.Now(),
		DepositAmount: 800000,
	},
	{
		DepositID:     "201808013851",
		AccountID:     "03212546",
		DepositDate:   clock.Now(),
		DepositAmount: 900000,
	},
}
var singleDeposit = listDeposits[0]
var singleDepositEmpty = cashdeposit.CashDeposit{}

func TestListDepositsByAccount(t *testing.T) {
	tests := map[string]struct {
		input      string
		output     []cashdeposit.CashDeposit
		err        error
		goldenFile string
	}{
		"Success": {
			input:      listDepositsByAccount[0].AccountID,
			output:     listDepositsByAccount,
			err:        nil,
			goldenFile: "testdata/list-deposits-by-account/success",
		},
		"Success while list empty": {
			input:      listDepositsByAccount[0].AccountID,
			output:     listDepositsEmpty,
			err:        nil,
			goldenFile: "testdata/list-deposits-by-account/success-whlie-empty",
		},
		"Failure While db error": {
			input:      listDepositsByAccount[0].AccountID,
			output:     listDepositsEmpty,
			err:        cashdeposit.ErrDatabase,
			goldenFile: "testdata/list-deposits-by-account/failure-db-error",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			res := httptest.NewRecorder()
			url := baseUrl + "/account/" + test.input
			req, err := http.NewRequest("GET", url, nil)
			assert.Nil(t, err)

			s := new(mockStore.Store)
			s.On("GetListDepositsByAccount", test.input).
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

func TestNewCashDeposit(t *testing.T) {

	singleDepositNoAccountID := singleDeposit
	singleDepositNoAccountID.AccountID = ""

	singleDepositZeroAmount := singleDeposit
	singleDepositZeroAmount.DepositAmount = 0

	tests := map[string]struct {
		input      cashdeposit.CashDeposit
		err        error
		goldenFile string
	}{
		"Success": {
			input:      singleDeposit,
			err:        nil,
			goldenFile: "testdata/new-cash-deposit/success",
		},
		"Failure While db error": {
			input:      singleDeposit,
			err:        accounts.ErrDatabase,
			goldenFile: "testdata/new-cash-deposit/failure-db-error",
		},
		"Failure While accountID not exits": {
			input:      singleDepositNoAccountID,
			err:        nil,
			goldenFile: "testdata/new-cash-deposit/failure-accountid-not-exist",
		},
		"Failure While deposit amount is zero": {
			input:      singleDepositZeroAmount,
			err:        nil,
			goldenFile: "testdata/new-cash-deposit/failure-amount-zero",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			res := httptest.NewRecorder()
			jsonString, err := json.Marshal(test.input)
			assert.Nil(t, err)
			req, err := http.NewRequest("POST", baseUrl, bytes.NewBuffer(jsonString))
			assert.Nil(t, err)

			s := new(mockStore.Store)
			s.On("PostDeposit", test.input).
				Return(test.err).
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
