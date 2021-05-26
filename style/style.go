package style

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

type Style struct {
	Id         string                 `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Version    int                    `json:"version"`
	Center     [2]float64             `json:"center,omitempty"`
	Zoom       float64                `json:"zoom,omitempty"`
	Bearing    int64                  `json:"bearing,omitempty"`
	Pitch      int64                  `json:"pitch,omitempty"`
	Transition *StyleTransition       `json:"transition,omitempty"`
	Created    string                 `json:"created,omitempty"`
	Modified   string                 `json:"modified,omitempty"`
	Owner      string                 `json:"owner,omitempty"`
	Sources    map[string]StyleSource `json:"sources"`
	Layers     []StyleLayer           `json:"layers"`
}

type StyleTransition struct {
	Duration float32 `json:"duration"`
	Delay    float32 `json:"delay"`
}

type StyleLight struct {
	Anchor    string     `json:"anchor"`
	Color     string     `json:"color"`
	Intensity float32    `json:"intensity"`
	Position  [3]float32 `json:"position"`
}

type StyleSource struct {
	Type        string             `json:"type"`
	Tiles       []string           `json:"tiles,omitempty"`
	Url         string             `json:"url,omitempty"`
	MinZoom     float32            `json:"min_zoom,omitempty"`
	MaxZoom     float32            `json:"max_zoom,omitempty"`
	Bounds      [4]float32         `json:"bounds,omitempty"`
	Scheme      string             `json:"scheme,omitempty"` // xyz or tms
	Attribution *map[string]string `json:"attribution,omitempty"`
	TileSize    int                `json:"tileSize,omitempty"`
	Encoding    string             `json:"encoding,omitempty"`
	Data        string             `json:"data,omitempty"`
}

type StyleLayer struct {
	Source      string                 `json:"source,omitempty"`
	SourceLayer string                 `json:"source-layer,omitempty"`
	Type        string                 `json:"type,omitempty"`
	MinZoom     float32                `json:"min_zoom"`
	MaxZoom     float32                `json:"max_zoom"`
	Filter      *StyleFilter           `json:"filter"`
	Layout      *StyleLayout           `json:"layout"`
	Paint       *StylePaint            `json:"paint"`
	Interactive bool                   `json:"interactive,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Ref         string                 `json:"ref,omitempty"`
}

type StyleFilter [][]string

type StyleLayout map[string]string

type StylePaint map[string]string
