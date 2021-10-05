package style

import (
	"encoding/json"
	"image/color"
)

type colorStop struct {
	ZoomLevel ZoomLevel
	Color     color.Color
}

func (c *colorStop) UnmarshalJSON(data []byte) error {
	type colorStopJSONType [2]interface{}

	var colorStopJSON colorStopJSONType
	err := json.Unmarshal(data, &colorStopJSON)
	if err != nil {
		return err
	}

	zoomLevel := ZoomLevel(colorStopJSON[0].(float64))
	colorStr := colorStopJSON[1].(string)
	stopColor, err := strToColor(colorStr, defaultColorAlpha)
	if err != nil {
		return err
	}
	c.ZoomLevel = zoomLevel
	c.Color = stopColor

	return nil
}

type ColorStopsType struct {
	Stops []*colorStop `json:"stops"`
	Base  *float64     `json:"base"`
}

func (c *ColorStopsType) GetValueAtZoomLevel(zoomLevel ZoomLevel) color.Color {
	stopsLen := len(c.Stops)
	if stopsLen == 0 {
		panic("found no stops")
	}

	if zoomLevel < c.Stops[0].ZoomLevel {
		return nil
	}

	for i := 0; i < stopsLen; i++ {
		thisStop := c.Stops[i]
		isLastStop := stopsLen == i+1
		if isLastStop {
			value := thisStop.Color
			return value
		}

		nextStop := c.Stops[i+1]
		if zoomLevel >= nextStop.ZoomLevel {
			continue
		}

		base := 1.0
		if c.Base != nil {
			base = *c.Base
		}

		percentageThrough := getExponentialPercentage(zoomLevel, thisStop.ZoomLevel, nextStop.ZoomLevel, base)

		thisColor := thisStop.Color.(color.RGBA)
		nextColor := nextStop.Color.(color.RGBA)

		value := color.RGBA{
			R: getColorValueBetweenStops(percentageThrough, thisColor.R, nextColor.R),
			G: getColorValueBetweenStops(percentageThrough, thisColor.G, nextColor.G),
			B: getColorValueBetweenStops(percentageThrough, thisColor.B, nextColor.B),
			A: getColorValueBetweenStops(percentageThrough, thisColor.A, nextColor.A),
		}

		return value
	}

	panic("shouldn't get here!")
}

func getColorValueBetweenStops(percentageThrough float64, this, next uint8) uint8 {
	thisFloat64 := float64(this)
	nextFloat64 := float64(next)

	value := getValueThroughStop(thisFloat64, nextFloat64, percentageThrough)

	return uint8(value)
}
