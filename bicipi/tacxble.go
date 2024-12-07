package bicipi

import (
	"fmt"
	"log"
	"time"

	"go.bug.st/serial"
	"tinygo.org/x/bluetooth"

	"github.com/rcambrj/tacxble/ftms"
	"github.com/rcambrj/tacxble/tacx"
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

	command, err := tacx.SerializeCommand([]byte{0x02, 0x00, 0x00, 0x00})
	if err != nil {
		log.Fatal(err)
	}

	n, err := port.Write(command)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)

	buff := make([]byte, 64)
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
		fmt.Printf("%v\n", string(buff[:n]))
		response, err := tacx.DeserializeResponse(buff)
		if err != nil {
			fmt.Printf("unable to deserialize response: %v", err)
		} else {
			fmt.Printf("%v", response)
		}
	}

	// //////////////////////////

	adapter := bluetooth.DefaultAdapter

	must("enable BLE stack", adapter.Enable())

	serviceManager := ftms.NewServiceManager()

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

	ftms.WriteFakeData(
		"HeartRate",
		&serviceManager,
		bluetooth.ServiceUUIDHeartRate,
		bluetooth.CharacteristicUUIDHeartRateMeasurement,
		ftms.HeartRateDataGenerator(),
	)

	ftms.WriteFakeData(
		"Cadence",
		&serviceManager,
		bluetooth.ServiceUUIDCyclingSpeedAndCadence,
		bluetooth.CharacteristicUUIDCSCMeasurement,
		ftms.CadenceDataGenerator(),
	)

	for {
		// Sleep forever.
		time.Sleep(time.Hour)
	}
}

func registerServices(serviceManager *ftms.ServiceManager) {
	must("declare HeartRate service", serviceManager.AddService(
		bluetooth.ServiceUUIDHeartRate,
		ftms.CreateHeartRateCharacteristics()...,
	))

	must("declare FTMS service", serviceManager.AddService(
		bluetooth.ServiceUUIDFitnessMachine,
		ftms.CreateFitnessMachineCharacteristics()...,
	))

	must("declare Cycling Power service", serviceManager.AddService(
		bluetooth.ServiceUUIDCyclingPower,
		ftms.CreateCyclingPowerCharacteristics()...,
	))

	must("declare Cycling Speed and Cadence service", serviceManager.AddService(
		bluetooth.ServiceUUIDCyclingSpeedAndCadence,
		ftms.CreateCyclingSpeedCadenceCharacteristics()...,
	))

	must("declare Cycling Steering service", serviceManager.AddService(
		ftms.ServiceUUIDCyclingSteering,
		ftms.CreateCyclingSteeringCharacteristics()...,
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
