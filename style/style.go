package style

import (
	"encoding/json"
	"image/color"
	"io"

	"github.com/pkg/errors"
)

const (
	ObjectTypePoint      = "Point"
	ObjectTypeLineString = "LineString"
	ObjectTypePolygon    = "Polygon"

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
	Bearing    float64    `json:"bearing"`
	Center     []float64  `json:"center"`
	Glyphs     string     `json:"glyphs"`
	Layers     []*Layer   `json:"layers"`
	Light      Light      `json:"light"`
	Metadata   Metadata   `json:"metadata"`
	Name       string     `json:"name"`
	Pitch      float64    `json:"pitch"`
	Sources    Sources    `json:"sources"`
	Sprite     string     `json:"sprite"`
	Transition Transition `json:"transition"`
	Version    int        `json:"version"`
	Zoom       float64    `json:"zoom"`
	ID         string     `json:"id"`
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
