package mapboxglstyle

import (
	"encoding/json"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestColorStopsType_UnmarshalJSON(t *testing.T) {
	body := `{
		"stops": [
			[10, "rgba(100, 100, 100, 0)"],
			[18, "rgba(100, 100, 100, 0.5)"]
		]
	}`

	c := new(ColorStopsType)
	err := json.Unmarshal([]byte(body), c)
	require.NoError(t, err)

	halfwayColor := c.GetValueAtZoomLevel(14)
	assert.Equal(t, color.RGBA{100, 100, 100, 191}, halfwayColor)

}

func Test_getColorValueBetweenStops(t *testing.T) {
	type args struct {
		percentageThrough float64
		this              uint8
		next              uint8
	}
	tests := []struct {
		name string
		args args
		want uint8
	}{
		{
			name: "half way forwards",
			args: args{
				percentageThrough: 0.5,
				this:              0,
				next:              255,
			},
			want: 127,
		}, {
			name: "half way backwards",
			args: args{
				percentageThrough: 0.5,
				this:              255,
				next:              0,
			},
			want: 127,
		}, {
			name: "quarter way forwards",
			args: args{
				percentageThrough: 0.25,
				this:              0,
				next:              255,
			},
			want: 63,
		}, {
			name: "quarter way backwards",
			args: args{
				percentageThrough: 0.25,
				this:              255,
				next:              0,
			},
			want: 191,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getColorValueBetweenStops(tt.args.percentageThrough, tt.args.this, tt.args.next); got != tt.want {
				t.Errorf("getColorValueBetweenStops() = %v, want %v", got, tt.want)
			}
		})
	}
}
