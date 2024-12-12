package tacx

import (
	"time"

	"github.com/montanaflynn/stats"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Device           string
	Calibrate        bool
	Slow             bool
	CalibrationSpeed int
	CalibrationMin   int
	CalibrationMax   int
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
	var calibrationSpeed = int16(getRawSpeed(float64(config.CalibrationSpeed)))
	var calibrationDurationMin = time.Duration(config.CalibrationMin) * time.Second
	var calibrationDurationMax = time.Duration(config.CalibrationMax) * time.Second

	// var enabled = true           // TODO: get this from BLE
	// var behaviour = BehaviourERG // TODO: get this from BLE
	// var targetWatts = 100.0      // TODO: get this from BLE
	// var weight = 80              // TODO: get this from BLE
	// var windSpeed = 0            // TODO: get this from BLE
	// var draftingFactor = 1       // TODO: get this from BLE
	// var gradient = 3             // TODO: get this from BLE

	var enabled = true             // TODO: get this from BLE
	var behaviour = BehaviourSlope // TODO: get this from BLE
	var targetWatts = 0.0          // TODO: get this from BLE
	var weight = 80                // TODO: get this from BLE
	var windSpeed = 0              // TODO: get this from BLE
	var draftingFactor = 1         // TODO: get this from BLE
	var gradient = 3               // TODO: get this from BLE

	var controlCommandsPerSecond = 5
	if config.Slow {
		controlCommandsPerSecond = 1
	}

	var calibrating = config.Calibrate
	var calibrationStartedAt time.Time
	var calibrationLastLoads = make([]float64, 0, controlCommandsPerSecond*10)
	var calibrationTolerance = 10.0     // acceptable difference between percentiles to conclude
	var calibrationResult uint16 = 1040 // a sensible default, in case calibration is disabled

	lastResponse := controlResponse{}
	for {
		startTime := time.Now()

		command := controlCommand{
			keepalive: lastResponse.keepalive,
			adjust:    calibrationResult,
		}

		if calibrating {
			command.mode = modeCalibrating
			command.targetSpeed = calibrationSpeed
			command.weight = lowestWeight
			log.Debug("mode: calibrating")
		} else if !enabled {
			command.mode = modeOff
			log.Debug("mode: off")
		} else {
			command.mode = modeNormal
			switch behaviour {
			case BehaviourERG:
				command.weight = lowestWeight
				command.targetLoad = getTargetLoad(targetLoadArgs{
					targetWatts:  targetWatts,
					currentSpeed: lastResponse.speed,
				})
				log.Warnf("mode: normal; behaviour: erg; watts: %v; speed: %v; target %v", targetWatts, lastResponse.speed, command.targetLoad)
			case BehaviourSlope:
				command.weight = uint8(weight)
				targetWatts := getWattsForSlope(targetLoadForSlopeArgs{
					currentSpeed:   lastResponse.speed,
					weight:         weight,
					windSpeed:      windSpeed,
					draftingFactor: draftingFactor,
					gradient:       gradient,
				})
				command.targetLoad = getTargetLoad(targetLoadArgs{
					targetWatts:  targetWatts,
					currentSpeed: lastResponse.speed,
				})
				log.Warnf("mode: normal; behaviour: slope; gradient: %v; speed: %v; target %v", gradient, lastResponse.speed, command.targetLoad)
			}
		}

		controlResponse, err := sendControl(commander, command)
		if err != nil {
			log.Fatalf("unable to execute main command: %+v", err)
		}

		if calibrating {
			// calibrating mode means:
			// 1. waiting for the user to push the pedal
			// 2. after, wait for the wheel to spin up
			// 3. after, wait for a minimum duration (motor+tyre warm up)
			// 4. after, wait for current load to stabilise
			// start the timer once step 2 begins
			if calibrationStartedAt.IsZero() {
				if controlResponse.speed > uint16(calibrationSpeed)/2 {
					calibrationStartedAt = time.Now()
				} else {
					log.Warnf("waiting for calibration: pedal once then stop")
				}
			} else {
				untilMinimum := time.Until(calibrationStartedAt.Add(calibrationDurationMin))
				untilMaximum := time.Until(calibrationStartedAt.Add(calibrationDurationMax))

				if len(calibrationLastLoads) == cap(calibrationLastLoads) {
					calibrationLastLoads = calibrationLastLoads[1:]
				}
				calibrationLastLoads = append(calibrationLastLoads, float64(controlResponse.currentLoad))

				if untilMaximum > 0 {
					if untilMinimum < 0 {
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
					log.Warnf("calibrating. minimum: %t; remaining: %.0f; speed: %v; resistance: %v;", untilMinimum < 0, untilMinimum.Seconds(), controlResponse.speed, controlResponse.currentLoad)
				} else {
					calibrating = false
					log.Fatal("calibration aborted: maximum time reached.")
				}
			}
		}

		lastResponse = controlResponse

		period := time.Second / time.Duration(controlCommandsPerSecond)
		time.Sleep(time.Until(startTime.Add(period)))
	}
}
