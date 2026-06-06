package recipe

import (
	"fmt"
	"strings"
)

const CurrentVersion = 1

// RecipeType identifies the type of recipe.
type RecipeType string

const (
	RecipeVector      RecipeType = ""
	RecipeRaster      RecipeType = "raster"
	RecipeRasterArray RecipeType = "rasterarray"
)

// Recipe is the top-level structure for any Mapbox Tiling Service recipe.
type Recipe struct {
	Version     int               `json:"version"`
	Type        RecipeType        `json:"type,omitempty"`
	FillZoom    *int              `json:"fillzoom,omitempty"`
	Incremental *bool             `json:"incremental,omitempty"`
	Sources     []Source          `json:"sources,omitempty"`
	MinZoom     *int              `json:"minzoom,omitempty"`
	MaxZoom     *int              `json:"maxzoom,omitempty"`
	Layers      map[string]*Layer `json:"layers"`
}

// Source defines an input source for raster/rasterarray recipes.
type Source struct {
	URI string `json:"uri"`
	CRS string `json:"crs,omitempty"`
}

// Layer defines a single layer in a recipe.
type Layer struct {
	Source      string          `json:"source,omitempty"`
	MinZoom     *uint           `json:"minzoom,omitempty"`
	MaxZoom     *uint           `json:"maxzoom,omitempty"`
	Features    *FeaturesConfig `json:"features,omitempty"`
	Tiles       *TilesConfig    `json:"tiles,omitempty"`
	TileSize    *int            `json:"tilesize,omitempty"`
	Buffer      *int            `json:"buffer,omitempty"`
	Units       string          `json:"units,omitempty"`
	InputNoData *float64        `json:"input_no_data_value,omitempty"`
	Resampling  string          `json:"resampling,omitempty"`
	Offset      *float64        `json:"offset,omitempty"`
	Scale       *float64        `json:"scale,omitempty"`
	SourceRules *SourceRules    `json:"source_rules,omitempty"`
}

// FeaturesConfig describes per-feature processing options (vector recipe).
type FeaturesConfig struct {
	ID             interface{}         `json:"id,omitempty"`
	BBox           []float64           `json:"bbox,omitempty"`
	Attributes     *FeaturesAttributes `json:"attributes,omitempty"`
	Filter         interface{}         `json:"filter,omitempty"`
	Simplification interface{}         `json:"simplification,omitempty"`
}

// FeaturesAttributes describes feature attribute operations.
type FeaturesAttributes struct {
	ZoomElement   []string               `json:"zoom_element,omitempty"`
	Set           map[string]interface{} `json:"set,omitempty"`
	AllowedOutput []string               `json:"allowed_output,omitempty"`
}

// Simplification controls geometry simplification.
type Simplification struct {
	Distance    interface{} `json:"distance,omitempty"`
	OutwardOnly interface{} `json:"outward_only,omitempty"`
}

// TilesConfig describes per-tile processing options (vector recipe).
type TilesConfig struct {
	BBox         []float64        `json:"bbox,omitempty"`
	Extent       interface{}      `json:"extent,omitempty"`
	BufferSize   interface{}      `json:"buffer_size,omitempty"`
	Limit        []LimitRule      `json:"limit,omitempty"`
	Union        []UnionObject    `json:"union,omitempty"`
	Filter       interface{}      `json:"filter,omitempty"`
	Attributes   *TilesAttributes `json:"attributes,omitempty"`
	Order        string           `json:"order,omitempty"`
	RemoveFilled interface{}      `json:"remove_filled,omitempty"`
	ID           interface{}      `json:"id,omitempty"`
	LayerSize    *int             `json:"layer_size,omitempty"`
}

// TilesAttributes describes tile-level attribute operations.
type TilesAttributes struct {
	Set map[string]interface{} `json:"set,omitempty"`
}

// LimitRule limits the number of features in a tile.
type LimitRule []interface{}

// UnionObject joins features with matching attributes.
type UnionObject struct {
	Where             interface{}       `json:"where,omitempty"`
	GroupBy           []string          `json:"group_by,omitempty"`
	Cluster           *bool             `json:"cluster,omitempty"`
	RegionCount       *int              `json:"region_count,omitempty"`
	Aggregate         map[string]string `json:"aggregate,omitempty"`
	MaintainDirection *bool             `json:"maintain_direction,omitempty"`
	Simplification    interface{}       `json:"simplification,omitempty"`
}

// SourceRules defines rules for mapping sources to bands (raster/rasterarray).
type SourceRules struct {
	Filter  interface{} `json:"filter,omitempty"`
	SortKey interface{} `json:"sort_key,omitempty"`
	Name    interface{} `json:"name,omitempty"`
	Order   string      `json:"order,omitempty"`
}

// ─── Validation ─────────────────────────────────────────────────────────────

func (r *Recipe) Validate() error {
	if r.Version != CurrentVersion {
		return fmt.Errorf("unsupported recipe version %d (expected %d)", r.Version, CurrentVersion)
	}
	if len(r.Layers) == 0 {
		return fmt.Errorf("recipe must have at least one layer")
	}
	if len(r.Layers) > 20 {
		return fmt.Errorf("recipe cannot have more than 20 layers (got %d)", len(r.Layers))
	}

	switch r.Type {
	case RecipeVector:
		return r.validateVector()
	case RecipeRaster:
		return r.validateRaster()
	case RecipeRasterArray:
		return r.validateRasterArray()
	default:
		return fmt.Errorf("unknown recipe type %q", r.Type)
	}
}

func (r *Recipe) validateVector() error {
	for name, layer := range r.Layers {
		if layer == nil {
			return fmt.Errorf("layer %q is nil", name)
		}
		if layer.Source == "" {
			return fmt.Errorf("layer %q: source is required", name)
		}
		if layer.MinZoom == nil {
			return fmt.Errorf("layer %q: minzoom is required", name)
		}
		if layer.MaxZoom == nil {
			return fmt.Errorf("layer %q: maxzoom is required", name)
		}
		if *layer.MinZoom > *layer.MaxZoom {
			return fmt.Errorf("layer %q: minzoom (%d) > maxzoom (%d)", name, *layer.MinZoom, *layer.MaxZoom)
		}
		if *layer.MaxZoom > 16 {
			return fmt.Errorf("layer %q: maxzoom (%d) > 16", name, *layer.MaxZoom)
		}
	}
	return nil
}

func (r *Recipe) validateRaster() error {
	if len(r.Sources) == 0 {
		return fmt.Errorf("raster recipe requires at least one source")
	}
	if r.MinZoom == nil {
		return fmt.Errorf("raster recipe requires minzoom")
	}
	if r.MaxZoom == nil {
		return fmt.Errorf("raster recipe requires maxzoom")
	}
	if *r.MinZoom > *r.MaxZoom {
		return fmt.Errorf("minzoom (%d) > maxzoom (%d)", *r.MinZoom, *r.MaxZoom)
	}
	return nil
}

func (r *Recipe) validateRasterArray() error {
	if len(r.Sources) == 0 {
		return fmt.Errorf("rasterarray recipe requires at least one source")
	}
	if r.MinZoom != nil && r.MaxZoom != nil && *r.MinZoom > *r.MaxZoom {
		return fmt.Errorf("minzoom (%d) > maxzoom (%d)", *r.MinZoom, *r.MaxZoom)
	}
	return nil
}

// ─── Helpers ────────────────────────────────────────────────────────────────

// IsVector returns true if this is a vector recipe.
func (r *Recipe) IsVector() bool { return r.Type == RecipeVector || r.Type == "" }

// IsRaster returns true if this is a raster recipe.
func (r *Recipe) IsRaster() bool { return r.Type == RecipeRaster }

// IsRasterArray returns true if this is a rasterarray recipe.
func (r *Recipe) IsRasterArray() bool { return r.Type == RecipeRasterArray }

// LayerNames returns the names of all layers in the recipe.
func (r *Recipe) LayerNames() []string {
	names := make([]string, 0, len(r.Layers))
	for name := range r.Layers {
		names = append(names, name)
	}
	return names
}

// AddLayer adds a layer to the recipe. Returns error if layer name is invalid.
func (r *Recipe) AddLayer(name string, layer *Layer) error {
	if !validLayerName(name) {
		return fmt.Errorf("invalid layer name %q: only underscores and alphanumeric characters allowed", name)
	}
	if r.Layers == nil {
		r.Layers = make(map[string]*Layer)
	}
	if _, ok := r.Layers[name]; ok {
		return fmt.Errorf("layer %q already exists", name)
	}
	r.Layers[name] = layer
	return nil
}

// RemoveLayer removes a layer by name.
func (r *Recipe) RemoveLayer(name string) {
	delete(r.Layers, name)
}

// NewVectorRecipe creates a vector recipe with the given layers.
func NewVectorRecipe(layers map[string]*Layer) *Recipe {
	return &Recipe{
		Version: CurrentVersion,
		Layers:  layers,
	}
}

// NewRasterRecipe creates a raster recipe with the given sources, zoom range, and layers.
func NewRasterRecipe(sources []Source, minzoom, maxzoom int, layers map[string]*Layer) *Recipe {
	mz := minzoom
	xz := maxzoom
	return &Recipe{
		Version: CurrentVersion,
		Type:    RecipeRaster,
		Sources: sources,
		MinZoom: &mz,
		MaxZoom: &xz,
		Layers:  layers,
	}
}

// NewRasterArrayRecipe creates a rasterarray recipe.
func NewRasterArrayRecipe(sources []Source, layers map[string]*Layer) *Recipe {
	return &Recipe{
		Version: CurrentVersion,
		Type:    RecipeRasterArray,
		Sources: sources,
		Layers:  layers,
	}
}

// ─── Helpers for Layer creation ─────────────────────────────────────────────

// NewLayer creates a basic vector layer with the required fields.
func NewLayer(source string, minzoom, maxzoom uint) *Layer {
	mz := minzoom
	xz := maxzoom
	return &Layer{
		Source:  source,
		MinZoom: &mz,
		MaxZoom: &xz,
	}
}

// NewRasterLayer creates a basic raster/rasterarray layer.
func NewRasterLayer(opts ...RasterLayerOption) *Layer {
	l := &Layer{}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// RasterLayerOption is a function that configures a raster layer.
type RasterLayerOption func(*Layer)

// WithTileSize sets the tilesize for a raster layer.
func WithTileSize(size int) RasterLayerOption {
	return func(l *Layer) { l.TileSize = &size }
}

// WithBuffer sets the buffer for a raster layer.
func WithBuffer(buf int) RasterLayerOption {
	return func(l *Layer) { l.Buffer = &buf }
}

// WithUnits sets the units for a raster layer.
func WithUnits(units string) RasterLayerOption {
	return func(l *Layer) { l.Units = units }
}

// WithOffset sets the offset for a raster layer.
func WithOffset(offset float64) RasterLayerOption {
	return func(l *Layer) { l.Offset = &offset }
}

// WithScale sets the scale for a raster layer.
func WithScale(scale float64) RasterLayerOption {
	return func(l *Layer) { l.Scale = &scale }
}

// WithResampling sets the resampling method for a raster layer.
func WithResampling(method string) RasterLayerOption {
	return func(l *Layer) { l.Resampling = method }
}

// WithInputNoData sets the input_no_data_value for a raster layer.
func WithInputNoData(value float64) RasterLayerOption {
	return func(l *Layer) { l.InputNoData = &value }
}

// WithSourceRules sets the source_rules for a raster layer.
func WithSourceRules(rules *SourceRules) RasterLayerOption {
	return func(l *Layer) { l.SourceRules = rules }
}

// ─── Internal helpers ───────────────────────────────────────────────────────

func validLayerName(name string) bool {
	if name == "" {
		return false
	}
	for _, r := range name {
		if !(r == '_' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
			return false
		}
	}
	return true
}

// String returns a human-readable summary of the recipe.
func (r *Recipe) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Recipe v%d", r.Version))
	switch r.Type {
	case RecipeRaster:
		b.WriteString(" [raster]")
	case RecipeRasterArray:
		b.WriteString(" [rasterarray]")
	}
	b.WriteString(fmt.Sprintf(" (%d layers)", len(r.Layers)))
	return b.String()
}
