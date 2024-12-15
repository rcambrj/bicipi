package tacx

import (
	"sync"
	"time"

	"github.com/montanaflynn/stats"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Device               string
	Calibrate            bool
	Slow                 bool
	CalibrationSpeed     int
	CalibrationMin       int
	CalibrationMax       int
	CalibrationTolerance int
}

type State struct {
	Enabled        bool
	Behaviour      Behaviour
	TargetWatts    float64
	Weight         int
	WindSpeed      int
	DraftingFactor int
	Gradient       int
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

	t.state.Enabled = state.Enabled
	t.state.Behaviour = state.Behaviour
	t.state.TargetWatts = state.TargetWatts
	t.state.Weight = state.Weight
	t.state.WindSpeed = state.WindSpeed
	t.state.DraftingFactor = state.DraftingFactor
	t.state.Gradient = state.Gradient
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

		if calibrating {
			command.mode = modeCalibrating
			command.targetSpeed = calibrationSpeed
			command.weight = lowestWeight
			log.Infof("mode: calibrating")
		} else if !state.Enabled {
			command.mode = modeOff
			log.Infof("mode: off")
		} else {
			command.mode = modeNormal
			switch state.Behaviour {
			case BehaviourERG:
				command.weight = lowestWeight
				command.targetLoad = getTargetLoad(targetLoadArgs{
					targetWatts:  state.TargetWatts,
					currentSpeed: lastResponse.speed,
				})
				log.Infof("mode: normal; behaviour: erg; watts: %v; speed: %v; target %v", state.TargetWatts, lastResponse.speed, command.targetLoad)
			case BehaviourSimulator:
				command.weight = uint8(state.Weight)
				targetWattsForSimulator := getWattsForSimulator(targetLoadForSimulatorArgs{
					currentSpeed:   lastResponse.speed,
					weight:         state.Weight,
					windSpeed:      state.WindSpeed,
					draftingFactor: state.DraftingFactor,
					gradient:       state.Gradient,
				})
				command.targetLoad = getTargetLoad(targetLoadArgs{
					targetWatts:  targetWattsForSimulator,
					currentSpeed: lastResponse.speed,
				})
				log.Infof("mode: normal; behaviour: simulator; gradient: %v; speed: %v; target %v", state.Gradient, lastResponse.speed, command.targetLoad)
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
					log.Infof("waiting for calibration: pedal once then stop")
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
						log.Infof("calculating. stable: %t; quartile1: %.2f; quartile3: %.2f; average: %.2f", stable, quartile1, quartile3, average)

						if stable {
							calibrationResult = uint16(average)
							calibrating = false
						}
					}
					log.Infof("calibrating. minimum: %t; remaining: %.0f; speed: %v; resistance: %v;", untilMinimum < 0, untilMinimum.Seconds(), controlResponse.speed, controlResponse.currentLoad)
				} else {
					calibrating = false
					log.Fatal("calibration aborted: maximum time reached.")
				}
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
