package ftms

import "tinygo.org/x/bluetooth"

func CreateCyclingPowerCharacteristics(controlPointDataHandler bluetooth.WriteEvent) []bluetooth.CharacteristicConfig {
	return []bluetooth.CharacteristicConfig{
		{
			UUID:  bluetooth.CharacteristicUUIDCyclingPowerMeasurement,
			Value: []byte{0x00},
			Flags: bluetooth.CharacteristicNotifyPermission,
		},
		{
			UUID:  bluetooth.CharacteristicUUIDCyclingPowerFeature,
			Value: []byte{0x00},
			Flags: bluetooth.CharacteristicReadPermission,
		},
		{
			UUID:  bluetooth.CharacteristicUUIDSensorLocation,
			Value: []byte{0x00},
			Flags: bluetooth.CharacteristicReadPermission,
		},
		{
			UUID:       bluetooth.CharacteristicUUIDCyclingPowerControlPoint,
			Flags:      bluetooth.CharacteristicIndicatePermission | bluetooth.CharacteristicWritePermission,
			WriteEvent: controlPointDataHandler,
		},
	}
}
