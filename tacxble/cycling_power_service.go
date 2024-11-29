package tacxble

import "tinygo.org/x/bluetooth"

func createCyclingPowerCharacteristics() []bluetooth.CharacteristicConfig {
	return []bluetooth.CharacteristicConfig{
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
				unhandledWriteEvent("CharacteristicUUIDCyclingPowerControlPoint", offset, value)
			},
		},
	}
}
