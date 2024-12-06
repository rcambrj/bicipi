package tacxble

import (
	"reflect"
	"testing"
)

func TestGetHex(t *testing.T) {
	// Test cases where the function should succeed
	cases := []struct {
		Name  string
		Input byte
		Want  byte
	}{
		{"Zero should convert to '0'", 0, '0'},
		{"Five should convert to '5'", 5, '5'},
		{"Nine should convert to '9'", 9, '9'},
		{"Ten should convert to 'A'", 10, 'A'},
		{"Fifteen should convert to 'F'", 15, 'F'},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			got, err := getHex(c.Input)
			if err != nil {
				t.Errorf("%s returned an error: %v", c.Name, err)
			}
			if got != c.Want {
				t.Errorf("%s = %c, want %c", c.Name, got, c.Want)
			}
		})
	}

	// Test error cases
	errorCases := []struct {
		Name  string
		Input byte
	}{
		{"Too high should error", 16},
	}

	for _, c := range errorCases {
		t.Run(c.Name, func(t *testing.T) {
			_, err := getHex(c.Input)
			if err == nil {
				t.Errorf("expected an error for %s, but got none", c.Name)
			}
			if err != nil && err.Error() != "only 4bit nibbles allowed" {
				t.Errorf("wrong error message: got '%v' want 'only 4bit nibbles allowed'", err)
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
			input:  []byte{0x01, 0x08, 0x01, 0x00, 0xf8, 0xfd, 0x00, 0x00, 0x02, 0x52, 0x10, 0x04},
			want: []byte{0x01, 0x30, 0x31, 0x30, 0x38, 0x30, 0x31, 0x30, 0x30, 0x46, 0x38, 0x46, 0x44, 0x30, 0x30, 0x30, 0x30, 0x30, 0x32, 0x35, 0x32, 0x31, 0x30, 0x30, 0x34, 0x44, 0x38, 0x38, 0x30, 0x17},
		},

	for _, tc := range tests {
		if got, err := SerializeCommand(tc.input); !reflect.DeepEqual(got, tc.want) || err != nil {
			t.Errorf("SerializeCommand(%#v) = %#v, %v; want %#v", tc.input, got, err, tc.want)
		}
	}
}
