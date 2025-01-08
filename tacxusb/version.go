package tacxusb

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/rcambrj/bicipi/tacxcommon"
	log "github.com/sirupsen/logrus"
)

func parseVersionResponseBytes(response []byte) (versionResponseRaw, error) {
	buf := bytes.NewReader(response)
	out := versionResponseRaw{}
	if err := binary.Read(buf, binary.LittleEndian, &out); err != nil {
		return versionResponseRaw{}, err
	}
	return out, nil
}

type versionResponseRaw struct {
	_         uint32
	Firmware0 uint8
	Firmware1 uint8
	Firmware2 uint8
	Firmware3 uint8
	Serial    uint32
	DateD     uint8
	DateM     uint8
}

func getVersion(t commander) (tacxcommon.Version, error) {
	log.Info("requesting tacx version...")
	response, err := t.sendCommand([]byte{0x02, 0x00, 0x00, 0x00})
	if err != nil {
		return tacxcommon.Version{}, fmt.Errorf("unable to get version: %w", err)
	}

	// head unit
	// 18 204 18 0 0 2 0 0 9 0 0 0 0 0 0 0 0 0 13 10 0 0 0 0
	// motor brake
	// 0 0 0 0 0 0 0 0 7 5 2 2 64 0 0 0 0 0 0 0 0 0 0 0
	// padding
	// 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0

	versionRaw, err := parseVersionResponseBytes(response)
	if err != nil {
		return tacxcommon.Version{}, fmt.Errorf("unable to parse version: %w", err)
	}

	firmwareVersion := fmt.Sprintf("%02X.%02X.%02X.%02X",
		versionRaw.Firmware3,
		versionRaw.Firmware2,
		versionRaw.Firmware1,
		versionRaw.Firmware0,
	)
	serial := int32(versionRaw.Serial)
	model := fmt.Sprintf("T19%v", int(serial/10000000))
	manufactureNumber := int(serial % 100000)
	year := 2000 + int(serial/100000%100)
	date := fmt.Sprintf("%v-%02X-%02X",
		year,
		versionRaw.DateM,
		versionRaw.DateD,
	)

	version := tacxcommon.Version{
		Model:             model,
		ManufactureNumber: manufactureNumber,
		FirmwareVersion:   firmwareVersion,
		Serial:            serial,
		Date:              date,
	}
	log.WithFields(log.Fields{"version": fmt.Sprintf("%+v", version)}).Info("received tacx version")
	return version, nil
}
