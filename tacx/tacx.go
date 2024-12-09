package tacx

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Device    string
	Calibrate bool
}

func Start(config Config) {
	port, err := connect(config.Device)
	if err != nil {
		log.Fatalf("unable to connect to tacx: %v", err)
	}

	commander := makeCommander(port)

	_, err = getVersion(commander)
	if err != nil {
		log.Fatal(err)
	}

	var ergModeWeight = 10
	var calibrationSpeed = 20.0                // km/h TODO: allow setting this via CLI
	var calibrationDuration = 15 * time.Second // TODO: allow setting this via CLI
	var calibrationStableRange = 0.01
	var adjust int8 = 0 // a neutral value. TODO: allow setting this via CLI

	var enabled = true           // TODO: get this from BLE
	var behaviour = BehaviourERG // TODO: get this from BLE
	var weight = 80              // TODO: get this from BLE
	var targetLoad = 100.0       // TODO: get this from BLE

	var calibrating = config.Calibrate
	var calibrationStartedAt time.Time

	lastResponse := controlResponse{}
	for {
		startTime := time.Now()

		var mode mode
		if enabled {
			if calibrating {
				mode = modeCalibrating
				log.Info("mode: calibrating")
			} else {
				mode = modeRunning
				log.Info("mode: running")
			}
		} else {
			mode = modeOff
			log.Info("mode: off")
		}

		command := controlCommand{
			targetLoad:  targetLoad,       // sendControl() will ignore this while calibrating
			targetSpeed: calibrationSpeed, // sendControl() will ignore this while not calibrating
			mode:        mode,
			keepalive:   lastResponse.keepalive,
			weight:      uint8(weight),
			adjust:      adjust,
		}

		if behaviour == BehaviourERG {
			command.weight = uint8(ergModeWeight)
		}

		controlResponse, err := sendControl(commander, command)
		if err != nil {
			log.Fatalf("unable to execute main command: %+v", err)
		}

		if calibrating {
			if calibrationStartedAt.IsZero() {
				if controlResponse.speed > calibrationSpeed/2 {
					// calibrating mode means:
					// * waiting for the user to push the pedal
					// * after the pedal is pushed, waiting for the motor to calibrate itself
					// start the timer once the second phase begins
					calibrationStartedAt = time.Now()
				} else {
					log.Warnf("waiting for calibration: pedal once then stop")
				}
			} else {

				stable := calibrationSpeed*(1-calibrationStableRange) < controlResponse.speed && controlResponse.speed < calibrationSpeed*(1+calibrationStableRange)
				remaining := time.Until(calibrationStartedAt.Add(calibrationDuration))

				if remaining > 0 {
					log.Warnf("calibrating. remaining: %v; speed: %v; stable: %v; resistance: %v;", remaining, controlResponse.speed, stable, controlResponse.currentLoad)
				} else {
					calibrating = false
					log.Warnf("calibration complete")
				}
			}
		}

		lastResponse = controlResponse

		period := 1000 * time.Millisecond
		if mode == modeOff {
			// save some power
			period = 2 * time.Second
		}
		time.Sleep(time.Until(startTime.Add(period)))
	}
}
