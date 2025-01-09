package common

import (
	"reflect"
	"testing"
)

func TestGetVersion(t *testing.T) {
	type test struct {
		description string
		response    []byte
		error       error
		want        Version
	}

	tests := []test{
		{
			description: "valid",
			response:    []byte{0x03, 0x0C, 0x00, 0x00, 0x07, 0x10, 0x00, 0x00, 0xae, 0x25, 0x7e, 0x18, 0x15, 0x08, 0x00, 0x00},
			error:       nil,
			want: Version{
				Model:             "T1941",
				ManufactureNumber: 20366,
				FirmwareVersion:   "00.00.10.07",
				Serial:            410920366, // sticker shows 41092366 though?
				Date:              "2009-08-15",
			},
		},
	}
	for _, tc := range tests {

		if got, err := GetVersionFromResponseBytes(tc.response); !reflect.DeepEqual(got, tc.want) || err != nil {
			t.Errorf("getVersion() [%v] => %#v, %v; want %#v", tc.description, got, err, tc.want)
		}
	}
}
