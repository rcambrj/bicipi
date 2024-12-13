package ftms

import (
	"fmt"
	"sync"

	"tinygo.org/x/bluetooth"
)

type Config struct {
	BluetoothName string
}

type State struct {
}

type FTMSEvent struct {
}

type Listener = func(event FTMSEvent)

func MakeService(config Config) FTMS {
	return FTMS{
		config:  config,
		channel: make(chan FTMSEvent),
	}
}

type FTMS struct {
	config    Config
	stateLock sync.Mutex
	state     State
	channel   chan FTMSEvent
	listeners []Listener

	connected      bool
	serviceManager ServiceManager
}

func (f *FTMS) SetState(state State) {
	f.stateLock.Lock()
	defer f.stateLock.Unlock()

	// t.state.foo = state.foo
}

func (f *FTMS) getState() State {
	f.stateLock.Lock()
	defer f.stateLock.Unlock()

	return f.state
}

func (f *FTMS) On(listener Listener) {
	f.listeners = append(f.listeners, listener)
}

func (f *FTMS) requestControl() error {
	// Do we even need this?
	if !f.connected {
		return fmt.Errorf("not connected")
	}

	char, err := f.serviceManager.GetCharacteristic(bluetooth.ServiceUUIDFitnessMachine, bluetooth.CharacteristicUUIDFitnessMachineControlPoint)
	if err != nil {
		return fmt.Errorf("unable to get characteristic: %w", err)
	}

	fmt.Printf("%v", char)
	// isLegit := char.Write(getFTMSMode(mode))
	return nil
}

func (f *FTMS) Start() {
	go f.startBLELoop()
}

func (f *FTMS) startBLELoop() {
	adapter := bluetooth.DefaultAdapter

	must("enable BLE stack", adapter.Enable())

	serviceManager := NewServiceManager()

	registerServices(&serviceManager)

	serviceUUIDs := serviceManager.GetServiceIds()
	must("register services", serviceManager.RegisterServices(adapter))

	adv := adapter.DefaultAdvertisement()
	must("configure advertisement", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName:    f.config.BluetoothName,
		ServiceUUIDs: serviceUUIDs,
	}))

	for _, svcUUID := range serviceUUIDs {
		fmt.Printf("- %s\n", svcUUID)
	}

	must("start advertising BLE", adv.Start())

	println("advertising BLE...")

	WriteFakeData(
		"HeartRate",
		&serviceManager,
		bluetooth.ServiceUUIDHeartRate,
		bluetooth.CharacteristicUUIDHeartRateMeasurement,
		HeartRateDataGenerator(),
	)

	WriteFakeData(
		"Cadence",
		&serviceManager,
		bluetooth.ServiceUUIDCyclingSpeedAndCadence,
		bluetooth.CharacteristicUUIDCSCMeasurement,
		CadenceDataGenerator(),
	)
}

func registerServices(serviceManager *ServiceManager) {
	must("declare HeartRate service", serviceManager.AddService(
		bluetooth.ServiceUUIDHeartRate,
		CreateHeartRateCharacteristics()...,
	))

	must("declare FTMS service", serviceManager.AddService(
		bluetooth.ServiceUUIDFitnessMachine,
		CreateFitnessMachineCharacteristics()...,
	))

	must("declare Cycling Power service", serviceManager.AddService(
		bluetooth.ServiceUUIDCyclingPower,
		CreateCyclingPowerCharacteristics()...,
	))

	must("declare Cycling Speed and Cadence service", serviceManager.AddService(
		bluetooth.ServiceUUIDCyclingSpeedAndCadence,
		CreateCyclingSpeedCadenceCharacteristics()...,
	))

	must("declare Cycling Steering service", serviceManager.AddService(
		ServiceUUIDCyclingSteering,
		CreateCyclingSteeringCharacteristics()...,
	))
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
