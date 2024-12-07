package main

import (
	"flag"
	"log"

	"github.com/rcambrj/tacxble/bicipi"
)

func main() {
	serialDevice := flag.String("serial", "", "The serial device to which Tacx motorbrake is connected. Defaults to the first one found.")
	bluetoothDevice := flag.String("bluetooth", "", "The bluetooth device on which the FTMS will be advertised. Defaults to the first one found.")

	logger := log.Logger{}

	config := bicipi.StartConfig{
		SerialDevice:    *serialDevice,
		BluetoothDevice: *bluetoothDevice,
		Logger:          &logger,
	}

	bicipi.Start(config)
}
