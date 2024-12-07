package tacx

import (
	"encoding/binary"
	"fmt"
)

var startOfFrame byte = 0x01
var endOfFrame byte = 0x17

// converts the message
// calculates a checksum and appends it
// prepends start of frame
// appends end of frame
func SerializeCommand(message []byte) ([]byte, error) {
	serialized := make([]byte, 0, 30)
	for _, b := range message {
		for _, nibble := range []byte{b >> 4 & 0xf, b >> 0 & 0xf} {
			h, err := getHexFromBin(nibble)
			if err != nil {
				return nil, fmt.Errorf("unable to serialize command: %w", err)
			}
			serialized = append(serialized, h)
		}
	}

	checksum := getChecksum(serialized)
	for _, nibble := range []uint16{checksum >> 4 & 0xf, checksum >> 0 & 0xf, checksum >> 12 & 0xf, checksum >> 8 & 0xf} {
		bytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(bytes, nibble)
		h, err := getHexFromBin(bytes[0])
		if err != nil {
			return nil, fmt.Errorf("unable to serialize command: %w", err)
		}
		serialized = append(serialized, h)
	}

	serialized = append([]byte{startOfFrame}, serialized...)
	serialized = append(serialized, endOfFrame)

	return serialized, nil
}

func getParity16(b uint16) int {
	b ^= b >> 8
	b ^= b >> 4
	b &= 0xf
	return int((0x6996 >> b) & 1)
}

func getHexFromBin(b byte) (byte, error) {
	if b >= 0 && b < 10 {
		return b + 0x30, nil // '0'
	} else if b >= 10 && b < 16 {
		return b - 10 + 0x41, nil // 'A'
	} else {
		return 0x0, fmt.Errorf("only 4bit nibbles allowed")
	}
}

func getBinFromHex(b byte) (byte, error) {
	if b >= 0x30 && b <= 0x39 {
		return b - 0x30, nil // '0'..'9'
	} else if b >= 0x41 && b <= 0x46 {
		return b + 10 - 0x41, nil // 'A'..'F'
	} else if b >= 0x61 && b <= 0x66 {
		return b + 10 - 0x61, nil // 'a'..'f'
	} else if b == 0x0 {
		// special fallback to handle case with uninitialized brake
		return 0, nil
	}

	return 0x0, fmt.Errorf("only hex code characters allowed")
}

func getChecksum(buffer []byte) uint16 {
	shiftreg := uint16(0xc0c1)
	poly := uint16(0xc001)

	for _, a := range buffer {
		tmp := uint16(a) ^ (shiftreg & 0xff)
		shiftreg >>= 8

		if getParity16(tmp) == 1 {
			shiftreg ^= poly
		}

		tmp ^= tmp << 1
		shiftreg ^= tmp << 6
	}

	return shiftreg
}

func DeserializeResponse(response []byte) ([]byte, error) {
	l := len(response)
	if len(response) < 6 || response[0] != startOfFrame || response[l-1] != endOfFrame {
		return []byte{}, fmt.Errorf("invalid frame")
	}

	type checksumPart struct {
		hex   byte
		shift int
	}

	checksumParts := []checksumPart{
		{hex: response[l-5], shift: 4},
		{hex: response[l-4], shift: 0},
		{hex: response[l-3], shift: 12},
		{hex: response[l-2], shift: 8},
	}
	var checksumCalculated uint16
	for _, part := range checksumParts {
		b, err := getBinFromHex(part.hex)
		if err != nil {
			return []byte{}, fmt.Errorf("invalid checksum: %w", err)
		}
		checksumCalculated += uint16(b) << part.shift
	}
	checksumReceived := getChecksum(response[1 : l-5])
	checksumMatches := checksumReceived == checksumCalculated

	fmt.Printf("checksum: matches %v received %#v calculated %#v", checksumMatches, checksumReceived, checksumCalculated)
	if !checksumMatches {
		return []byte{}, fmt.Errorf("checksum does not match")
	}

	// message := response[1:-5]

	deserialized := make([]byte, 0, 32)
	for i := range response[1 : l-5] {
		if i%2 == 0 {
			continue
		}
		thisByte, err := getBinFromHex(response[i])
		if err != nil {
			return []byte{}, fmt.Errorf("unable to parse message: %w", err)
		}
		nextByte, err := getBinFromHex(response[i+1])
		if err != nil {
			return []byte{}, fmt.Errorf("unable to parse message: %w", err)
		}
		deserialized = append(deserialized, thisByte<<4+nextByte)
	}

	return deserialized, nil
}
