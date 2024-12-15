package ftms

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"tinygo.org/x/bluetooth"
)

type Config struct {
	BluetoothName   string
	BluetoothDevice string
}

type State struct {
	Speed   uint16
	Load    int16
	Cadence uint16
}

type FTMSEvent struct {
	Mode Mode

	// ModeTargetPower
	TargetPower int16

	// ModeIndoorBikeSimulation
	WindSpeed         float64
	TargetGrade       float64
	RollingResistance float64
	WindResistance    float64
}

type Listener = func(event FTMSEvent)

func MakeService(config Config) FTMS {
	return FTMS{
		config:         config,
		channel:        make(chan FTMSEvent),
		serviceManager: NewServiceManager(),
	}
}

type FTMS struct {
	config         Config
	channel        chan FTMSEvent
	listeners      []Listener
	serviceManager ServiceManager
}

func (f *FTMS) SetState(state State) {
	err := writeFMIndoorBikeData(&f.serviceManager, state.Speed, state.Cadence, state.Load)
	if err != nil {
		log.Fatalf("unable to write to characteristic: %v", err)
	}
}

func (f *FTMS) On(listener Listener) {
	f.listeners = append(f.listeners, listener)
}

func (f *FTMS) Start() {
	go f.startEventLoop()
	f.startBLE()
}

func (f *FTMS) startEventLoop() {
	for {
		select {
		case msg := <-f.channel:
			for _, listener := range f.listeners {
				listener(msg)
			}
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (f *FTMS) startBLE() {
	adapter := bluetooth.DefaultAdapter
	// TODO: bluetooth.Adapter.id is private, so can't use non-default adapter

	err := adapter.Enable()
	if err != nil {
		log.Fatalf("unable to use bluetooth adapter: %v", err)
	}
	err = f.registerServices(&f.serviceManager)
	if err != nil {
		log.Fatalf("unable to register ble services: %v", err)
	}
	err = f.serviceManager.PublishServices(adapter)
	if err != nil {
		log.Fatalf("unable to publish ble services: %v", err)
	}
	adv := adapter.DefaultAdvertisement()
	serviceUUIDs := f.serviceManager.GetServiceIds()
	err = adv.Configure(bluetooth.AdvertisementOptions{
		LocalName:    f.config.BluetoothName,
		ServiceUUIDs: serviceUUIDs,
	})
	if err != nil {
		log.Fatalf("unable to configure ble advertisement: %v", err)
	}
	adapter.SetConnectHandler(handleConnect)
	err = adv.Start()
	if err != nil {
		log.Fatalf("unable to start advertising ble: %v", err)
	}
	log.Info("advertising ble")
}

func handleConnect(device bluetooth.Device, connected bool) {
	log.Warnf("handleConnect called: %v", device.Address.MACAddress.MAC)
}

func (f *FTMS) registerServices(serviceManager *ServiceManager) error {
	err := serviceManager.AddService(
		bluetooth.ServiceUUIDFitnessMachine,
		CreateFitnessMachineCharacteristics(f.receiveFTMSOperation)...,
	)
	if err != nil {
		return fmt.Errorf("unable to add ftms service: %w", err)
	}

	err = serviceManager.AddService(
		bluetooth.ServiceUUIDCyclingPower,
		CreateCyclingPowerCharacteristics(f.receiveCyclingPowerOperation)...,
	)
	if err != nil {
		return fmt.Errorf("unable to add cycling power service: %w", err)
	}

	err = serviceManager.AddService(
		bluetooth.ServiceUUIDCyclingSpeedAndCadence,
		CreateCyclingSpeedCadenceCharacteristics()...,
	)
	if err != nil {
		return fmt.Errorf("unable to add cycling speed and cadence service: %w", err)
	}

	return nil
}

func (f *FTMS) receiveFTMSOperation(client bluetooth.Connection, offset int, value []byte) {
	log.WithFields(log.Fields{
		"operation": value,
	}).Trace("received FTMS operation")

	switch value[0] {
	case FMCPOpCodeResponseCode:
		return
	case FMCPOpCodeRequestControl:
		err := writeFMCPResultCode(&f.serviceManager, FMCPOpCodeRequestControl, FMCPResultCodeSuccess)
		if err != nil {
			log.Fatalf("unable to accept control: %v", err)
		}
	case FMCPOpCodeSetTargetPower:
		targetPowerCommand, err := readFMCPTargetPower(value)
		if err != nil {
			log.Fatalf("unable to read ble SetTargetPower command: %v", err)
		}
		log.Infof("ble command SetTargetPower to %v", targetPowerCommand.TargetPower)

		f.channel <- FTMSEvent{
			Mode:        ModeTargetPower,
			TargetPower: getWatts(targetPowerCommand.TargetPower),
		}

		err = writeFMSTargetPower(&f.serviceManager, targetPowerCommand.TargetPower)
		if err != nil {
			log.Errorf("unable to SetTargetPower on FMS: %v", err)
		}
		err = writeFMCPResultCode(&f.serviceManager, FMCPOpCodeSetTargetPower, FMCPResultCodeSuccess)
		if err != nil {
			log.Errorf("unable to SetTargetPower on FMCP: %v", err)
		}
	case FMCPOpCodeSetIndoorBikeSimulation:
		f.channel <- FTMSEvent{
			Mode:              ModeIndoorBikeSimulation,
			WindSpeed:         float64(value[1]) * 0.001,
			TargetGrade:       float64(value[2]) * 0.01,
			RollingResistance: float64(value[3]) * 0.0001,
			WindResistance:    float64(value[4]) * 0.01,
		}
	default:
		log.WithFields(log.Fields{
			"offset": offset,
			"value":  value,
		}).Fatalf("FTMS operation opcode not implemented")
	}
}

func (f *FTMS) receiveCyclingPowerOperation(client bluetooth.Connection, offset int, value []byte) {
	log.WithFields(log.Fields{
		"offset": offset,
		"value":  value,
	}).Fatalf("Cycling Power operation opcode not implemented")
}
