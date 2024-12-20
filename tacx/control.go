package tacx

import (
	"bytes"
	"encoding/binary"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type controlCommand struct {
	mode        mode   // see const `mode`
	targetSpeed int16  // modeCalibrating: raw speed
	targetLoad  int16  // modeNormal: raw load
	keepalive   uint8  // value from the last response
	weight      uint8  // kg
	adjust      uint16 // adjustment resulting from calibration
}

type controlCommandRaw struct {
	target    int16
	keepalive uint8
	mode      uint8
	weight    uint8
	adjust    uint16
}

func getControlCommandBytes(args controlCommandRaw) ([]byte, error) {
	buf := new(bytes.Buffer)
	parts := []any{
		uint8(0x01),
		uint8(0x08),
		uint8(0x01),
		uint8(0x00),
		args.target,
		args.keepalive,
		uint8(0x00),
		args.mode,
		args.weight,
		args.adjust,
	}
	for i, v := range parts {
		err := binary.Write(buf, binary.LittleEndian, v)
		if err != nil {
			return []byte{}, fmt.Errorf("unable to write part %v, %v: %w", i, v, err)
		}
	}
	return buf.Bytes(), nil
}

func parseControlResponseBytes(response []byte) (controlResponseRaw, error) {
	buf := bytes.NewReader(response)
	out := controlResponseRaw{}
	if err := binary.Read(buf, binary.LittleEndian, &out); err != nil {
		return controlResponseRaw{}, err
	}
	return out, nil
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

type controlResponse struct {
	speed       uint16 // tacx speed units
	currentLoad int16  // tacx load units
	targetLoad  int16  // tacx load units
	keepalive   uint8  // value to send in the next control
	cadence     uint8  // rpm
}

// this is the main function to send and receive data from tacx
// it both sends the target status and receives the reported status
func sendControl(t Commander, command controlCommand) (controlResponse, error) {
	log.WithFields(log.Fields{"command": fmt.Sprintf("%+v", command)}).Debugf("sending tacx status")

	var target int16
	switch command.mode {
	case modeCalibrating:
		target = command.targetSpeed
	case modeNormal:
		target = command.targetLoad
	}

	commandRaw := controlCommandRaw{
		target:    target,
		keepalive: command.keepalive,
		mode:      uint8(command.mode),
		weight:    command.weight,
		adjust:    command.adjust,
	}
	log.WithFields(log.Fields{"commandRaw": fmt.Sprintf("%+v", commandRaw)}).Tracef("sending tacx status raw")

	commandBytes, err := getControlCommandBytes(commandRaw)
	if err != nil {
		return controlResponse{}, fmt.Errorf("unable to process tacx control command: %w", err)
	}

	responseBytes, err := t.sendCommand(commandBytes)
	if err != nil {
		return controlResponse{}, fmt.Errorf("unable to send tacx control command: %w", err)
	}

	responseRaw, err := parseControlResponseBytes(responseBytes)
	if err != nil {
		return controlResponse{}, fmt.Errorf("unable to process tacx control response: %w", err)
	}
	log.WithFields(log.Fields{"responseRaw": fmt.Sprintf("%+v", responseRaw)}).Tracef("received tacx status raw")

	response := controlResponse{
		speed:       responseRaw.Speed,
		currentLoad: responseRaw.CurrentLoad,
		targetLoad:  responseRaw.TargetLoad,
		keepalive:   responseRaw.KeepAlive,
		cadence:     responseRaw.Cadence,
	}
	log.WithFields(log.Fields{"response": fmt.Sprintf("%+v", response)}).Debugf("received tacx status")
	return response, nil
}
