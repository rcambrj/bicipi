package tacx

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type version struct {
	model             string
	manufactureYear   int
	manufactureNumber int
	firmwareVersion   string
	serial            int32
	date              string
}

func getVersion(t Commander) (version, error) {
	log.Info("requesting tacx version...")
	response, err := t.sendCommand([]byte{0x02, 0x00, 0x00, 0x00})
	if err != nil {
		return version{}, fmt.Errorf("unable to get version: %w", err)
	}

	firmwareVersion := fmt.Sprintf("%02X.%02X.%02X.%02X", response[7], response[6], response[5], response[4])
	serial := int32(response[8]) | int32(response[9])<<8 | int32(response[10])<<16 | int32(response[11])<<24
	date := fmt.Sprintf("%02X-%02X", response[13], response[12])
	// serial-based properties
	manufactureYear := 2000 + int(serial/100000%100)
	manufactureNumber := int(serial % 100000)
	model := fmt.Sprintf("T19%v", int(serial/10000000))

	version := version{
		model:             model,
		manufactureYear:   manufactureYear,
		manufactureNumber: manufactureNumber,
		firmwareVersion:   firmwareVersion,
		serial:            serial,
		date:              date,
	}
	log.WithFields(log.Fields{"version": fmt.Sprintf("%+v", version)}).Infof("received tacx version")
	return version, nil
}
