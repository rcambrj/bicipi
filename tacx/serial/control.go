package serial

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/rcambrj/bicipi/tacx/common"
	log "github.com/sirupsen/logrus"
)

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
func sendControl(t commander, command common.ControlCommand) (common.ControlResponse, error) {
	log.WithFields(log.Fields{"command": fmt.Sprintf("%+v", command)}).Debug("sending tacx status")

	commandBytes, err := common.GetControlCommandBytes(command)
	if err != nil {
		return common.ControlResponse{}, fmt.Errorf("unable to process tacx control command: %w", err)
	}

	responseBytes, err := t.sendCommand(commandBytes)
	if err != nil {
		return common.ControlResponse{}, fmt.Errorf("unable to send tacx control command: %w", err)
	}

	responseRaw, err := parseControlResponseBytes(responseBytes)
	if err != nil {
		return common.ControlResponse{}, fmt.Errorf("unable to process tacx control response: %w", err)
	}
	log.WithFields(log.Fields{"responseRaw": fmt.Sprintf("%+v", responseRaw)}).Trace("received tacx status raw")

	response := common.ControlResponse{
		Speed:       responseRaw.Speed,
		CurrentLoad: responseRaw.CurrentLoad,
		TargetLoad:  responseRaw.TargetLoad,
		Keepalive:   responseRaw.KeepAlive,
		Cadence:     responseRaw.Cadence,
	}
	log.WithFields(log.Fields{"response": fmt.Sprintf("%+v", response)}).Debug("received tacx status")
	return response, nil
}
