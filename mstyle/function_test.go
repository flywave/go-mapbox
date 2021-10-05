package mapboxglstyle

import (
	"testing"

	"github.com/jamesrr39/ownmap-app/ownmap"
)

func Test_handleBaseLevel(t *testing.T) {
	stops := []numericStopType{
		{4, 0.25},
		{20, 30},
	}

	type args struct {
		linearValue   ownmap.ZoomLevel
		lowerStopVal  ownmap.ZoomLevel
		higherStopVal ownmap.ZoomLevel
		base          float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "bottom end",
			args: args{
				linearValue:   stops[0].ZoomLevel,
				lowerStopVal:  stops[0].ZoomLevel,
				higherStopVal: stops[1].ZoomLevel,
				base:          1.55,
			},
			want: 0,
		}, {
			name: "bottom end",
			args: args{
				linearValue:   stops[1].ZoomLevel,
				lowerStopVal:  stops[0].ZoomLevel,
				higherStopVal: stops[1].ZoomLevel,
				base:          1.55,
			},
			want: 1,
		}, {
			name: "half way",
			args: args{
				linearValue:   (stops[0].ZoomLevel + stops[1].ZoomLevel) / 2,
				lowerStopVal:  stops[0].ZoomLevel,
				higherStopVal: stops[1].ZoomLevel,
				base:          1,
			},
			want: 0.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getExponentialPercentage(tt.args.linearValue, tt.args.lowerStopVal, tt.args.higherStopVal, tt.args.base); got != tt.want {
				t.Errorf("handleBaseLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}
