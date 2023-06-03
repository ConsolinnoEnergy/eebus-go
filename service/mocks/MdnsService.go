// Code generated by mockery v2.26.1. DO NOT EDIT.

package mocks

import (
	service "github.com/enbility/eebus-go/service"
	mock "github.com/stretchr/testify/mock"
)

// MdnsService is an autogenerated mock type for the MdnsService type
type MdnsService struct {
	mock.Mock
}

// AnnounceMdnsEntry provides a mock function with given fields:
func (_m *MdnsService) AnnounceMdnsEntry() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RegisterMdnsSearch provides a mock function with given fields: cb
func (_m *MdnsService) RegisterMdnsSearch(cb service.MdnsSearch) {
	_m.Called(cb)
}

// SetupMdnsService provides a mock function with given fields:
func (_m *MdnsService) SetupMdnsService() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ShutdownMdnsService provides a mock function with given fields:
func (_m *MdnsService) ShutdownMdnsService() {
	_m.Called()
}

// UnannounceMdnsEntry provides a mock function with given fields:
func (_m *MdnsService) UnannounceMdnsEntry() {
	_m.Called()
}

// UnregisterMdnsSearch provides a mock function with given fields: cb
func (_m *MdnsService) UnregisterMdnsSearch(cb service.MdnsSearch) {
	_m.Called(cb)
}

type mockConstructorTestingTNewMdnsService interface {
	mock.TestingT
	Cleanup(func())
}

// NewMdnsService creates a new instance of MdnsService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMdnsService(t mockConstructorTestingTNewMdnsService) *MdnsService {
	mock := &MdnsService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
