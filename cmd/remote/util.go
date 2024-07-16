package main

import (
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"go/token"
	"os"
	"reflect"
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
