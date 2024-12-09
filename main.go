package main

import (
	"flag"
	"fmt"

	"github.com/rcambrj/tacxble/bicipi"
	log "github.com/sirupsen/logrus"
)

func main() {
	logLevels := []string{"trace", "debug", "info", "warn", "error"}

	serialDevice := flag.String("serial", "", "The serial device to which Tacx motorbrake is connected. (default is first one found)")
	bluetoothDevice := flag.String("bluetooth", "", "The bluetooth device on which the FTMS will be advertised. (default is first one found)")
	logLevel := flag.String("loglevel", "info", fmt.Sprintf("The log level. May be one of %v.", logLevels))
	calibrate := flag.Bool("calibrate", true, "Whether to enable initial calibration. Defaults to true.") // --calibrate=false
	flag.Parse()

	validLogLevel, err := log.ParseLevel(*logLevel)
	if err != nil {
		panic("invalid log level")
	}
	log.SetLevel(validLogLevel)

	config := bicipi.Config{
		SerialDevice:    *serialDevice,
		BluetoothDevice: *bluetoothDevice,
		Calibrate:       *calibrate,
	}

	bicipi.Start(config)
}
