package tacxble

import "fmt"

var startOfFrame = 0x01
var endOfFrame = 0x17

// converts the message
// calculates a checksum and appends it
// prepends start of frame
// appends end of frame
func SerializeCommand(message []byte) ([]byte, error) {
	serialized := make([]byte, 4, 32)
	for _, b := range message {
		for _, bshifted := range []byte{b >> 4, b >> 0} {
			h, err := getHex(bshifted)
			if err != nil {
				return nil, fmt.Errorf("unable to serialize command: %w", err)
			}
			serialized = append(serialized, h&0xf)
		}
	}

	checksum := getChecksum(serialized)
	for _, bshifted := range []byte{checksum >> 4, checksum >> 0, checksum >> 12, checksum >> 8} {
		h, err := getHex(bshifted)
		if err != nil {
			return nil, fmt.Errorf("unable to serialize command: %w", err)
		}
		serialized = append(serialized, h&0xf)
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

func getHex(b byte) (byte, error) {
	if b >= 0 && b < 10 {
		return b + 0x30, nil // '0'
	} else if b >= 10 && b < 16 {
		return b - 10 + 0x41, nil // 'A'
	} else {
		return 0x0, fmt.Errorf("only 4bit nibbles allowed")
	}
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
