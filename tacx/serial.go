package tacx

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

func connect(device string) (serial.Port, error) {
	if device == "" {
		log.Info("searching for serial ports...")
		devices, err := serial.GetPortsList()
		if err != nil {
			return nil, fmt.Errorf("unable to list serial ports: %w", err)
		}
		if len(devices) == 0 {
			return nil, fmt.Errorf("no serial ports found")
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
		return nil, fmt.Errorf("unable to open serial port: %w", err)
	}

	// timeout doesn't affect how quickly data will be received.
	// port.Read() will return based on some internal trigger once some data is
	// received. this timeout only affects how quickly port.Read() will return
	// when there is no data being received. this shouldn't happen under normal
	// operation because port.Read() should not be called again once a valid
	// frame has been identified (start of frame byte ... end of frame byte)
	err = port.SetReadTimeout(10 * time.Millisecond)
	if err != nil {
		return nil, fmt.Errorf("unable to configure serial timeout: %w", err)
	}

	log.Infof("connected to serial port %v", device)

	return port, nil
}
