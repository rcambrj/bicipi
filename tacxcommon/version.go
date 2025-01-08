package tacxcommon

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func GetVersionCommand() []byte {
	return []byte{0x02, 0x00, 0x00, 0x00}
}

type Version struct {
	Model             string
	ManufactureNumber int
	FirmwareVersion   string
	Serial            int32
	Date              string
}

func parseVersionResponseBytes(response []byte) (versionResponseRaw, error) {
	buf := bytes.NewReader(response)
	out := versionResponseRaw{}
	if err := binary.Read(buf, binary.LittleEndian, &out); err != nil {
		return versionResponseRaw{}, err
	}
	return out, nil
}

type versionResponseRaw struct {
	Firmware0 uint8
	Firmware1 uint8
	Firmware2 uint8
	Firmware3 uint8
	Serial    uint32
	DateD     uint8
	DateM     uint8
}

func GetVersionFromResponseBytes(response []byte) (Version, error) {
	versionRaw, err := parseVersionResponseBytes(response)
	if err != nil {
		return Version{}, fmt.Errorf("unable to parse version: %w", err)
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

	version := Version{
		Model:             model,
		ManufactureNumber: manufactureNumber,
		FirmwareVersion:   firmwareVersion,
		Serial:            serial,
		Date:              date,
	}

	return version, nil
}
