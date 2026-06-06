package style

import (
	"encoding/json"
	"fmt"
	"image/color"
	"math"
	"regexp"
	"strconv"
	"strings"

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

		c.internalType = plainColorType{Color: colorValue, raw: val}
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

func (c *ColorType) MarshalJSON() ([]byte, error) {
	if c.internalType == nil {
		return json.Marshal(nil)
	}
	switch v := c.internalType.(type) {
	case plainColorType:
		return json.Marshal(v.raw)
	case *ColorStopsType:
		return json.Marshal(v)
	default:
		return json.Marshal(nil)
	}
}

type internalColorType interface {
	GetValueAtZoomLevel(zoomLevel ZoomLevel) color.Color
}

type plainColorType struct {
	Color color.Color
	raw   string
}

func (p plainColorType) GetValueAtZoomLevel(zoomLevel ZoomLevel) color.Color {
	return p.Color
}

var (
	strColorHSLRegexp  = regexp.MustCompile(`hsl\(\s*(\d+)\s*,\s*(\d+)%\s*,\s*(\d+)%\s*\)`)
	strColorHSLARegexp = regexp.MustCompile(`hsla\(\s*(\d+)\s*,\s*(\d+)%\s*,\s*(\d+)%\s*,\s*(\d*\.?\d*)\s*\)`)
	strColorRGBRegexp  = regexp.MustCompile(`rgb\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*\)`)
	strColorRGBARegexp = regexp.MustCompile(`rgba\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d*\.?\d*)\s*\)`)

	namedColors = map[string]string{
		"black": "#000000", "silver": "#c0c0c0", "gray": "#808080",
		"white": "#ffffff", "maroon": "#800000", "red": "#ff0000",
		"purple": "#800080", "fuchsia": "#ff00ff", "green": "#008000",
		"lime": "#00ff00", "olive": "#808000", "yellow": "#ffff00",
		"navy": "#000080", "blue": "#0000ff", "teal": "#008080",
		"aqua": "#00ffff", "orange": "#ffa500", "aliceblue": "#f0f8ff",
		"antiquewhite": "#faebd7", "aquamarine": "#7fffd4", "azure": "#f0ffff",
		"beige": "#f5f5dc", "bisque": "#ffe4c4", "blanchedalmond": "#ffebcd",
		"blueviolet": "#8a2be2", "brown": "#a52a2a", "burlywood": "#deb887",
		"cadetblue": "#5f9ea0", "chartreuse": "#7fff00", "chocolate": "#d2691e",
		"coral": "#ff7f50", "cornflowerblue": "#6495ed", "cornsilk": "#fff8dc",
		"crimson": "#dc143c", "cyan": "#00ffff", "darkblue": "#00008b",
		"darkcyan": "#008b8b", "darkgoldenrod": "#b8860b", "darkgray": "#a9a9a9",
		"darkgreen": "#006400", "darkgrey": "#a9a9a9", "darkkhaki": "#bdb76b",
		"darkmagenta": "#8b008b", "darkolivegreen": "#556b2f", "darkorange": "#ff8c00",
		"darkorchid": "#9932cc", "darkred": "#8b0000", "darksalmon": "#e9967a",
		"darkseagreen": "#8fbc8f", "darkslateblue": "#483d8b", "darkslategray": "#2f4f4f",
		"darkturquoise": "#00ced1", "darkviolet": "#9400d3", "deeppink": "#ff1493",
		"deepskyblue": "#00bfff", "dimgray": "#696969", "dimgrey": "#696969",
		"dodgerblue": "#1e90ff", "firebrick": "#b22222", "floralwhite": "#fffaf0",
		"forestgreen": "#228b22", "gainsboro": "#dcdcdc", "ghostwhite": "#f8f8ff",
		"gold": "#ffd700", "goldenrod": "#daa520", "greenyellow": "#adff2f",
		"grey": "#808080", "honeydew": "#f0fff0", "hotpink": "#ff69b4",
		"indianred": "#cd5c5c", "indigo": "#4b0082", "ivory": "#fffff0",
		"khaki": "#f0e68c", "lavender": "#e6e6fa", "lavenderblush": "#fff0f5",
		"lawngreen": "#7cfc00", "lemonchiffon": "#fffacd", "lightblue": "#add8e6",
		"lightcoral": "#f08080", "lightcyan": "#e0ffff", "lightgoldenrodyellow": "#fafad2",
		"lightgray": "#d3d3d3", "lightgreen": "#90ee90", "lightgrey": "#d3d3d3",
		"lightpink": "#ffb6c1", "lightsalmon": "#ffa07a", "lightseagreen": "#20b2aa",
		"lightskyblue": "#87cefa", "lightslategray": "#778899", "lightslategrey": "#778899",
		"lightsteelblue": "#b0c4de", "lightyellow": "#ffffe0", "limegreen": "#32cd32",
		"linen": "#faf0e6", "magenta": "#ff00ff", "mediumaquamarine": "#66cdaa",
		"mediumblue": "#0000cd", "mediumorchid": "#ba55d3", "mediumpurple": "#9370db",
		"mediumseagreen": "#3cb371", "mediumslateblue": "#7b68ee", "mediumspringgreen": "#00fa9a",
		"mediumturquoise": "#48d1cc", "mediumvioletred": "#c71585", "midnightblue": "#191970",
		"mintcream": "#f5fffa", "mistyrose": "#ffe4e1", "moccasin": "#ffe4b5",
		"navajowhite": "#ffdead", "oldlace": "#fdf5e6", "olivedrab": "#6b8e23",
		"orangered": "#ff4500", "orchid": "#da70d6", "palegoldenrod": "#eee8aa",
		"palegreen": "#98fb98", "paleturquoise": "#afeeee", "palevioletred": "#db7093",
		"papayawhip": "#ffefd5", "peachpuff": "#ffdab9", "peru": "#cd853f",
		"pink": "#ffc0cb", "plum": "#dda0dd", "powderblue": "#b0e0e6",
		"rosybrown": "#bc8f8f", "royalblue": "#4169e1", "saddlebrown": "#8b4513",
		"salmon": "#fa8072", "sandybrown": "#f4a460", "seagreen": "#2e8b57",
		"seashell": "#fff5ee", "sienna": "#a0522d", "skyblue": "#87ceeb",
		"slateblue": "#6a5acd", "slategray": "#708090", "slategrey": "#708090",
		"snow": "#fffafa", "springgreen": "#00ff7f", "steelblue": "#4682b4",
		"tan": "#d2b48c", "thistle": "#d8bfd8", "tomato": "#ff6347",
		"turquoise": "#40e0d0", "violet": "#ee82ee", "wheat": "#f5deb3",
		"whitesmoke": "#f5f5f5", "yellowgreen": "#9acd32",
	}
)

func strToColor(str string, defaultAlpha uint8) (color.Color, error) {
	if str == "" {
		return nil, nil
	}

	switch strings.ToLower(str) {
	case "transparent":
		return color.RGBA{A: 0}, nil
	case "none":
		return nil, nil
	}

	if strings.HasPrefix(str, "#") {
		return getHexColor(str, defaultAlpha)
	}

	if hex, ok := namedColors[strings.ToLower(str)]; ok {
		return getHexColor(hex, defaultAlpha)
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
	case 9:
		hex6Char = str[1:7]
		alphaStr := str[7:9]
		alphaVal, err := strconv.ParseUint(alphaStr, 16, 64)
		if err != nil {
			return nil, err
		}
		defaultAlpha = uint8(alphaVal)
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

// hslToRGBA converts HSL values (hue 0-360, saturation 0-1, lightness 0-1) to RGBA.
func hslToRGBA(h, s, l float64, a uint8) color.RGBA {
	if s == 0 {
		v := uint8(math.Round(l * 255))
		return color.RGBA{R: v, G: v, B: v, A: a}
	}

	var q float64
	if l < 0.5 {
		q = l * (1.0 + s)
	} else {
		q = l + s - (l * s)
	}
	p := (2.0 * l) - q

	r := hueToRGB(p, q, h+120)
	g := hueToRGB(p, q, h)
	b := hueToRGB(p, q, h-120)

	return color.RGBA{
		R: uint8(math.Round(r * 255)),
		G: uint8(math.Round(g * 255)),
		B: uint8(math.Round(b * 255)),
		A: a,
	}
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 360
	}
	if t > 360 {
		t -= 360
	}
	if t < 60 {
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

	if colorFragments[0] < 0 || colorFragments[0] > 360 ||
		colorFragments[1] < 0 || colorFragments[1] > 100 ||
		colorFragments[2] < 0 || colorFragments[2] > 100 {
		return nil, errors.Errorf("hsl value out of range: h=%f s=%f l=%f", colorFragments[0], colorFragments[1], colorFragments[2])
	}

	return hslToRGBA(colorFragments[0], colorFragments[1]/100, colorFragments[2]/100, defaultAlpha), nil
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

	if colorFragments[0] < 0 || colorFragments[0] > 360 ||
		colorFragments[1] < 0 || colorFragments[1] > 100 ||
		colorFragments[2] < 0 || colorFragments[2] > 100 ||
		colorFragments[3] < 0 || colorFragments[3] > 1 {
		return nil, errors.Errorf("hsla value out of range")
	}

	return hslToRGBA(colorFragments[0], colorFragments[1]/100, colorFragments[2]/100,
		uint8(math.Round(colorFragments[3]*256)-1)), nil
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
