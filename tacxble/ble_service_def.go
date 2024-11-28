package tacxble

import (
	"encoding/binary"
	"fmt"

	"tinygo.org/x/bluetooth"
)

func getFitnessMachineServiceDefinition() bluetooth.Service {
	ftmService := bluetooth.Service{
		UUID: bluetooth.ServiceUUIDFitnessMachine,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				UUID:  bluetooth.CharacteristicUUIDFitnessMachineFeature,
				Value: getFitnessMachineFeatures(),
				Flags: bluetooth.CharacteristicReadPermission,
			},
			{
				UUID:  bluetooth.CharacteristicUUIDIndoorBikeData,
				Value: getIndoorBikeData(),
				Flags: bluetooth.CharacteristicNotifyPermission,
			},
			{
				UUID:  bluetooth.CharacteristicUUIDFitnessMachineStatus,
				Value: getFitnessMachineStatus(),
				Flags: bluetooth.CharacteristicNotifyPermission,
			},
			{
				UUID:  bluetooth.CharacteristicUUIDFitnessMachineControlPoint,
				Value: getFitnessMachineControlPoint(),
				Flags: bluetooth.CharacteristicIndicatePermission | bluetooth.CharacteristicWritePermission,
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					printBinary(value)
					//getFitnessMachineControlPoint()
				},
			},
			{
				UUID:  bluetooth.CharacteristicUUIDSupportedPowerRange,
				Value: getSupportedPowerRange(),
				Flags: bluetooth.CharacteristicReadPermission,
			},
			// TODO: 0x2AD6
			//{
			//UUID:  bluetooth.CharacteristicUUIDSupportedResistanceLevelRange,
			//Value: TODO,
			//Flags: bluetooth.CharacteristicReadPermission,
			//},
			// TODO: 0x2AD3
			{
				UUID:  bluetooth.CharacteristicUUIDTrainingStatus,
				Value: []byte{0x01},
				Flags: bluetooth.CharacteristicNotifyPermission,
			},
		},
	}

	return ftmService
}

func getCyclingSpeedAndCadenceServiceDefinition() bluetooth.Service {
	service := bluetooth.Service{
		UUID: bluetooth.ServiceUUIDCyclingSpeedAndCadence,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				UUID: bluetooth.CharacteristicUUIDCSCMeasurement,
				//Value: TODO
				Value: []byte{0x00},
				Flags: bluetooth.CharacteristicNotifyPermission,
			},
			{
				UUID: bluetooth.CharacteristicUUIDCSCFeature,
				//Value: TODO
				Value: []byte{0x00},
				Flags: bluetooth.CharacteristicReadPermission,
			},
			{
				UUID: bluetooth.CharacteristicUUIDSensorLocation,
				//Value: TODO
				Value: []byte{0x00},
				Flags: bluetooth.CharacteristicReadPermission,
			},
		},
	}

	return service
}

func getCyclingPowerServiceDefinition() bluetooth.Service {
	service := bluetooth.Service{
		UUID: bluetooth.ServiceUUIDCyclingPower,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				UUID: bluetooth.CharacteristicUUIDCyclingPowerMeasurement,
				//Value: TODO
				Value: []byte{0x00},
				Flags: bluetooth.CharacteristicNotifyPermission,
			},
			{
				UUID: bluetooth.CharacteristicUUIDCyclingPowerFeature,
				//Value: TODO
				Value: []byte{0x00},
				Flags: bluetooth.CharacteristicReadPermission,
			},
			{
				UUID: bluetooth.CharacteristicUUIDSensorLocation,
				//Value: TODO
				Value: []byte{0x00},
				Flags: bluetooth.CharacteristicReadPermission,
			},
			{
				UUID: bluetooth.CharacteristicUUIDCyclingPowerControlPoint,
				//Value: TODO
				Flags: bluetooth.CharacteristicIndicatePermission | bluetooth.CharacteristicWritePermission,
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					printBinary(value)
				},
			},
		},
	}

	return service
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
