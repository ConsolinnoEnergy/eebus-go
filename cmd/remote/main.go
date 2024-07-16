package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/usecases/eg/lpc"
	"github.com/enbility/ship-go/cert"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

func loadCertificate(crtPath, keyPath string) tls.Certificate {
	certificate, err := tls.LoadX509KeyPair(crtPath, keyPath)
	if err != nil {
		certificate, err = cert.CreateCertificate("Demo", "Demo", "DE", "Demo-Unit-02")
		if err != nil {
			log.Fatal(err)
		}

		if err = WriteKey(certificate, keyPath); err != nil {
			log.Fatal(err)
		}
		if err = WriteCertificate(certificate, crtPath); err != nil {
			log.Fatal(err)
		}
	}

	return certificate
}

func main() {
	certificate := loadCertificate("cert.pem", "key.pem")

	configuration, err := api.NewConfiguration(
		"Demo", "Demo", "HEMS", "898237",
		model.DeviceTypeTypeEnergyManagementSystem,
		[]model.EntityTypeType{model.EntityTypeTypeGridGuard, model.EntityTypeTypeCEM},
		23292, certificate, time.Second*4)

	r, err := NewRemote(configuration)
	if err != nil {
		log.Fatal(err)
	}

	r.RegisterUseCase(model.EntityTypeTypeCEM, "LPC", func(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) api.UseCaseInterface {
		return lpc.NewLPC(localEntity, eventCB)
	})

	ctx, cancelCtx := context.WithCancel(context.Background())
	if err = r.Listen(ctx, "tcp", net.JoinHostPort("::", strconv.Itoa(3393))); err != nil {
		log.Fatal(err)
	}
	log.Print("Started")

	// Clean exit to make sure mdns shutdown is invoked
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	// User exit

	cancelCtx()
}
