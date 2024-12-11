package tacx

import (
	"time"

	"github.com/montanaflynn/stats"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Device    string
	Calibrate bool
	Slow      bool
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

	var lowestWeight = uint8(0x0a)
	var calibrationSpeed = 20.0                   // km/h TODO: allow setting this via CLI
	var calibrationDurationMin = 30 * time.Second // TODO: allow setting this via CLI and set default to 5 minutes to allow motor+tyre to warm up
	var calibrationDurationMax = 8 * time.Minute  // TODO: allow setting this via CLI

	var enabled = false          // TODO: get this from BLE
	var behaviour = BehaviourERG // TODO: get this from BLE
	var weight = 80              // TODO: get this from BLE
	// var targetWatts = 100.0      // TODO: get this from BLE

	var calibrating = config.Calibrate
	var calibrationStartedAt time.Time
	var calibrationLastLoads = make([]float64, 0, 50)
	var calibrationTolerance = 10.0     // difference between 10th & 90th percentile values before concluding
	var calibrationResult uint16 = 1040 // a sensible default, in case calibration is disabled

	lastResponse := controlResponse{}
	for {
		startTime := time.Now()

		command := controlCommand{
			keepalive: lastResponse.keepalive,
			adjust:    calibrationResult,
		}

		if calibrating {
			log.Debug("mode: calibrating")
			command.mode = modeCalibrating
			command.targetSpeed = calibrationSpeed
			command.weight = lowestWeight
		} else if !enabled {
			log.Debug("mode: off")
			command.mode = modeOff
		} else {
			log.Debug("mode: running")
			command.mode = modeNormal
			command.targetLoad = 0
			// command.targetLoad = targetWatts * rawLoadFactor / float64(lastResponse.raw.speed) // TODO
			switch behaviour {
			case BehaviourERG:
				command.weight = lowestWeight
			case BehaviourSlope:
				command.weight = uint8(weight)
			}
		}

		controlResponse, err := sendControl(commander, command)
		if err != nil {
			log.Fatalf("unable to execute main command: %+v", err)
		}

		if calibrating {
			// calibrating mode means:
			// * waiting for the user to push the pedal
			// * after the pedal is pushed, waiting for the wheel to spin up
			// * after the wheel is at speed, waiting for a set duration
			// start the timer once the second phase begins
			if calibrationStartedAt.IsZero() {
				if controlResponse.speed > calibrationSpeed/2 {
					calibrationStartedAt = time.Now()
				} else {
					log.Warnf("waiting for calibration: pedal once then stop")
				}
			} else {
				minimumCrossed := time.Now().After(calibrationStartedAt.Add(calibrationDurationMin))
				remainingTime := time.Until(calibrationStartedAt.Add(calibrationDurationMax))

				if len(calibrationLastLoads) == cap(calibrationLastLoads) {
					calibrationLastLoads = calibrationLastLoads[1:]
				}
				calibrationLastLoads = append(calibrationLastLoads, float64(controlResponse.raw.currentLoad))

				if err != nil {
					log.Warnf("unable to calculate calibration mean: %v", err)
				}
				if remainingTime > 0 {
					if minimumCrossed {
						quartile1, _ := stats.Percentile(calibrationLastLoads, 25)
						quartile3, _ := stats.Percentile(calibrationLastLoads, 75)
						stable := quartile3-quartile1 <= calibrationTolerance
						average := quartile1 + ((quartile3 - quartile1) / 2)
						log.Warnf("calculating. stable: %t; quartile1: %.2f; quartile3: %.2f; average: %.2f", stable, quartile1, quartile3, average)

						if stable {
							calibrationResult = uint16(average)
							calibrating = false
						}
					}
					log.Warnf("calibrating. minimum: %t; remaining: %.0f; speed: %.2f; resistance: %v;", minimumCrossed, remainingTime.Seconds(), controlResponse.speed, controlResponse.raw.currentLoad)
				} else {
					log.Warnf("calibration aborted: maximum time reached.")
					calibrating = false
				}
			}
		}

		lastResponse = controlResponse

		period := 200 * time.Millisecond
		if config.Slow {
			period = period * 5
		}
		if !enabled {
			// save some power
			period = 2 * time.Second
		}
		time.Sleep(time.Until(startTime.Add(period)))
	}
}
