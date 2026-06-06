package style

import (
	"encoding/json"
	"image/color"
	"io"

	"github.com/pkg/errors"
)

const (
	backgroundLayerID = "background"
)

func (s *MapboxGLStyle) GetStyleID() string {
	return s.style.ID
}

func (s *MapboxGLStyle) GetBackground() color.Color {
	return s.backgroundColor
}

func (s *Style) calculateBackgroundColor() (color.Color, error) {
	if len(s.Layers) == 0 {
		return color.White, nil
	}

	backgroundLayer := s.Layers[0]
	if backgroundLayer.ID != backgroundLayerID {
		return color.White, nil
	}

	return backgroundLayer.Paint.BackgroundColor.GetColorAtZoomLevel(0), nil
}

type MapboxGLStyle struct {
	style           *Style
	backgroundColor color.Color
}

type Style struct {
	Bearing    float64                `json:"bearing,omitempty"`
	Camera     *Camera                `json:"camera,omitempty"`
	Center     []float64              `json:"center,omitempty"`
	ColorTheme *ColorTheme            `json:"color-theme,omitempty"`
	Featuresets map[string]Featureset `json:"featuresets,omitempty"`
	Fog        *Fog                   `json:"fog,omitempty"`
	Fragment   *bool                  `json:"fragment,omitempty"`
	Glyphs     string                 `json:"glyphs,omitempty"`
	Iconsets   map[string]Iconset     `json:"iconsets,omitempty"`
	ID         string                 `json:"id,omitempty"`
	Imports    []Import               `json:"imports,omitempty"`
	Layers     []*Layer               `json:"layers"`
	Light      *Light                 `json:"light,omitempty"`
	Lights     []Light3D              `json:"lights,omitempty"`
	Metadata   Metadata               `json:"metadata,omitempty"`
	Models     map[string]string      `json:"models,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Pitch      float64                `json:"pitch,omitempty"`
	Projection *Projection            `json:"projection,omitempty"`
	Rain       *Rain                  `json:"rain,omitempty"`
	Schema     map[string]SchemaOption `json:"schema,omitempty"`
	Snow       *Snow                  `json:"snow,omitempty"`
	Sources    Sources                `json:"sources"`
	Sprite     string                 `json:"sprite,omitempty"`
	Terrain    *Terrain               `json:"terrain,omitempty"`
	Transition *Transition            `json:"transition,omitempty"`
	Version    int                    `json:"version"`
	Zoom       float64                `json:"zoom,omitempty"`
}

type Featureset struct {
	Metadata  Metadata    `json:"metadata,omitempty"`
	Selectors []Selector  `json:"selectors,omitempty"`
}

type Selector struct {
	Layer             string                 `json:"layer"`
	FeatureNamespace  string                 `json:"featureNamespace,omitempty"`
	Properties        map[string]interface{} `json:"properties,omitempty"`
}

type Iconset struct {
	Type   string `json:"type,omitempty"`
	URL    string `json:"url,omitempty"`
	Source string `json:"source,omitempty"`
}

type Rain struct {
	CenterThinning  *float64   `json:"center-thinning,omitempty"`
	Color           *ColorType `json:"color,omitempty"`
	Density         *float64   `json:"density,omitempty"`
	Direction       *float64   `json:"direction,omitempty"`
	DistortionStrength *float64 `json:"distortion-strength,omitempty"`
	DropletSize     *float64   `json:"droplet-size,omitempty"`
	Intensity       *float64   `json:"intensity,omitempty"`
	Opacity         *float64   `json:"opacity,omitempty"`
	Vignette        *float64   `json:"vignette,omitempty"`
	VignetteColor   *ColorType `json:"vignette-color,omitempty"`
}

type Snow struct {
	CenterThinning *float64   `json:"center-thinning,omitempty"`
	Color          *ColorType `json:"color,omitempty"`
	Density        *float64   `json:"density,omitempty"`
	Direction      *float64   `json:"direction,omitempty"`
	FlakeSize      *float64   `json:"flake-size,omitempty"`
	Intensity      *float64   `json:"intensity,omitempty"`
	Opacity        *float64   `json:"opacity,omitempty"`
	Vignette       *float64   `json:"vignette,omitempty"`
	VignetteColor  *ColorType `json:"vignette-color,omitempty"`
}

func Parse(reader io.Reader) (*MapboxGLStyle, error) {
	s := new(Style)
	dec := json.NewDecoder(reader)
	err := dec.Decode(s)
	if err != nil {
		return nil, err
	}

	bgColor, err := s.calculateBackgroundColor()
	if err != nil {
		return nil, err
	}

	return &MapboxGLStyle{s, bgColor}, nil
}

func (s *Style) Validate() error {
	const expectedVersion = 8
	if s.Version != expectedVersion {
		return errors.Errorf("version: expected %d but was %d", expectedVersion, s.Version)
	}

	return nil
}
