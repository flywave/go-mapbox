package mapboxglstyle

import (
	"testing"

	"github.com/jamesrr39/ownmap-app/ownmap"
	"github.com/stretchr/testify/assert"
)

func Test_strColorHSLRegexp(t *testing.T) {
	assert.True(t, strColorHSLRegexp.MatchString("hsl(47,26%,88%)"))
}

func Test_baseAndStopsType_GetValueAtZoomLevel(t *testing.T) {
	base1Point2 := 1.2

	type fields struct {
		Base  *float64
		Stops []numericStopType
	}
	type args struct {
		zoomLevel ownmap.ZoomLevel
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name: "simple increasing example",
			fields: fields{
				Base: &base1Point2,
				Stops: []numericStopType{
					{14.0, 1.0},
					{20.0, 10.0},
				},
			},
			args: args{
				zoomLevel: ownmap.ZoomLevel(14.16),
			},
			want: 1.134145054348234,
		}, {
			name: "zoom level too low",
			fields: fields{
				Base: &base1Point2,
				Stops: []numericStopType{
					{14.0, 1.0},
					{20.0, 10.0},
				},
			},
			args: args{
				zoomLevel: ownmap.ZoomLevel(13.9),
			},
			want: 0,
		}, {
			name: "zoom level means it falls on last stop",
			fields: fields{
				Base: &base1Point2,
				Stops: []numericStopType{
					{14.0, 1.0},
					{20.0, 10.0},
				},
			},
			args: args{
				zoomLevel: ownmap.ZoomLevel(22),
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NumericFunctionType{
				Base:  tt.fields.Base,
				Stops: tt.fields.Stops,
			}
			if got := s.GetValueAtZoomLevel(tt.args.zoomLevel); got != tt.want {
				t.Errorf("baseAndStopsType.GetValueAtZoomLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}
