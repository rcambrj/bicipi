package tacxserial

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/rcambrj/bicipi/tacxcommon"
	log "github.com/sirupsen/logrus"
)

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

// this is the main function to send and receive data from tacx
// it both sends the target status and receives the reported status
func sendControl(t commander, command tacxcommon.ControlCommand) (tacxcommon.ControlResponse, error) {
	log.WithFields(log.Fields{"command": fmt.Sprintf("%+v", command)}).Debug("sending tacx status")

	var target int16
	switch command.Mode {
	case tacxcommon.ModeCalibrating:
		target = command.TargetSpeed
	case tacxcommon.ModeNormal:
		target = command.TargetLoad
	}

	commandRaw := controlCommandRaw{
		target:    target,
		keepalive: command.Keepalive,
		mode:      uint8(command.Mode),
		weight:    command.Weight,
		adjust:    command.Adjust,
	}
	log.WithFields(log.Fields{"commandRaw": fmt.Sprintf("%+v", commandRaw)}).Trace("sending tacx status raw")

	commandBytes, err := getControlCommandBytes(commandRaw)
	if err != nil {
		return tacxcommon.ControlResponse{}, fmt.Errorf("unable to process tacx control command: %w", err)
	}

	responseBytes, err := t.sendCommand(commandBytes)
	if err != nil {
		return tacxcommon.ControlResponse{}, fmt.Errorf("unable to send tacx control command: %w", err)
	}

	responseRaw, err := parseControlResponseBytes(responseBytes)
	if err != nil {
		return tacxcommon.ControlResponse{}, fmt.Errorf("unable to process tacx control response: %w", err)
	}
	log.WithFields(log.Fields{"responseRaw": fmt.Sprintf("%+v", responseRaw)}).Trace("received tacx status raw")

	response := tacxcommon.ControlResponse{
		Speed:       responseRaw.Speed,
		CurrentLoad: responseRaw.CurrentLoad,
		TargetLoad:  responseRaw.TargetLoad,
		Keepalive:   responseRaw.KeepAlive,
		Cadence:     responseRaw.Cadence,
	}
	log.WithFields(log.Fields{"response": fmt.Sprintf("%+v", response)}).Debug("received tacx status")
	return response, nil
}
