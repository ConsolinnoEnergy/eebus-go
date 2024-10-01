// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	api "github.com/enbility/eebus-go/api"
	mock "github.com/stretchr/testify/mock"

	model "github.com/enbility/spine-go/model"
)

// MeasurementServerInterface is an autogenerated mock type for the MeasurementServerInterface type
type MeasurementServerInterface struct {
	mock.Mock
}

type MeasurementServerInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *MeasurementServerInterface) EXPECT() *MeasurementServerInterface_Expecter {
	return &MeasurementServerInterface_Expecter{mock: &_m.Mock}
}

// AddDescription provides a mock function with given fields: description
func (_m *MeasurementServerInterface) AddDescription(description model.MeasurementDescriptionDataType) *model.MeasurementIdType {
	ret := _m.Called(description)

	if len(ret) == 0 {
		panic("no return value specified for AddDescription")
	}

	var r0 *model.MeasurementIdType
	if rf, ok := ret.Get(0).(func(model.MeasurementDescriptionDataType) *model.MeasurementIdType); ok {
		r0 = rf(description)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.MeasurementIdType)
		}
	}

	return r0
}

// MeasurementServerInterface_AddDescription_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddDescription'
type MeasurementServerInterface_AddDescription_Call struct {
	*mock.Call
}

// AddDescription is a helper method to define mock.On call
//   - description model.MeasurementDescriptionDataType
func (_e *MeasurementServerInterface_Expecter) AddDescription(description interface{}) *MeasurementServerInterface_AddDescription_Call {
	return &MeasurementServerInterface_AddDescription_Call{Call: _e.mock.On("AddDescription", description)}
}

func (_c *MeasurementServerInterface_AddDescription_Call) Run(run func(description model.MeasurementDescriptionDataType)) *MeasurementServerInterface_AddDescription_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(model.MeasurementDescriptionDataType))
	})
	return _c
}

func (_c *MeasurementServerInterface_AddDescription_Call) Return(_a0 *model.MeasurementIdType) *MeasurementServerInterface_AddDescription_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MeasurementServerInterface_AddDescription_Call) RunAndReturn(run func(model.MeasurementDescriptionDataType) *model.MeasurementIdType) *MeasurementServerInterface_AddDescription_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateDataForFilters provides a mock function with given fields: data, deleteSelector, deleteElements
func (_m *MeasurementServerInterface) UpdateDataForFilters(data []api.MeasurementDataForFilter, deleteSelector *model.MeasurementListDataSelectorsType, deleteElements *model.MeasurementDataElementsType) error {
	ret := _m.Called(data, deleteSelector, deleteElements)

	if len(ret) == 0 {
		panic("no return value specified for UpdateDataForFilters")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]api.MeasurementDataForFilter, *model.MeasurementListDataSelectorsType, *model.MeasurementDataElementsType) error); ok {
		r0 = rf(data, deleteSelector, deleteElements)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MeasurementServerInterface_UpdateDataForFilters_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateDataForFilters'
type MeasurementServerInterface_UpdateDataForFilters_Call struct {
	*mock.Call
}

// UpdateDataForFilters is a helper method to define mock.On call
//   - data []api.MeasurementDataForFilter
//   - deleteSelector *model.MeasurementListDataSelectorsType
//   - deleteElements *model.MeasurementDataElementsType
func (_e *MeasurementServerInterface_Expecter) UpdateDataForFilters(data interface{}, deleteSelector interface{}, deleteElements interface{}) *MeasurementServerInterface_UpdateDataForFilters_Call {
	return &MeasurementServerInterface_UpdateDataForFilters_Call{Call: _e.mock.On("UpdateDataForFilters", data, deleteSelector, deleteElements)}
}

func (_c *MeasurementServerInterface_UpdateDataForFilters_Call) Run(run func(data []api.MeasurementDataForFilter, deleteSelector *model.MeasurementListDataSelectorsType, deleteElements *model.MeasurementDataElementsType)) *MeasurementServerInterface_UpdateDataForFilters_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]api.MeasurementDataForFilter), args[1].(*model.MeasurementListDataSelectorsType), args[2].(*model.MeasurementDataElementsType))
	})
	return _c
}

func (_c *MeasurementServerInterface_UpdateDataForFilters_Call) Return(_a0 error) *MeasurementServerInterface_UpdateDataForFilters_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MeasurementServerInterface_UpdateDataForFilters_Call) RunAndReturn(run func([]api.MeasurementDataForFilter, *model.MeasurementListDataSelectorsType, *model.MeasurementDataElementsType) error) *MeasurementServerInterface_UpdateDataForFilters_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateDataForIds provides a mock function with given fields: data, deleteId, deleteElements
func (_m *MeasurementServerInterface) UpdateDataForIds(data []api.MeasurementDataForID, deleteId *model.MeasurementIdType, deleteElements *model.MeasurementDataElementsType) error {
	ret := _m.Called(data, deleteId, deleteElements)

	if len(ret) == 0 {
		panic("no return value specified for UpdateDataForIds")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]api.MeasurementDataForID, *model.MeasurementIdType, *model.MeasurementDataElementsType) error); ok {
		r0 = rf(data, deleteId, deleteElements)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MeasurementServerInterface_UpdateDataForIds_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateDataForIds'
type MeasurementServerInterface_UpdateDataForIds_Call struct {
	*mock.Call
}

// UpdateDataForIds is a helper method to define mock.On call
//   - data []api.MeasurementDataForID
//   - deleteId *model.MeasurementIdType
//   - deleteElements *model.MeasurementDataElementsType
func (_e *MeasurementServerInterface_Expecter) UpdateDataForIds(data interface{}, deleteId interface{}, deleteElements interface{}) *MeasurementServerInterface_UpdateDataForIds_Call {
	return &MeasurementServerInterface_UpdateDataForIds_Call{Call: _e.mock.On("UpdateDataForIds", data, deleteId, deleteElements)}
}

func (_c *MeasurementServerInterface_UpdateDataForIds_Call) Run(run func(data []api.MeasurementDataForID, deleteId *model.MeasurementIdType, deleteElements *model.MeasurementDataElementsType)) *MeasurementServerInterface_UpdateDataForIds_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]api.MeasurementDataForID), args[1].(*model.MeasurementIdType), args[2].(*model.MeasurementDataElementsType))
	})
	return _c
}

func (_c *MeasurementServerInterface_UpdateDataForIds_Call) Return(_a0 error) *MeasurementServerInterface_UpdateDataForIds_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MeasurementServerInterface_UpdateDataForIds_Call) RunAndReturn(run func([]api.MeasurementDataForID, *model.MeasurementIdType, *model.MeasurementDataElementsType) error) *MeasurementServerInterface_UpdateDataForIds_Call {
	_c.Call.Return(run)
	return _c
}

// NewMeasurementServerInterface creates a new instance of MeasurementServerInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMeasurementServerInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *MeasurementServerInterface {
	mock := &MeasurementServerInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
