package mapboxglstyle

import (
	"image/color"
	"reflect"
	"testing"

	"github.com/jamesrr39/goutil/colorextra"
	"github.com/jamesrr39/goutil/errorsx"
)

func Test_strToColor(t *testing.T) {
	type args struct {
		str          string
		defaultAlpha uint8
	}
	tests := []struct {
		name  string
		args  args
		want  color.Color
		want1 errorsx.Error
	}{
		{
			"hex",
			args{
				"#0369CF",
				0,
			},
			color.RGBA{0x03, 0x69, 0xCF, 0},
			nil,
		}, {
			"hex: short",
			args{
				"#3CF",
				0,
			},
			color.RGBA{0x33, 0xCC, 0xFF, 0},
			nil,
		}, {
			"hsl no spaces",
			args{
				"hsl(47,26%,88%)",
				0,
			},
			colorextra.HSLColor{H: 47, S: 0.26, L: 0.88},
			nil,
		}, {
			"hsl with spaces",
			args{
				"hsl(47, 26%, 88%)",
				0,
			},
			colorextra.HSLColor{H: 47, S: 0.26, L: 0.88},
			nil,
		}, {
			"hsla",
			args{
				"hsla(20, 10%, 60%, 0.5)",
				0,
			},
			colorextra.HSLColor{H: 20, S: 0.10, L: 0.60, A: 127},
			nil,
		}, {
			"rgb",
			args{
				"rgb(100, 200, 150)",
				63,
			},
			color.RGBA{R: 100, G: 200, B: 150, A: 63},
			nil,
		}, {
			"rgba",
			args{
				"rgba(100, 200, 150, 1)",
				0,
			},
			color.RGBA{R: 100, G: 200, B: 150, A: 255},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := strToColor(tt.args.str, tt.args.defaultAlpha)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("strToColor() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("strToColor() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
