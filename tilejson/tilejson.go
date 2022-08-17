package tilejson

import "time"

type SortBy string

const (
	SortByCreated  SortBy = "created"
	SortByModified SortBy = "modified"
)

const TileJSONVersion = "2.2.0"

type TilesetType string

const (
	VectorTileset TilesetType = "vector"
	RasterTileset TilesetType = "raster"
)

type TilesetVisibility string

const (
	PublicTileset  TilesetVisibility = "public"
	PrivateTileset TilesetVisibility = "private"
)

type ListTilesetsParams struct {
	Type       *TilesetType
	Visibility *TilesetVisibility
	SortBy     *SortBy
	Limit      *int
}

type GeomType string

const (
	GeomTypePoint   GeomType = "point"
	GeomTypeLine    GeomType = "line"
	GeomTypePolygon GeomType = "polygon"
	GeomTypeUnknown GeomType = "unknown"
)

const (
	SchemeXYZ  = "xyz"
	SchemeTMLS = "tms"
)

type FieldType string

const (
	FieldTypeString  FieldType = "String"
	FieldTypeNumber  FieldType = "Number"
	FieldTypeBoolean FieldType = "Boolean"
)

func NewTileBounds(minLon, minLat, maxLon, maxLat float64) *[4]float64 {
	v := [4]float64{minLon, minLat, maxLon, maxLat}
	return &v
}

func NewTileCenter(lon, lat float64, zoom float64) *[3]float64 {
	v := [3]float64{lon, lat, zoom}
	return &v
}

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

type TileJSON struct {
	Id           string        `json:"id"`
	Attribution  *string       `json:"attribution"`
	Bounds       [4]float64    `json:"bounds"`
	Center       [3]float64    `json:"center"`
	Format       string        `json:"format"`
	MinZoom      uint          `json:"minzoom"`
	MaxZoom      uint          `json:"maxzoom"`
	Name         *string       `json:"name"`
	Description  *string       `json:"description"`
	Scheme       string        `json:"scheme"`
	TileJSON     string        `json:"tilejson"`
	Tiles        []string      `json:"tiles"`
	Grids        []string      `json:"grids"`
	Data         []string      `json:"data"`
	Version      string        `json:"version"`
	Template     *string       `json:"template"`
	Legend       *string       `json:"legend"`
	VectorLayers []VectorLayer `json:"vector_layers,omitempty"`
	Type         string        `json:"type"`
	TileSize     uint32        `json:"tileSize"`
}

type VectorLayer struct {
	Version      int                  `json:"version"`
	Extent       int                  `json:"extent"`
	ID           string               `json:"id"`
	Source       string               `json:"source"`
	Name         string               `json:"source_name"`
	Fields       map[string]FieldType `json:"fields"`
	FeatureTags  []string             `json:"feature_tags,omitempty"`
	GeometryType GeomType             `json:"geometry_type,omitempty"`
	MinZoom      uint                 `json:"minzoom"`
	MaxZoom      uint                 `json:"maxzoom"`
	Tiles        []string             `json:"tiles"`
}
