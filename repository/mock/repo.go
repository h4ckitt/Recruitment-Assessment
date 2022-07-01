// Code generated by MockGen. DO NOT EDIT.
// Source: repository/repo_interface.go

// Package repo_mock is a generated GoMock package.
package repo_mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPhoneNumberRepository is a mock of PhoneNumberRepository interface.
type MockPhoneNumberRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPhoneNumberRepositoryMockRecorder
}

// MockPhoneNumberRepositoryMockRecorder is the mock recorder for MockPhoneNumberRepository.
type MockPhoneNumberRepositoryMockRecorder struct {
	mock *MockPhoneNumberRepository
}

// NewMockPhoneNumberRepository creates a new mock instance.
func NewMockPhoneNumberRepository(ctrl *gomock.Controller) *MockPhoneNumberRepository {
	mock := &MockPhoneNumberRepository{ctrl: ctrl}
	mock.recorder = &MockPhoneNumberRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPhoneNumberRepository) EXPECT() *MockPhoneNumberRepositoryMockRecorder {
	return m.recorder
}

// FetchPaginatedPhoneNumbers mocks base method.
func (m *MockPhoneNumberRepository) FetchPaginatedPhoneNumbers(offset, limit int) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchPaginatedPhoneNumbers", offset, limit)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchPaginatedPhoneNumbers indicates an expected call of FetchPaginatedPhoneNumbers.
func (mr *MockPhoneNumberRepositoryMockRecorder) FetchPaginatedPhoneNumbers(offset, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchPaginatedPhoneNumbers", reflect.TypeOf((*MockPhoneNumberRepository)(nil).FetchPaginatedPhoneNumbers), offset, limit)
}

// FetchPaginatedPhoneNumbersByCode mocks base method.
func (m *MockPhoneNumberRepository) FetchPaginatedPhoneNumbersByCode(code string, offset, limit int) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchPaginatedPhoneNumbersByCode", code, offset, limit)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchPaginatedPhoneNumbersByCode indicates an expected call of FetchPaginatedPhoneNumbersByCode.
func (mr *MockPhoneNumberRepositoryMockRecorder) FetchPaginatedPhoneNumbersByCode(code, offset, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchPaginatedPhoneNumbersByCode", reflect.TypeOf((*MockPhoneNumberRepository)(nil).FetchPaginatedPhoneNumbersByCode), code, offset, limit)
}