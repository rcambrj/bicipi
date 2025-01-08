package tacxcommon

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

// ControlResponse could contain more information if the connection is USB.
// in order to reduce maintenance, ControlResponse contains only what's
// supported on all connection protocols.
// If the extra USB data is necessary, then ControlResponse should not live here
// as the responses for serial vs USB will differ
type ControlResponse struct {
	Speed       uint16 // tacx speed units
	CurrentLoad int16  // tacx load units
	TargetLoad  int16  // tacx load units
	Keepalive   uint8  // value to send in the next control
	Cadence     uint8  // rpm
}
