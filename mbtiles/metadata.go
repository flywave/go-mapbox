package mbtiles

import (
	"encoding/json"
	"strconv"
	"strings"
)

const (
	ORIGIN_UL = "ul"
	ORIGIN_LL = "ll"
	ORIGIN_NW = "nw"
	ORIGIN_SW = "sw"

	RES_FACTOR_SQRT2 = "sqrt2"
	RES_FACTOR_2     = float64(2.0)
)

type Metadata struct {
	Name            string      `json:"name"`
	Format          TileFormat  `json:"format"`
	Bounds          [4]float64  `json:"bounds,omitempty"`
	Center          [3]float64  `json:"center,omitempty"`
	MinZoom         int         `json:"minzoom,omitempty"`
	MaxZoom         int         `json:"maxzoom,omitempty"`
	Description     string      `json:"description,omitempty"`
	Version         string      `json:"version,omitempty"`
	Type            LayerType   `json:"type,omitempty"`
	Attribution     string      `json:"attribution,omitempty"`
	LayerData       *LayerData  `json:"json,omitempty"`
	DirectoryLayout string      `json:"directory_layout,omitempty"`
	Origin          string      `json:"origin,omitempty"`
	Srs             string      `json:"srs,omitempty"`
	BoundsSrs       string      `json:"bounds_srs,omitempty"`
	ResFactor       interface{} `json:"res_factor,omitempty"`
	TileSize        *[2]int     `json:"tile_size,omitempty"`
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
		ret["json"] = string(data)
	}

	if m.DirectoryLayout != "" {
		ret["directory_layout"] = m.DirectoryLayout
	}

	if m.Origin != "" {
		ret["origin"] = m.Origin
	}

	if m.Srs != "" {
		ret["srs"] = m.Srs
	}

	if m.BoundsSrs != "" {
		ret["bounds_srs"] = m.BoundsSrs
	}

	if m.ResFactor != nil {
		ret["res_factor"] = resFactorToString(m.ResFactor)
	}

	if m.TileSize != nil {
		ret["tile_size"] = tileSizeToString(*m.TileSize)
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

func resFactorToString(fac interface{}) string {
	switch v := fac.(type) {
	case string:
		return v
	case float64:
		s := strconv.FormatFloat(v, 'f', 8, 64)
		return s
	}
	return ""
}

func stringToResFactor(fac string) interface{} {
	if fac == "sqrt2" {
		return fac
	} else {
		f, err := strconv.ParseFloat(strings.TrimSpace(fac), 64)
		if err != nil {
			return nil
		}
		return f
	}
}

func tileSizeToString(ts [2]int) string {
	elem := []string{}
	for i := range ts {
		s := strconv.Itoa(ts[i])
		elem = append(elem, s)
	}
	return strings.Join(elem, ",")
}

func stringToTileSize(str string) (tilesize *[2]int, err error) {
	tilesize = new([2]int)
	for i, v := range strings.Split(str, ",") {
		var i64 int64
		i64, err = strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		tilesize[i] = int(i64)
	}

	return
}
