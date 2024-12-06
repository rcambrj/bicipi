package tacxble

import (
	"fmt"
	"log"
	"time"

	"go.bug.st/serial"
	"tinygo.org/x/bluetooth"
)

func Start() {
	fmt.Println("starting...")

	ports, err := serial.GetPortsList()
	must("get ports list", err)
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}
	for _, port := range ports {
		fmt.Printf("Found port: %v\n", port)
	}
	if len(ports) > 1 {
		log.Fatal("Found more than one port. TODO: allow specifying port on cli")
	}

	mode := &serial.Mode{
		BaudRate: 19200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open(ports[0], mode)
	if err != nil {
		log.Fatal(err)
	}

	n, err := port.Write([]byte{0x02, 0x00, 0x00, 0x00})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)

	buff := make([]byte, 1)
	for {
		n, err := port.Read(buff)
		if err != nil {
			log.Fatal(err)
			break
		}
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}
		fmt.Printf("%v", string(buff[:n]))
	}

	// //////////////////////////

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
