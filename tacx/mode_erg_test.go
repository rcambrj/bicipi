package tacx

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetTargetLoad(t *testing.T) {
	type test struct {
		description string
		input       targetLoadArgs
		wantOut     int16
		wantErr     error
	}

	tests := []test{
		{
			description: "zero everything",
			input: targetLoadArgs{
				targetWatts:  0,
				currentSpeed: 0,
			},
			wantOut: 0,
			wantErr: nil,
		},
		{
			description: "no watts, no out",
			input: targetLoadArgs{
				targetWatts:  0,
				currentSpeed: 1000,
			},
			wantOut: 0,
			wantErr: nil,
		},
		{
			description: "zero speed, no out",
			input: targetLoadArgs{
				targetWatts:  100,
				currentSpeed: 0,
			},
			wantOut: 0,
			wantErr: nil,
		},
		{
			description: "slow speed 1, reduced out",
			input: targetLoadArgs{
				targetWatts:  100,
				currentSpeed: 1000,
			},
			wantOut: 357,
			wantErr: nil,
		},
		{
			description: "slow speed 2, reduced out",
			input: targetLoadArgs{
				targetWatts:  100,
				currentSpeed: 2000,
			},
			wantOut: 715,
			wantErr: nil,
		},
		{
			description: "slow speed 3, reduced out",
			input: targetLoadArgs{
				targetWatts:  100,
				currentSpeed: 4000,
			},
			wantOut: 1431,
			wantErr: nil,
		},
		{
			description: "before transition, reduced out",
			input: targetLoadArgs{
				targetWatts:  100,
				currentSpeed: 6000,
			},
			wantOut: 2147, // this value should be very close to the *after* transition value
			wantErr: nil,
		},
		{
			description: "after transition, normal out",
			input: targetLoadArgs{
				targetWatts:  100,
				currentSpeed: 6001,
			},
			wantOut: 2147, // this value should be very close to the *before* transition value
			wantErr: nil,
		},
		{
			description: "full speed, normal out",
			input: targetLoadArgs{
				targetWatts:  100,
				currentSpeed: 8692, // ~30km/h
			},
			wantOut: 1482,
			wantErr: nil,
		},
	}
	for _, tc := range tests {
		if out := getTargetLoad(tc.input); !cmp.Equal(out, tc.wantOut) {
			t.Errorf("[%v] getTargetLoad(%#v) = %#v; want %#v; %v", tc.description, tc.input, out, tc.wantOut, cmp.Diff(tc.wantOut, out))
		}
	}
}
