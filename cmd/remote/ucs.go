package main

import (
	"context"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/service"
	spineapi "github.com/enbility/spine-go/api"
)

type UseCaseId string
type UseCaseTypeType string

const (
	UseCaseTypeLPC UseCaseTypeType = "LPC"
)

type UseCaseBuilder func(*service.Service, api.EntityEventCallback) api.UseCaseInterface

func (r *Remote) RegisterUseCase(builder UseCaseBuilder) (UseCaseId, error) {
	var id UseCaseId = UseCaseId("test")

	uc := builder(r.service, func(
		ski string,
		device spineapi.DeviceRemoteInterface,
		entity spineapi.EntityRemoteInterface,
		event api.EventType,
	) {
		r.PropagateEvent(id, ski, device, entity, event)
	})
	r.service.AddUseCase(uc)

	return id, nil
}

func (r *Remote) PropagateEvent(
	id UseCaseId,
	ski string,
	device spineapi.DeviceRemoteInterface,
	entity spineapi.EntityRemoteInterface,
	event api.EventType,
) {
	for _, conn := range r.connections {
		conn.Notify(context.Background(), string(event), nil)
	}
}
