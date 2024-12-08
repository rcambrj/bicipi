package tacx

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Version struct {
	Model             string
	ManufactureYear   int
	ManufactureNumber int
	FirmwareVersion   string
	Serial            int32
	Date              string
	Other             string
}

func getVersion(t Commander) (Version, error) {
	log.Info("requesting tacx version...")
	response, err := t.sendCommand([]byte{0x02, 0x00, 0x00, 0x00})
	if err != nil {
		return Version{}, fmt.Errorf("unable to get version: %w", err)
	}

	firmwareVersion := fmt.Sprintf("%02X.%02X.%02X.%02X", response[7], response[6], response[5], response[4])
	date := fmt.Sprintf("%02X-%02X", response[13], response[12])
	other := fmt.Sprintf("%02X.%02X", response[15], response[14])
	serial := int32(response[8]) | int32(response[9])<<8 | int32(response[10])<<16 | int32(response[11])<<24
	// serial-based properties
	manufactureYear := 2000 + int(serial/100000%100)
	manufactureNumber := int(serial % 100000)
	model := fmt.Sprintf("T19%v", int(serial/10000000))

	version := Version{
		Model:             model,
		ManufactureYear:   manufactureYear,
		ManufactureNumber: manufactureNumber,
		FirmwareVersion:   firmwareVersion,
		Serial:            serial,
		Date:              date,
		Other:             other,
	}
	log.Infof("received tacx version: %+v", version)
	return version, nil
}
