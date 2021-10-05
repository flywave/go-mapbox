package colorextra

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*

RGB to HSV color table
Color	Color
name
Hex	(R,G,B)	(H,S,V)
 	Black	#000000	(0,0,0)	(0°,0%,0%)
 	White	#FFFFFF	(255,255,255)	(0°,0%,100%)
 	Red	#FF0000	(255,0,0)	(0°,100%,100%)
 	Lime	#00FF00	(0,255,0)	(120°,100%,100%)
 	Blue	#0000FF	(0,0,255)	(240°,100%,100%)
 	Yellow	#FFFF00	(255,255,0)	(60°,100%,100%)
 	Cyan	#00FFFF	(0,255,255)	(180°,100%,100%)
 	Magenta	#FF00FF	(255,0,255)	(300°,100%,100%)
 	Silver	#C0C0C0	(192,192,192)	(0°,0%,75%)
 	Gray	#808080	(128,128,128)	(0°,0%,50%)
 	Maroon	#800000	(128,0,0)	(0°,100%,50%)
 	Olive	#808000	(128,128,0)	(60°,100%,50%)
 	Green	#008000	(0,128,0)	(120°,100%,50%)
 	Purple	#800080	(128,0,128)	(300°,100%,50%)
 	Teal	#008080	(0,128,128)	(180°,100%,50%)
 	Navy	#000080	(0,0,128)	(240°,100%,50%)

*/

const table = `
 	Black	#000000	(0,0,0)	(0°,0%,0%)
 	White	#FFFFFF	(255,255,255)	(0°,0%,100%)
 	Red	#FF0000	(255,0,0)	(0°,100%,100%)
 	Lime	#00FF00	(0,255,0)	(120°,100%,100%)
 	Blue	#0000FF	(0,0,255)	(240°,100%,100%)
 	Yellow	#FFFF00	(255,255,0)	(60°,100%,100%)
 	Cyan	#00FFFF	(0,255,255)	(180°,100%,100%)
 	Magenta	#FF00FF	(255,0,255)	(300°,100%,100%)
 	Silver	#C0C0C0	(192,192,192)	(0°,0%,75%)
 	Gray	#808080	(128,128,128)	(0°,0%,50%)
 	Maroon	#800000	(128,0,0)	(0°,100%,50%)
 	Olive	#808000	(128,128,0)	(60°,100%,50%)
 	Green	#008000	(0,128,0)	(120°,100%,50%)
 	Purple	#800080	(128,0,128)	(300°,100%,50%)
 	Teal	#008080	(0,128,128)	(180°,100%,50%)
 	Navy	#000080	(0,0,128)	(240°,100%,50%)
`

type testColor struct {
	name string
	rgba color.RGBA
	hsv  HSVColor
}

func generateTestColors() []testColor {
	var colors []testColor

	colors = append(colors,
		testColor{"Black", color.RGBA{0, 0, 0, 0}, HSVColor{0, 0, 0, 0}},
		testColor{"White", color.RGBA{255, 255, 255, 0}, HSVColor{0, 0, 100, 0}},
		testColor{"Red", color.RGBA{255, 0, 0, 0}, HSVColor{0, 100, 100, 0}},
		testColor{"Lime", color.RGBA{0, 255, 0, 0}, HSVColor{120, 100, 100, 0}},
		testColor{"Blue", color.RGBA{0, 0, 255, 0}, HSVColor{240, 100, 100, 0}},
		testColor{"Yellow", color.RGBA{255, 255, 0, 0}, HSVColor{60, 100, 100, 0}},
		testColor{"Cyan", color.RGBA{0, 255, 255, 0}, HSVColor{180, 100, 100, 0}},
		testColor{"Magenta", color.RGBA{255, 0, 255, 0}, HSVColor{300, 100, 100, 0}},
		testColor{"Silver", color.RGBA{192, 192, 192, 0}, HSVColor{0, 0, 75, 0}},
		testColor{"Gray", color.RGBA{128, 128, 128, 0}, HSVColor{0, 0, 50, 0}},
		testColor{"Maroon", color.RGBA{128, 0, 0, 0}, HSVColor{0, 100, 50, 0}},
		testColor{"Olive", color.RGBA{128, 128, 0, 0}, HSVColor{60, 100, 50, 0}},
		testColor{"Green", color.RGBA{0, 128, 0, 0}, HSVColor{120, 100, 50, 0}},
		testColor{"Purple", color.RGBA{128, 0, 128, 0}, HSVColor{300, 100, 50, 0}},
		testColor{"Teal", color.RGBA{0, 128, 128, 0}, HSVColor{180, 100, 50, 0}},
		testColor{"Navy", color.RGBA{0, 0, 128, 0}, HSVColor{240, 100, 50, 0}},
	)

	return colors

}

func round(n float64) int64 {
	if n < 0 {
		return int64(n - 0.5)
	}
	return int64(n + 0.5)
}

func Test_round(t *testing.T) {
	assert.Equal(t, int64(1), round(0.5))
	assert.Equal(t, int64(3), round(3.2))
	assert.Equal(t, int64(0), round(-0.2))
	assert.Equal(t, int64(-1), round(-0.6))
	assert.Equal(t, int64(-3), round(-3.2))
	assert.Equal(t, int64(-4), round(-3.5))
}

func Test_NewHSVFromRGB(t *testing.T) {

	for _, c := range generateTestColors() {
		calculatedHsv := NewHSVFromRGB(c.rgba)

		assert.Equal(t, int64(c.hsv.H), round(calculatedHsv.H), c.name) // round to percentage
		assert.Equal(t, int64(c.hsv.S), round(calculatedHsv.S*100))     // round to percentage
		assert.Equal(t, int64(c.hsv.V), round(calculatedHsv.V*100))     // round to percentage
	}

}
