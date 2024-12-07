package ftms

import "tinygo.org/x/bluetooth"

// 347b0001-7635-408b-8918-8ff3949ce592
var ServiceUUIDCyclingSteering = bluetooth.New32BitUUID(0x347B0001)

func CreateCyclingSteeringCharacteristics() []bluetooth.CharacteristicConfig {
	return []bluetooth.CharacteristicConfig{
		//{
		//  UUID: bluetooth.Charater,
		//  //Value: TODO
		//  Value: []byte{0x00},
		//  Flags: bluetooth.CharacteristicNotifyPermission,
		//},
		//{
		//	UUID: bluetooth.CharacteristicUUIDCSCFeature,
		//	//Value: TODO
		//	Value: []byte{0x00},
		//	Flags: bluetooth.CharacteristicReadPermission,
		//},
		//{
		//	UUID: bluetooth.CharacteristicUUIDSensorLocation,
		//	//Value: TODO
		//	Value: []byte{0x00},
		//	Flags: bluetooth.CharacteristicReadPermission,
		//},
	}
}
