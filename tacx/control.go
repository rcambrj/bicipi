package tacx

import (
	"bytes"
	"encoding/binary"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type controlCommand struct {
	mode        mode        // see const `mode`
	behaviour   Behaviour   // see const `Behaviour`
	targetSpeed int16       // modeCalibrating: raw speed
	targetLoad  int16       // modeNormal + BehaviourERG: raw load
	targetSlope targetSlope // modeNormal + BehaviourSlope: struct
	keepalive   uint8       // value from the last response
	weight      uint8       // kg
	adjust      uint16      // adjustment resulting from calibration
}

type targetSlope struct {
	degrees int8
	wind    int8
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
	log.Debugf("sending tacx status: %+v", command)

	var target int16
	switch command.mode {
	case modeCalibrating:
		target = command.targetSpeed
	case modeNormal:
		switch command.behaviour {
		case BehaviourERG:
			target = command.targetLoad
		case BehaviourSlope:
			// TODO
		}
	}

	commandRaw := controlCommandRaw{
		target:    target,
		keepalive: command.keepalive,
		mode:      uint8(command.mode),
		weight:    command.weight,
		adjust:    command.adjust,
	}
	log.Tracef("sending tacx status raw: %+v", commandRaw)

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
	log.Tracef("received tacx status raw: %+v", responseRaw)

	response := controlResponse{
		speed:       responseRaw.Speed,
		currentLoad: responseRaw.CurrentLoad,
		targetLoad:  responseRaw.TargetLoad,
		keepalive:   responseRaw.KeepAlive,
		cadence:     responseRaw.Cadence,
	}
	log.Debugf("received tacx status: %+v", response)
	return response, nil
}
