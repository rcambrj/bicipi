package tacxusb

import (
	"fmt"

	"github.com/rcambrj/bicipi/tacxcommon"
)

func MakeTacxDevice() (TacxUSBDevice, error) {
	commander, err := makeCommander()
	if err != nil {
		return TacxUSBDevice{}, fmt.Errorf("unable to create tacx usb device: %w", err)
	}
	return TacxUSBDevice{
		commander: *commander,
	}, nil
}

type TacxUSBDevice struct {
	commander commander
}

type Version struct {
	Model             string
	ManufactureYear   int
	ManufactureNumber int
	FirmwareVersion   string
	Serial            int32
	Date              string
}

type ControlCommand struct {
	Mode        tacxcommon.Mode // see const `Mode`
	TargetSpeed int16           // modeCalibrating: raw speed
	TargetLoad  int16           // modeNormal: raw load
	Keepalive   uint8           // value from the last response
	Weight      uint8           // kg
	Adjust      uint16          // adjustment resulting from calibration
}

type ControlResponse struct {
	Speed       uint16 // tacx speed units
	CurrentLoad int16  // tacx load units
	TargetLoad  int16  // tacx load units
	Keepalive   uint8  // value to send in the next control
	Cadence     uint8  // rpm
}

func (td *TacxUSBDevice) GetVersion() (Version, error) {
	return Version{}, nil
}
func (td *TacxUSBDevice) SendControl(command ControlCommand) (ControlResponse, error) {
	return ControlResponse{}, nil
}
func (td *TacxUSBDevice) Close() error {
	return nil
}
