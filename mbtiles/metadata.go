package mbtiles

import "encoding/json"

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

func NewMetadata(md map[string]string) *Metadata {
	ret := &Metadata{}
	if v, ok := md["name"]; ok {
		ret.Name = v
	}
	if v, ok := md["description"]; ok {
		ret.Description = v
	}
	if v, ok := md["version"]; ok {
		ret.Version = v
	}
	if v, ok := md["minzoom"]; ok {
		ret.Minzoom = v
	}
	if v, ok := md["maxzoom"]; ok {
		ret.Maxzoom = v
	}
	if v, ok := md["center"]; ok {
		ret.Center = v
	}
	if v, ok := md["bounds"]; ok {
		ret.Bounds = v
	}
	if v, ok := md["type"]; ok {
		ret.Type = v
	}
	if v, ok := md["format"]; ok {
		ret.Format = v
	}
	if v, ok := md["json"]; ok {
		ret.Json = v
	}
	return ret
}

func (m *Metadata) ToMap() map[string]string {
	ret := make(map[string]string)
	if m.Name != "" {
		ret["name"] = m.Name
	}
	if m.Description != "" {
		ret["description"] = m.Description
	}
	if m.Minzoom != "" {
		ret["minzoom"] = m.Minzoom
	}
	if m.Maxzoom != "" {
		ret["maxzoom"] = m.Maxzoom
	}
	if m.Center != "" {
		ret["center"] = m.Center
	}
	if m.Bounds != "" {
		ret["bounds"] = m.Bounds
	}
	if m.Type != "" {
		ret["type"] = m.Type
	}
	if m.Format != "" {
		ret["format"] = m.Format
	}
	if m.Json != "" {
		ret["json"] = m.Json
	}
	return ret
}

func (m *Metadata) SetMetadataJson(md *MetadataJson) {
	data, _ := json.Marshal(md)
	m.Json = string(data)
}

func (m *Metadata) GetMetadataJson() *MetadataJson {
	if m.Json != "" {
		var ret *MetadataJson
		json.Unmarshal([]byte(m.Json), ret)
		return ret
	}
	return nil
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
