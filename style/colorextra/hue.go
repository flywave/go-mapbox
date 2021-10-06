package colorextra

import (
	"image/color"
	"math"
)

func getHue(rgb color.RGBA) float64 {
	// h
	var h float64
	/*

		if G >= B
			H = cos -1 [ (R - 1⁄2G - 1⁄2B)/√R2 + G2 + B2 - RG - RB - GB ]
		else
			H = 360 -	cos -1 [	(R - 1⁄2G - 1⁄2B)/√R2 + G2 + B2 - RG - RB - GB ]
	*/

	r := float64(rgb.R)
	g := float64(rgb.G)
	b := float64(rgb.B)
	underSqRoot := math.Pow(r, 2) + math.Pow(g, 2) + math.Pow(b, 2) - ((r * g) + (r * b) + (g * b))
	beforeSqRoot := r - ((g / 2) + (b / 2))
	var inBrackets float64

	// avoid dividing by 0
	if underSqRoot > 0 {
		inBrackets = beforeSqRoot / math.Sqrt(underSqRoot)
	} else {
		inBrackets = 1
	}
	arccosRadians := math.Acos(inBrackets)

	// pi radians = 180 deg
	arccosDeg := arccosRadians * 180 / math.Pi
	if g >= b {
		h = arccosDeg
	} else {
		h = 360 - arccosDeg
	}

	return h
}
