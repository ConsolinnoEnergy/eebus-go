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
) (map[model.ElectricalConnectionPhaseNameType]float64, error) {
	measurement, err := client.NewMeasurement(localEntity, remoteEntity)
	electricalConnection, err1 := client.NewElectricalConnection(localEntity, remoteEntity)
	if err != nil || err1 != nil {
		return nil, api.ErrMetadataNotAvailable
	}

	data, err := measurement.GetDataForFilter(measurementFilter)
	if err != nil || len(data) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	result := make(map[model.ElectricalConnectionPhaseNameType]float64, len(validPhaseNameTypes))

	for _, item := range data {
		if item.Value == nil || item.MeasurementId == nil {
			continue
		}

		filter := model.ElectricalConnectionParameterDescriptionDataType{
			MeasurementId: item.MeasurementId,
		}
		param, err := electricalConnection.GetParameterDescriptionsForFilter(filter)
		if err != nil || len(param) == 0 || param[0].AcMeasuredPhases == nil {
			// error getting parameter description
			continue
		}

		// calculate the offset into result for the measured phase
		phaseName := *param[0].AcMeasuredPhases
		if validPhaseNameTypes != nil &&
			!slices.Contains(validPhaseNameTypes, phaseName) {
			// ignore phase measurements not specified in validPhaseNameTypes
			continue
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

		result[phaseName] = value
	}

	return result, nil
}
