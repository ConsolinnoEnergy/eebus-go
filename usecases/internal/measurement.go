package internal

import (
	"slices"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// return the phase specific measurement data
func MeasurementPhaseSpecificDataForFilter(
	localEntity spineapi.EntityLocalInterface,
	remoteEntity spineapi.EntityRemoteInterface,
	measurementFilter model.MeasurementDescriptionDataType,
	energyDirection model.EnergyDirectionType,
	validPhaseNameTypes []model.ElectricalConnectionPhaseNameType,
) ([]float64, error) {
	measurement, err := client.NewMeasurement(localEntity, remoteEntity)
	electricalConnection, err1 := client.NewElectricalConnection(localEntity, remoteEntity)
	if err != nil || err1 != nil {
		return nil, api.ErrMetadataNotAvailable
	}

	data, err := measurement.GetDataForFilter(measurementFilter)
	if err != nil || len(data) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	var result []float64
	if validPhaseNameTypes != nil {
		// pre-allocate result array for each possible phase so we can add phases to it in arbitrary order
		result = make([]float64, len(validPhaseNameTypes))
	}

	for _, item := range data {
		if item.Value == nil || item.MeasurementId == nil {
			continue
		}

		phaseIndex := -1
		if validPhaseNameTypes != nil {
			filter := model.ElectricalConnectionParameterDescriptionDataType{
				MeasurementId: item.MeasurementId,
			}
			param, err := electricalConnection.GetParameterDescriptionsForFilter(filter)
			if err != nil || len(param) == 0 || param[0].AcMeasuredPhases == nil {
				// error getting parameter description
				continue
			}

			// calculate the offset into result for the measured phase
			phaseIndex = slices.Index(validPhaseNameTypes, *param[0].AcMeasuredPhases)
			if phaseIndex == -1 {
				// ignore phase measurements not specified in validPhaseNameTypes
				continue
			}
		}

		if energyDirection != "" {
			filter := model.ElectricalConnectionParameterDescriptionDataType{
				MeasurementId: item.MeasurementId,
			}
			desc, err := electricalConnection.GetDescriptionForParameterDescriptionFilter(filter)
			if err != nil || desc == nil {
				continue
			}

			// if energy direction is not consume
			if desc.PositiveEnergyDirection == nil || *desc.PositiveEnergyDirection != energyDirection {
				return nil, err
			}
		}

		// if the value state is set and not normal, the value is not valid and should be ignored
		// therefore we return an error
		if item.ValueState != nil && *item.ValueState != model.MeasurementValueStateTypeNormal {
			return nil, api.ErrDataInvalid
		}

		value := item.Value.GetValue()

		if phaseIndex < 0 {
			// measurement is not for a specific phase
			result = append(result, value)
		} else {
			// measurement is for a specifyc phase, store the value at the corresponding phaseIndex
			result[phaseIndex] = value
		}
	}

	return result, nil
}
