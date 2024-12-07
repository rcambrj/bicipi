package main

import (
	"flag"
	"fmt"

	"github.com/rcambrj/tacxble/bicipi"
	log "github.com/sirupsen/logrus"
)

func main() {
	logLevels := []string{"debug", "info", "warn", "error"}

	serialDevice := flag.String("serial", "", "The serial device to which Tacx motorbrake is connected. Defaults to the first one found.")
	bluetoothDevice := flag.String("bluetooth", "", "The bluetooth device on which the FTMS will be advertised. Defaults to the first one found.")
	logLevel := flag.String("loglevel", "info", fmt.Sprintf("The log level. May be one of %v", logLevels))
	flag.Parse()

	validLogLevel, err := log.ParseLevel(*logLevel)
	if err != nil {
		panic("invalid log level")
	}

	log.SetLevel(validLogLevel)

	config := bicipi.Config{
		SerialDevice:    *serialDevice,
		BluetoothDevice: *bluetoothDevice,
	}

	bicipi.Start(config)
}
