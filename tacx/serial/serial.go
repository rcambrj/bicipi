package serial

import (
	"fmt"

	"github.com/rcambrj/bicipi/tacx/common"
)

func MakeTacxDevice(port string) (*TacxSerialDevice, error) {
	commander, err := makeCommander(port)
	if err != nil {
		return &TacxSerialDevice{}, fmt.Errorf("unable to create tacx serial device: %w", err)
	}
	return &TacxSerialDevice{
		commander: commander,
	}, nil
}

type commander interface {
	sendCommand(command []byte) ([]byte, error)
	close() error
}

type TacxSerialDevice struct {
	commander commander
}

func (td *TacxSerialDevice) GetVersion() (common.Version, error) {
	return getVersion(td.commander)
}
func (td *TacxSerialDevice) SendControl(command common.ControlCommand) (common.ControlResponse, error) {
	return sendControl(td.commander, command)
}
func (td *TacxSerialDevice) Close() error {
	return td.commander.close()
}
