package tacxserial

import (
	"reflect"
	"testing"
)

func TestGetHexFromBin(t *testing.T) {
	// Test cases where the function should succeed
	cases := []struct {
		Name  string
		Input byte
		Want  byte
	}{
		{"0 should convert to '0'", 0, '0'},
		{"5 should convert to '5'", 5, '5'},
		{"9 should convert to '9'", 9, '9'},
		{"10 should convert to 'A'", 10, 'A'},
		{"15 should convert to 'F'", 15, 'F'},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			got, err := getHexFromBin(tc.Input)
			if err != nil {
				t.Errorf("%s returned an error: %v", tc.Name, err)
			}
			if got != tc.Want {
				t.Errorf("%s = %c, want %c", tc.Name, got, tc.Want)
			}
		})
	}

	// Test error cases
	errorCases := []struct {
		Name  string
		Input byte
		Error string
	}{
		{"Too high", 16, "only 4bit nibbles allowed"},
	}

	for _, tc := range errorCases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := getHexFromBin(tc.Input)
			if err == nil {
				t.Errorf("expected an error for %s, but got none", tc.Name)
			}
			if err != nil && err.Error() != tc.Error {
				t.Errorf("wrong error message: got '%v' want '%v'", err, tc.Error)
			}
		})
	}
}

func TestGetBinFromHex(t *testing.T) {
	// Test cases where the function should succeed
	cases := []struct {
		Name  string
		Input byte
		Want  byte
	}{
		{"'0' should convert to 0", '0', 0},
		{"'5' should convert to 5", '5', 5},
		{"'9 should convert to 9", '9', 9},
		{"'A' should convert to 10", 'A', 10},
		{"'F' should convert to 15", 'F', 15},
		{"'a' should convert to 10", 'A', 10},
		{"'a' should convert to 15", 'F', 15},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			got, err := getBinFromHex(tc.Input)
			if err != nil {
				t.Errorf("%s returned an error: %v", tc.Name, err)
			}
			if got != tc.Want {
				t.Errorf("%s = %c, want %c", tc.Name, got, tc.Want)
			}
		})
	}

	// Test error cases
	errorCases := []struct {
		Name  string
		Input byte
		Error string
	}{
		{"char 01 outside range", 0x01, "only hex code characters allowed"},
		{"char 29 outside range", 0x29, "only hex code characters allowed"},
		{"char 47 outside range", 0x47, "only hex code characters allowed"},
		{"char 60 outside range", 0x60, "only hex code characters allowed"},
		{"char 67 outside range", 0x67, "only hex code characters allowed"},
		{"char 99 outside range", 0x99, "only hex code characters allowed"},
	}

	for _, tc := range errorCases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := getBinFromHex(tc.Input)
			if err == nil {
				t.Errorf("expected an error for %s, but got none", tc.Name)
			}
			if err != nil && err.Error() != tc.Error {
				t.Errorf("wrong error message: got '%v' want '%v'", err, tc.Error)
			}
		})
	}
}

func TestGetParity16(t *testing.T) {
	type test struct {
		input uint16
		want  int
	}

	// Test cases, calculated by checking parity (even - 0, odd - 1)
	testCases := []test{
		{uint16(0), 0},
		{uint16(1), 1},
		{uint16(2), 1},
		{uint16(3), 0},
		{uint16(0xFFFF), 0},
		{uint16(0xAAAB), 1},
		{uint16(0b1110111011101110), 0},
		{uint16(0b1010101010101011), 1},
	}

	for _, tc := range testCases {
		if got := getParity16(tc.input); got != tc.want {
			t.Errorf("getParity16(%#016b) = %v; want %v", tc.input, got, tc.want)
		}
	}
}

func TestGetChecksum(t *testing.T) {
	type test struct {
		input []byte
		want  uint16
	}

	// Test cases where the checksum was calculated for the inputs
	testCases := []test{
		{[]byte{}, 0xc0c1}, // No change when buffer is empty
		{[]byte{0}, 0x9001},
		{[]byte{1}, 0x50c0},
		{[]byte{0x02, 0x00}, 0xa050},
		{[]byte{0xff, 0xff, 0xff, 0xff}, 0x543c},
	}

	for _, tc := range testCases {
		if got := getChecksum(tc.input); got != tc.want {
			t.Errorf("getChecksum(%v) = %#04x; want %#04x", tc.input, got, tc.want)
		}
	}
}

func TestSerializeCommand(t *testing.T) {
	type test struct {
		input []byte
		want  []byte
	}

	tests := []test{
		{
			input: []byte{0x02, 0x00, 0x00, 0x00},
			want:  []byte{0x01, 0x30, 0x32, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x34, 0x33, 0x41, 0x38, 0x17},
		},
		{
			input: []byte{0x01, 0x08, 0x01, 0x00, 0xf8, 0xfd, 0x00, 0x00, 0x02, 0x52, 0x10, 0x04},
			want:  []byte{0x01, 0x30, 0x31, 0x30, 0x38, 0x30, 0x31, 0x30, 0x30, 0x46, 0x38, 0x46, 0x44, 0x30, 0x30, 0x30, 0x30, 0x30, 0x32, 0x35, 0x32, 0x31, 0x30, 0x30, 0x34, 0x44, 0x38, 0x38, 0x30, 0x17},
		},
	}
	for _, tc := range tests {
		if got, err := serializeCommand(tc.input); !reflect.DeepEqual(got, tc.want) || err != nil {
			t.Errorf("SerializeCommand(%#v) = %#v, %v; want %#v", tc.input, got, err, tc.want)
		}
	}
}

func TestDeserializeResponse(t *testing.T) {
	type test struct {
		input []byte
		want  []byte
	}

	tests := []test{
		{
			// version
			input: []byte{0x01, 0x30, 0x33, 0x30, 0x43, 0x30, 0x30, 0x30, 0x30, 0x30, 0x37, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x41, 0x45, 0x32, 0x35, 0x37, 0x45, 0x31, 0x38, 0x31, 0x35, 0x30, 0x38, 0x30, 0x30, 0x30, 0x30, 0x42, 0x30, 0x39, 0x45, 0x17},
			want:  []byte{0x03, 0x0c, 0x00, 0x00, 0x07, 0x10, 0x00, 0x00, 0xae, 0x25, 0x7e, 0x18, 0x15, 0x08, 0x00, 0x00},
		},
		{
			// status
			input: []byte{0x01, 0x30, 0x33, 0x31, 0x33, 0x30, 0x32, 0x30, 0x30, 0x30, 0x46, 0x33, 0x39, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x42, 0x30, 0x35, 0x45, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x46, 0x38, 0x46, 0x44, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x32, 0x36, 0x31, 0x38, 0x31, 0x17},
			want:  []byte{0x03, 0x13, 0x02, 0x00, 0x0f, 0x39, 0x00, 0x00, 0x00, 0x00, 0xb0, 0x5e, 0x00, 0x00, 0x00, 0x00, 0xf8, 0xfd, 0x00, 0x00, 0x00, 0x00, 0x02},
		},
	}
	for _, tc := range tests {
		if got, err := deserializeResponse(tc.input); !reflect.DeepEqual(got, tc.want) || err != nil {
			t.Errorf("DeserializeResponse(%#v) = %#v, %v; want %#v", tc.input, got, err, tc.want)
		}
	}

	// TODO: test malformed messages
}
