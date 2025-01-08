package tacxserial

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

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
