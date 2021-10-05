package style

import "math"

type functionTypeName string

const (
	functionTypeNameIdentity    functionTypeName = "identity"
	functionTypeNameExponential functionTypeName = "exponential"
	functionTypeNameInterval    functionTypeName = "interval"
	functionTypeNameCategory    functionTypeName = "categorical"
)

const defaultFunctionTypeName = functionTypeNameIdentity

func getValueThroughStop(this, next, progressThroughStop float64) float64 {
	diffBetweenValues := next - this
	xProgressThrough := math.Floor(diffBetweenValues * progressThroughStop)

	return xProgressThrough + this
}

func getExponentialPercentage(zoomLevel, lowerStopVal, higherStopVal ZoomLevel, base float64) float64 {
	differenceBetweenLevels := float64(higherStopVal - lowerStopVal)
	if differenceBetweenLevels == 0 {
		return 0
	}

	progressThroughStop := float64(zoomLevel - lowerStopVal)

	if base == 1 {
		return progressThroughStop / differenceBetweenLevels
	}

	top := math.Pow(base, progressThroughStop) - 1
	bottom := math.Pow(base, differenceBetweenLevels) - 1

	return top / bottom
}
