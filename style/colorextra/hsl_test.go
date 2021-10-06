package colorextra

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHSLFromRGB(t *testing.T) {
	for _, testColor := range testColors {
		hsl := NewHSLFromRGB(testColor.RGB)
		assert.InDelta(t, testColor.HSL.H, hsl.H, 0.01, fmt.Sprintf("color: %q", testColor.Name))
		assert.InDelta(t, testColor.HSL.S, hsl.S, 0.01, fmt.Sprintf("color: %q", testColor.Name))
		assert.InDelta(t, testColor.HSL.L, hsl.L, 0.01, fmt.Sprintf("color: %q", testColor.Name))
	}
}

type colorTableEntry struct {
	Name string
	HSL  HSLColor
	RGB  color.RGBA
}

var testColors = []colorTableEntry{
	{Name: "Black", HSL: HSLColor{H: 0.000000, S: 0.000000, L: 0.000000}, RGB: color.RGBA{R: 0, G: 0, B: 0}},
	{Name: "White", HSL: HSLColor{H: 0.000000, S: 0.000000, L: 1.000000}, RGB: color.RGBA{R: 255, G: 255, B: 255}},
	{Name: "Red", HSL: HSLColor{H: 0.000000, S: 1.000000, L: 0.500000}, RGB: color.RGBA{R: 255, G: 0, B: 0}},
	{Name: "Lime", HSL: HSLColor{H: 120.000000, S: 1.000000, L: 0.500000}, RGB: color.RGBA{R: 0, G: 255, B: 0}},
	{Name: "Blue", HSL: HSLColor{H: 240.000000, S: 1.000000, L: 0.500000}, RGB: color.RGBA{R: 0, G: 0, B: 255}},
	{Name: "Yellow", HSL: HSLColor{H: 60.000000, S: 1.000000, L: 0.500000}, RGB: color.RGBA{R: 255, G: 255, B: 0}},
	{Name: "Cyan", HSL: HSLColor{H: 180.000000, S: 1.000000, L: 0.500000}, RGB: color.RGBA{R: 0, G: 255, B: 255}},
	{Name: "Magenta", HSL: HSLColor{H: 300.000000, S: 1.000000, L: 0.500000}, RGB: color.RGBA{R: 255, G: 0, B: 255}},
	{Name: "Silver", HSL: HSLColor{H: 0.000000, S: 0.000000, L: 0.750000}, RGB: color.RGBA{R: 191, G: 191, B: 191}},
	{Name: "Gray", HSL: HSLColor{H: 0.000000, S: 0.000000, L: 0.500000}, RGB: color.RGBA{R: 128, G: 128, B: 128}},
	{Name: "Maroon", HSL: HSLColor{H: 0.000000, S: 1.000000, L: 0.250000}, RGB: color.RGBA{R: 128, G: 0, B: 0}},
	{Name: "Olive", HSL: HSLColor{H: 60.000000, S: 1.000000, L: 0.250000}, RGB: color.RGBA{R: 128, G: 128, B: 0}},
	{Name: "Green", HSL: HSLColor{H: 120.000000, S: 1.000000, L: 0.250000}, RGB: color.RGBA{R: 0, G: 128, B: 0}},
	{Name: "Purple", HSL: HSLColor{H: 300.000000, S: 1.000000, L: 0.250000}, RGB: color.RGBA{R: 128, G: 0, B: 128}},
	{Name: "Teal", HSL: HSLColor{H: 180.000000, S: 1.000000, L: 0.250000}, RGB: color.RGBA{R: 0, G: 128, B: 128}},
	{Name: "Navy", HSL: HSLColor{H: 240.000000, S: 1.000000, L: 0.250000}, RGB: color.RGBA{R: 0, G: 0, B: 128}},
}

func TestHSLColor_RGBA(t *testing.T) {
	for _, testColor := range testColors {
		r, g, b, _ := testColor.HSL.RGBA()
		expectedR, expectedG, expectedB, _ := testColor.RGB.RGBA()
		assert.Equal(t, expectedR, r, fmt.Sprintf("color: %q, red component, %d :: %d", testColor.Name, r, expectedR))
		assert.Equal(t, expectedG, g, fmt.Sprintf("color: %q, green component", testColor.Name))
		assert.Equal(t, expectedB, b, fmt.Sprintf("color: %q, blue component", testColor.Name))
	}

	t.Run("alpha", func(t *testing.T) {
		alpha, err := NewHSLColor(0, 0, 0, 0xff)
		require.NoError(t, err)

		_, _, _, a := alpha.RGBA()
		assert.Equal(t, uint32(65535), a)
	})
}
