package ftms

import (
	"math/rand"

	"tinygo.org/x/bluetooth"
)

func CreateHeartRateCharacteristics() []bluetooth.CharacteristicConfig {
	return []bluetooth.CharacteristicConfig{
		{
			UUID:  bluetooth.CharacteristicUUIDHeartRateMeasurement,
			Value: getHeartRate(69),
			Flags: bluetooth.CharacteristicNotifyPermission,
		},
		{
			UUID: bluetooth.New16BitUUID(0x2902),
			//Value: TODO
			Value: []byte{0x00},
			Flags: bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission,
		},
	}
}

func HeartRateDataGenerator() WriteValue {
	var currentRate uint8 = 60
	rateFluctuation := 10

	return func() []byte {
		rateOffset := rateFluctuation/2 - 1 - rand.Intn(rateFluctuation)
		currentRate = uint8(min(max(int(currentRate)+rateOffset, 55), 110))

		return getHeartRate(currentRate)
	}
}

func getHeartRate(heartRate uint8) []byte {
	bytes := []byte{
		0,         // flags
		heartRate, // heartrate
	}
	return bytes
}
