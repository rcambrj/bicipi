package ftms

import (
	"fmt"

	"tinygo.org/x/bluetooth"
)

func Start() {
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
