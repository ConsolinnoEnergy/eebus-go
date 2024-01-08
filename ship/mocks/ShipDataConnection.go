// Code generated by mockery v2.39.1. DO NOT EDIT.

package mocks

import (
	ship "github.com/enbility/eebus-go/ship"
	mock "github.com/stretchr/testify/mock"
)

// ShipDataConnection is an autogenerated mock type for the ShipDataConnection type
type ShipDataConnection struct {
	mock.Mock
}

// CloseDataConnection provides a mock function with given fields: closeCode, reason
func (_m *ShipDataConnection) CloseDataConnection(closeCode int, reason string) {
	_m.Called(closeCode, reason)
}

// InitDataProcessing provides a mock function with given fields: _a0
func (_m *ShipDataConnection) InitDataProcessing(_a0 ship.ShipDataProcessing) {
	_m.Called(_a0)
}

// IsDataConnectionClosed provides a mock function with given fields:
func (_m *ShipDataConnection) IsDataConnectionClosed() (bool, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsDataConnectionClosed")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func() (bool, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WriteMessageToDataConnection provides a mock function with given fields: _a0
func (_m *ShipDataConnection) WriteMessageToDataConnection(_a0 []byte) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for WriteMessageToDataConnection")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewShipDataConnection creates a new instance of ShipDataConnection. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewShipDataConnection(t interface {
	mock.TestingT
	Cleanup(func())
}) *ShipDataConnection {
	mock := &ShipDataConnection{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}