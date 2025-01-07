package tacx

import (
	"fmt"
	"sync"
	"time"

	"github.com/montanaflynn/stats"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Weight               uint8
	SerialDevice         string
	Calibrate            bool
	Slow                 bool
	CalibrationSpeed     int
	CalibrationMin       int
	CalibrationMax       int
	CalibrationTolerance int
}

type State struct {
	Enabled   bool
	Behaviour Behaviour

	// BehaviourERG
	TargetWatts float64

	// BehaviourSimulator
	WindSpeed         float64
	Gradient          float64
	RollingResistance float64
	WindResistance    float64
}

type TacxEvent struct {
	Ready   bool
	Speed   float64 // km/h
	Load    float64 // watts
	Cadence uint8   // rpm
}

type Listener = func(event TacxEvent)

func MakeService(config Config) Tacx {
	return Tacx{
		config:  config,
		channel: make(chan TacxEvent),
	}
}

type Tacx struct {
	config    Config
	stateLock sync.Mutex
	state     State
	channel   chan TacxEvent
	listeners []Listener
}

func (t *Tacx) SetState(state State) {
	t.stateLock.Lock()
	defer t.stateLock.Unlock()

	t.state = state
}

func (t *Tacx) getState() State {
	t.stateLock.Lock()
	defer t.stateLock.Unlock()

	return t.state
}

func (t *Tacx) On(listener Listener) {
	t.listeners = append(t.listeners, listener)
}

func (t *Tacx) Start() {
	go t.startEventLoop()
	go t.startSerialLoop()
}

func (t *Tacx) startEventLoop() {
	for {
		select {
		case msg := <-t.channel:
			for _, listener := range t.listeners {
				listener(msg)
			}
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (t *Tacx) startSerialLoop() {
	config := t.config
	port, err := connect(config.SerialDevice)

	if err != nil {
		log.Fatalf("unable to connect to tacx: %v", err)
	}

	commander := makeCommander(port)

	_, err = getVersion(commander)
	if err != nil {
		log.Fatalf("unable to retrieve tacx version: %v", err)
	}

	var lowestWeight = uint8(0x0a)
	var calibrationSpeed = int16(getRawSpeed(float64(config.CalibrationSpeed)))
	var calibrationDurationMin = time.Duration(config.CalibrationMin) * time.Second
	var calibrationDurationMax = time.Duration(config.CalibrationMax) * time.Second

	var controlCommandsPerSecond = 5
	if config.Slow {
		controlCommandsPerSecond = 1
	}

	var calibrating = config.Calibrate
	var calibrationStartedAt time.Time
	var calibrationLastLoads = make([]float64, 0, controlCommandsPerSecond*10)
	var calibrationTolerance = float64(config.CalibrationTolerance)
	var calibrationResult uint16 = 1040 // a sensible default, in case calibration is disabled

	lastResponse := controlResponse{}
	for {
		state := t.getState()
		startTime := time.Now()

		command := controlCommand{
			keepalive: lastResponse.keepalive,
			adjust:    calibrationResult,
		}

		logLine := "mode"
		if calibrating {
			command.mode = modeCalibrating
			command.targetSpeed = calibrationSpeed
			command.weight = lowestWeight
			log.WithFields(log.Fields{
				"mode": "calibrating",
			}).Debug(logLine)
		} else if !state.Enabled {
			command.mode = modeOff
			log.WithFields(log.Fields{
				"mode": "off",
			}).Debug(logLine)
		} else {
			command.mode = modeNormal
			switch state.Behaviour {
			case BehaviourERG:
				command.weight = lowestWeight
				command.targetLoad = getTargetLoad(targetLoadArgs{
					targetWatts:  state.TargetWatts,
					currentSpeed: lastResponse.speed,
				})
				log.WithFields(log.Fields{
					"mode":      "normal",
					"behaviour": "erg",
				}).Debug(logLine)
			case BehaviourSimulator:
				command.weight = config.Weight
				targetWattsForSimulator := getWattsForSimulator(targetLoadForSimulatorArgs{
					currentSpeed:      lastResponse.speed,
					weight:            config.Weight,
					windSpeed:         state.WindSpeed,
					gradient:          state.Gradient,
					rollingResistance: state.RollingResistance,
					windResistance:    state.WindResistance,
				})
				command.targetLoad = getTargetLoad(targetLoadArgs{
					targetWatts:  targetWattsForSimulator,
					currentSpeed: lastResponse.speed,
				})
				log.WithFields(log.Fields{
					"mode":      "normal",
					"behaviour": "sim",
				}).Debug(logLine)
			}
		}

		controlResponse, err := sendControl(commander, command)
		if err != nil {
			log.Errorf("unable to execute main command: %+v", err)
			// allow this to occasionally fail
			// TODO: count failures and exit after reaching limit
			continue
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
					log.Info("waiting for calibration: pedal once then stop")
				}
			} else {
				untilMinimum := time.Until(calibrationStartedAt.Add(calibrationDurationMin))
				untilMaximum := time.Until(calibrationStartedAt.Add(calibrationDurationMax))

				if len(calibrationLastLoads) == cap(calibrationLastLoads) {
					calibrationLastLoads = calibrationLastLoads[1:]
				}
				calibrationLastLoads = append(calibrationLastLoads, float64(controlResponse.currentLoad))

				if untilMinimum < 0 {
					quartile1, _ := stats.Percentile(calibrationLastLoads, 25)
					quartile3, _ := stats.Percentile(calibrationLastLoads, 75)
					stable := quartile3-quartile1 <= calibrationTolerance
					average := quartile1 + ((quartile3 - quartile1) / 2)
					log.WithFields(log.Fields{
						"stable":    fmt.Sprintf("%t", stable),
						"quartile1": fmt.Sprintf("%.2f", quartile1),
						"quartile3": fmt.Sprintf("%.2f", quartile3),
						"average":   fmt.Sprintf("%.2f", average),
					}).Info("calibrated?")

					if stable || untilMaximum > 0 {
						calibrationResult = uint16(average)
						calibrating = false
						if !stable {
							log.Error("calibration aborted: maximum time reached. using last average")
						}
					}
				}
				log.WithFields(log.Fields{
					"minimum":   fmt.Sprintf("%t", untilMinimum < 0),
					"remaining": fmt.Sprintf("%.0f", untilMinimum.Seconds()),
					"speed":     fmt.Sprintf("%v", controlResponse.speed),
					"load":      fmt.Sprintf("%v", controlResponse.currentLoad),
				}).Info("calibrating")
			}
		}

		t.channel <- TacxEvent{
			Ready:   !calibrating,
			Speed:   getKilometers(controlResponse.speed),
			Load:    getWatts(controlResponse.currentLoad) * float64(controlResponse.speed),
			Cadence: controlResponse.cadence,
		}

		lastResponse = controlResponse

		period := time.Second / time.Duration(controlCommandsPerSecond)
		time.Sleep(time.Until(startTime.Add(period)))
	}
}
