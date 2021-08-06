package mvt

import (
	"fmt"
	"strings"
	"testing"

	"github.com/flywave/go-geom"
)

func TestWriteTile(t *testing.T) {
	feats1, _ := ReadTile(bytevals2, tileid2, PROTO_LK)
	for _, feat := range feats1 {
		for k, v := range feat.Properties {
			if s, ok := v.(string); ok {
				if strings.Contains(s, "tan zhe") {
					fmt.Println(k, v)
				}
			}
			fmt.Println(k, v)
		}
	}

	conf := NewConfig("LK", tileid2, PROTO_MAPBOX)

	data := WriteLayer(feats1, conf)

	feats2, _ := ReadTile(data, tileid2, PROTO_MAPBOX)
	for _, feat := range feats2 {
		for k, v := range feat.Properties {
			fmt.Println(k, v)
		}
	}

	if len(feats1) != len(feats2) {
		t.FailNow()
	}
}

func TestWriteLayer(t *testing.T) {
	data := []byte{}
	tile, _ := NewTile(bytevals2, PROTO_LK)
	for _, layer := range tile.LayerMap {
		conf := NewConfig(layer.Name, tileid2, PROTO_MAPBOX)
		feats := []*geom.Feature{}
		for layer.Next() {
			feat, _ := layer.Feature()
			featg, _ := feat.ToGeoJSON(tileid)
			feats = append(feats, featg)
		}

		data = append(data, WriteLayer(feats, conf)...)
	}

	feats2, _ := ReadTile(data, tileid2, PROTO_MAPBOX)
	for _, feat := range feats2 {
		for k, v := range feat.Properties {
			if s, ok := v.(string); ok {
				if strings.Contains(s, "tan") {
					fmt.Println(k, v)
				}
			}
			fmt.Println(k, v)
		}
	}
}
