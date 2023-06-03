// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/enbility/eebus-go/service (interfaces: ServiceProvider)

// Package service is a generated GoMock package.
package service

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockServiceProvider is a mock of ServiceProvider interface.
type MockServiceProvider struct {
	ctrl     *gomock.Controller
	recorder *MockServiceProviderMockRecorder
}

// MockServiceProviderMockRecorder is the mock recorder for MockServiceProvider.
type MockServiceProviderMockRecorder struct {
	mock *MockServiceProvider
}

// NewMockServiceProvider creates a new mock instance.
func NewMockServiceProvider(ctrl *gomock.Controller) *MockServiceProvider {
	mock := &MockServiceProvider{ctrl: ctrl}
	mock.recorder = &MockServiceProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceProvider) EXPECT() *MockServiceProviderMockRecorder {
	return m.recorder
}

// AllowWaitingForTrust mocks base method.
func (m *MockServiceProvider) AllowWaitingForTrust(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllowWaitingForTrust", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// AllowWaitingForTrust indicates an expected call of AllowWaitingForTrust.
func (mr *MockServiceProviderMockRecorder) AllowWaitingForTrust(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllowWaitingForTrust", reflect.TypeOf((*MockServiceProvider)(nil).AllowWaitingForTrust), arg0)
}

// RemoteSKIConnected mocks base method.
func (m *MockServiceProvider) RemoteSKIConnected(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoteSKIConnected", arg0)
}

// RemoteSKIConnected indicates an expected call of RemoteSKIConnected.
func (mr *MockServiceProviderMockRecorder) RemoteSKIConnected(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoteSKIConnected", reflect.TypeOf((*MockServiceProvider)(nil).RemoteSKIConnected), arg0)
}

// RemoteSKIDisconnected mocks base method.
func (m *MockServiceProvider) RemoteSKIDisconnected(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoteSKIDisconnected", arg0)
}

// RemoteSKIDisconnected indicates an expected call of RemoteSKIDisconnected.
func (mr *MockServiceProviderMockRecorder) RemoteSKIDisconnected(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoteSKIDisconnected", reflect.TypeOf((*MockServiceProvider)(nil).RemoteSKIDisconnected), arg0)
}

// ServicePairingDetailUpdate mocks base method.
func (m *MockServiceProvider) ServicePairingDetailUpdate(arg0 string, arg1 ConnectionStateDetail) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ServicePairingDetailUpdate", arg0, arg1)
}

// ServicePairingDetailUpdate indicates an expected call of ServicePairingDetailUpdate.
func (mr *MockServiceProviderMockRecorder) ServicePairingDetailUpdate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServicePairingDetailUpdate", reflect.TypeOf((*MockServiceProvider)(nil).ServicePairingDetailUpdate), arg0, arg1)
}

// ServiceShipIDUpdate mocks base method.
func (m *MockServiceProvider) ServiceShipIDUpdate(arg0, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ServiceShipIDUpdate", arg0, arg1)
}

// ServiceShipIDUpdate indicates an expected call of ServiceShipIDUpdate.
func (mr *MockServiceProviderMockRecorder) ServiceShipIDUpdate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServiceShipIDUpdate", reflect.TypeOf((*MockServiceProvider)(nil).ServiceShipIDUpdate), arg0, arg1)
}

// VisibleMDNSRecordsUpdated mocks base method.
func (m *MockServiceProvider) VisibleMDNSRecordsUpdated(arg0 []MdnsEntry) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "VisibleMDNSRecordsUpdated", arg0)
}

// VisibleMDNSRecordsUpdated indicates an expected call of VisibleMDNSRecordsUpdated.
func (mr *MockServiceProviderMockRecorder) VisibleMDNSRecordsUpdated(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VisibleMDNSRecordsUpdated", reflect.TypeOf((*MockServiceProvider)(nil).VisibleMDNSRecordsUpdated), arg0)
}
