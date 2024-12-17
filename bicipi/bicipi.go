package bicipi

import (
	"fmt"
	"time"

	"github.com/rcambrj/tacxble/ftms"
	"github.com/rcambrj/tacxble/tacx"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Weight               uint8
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

	tacxReady := false
	ftmsStarted := false

	tacxService := tacx.MakeService(tacx.Config{
		Weight:               config.Weight,
		Device:               config.SerialDevice,
		Calibrate:            config.Calibrate,
		Slow:                 config.Slow,
		CalibrationSpeed:     config.CalibrationSpeed,
		CalibrationMin:       config.CalibrationMin,
		CalibrationMax:       config.CalibrationMax,
		CalibrationTolerance: config.CalibrationTolerance,
	})

	ftmsService := ftms.MakeService(ftms.Config{
		BluetoothName: config.BluetoothName,
	})

	tacxService.On(func(event tacx.TacxEvent) {
		log.WithFields(log.Fields{"event": fmt.Sprintf("%+v", event)}).Infof("tacx event")
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
		log.WithFields(log.Fields{"event": fmt.Sprintf("%+v", event)}).Infof("ble event")

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
				Enabled:           true,
				Behaviour:         tacx.BehaviourSimulator,
				WindSpeed:         event.WindSpeed,
				Gradient:          event.TargetGrade,
				RollingResistance: event.RollingResistance,
				WindResistance:    event.WindResistance,
			})
		}
	})

	tacxService.Start()

	for {
		time.Sleep(time.Hour)
	}
}
