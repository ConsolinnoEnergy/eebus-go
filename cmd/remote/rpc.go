package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/service"
	"golang.org/x/exp/jsonrpc2"

	shipapi "github.com/enbility/ship-go/api"
)

type Remote struct {
	rpc     *jsonrpc2.Server
	service *service.Service

	connections    []*jsonrpc2.Connection
	remoteServices []shipapi.RemoteService

	calls map[string]any
}

func NewRemote(configuration *api.Configuration) (*Remote, error) {
	r := Remote{
		calls: make(map[string]any),
	}
	r.service = service.NewService(configuration, &r)
	r.registerCall("service", "RegisterRemoteSKI", r.service.RegisterRemoteSKI)
	r.registerCall("service", "UnregisterRemoteSKI", r.service.UnregisterRemoteSKI)

	r.registerCall("remote", "RemoteServices", func() []shipapi.RemoteService {
		return r.remoteServices
	})

	return &r, nil
}

func (r *Remote) Listen(context context.Context, network, address string) error {
	listener, err := jsonrpc2.NetListener(context, network, address, jsonrpc2.NetListenOptions{})
	if err != nil {
		return err
	}

	connOpts := jsonrpc2.ConnectionOptions{
		Framer:    jsonrpc2.RawFramer(),
		Preempter: nil,
		Handler:   jsonrpc2.HandlerFunc(r.handleRPC),
	}

	conn, err := jsonrpc2.Serve(context, listener, connOpts)
	if err != nil {
		return err
	}
	r.rpc = conn

	if err := r.service.Setup(); err != nil {
		return err
	}

	r.service.Start()
	go func() {
		<-context.Done()
		r.service.Shutdown()
	}()

	return nil
}

func (r *Remote) registerCall(group, name string, method any) {
	methodValue := reflect.ValueOf(method)
	if methodValue.Kind() != reflect.Func {
		panic(fmt.Sprintf("registerCall must be called with a function argument, found: %s", methodValue.Kind().String()))
	}

	r.calls[fmt.Sprintf("%s/%s", group, name)] = method
}

func (r *Remote) handleRPC(ctx context.Context, req *jsonrpc2.Request) (interface{}, error) {
	if req.IsCall() {
		method, found := r.calls[req.Method]
		if !found {
			// TODO: unify error reporting
			return jsonrpc2.NewResponse(req.ID, nil, errors.New("unknown method"))
		}

		var params []interface{}
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return jsonrpc2.NewResponse(req.ID, nil, err)
		}

		methodType := reflect.TypeOf(method)
		neededParams := methodType.NumIn()
		if len(params) != neededParams {
			// TODO: unify error reporting
			return jsonrpc2.NewResponse(req.ID, nil, errors.New("type error"))
		}

		methodParams := make([]reflect.Value, neededParams)
		for i := 0; i < neededParams; i++ {
			paramType := methodType.In(i)
			paramValue := reflect.ValueOf(params[i])

			if !paramValue.CanConvert(paramType) {
				// TODO: unify error reporting
				return jsonrpc2.NewResponse(req.ID, nil, errors.New("type error"))
			}
			methodParams[i] = paramValue.Convert(paramType)
		}

		output := reflect.ValueOf(method).Call(methodParams)
		switch methodType.NumOut() {
		case 0:
			return jsonrpc2.NewResponse(req.ID, []interface{}{}, nil)
		case 1:
			// TODO: verify jsonrpc output on reflect []Values
			return jsonrpc2.NewResponse(req.ID, output[0], nil)
		default:
			// TODO: verify jsonrpc output on reflect []Values
			return jsonrpc2.NewResponse(req.ID, output, nil)
		}
	} else {
		// RPC Notification
		// TODO: implement
	}

	return nil, nil
}

// Implement api.ServiceReaderInterface
func (r Remote) RemoteSKIConnected(service api.ServiceInterface, ski string) {
	for _, conn := range r.connections {
		conn.Notify(context.Background(), "RemoteSKIConnected", []string{ski})
	}
}

func (r Remote) RemoteSKIDisconnected(service api.ServiceInterface, ski string) {
	for _, conn := range r.connections {
		conn.Notify(context.Background(), "RemoteSKIDisconnected", []string{ski})
	}
}

func (r *Remote) VisibleRemoteServicesUpdated(service api.ServiceInterface, entries []shipapi.RemoteService) {
	r.remoteServices = entries
}

func (r Remote) ServiceShipIDUpdate(ski string, shipdID string) {
}

func (r Remote) ServicePairingDetailUpdate(ski string, detail *shipapi.ConnectionStateDetail) {
}
