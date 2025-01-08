package tacxcommon

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Mode uint8

const (
	ModeOff Mode = iota
	_            // 1 is unused
	ModeNormal
	ModeCalibrating
)

type ControlCommand struct {
	Mode        Mode   // see const `Mode`
	TargetSpeed int16  // modeCalibrating: raw speed
	TargetLoad  int16  // modeNormal: raw load
	Keepalive   uint8  // value from the last response
	Weight      uint8  // kg
	Adjust      uint16 // adjustment resulting from calibration
}

func GetControlCommandBytes(command ControlCommand) ([]byte, error) {
	var target int16
	switch command.Mode {
	case ModeCalibrating:
		target = command.TargetSpeed
	case ModeNormal:
		target = command.TargetLoad
	}

	buf := new(bytes.Buffer)
	parts := []any{
		uint8(0x01),
		uint8(0x08),
		uint8(0x01),
		uint8(0x00),
		target,
		command.Keepalive,
		uint8(0x00),
		uint8(command.Mode),
		command.Weight,
		command.Adjust,
	}
	for i, v := range parts {
		err := binary.Write(buf, binary.LittleEndian, v)
		if err != nil {
			return []byte{}, fmt.Errorf("unable to write part %v, %v: %w", i, v, err)
		}
	}
	return buf.Bytes(), nil
}

type controlResponseRaw struct {
	_           uint32
	Distance    uint32
	Speed       uint16
	_           uint16
	AverageLoad int16
	CurrentLoad int16
	TargetLoad  int16
	KeepAlive   uint8
	_           uint8
	Cadence     uint8
}

func GetControlResponseFromBytes(response []byte) (ControlResponse, error) {
	buf := bytes.NewReader(response)
	responseRaw := controlResponseRaw{}
	if err := binary.Read(buf, binary.LittleEndian, &responseRaw); err != nil {
		return ControlResponse{}, err
	}

	return ControlResponse{
		Speed:       responseRaw.Speed,
		CurrentLoad: responseRaw.CurrentLoad,
		TargetLoad:  responseRaw.TargetLoad,
		Keepalive:   responseRaw.KeepAlive,
		Cadence:     responseRaw.Cadence,
	}, nil
}

// ControlResponse could contain more information if the connection is USB.
// in order to reduce maintenance, ControlResponse contains only what's
// supported on all connection protocols (ie. what comes from the motor brake)
// If the extra USB data is necessary, then ControlResponse should not live here
// as the responses for serial vs USB will differ
type ControlResponse struct {
	Speed       uint16 // tacx speed units
	CurrentLoad int16  // tacx load units
	TargetLoad  int16  // tacx load units
	Keepalive   uint8  // value to send in the next control
	Cadence     uint8  // rpm
}
