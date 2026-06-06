package style

import (
	"github.com/pkg/errors"
)

type LayerType string

const (
	LayerTypeBackground    LayerType = "background"
	LayerTypeBuilding      LayerType = "building"
	LayerTypeCircle        LayerType = "circle"
	LayerTypeClip          LayerType = "clip"
	LayerTypeFill          LayerType = "fill"
	LayerTypeFillExtrusion LayerType = "fill-extrusion"
	LayerTypeHeatmap       LayerType = "heatmap"
	LayerTypeHillshade     LayerType = "hillshade"
	LayerTypeLine          LayerType = "line"
	LayerTypeModel         LayerType = "model"
	LayerTypeRaster        LayerType = "raster"
	LayerTypeRasterParticle LayerType = "raster-particle"
	LayerTypeSky           LayerType = "sky"
	LayerTypeSlot          LayerType = "slot"
	LayerTypeSymbol        LayerType = "symbol"
)

type Layer struct {
	Appearances []Appearance      `json:"appearances,omitempty"`
	Filter      *FilterContainer `json:"filter,omitempty"`
	ID          string           `json:"id"`
	Layout      *Layout          `json:"layout,omitempty"`
	MaxZoom     *float64         `json:"maxzoom,omitempty"`
	Metadata    Metadata         `json:"metadata,omitempty"`
	MinZoom     *float64         `json:"minzoom,omitempty"`
	Paint       *Paint           `json:"paint,omitempty"`
	Slot        *string          `json:"slot,omitempty"`
	Source      *string          `json:"source,omitempty"`
	SourceLayer *string          `json:"source-layer,omitempty"`
	Type        LayerType        `json:"type"`
}

func validLayerType(t LayerType) bool {
	switch t {
	case LayerTypeBackground, LayerTypeBuilding, LayerTypeCircle, LayerTypeClip,
		LayerTypeFill, LayerTypeFillExtrusion, LayerTypeHeatmap, LayerTypeHillshade,
		LayerTypeLine, LayerTypeModel, LayerTypeRaster, LayerTypeRasterParticle,
		LayerTypeSky, LayerTypeSlot, LayerTypeSymbol:
		return true
	default:
		return false
	}
}

func (l *Layer) Validate() error {
	if l.ID == "" {
		return errors.Errorf("layer id is required")
	}

	if !validLayerType(l.Type) {
		return errors.Errorf("unknown layer type: %q", l.Type)
	}

	if l.MaxZoom != nil && l.MinZoom != nil {
		if *l.MaxZoom < *l.MinZoom {
			return errors.Errorf("max zoom is smaller than min zoom")
		}
	}

	if l.MaxZoom != nil && (*l.MaxZoom < 0 || *l.MaxZoom > 24) {
		return errors.Errorf("max zoom must be between 0 and 24 (inclusive) but was %f", *l.MaxZoom)
	}

	if l.MinZoom != nil && (*l.MinZoom < 0 || *l.MinZoom > 24) {
		return errors.Errorf("min zoom must be between 0 and 24 (inclusive) but was %f", *l.MinZoom)
	}

	return nil
}

type Light struct {
	Anchor    string     `json:"anchor,omitempty"`
	Color     *ColorType `json:"color,omitempty"`
	Intensity float64    `json:"intensity,omitempty"`
	Position  []float64  `json:"position,omitempty"`
}

type Transition struct {
	Delay    int `json:"delay,omitempty"`
	Duration int `json:"duration,omitempty"`
}
