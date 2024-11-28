package tacxble

import (
	"encoding/binary"

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

func getBytesFromBitmask(bitmask uint32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bitmask)
	return bytes
}

func getFitnessMachineFeatures() []byte {
	var bitmask uint32 = 0
	bitmask |= FMFCadenceSupported
	return getBytesFromBitmask(bitmask)
}
func getIndoorBikeData() []byte {
	var bitmask uint32 = 0
	return getBytesFromBitmask(bitmask)
}
func getFitnessMachineControlPoint() []byte {
	var bitmask uint32 = 0
	return getBytesFromBitmask(bitmask)
}
func getFitnessMachineStatus() []byte {
	var bitmask uint32 = 0
	return getBytesFromBitmask(bitmask)
}
func getSupportedPowerRange() []byte {
	var bitmask uint32 = 0
	return getBytesFromBitmask(bitmask)
}
