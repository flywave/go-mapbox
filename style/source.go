package style

type SourceType string

const (
	SourceTypeVector  SourceType = "vector"
	SourceTypeRaster  SourceType = "raster"
	SourceTypeGeoJSON SourceType = "geojson"
	SourceTypeImage   SourceType = "image"
	SourceTypeVideo   SourceType = "video"
	SourceTypeCanvas  SourceType = "canvas"
)

type Source struct {
	Type    string   `json:"type"`
	Tiles   []string `json:"tiles,omitempty"`
	MinZoom int      `json:"minzoom,omitempty"`
	MaxZoom int      `json:"maxzoom,omitempty"`
	URL     string   `json:"url,omitempty"`
}
