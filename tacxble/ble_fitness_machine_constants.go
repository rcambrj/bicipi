package tacxble

// 32-bit bitmask
//
// 44.3 Fitness Machine Feature
//
// https://www.bluetooth.com/specifications/specs/fitness-machine-service-1-0/
const (
	FMFAverageSpeedSupported = 1 << iota
	FMFCadenceSupported
	FMFTotalDistanceSupported
	FMFInclinationSupported
	FMFElevationGainSupported
	FMFPaceSupported
	FMFStepCountSupported
	FMFResistanceLevelSupported
	FMFStrideCountSupported
	FMFExpendedEnergySupported
	FMFHeartRateMeasurementSupported
	FMFMetabolicEquivalentSupported
	FMFElapsedTimeSupported
	FMFRemainingTimeSupported
	FMFPowerMeasurementSupported
	FMFForceOnBeltandPowerOutputSupported
	FMFUserDataRetentionSupported
)

// ???-bit bitmask
// TODO: how long is this bitmask?
//
// 4.9 Indoor Bike Data
//
// https://www.bluetooth.com/specifications/specs/fitness-machine-service-1-0/
const (
	IBDMoreData = 1 << iota
	IBDAverageSpeedPresent
	IBDInstantaneousCadence
	IBDAverageCadencepresent
	IBDTotalDistancePresent
	IBDResistanceLevelPresent
	IBDInstantaneousPowerPresent
	IBDAveragePowerPresent
	IBDExpendedEnergyPresent
	IBDHeartRatePresent
	IBDMetabolicEquivalentPresent
	IBDElapsedTimePresent
	IBDRemainingTimePresent
)

// bitmask
//
// 4.16 Fitness Machine Control Point
//
// https://www.bluetooth.com/specifications/specs/fitness-machine-service-1-0/
const (
	RequestControl                       = 0x00
	Reset                                = 0x01
	SetTargetSpeed                       = 0x02
	SetTargetInclination                 = 0x03
	SetTargetResistanceLevel             = 0x04
	SetTargetPower                       = 0x05
	SetTargetHeartRate                   = 0x06
	StartorResume                        = 0x07
	StoporPause                          = 0x08
	SetTargetedExpendedEnergy            = 0x09
	SetTargetedNumberofSteps             = 0x0A
	SetTargetedNumberofStrides           = 0x0B
	SetTargetedDistance                  = 0x0C
	SetTargetedTrainingTime              = 0x0D
	SetTargetedTimeinTwoHeartRateZones   = 0x0E
	SetTargetedTimeinThreeHeartRateZones = 0x0F
	SetTargetedTimeinFiveHeartRateZones  = 0x10
	SetIndoorBikeSimulationParameters    = 0x11
	SetWheelCircumference                = 0x12
	SpinDownControl                      = 0x13
	SetTargetedCadence                   = 0x14
	// 0x15-0x7F Reserved for Future Use
	ResponseCode = 0x80
	// 0x81-0xFF Reserved for Future Use
)

// variable-length bitmask
//
// 4.17 Fitness Machine Status
//
// https://www.bluetooth.com/specifications/specs/fitness-machine-service-1-0/
const (
	// 0x00 ReservedforFutureUse
	FMSReset                                   = 0x01
	FMSFitnessMachineStoppedorPausedbytheUser  = 0x02
	FMSFitnessMachineStoppedbySafetyKey        = 0x03
	FMSFitnessMachineStartedorResumedbytheUser = 0x04
	FMSTargetSpeedChanged                      = 0x05
	FMSTargetInclineChanged                    = 0x06
	FMSTargetResistanceLevelChanged            = 0x07
	FMSTargetPowerChanged                      = 0x08
	FMSTargetHeartRateChanged                  = 0x09
	FMSTargetedExpendedEnergyChanged           = 0x0A
	FMSTargetedNumberofStepsChanged            = 0x0B
	FMSTargetedNumberofStridesChanged          = 0x0C
	FMSTargetedDistanceChanged                 = 0x0D
	FMSTargetedTrainingTimeChanged             = 0x0E
	FMSTargetedTimeinTwoHeartRateZones         = 0x0F
	FMSTargetedTimeinThreeHeartRateZones       = 0x10
	FMSTargetedTimeinFiveHeartRateZones        = 0x11
	FMSIndoorBikeSimulationParametersChanged   = 0x12
	FMSWheelCircumferenceChanged               = 0x13
	FMSSpinDownStatus                          = 0x14
	FMSTargetedCadenceChanged                  = 0x15
	// 0xFE - 0x16 ReservedforFutureUse
	FMSControlPermissionLost = 0xFF
)
