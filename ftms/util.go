package ftms

import (
	"fmt"
	"strings"
)

type Mode uint8

const (
	ModeTargetPower Mode = iota
	ModeIndoorBikeSimulation
)

func formatBinary(bytes []byte) string {
	var output string
	for _, n := range bytes {
		output += fmt.Sprintf("%08b ", n) // format each byte to a binary octet
	}
	return strings.TrimSpace(output)
}

func getWatts(rawPower int16) int16 {
	return rawPower // no conversion: ble uses watts
}
func getRawPower(watts int16) int16 {
	return watts
}

func getRawWindSpeed() {
	// TODO
}
