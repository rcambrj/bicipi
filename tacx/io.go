package tacx

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// https://pkg.go.dev/go.bug.st/serial#Port
type SerialPort interface {
	ResetInputBuffer() error
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
}

func readSerial(port SerialPort) ([]byte, error) {
	buff := make([]byte, 64)
	n, err := port.Read(buff)
	if err != nil {
		return []byte{}, err
	}
	return buff[:n], nil
}

func getResponse(port SerialPort) ([]byte, error) {
	var frame = make([]byte, 0, 64)
	tries := 3
	for {
		extra, err := readSerial(port)
		if err != nil {
			return []byte{}, fmt.Errorf("unable to read from serial port: %w", err)
		}
		frame = append(frame, extra...)

		// os.Stderr.WriteString(string(extra))
		if !isValidFrame(frame) {
			if tries == 0 {
				return []byte{}, fmt.Errorf("no serial response received")
			}
			log.Debugf("received partial frame: %v", frame)
			tries--
			continue
		}

		log.Debugf("received whole frame: %v", frame)
		return frame, nil
	}
}

func sendCommand(port SerialPort, command []byte) ([]byte, error) {
	log.Debugf("sending serial command: %v", command)
	outFrame, err := serializeCommand(command)
	if err != nil {
		return []byte{}, fmt.Errorf("unable to serialize command: %w", err)
	}

	port.ResetInputBuffer()

	_, err = port.Write(outFrame)
	if err != nil {
		return []byte{}, fmt.Errorf("unable to write to serial port: %w", err)
	}
	log.Debugf("sent serial data: %v", outFrame)

	inFrame, err := getResponse(port)
	if err != nil {
		return []byte{}, fmt.Errorf("unable to read from serial port: %w", err)
	}
	response, err := deserializeResponse(inFrame)
	if err != nil {
		return []byte{}, fmt.Errorf("unable to deserialize response: %w", err)
	}

	log.Debugf("received serial response: %v", response)
	return response, nil
}

type Version struct {
	Model             string
	ManufactureYear   int
	ManufactureNumber int
	FirmwareVersion   string
	Serial            uint32
	Date              string
	Other             string
}

func getVersion(port SerialPort) (Version, error) {
	response, err := sendCommand(port, []byte{0x02, 0x00, 0x00, 0x00})
	if err != nil {
		return Version{}, fmt.Errorf("unable to get version: %w", err)
	}

	firmwareVersion := fmt.Sprintf("%v.%v.%v.%v", response[31-24], response[30-24], response[29-24], response[28-24])
	date := fmt.Sprintf("%v-%v", response[37-24], response[36-24])
	other := fmt.Sprintf("%v.%v", response[39-24], response[38-24])
	serial := uint32(response[32-24] | (response[33-24] << 8) | (response[34-24] << 16) | (response[35-24] << 24))

	// serial-based properties
	manufactureYear := int(serial / 100000 % 100)
	manufactureNumber := int(serial % 100000)
	model := fmt.Sprintf("T19%v", int(serial/10000000))

	return Version{
		Model:             model,
		ManufactureYear:   manufactureYear,
		ManufactureNumber: manufactureNumber,
		FirmwareVersion:   firmwareVersion,
		Serial:            serial,
		Date:              date,
		Other:             other,
	}, nil
}
