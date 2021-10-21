package mbtiles

import (
	"encoding/json"
	"strconv"
	"strings"
)

type Metadata struct {
	Name        string     `json:"name"`
	Format      TileFormat `json:"format"`
	Bounds      [4]float64 `json:"bounds,omitempty"`
	Center      [3]float64 `json:"center,omitempty"`
	MinZoom     int        `json:"minzoom,omitempty"`
	MaxZoom     int        `json:"maxzoom,omitempty"`
	Description string     `json:"description,omitempty"`
	Version     string     `json:"version,omitempty"`
	Type        LayerType  `json:"type,omitempty"`
	Attribution string     `json:"attribution,omitempty"`
	LayerData   *LayerData `json:"layerData,omitempty"`
}

func (m *Metadata) ToMap() map[string]string {
	ret := make(map[string]string)
	if m.Name != "" {
		ret["name"] = m.Name
	}
	if m.Description != "" {
		ret["description"] = m.Description
	}
	ret["minzoom"] = strconv.Itoa(m.MinZoom)
	ret["maxzoom"] = strconv.Itoa(m.MaxZoom)

	ret["center"], _ = centerToString(m.Center)
	ret["bounds"], _ = boundsToString(m.Bounds)

	ret["type"] = layerTypeToString(m.Type)
	ret["format"] = tileFormatToString(m.Format)

	if m.LayerData != nil {
		data, _ := json.Marshal(m.LayerData)
		ret["layerData"] = string(data)
	}
	return ret
}

type LayerData struct {
	VectorLayers *[]VectorLayer `json:"vector_layers,omitempty"`
	TileStats    *TileStats     `json:"tilestats,omitempty"`
}

type VectorLayer struct {
	ID          string                 `json:"id"`
	Fields      map[string]interface{} `json:"fields"`
	Description string                 `json:"description,omitempty"`
	MinZoom     int                    `json:"minzoom,omitempty"`
	MaxZoom     int                    `json:"maxzoom,omitempty"`
}

type TileStats struct {
	LayerCount int     `json:"layerCount"`
	Layers     []Layer `json:"layers"`
}

type Layer struct {
	Name           string      `json:"layer"`
	Count          int64       `json:"count"`
	Geometry       string      `json:"geometry"`
	AttributeCount int         `json:"attributeCount"`
	Attributes     []Attribute `json:"attributes,omitempty"`
}

type Attribute struct {
	Name   string        `json:"attribute"`
	Count  int           `json:"count"`
	Type   string        `json:"type"`
	Values []interface{} `json:"values"`
}

func stringToTileFormat(s string) TileFormat {
	for i, k := range formatStrings {
		if k == s {
			return TileFormat(i)
		}
	}

	return UNKNOWN
}

func tileFormatToString(s TileFormat) string {
	return formatStrings[int(s)]
}

func stringToBounds(str string) (bounds [4]float64, err error) {
	for i, v := range strings.Split(str, ",") {
		bounds[i], err = strconv.ParseFloat(strings.TrimSpace(v), 64)
	}

	return
}

func boundsToString(bounds [4]float64) (str string, err error) {
	elem := []string{}
	for i := range bounds {
		s := strconv.FormatFloat(bounds[i], 'f', 8, 64)
		elem = append(elem, s)
	}
	return strings.Join(elem, ","), nil
}

func stringToCenter(str string) (center [3]float64, err error) {
	for i, v := range strings.Split(str, ",") {
		center[i], err = strconv.ParseFloat(strings.TrimSpace(v), 64)
	}

	return
}

func centerToString(center [3]float64) (str string, err error) {
	elem := []string{}
	for i := range center {
		s := strconv.FormatFloat(center[i], 'f', 8, 64)
		elem = append(elem, s)
	}
	return strings.Join(elem, ","), nil
}
