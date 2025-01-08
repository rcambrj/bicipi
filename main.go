package main

import (
	"flag"
	"fmt"

	"github.com/rcambrj/bicipi/bicipi"
	log "github.com/sirupsen/logrus"
)

func main() {
	logLevels := []string{"trace", "debug", "info", "warn", "error"}

	weight := flag.Uint("weight", 80, "The approximate weight of the rider + bicycle, used only in simulator mode (Zwift / MyWhoosh).")
	serialDevice := flag.String("serial", "", "The serial device to which Tacx motorbrake is connected. (default is to use USB)")
	bluetoothName := flag.String("bluetooth-name", "bicipi", "The bluetooth device name to advertise")
	logLevel := flag.String("loglevel", "info", fmt.Sprintf("The log level. May be one of %v.", logLevels))
	calibrate := flag.Bool("calibrate", true, "Whether to enable initial calibration. (--calibrate=false to disable)")
	slow := flag.Bool("slow", false, "Whether to poll slowly so that logs are easier to follow.")
	calibrationSpeed := flag.Int("calibration-speed", 20, "How fast in km/h to spin the tyre during calibration.")
	calibrationMin := flag.Int("calibration-min", 300, "How long in seconds to warm up the motor and tyre during calibration.")
	calibrationMax := flag.Int("calibration-max", 600, "How long in seconds before calibration is abandoned.")
	calibrationTolerance := flag.Int("calibration-tolerance", 10, "How fussy to be when considering calibration complete. Lower is more fussy.")
	flag.Parse()

	validLogLevel, err := log.ParseLevel(*logLevel)
	if err != nil {
		panic("invalid log level")
	}
	log.SetLevel(validLogLevel)

	config := bicipi.Config{
		Weight:               uint8(*weight),
		SerialDevice:         *serialDevice,
		BluetoothName:        *bluetoothName,
		Calibrate:            *calibrate,
		Slow:                 *slow,
		CalibrationSpeed:     *calibrationSpeed,
		CalibrationMin:       *calibrationMin,
		CalibrationMax:       *calibrationMax,
		CalibrationTolerance: *calibrationTolerance,
	}

	bicipi.Start(config)
}
