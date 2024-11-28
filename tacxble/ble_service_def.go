package tacxble

import (
	"encoding/binary"
	"fmt"

	"tinygo.org/x/bluetooth"
)

func getBLEServiceDefinition() bluetooth.Service {
	ftmService := bluetooth.Service{
		UUID: bluetooth.ServiceUUIDFitnessMachine,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				UUID:  bluetooth.CharacteristicUUIDFitnessMachineFeature,
				Value: getFitnessMachineFeatures(),
				Flags: bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission,
			},
			{
				UUID:  bluetooth.CharacteristicUUIDIndoorBikeData,
				Value: getIndoorBikeData(),
				Flags: bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission,
			},
			{
				UUID:  bluetooth.CharacteristicUUIDFitnessMachineControlPoint,
				Value: getFitnessMachineControlPoint(),
				Flags: bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission,
			},
			{
				UUID:  bluetooth.CharacteristicUUIDFitnessMachineStatus,
				Value: getFitnessMachineStatus(),
				Flags: bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission,
			},
			{
				UUID:  bluetooth.CharacteristicUUIDSupportedPowerRange,
				Value: getSupportedPowerRange(),
				Flags: bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission,
			},
		},
	}

	return ftmService
}

func printBinary(bytes []byte) {
	for _, n := range bytes {
		fmt.Printf("%08b ", n) // prints 00000000 11111101
	}
	fmt.Printf("\n")
}

func getFitnessMachineFeatures() []byte {
	// confusing: this contains 4.3.1.1 Fitness Machine Features & 4.3.1.2 Target Setting Features
	var featuresBitmask uint32 = FMFCadenceSupported | FMFPowerMeasurementSupported
	var targetSettingsBitmask uint32 = TSFPowerTargetSettingSupported | TSFIndoorBikeSimulationParametersSupported
	bytes := make([]byte, 0, 64)
	bytes = binary.LittleEndian.AppendUint32(bytes, featuresBitmask)
	bytes = binary.LittleEndian.AppendUint32(bytes, targetSettingsBitmask)
	fmt.Println("getFitnessMachineFeatures")
	// FortiusAnt: 00000010 01000000 00000000 00000000 00001000 00100000 00000000 00000000
	printBinary(bytes)
	return bytes
}
func getIndoorBikeData() []byte {
	var bitmask uint16 = IBDInstantaneousPowerPresent | IBDHeartRatePresent
	bytes := make([]byte, 4*16/8)
	binary.LittleEndian.PutUint16(bytes, bitmask)
	fmt.Println("getIndoorBikeData")
	// FortiusAnt: 01000000 00000010 01111011 00000000 11001000 00000001 01011001 00000000
	printBinary(bytes)
	return bytes
}
func getFitnessMachineControlPoint() []byte {
	// Collector commands sent to Server
	bytes := []byte{0x00, 0x00}
	fmt.Println("getFitnessMachineControlPoint")
	// FortiusAnt: 00000000 00000000
	printBinary(bytes)
	return bytes
}
func getFitnessMachineStatus() []byte {
	// Server status sent to Collector
	bytes := []byte{0x00, 0x00}
	fmt.Println("getFitnessMachineStatus")
	// FortiusAnt: 00000000 00000000
	printBinary(bytes)
	return bytes
}
func getSupportedPowerRange() []byte {
	bytes := make([]byte, 0, 3*16/8)
	bytes = binary.LittleEndian.AppendUint16(bytes, 0)    // min
	bytes = binary.LittleEndian.AppendUint16(bytes, 1000) // max
	bytes = binary.LittleEndian.AppendUint16(bytes, 1)    // step
	fmt.Println("getSupportedPowerRange")
	// FortiusAnt: 00000000 00000000 11101000 00000011 00000001 00000000
	printBinary(bytes)
	return bytes
}
