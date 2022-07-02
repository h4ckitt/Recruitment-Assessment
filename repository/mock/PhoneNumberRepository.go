// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// PhoneNumberRepository is an autogenerated mock type for the PhoneNumberRepository type
type PhoneNumberRepository struct {
	mock.Mock
}

// FetchPaginatedPhoneNumbers provides a mock function with given fields: offset, limit
func (_m *PhoneNumberRepository) FetchPaginatedPhoneNumbers(offset int, limit int) ([]string, error) {
	ret := _m.Called(offset, limit)

	var r0 []string
	if rf, ok := ret.Get(0).(func(int, int) []string); ok {
		r0 = rf(offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FetchPaginatedPhoneNumbersByCode provides a mock function with given fields: code, offset, limit
func (_m *PhoneNumberRepository) FetchPaginatedPhoneNumbersByCode(code string, offset int, limit int) ([]string, error) {
	ret := _m.Called(code, offset, limit)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string, int, int) []string); ok {
		r0 = rf(code, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int, int) error); ok {
		r1 = rf(code, offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewPhoneNumberRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewPhoneNumberRepository creates a new instance of PhoneNumberRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPhoneNumberRepository(t mockConstructorTestingTNewPhoneNumberRepository) *PhoneNumberRepository {
	mock := &PhoneNumberRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
