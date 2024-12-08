package tacx

import (
	"fmt"
	"time"

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
	err = port.SetReadTimeout(100 * time.Millisecond)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to configure serial timeout: %w", err))
	}

	commander := makeCommander(port)

	version, err := getVersion(commander)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("done %+v", version)
}
