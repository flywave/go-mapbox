package mapboxglstyle

import (
	"testing"

	"github.com/jamesrr39/ownmap-app/ownmap"
)

func Test_getNumericValueBetweenStops(t *testing.T) {
	stops := []numericStopType{
		{4, 0.25},
		{20, 30},
	}

	base := float64(1)

	type args struct {
		thisStop  numericStopType
		nextStop  numericStopType
		zoomLevel ownmap.ZoomLevel
		base      *float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "bottom end",
			args: args{
				thisStop:  stops[0],
				nextStop:  stops[1],
				zoomLevel: stops[0].ZoomLevel,
				base:      &base,
			},
			want: stops[0].Value,
		}, {
			name: "top end",
			args: args{
				thisStop:  stops[0],
				nextStop:  stops[1],
				zoomLevel: stops[1].ZoomLevel,
				base:      &base,
			},
			want: stops[1].Value,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNumericValueBetweenStops(tt.args.thisStop, tt.args.nextStop, tt.args.zoomLevel, tt.args.base); got != tt.want {
				t.Errorf("getNumericValueBetweenStops() = %v, want %v", got, tt.want)
			}
		})
	}
}
