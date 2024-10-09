// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	eebus_goapi "github.com/enbility/eebus-go/api"
	api "github.com/enbility/eebus-go/usecases/api"

	mock "github.com/stretchr/testify/mock"

	model "github.com/enbility/spine-go/model"

	spine_goapi "github.com/enbility/spine-go/api"
)

// CemOPEVInterface is an autogenerated mock type for the CemOPEVInterface type
type CemOPEVInterface struct {
	mock.Mock
}

type CemOPEVInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *CemOPEVInterface) EXPECT() *CemOPEVInterface_Expecter {
	return &CemOPEVInterface_Expecter{mock: &_m.Mock}
}

// AddFeatures provides a mock function with given fields:
func (_m *CemOPEVInterface) AddFeatures() {
	_m.Called()
}

// CemOPEVInterface_AddFeatures_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddFeatures'
type CemOPEVInterface_AddFeatures_Call struct {
	*mock.Call
}

// AddFeatures is a helper method to define mock.On call
func (_e *CemOPEVInterface_Expecter) AddFeatures() *CemOPEVInterface_AddFeatures_Call {
	return &CemOPEVInterface_AddFeatures_Call{Call: _e.mock.On("AddFeatures")}
}

func (_c *CemOPEVInterface_AddFeatures_Call) Run(run func()) *CemOPEVInterface_AddFeatures_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *CemOPEVInterface_AddFeatures_Call) Return() *CemOPEVInterface_AddFeatures_Call {
	_c.Call.Return()
	return _c
}

func (_c *CemOPEVInterface_AddFeatures_Call) RunAndReturn(run func()) *CemOPEVInterface_AddFeatures_Call {
	_c.Call.Return(run)
	return _c
}

// AddUseCase provides a mock function with given fields:
func (_m *CemOPEVInterface) AddUseCase() {
	_m.Called()
}

// CemOPEVInterface_AddUseCase_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddUseCase'
type CemOPEVInterface_AddUseCase_Call struct {
	*mock.Call
}

// AddUseCase is a helper method to define mock.On call
func (_e *CemOPEVInterface_Expecter) AddUseCase() *CemOPEVInterface_AddUseCase_Call {
	return &CemOPEVInterface_AddUseCase_Call{Call: _e.mock.On("AddUseCase")}
}

func (_c *CemOPEVInterface_AddUseCase_Call) Run(run func()) *CemOPEVInterface_AddUseCase_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *CemOPEVInterface_AddUseCase_Call) Return() *CemOPEVInterface_AddUseCase_Call {
	_c.Call.Return()
	return _c
}

func (_c *CemOPEVInterface_AddUseCase_Call) RunAndReturn(run func()) *CemOPEVInterface_AddUseCase_Call {
	_c.Call.Return(run)
	return _c
}

// AvailableScenariosForEntity provides a mock function with given fields: entity
func (_m *CemOPEVInterface) AvailableScenariosForEntity(entity spine_goapi.EntityRemoteInterface) []uint {
	ret := _m.Called(entity)

	if len(ret) == 0 {
		panic("no return value specified for AvailableScenariosForEntity")
	}

	var r0 []uint
	if rf, ok := ret.Get(0).(func(spine_goapi.EntityRemoteInterface) []uint); ok {
		r0 = rf(entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]uint)
		}
	}

	return r0
}

// CemOPEVInterface_AvailableScenariosForEntity_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AvailableScenariosForEntity'
type CemOPEVInterface_AvailableScenariosForEntity_Call struct {
	*mock.Call
}

// AvailableScenariosForEntity is a helper method to define mock.On call
//   - entity spine_goapi.EntityRemoteInterface
func (_e *CemOPEVInterface_Expecter) AvailableScenariosForEntity(entity interface{}) *CemOPEVInterface_AvailableScenariosForEntity_Call {
	return &CemOPEVInterface_AvailableScenariosForEntity_Call{Call: _e.mock.On("AvailableScenariosForEntity", entity)}
}

func (_c *CemOPEVInterface_AvailableScenariosForEntity_Call) Run(run func(entity spine_goapi.EntityRemoteInterface)) *CemOPEVInterface_AvailableScenariosForEntity_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(spine_goapi.EntityRemoteInterface))
	})
	return _c
}

func (_c *CemOPEVInterface_AvailableScenariosForEntity_Call) Return(_a0 []uint) *CemOPEVInterface_AvailableScenariosForEntity_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CemOPEVInterface_AvailableScenariosForEntity_Call) RunAndReturn(run func(spine_goapi.EntityRemoteInterface) []uint) *CemOPEVInterface_AvailableScenariosForEntity_Call {
	_c.Call.Return(run)
	return _c
}

// CurrentLimits provides a mock function with given fields: entity
func (_m *CemOPEVInterface) CurrentLimits(entity spine_goapi.EntityRemoteInterface) ([]float64, []float64, []float64, error) {
	ret := _m.Called(entity)

	if len(ret) == 0 {
		panic("no return value specified for CurrentLimits")
	}

	var r0 []float64
	var r1 []float64
	var r2 []float64
	var r3 error
	if rf, ok := ret.Get(0).(func(spine_goapi.EntityRemoteInterface) ([]float64, []float64, []float64, error)); ok {
		return rf(entity)
	}
	if rf, ok := ret.Get(0).(func(spine_goapi.EntityRemoteInterface) []float64); ok {
		r0 = rf(entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]float64)
		}
	}

	if rf, ok := ret.Get(1).(func(spine_goapi.EntityRemoteInterface) []float64); ok {
		r1 = rf(entity)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]float64)
		}
	}

	if rf, ok := ret.Get(2).(func(spine_goapi.EntityRemoteInterface) []float64); ok {
		r2 = rf(entity)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).([]float64)
		}
	}

	if rf, ok := ret.Get(3).(func(spine_goapi.EntityRemoteInterface) error); ok {
		r3 = rf(entity)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// CemOPEVInterface_CurrentLimits_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CurrentLimits'
type CemOPEVInterface_CurrentLimits_Call struct {
	*mock.Call
}

// CurrentLimits is a helper method to define mock.On call
//   - entity spine_goapi.EntityRemoteInterface
func (_e *CemOPEVInterface_Expecter) CurrentLimits(entity interface{}) *CemOPEVInterface_CurrentLimits_Call {
	return &CemOPEVInterface_CurrentLimits_Call{Call: _e.mock.On("CurrentLimits", entity)}
}

func (_c *CemOPEVInterface_CurrentLimits_Call) Run(run func(entity spine_goapi.EntityRemoteInterface)) *CemOPEVInterface_CurrentLimits_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(spine_goapi.EntityRemoteInterface))
	})
	return _c
}

func (_c *CemOPEVInterface_CurrentLimits_Call) Return(_a0 []float64, _a1 []float64, _a2 []float64, _a3 error) *CemOPEVInterface_CurrentLimits_Call {
	_c.Call.Return(_a0, _a1, _a2, _a3)
	return _c
}

func (_c *CemOPEVInterface_CurrentLimits_Call) RunAndReturn(run func(spine_goapi.EntityRemoteInterface) ([]float64, []float64, []float64, error)) *CemOPEVInterface_CurrentLimits_Call {
	_c.Call.Return(run)
	return _c
}

// IsCompatibleEntityType provides a mock function with given fields: entity
func (_m *CemOPEVInterface) IsCompatibleEntityType(entity spine_goapi.EntityRemoteInterface) bool {
	ret := _m.Called(entity)

	if len(ret) == 0 {
		panic("no return value specified for IsCompatibleEntityType")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(spine_goapi.EntityRemoteInterface) bool); ok {
		r0 = rf(entity)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// CemOPEVInterface_IsCompatibleEntityType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsCompatibleEntityType'
type CemOPEVInterface_IsCompatibleEntityType_Call struct {
	*mock.Call
}

// IsCompatibleEntityType is a helper method to define mock.On call
//   - entity spine_goapi.EntityRemoteInterface
func (_e *CemOPEVInterface_Expecter) IsCompatibleEntityType(entity interface{}) *CemOPEVInterface_IsCompatibleEntityType_Call {
	return &CemOPEVInterface_IsCompatibleEntityType_Call{Call: _e.mock.On("IsCompatibleEntityType", entity)}
}

func (_c *CemOPEVInterface_IsCompatibleEntityType_Call) Run(run func(entity spine_goapi.EntityRemoteInterface)) *CemOPEVInterface_IsCompatibleEntityType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(spine_goapi.EntityRemoteInterface))
	})
	return _c
}

func (_c *CemOPEVInterface_IsCompatibleEntityType_Call) Return(_a0 bool) *CemOPEVInterface_IsCompatibleEntityType_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CemOPEVInterface_IsCompatibleEntityType_Call) RunAndReturn(run func(spine_goapi.EntityRemoteInterface) bool) *CemOPEVInterface_IsCompatibleEntityType_Call {
	_c.Call.Return(run)
	return _c
}

// IsScenarioAvailableAtEntity provides a mock function with given fields: entity, scenario
func (_m *CemOPEVInterface) IsScenarioAvailableAtEntity(entity spine_goapi.EntityRemoteInterface, scenario uint) bool {
	ret := _m.Called(entity, scenario)

	if len(ret) == 0 {
		panic("no return value specified for IsScenarioAvailableAtEntity")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(spine_goapi.EntityRemoteInterface, uint) bool); ok {
		r0 = rf(entity, scenario)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// CemOPEVInterface_IsScenarioAvailableAtEntity_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsScenarioAvailableAtEntity'
type CemOPEVInterface_IsScenarioAvailableAtEntity_Call struct {
	*mock.Call
}

// IsScenarioAvailableAtEntity is a helper method to define mock.On call
//   - entity spine_goapi.EntityRemoteInterface
//   - scenario uint
func (_e *CemOPEVInterface_Expecter) IsScenarioAvailableAtEntity(entity interface{}, scenario interface{}) *CemOPEVInterface_IsScenarioAvailableAtEntity_Call {
	return &CemOPEVInterface_IsScenarioAvailableAtEntity_Call{Call: _e.mock.On("IsScenarioAvailableAtEntity", entity, scenario)}
}

func (_c *CemOPEVInterface_IsScenarioAvailableAtEntity_Call) Run(run func(entity spine_goapi.EntityRemoteInterface, scenario uint)) *CemOPEVInterface_IsScenarioAvailableAtEntity_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(spine_goapi.EntityRemoteInterface), args[1].(uint))
	})
	return _c
}

func (_c *CemOPEVInterface_IsScenarioAvailableAtEntity_Call) Return(_a0 bool) *CemOPEVInterface_IsScenarioAvailableAtEntity_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CemOPEVInterface_IsScenarioAvailableAtEntity_Call) RunAndReturn(run func(spine_goapi.EntityRemoteInterface, uint) bool) *CemOPEVInterface_IsScenarioAvailableAtEntity_Call {
	_c.Call.Return(run)
	return _c
}

// LoadControlLimits provides a mock function with given fields: entity
func (_m *CemOPEVInterface) LoadControlLimits(entity spine_goapi.EntityRemoteInterface) ([]api.LoadLimitsPhase, error) {
	ret := _m.Called(entity)

	if len(ret) == 0 {
		panic("no return value specified for LoadControlLimits")
	}

	var r0 []api.LoadLimitsPhase
	var r1 error
	if rf, ok := ret.Get(0).(func(spine_goapi.EntityRemoteInterface) ([]api.LoadLimitsPhase, error)); ok {
		return rf(entity)
	}
	if rf, ok := ret.Get(0).(func(spine_goapi.EntityRemoteInterface) []api.LoadLimitsPhase); ok {
		r0 = rf(entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]api.LoadLimitsPhase)
		}
	}

	if rf, ok := ret.Get(1).(func(spine_goapi.EntityRemoteInterface) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CemOPEVInterface_LoadControlLimits_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LoadControlLimits'
type CemOPEVInterface_LoadControlLimits_Call struct {
	*mock.Call
}

// LoadControlLimits is a helper method to define mock.On call
//   - entity spine_goapi.EntityRemoteInterface
func (_e *CemOPEVInterface_Expecter) LoadControlLimits(entity interface{}) *CemOPEVInterface_LoadControlLimits_Call {
	return &CemOPEVInterface_LoadControlLimits_Call{Call: _e.mock.On("LoadControlLimits", entity)}
}

func (_c *CemOPEVInterface_LoadControlLimits_Call) Run(run func(entity spine_goapi.EntityRemoteInterface)) *CemOPEVInterface_LoadControlLimits_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(spine_goapi.EntityRemoteInterface))
	})
	return _c
}

func (_c *CemOPEVInterface_LoadControlLimits_Call) Return(limits []api.LoadLimitsPhase, resultErr error) *CemOPEVInterface_LoadControlLimits_Call {
	_c.Call.Return(limits, resultErr)
	return _c
}

func (_c *CemOPEVInterface_LoadControlLimits_Call) RunAndReturn(run func(spine_goapi.EntityRemoteInterface) ([]api.LoadLimitsPhase, error)) *CemOPEVInterface_LoadControlLimits_Call {
	_c.Call.Return(run)
	return _c
}

// RemoteEntitiesScenarios provides a mock function with given fields:
func (_m *CemOPEVInterface) RemoteEntitiesScenarios() []eebus_goapi.RemoteEntityScenarios {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for RemoteEntitiesScenarios")
	}

	var r0 []eebus_goapi.RemoteEntityScenarios
	if rf, ok := ret.Get(0).(func() []eebus_goapi.RemoteEntityScenarios); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]eebus_goapi.RemoteEntityScenarios)
		}
	}

	return r0
}

// CemOPEVInterface_RemoteEntitiesScenarios_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoteEntitiesScenarios'
type CemOPEVInterface_RemoteEntitiesScenarios_Call struct {
	*mock.Call
}

// RemoteEntitiesScenarios is a helper method to define mock.On call
func (_e *CemOPEVInterface_Expecter) RemoteEntitiesScenarios() *CemOPEVInterface_RemoteEntitiesScenarios_Call {
	return &CemOPEVInterface_RemoteEntitiesScenarios_Call{Call: _e.mock.On("RemoteEntitiesScenarios")}
}

func (_c *CemOPEVInterface_RemoteEntitiesScenarios_Call) Run(run func()) *CemOPEVInterface_RemoteEntitiesScenarios_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *CemOPEVInterface_RemoteEntitiesScenarios_Call) Return(_a0 []eebus_goapi.RemoteEntityScenarios) *CemOPEVInterface_RemoteEntitiesScenarios_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CemOPEVInterface_RemoteEntitiesScenarios_Call) RunAndReturn(run func() []eebus_goapi.RemoteEntityScenarios) *CemOPEVInterface_RemoteEntitiesScenarios_Call {
	_c.Call.Return(run)
	return _c
}

// RemoveUseCase provides a mock function with given fields:
func (_m *CemOPEVInterface) RemoveUseCase() {
	_m.Called()
}

// CemOPEVInterface_RemoveUseCase_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoveUseCase'
type CemOPEVInterface_RemoveUseCase_Call struct {
	*mock.Call
}

// RemoveUseCase is a helper method to define mock.On call
func (_e *CemOPEVInterface_Expecter) RemoveUseCase() *CemOPEVInterface_RemoveUseCase_Call {
	return &CemOPEVInterface_RemoveUseCase_Call{Call: _e.mock.On("RemoveUseCase")}
}

func (_c *CemOPEVInterface_RemoveUseCase_Call) Run(run func()) *CemOPEVInterface_RemoveUseCase_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *CemOPEVInterface_RemoveUseCase_Call) Return() *CemOPEVInterface_RemoveUseCase_Call {
	_c.Call.Return()
	return _c
}

func (_c *CemOPEVInterface_RemoveUseCase_Call) RunAndReturn(run func()) *CemOPEVInterface_RemoveUseCase_Call {
	_c.Call.Return(run)
	return _c
}

// SetOperatingState provides a mock function with given fields: failureState
func (_m *CemOPEVInterface) SetOperatingState(failureState bool) error {
	ret := _m.Called(failureState)

	if len(ret) == 0 {
		panic("no return value specified for SetOperatingState")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(bool) error); ok {
		r0 = rf(failureState)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CemOPEVInterface_SetOperatingState_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetOperatingState'
type CemOPEVInterface_SetOperatingState_Call struct {
	*mock.Call
}

// SetOperatingState is a helper method to define mock.On call
//   - failureState bool
func (_e *CemOPEVInterface_Expecter) SetOperatingState(failureState interface{}) *CemOPEVInterface_SetOperatingState_Call {
	return &CemOPEVInterface_SetOperatingState_Call{Call: _e.mock.On("SetOperatingState", failureState)}
}

func (_c *CemOPEVInterface_SetOperatingState_Call) Run(run func(failureState bool)) *CemOPEVInterface_SetOperatingState_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool))
	})
	return _c
}

func (_c *CemOPEVInterface_SetOperatingState_Call) Return(_a0 error) *CemOPEVInterface_SetOperatingState_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CemOPEVInterface_SetOperatingState_Call) RunAndReturn(run func(bool) error) *CemOPEVInterface_SetOperatingState_Call {
	_c.Call.Return(run)
	return _c
}

// StartHeartbeat provides a mock function with given fields:
func (_m *CemOPEVInterface) StartHeartbeat() {
	_m.Called()
}

// CemOPEVInterface_StartHeartbeat_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StartHeartbeat'
type CemOPEVInterface_StartHeartbeat_Call struct {
	*mock.Call
}

// StartHeartbeat is a helper method to define mock.On call
func (_e *CemOPEVInterface_Expecter) StartHeartbeat() *CemOPEVInterface_StartHeartbeat_Call {
	return &CemOPEVInterface_StartHeartbeat_Call{Call: _e.mock.On("StartHeartbeat")}
}

func (_c *CemOPEVInterface_StartHeartbeat_Call) Run(run func()) *CemOPEVInterface_StartHeartbeat_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *CemOPEVInterface_StartHeartbeat_Call) Return() *CemOPEVInterface_StartHeartbeat_Call {
	_c.Call.Return()
	return _c
}

func (_c *CemOPEVInterface_StartHeartbeat_Call) RunAndReturn(run func()) *CemOPEVInterface_StartHeartbeat_Call {
	_c.Call.Return(run)
	return _c
}

// StopHeartbeat provides a mock function with given fields:
func (_m *CemOPEVInterface) StopHeartbeat() {
	_m.Called()
}

// CemOPEVInterface_StopHeartbeat_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StopHeartbeat'
type CemOPEVInterface_StopHeartbeat_Call struct {
	*mock.Call
}

// StopHeartbeat is a helper method to define mock.On call
func (_e *CemOPEVInterface_Expecter) StopHeartbeat() *CemOPEVInterface_StopHeartbeat_Call {
	return &CemOPEVInterface_StopHeartbeat_Call{Call: _e.mock.On("StopHeartbeat")}
}

func (_c *CemOPEVInterface_StopHeartbeat_Call) Run(run func()) *CemOPEVInterface_StopHeartbeat_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *CemOPEVInterface_StopHeartbeat_Call) Return() *CemOPEVInterface_StopHeartbeat_Call {
	_c.Call.Return()
	return _c
}

func (_c *CemOPEVInterface_StopHeartbeat_Call) RunAndReturn(run func()) *CemOPEVInterface_StopHeartbeat_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateUseCaseAvailability provides a mock function with given fields: available
func (_m *CemOPEVInterface) UpdateUseCaseAvailability(available bool) {
	_m.Called(available)
}

// CemOPEVInterface_UpdateUseCaseAvailability_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateUseCaseAvailability'
type CemOPEVInterface_UpdateUseCaseAvailability_Call struct {
	*mock.Call
}

// UpdateUseCaseAvailability is a helper method to define mock.On call
//   - available bool
func (_e *CemOPEVInterface_Expecter) UpdateUseCaseAvailability(available interface{}) *CemOPEVInterface_UpdateUseCaseAvailability_Call {
	return &CemOPEVInterface_UpdateUseCaseAvailability_Call{Call: _e.mock.On("UpdateUseCaseAvailability", available)}
}

func (_c *CemOPEVInterface_UpdateUseCaseAvailability_Call) Run(run func(available bool)) *CemOPEVInterface_UpdateUseCaseAvailability_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool))
	})
	return _c
}

func (_c *CemOPEVInterface_UpdateUseCaseAvailability_Call) Return() *CemOPEVInterface_UpdateUseCaseAvailability_Call {
	_c.Call.Return()
	return _c
}

func (_c *CemOPEVInterface_UpdateUseCaseAvailability_Call) RunAndReturn(run func(bool)) *CemOPEVInterface_UpdateUseCaseAvailability_Call {
	_c.Call.Return(run)
	return _c
}

// WriteLoadControlLimits provides a mock function with given fields: entity, limits, resultCB
func (_m *CemOPEVInterface) WriteLoadControlLimits(entity spine_goapi.EntityRemoteInterface, limits []api.LoadLimitsPhase, resultCB func(model.ResultDataType)) (*model.MsgCounterType, error) {
	ret := _m.Called(entity, limits, resultCB)

	if len(ret) == 0 {
		panic("no return value specified for WriteLoadControlLimits")
	}

	var r0 *model.MsgCounterType
	var r1 error
	if rf, ok := ret.Get(0).(func(spine_goapi.EntityRemoteInterface, []api.LoadLimitsPhase, func(model.ResultDataType)) (*model.MsgCounterType, error)); ok {
		return rf(entity, limits, resultCB)
	}
	if rf, ok := ret.Get(0).(func(spine_goapi.EntityRemoteInterface, []api.LoadLimitsPhase, func(model.ResultDataType)) *model.MsgCounterType); ok {
		r0 = rf(entity, limits, resultCB)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.MsgCounterType)
		}
	}

	if rf, ok := ret.Get(1).(func(spine_goapi.EntityRemoteInterface, []api.LoadLimitsPhase, func(model.ResultDataType)) error); ok {
		r1 = rf(entity, limits, resultCB)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CemOPEVInterface_WriteLoadControlLimits_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WriteLoadControlLimits'
type CemOPEVInterface_WriteLoadControlLimits_Call struct {
	*mock.Call
}

// WriteLoadControlLimits is a helper method to define mock.On call
//   - entity spine_goapi.EntityRemoteInterface
//   - limits []api.LoadLimitsPhase
//   - resultCB func(model.ResultDataType)
func (_e *CemOPEVInterface_Expecter) WriteLoadControlLimits(entity interface{}, limits interface{}, resultCB interface{}) *CemOPEVInterface_WriteLoadControlLimits_Call {
	return &CemOPEVInterface_WriteLoadControlLimits_Call{Call: _e.mock.On("WriteLoadControlLimits", entity, limits, resultCB)}
}

func (_c *CemOPEVInterface_WriteLoadControlLimits_Call) Run(run func(entity spine_goapi.EntityRemoteInterface, limits []api.LoadLimitsPhase, resultCB func(model.ResultDataType))) *CemOPEVInterface_WriteLoadControlLimits_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(spine_goapi.EntityRemoteInterface), args[1].([]api.LoadLimitsPhase), args[2].(func(model.ResultDataType)))
	})
	return _c
}

func (_c *CemOPEVInterface_WriteLoadControlLimits_Call) Return(_a0 *model.MsgCounterType, _a1 error) *CemOPEVInterface_WriteLoadControlLimits_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CemOPEVInterface_WriteLoadControlLimits_Call) RunAndReturn(run func(spine_goapi.EntityRemoteInterface, []api.LoadLimitsPhase, func(model.ResultDataType)) (*model.MsgCounterType, error)) *CemOPEVInterface_WriteLoadControlLimits_Call {
	_c.Call.Return(run)
	return _c
}

// NewCemOPEVInterface creates a new instance of CemOPEVInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCemOPEVInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *CemOPEVInterface {
	mock := &CemOPEVInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
