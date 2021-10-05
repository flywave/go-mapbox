package style

import (
	"encoding/json"
	"math"

	"github.com/pkg/errors"
)

type NumberOrFunctionWrapperType struct {
	internalType numberOfBaseAndStopsType
}

func (n *NumberOrFunctionWrapperType) GetValueAtZoomLevel(zoomLevel ZoomLevel) float64 {
	if n == nil || n.internalType == nil {
		return 0
	}
	return n.internalType.GetValueAtZoomLevel(zoomLevel)
}

func (n *NumberOrFunctionWrapperType) UnmarshalJSON(data []byte) error {
	var i interface{}

	err := json.Unmarshal(data, &i)
	if err != nil {
		return err
	}

	switch val := i.(type) {
	case float64:
		n.internalType = plainNumberType(val)
		return nil
	case map[string]interface{}:
		internalType := new(NumericFunctionType)
		err = json.Unmarshal(data, internalType)
		if err != nil {
			return errors.Wrap(err, string(data))
		}

		n.internalType = internalType

		return nil
	}

	panic("unknown type??")
}

type plainNumberType float64

func (p plainNumberType) GetValueAtZoomLevel(zoomLevel ZoomLevel) float64 {
	return float64(p)
}

type numberOfBaseAndStopsType interface {
	GetValueAtZoomLevel(zoomLevel ZoomLevel) float64
}

type NumericFunctionType struct {
	Type  functionTypeName  `json:"type"`
	Base  *float64          `json:"base"`
	Stops []numericStopType `json:"stops"`
}

func (s NumericFunctionType) GetValueAtZoomLevel(zoomLevel ZoomLevel) float64 {
	stopsLen := len(s.Stops)
	if stopsLen == 0 {
		panic("no stops found")
	}

	if zoomLevel < s.Stops[0].ZoomLevel {
		return 0
	}

	for i := 0; i < stopsLen; i++ {
		thisStop := s.Stops[i]
		isLastStop := stopsLen == i+1
		if isLastStop {
			return thisStop.Value
		}

		nextStop := s.Stops[i+1]
		nextStopZoomLevel := nextStop.ZoomLevel
		if zoomLevel >= nextStopZoomLevel {
			continue
		}

		return getNumericValueBetweenStops(thisStop, nextStop, zoomLevel, s.Base)
	}

	panic("shouldn't get here!")
}

func getNumericValueBetweenStops(thisStop, nextStop numericStopType, zoomLevel ZoomLevel, base *float64) float64 {
	if base == nil {
		one := 1.0
		base = &one
	}

	progressThroughStop := getExponentialPercentage(zoomLevel, thisStop.ZoomLevel, nextStop.ZoomLevel, *base)

	value := (progressThroughStop * math.Abs(nextStop.Value-thisStop.Value)) + thisStop.Value

	return value
}
