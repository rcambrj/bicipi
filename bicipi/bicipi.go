package bicipi

import (
	"time"

	"github.com/rcambrj/tacxble/tacx"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	SerialDevice    string
	BluetoothDevice string
	Calibrate       bool
}

func Start(config Config) {
	log.Info("starting...")

	tacx.Start(tacx.Config{
		Device:    config.SerialDevice,
		Calibrate: config.Calibrate,
	})
	// TODO: wait for tacx to be ready then advertise FTMS
	// ftms.Start()

	for {
		// Sleep forever.
		time.Sleep(time.Hour)
	}
}
