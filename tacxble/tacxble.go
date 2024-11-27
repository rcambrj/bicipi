package tacxble

import (
	"fmt"
	"time"

	"tinygo.org/x/bluetooth"
)

func Start() {
	fmt.Println("starting...")

	adapter := bluetooth.DefaultAdapter

	must("enable BLE stack", adapter.Enable())

	// bluetooth.CharacteristicUUIDFitnessMachineFeature
	//
	// fmf_CadenceSupported                        = 1 <<  1
	// fmf_HeartRateMeasurementSupported           = 0       # 1 << 10; CTP's do not expect heartrate to be supplied by Fitness Machine
	// fmf_PowerMeasurementSupported               = 1 << 14
	// fmf_PowerTargetSettingSupported             = 1 <<  3
	// fmf_IndoorBikeSimulationParametersSupported = 1 << 13
	// fmf_Info                        = struct.pack(little_endian + unsigned_long * 2,  fmf_CadenceSupported                        |
	//                                                                                   fmf_HeartRateMeasurementSupported           |
	//                                                                                   fmf_PowerMeasurementSupported,
	//                                                                                   fmf_PowerTargetSettingSupported             |
	//                                                                                   fmf_IndoorBikeSimulationParametersSupported )
	// FM Service, section 4.3 p 19
	fitnessMachineFeatureValue := []byte{}

	// bluetooth.CharacteristicUUIDIndoorBikeData
	//
	// ibd_InstantaneousSpeedIsNotPresent  = 0         # Bit 0     # Present unless flagged that it's not
	// ibd_InstantaneousCadencePresent     = 1 << 2    # Bit 2
	// ibd_InstantaneousPowerPresent       = 1 << 6    # Bit 6
	// ibd_HeartRatePresent                = 1 << 9    # Bit 9
	// ibd_Flags                           = 0
	// 								# FM Service, section 4.9 p 44: Flags, Cadence, Power, HeartRate
	// ibd_Info                            = struct.pack(little_endian + unsigned_short * 4,
	// 											ibd_InstantaneousPowerPresent | ibd_HeartRatePresent,
	// 											123, 456, 89
	// 											)
	// FM Service, section 4.9 p 44
	indoorBikeDataValue := []byte{}

	// bluetooth.CharacteristicUUIDFitnessMachineStatus
	//
	// b'\x00\x00'
	// FM Service, section 4.17 p 56
	fitnessMachineStatusValue := []byte{}

	// bluetooth.CharacteristicUUIDFitnessMachineControlPoint
	//
	// b'\x00\x00'
	// FM Service, section 4.16 p 50
	fitnessMachineControlPointValue := []byte{}

	// bluetooth.CharacteristicUUIDSupportedPowerRange
	//
	// spr_Info                        = struct.pack(little_endian + unsigned_short * 3, 0,   1000, 1)
	// FM Service, section 4.14 p 49
	supportedPowerRangeValue := []byte{}

	ftmsUUID := bluetooth.ServiceUUIDFitnessMachine
	ftmsService := bluetooth.Service{
		UUID: ftmsUUID,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				UUID:  bluetooth.CharacteristicUUIDFitnessMachineFeature,
				Value: fitnessMachineFeatureValue,
				Flags: bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission,
			},
			{
				UUID:  bluetooth.CharacteristicUUIDIndoorBikeData,
				Value: indoorBikeDataValue,
				Flags: bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission,
			},
			{
				UUID:  bluetooth.CharacteristicUUIDFitnessMachineStatus,
				Value: fitnessMachineStatusValue,
				Flags: bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission,
			},
			{
				UUID:  bluetooth.CharacteristicUUIDFitnessMachineControlPoint,
				Value: fitnessMachineControlPointValue,
				Flags: bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission,
			},
			{
				UUID:  bluetooth.CharacteristicUUIDSupportedPowerRange,
				Value: supportedPowerRangeValue,
				Flags: bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission,
			},
		},
	}

	must("declare FTMS service", adapter.AddService(&ftmsService))

	adv := adapter.DefaultAdvertisement()
	must("configure advertisement", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: "Tacx BLE Trainer",
	}))

	adapter.SetConnectHandler(handleConnect)

	must("start advertising BLE", adv.Start())

	println("advertising BLE...")

	for {
		// Sleep forever.
		time.Sleep(time.Hour)
	}
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}

func handleConnect(device bluetooth.Device, connected bool) {
	fmt.Println("received connection")
}
