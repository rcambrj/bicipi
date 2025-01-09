package usb

import (
	"fmt"

	"github.com/rcambrj/bicipi/tacx/common"
	log "github.com/sirupsen/logrus"
)

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
	if !isValidFrame(responseBytes, frameTypeControl) {
		log.Warn("received invalid frame")
		return common.ControlResponse{}, ErrReceivedInvalidFrame
	}

	response, err := common.GetControlResponseFromBytes(responseBytes[24:48])
	if err != nil {
		return common.ControlResponse{}, fmt.Errorf("unable to process tacx control response: %w", err)
	}

	log.WithFields(log.Fields{"response": fmt.Sprintf("%+v", response)}).Debug("received tacx status")
	return response, nil
}
