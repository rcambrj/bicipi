package ftms

// 32-bit bitmask
//
// 4.3.1.1 Fitness Machine Features
//
// https://www.bluetooth.com/specifications/specs/fitness-machine-service-1-0/
const (
	FMFAverageSpeedSupported              = 1 << 0
	FMFCadenceSupported                   = 1 << 1
	FMFTotalDistanceSupported             = 1 << 2
	FMFInclinationSupported               = 1 << 3
	FMFElevationGainSupported             = 1 << 4
	FMFPaceSupported                      = 1 << 5
	FMFStepCountSupported                 = 1 << 6
	FMFResistanceLevelSupported           = 1 << 7
	FMFStrideCountSupported               = 1 << 8
	FMFExpendedEnergySupported            = 1 << 9
	FMFHeartRateMeasurementSupported      = 1 << 10
	FMFMetabolicEquivalentSupported       = 1 << 11
	FMFElapsedTimeSupported               = 1 << 12
	FMFRemainingTimeSupported             = 1 << 13
	FMFPowerMeasurementSupported          = 1 << 14
	FMFForceOnBeltandPowerOutputSupported = 1 << 15
	FMFUserDataRetentionSupported         = 1 << 16
	// 17-31 Reserved for Future Use
)

// 32-bit bitmask
//
// 4.3.1.2 Target Setting Features
//
// https://www.bluetooth.com/specifications/specs/fitness-machine-service-1-0/
const (
	TSFSpeedTargetSettingSupported = 1 << iota
	TSFInclinationTargetSettingSupported
	TSFResistanceTargetSettingSupported
	TSFPowerTargetSettingSupported
	TSFHeartRateTargetSettingSupported
	TSFTargetedExpendedEnergyConfigurationSupported
	TSFTargetedStepNumberConfigurationSupported
	TSFTargetedStrideNumberConfigurationSupported
	TSFTargetedDistanceConfigurationSupported
	TSFTargetedTrainingTimeConfigurationSupported
	TSFTargetedTimeinTwoHeartRateZonesConfigurationSupported
	TSFTargetedTimeinThreeHeartRateZonesConfigurationSupported
	TSFTargetedTimeinFiveHeartRateZonesConfigurationSupported
	TSFIndoorBikeSimulationParametersSupported
	TSFWheelCircumferenceConfigurationSupported
	TSFSpinDownControlSupported
	TSFTargetedCadenceConfigurationSupported
	// 17-31 Reserved for Future Use
)

// 16-bit bitmask
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
	// No bits reserved for future use
)
