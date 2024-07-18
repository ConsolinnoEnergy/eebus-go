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

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/service"
	"golang.org/x/exp/jsonrpc2"

	shipapi "github.com/enbility/ship-go/api"
)

type rpcService interface {
	Call(string, []interface{}) ([]reflect.Value, error)
}

type methodProxy struct {
	name   string
	rcvr   reflect.Value
	typ    reflect.Type
	method map[string]reflect.Method
}

func (svc *methodProxy) Call(methodName string, params []interface{}) ([]reflect.Value, error) {
	method, found := svc.method[methodName]
	if !found {
		return nil, jsonrpc2.ErrNotHandled
	}

	methodType := method.Type
	neededParams := methodType.NumIn()
	// don't count receiver as needed -> neededParams - 1
	if len(params) != (neededParams - 1) {
		return nil, jsonrpc2.ErrInvalidParams
	}

	methodParams := make([]reflect.Value, neededParams)
	methodParams[0] = svc.rcvr
	for dstIndex := 1; dstIndex < neededParams; dstIndex++ {
		paramType := methodType.In(dstIndex)
		// i - 1 due to receiver offset
		paramIndex := dstIndex - 1
		paramValue := reflect.ValueOf(params[paramIndex])

		if !paramValue.CanConvert(paramType) {
			return nil, jsonrpc2.ErrInvalidParams
		}
		methodParams[dstIndex] = paramValue.Convert(paramType)
	}

	output := method.Func.Call(methodParams)

	return output, nil
}

type Remote struct {
	rpc     *jsonrpc2.Server
	service *service.Service

	connections    []*jsonrpc2.Connection
	remoteServices []shipapi.RemoteService

	rpcServices map[string]rpcService
}

func (r Remote) RemoteServices() []shipapi.RemoteService {
	return r.remoteServices
}

func (r Remote) LocalSKI() string {
	return r.service.LocalService().SKI()
}

func NewRemote(configuration *api.Configuration) (*Remote, error) {
	r := Remote{
		connections:    []*jsonrpc2.Connection{},
		remoteServices: []shipapi.RemoteService{},

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
	return r.registerMethods(rcvr, name, false)
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

		var params []interface{}
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return nil, jsonrpc2.ErrParse
		}

		output, err := svc.Call(methodName, params)
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
			// TODO: handle output[0] == error specially
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
	for _, conn := range r.connections {
		conn.Notify(context.Background(), "remote/RemoteSKIConnected", []string{ski})
	}
}

func (r Remote) RemoteSKIDisconnected(service api.ServiceInterface, ski string) {
	for _, conn := range r.connections {
		conn.Notify(context.Background(), "remote/RemoteSKIDisconnected", []string{ski})
	}
}

func (r *Remote) VisibleRemoteServicesUpdated(service api.ServiceInterface, entries []shipapi.RemoteService) {
	r.remoteServices = entries
}

func (r Remote) ServiceShipIDUpdate(ski string, shipdID string) {
}

func (r Remote) ServicePairingDetailUpdate(ski string, detail *shipapi.ConnectionStateDetail) {
}
