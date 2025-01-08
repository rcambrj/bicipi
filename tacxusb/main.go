package tacxusb

import (
	"fmt"

	"github.com/rcambrj/bicipi/tacxcommon"
)

func MakeTacxDevice() (*TacxUSBDevice, error) {
	commander, err := makeCommander()
	if err != nil {
		return &TacxUSBDevice{}, fmt.Errorf("unable to create tacx usb device: %w", err)
	}
	return &TacxUSBDevice{
		commander: commander,
	}, nil
}

type commander interface {
	sendCommand(command []byte) ([]byte, error)
	close() error
}

type TacxUSBDevice struct {
	commander commander
}

func (td *TacxUSBDevice) GetVersion() (tacxcommon.Version, error) {
	return getVersion(td.commander)
}
func (td *TacxUSBDevice) SendControl(command tacxcommon.ControlCommand) (tacxcommon.ControlResponse, error) {
	return tacxcommon.ControlResponse{}, nil
}
func (td *TacxUSBDevice) Close() error {
	return td.commander.close()
}
