package ohpcf

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"time"
)

// Scenario 1

// The availability of an optional consumption of power [OHPCF-011/1].
//
// return true if the optional consumption of power is available
func (o *OHPCF) OptionalPowerConsumptionAvailable(entity spineapi.EntityRemoteInterface) (bool, error) {
	data, err := o.checkEntityTypeAndGetData(entity)
	if err != nil {
		return false, err
	}

	if !o.isDataAvailable(data) {
		return false, nil
	}

	if data.Alternatives[0].PowerSequence[0].State != nil &&
		data.Alternatives[0].PowerSequence[0].State.State != nil {
		return *data.Alternatives[0].PowerSequence[0].State.State == model.PowerSequenceStateTypeInactive ||
			*data.Alternatives[0].PowerSequence[0].State.State == model.PowerSequenceStateTypeScheduled, nil
	}

	return false, nil
}

// The power value [OHPCF-011/2/1].
//
// return the power value
func (o *OHPCF) RequestedPowerEstimate(entity spineapi.EntityRemoteInterface) (float64, error) {
	return o.powerOfType(entity, model.PowerTimeSlotValueTypeTypePower)
}

// The max power value [OHPCF-011/2/2].
//
// return the maximal power value
func (o *OHPCF) RequestedPowerMax(entity spineapi.EntityRemoteInterface) (float64, error) {
	return o.powerOfType(entity, model.PowerTimeSlotValueTypeTypePowerMax)
}

// Indication whether the consumption may be stopped by the CEM [OHPCF-011/5].
//
// return true if the consumption may be stopped
func (o *OHPCF) ConsumptionIsStoppable(entity spineapi.EntityRemoteInterface) (bool, error) {
	data, err := o.checkEntityTypeAndGetData(entity)

	if err != nil {
		return false, err
	}

	if !o.isDataAvailable(data) {
		return false, api.ErrDataNotAvailable
	}

	if data.Alternatives[0].PowerSequence[0].OperatingConstraintsInterrupt != nil &&
		data.Alternatives[0].PowerSequence[0].OperatingConstraintsInterrupt.IsStoppable != nil {
		return *data.Alternatives[0].PowerSequence[0].OperatingConstraintsInterrupt.IsStoppable, nil
	}

	return false, api.ErrDataNotAvailable
}

// Indication whether the consumption may be paused and resumed by the CEM [OHPCF-011/6]
//
// return true if the consumption may be paused
func (o *OHPCF) ConsumptionIsPausable(entity spineapi.EntityRemoteInterface) (bool, error) {
	data, err := o.checkEntityTypeAndGetData(entity)
	if err != nil {
		return false, err
	}

	if !o.isDataAvailable(data) {
		return false, api.ErrDataNotAvailable
	}

	if data.Alternatives[0].PowerSequence[0].OperatingConstraintsInterrupt != nil &&
		data.Alternatives[0].PowerSequence[0].OperatingConstraintsInterrupt.IsPausable != nil {
		return *data.Alternatives[0].PowerSequence[0].OperatingConstraintsInterrupt.IsPausable, nil
	}

	return false, api.ErrDataNotAvailable
}

// The start time of the process [OHPCF-012/1].
//
// return the start time of the process
func (o *OHPCF) PowerConsumptionProcessStartTime(entity spineapi.EntityRemoteInterface) (time.Time, error) {
	data, err := o.checkEntityTypeAndGetData(entity)
	if err != nil {
		return time.Time{}, err
	}

	if !o.isDataAvailable(data) {
		return time.Time{}, api.ErrDataNotAvailable
	}

	if data.Alternatives[0].PowerSequence[0].Schedule != nil &&
		data.Alternatives[0].PowerSequence[0].Schedule.StartTime != nil {
		return data.Alternatives[0].PowerSequence[0].Schedule.StartTime.GetTime()
	}

	return time.Time{}, api.ErrDataNotAvailable
}

// The current state of this power consumption process [OHPCF-012/2].
//
// return the current state of this power consumption process
func (o *OHPCF) PowerConsumptionProcessState(entity spineapi.EntityRemoteInterface) (ucapi.CompressorPowerConsumptionStateType, error) {
	data, err := o.checkEntityTypeAndGetData(entity)
	if err != nil {
		return "", err
	}

	if !o.isDataAvailable(data) {
		return ucapi.CompressorPowerConsumptionStateStopped, nil
	}

	if data.Alternatives[0].PowerSequence[0].State != nil &&
		data.Alternatives[0].PowerSequence[0].State.State != nil {
		switch *data.Alternatives[0].PowerSequence[0].State.State {
		case model.PowerSequenceStateTypeInactive:
			return ucapi.CompressorPowerConsumptionStateAvailable, nil
		case model.PowerSequenceStateTypeScheduled:
			return ucapi.CompressorPowerConsumptionStateScheduled, nil
		case model.PowerSequenceStateTypeRunning:
			return ucapi.CompressorPowerConsumptionStateRunning, nil
		case model.PowerSequenceStateTypeScheduledPaused:
			return ucapi.CompressorPowerConsumptionStatePaused, nil
		case model.PowerSequenceStateTypeCompleted:
			return ucapi.CompressorPowerConsumptionStateCompleted, nil
		case model.PowerSequenceStateTypeInvalid:
			return ucapi.CompressorPowerConsumptionStateStopped, nil
		default:
			return "", api.ErrNotSupported
		}
	}

	return "", api.ErrDataNotAvailable
}

// The minimal time a consumption process must last [OHPCF-008].
//
// return the minimal time a consumption process must last
func (o *OHPCF) PowerConsumptionMinimalRunDuration(entity spineapi.EntityRemoteInterface) (time.Duration, error) {
	data, err := o.checkEntityTypeAndGetData(entity)
	if err != nil {
		return time.Duration(0), err
	}

	if !o.isDataAvailable(data) {
		return time.Duration(0), api.ErrDataNotAvailable
	}

	if data.Alternatives[0].PowerSequence[0].OperatingConstraintsDuration != nil &&
		data.Alternatives[0].PowerSequence[0].OperatingConstraintsDuration.ActiveDurationMin != nil {
		return data.Alternatives[0].PowerSequence[0].OperatingConstraintsDuration.ActiveDurationMin.GetTimeDuration()
	}

	return time.Duration(0), api.ErrDataNotAvailable
}

// The minimal time a pause of a consumption process must last [OHPCF-009].
//
// return the minimal time a pause of a consumption process must last
func (o *OHPCF) PowerConsumptionMinimalPauseDuration(entity spineapi.EntityRemoteInterface) (time.Duration, error) {
	data, err := o.checkEntityTypeAndGetData(entity)
	if err != nil {
		return time.Duration(0), err
	}

	if !o.isDataAvailable(data) {
		return time.Duration(0), api.ErrDataNotAvailable
	}

	if data.Alternatives[0].PowerSequence[0].OperatingConstraintsDuration != nil &&
		data.Alternatives[0].PowerSequence[0].OperatingConstraintsDuration.PauseDurationMin != nil {
		return data.Alternatives[0].PowerSequence[0].OperatingConstraintsDuration.PauseDurationMin.GetTimeDuration()
	}

	return time.Duration(0), api.ErrDataNotAvailable
}

// Scenario 2

// Schedule an optional power consumption process [OHPCF-004].
//
// note:
// A re-schedule of an already scheduled power consumption process is possible as long as the
// scheduled process did not start.
//
// parameters:
//   - start: The start time of the power consumption
func (o *OHPCF) SchedulePowerConsumptionProcess(entity spineapi.EntityRemoteInterface, start time.Time, resultCB func(result model.ResultDataType)) (*model.MsgCounterType, error) {
	data := &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{{
				Schedule: &model.PowerSequenceScheduleDataType{
					StartTime: model.NewAbsoluteOrRelativeTimeTypeFromTime(start),
				},
			}},
		}},
	}

	return o.writeSmartEnergyManagementData(entity, data, resultCB)
}

// stop (abort) the process [OHPCF-022/1].
func (o *OHPCF) AbortPowerConsumptionProcess(entity spineapi.EntityRemoteInterface, resultCB func(result model.ResultDataType)) (*model.MsgCounterType, error) {
	data := &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{{
				State: &model.PowerSequenceStateDataType{
					State: util.Ptr(model.PowerSequenceStateTypeInvalid),
				},
			}},
		}},
	}

	return o.writeSmartEnergyManagementData(entity, data, resultCB)
}

// pause the process [OHPCF-022/2].
func (o *OHPCF) PausePowerConsumptionProcess(entity spineapi.EntityRemoteInterface, resultCB func(result model.ResultDataType)) (*model.MsgCounterType, error) {
	data := &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{{
				State: &model.PowerSequenceStateDataType{
					State: util.Ptr(model.PowerSequenceStateTypeScheduledPaused),
				},
			}},
		}},
	}

	return o.writeSmartEnergyManagementData(entity, data, resultCB)
}

// resume the process [OHPCF-022/3].
func (o *OHPCF) ResumePowerConsumptionProcess(entity spineapi.EntityRemoteInterface, resultCB func(result model.ResultDataType)) (*model.MsgCounterType, error) {
	data := &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{{
				State: &model.PowerSequenceStateDataType{
					State: util.Ptr(model.PowerSequenceStateTypeRunning),
				},
			}},
		}},
	}

	return o.writeSmartEnergyManagementData(entity, data, resultCB)
}

// ------------------------ helper methods ------------------------ //

func (o *OHPCF) checkEntityTypeAndGetData(entity spineapi.EntityRemoteInterface) (*model.SmartEnergyManagementPsDataType, error) {
	if !o.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	smartEnergyManagementPs, err := client.NewSmartEnergyManagementPs(o.LocalEntity, entity)
	if err != nil {
		return nil, err
	}

	data, err := smartEnergyManagementPs.GetData()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (o *OHPCF) isDataAvailable(data *model.SmartEnergyManagementPsDataType) bool {
	return len(data.Alternatives) == 1 &&
		data.Alternatives[0].PowerSequence != nil &&
		len(data.Alternatives[0].PowerSequence) == 1
}

func (o *OHPCF) powerOfType(entity spineapi.EntityRemoteInterface, valueType model.PowerTimeSlotValueTypeType) (float64, error) {
	data, err := o.checkEntityTypeAndGetData(entity)
	if err != nil {
		return 0, err
	}

	if !o.isDataAvailable(data) {
		return 0, api.ErrDataNotAvailable
	}

	if data.Alternatives[0].PowerSequence[0].PowerTimeSlot != nil &&
		len(data.Alternatives[0].PowerSequence[0].PowerTimeSlot) == 1 &&
		data.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ValueList != nil &&
		data.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ValueList.Value != nil {

		for _, value := range data.Alternatives[0].PowerSequence[0].PowerTimeSlot[0].ValueList.Value {
			if value.Value != nil && value.ValueType != nil && *value.ValueType == valueType {
				return value.Value.GetValue(), nil
			}
		}
	}

	return 0, api.ErrDataNotAvailable
}

func (o *OHPCF) writeSmartEnergyManagementData(entity spineapi.EntityRemoteInterface, data *model.SmartEnergyManagementPsDataType, resultCB func(result model.ResultDataType)) (*model.MsgCounterType, error) {
	if !o.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	smartEnergyManagementPs, err := client.NewSmartEnergyManagementPs(o.LocalEntity, entity)
	if err != nil {
		return nil, err
	}

	msgCounter, err := smartEnergyManagementPs.WriteData(data)

	if err != nil {
		return nil, err
	}

	if resultCB != nil && msgCounter != nil {
		cb := func(msg spineapi.ResponseMessage) {
			response, ok := msg.Data.(*model.ResultDataType)
			if ok {
				resultCB(*response)
			}
		}
		if errCB := smartEnergyManagementPs.AddResponseCallback(*msgCounter, cb); errCB != nil {
			logging.Log().Debug("Failed to add response callback for msgCounter %v: %v", msgCounter, errCB)
		}
	}

	return msgCounter, nil
}
