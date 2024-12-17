package ftms

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"tinygo.org/x/bluetooth"
)

// this is the way trainer data goes out through ble
func writeFMIndoorBikeData(serviceManager *ServiceManager, speed uint16, cadence uint16, load int16) error {
	bytes := make([]byte, 0, 4*16/8)
	bytes = binary.LittleEndian.AppendUint16(bytes, IBDInstantaneousCadence|IBDInstantaneousPowerPresent)
	bytes = binary.LittleEndian.AppendUint16(bytes, speed) // always present
	bytes = binary.LittleEndian.AppendUint16(bytes, cadence)
	bytes = binary.LittleEndian.AppendUint16(bytes, uint16(load))

	_, err := serviceManager.WriteToCharacteristic(bluetooth.ServiceUUIDFitnessMachine, bluetooth.CharacteristicUUIDIndoorBikeData, bytes)
	if err != nil {
		return fmt.Errorf("unable to write to IndoorBikeData: %w", err)
	}
	return nil

}

// Fitness Machine Control Point op code
const (
	FMCPOpCodeRequestControl          = 0x00
	FMCPOpCodeReset                   = 0x01
	FMCPOpCodeSetTargetPower          = 0x05
	FMCPOpCodeStartOrResume           = 0x07
	FMCPOpCodeStopOrPause             = 0x08
	FMCPOpCodeSetIndoorBikeSimulation = 0x11
	FMCPOpCodeResponseCode            = 0x80
)

// Fitness Machine Control Point result code
const (
	FMCPResultCodeSuccess             = 0x01
	FMCPResultCodeOpCodeNotSupported  = 0x02
	FMCPResultCodeInvalidParameter    = 0x03
	FMCPResultCodeOperationFailed     = 0x04
	FMCPResultCodeControlNotPermitted = 0x05
)

// helper to write result codes on the Fitness Machine Control Point
func writeFMCPResultCode(serviceManager *ServiceManager, opCode uint8, resultCode uint8) error {
	bytes := make([]byte, 0, 3*8/8)
	bytes = append(bytes, FMCPOpCodeResponseCode)
	bytes = append(bytes, opCode)
	bytes = append(bytes, resultCode)

	_, err := serviceManager.WriteToCharacteristic(bluetooth.ServiceUUIDFitnessMachine, bluetooth.CharacteristicUUIDFitnessMachineControlPoint, bytes)
	if err != nil {
		return fmt.Errorf("unable to write to FitnessMachineControlPoint: %w", err)
	}
	return nil
}

// Fitness Machine Status op code
const (
	FMSOpCodeReset                                 = 0x01
	FMSOpCodeFitnessMachineStoppedOrPausedByUser   = 0x02
	FMSOpCodeFitnessMachineStartedOrResumedByUser  = 0x04
	FMSOpCodeTargetPowerChanged                    = 0x08
	FMSOpCodeIndoorBikeSimulationParametersChanged = 0x12
)

// the shape of the Set Target Power operation received on Fitness Machine Control Point
type fmcpSetTargetPower struct {
	_           uint8
	TargetPower int16 // watts
}

// helper to read the Set Target Power operation received on Fitness Machine Control Point
func readFMCPTargetPower(message []byte) (fmcpSetTargetPower, error) {
	buf := bytes.NewReader(message)
	out := fmcpSetTargetPower{}
	if err := binary.Read(buf, binary.LittleEndian, &out); err != nil {
		return fmcpSetTargetPower{}, err
	}
	return out, nil
}

// helper to write the Target Power Changed operation to Fitness Machine Status
func writeFMSTargetPower(serviceManager *ServiceManager, command fmcpSetTargetPower) error {
	bytes := make([]byte, 0, (8+16)/8)
	bytes = append(bytes, FMSOpCodeTargetPowerChanged)
	binary.LittleEndian.AppendUint16(bytes, uint16(command.TargetPower))

	_, err := serviceManager.WriteToCharacteristic(bluetooth.ServiceUUIDFitnessMachine, bluetooth.CharacteristicUUIDFitnessMachineStatus, bytes)
	if err != nil {
		return fmt.Errorf("unable to write to FitnessMachineStatus: %w", err)
	}
	return nil
}

// the shape of the Set Indoor Bike Simulation operation received on Fitness Machine Control Point
type fmcpIndoorBikeSimulation struct {
	_                 uint8
	WindSpeed         int16
	TargetGrade       int16
	RollingResistance uint8
	WindResistance    uint8
}

// helper to read the Set Indoor Bike Simulation operation received on Fitness Machine Control Point
func readFMCPIndoorBikeSimulation(message []byte) (fmcpIndoorBikeSimulation, error) {
	buf := bytes.NewReader(message)
	out := fmcpIndoorBikeSimulation{}
	if err := binary.Read(buf, binary.LittleEndian, &out); err != nil {
		return fmcpIndoorBikeSimulation{}, err
	}
	return out, nil
}

// helper to write the Indoor Bike Simulation Parameters Changed operation to Fitness Machine Status
func writeFMSIndoorBikeSimulation(serviceManager *ServiceManager, command fmcpIndoorBikeSimulation) error {
	bytes := make([]byte, 0, (8*2+16*2)/8)
	bytes = append(bytes, FMSOpCodeTargetPowerChanged)
	binary.LittleEndian.AppendUint16(bytes, uint16(command.WindSpeed))
	binary.LittleEndian.AppendUint16(bytes, uint16(command.TargetGrade))
	bytes = append(bytes, command.RollingResistance)
	bytes = append(bytes, command.WindResistance)

	_, err := serviceManager.WriteToCharacteristic(bluetooth.ServiceUUIDFitnessMachine, bluetooth.CharacteristicUUIDFitnessMachineStatus, bytes)
	if err != nil {
		return fmt.Errorf("unable to write to FitnessMachineStatus: %w", err)
	}
	return nil
}
