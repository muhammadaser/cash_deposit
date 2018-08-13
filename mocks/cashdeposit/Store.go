// Code generated by mockery v1.0.0
package mocks

import cashdeposit "github.com/muhammadaser/cash_deposit/cashdeposit"
import mock "github.com/stretchr/testify/mock"

// Store is an autogenerated mock type for the Store type
type Store struct {
	mock.Mock
}

// GetListDeposits provides a mock function with given fields:
func (_m *Store) GetListDeposits() ([]cashdeposit.CashDeposit, error) {
	ret := _m.Called()

	var r0 []cashdeposit.CashDeposit
	if rf, ok := ret.Get(0).(func() []cashdeposit.CashDeposit); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]cashdeposit.CashDeposit)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}