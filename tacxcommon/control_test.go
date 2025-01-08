package tacxcommon

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetControlCommandBytes(t *testing.T) {
	type test struct {
		input   ControlCommand
		wantOut []byte
		wantErr error
	}

	tests := []test{
		{
			input: ControlCommand{
				TargetSpeed: 5432,
				Keepalive:   1,
				Mode:        ModeCalibrating,
				Weight:      10,
				Adjust:      1040,
			},
			wantOut: []byte{0x1, 0x8, 0x1, 0x0, 0x38, 0x15, 0x1, 0x0, 0x3, 0x0a, 0x10, 0x4},
			wantErr: nil,
		},
		{
			input: ControlCommand{
				TargetLoad: 32767,
				Keepalive:  1,
				Mode:       ModeNormal,
				Weight:     10,
				Adjust:     1040,
			},
			wantOut: []byte{0x1, 0x8, 0x1, 0x0, 0xff, 0x7f, 0x1, 0x0, 0x2, 0x0a, 0x10, 0x4},
			wantErr: nil,
		},
	}
	for _, tc := range tests {
		if out, err := GetControlCommandBytes(tc.input); !cmp.Equal(out, tc.wantOut) || !cmp.Equal(err, tc.wantErr) {
			t.Errorf("GetControlCommandBytes(%#v) = %#v, %v; want %#v, %v; %v %v", tc.input, out, err, tc.wantOut, tc.wantErr, cmp.Diff(tc.wantOut, out), cmp.Diff(tc.wantErr, err))
		}
	}
}

func TestGetControlResponseFromBytes(t *testing.T) {
	type test struct {
		input   []byte
		wantOut ControlResponse
		wantErr error
	}

	tests := []test{
		{
			input: []byte{3, 19, 2, 0, 192, 2, 0, 0, 44, 24, 24, 91, 32, 1, 76, 3, 0, 0, 1, 0, 27, 0, 2},
			wantOut: ControlResponse{
				Speed:       6188,
				CurrentLoad: 844,
				TargetLoad:  0,
				Keepalive:   1,
				Cadence:     27,
			},
			wantErr: nil,
		},
	}
	for _, tc := range tests {
		if out, err := GetControlResponseFromBytes(tc.input); !cmp.Equal(out, tc.wantOut) || !cmp.Equal(err, tc.wantErr) {
			t.Errorf("GetControlResponseFromBytes(%#v) = %#v, %v; want %#v, %v; %v %v", tc.input, out, err, tc.wantOut, tc.wantErr, cmp.Diff(tc.wantOut, out), cmp.Diff(tc.wantErr, err))
		}
	}
}
