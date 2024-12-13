package tacx

type targetLoadArgs struct {
	targetWatts  float64
	currentSpeed uint16
}

// the speed above which watts are applied correctly
// below this speed, the T1941 doesn't behave very well
var transitionSpeed uint16 = 4636

// given a desired wattage, calculates the load which should be sent to the
// trainer. also:
// * ensures that the trainer is easy to move from a static position
// * ensures that the power is smooth between 0 - 20km/h where the T1941 judders
func getTargetLoad(args targetLoadArgs) int16 {
	if args.targetWatts == 0 {
		return 0
	}

	if args.currentSpeed <= transitionSpeed {
		transitionLoadValue := float64(getRawLoad(args.targetWatts)) / float64(transitionSpeed)
		return int16(float64(args.currentSpeed) / float64(transitionSpeed) * float64(transitionLoadValue))
	}

	return int16(float64(getRawLoad(args.targetWatts)) / float64(args.currentSpeed))
}
