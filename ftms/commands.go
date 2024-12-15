package ftms

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"tinygo.org/x/bluetooth"
)

const (
	FMCPOpCodeRequestControl          = 0x00
	FMCPOpCodeReset                   = 0x01
	FMCPOpCodeSetTargetPower          = 0x05
	FMCPOpCodeStartOrResume           = 0x07
	FMCPOpCodeStopOrPause             = 0x08
	FMCPOpCodeSetIndoorBikeSimulation = 0x11
	FMCPOpCodeResponseCode            = 0x80
)
const (
	FMCPResultCodeSuccess             = 0x01
	FMCPResultCodeOpCodeNotSupported  = 0x02
	FMCPResultCodeInvalidParameter    = 0x03
	FMCPResultCodeOperationFailed     = 0x04
	FMCPResultCodeControlNotPermitted = 0x05
)

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

const (
	FMSReset                                 = 0x01
	FMSFitnessMachineStoppedOrPausedByUser   = 0x02
	FMSFitnessMachineStartedOrResumedByUser  = 0x04
	FMSTargetPowerChanged                    = 0x08
	FMSIndoorBikeSimulationParametersChanged = 0x12
)

type fmcpSetTargetPower struct {
	_           uint8
	TargetPower int16 // watts
}

func readFMCPTargetPower(message []byte) (fmcpSetTargetPower, error) {
	buf := bytes.NewReader(message)
	out := fmcpSetTargetPower{}
	if err := binary.Read(buf, binary.LittleEndian, &out); err != nil {
		return fmcpSetTargetPower{}, err
	}
	return out, nil
}

func writeFMSTargetPower(serviceManager *ServiceManager, targetPower int16) error {
	bytes := make([]byte, 0, (8+16)/8)
	bytes = append(bytes, FMSTargetPowerChanged)
	binary.LittleEndian.AppendUint16(bytes, uint16(targetPower))

	_, err := serviceManager.WriteToCharacteristic(bluetooth.ServiceUUIDFitnessMachine, bluetooth.CharacteristicUUIDFitnessMachineStatus, bytes)
	if err != nil {
		return fmt.Errorf("unable to write to FitnessMachineStatus: %w", err)
	}
	return nil
}

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
