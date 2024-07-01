package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/ship-go/cert"
	"github.com/enbility/spine-go/model"
)

func main() {
	certificate, err := cert.CreateCertificate("Demo", "Demo", "DE", "Demo-Unit-02")
	if err != nil {
		log.Fatal(err)
	}

	configuration, err := api.NewConfiguration(
		"Demo", "Demo", "HEMS", "898237",
		model.DeviceTypeTypeEnergyManagementSystem,
		[]model.EntityTypeType{model.EntityTypeTypeGridGuard, model.EntityTypeTypeCEM},
		23292, certificate, time.Second*4)

	r, err := NewRemote(configuration)
	if err != nil {
		log.Fatal(err)
	}

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
