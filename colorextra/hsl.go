package colorextra

import (
	"image/color"
	"math"
)

type HSLColor struct {
	H float64 // 0 <= H <= 360
	S float64 // 0 <= S <= 1
	L float64 // 0 <= L <= 1
	A uint8
}

func NewHSLColor(h, s, l float64, a uint8) (HSLColor, error) {
	hslRangeValidations := []RangeValidation{
		{
			Name:       "hue",
			Value:      h,
			LowerBound: 0,
			UpperBound: 360,
		},
		{
			Name:       "saturation",
			Value:      s,
			LowerBound: 0,
			UpperBound: 1,
		},
		{
			Name:       "luminosity",
			Value:      l,
			LowerBound: 0,
			UpperBound: 1,
		},
	}

	for _, rangeValidation := range hslRangeValidations {
		err := AssertInRangeFloat64(rangeValidation)
		if err != nil {
			return HSLColor{}, err
		}
	}

	return HSLColor{h, s, l, a}, nil
}

func NewHSLFromRGB(rgb color.RGBA) HSLColor {
	r := float64(rgb.R) / float64(255)
	g := float64(rgb.G) / float64(255)
	b := float64(rgb.B) / float64(255)

	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))

	h := getHue(rgb)

	l := (max + min) / 2

	if max == min {
		return HSLColor{0, 0, l, rgb.A}
	}

	d := max - min

	var s float64
	if l > 0.5 {
		s = d / (2 - max - min)
	} else {
		s = d / (max + min)
	}

	return HSLColor{h, s, l, rgb.A}
}

func zeroToOneColorToUint32(in float64) uint32 {
	r := uint32(math.Round(255.0 * in))
	r |= r << 8

	return uint32(r)
}

func (hsl HSLColor) RGBA() (uint32, uint32, uint32, uint32) {
	a := ConvertUint8ToUint32Color(hsl.A)
	if hsl.S == 0 {
		// achromatic (gray)

		result := zeroToOneColorToUint32(hsl.L)

		return result, result, result, a
	}

	var q float64
	if hsl.L < 0.5 {
		q = hsl.L * (1.0 + hsl.S)
	} else {
		q = hsl.L + hsl.S - (hsl.L * hsl.S)
	}
	p := (2.0 * hsl.L) - q

	r := hueToRGB(p, q, hsl.H+120)
	g := hueToRGB(p, q, hsl.H)
	b := hueToRGB(p, q, hsl.H-120)

	return zeroToOneColorToUint32(r), zeroToOneColorToUint32(g), zeroToOneColorToUint32(b), a
}

const (
	maxUint24 = 1<<24 - 1
)

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 360
	}
	if t > 360 {
		t -= 360
	}
	if t < 60 {
		// return p + (q-p)*6.0*t
		return p + (q-p)*(t/60.0)
	}
	if t < 180 {
		return q
	}
	if t < 240 {
		return p + ((q - p) * (240.0 - t) / 60.0)
	}
	return p
}
