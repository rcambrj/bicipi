package tacxserial

import (
	"fmt"

	"github.com/rcambrj/bicipi/tacxcommon"
)

func MakeTacxDevice(port string) (TacxSerialDevice, error) {
	commander, err := makeCommander(port)
	if err != nil {
		return TacxSerialDevice{}, fmt.Errorf("unable to create tacx serial device: %w", err)
	}
	return TacxSerialDevice{
		commander: commander,
	}, nil
}

type TacxSerialDevice struct {
	commander commander
}

type commander interface {
	sendCommand(command []byte) ([]byte, error)
	close() error
}

func (td *TacxSerialDevice) GetVersion() (tacxcommon.Version, error) {
	return GetVersion(td.commander)
}
func (td *TacxSerialDevice) SendControl(command tacxcommon.ControlCommand) (tacxcommon.ControlResponse, error) {
	return SendControl(td.commander, command)
}
func (td *TacxSerialDevice) Close() error {
	return td.commander.close()
}
