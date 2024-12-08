package tacx

import (
	"reflect"
	"testing"
)

type mockSerialPort struct {
	mockRead []byte
}

func (port mockSerialPort) ResetInputBuffer() error {
	return nil
}
func (port mockSerialPort) Read(p []byte) (n int, err error) {
	copy(p, port.mockRead)
	return len(port.mockRead), nil
}
func (port mockSerialPort) Write(p []byte) (n int, err error) {
	return 0, nil
}

func TestGetVersion(t *testing.T) {
	type test struct {
		response []byte
		want     Version
	}

	tests := []test{
		{
			response: []byte{0x01, 0x30, 0x33, 0x31, 0x33, 0x30, 0x32, 0x30, 0x30, 0x30, 0x46, 0x33, 0x39, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x42, 0x30, 0x35, 0x45, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x46, 0x38, 0x46, 0x44, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x32, 0x36, 0x31, 0x38, 0x31, 0x17},
			want: Version{
				Model:             "T1941",
				ManufactureYear:   2009,
				ManufactureNumber: 20366,
				FirmwareVersion:   "00.00.10.07",
				Serial:            410920366,
				Date:              "08-15",
				Other:             "00.00",
			},
		},
	}
	for _, tc := range tests {
		port := mockSerialPort{
			mockRead: tc.response,
		}

		if got, err := getVersion(port); !reflect.DeepEqual(got, tc.want) || err != nil {
			t.Errorf("getVersion() serial => %#v = %#v, %v; want %#v", tc.response, got, err, tc.want)
		}
	}
}
