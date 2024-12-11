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
func getKilometers(rawSpeed uint32) float64 {
	return float64(rawSpeed) / rawSpeedFactor
}

var rawLoadFactor = 137.0

func getRawLoad(newtons float64) uint16 {
	return uint16(newtons * rawLoadFactor)
}
func getNewtons(rawLoad uint16) float64 {
	return float64(rawLoad) / rawLoadFactor
}

func getRawAdjust(v int8) uint16 {
	return (uint16(v) + 8) * 130
}
func getNiceAdjust(v uint16) int16 {
	return int16(v)/130 - 8
}

func getMeanAverageUint16(list []uint16) float64 {
	var total uint16 = 0
	for _, num := range list {
		total += num
	}
	return float64(total) / float64(len(list))
}
