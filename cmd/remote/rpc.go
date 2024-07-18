package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go/token"
	"log"
	"reflect"
	"strings"
	"sync"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/service"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"golang.org/x/exp/jsonrpc2"

	shipapi "github.com/enbility/ship-go/api"
)

type rpcService interface {
	Call(*Remote, string, json.RawMessage) ([]reflect.Value, error)
}

type methodProxy struct {
	name   string
	rcvr   reflect.Value
	typ    reflect.Type
	method map[string]reflect.Method
}

func (svc *methodProxy) Call(remote *Remote, methodName string, params json.RawMessage) ([]reflect.Value, error) {
	method, found := svc.method[methodName]
	if !found {
		return nil, jsonrpc2.ErrNotHandled
	}

	methodType := method.Type
	neededParams := methodType.NumIn()

	var decodedParams []interface{}
	for idx := 1; idx < neededParams; idx++ {
		paramType := methodType.In(idx)

		var paramValue reflect.Value
		if paramType == reflect.TypeFor[spineapi.EntityRemoteInterface]() {
			// convert between EntityRemoteInterface and EntityAddressType
			paramValue = reflect.New(reflect.TypeFor[model.EntityAddressType]())
		} else {
			paramValue = reflect.New(paramType)
		}

		decodedParams = append(decodedParams, paramValue.Interface())
	}

	log.Printf("pre: decodedParams %#v", decodedParams)
	if err := json.Unmarshal(params, &decodedParams); err != nil {
		return nil, jsonrpc2.ErrParse
	}
	log.Printf("post: decodedParams %v", decodedParams)

	methodParams := make([]reflect.Value, neededParams)
	methodParams[0] = svc.rcvr
	for dstIndex := 1; dstIndex < neededParams; dstIndex++ {
		paramType := methodType.In(dstIndex)
		// i - 1 due to receiver offset
		paramIndex := dstIndex - 1

		if paramType == reflect.TypeFor[spineapi.EntityRemoteInterface]() {
			// convert between EntityRemoteInterface and EntityAddressType
			address, ok := decodedParams[paramIndex].(*model.EntityAddressType)
			if !ok {
				return nil, jsonrpc2.ErrInvalidParams
			}
			log.Printf("entityInterfaces: %v", remote.entityInterfaces)
			log.Printf("address: %v", address)
			log.Printf("map: %v", remote.entityInterfaces[fmt.Sprintf("%s", address)])
			methodParams[dstIndex] = reflect.ValueOf(remote.entityInterfaces[fmt.Sprintf("%s", address)])
		} else if decodedParams[paramIndex] == nil {
			// some parameters are optional and allowed to be nil
			methodParams[dstIndex] = reflect.New(paramType).Elem()
		} else {
			methodParams[dstIndex] = reflect.ValueOf(decodedParams[paramIndex]).Elem()
		}
	}

	output := method.Func.Call(methodParams)

	return output, nil
}

type Remote struct {
	rpc     *jsonrpc2.Server
	service *service.Service

	connections         []*jsonrpc2.Connection
	remoteServices      []shipapi.RemoteService
	entityInterfaces    map[string]spineapi.EntityRemoteInterface
	entityInterfaceLock *sync.Mutex

	rpcServices map[string]rpcService
}

func (r Remote) RemoteServices() []shipapi.RemoteService {
	return r.remoteServices
}

func (r Remote) ConnectedDevices() []string {
	remoteDevices := r.service.LocalDevice().RemoteDevices()
	skiList := make([]string, len(remoteDevices))

	for i, dev := range remoteDevices {
		skiList[i] = dev.Ski()
	}

	return skiList
}

func (r Remote) LocalSKI() string {
	return r.service.LocalService().SKI()
}

func NewRemote(configuration *api.Configuration) (*Remote, error) {
	r := Remote{
		connections:         []*jsonrpc2.Connection{},
		remoteServices:      []shipapi.RemoteService{},
		entityInterfaces:    make(map[string]spineapi.EntityRemoteInterface),
		entityInterfaceLock: &sync.Mutex{},

		rpcServices: make(map[string]rpcService),
	}
	r.service = service.NewService(configuration, &r)

	r.RegisterMethods("service", r.service)
	r.RegisterMethods("remote", &r)

	if err := r.service.Setup(); err != nil {
		return nil, err
	}

	return &r, nil
}

func (r *Remote) Bind(context context.Context, conn *jsonrpc2.Connection) (jsonrpc2.ConnectionOptions, error) {
	connOpts := jsonrpc2.ConnectionOptions{
		Framer:    NewlineFramer{},
		Preempter: nil,
		Handler:   jsonrpc2.HandlerFunc(r.handleRPC),
	}

	r.connections = append(r.connections, conn)
	return connOpts, nil
}

func (r *Remote) Listen(context context.Context, network, address string) error {
	listener, err := jsonrpc2.NetListener(context, network, address, jsonrpc2.NetListenOptions{})
	if err != nil {
		return err
	}

	conn, err := jsonrpc2.Serve(context, listener, r)
	if err != nil {
		return err
	}
	r.rpc = conn

	r.service.Start()
	go func() {
		<-context.Done()
		r.service.Shutdown()
	}()

	return nil
}

func (r *Remote) RegisterMethods(name string, rcvr any) error {
	return r.registerMethods(rcvr, name, true)
}

func (r *Remote) registerMethods(rcvr any, name string, useName bool) error {
	c := new(methodProxy)
	c.typ = reflect.TypeOf(rcvr)
	c.rcvr = reflect.ValueOf(rcvr)
	sname := name
	if !useName {
		sname = reflect.Indirect(c.rcvr).Type().Name()
	}
	if sname == "" {
		s := "rpc.Register: no service name for type " + c.typ.String()
		log.Print(s)
		return errors.New(s)
	}
	if !useName && !token.IsExported(sname) {
		s := "rpc.Register: type " + sname + " is not exported"
		log.Print(s)
		return errors.New(s)
	}
	sname = strings.ToLower(sname)
	c.name = sname

	c.method = make(map[string]reflect.Method)
	for m := 0; m < c.typ.NumMethod(); m++ {
		method := c.typ.Method(m)
		mtype := method.Type
		mname := method.Name

		// Method bust be xeported
		if !method.IsExported() {
			continue
		}

		// all (non-receiver) arguments must be builtin or exported
		for i := 1; i < mtype.NumIn(); i++ {
			argType := mtype.In(i)
			if !isExportedOrBuiltinType(argType) {
				panic(fmt.Sprintf("UseCaseProxy.Register: argument type of method %q is not exported: %q\n", mname, argType))
			}
			continue
		}
		for i := 1; i < mtype.NumOut(); i++ {
			argType := mtype.Out(i)
			if !isExportedOrBuiltinType(argType) {
				panic(fmt.Sprintf("UseCaseProxy.Register: return type of method %q is not exported: %q\n", mname, argType))
			}
			continue
		}

		log.Printf("registering method %s/%s", sname, mname)
		c.method[strings.ToLower(mname)] = method
	}

	r.rpcServices[sname] = c
	return nil
}

func (r *Remote) handleRPC(ctx context.Context, req *jsonrpc2.Request) (interface{}, error) {
	if req.IsCall() {
		slash := strings.LastIndex(req.Method, "/")
		if slash < 0 {
			return nil, jsonrpc2.ErrMethodNotFound
		}
		serviceName := strings.ToLower(req.Method[:slash])
		methodName := strings.ToLower(req.Method[slash+1:])

		svc, found := r.rpcServices[serviceName]
		if !found {
			return nil, jsonrpc2.ErrMethodNotFound
		}

		output, err := svc.Call(r, methodName, req.Params)
		if err != nil {
			return nil, err
		}

		var resp interface{}
		numOut := len(output)
		switch numOut {
		case 0:
			resp = []interface{}{}
		case 1:
			resp = output[0].Interface()
		case 2:
			if output[1].Type().Implements(reflect.TypeFor[error]()) {
				if !output[1].IsNil() {
					log.Printf("error handling %v: %v", req.Method, output[1].Interface())
					return nil, jsonrpc2.ErrInternal
				}
				resp = output[0].Interface()
			} else {
				r := make([]interface{}, numOut)
				for i, e := range output {
					r[i] = e.Interface()
				}
				resp = r
			}
		default:
			r := make([]interface{}, numOut)
			for i, e := range output {
				r[i] = e.Interface()
			}
			resp = r
		}
		return resp, nil
	} else {
		// RPC Notification
		// TODO: implement
	}

	return nil, nil
}

// Implement api.ServiceReaderInterface
func (r Remote) RemoteSKIConnected(service api.ServiceInterface, ski string) {
	params := make(map[string]interface{}, 1)
	params["device"] = ski
	for _, conn := range r.connections {
		conn.Notify(context.Background(), "remote/RemoteSKIConnected", params)
	}
}

func (r Remote) RemoteSKIDisconnected(service api.ServiceInterface, ski string) {
	params := make(map[string]interface{}, 1)
	params["device"] = ski
	for _, conn := range r.connections {
		conn.Notify(context.Background(), "remote/RemoteSKIDisconnected", params)
	}
}

func (r *Remote) VisibleRemoteServicesUpdated(service api.ServiceInterface, entries []shipapi.RemoteService) {
	r.remoteServices = entries
}

func (r Remote) ServiceShipIDUpdate(ski string, shipdID string) {
}

func (r Remote) ServicePairingDetailUpdate(ski string, detail *shipapi.ConnectionStateDetail) {
}
