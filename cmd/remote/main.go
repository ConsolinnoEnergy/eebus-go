package main

import (
	"context"
	"crypto/tls"
	"flag"
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

type eebusConfiguration struct {
	vendorCode   string
	deviceBrand  string
	deviceModel  string
	serialNumber string
}

func loadCertificate(config eebusConfiguration, crtPath, keyPath string) tls.Certificate {
	certificate, err := tls.LoadX509KeyPair(crtPath, keyPath)
	if err != nil {
		certificate, err = cert.CreateCertificate(config.vendorCode, config.deviceModel, "DE", config.serialNumber)
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
	config := eebusConfiguration{}

	flag.StringVar(&config.vendorCode, "vendor", "", "EEBus vendor code")
	flag.StringVar(&config.deviceBrand, "brand", "", "EEBus device brand")
	flag.StringVar(&config.deviceModel, "model", "", "EEBus device model")
	flag.StringVar(&config.serialNumber, "serial", "", "EEBus device serial")

	flag.Parse()

	if config.serialNumber == "" {
		serialNumber, err := os.Hostname()
		if err != nil {
			log.Fatal(err)
		}
		config.serialNumber = serialNumber
	}

	if config.vendorCode == "" || config.deviceBrand == "" || config.deviceModel == "" {
		flag.Usage()
		return
	}

	certificate := loadCertificate(config, "cert.pem", "key.pem")

	configuration, err := api.NewConfiguration(
		config.vendorCode, config.deviceBrand, config.deviceModel, config.serialNumber,
		model.DeviceTypeTypeEnergyManagementSystem,
		[]model.EntityTypeType{model.EntityTypeTypeGridGuard, model.EntityTypeTypeCEM},
		23292, certificate, time.Second*4)

	r, err := NewRemote(configuration)
	if err != nil {
		log.Fatal(err)
	}

	r.RegisterUseCase(model.EntityTypeTypeCEM, "EG-LPC", func(localEntity spineapi.EntityLocalInterface, eventCB api.EntityEventCallback) api.UseCaseInterface {
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
