// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mi4r/gophermart/internal/storage (interfaces: StorageAccrualSystem)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockStorageAccrualSystem is a mock of StorageAccrualSystem interface.
type MockStorageAccrualSystem struct {
	ctrl     *gomock.Controller
	recorder *MockStorageAccrualSystemMockRecorder
}

// MockStorageAccrualSystemMockRecorder is the mock recorder for MockStorageAccrualSystem.
type MockStorageAccrualSystemMockRecorder struct {
	mock *MockStorageAccrualSystem
}

// NewMockStorageAccrualSystem creates a new mock instance.
func NewMockStorageAccrualSystem(ctrl *gomock.Controller) *MockStorageAccrualSystem {
	mock := &MockStorageAccrualSystem{ctrl: ctrl}
	mock.recorder = &MockStorageAccrualSystemMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorageAccrualSystem) EXPECT() *MockStorageAccrualSystemMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockStorageAccrualSystem) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockStorageAccrualSystemMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStorageAccrualSystem)(nil).Close))
}

// Migrate mocks base method.
func (m *MockStorageAccrualSystem) Migrate(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Migrate", arg0)
}

// Migrate indicates an expected call of Migrate.
func (mr *MockStorageAccrualSystemMockRecorder) Migrate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Migrate", reflect.TypeOf((*MockStorageAccrualSystem)(nil).Migrate), arg0)
}

// Open mocks base method.
func (m *MockStorageAccrualSystem) Open() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Open")
	ret0, _ := ret[0].(error)
	return ret0
}

// Open indicates an expected call of Open.
func (mr *MockStorageAccrualSystemMockRecorder) Open() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Open", reflect.TypeOf((*MockStorageAccrualSystem)(nil).Open))
}

// Ping mocks base method.
func (m *MockStorageAccrualSystem) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockStorageAccrualSystemMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockStorageAccrualSystem)(nil).Ping))
}
