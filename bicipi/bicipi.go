package bicipi

import (
	"time"

	"github.com/rcambrj/tacxble/ftms"
	"github.com/rcambrj/tacxble/tacx"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	SerialDevice         string
	BluetoothDevice      string
	BluetoothName        string
	Calibrate            bool
	Slow                 bool
	CalibrationSpeed     int
	CalibrationMin       int
	CalibrationMax       int
	CalibrationTolerance int
}

func Start(config Config) {
	log.Info("starting...")

	tacxService := tacx.MakeService(tacx.Config{
		Device:               config.SerialDevice,
		Calibrate:            config.Calibrate,
		Slow:                 config.Slow,
		CalibrationSpeed:     config.CalibrationSpeed,
		CalibrationMin:       config.CalibrationMin,
		CalibrationMax:       config.CalibrationMax,
		CalibrationTolerance: config.CalibrationTolerance,
	})
	tacxService.Start()

	ftmsService := ftms.MakeService(ftms.Config{
		BluetoothName: config.BluetoothName,
	})
	ftmsService.Start()

	// TODO: wait for tacx to be ready then advertise FTMS

	tacxService.On(func(event tacx.TacxEvent) {
		// TODO
	})
	ftmsService.On(func(event ftms.FTMSEvent) {
		// TODO

		// var enabled = true           // TODO: get this from BLE
		// var behaviour = BehaviourERG // TODO: get this from BLE
		// var targetWatts = 100.0      // TODO: get this from BLE
		// var weight = 80              // TODO: get this from BLE
		// var windSpeed = 0            // TODO: get this from BLE
		// var draftingFactor = 1       // TODO: get this from BLE
		// var gradient = 3             // TODO: get this from BLE

		// var enabled = true             // TODO: get this from BLE
		// var behaviour = BehaviourSlope // TODO: get this from BLE
		// var targetWatts = 0.0          // TODO: get this from BLE
		// var weight = 80                // TODO: get this from BLE
		// var windSpeed = 0              // TODO: get this from BLE
		// var draftingFactor = 1         // TODO: get this from BLE
		// var gradient = 3               // TODO: get this from BLE
	})

	for {
		time.Sleep(10 * time.Hour)
	}
}
