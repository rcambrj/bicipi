package tacx

import (
	"math"

	log "github.com/sirupsen/logrus"
)

type targetLoadForSlopeArgs struct {
	currentSpeed   uint16
	weight         int // kg
	windSpeed      int // m/s
	draftingFactor int // multiplier (default 1)
	gradient       int // Percentage 0...100
}

func getWattsForSlope(args targetLoadForSlopeArgs) float64 {
	rollingResistance := 0.004
	weight := args.weight                           // kg
	gravity := 9.81                                 // m/s2
	speed := getKilometers(args.currentSpeed) / 3.6 // m/s
	rollingWatts := float64(rollingResistance) * float64(weight) * gravity * speed

	windResistance := 0.51 // default=0.51
	windSpeed := args.windSpeed
	draftingFactor := args.draftingFactor
	// without abs a strong tailwind would result in a higher power
	airWatts := 0.5 * windResistance * (float64(speed) + float64(windSpeed)) * math.Abs(float64(speed)+float64(windSpeed)) * float64(draftingFactor) * speed

	gradient := float64(args.gradient)
	gravityWatts := gradient / 100 * float64(weight) * gravity * speed

	log.Debugf("rollingWatts: %v; airWatts: %v; gravityWatts: %v", rollingWatts, airWatts, gravityWatts)

	return rollingWatts + airWatts + gravityWatts
}
