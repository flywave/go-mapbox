package mbtiles

type Metadata struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version,omitempty"`
	Minzoom     string `json:"minzoom,omitempty"`
	Maxzoom     string `json:"maxzoom,omitempty"`
	Center      string `json:"center,omitempty"`
	Bounds      string `json:"bounds,omitempty"`
	Type        string `json:"type,omitempty"`
	Format      string `json:"format,omitempty"`
	Json        string `json:"json,omitempty"`
}

type MetadataJson struct {
	VectorLayers []VectorLayer `json:"vector_layers,omitempty"`
	TileStats    TileStats     `json:"tilestats,omitempty"`
}

type VectorLayer struct {
	Id          string            `json:"id,omitempty"`
	Description string            `json:"description,omitempty"`
	Minzoom     int               `json:"minzoom,omitempty"`
	Maxzoom     int               `json:"maxzoom,omitempty"`
	Fields      map[string]string `json:"fields,omitempty"`
}

type TileStats struct {
	LayerCount int `json:"layerCount,omitempty"`
	Layers     []struct {
		Layer          string `json:"layer,omitempty"`
		Count          int    `json:"count,omitempty"`
		Geometry       string `json:"geometry,omitempty"`
		AttributeCount int    `json:"attributeCount,omitempty"`
		Attributes     []struct {
			Attribute string        `json:"attribute,omitempty"`
			Count     int           `json:"count,omitempty"`
			Type      string        `json:"type,omitempty"`
			Values    []interface{} `json:"values,omitempty"`
			Min       interface{}   `json:"min,omitempty"`
			Max       interface{}   `json:"max,omitempty"`
		} `json:"attributes,omitempty"`
	} `json:"layers,omitempty"`
}
