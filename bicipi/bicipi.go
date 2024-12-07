package bicipi

import (
	"fmt"
	"log"
	"time"

	"github.com/rcambrj/tacxble/ftms"
	"github.com/rcambrj/tacxble/tacx"
)

type StartConfig struct {
	SerialDevice    string
	BluetoothDevice string
	Logger          *log.Logger
}

func Start(config StartConfig) {
	fmt.Println("starting...")

	tacx.Start()
	// TODO: wait for tacx to be ready then advertise FTMS
	ftms.Start()

	for {
		// Sleep forever.
		time.Sleep(time.Hour)
	}
}
