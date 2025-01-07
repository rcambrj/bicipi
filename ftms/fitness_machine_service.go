package ftms

import (
	"encoding/binary"

	log "github.com/sirupsen/logrus"
	"tinygo.org/x/bluetooth"
)

func CreateFitnessMachineCharacteristics(controlPointDataHandler bluetooth.WriteEvent) []bluetooth.CharacteristicConfig {
	return []bluetooth.CharacteristicConfig{
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
			UUID:       bluetooth.CharacteristicUUIDFitnessMachineControlPoint,
			Value:      getFitnessMachineControlPoint(),
			Flags:      bluetooth.CharacteristicIndicatePermission | bluetooth.CharacteristicWritePermission,
			WriteEvent: controlPointDataHandler,
		},
		{
			UUID:  bluetooth.CharacteristicUUIDSupportedPowerRange,
			Value: getSupportedPowerRange(),
			Flags: bluetooth.CharacteristicReadPermission,
		},
		////TODO: 0x2AD6
		//{
		//  UUID:  bluetooth.CharacteristicUUIDSupportedResistanceLevelRange,
		//  Value: TODO,
		//  Flags: bluetooth.CharacteristicReadPermission,
		//},
		////TODO: 0x2AD3
		{
			UUID:  bluetooth.CharacteristicUUIDTrainingStatus,
			Value: []byte{0x01},
			Flags: bluetooth.CharacteristicNotifyPermission,
		},
	}
}

func getFitnessMachineFeatures() []byte {
	// confusing: this contains 4.3.1.1 Fitness Machine Features & 4.3.1.2 Target Setting Features
	var featuresBitmask uint32 = FMFCadenceSupported | FMFPowerMeasurementSupported
	var targetSettingsBitmask uint32 = TSFPowerTargetSettingSupported | TSFIndoorBikeSimulationParametersSupported
	bytes := make([]byte, 0, 64)
	bytes = binary.LittleEndian.AppendUint32(bytes, featuresBitmask)
	bytes = binary.LittleEndian.AppendUint32(bytes, targetSettingsBitmask)
	// FortiusAnt: 00000010 01000000 00000000 00000000 00001000 00100000 00000000 00000000
	log.WithFields(log.Fields{"characteristic": formatBinary(bytes)}).Trace("defining ble FitnessMachineFeatures")
	return bytes
}

func getIndoorBikeData() []byte {
	var bitmask uint16 = IBDInstantaneousCadence | IBDInstantaneousPowerPresent
	bytes := make([]byte, 4*16/8)
	binary.LittleEndian.PutUint16(bytes, bitmask)
	// FortiusAnt: 01000000 00000010 01111011 00000000 11001000 00000001 01011001 00000000
	log.WithFields(log.Fields{"characteristic": formatBinary(bytes)}).Trace("defining ble IndoorBikeData")
	return bytes
}

func getFitnessMachineStatus() []byte {
	// Server status sent to Collector
	bytes := []byte{0x00, 0x00}
	// FortiusAnt: 00000000 00000000
	log.WithFields(log.Fields{"characteristic": formatBinary(bytes)}).Trace("defining ble FitnessMachineStatus")
	return bytes
}

func getFitnessMachineControlPoint() []byte {
	// Collector commands sent to Server
	bytes := []byte{0x00, 0x00}
	// FortiusAnt: 00000000 00000000
	log.WithFields(log.Fields{"characteristic": formatBinary(bytes)}).Trace("defining ble FitnessMachineControlPoint")
	return bytes
}

func getSupportedPowerRange() []byte {
	bytes := make([]byte, 0, 3*16/8)
	bytes = binary.LittleEndian.AppendUint16(bytes, 0)    // min
	bytes = binary.LittleEndian.AppendUint16(bytes, 1000) // max
	bytes = binary.LittleEndian.AppendUint16(bytes, 1)    // step
	// FortiusAnt: 00000000 00000000 11101000 00000011 00000001 00000000
	log.WithFields(log.Fields{"characteristic": formatBinary(bytes)}).Trace("defining ble SupportedPowerRange")
	return bytes
}
