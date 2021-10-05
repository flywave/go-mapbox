package colorextra

import "fmt"

type RangeValidation struct {
	Name       string
	Value      float64
	LowerBound float64
	UpperBound float64
}

func AssertInRangeFloat64(rangeValidation RangeValidation) error {
	if rangeValidation.Value < rangeValidation.LowerBound || rangeValidation.Value > rangeValidation.UpperBound {
		return fmt.Errorf("expected %q value to be between %v and %v, but was %v", rangeValidation.Name, rangeValidation.LowerBound, rangeValidation.UpperBound, rangeValidation.Value)
	}
	return nil
}
