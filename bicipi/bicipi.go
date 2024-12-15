package bicipi

import (
	"time"

	"github.com/rcambrj/tacxble/ftms"
	"github.com/rcambrj/tacxble/tacx"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	SerialDevice         string
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

	// TODO: wait for tacx to be ready then advertise FTMS

	tacxReady := false
	ftmsStarted := false

	tacxService.On(func(event tacx.TacxEvent) {
		log.WithFields(log.Fields{"event": event}).Debugf("Tacx event")
		tacxReady = event.Ready

		if tacxReady {
			if !ftmsStarted {
				ftmsStarted = true
				ftmsService.Start()
			}
			ftmsService.SetState(ftms.State{
				Speed:   uint16(event.Speed * 100), // TODO what is this unit?
				Load:    int16(event.Load),
				Cadence: uint16(event.Cadence * 2), // TODO what is this unit?
			})
		}
	})
	ftmsService.On(func(event ftms.FTMSEvent) {
		log.WithFields(log.Fields{"event": event}).Debugf("BLE event")

		if !tacxReady {
			// this shouldn't happen as FTMS starts after ready
			log.Fatalf("unable to set tacx state: tacx not ready")
		}

		if event.Mode == ftms.ModeTargetPower {
			tacxService.SetState(tacx.State{
				Enabled:     true,
				Behaviour:   tacx.BehaviourERG,
				TargetWatts: float64(event.TargetPower),
			})
		}

		if event.Mode == ftms.ModeIndoorBikeSimulation {
			tacxService.SetState(tacx.State{
				Enabled:   true,
				Behaviour: tacx.BehaviourSimulator,
				// Weight: ??
				WindSpeed: 0,
				Gradient:  0,
				// RollingResistance: 0,
				// WindResistance:    0,

				// draftingFactor
			})
		}
	})

	for {
		time.Sleep(time.Hour)
	}
}
