package ftms

import (
	"tinygo.org/x/bluetooth"
)

func CreateCyclingSpeedCadenceCharacteristics() []bluetooth.CharacteristicConfig {
	return []bluetooth.CharacteristicConfig{
		{
			UUID: bluetooth.CharacteristicUUIDCSCMeasurement,
			//Value: TODO
			Value: []byte{0x0},
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
	}
}
