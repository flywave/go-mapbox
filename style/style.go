package style

import (
	"encoding/json"
	"time"
)

const (
	LayerTypeFill          = "fill"
	LayerTypeLine          = "line"
	LayerTypeSymbol        = "symbol"
	LayerTypeCircle        = "circle"
	LayerTypeFillExtrusion = "fill-extrusion"
	LayerTypeRaster        = "raster"
	LayerTypeBackground    = "background"
	LayerTypeHeatMap       = "heatmap"
	LayerTypeHillShade     = "hillshade"
	StyleVersion           = 8
	LayoutVisible          = "visible"
	LayoutVisibleNone      = "none"
	SourceTypeVector       = "vector"
	SourceTypeRaster       = "raster"
	SourceTypeGeoJSON      = "geojson"
	SourceTypeImage        = "image"
	SourceTypeVideo        = "video"
	SourceTypeCanvas       = "canvas"
)

const (
	LayerMaxZoomMax = 24
	LayerMaxZoomMin = 0
	LayerMinZoomMax = 24
	LayerMinZoomMin = 0
)

type ListStyle struct {
	Version  int64     `json:"version,omitempty"`
	Name     string    `json:"name,omitempty"`
	Created  time.Time `json:"created,omitempty"`
	ID       string    `json:"id,omitempty"`
	Modified time.Time `json:"modified,omitempty"`
	Owner    string    `json:"owner,omitempty"`
}

type Light struct {
	Anchor    string  `json:"anchor,omitempty"`
	Color     string  `json:"color,omitempty"`
	Intensity float64 `json:"intensity,omitempty"`
}

type Transition struct {
	Duration int64
	Delay    int64
}

type Style struct {
	Version    int                        `json:"version,omitempty"`
	Name       string                     `json:"name,omitempty"`
	Metadata   *json.RawMessage           `json:"metadata,omitempty"`
	Center     *[2]float64                `json:"center,omitempty"`
	Zoom       *float64                   `json:"zoom,omitempty"`
	Bearing    *float64                   `json:"bearing,omitempty"`
	Pitch      *float64                   `json:"pitch,omitempty"`
	Light      *Light                     `json:"light,omitempty"`
	Sources    map[string]json.RawMessage `json:"sources,omitempty"`
	Sprite     *string                    `json:"sprite,omitempty"`
	Glyphs     *string                    `json:"glyphs,omitempty"`
	Layers     []json.RawMessage          `json:"layers,omitempty"`
	Transition *Transition                `json:"transition,omitempty"`
	Created    time.Time                  `json:"created,omitempty"`
	Id         string                     `json:"id,omitempty"`
	Modified   time.Time                  `json:"modified,omitempty"`
	Owner      string                     `json:"owner,omitempty"`
	Visibility string
}
