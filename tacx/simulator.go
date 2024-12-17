package tacx

import (
	"math"

	log "github.com/sirupsen/logrus"
)

type targetLoadForSimulatorArgs struct {
	currentSpeed      uint16
	weight            uint8   // kg
	windSpeed         float64 // m/s
	gradient          float64 // Percentage 0...100
	rollingResistance float64 // ??
	windResistance    float64 // CdA
}

// for a set of resistances at a current speed, calculates watts necessary to overcome
// based on FortiusAnt's interpretation of https://www.gribble.org/cycling/power_v_speed.html
func getWattsForSimulator(args targetLoadForSimulatorArgs) float64 {
	speed := getKilometers(args.currentSpeed) / 3.6 // m/s

	rollingResistance := max(args.rollingResistance, 0.004) // ??
	weight := args.weight                                   // kg
	gravity := 9.81                                         // m/s2
	rollingWatts := rollingResistance * float64(weight) * gravity * speed

	windResistance := max(args.windResistance, 0.51) // CdA
	windSpeed := args.windSpeed                      // m/s
	draftingFactor := 1.0                            // multiplier, not supplied
	// without abs a strong tailwind would result in a higher power
	airWatts := 0.5 * windResistance * (speed + windSpeed) * math.Abs(speed+windSpeed) * draftingFactor * speed

	gradient := args.gradient
	gravityWatts := gradient / 100 * float64(weight) * gravity * speed

	totalWatts := rollingWatts + airWatts + gravityWatts
	log.Debugf("simulating rollingWatts: %v; airWatts: %v; gravityWatts: %v; total: %v", rollingWatts, airWatts, gravityWatts, totalWatts)

	return totalWatts
}
