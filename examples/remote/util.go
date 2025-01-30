package main

import (
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"go/token"
	"os"
	"reflect"

	"github.com/enbility/spine-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"golang.org/x/exp/jsonrpc2"
)

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return token.IsExported(t.Name()) || t.PkgPath() == ""
}

// json.Marshall won't marshall error types, marshall as string
func errorAsJson(v reflect.Value) interface{} {
	if v.IsNil() {
		// passthrough nil as nil, otherwise nil.(error) will panic
		return nil
	} else {
		return v.Interface().(error).Error()
	}
}

func decodeVariadicFunctionArgument(paramType reflect.Type, params []json.RawMessage, decodedParams []interface{}) ([]interface{}, error) {
	var err error

	paramType = paramType.Elem()
	for _, param := range params {
		decodedParams, err = decodeFunctionArgument(paramType, param, decodedParams)
		if err != nil {
			return nil, err
		}
	}

	return decodedParams, nil
}

func decodeFunctionArgument(paramType reflect.Type, param json.RawMessage, decodedParams []interface{}) ([]interface{}, error) {
	var paramValue reflect.Value
	if paramType == reflect.TypeFor[spineapi.DeviceRemoteInterface]() {
		// convert between DeviceRemoteInterface and DeviceAddressType
		paramValue = reflect.New(reflect.TypeFor[model.DeviceAddressType]())
	} else if paramType == reflect.TypeFor[spineapi.EntityRemoteInterface]() {
		// convert between EntityRemoteInterface and EntityAddressType
		paramValue = reflect.New(reflect.TypeFor[model.EntityAddressType]())
	} else {
		paramValue = reflect.New(paramType)
	}

	decodedParam := paramValue.Interface()
	if err := json.Unmarshal(param, &decodedParam); err != nil {
		return nil, jsonrpc2.ErrParse
	}
	decodedParams = append(decodedParams, decodedParam)

	return decodedParams, nil

}

func transformVariadicFunctionArgument(remote *Remote, paramType reflect.Type, params []interface{}, methodParams []reflect.Value) ([]reflect.Value, error) {
	var err error

	paramType = paramType.Elem()
	for _, param := range params {
		methodParams, err = transformFunctionArgument(remote, paramType, param, methodParams)
		if err != nil {
			return nil, err
		}
	}

	return methodParams, nil
}

func transformFunctionArgument(remote *Remote, paramType reflect.Type, param interface{}, methodParams []reflect.Value) ([]reflect.Value, error) {
	if paramType == reflect.TypeFor[spineapi.DeviceRemoteInterface]() {
		// convert between DeviceRemoteInterface and DeviceAddressType
		address, ok := param.(*model.DeviceAddressType)
		if !ok || address.Device == nil {
			return nil, jsonrpc2.ErrInvalidParams
		}

		deviceInterface := remote.service.LocalDevice().RemoteDeviceForAddress(*address.Device)
		if deviceInterface == nil {
			return nil, jsonrpc2.ErrInvalidParams
		}

		methodParams = append(methodParams, reflect.ValueOf(deviceInterface))
	} else if paramType == reflect.TypeFor[spineapi.EntityRemoteInterface]() {
		// convert between EntityRemoteInterface and EntityAddressType
		address, ok := param.(*model.EntityAddressType)
		if !ok || address.Device == nil {
			return nil, jsonrpc2.ErrInvalidParams
		}

		deviceInterface := remote.service.LocalDevice().RemoteDeviceForAddress(*address.Device)
		if deviceInterface == nil {
			return nil, jsonrpc2.ErrInvalidParams
		}

		entityInterface := deviceInterface.Entity(address.Entity)
		if entityInterface == nil {
			return nil, jsonrpc2.ErrInvalidParams
		}

		methodParams = append(methodParams, reflect.ValueOf(entityInterface))
	} else if param == nil {
		// some parameters are optional and allowed to be nil
		methodParams = append(methodParams, reflect.New(paramType).Elem())
	} else {
		methodParams = append(methodParams, reflect.ValueOf(param).Elem())
	}

	return methodParams, nil
}

func transformReturnValues(values []reflect.Value) []interface{} {
	result := make([]interface{}, len(values))

	for i, e := range values {
		valueType := e.Type()

		switch valueType {
		case reflect.TypeFor[spineapi.DeviceRemoteInterface]():
			result[i] = e.Interface().(spineapi.DeviceRemoteInterface).Address()
		case reflect.TypeFor[[]spineapi.DeviceRemoteInterface]():
			rawValues := e.Interface().([]api.DeviceRemoteInterface)
			transformedValues := make([]model.AddressDeviceType, len(rawValues))

			for j, r := range rawValues {
				transformedValues[j] = *r.Address()
			}
			result[i] = transformedValues
		case reflect.TypeFor[spineapi.EntityRemoteInterface]():
			result[i] = e.Interface().(spineapi.EntityRemoteInterface).Address()
		case reflect.TypeFor[[]spineapi.EntityRemoteInterface]():
			rawValues := e.Interface().([]api.EntityRemoteInterface)
			transformedValues := make([]model.EntityAddressType, len(rawValues))

			for j, r := range rawValues {
				transformedValues[j] = *r.Address()
			}
			result[i] = transformedValues
		default:
			switch {
			case valueType.Implements(reflect.TypeFor[error]()):
				result[i] = errorAsJson(e)
			default:
				result[i] = e.Interface()
			}
		}
	}

	return result
}

func WriteKey(cert tls.Certificate, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	switch v := cert.PrivateKey.(type) {
	case *ecdsa.PrivateKey:
		bytes, err := x509.MarshalECPrivateKey(v)
		if err != nil {
			return err
		}

		err = pem.Encode(file, &pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: bytes,
		})
	default:
		return fmt.Errorf("Unable to serialize private key of type %T", v)
	}

	return nil
}

func WriteCertificate(cert tls.Certificate, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, leaf := range cert.Certificate {
		err = pem.Encode(file, &pem.Block{Type: "CERTIFICATE", Bytes: leaf})
		if err != nil {
			return err
		}
	}

	return nil
}
