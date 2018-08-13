package accounts_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-pg/pg"

	"github.com/jonboulle/clockwork"

	"github.com/go-kit/kit/log"
	"github.com/muhammadaser/cash_deposit/accounts"
	mockStore "github.com/muhammadaser/cash_deposit/mocks/accounts"
	"github.com/sebdah/goldie"
	"github.com/stretchr/testify/assert"
)

var clock = clockwork.NewFakeClock()
var baseUrl = "/v1/accounts"
var listAccountsEmty = []accounts.Account{}
var listAccounts = []accounts.Account{
	{
		AccountID: "03212546",
		FirstName: "Jhone",
		LastName:  "Doe",
		Email:     "jhone.doe@gmail.com",
		PhoneNo:   "085263123456",
		Address:   "Jakarta",
	},
	{
		AccountID: "03253548",
		FirstName: "Nathan",
		LastName:  "Dyer",
		Email:     "nathan.dyer@gmail.com",
		PhoneNo:   "085263123457",
		Address:   "Jakarta",
	},
}
var singleAccount = listAccounts[0]
var singleAccountEmpty = accounts.Account{}

func TestListAccounts(t *testing.T) {
	tests := map[string]struct {
		output     []accounts.Account
		err        error
		goldenFile string
	}{
		"Success": {
			output:     listAccounts,
			err:        nil,
			goldenFile: "testdata/list-accounts/success",
		},
		"Success while list empty": {
			output:     listAccountsEmty,
			err:        nil,
			goldenFile: "testdata/list-accounts/success-whlie-empty",
		},
		"Failure While db error": {
			output:     listAccountsEmty,
			err:        accounts.ErrDatabase,
			goldenFile: "testdata/list-accounts/failure-db-error",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			res := httptest.NewRecorder()
			req, err := http.NewRequest("GET", baseUrl, nil)
			assert.Nil(t, err)

			s := new(mockStore.Store)
			s.On("GetListAccounts").
				Return(test.output, test.err).
				Maybe()

			svc := accounts.NewService(s, clock)
			svc = accounts.NewValidationMiddleware()(svc)

			endpoint := accounts.NewEndpoint(svc, log.NewNopLogger())
			handler := accounts.NewHTTPHandler(endpoint, log.NewNopLogger())

			handler.ServeHTTP(res, req)
			goldie.Assert(t, test.goldenFile, res.Body.Bytes())

			s.AssertExpectations(t)
		})
	}
}
func TestAccount(t *testing.T) {

	tests := map[string]struct {
		input      string
		output     accounts.Account
		err        error
		goldenFile string
	}{
		"Success": {
			input:      singleAccount.AccountID,
			output:     singleAccount,
			err:        nil,
			goldenFile: "testdata/account/success",
		},
		"Success while account not exist": {
			input:      singleAccount.AccountID,
			output:     singleAccountEmpty,
			err:        pg.ErrNoRows,
			goldenFile: "testdata/account/success-whlie-empty",
		},
		"Failure While db error": {
			input:      singleAccount.AccountID,
			output:     singleAccountEmpty,
			err:        accounts.ErrDatabase,
			goldenFile: "testdata/account/failure-db-error",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			res := httptest.NewRecorder()
			req, err := http.NewRequest("GET", baseUrl+"/"+test.input, nil)
			assert.Nil(t, err)

			s := new(mockStore.Store)
			s.On("GetAccount", test.input).
				Return(test.output, test.err).
				Maybe()

			svc := accounts.NewService(s, clock)
			svc = accounts.NewValidationMiddleware()(svc)

			endpoint := accounts.NewEndpoint(svc, log.NewNopLogger())
			handler := accounts.NewHTTPHandler(endpoint, log.NewNopLogger())

			handler.ServeHTTP(res, req)
			goldie.Assert(t, test.goldenFile, res.Body.Bytes())

			s.AssertExpectations(t)
		})
	}
}
func TestNewAccount(t *testing.T) {

	singleAccountNoAccountID := singleAccount
	singleAccountNoAccountID.AccountID = ""

	singleAccountNoName := singleAccount
	singleAccountNoName.FirstName = ""
	singleAccountNoName.LastName = ""

	tests := map[string]struct {
		input      accounts.Account
		err        error
		goldenFile string
	}{
		"Success": {
			input:      singleAccount,
			err:        nil,
			goldenFile: "testdata/new-account/success",
		},
		"Failure While db error": {
			input:      singleAccount,
			err:        accounts.ErrDatabase,
			goldenFile: "testdata/new-account/failure-db-error",
		},
		"Failure While accountID not exits": {
			input:      singleAccountNoAccountID,
			err:        nil,
			goldenFile: "testdata/new-account/failure-accountid-not-exist",
		},
		"Failure While first name and last name not exits": {
			input:      singleAccountNoName,
			err:        nil,
			goldenFile: "testdata/new-account/failure-name-not-exist",
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
			s.On("PostAccount", test.input).
				Return(test.err).
				Maybe()

			svc := accounts.NewService(s, clock)
			svc = accounts.NewValidationMiddleware()(svc)

			endpoint := accounts.NewEndpoint(svc, log.NewNopLogger())
			handler := accounts.NewHTTPHandler(endpoint, log.NewNopLogger())

			handler.ServeHTTP(res, req)
			goldie.Assert(t, test.goldenFile, res.Body.Bytes())

			s.AssertExpectations(t)
		})
	}
}
