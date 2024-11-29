package tacxble

import (
	"fmt"
	"time"

	"tinygo.org/x/bluetooth"
)

func Start() {
	fmt.Println("starting...")

	adapter := bluetooth.DefaultAdapter

	must("enable BLE stack", adapter.Enable())

	serviceManager := NewServiceManager()

	registerServices(&serviceManager)

	serviceUUIDs := serviceManager.GetServiceIds()
	must("register services", serviceManager.RegisterServices(adapter))

	adv := adapter.DefaultAdvertisement()
	must("configure advertisement", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName:    "Tacx BLE Trainer",
		ServiceUUIDs: serviceUUIDs,
	}))

	for _, svcUUID := range serviceUUIDs {
		fmt.Printf("- %s\n", svcUUID)
	}

	adapter.SetConnectHandler(handleConnect)

	must("start advertising BLE", adv.Start())

	println("advertising BLE...")

	writeFakeData(
		"HeartRate",
		&serviceManager,
		bluetooth.ServiceUUIDHeartRate,
		bluetooth.CharacteristicUUIDHeartRateMeasurement,
		heartRateDataGenerator(),
	)

	writeFakeData(
		"Cadence",
		&serviceManager,
		bluetooth.ServiceUUIDCyclingSpeedAndCadence,
		bluetooth.CharacteristicUUIDCSCMeasurement,
		cadenceDataGenerator(),
	)

	for {
		// Sleep forever.
		time.Sleep(time.Hour)
	}
}

func registerServices(serviceManager *ServiceManager) {
	must("declare HeartRate service", serviceManager.AddService(
		bluetooth.ServiceUUIDHeartRate,
		createHeartRateCharacteristics()...,
	))

	must("declare FTMS service", serviceManager.AddService(
		bluetooth.ServiceUUIDFitnessMachine,
		createFitnessMachineCharacteristics()...,
	))

	must("declare Cycling Power service", serviceManager.AddService(
		bluetooth.ServiceUUIDCyclingPower,
		createCyclingPowerCharacteristics()...,
	))

	must("declare Cycling Speed and Cadence service", serviceManager.AddService(
		bluetooth.ServiceUUIDCyclingSpeedAndCadence,
		createCyclingSpeedCadenceCharacteristics()...,
	))

	must("declare Cycling Steering service", serviceManager.AddService(
		ServiceUUIDCyclingSteering,
		createCyclingSteeringCharacteristics()...,
	))
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}

func handleConnect(device bluetooth.Device, connected bool) {
	fmt.Println("received connection")
}
