package tacxble

import (
	"fmt"
	"math/rand"
	"time"

	"tinygo.org/x/bluetooth"
)

func Start() {
	fmt.Println("starting...")

	adapter := bluetooth.DefaultAdapter

	must("enable BLE stack", adapter.Enable())

	serviceUUIDs := []bluetooth.UUID{}

	ftmService := getFitnessMachineServiceDefinition()
	must("declare FTMS service", adapter.AddService(&ftmService))
	serviceUUIDs = append(serviceUUIDs, ftmService.UUID)

	cyclingPowerService := getCyclingPowerServiceDefinition()
	must("declare Cycling service", adapter.AddService(&cyclingPowerService))
	serviceUUIDs = append(serviceUUIDs, cyclingPowerService.UUID)

	cyclingSpeedAndCadenceService := getCyclingSpeedAndCadenceServiceDefinition()
	must("declare Cycling service", adapter.AddService(&cyclingSpeedAndCadenceService))
	serviceUUIDs = append(serviceUUIDs, cyclingSpeedAndCadenceService.UUID)

	cyclingSteeringService := getVirtualSteeringService()
	must("declare Cycling service", adapter.AddService(&cyclingSteeringService))
	serviceUUIDs = append(serviceUUIDs, cyclingSteeringService.UUID)

	var heartRate bluetooth.Characteristic
	heartRateService := getHeartRateService(&heartRate)
	must("declare Cycling service", adapter.AddService(&heartRateService))
	serviceUUIDs = append(serviceUUIDs, heartRateService.UUID)

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

	startFakeHeartbeat(&heartRate)

	for {
		// Sleep forever.
		time.Sleep(time.Hour)
	}
}

func startFakeHeartbeat(heartRate *bluetooth.Characteristic) chan bool {
	var currentRate uint8 = 60
	nextBeat := time.Now()
	rateFluctuation := 10

	exitSignal := make(chan bool)

	go func() {
		for {
			select {
			case <-exitSignal:
				fmt.Println("exiting heartbeat")
				return
			default:
			}

			nextBeat = nextBeat.Add(time.Minute / time.Duration(currentRate))
			fmt.Printf("heartbeat at %s\n", time.Now().Format(time.RFC3339Nano))
			time.Sleep(nextBeat.Sub(time.Now()))

			rateOffset := rateFluctuation/2 - 1 - rand.Intn(rateFluctuation)
			currentRate = uint8(min(max(int(currentRate)+rateOffset, 55), 110))

			rateData := getHeartRate(currentRate)
			heartRate.Write(rateData)
		}
	}()

	return exitSignal
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}

func handleConnect(device bluetooth.Device, connected bool) {
	fmt.Println("received connection")
}
