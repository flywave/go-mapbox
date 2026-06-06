package style

// Source is a generic source with all possible fields across source types.
// The Type field determines which source type this is.
type Source struct {
	Type        string                   `json:"type"`
	Attribution string                   `json:"attribution,omitempty"`
	Bounds      []float64                `json:"bounds,omitempty"`
	Buffer      *float64                 `json:"buffer,omitempty"`
	Cluster     *bool                    `json:"cluster,omitempty"`
	ClusterMaxZoom    *float64           `json:"clusterMaxZoom,omitempty"`
	ClusterMinPoints  *float64           `json:"clusterMinPoints,omitempty"`
	ClusterProperties map[string]interface{} `json:"clusterProperties,omitempty"`
	ClusterRadius     *float64           `json:"clusterRadius,omitempty"`
	Coordinates [][]float64              `json:"coordinates,omitempty"`
	Data        interface{}              `json:"data,omitempty"`
	Dynamic     *bool                    `json:"dynamic,omitempty"`
	Encoding    string                   `json:"encoding,omitempty"`
	ExtraBounds [][]float64              `json:"extra_bounds,omitempty"`
	Filter      interface{}              `json:"filter,omitempty"`
	GenerateID  *bool                    `json:"generateId,omitempty"`
	LineMetrics *bool                    `json:"lineMetrics,omitempty"`
	MaxZoom     *float64                 `json:"maxzoom,omitempty"`
	MinZoom     *float64                 `json:"minzoom,omitempty"`
	Models      map[string]ModelSourceModel `json:"models,omitempty"`
	PromoteID   interface{}              `json:"promoteId,omitempty"`
	RasterLayers interface{}             `json:"rasterLayers,omitempty"`
	Scheme      string                   `json:"scheme,omitempty"`
	Tiles       []string                 `json:"tiles,omitempty"`
	TileSize    *float64                 `json:"tileSize,omitempty"`
	Tolerance   *float64                 `json:"tolerance,omitempty"`
	URL         string                   `json:"url,omitempty"`
	URLs        []string                 `json:"urls,omitempty"`
	Volatile    *bool                    `json:"volatile,omitempty"`
}

type Sources map[string]*Source

// ModelSourceModel defines properties of a single 3D model in a model source.
type ModelSourceModel struct {
	URI                   string                         `json:"uri"`
	Orientation           []float64                      `json:"orientation,omitempty"`
	Position              []float64                      `json:"position,omitempty"`
	FeatureProperties     map[string]interface{}         `json:"featureProperties,omitempty"`
	MaterialOverrideNames []string                       `json:"materialOverrideNames,omitempty"`
	MaterialOverrides     map[string]ModelMaterialOverride `json:"materialOverrides,omitempty"`
	NodeOverrideNames     []string                       `json:"nodeOverrideNames,omitempty"`
	NodeOverrides         map[string]ModelNodeOverride    `json:"nodeOverrides,omitempty"`
}

// ModelMaterialOverride defines overrides for a single material in a 3D model.
type ModelMaterialOverride struct {
	ModelColor             *ColorType `json:"model-color,omitempty"`
	ModelColorMixIntensity *float64   `json:"model-color-mix-intensity,omitempty"`
	ModelEmissiveStrength  *float64   `json:"model-emissive-strength,omitempty"`
	ModelOpacity           *float64   `json:"model-opacity,omitempty"`
}

// ModelNodeOverride defines transform overrides for a single node in a 3D model.
type ModelNodeOverride struct {
	Orientation []float64 `json:"orientation,omitempty"`
}
