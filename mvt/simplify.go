package mvt

import (
	"math"
)

func GetSqSegDist(px, py, x, y, bx, by float64) float64 {
	dx := bx - x
	dy := by - y

	if dx != 0 || dy != 0 {

		t := ((px-x)*dx + (py-y)*dy) / (dx*dx + dy*dy)

		if t > 1 {
			x = bx
			y = by
		} else if t > 0 {
			x += dx * t
			y += dy * t
		}
	}

	dx = px - x
	dy = py - y

	return dx*dx + dy*dy
}

func Simplify(coords []float64, first, last int, sqTolerance float64) {
	maxSqDist := sqTolerance
	mid := (last - first) >> 1
	minPosToMid := last - first
	var index int

	ax := coords[first]
	ay := coords[first+1]
	bx := coords[last]
	by := coords[last+1]

	for i := first + 3; i < last; i += 3 {
		d := GetSqSegDist(coords[i], coords[i+1], ax, ay, bx, by)

		if d > maxSqDist {
			index = i
			maxSqDist = d

		} else if d == maxSqDist {
			posToMid := int(math.Abs(float64(i - mid)))
			if posToMid < minPosToMid {
				index = i
				minPosToMid = posToMid
			}
		}
	}

	if maxSqDist > sqTolerance {
		if index-first > 3 {
			Simplify(coords, first, index, sqTolerance)
		}
		coords[index+2] = maxSqDist
		if last-index > 3 {
			Simplify(coords, index, last, sqTolerance)
		}
	}
}
