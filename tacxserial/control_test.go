package tacxserial

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetControlCommandBytes(t *testing.T) {
	type test struct {
		input   controlCommandRaw
		wantOut []byte
		wantErr error
	}

	tests := []test{
		{
			input: controlCommandRaw{
				target:    32767,
				keepalive: 1,
				mode:      2,
				weight:    10,
				adjust:    1040,
			},
			wantOut: []byte{0x1, 0x8, 0x1, 0x0, 0xff, 0x7f, 0x1, 0x0, 0x2, 0x0a, 0x10, 0x4},
			wantErr: nil,
		},
	}
	for _, tc := range tests {
		if out, err := getControlCommandBytes(tc.input); !cmp.Equal(out, tc.wantOut) || !cmp.Equal(err, tc.wantErr) {
			t.Errorf("getControlCommandBytes(%#v) = %#v, %v; want %#v, %v; %v %v", tc.input, out, err, tc.wantOut, tc.wantErr, cmp.Diff(tc.wantOut, out), cmp.Diff(tc.wantErr, err))
		}
	}
}

func TestParseControlResponseBytes(t *testing.T) {
	type test struct {
		input   []byte
		wantOut controlResponseRaw
		wantErr error
	}

	tests := []test{
		{
			input: []byte{3, 19, 2, 0, 192, 2, 0, 0, 44, 24, 24, 91, 32, 1, 76, 3, 0, 0, 1, 0, 27, 0, 2},
			wantOut: controlResponseRaw{
				Distance:    704,
				Speed:       6188,
				AverageLoad: 288,
				CurrentLoad: 844,
				TargetLoad:  0,
				KeepAlive:   1,
				Cadence:     27,
			},
			wantErr: nil,
		},
	}
	for _, tc := range tests {
		if out, err := parseControlResponseBytes(tc.input); !cmp.Equal(out, tc.wantOut) || !cmp.Equal(err, tc.wantErr) {
			t.Errorf("parseControlResponseBytes(%#v) = %#v, %v; want %#v, %v; %v %v", tc.input, out, err, tc.wantOut, tc.wantErr, cmp.Diff(tc.wantOut, out), cmp.Diff(tc.wantErr, err))
		}
	}
}
