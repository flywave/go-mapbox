package tilejson

import (
	"fmt"
	"time"
)

// TileJSON version constants.
const (
	SpecVersion    = "3.0.0"
	DefaultVersion = "1.0.0"
)

// Tileset type constants.
const (
	VectorTileset TilesetType = "vector"
	RasterTileset TilesetType = "raster"
)

// Geometry type constants.
const (
	GeomTypePoint   GeomType = "point"
	GeomTypeLine    GeomType = "line"
	GeomTypePolygon GeomType = "polygon"
	GeomTypeUnknown GeomType = "unknown"
)

// Scheme constants.
const (
	SchemeXYZ = "xyz"
	SchemeTMS = "tms"
)

// Field type constants.
const (
	FieldTypeString  FieldType = "String"
	FieldTypeNumber  FieldType = "Number"
	FieldTypeBoolean FieldType = "Boolean"
)

// Accessor constants.
const (
	SortByCreated  SortBy = "created"
	SortByModified SortBy = "modified"
)

// Visibility constants.
const (
	PublicTileset  TilesetVisibility = "public"
	PrivateTileset TilesetVisibility = "private"
)

type (
	SortBy             string
	TilesetType        string
	TilesetVisibility  string
	GeomType           string
	FieldType          string
)

// TileJSON implements the TileJSON 3.0.0 specification.
// See https://github.com/mapbox/tilejson-spec
type TileJSON struct {
	TileJSON     string         `json:"tilejson"`
	Tiles        []string       `json:"tiles"`
	VectorLayers []VectorLayer  `json:"vector_layers,omitempty"`

	Attribution  *string        `json:"attribution,omitempty"`
	Bounds       *[4]float64    `json:"bounds,omitempty"`
	Center       *[3]float64    `json:"center,omitempty"`
	Data         []string       `json:"data,omitempty"`
	Description  *string        `json:"description,omitempty"`
	FillZoom     *int           `json:"fillzoom,omitempty"`
	Grids        []string       `json:"grids,omitempty"`
	Legend       *string        `json:"legend,omitempty"`
	MaxZoom      *int           `json:"maxzoom,omitempty"`
	MinZoom      *int           `json:"minzoom,omitempty"`
	Name         *string        `json:"name,omitempty"`
	Scheme       string         `json:"scheme,omitempty"`
	Template     *string        `json:"template,omitempty"`
	Version      string         `json:"version,omitempty"`
}

// VectorLayer describes a layer of vector tile data.
type VectorLayer struct {
	ID          string              `json:"id"`
	Fields      map[string]string   `json:"fields"`
	Description *string             `json:"description,omitempty"`
	MinZoom     *int                `json:"minzoom,omitempty"`
	MaxZoom     *int                `json:"maxzoom,omitempty"`
}

// Tileset is a Mapbox API tileset listing entry.
type Tileset struct {
	Type        string     `json:"type,omitempty"`
	Center      [3]float64 `json:"center,omitempty"`
	Created     time.Time  `json:"created,omitempty"`
	Description string     `json:"description,omitempty"`
	Filesize    int64      `json:"filesize,omitempty"`
	ID          string     `json:"id,omitempty"`
	Modified    time.Time  `json:"modified,omitempty"`
	Name        string     `json:"name,omitempty"`
	Visibility  string     `json:"visibility,omitempty"`
	Status      string     `json:"status,omitempty"`
}

// ListTilesetsParams defines parameters for listing tilesets.
type ListTilesetsParams struct {
	Type       *TilesetType
	Visibility *TilesetVisibility
	SortBy     *SortBy
	Limit      *int
}

// ─── Constructors ──────────────────────────────────────────────────────────

// New creates a TileJSON with the required tiles array.
func New(tiles []string) *TileJSON {
	return &TileJSON{
		TileJSON: SpecVersion,
		Tiles:    tiles,
		Scheme:   SchemeXYZ,
		Version:  DefaultVersion,
	}
}

// NewTileBounds creates a bounds array [west, south, east, north].
func NewTileBounds(minLon, minLat, maxLon, maxLat float64) *[4]float64 {
	return &[4]float64{minLon, minLat, maxLon, maxLat}
}

// NewTileCenter creates a center array [lon, lat, zoom].
func NewTileCenter(lon, lat float64, zoom int) *[3]float64 {
	return &[3]float64{lon, lat, float64(zoom)}
}

// AddTile appends a tile endpoint.
func (t *TileJSON) AddTile(url string) {
	t.Tiles = append(t.Tiles, url)
}

// AddVectorLayer adds a vector layer description.
func (t *TileJSON) AddVectorLayer(vl VectorLayer) {
	t.VectorLayers = append(t.VectorLayers, vl)
}

// NewVectorLayer creates a VectorLayer with the required id and fields.
func NewVectorLayer(id string, fields map[string]string) VectorLayer {
	return VectorLayer{ID: id, Fields: fields}
}

// ─── Validation ────────────────────────────────────────────────────────────

func (t *TileJSON) Validate() error {
	if t.TileJSON != SpecVersion {
		return fmt.Errorf("tilejson: expected %q, got %q", SpecVersion, t.TileJSON)
	}
	if len(t.Tiles) == 0 {
		return fmt.Errorf("tiles: must contain at least one endpoint")
	}
	if t.MinZoom != nil && t.MaxZoom != nil && *t.MinZoom > *t.MaxZoom {
		return fmt.Errorf("minzoom (%d) > maxzoom (%d)", *t.MinZoom, *t.MaxZoom)
	}
	if t.MinZoom != nil && (*t.MinZoom < 0 || *t.MinZoom > 30) {
		return fmt.Errorf("minzoom (%d) out of range [0, 30]", *t.MinZoom)
	}
	if t.MaxZoom != nil && (*t.MaxZoom < 0 || *t.MaxZoom > 30) {
		return fmt.Errorf("maxzoom (%d) out of range [0, 30]", *t.MaxZoom)
	}
	if t.Scheme != "" && t.Scheme != SchemeXYZ && t.Scheme != SchemeTMS {
		return fmt.Errorf("scheme: expected %q or %q, got %q", SchemeXYZ, SchemeTMS, t.Scheme)
	}
	if t.Center != nil {
		if len(t.Center) != 3 {
			return fmt.Errorf("center: expected [lon, lat, zoom]")
		}
	}
	if t.Bounds != nil {
		if len(t.Bounds) != 4 {
			return fmt.Errorf("bounds: expected [west, south, east, north]")
		}
	}
	for i, vl := range t.VectorLayers {
		if vl.ID == "" {
			return fmt.Errorf("vector_layers[%d]: id is required", i)
		}
		if vl.Fields == nil {
			return fmt.Errorf("vector_layers[%d] (%q): fields is required", i, vl.ID)
		}
	}
	return nil
}
