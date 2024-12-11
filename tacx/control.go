package tacx

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type controlCommand struct {
	targetLoad  float64 // newtons (for modeNormal)
	targetSpeed float64 // km/h (for modeCalibrating)
	keepalive   uint8   // value from the last response
	mode        mode    // see const `mode`
	weight      uint8   // kg
	adjust      uint16  // adjustment resulting from calibration
}

type controlCommandRaw struct {
	target    uint16
	keepalive uint8
	mode      mode
	weight    uint8
	adjust    uint16
}

func getControlCommandBytes(args controlCommandRaw) []byte {
	return []byte{
		0x01,
		0x08,
		0x01,
		0x00,
		uint8(args.target & 0xff),
		uint8(args.target >> 8),
		args.keepalive,
		0x00,
		uint8(args.mode),
		args.weight,
		uint8(args.adjust & 0xff),
		uint8(args.adjust >> 8),
	}
}

func parseControlResponseBytes(response []byte) controlResponseRaw {
	return controlResponseRaw{
		distance:    uint32(response[4]) | uint32(response[5])<<8 | uint32(response[6])<<16 | uint32(response[7])<<24,
		speed:       uint16(response[8]) | uint16(response[9])<<8,
		averageLoad: uint16(response[12]) | uint16(response[13])<<8,
		currentLoad: uint16(response[14]) | uint16(response[15])<<8,
		targetLoad:  uint16(response[16]) | uint16(response[17])<<8,
		keepalive:   response[18],
		cadence:     response[20],
	}
}

type controlResponseRaw struct {
	distance    uint32
	speed       uint16
	averageLoad uint16
	currentLoad uint16
	targetLoad  uint16
	keepalive   uint8
	cadence     uint8
}

type controlResponse struct {
	distance    float64            // km?
	speed       float64            // km/h
	averageLoad float64            // newtons
	currentLoad float64            // newtons
	targetLoad  float64            // newtons
	keepalive   uint8              // value to send in the next control
	cadence     uint8              // rpm
	raw         controlResponseRaw // needed for some things like calibration
}

// this is the main function to send and receive data from tacx
// it both sends the target status and receives the reported status
func sendControl(t Commander, command controlCommand) (controlResponse, error) {
	log.Infof("sending tacx status: %+v", command)

	var target uint16
	if command.mode == modeNormal {
		target = getRawLoad(command.targetLoad)
	}
	if command.mode == modeCalibrating {
		target = getRawSpeed(command.targetSpeed)
	}

	commandRaw := controlCommandRaw{
		target:    target,
		keepalive: command.keepalive,
		mode:      command.mode,
		weight:    command.weight,
		adjust:    command.adjust,
	}
	log.Debugf("sending tacx status raw: %+v", commandRaw)

	commandBytes := getControlCommandBytes(commandRaw)

	responseBytes, err := t.sendCommand(commandBytes)
	if err != nil {
		return controlResponse{}, fmt.Errorf("unable to set status: %w", err)
	}

	responseRaw := parseControlResponseBytes(responseBytes)
	log.Debugf("received tacx status raw: %+v", responseRaw)

	response := controlResponse{
		distance:    getKilometers(responseRaw.distance),
		speed:       getKilometers(uint32(responseRaw.speed)),
		averageLoad: getNewtons(responseRaw.averageLoad),
		currentLoad: getNewtons(responseRaw.currentLoad),
		targetLoad:  getNewtons(responseRaw.targetLoad),
		keepalive:   responseRaw.keepalive,
		cadence:     responseRaw.cadence,
		raw:         responseRaw,
	}
	log.Infof("received tacx status: %+v", response)
	return response, nil
}
