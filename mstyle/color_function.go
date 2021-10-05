package mapboxglstyle

import (
	"encoding/json"
	"image/color"

	"github.com/jamesrr39/goutil/errorsx"
	"github.com/jamesrr39/ownmap-app/ownmap"
)

type colorStop struct {
	ZoomLevel ownmap.ZoomLevel
	Color     color.Color
}

func (c *colorStop) UnmarshalJSON(data []byte) error {
	type colorStopJSONType [2]interface{} // float64, string

	var colorStopJSON colorStopJSONType
	err := json.Unmarshal(data, &colorStopJSON)
	if err != nil {
		return errorsx.Wrap(err)
	}

	zoomLevel := ownmap.ZoomLevel(colorStopJSON[0].(float64))
	colorStr := colorStopJSON[1].(string)
	stopColor, err := strToColor(colorStr, defaultColorAlpha)
	if err != nil {
		return errorsx.Wrap(err, "color", colorStr)
	}
	c.ZoomLevel = zoomLevel
	c.Color = stopColor

	return nil
}

type ColorStopsType struct {
	Stops []*colorStop `json:"stops"`
	Base  *float64     `json:"base"`
}

// GetValueAtZoomLevel returns the color at a given zoom level for this item.
// If no color (i.e. not shown), it returns nil
func (c *ColorStopsType) GetValueAtZoomLevel(zoomLevel ownmap.ZoomLevel) color.Color {
	stopsLen := len(c.Stops)
	if stopsLen == 0 {
		panic("found no stops")
	}

	if zoomLevel < c.Stops[0].ZoomLevel {
		// too zoomed out to see this detail
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
			// go to the next stop
			continue
		}

		// this is the correct stop; use this one

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
