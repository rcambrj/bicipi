package tacx

import (
	"reflect"
	"testing"
)

type mockCommander struct {
	mockSendCommandOut []byte
	mockSendCommandErr error
}

func (c *mockCommander) sendCommand(command []byte) ([]byte, error) {
	return c.mockSendCommandOut, c.mockSendCommandErr
}

func TestGetVersion(t *testing.T) {
	type test struct {
		description string
		response    []byte
		error       error
		want        version
	}

	tests := []test{
		{
			description: "valid",
			response:    []byte{0x03, 0x0c, 0x00, 0x00, 0x07, 0x10, 0x00, 0x00, 0xae, 0x25, 0x7e, 0x18, 0x15, 0x08, 0x00, 0x00},
			error:       nil,
			want: version{
				model:             "T1941",
				manufactureYear:   2009,
				manufactureNumber: 20366,
				firmwareVersion:   "00.00.10.07",
				serial:            410920366,
				date:              "08-15",
				other:             "00.00",
			},
		},
	}
	for _, tc := range tests {
		mc := &mockCommander{
			mockSendCommandOut: tc.response,
			mockSendCommandErr: tc.error,
		}

		if got, err := getVersion(mc); !reflect.DeepEqual(got, tc.want) || err != nil {
			t.Errorf("getVersion() [%v] => %#v, %v; want %#v", tc.description, got, err, tc.want)
		}
	}
}
