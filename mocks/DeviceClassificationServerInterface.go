// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// DeviceClassificationServerInterface is an autogenerated mock type for the DeviceClassificationServerInterface type
type DeviceClassificationServerInterface struct {
	mock.Mock
}

type DeviceClassificationServerInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *DeviceClassificationServerInterface) EXPECT() *DeviceClassificationServerInterface_Expecter {
	return &DeviceClassificationServerInterface_Expecter{mock: &_m.Mock}
}

// NewDeviceClassificationServerInterface creates a new instance of DeviceClassificationServerInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDeviceClassificationServerInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *DeviceClassificationServerInterface {
	mock := &DeviceClassificationServerInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}