package tacxble

import (
	"fmt"
	"strings"
	"time"

	"tinygo.org/x/bluetooth"
)

func unhandledWriteEvent(name string, offset int, value []byte) {
	fmt.Printf("Recieved %s\nOffset %d\n", name, offset)
	fmt.Println(formatBinary(value))
}

func formatBinary(bytes []byte) string {
	var output string
	for _, n := range bytes {
		output += fmt.Sprintf("%08b ", n) // format each byte to a binary octet
	}
	return strings.TrimSpace(output)
}

type WriteValue = func() []byte

func writeFakeData(
	name string,
	serviceManager *ServiceManager,
	serviceUUID bluetooth.UUID,
	characteristicUUID bluetooth.UUID,
	getValue WriteValue,
) chan bool {
	characteristic, err := serviceManager.GetCharacteristic(
		serviceUUID,
		characteristicUUID,
	)
	if err != nil {
		panic(err.Error())
	}

	nextBeat := time.Now()

	exitSignal := make(chan bool)

	go func() {
		for {
			select {
			case <-exitSignal:
				fmt.Println(fmt.Sprintf("exiting %s", name))
				return
			default:
			}

			nextBeat = nextBeat.Add(time.Minute / time.Duration(60))
			time.Sleep(nextBeat.Sub(time.Now()))

			value := getValue()
			fmt.Printf(
				"%s at %s: %X\n",
				name,
				time.Now().Format(time.RFC3339Nano),
				value,
			)
			characteristic.Write(value)
		}
	}()

	return exitSignal
}
