package tacx

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

func read(ch chan []byte, port serial.Port) {
	for {
		buff := make([]byte, 64)
		n, err := port.Read(buff)
		if err != nil {
			log.Fatal(fmt.Errorf("unable to read from serial port: %w", err))
		}
		log.Debugf("received serial data: %v", buff[:n])
		ch <- buff[:n]
		if n == 0 {
			log.Fatal("serial port disconnected")
		}
	}
}

func waitForResponse(ch chan []byte, port serial.Port) []byte {
	var frame = make([]byte, 0, 64)
	for {
		extra := <-ch
		frame = append(frame, extra...)

		if !isValidFrame(frame) {
			log.Debugf("received partial frame: %v", frame)
			continue
		}

		log.Debugf("received whole frame: %v", frame)
		return frame
	}
}
