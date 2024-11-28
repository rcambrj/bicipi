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

	ftmService := getFitnessMachineServiceDefinition()
	must("declare FTMS service", adapter.AddService(&ftmService))

	cyclingPowerService := getCyclingPowerServiceDefinition()
	must("declare Cycling service", adapter.AddService(&cyclingPowerService))

	cyclingSpeedAndCadenceService := getCyclingSpeedAndCadenceServiceDefinition()
	must("declare Cycling service", adapter.AddService(&cyclingSpeedAndCadenceService))

	adv := adapter.DefaultAdvertisement()
	must("configure advertisement", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: "Tacx BLE Trainer",
	}))

	adapter.SetConnectHandler(handleConnect)

	must("start advertising BLE", adv.Start())

	println("advertising BLE...")

	for {
		// Sleep forever.
		time.Sleep(time.Hour)
	}
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}

func handleConnect(device bluetooth.Device, connected bool) {
	fmt.Println("received connection")
}
