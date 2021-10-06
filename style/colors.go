package style

import (
	"encoding/json"
	"fmt"
	"image/color"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/flywave/go-mapbox/style/colorextra"

	"github.com/pkg/errors"
)

const defaultColorAlpha = 0xff

type ColorType struct {
	internalType internalColorType
}

func (c *ColorType) GetColorAtZoomLevel(zoomLevel ZoomLevel) color.Color {
	if c == nil || c.internalType == nil {
		return nil
	}
	return c.internalType.GetValueAtZoomLevel(zoomLevel)
}

func (c *ColorType) UnmarshalJSON(data []byte) error {
	var i interface{}
	err := json.Unmarshal(data, &i)
	if err != nil {
		return err
	}

	switch val := i.(type) {
	case string:
		colorValue, err := strToColor(val, defaultColorAlpha)
		if err != nil {
			return err
		}

		c.internalType = plainColorType{colorValue}
	case map[string]interface{}:
		var colorStops ColorStopsType
		err = json.Unmarshal(data, &colorStops)
		if err != nil {
			return err
		}

		c.internalType = &colorStops
	default:
		panic(fmt.Sprintf("couldn't understand: %T :: %s", i, string(data)))
	}

	return nil
}

type internalColorType interface {
	GetValueAtZoomLevel(zoomLevel ZoomLevel) color.Color
}

type plainColorType struct {
	Color color.Color
}

func (p plainColorType) GetValueAtZoomLevel(zoomLevel ZoomLevel) color.Color {
	return p.Color
}

var (
	strColorHSLRegexp  = regexp.MustCompile(`hsl\(\s*(\d+)\s*,\s*(\d+)%\s*,\s*(\d+)%\s*\)`)
	strColorHSLARegexp = regexp.MustCompile(`hsla\(\s*(\d+)\s*,\s*(\d+)%\s*,\s*(\d+)%\s*,\s*(\d*\.?\d*)\s*\)`)
	strColorRGBRegexp  = regexp.MustCompile(`rgb\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*\)`)
	strColorRGBARegexp = regexp.MustCompile(`rgba\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d*\.?\d*)\s*\)`)
)

func strToColor(str string, defaultAlpha uint8) (color.Color, error) {
	if str == "" {
		return nil, nil
	}

	if strings.HasPrefix(str, "#") {
		return getHexColor(str, defaultAlpha)
	}

	if strColorHSLRegexp.MatchString(str) {
		return getHSLColor(str, defaultAlpha)
	}

	if strColorHSLARegexp.MatchString(str) {
		return getHSLAColor(str)
	}

	if strColorRGBRegexp.MatchString(str) {
		return getRGBColor(str, defaultAlpha)
	}

	if strColorRGBARegexp.MatchString(str) {
		return getRGBAColor(str)
	}

	return nil, errors.Errorf("unimplemented color handing: %q", str)
}

func getHexColor(str string, defaultAlpha uint8) (color.Color, error) {
	var hex6Char string

	switch len(str) {
	case 7:
		hex6Char = str[1:7]
	case 4:
		for _, ch := range str[1:4] {
			hex6Char += (string(ch) + string(ch))
		}
	default:
		return nil, errors.Errorf("unknown hex format: %q", str)
	}

	colorFragments := [3]uint8{}
	for i, match := range []string{hex6Char[0:2], hex6Char[2:4], hex6Char[4:6]} {
		c, err := strconv.ParseUint(match, 16, 64)
		if err != nil {
			return nil, err
		}
		colorFragments[i] = uint8(c)
	}

	return color.RGBA{
		R: colorFragments[0],
		G: colorFragments[1],
		B: colorFragments[2],
		A: defaultAlpha,
	}, nil
}

func getHSLColor(str string, defaultAlpha uint8) (color.Color, error) {
	matches := strColorHSLRegexp.FindStringSubmatch(str)
	colorFragments := [3]float64{}
	for i, match := range matches[1:4] {
		c, err := strconv.ParseFloat(match, 64)
		if err != nil {
			return nil, err
		}
		colorFragments[i] = c
	}

	c, err := colorextra.NewHSLColor(
		colorFragments[0],
		colorFragments[1]/100,
		colorFragments[2]/100,
		defaultAlpha,
	)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func getHSLAColor(str string) (color.Color, error) {
	matches := strColorHSLARegexp.FindStringSubmatch(str)
	colorFragments := [4]float64{}
	for i, match := range matches[1:5] {
		c, err := strconv.ParseFloat(match, 64)
		if err != nil {
			return nil, err
		}
		colorFragments[i] = c
	}

	c, err := colorextra.NewHSLColor(
		colorFragments[0],
		colorFragments[1]/100,
		colorFragments[2]/100,
		uint8(math.Round(colorFragments[3]*256)-1),
	)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func getRGBAColor(str string) (color.Color, error) {
	matches := strColorRGBARegexp.FindStringSubmatch(str)
	colorFragments := [3]uint8{}
	for i, match := range matches[1:4] {
		c, err := strconv.ParseUint(match, 10, 8)
		if err != nil {
			return nil, err
		}
		colorFragments[i] = uint8(c)
	}
	alpha, err := strconv.ParseFloat(matches[4], 64)
	if err != nil {
		return nil, err
	}

	return color.RGBA{
		R: colorFragments[0],
		G: colorFragments[1],
		B: colorFragments[2],
		A: uint8(math.Round(alpha*256) - 1),
	}, nil
}

func getRGBColor(str string, defaultAlpha uint8) (color.Color, error) {
	matches := strColorRGBRegexp.FindStringSubmatch(str)
	colorFragments := [3]uint8{}
	for i, match := range matches[1:4] {
		c, err := strconv.ParseUint(match, 10, 8)
		if err != nil {
			return nil, err
		}
		colorFragments[i] = uint8(c)
	}

	return color.RGBA{
		R: colorFragments[0],
		G: colorFragments[1],
		B: colorFragments[2],
		A: defaultAlpha,
	}, nil
}
