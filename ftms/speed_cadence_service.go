package ftms

import (
	"math/rand"

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

func CadenceDataGenerator() WriteValue {
	var currentRate uint8 = 60
	rateFluctuation := 10

	return func() []byte {
		rateOffset := rateFluctuation/2 - 1 - rand.Intn(rateFluctuation)
		currentRate = uint8(min(max(int(currentRate)+rateOffset, 55), 110))

		return getHeartRate(currentRate)
	}
}
