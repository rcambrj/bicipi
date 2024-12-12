package tacx

type mode uint8

const (
	modeOff mode = iota
	_            // 1 is unused
	modeNormal
	modeCalibrating
)

type Behaviour uint8

const (
	BehaviourERG Behaviour = iota
	BehaviourSlope
)

var rawSpeedFactor = 289.75

func getRawSpeed(kilometers float64) uint16 {
	return uint16(kilometers * rawSpeedFactor)
}
func getKilometers(rawSpeed uint16) float64 {
	return float64(rawSpeed) / rawSpeedFactor
}

// var rawLoadFactor = 137.0

// func getRawLoad(newtons float64) int16 {
// 	return int16(newtons * rawLoadFactor)
// }
// func getNewtons(rawLoad int16) float64 {
// 	return float64(rawLoad) / rawLoadFactor
// }
