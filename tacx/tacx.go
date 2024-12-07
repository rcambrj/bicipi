package tacx

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

type Config struct {
	Device string
}

func Start(config Config) {
	var device string

	if config.Device != "" {
		device = config.Device
	} else {
		log.Info("searching for serial ports...")
		devices, err := serial.GetPortsList()
		if err != nil {
			log.Fatalf("unable to list serial ports: %v", err)
		}
		if len(devices) == 0 {
			log.Fatal("no serial ports found")
		}
		device = devices[0]
	}
	log.Infof("connecting to serial port %v...", device)

	mode := &serial.Mode{
		BaudRate: 19200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open(device, mode)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to open serial port: %w", err))
	}

	command := []byte{0x02, 0x00, 0x00, 0x00}
	log.Debugf("sending serial command: %v", command)
	outFrame, err := serializeCommand(command)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to serialize command: %w", err))
	}

	readChan := make(chan []byte)

	// start reading before sending the first command
	go read(readChan, port)

	port.ResetInputBuffer()

	_, err = port.Write(outFrame)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to write to serial port: %w", err))
	}
	log.Debugf("sent serial data: %v", outFrame)

	inFrame := waitForResponse(readChan, port)
	response, err := deserializeResponse(inFrame)
	if err != nil {
		log.Warnf("unable to deserialize response: %v", err)
	} else {
		log.Debugf("received serial response: %v", response)
	}
}
