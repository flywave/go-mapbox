package style

import (
	"image/color"

	"github.com/pkg/errors"
)

type LayerType string

const (
	LayerTypeBackground    LayerType = "background"
	LayerTypeFill          LayerType = "fill"
	LayerTypeLine          LayerType = "line"
	LayerTypeSymbol        LayerType = "symbol"
	LayerTypeRaster        LayerType = "raster"
	LayerTypeCircle        LayerType = "circle"
	LayerTypeFillExtrusion LayerType = "fill-extrusion"
	LayerTypeHeatmap       LayerType = "heatmap"
	LayerTypeHillshade     LayerType = "hillshade"
)

type Layer struct {
	Filter      FilterContainer `json:"filter"`
	ID          LayerType       `json:"id"`
	Layout      Layout          `json:"layout"`
	MaxZoom     *float64        `json:"maxzoom"`
	Metadata    Metadata        `json:"metadata"`
	MinZoom     *float64        `json:"minzoom"`
	Paint       *Paint          `json:"paint"`
	Source      string          `json:"source"`
	SourceLayer string          `json:"source-layer"`
	Type        LayerType       `json:"type"`
}

func (l *Layer) Validate() error {
	if l.MaxZoom != nil && l.MinZoom != nil {
		if *l.MaxZoom < *l.MinZoom {
			return errors.Errorf("max zoom is smaller than min zoom")
		}
	}

	if l.MaxZoom != nil && *l.MaxZoom < 0 || *l.MaxZoom > 24 {
		return errors.Errorf("max zoom must be between 0 and 24 (inclusive) but was %f", *l.MaxZoom)
	}

	if l.MinZoom != nil && *l.MinZoom < 0 || *l.MinZoom > 24 {
		return errors.Errorf("min zoom must be between 0 and 24 (inclusive) but was %f", *l.MinZoom)
	}

	return nil
}

type Light struct {
	Anchor    string      `json:"anchor"`
	Color     color.Color `json:"color"`
	Intensity float64     `json:"intensity"`
	Position  []float64   `json:"position"`
}

type Source struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type Sources map[string]Source

type Transition struct {
	Delay    int `json:"delay"`
	Duration int `json:"duration"`
}
