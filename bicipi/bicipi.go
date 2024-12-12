package bicipi

import (
	"time"

	"github.com/rcambrj/tacxble/tacx"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	SerialDevice     string
	BluetoothDevice  string
	Calibrate        bool
	Slow             bool
	CalibrationSpeed int
	CalibrationMin   int
	CalibrationMax   int
}

func Start(config Config) {
	log.Info("starting...")

	tacx.Start(tacx.Config{
		Device:           config.SerialDevice,
		Calibrate:        config.Calibrate,
		Slow:             config.Slow,
		CalibrationSpeed: config.CalibrationSpeed,
		CalibrationMin:   config.CalibrationMin,
		CalibrationMax:   config.CalibrationMax,
	})
	// TODO: wait for tacx to be ready then advertise FTMS
	// ftms.Start()

	for {
		// Sleep forever.
		time.Sleep(time.Hour)
	}
}
