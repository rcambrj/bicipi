package tacx

import (
	"fmt"
	"log"

	"go.bug.st/serial"
)

func Start() {
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

	command, err := SerializeCommand([]byte{0x02, 0x00, 0x00, 0x00})
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
		response, err := DeserializeResponse(buff)
		if err != nil {
			fmt.Printf("unable to deserialize response: %v", err)
		} else {
			fmt.Printf("%v", response)
		}
	}

}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
