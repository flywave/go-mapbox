package mvt

import (
	"fmt"
	"io/ioutil"

	"github.com/flywave/go-geom"
	m "github.com/flywave/go-mapbox/tileid"

	"testing"
)

var bytevals, _ = ioutil.ReadFile("../data/703_1635_12.pbf")
var tileid = m.TileID{X: 701, Y: 1635, Z: 12}

func TestReads(t *testing.T) {
	feats1, _ := ReadTile(bytevals, tileid, PROTO_MAPBOX)
	m1, m2 := map[interface{}]*geom.Feature{}, map[interface{}]*geom.Feature{}
	for _, feat := range feats1 {
		delete(feat.Properties, "layer")
		m1[feat.Properties["@id"]] = feat
	}
	tile, _ := NewTile(bytevals, PROTO_MAPBOX)
	for _, layer := range tile.LayerMap {
		for layer.Next() {
			feat, _ := layer.Feature()
			featg, _ := feat.ToGeoJSON(tileid)
			delete(featg.Properties, "layer")

			m2[featg.Properties["@id"]] = featg
		}
	}
	if len(m2) != len(m2) {
		t.Errorf("Map sizes are different.")
	}
	i := 0
	for k := range m1 {
		i++
		v1, b1 := m1[k]
		v2, b2 := m2[k]
		if b1 && b2 {
			err := geom.IsFeatureEqual(*v1, *v2)
			if !err {
				t.Errorf("freeature not eq")
			}
		} else {
			t.Errorf("Both geojson features weren't in map.")
		}

	}

	fmt.Printf("Lazy Reads Are Exactly the same as bulk reads for %d features in tile.\n", len(feats1))
}

func TestReadsWrites(t *testing.T) {
	feats1, _ := ReadTile(bytevals, tileid, PROTO_MAPBOX)
	con := NewConfig("new", tileid, PROTO_MAPBOX)
	bs := WriteLayer(feats1, con)
	feats1, _ = ReadTile(bs, tileid, PROTO_MAPBOX)

	m1, m2 := map[interface{}]*geom.Feature{}, map[interface{}]*geom.Feature{}
	for _, feat := range feats1 {
		delete(feat.Properties, "layer")
		m1[feat.Properties["@id"]] = feat
	}
	tile, _ := NewTile(bs, PROTO_MAPBOX)
	for _, layer := range tile.LayerMap {
		for layer.Next() {
			feat, _ := layer.Feature()
			featg, _ := feat.ToGeoJSON(tileid)
			delete(featg.Properties, "layer")

			m2[featg.Properties["@id"]] = featg
		}
	}
	if len(m2) != len(m2) {
		t.Errorf("Map sizes are different.")
	}
	i := 0
	for k := range m1 {
		i++
		v1, b1 := m1[k]
		v2, b2 := m2[k]
		if b1 && b2 {
			err := geom.IsFeatureEqual(*v1, *v2)
			if !err {
				t.Errorf("freeature not eq")
			}
		} else {
			t.Errorf("Both geojson features weren't in map.")
		}

	}

	fmt.Printf("Lazy Reads Are Exactly the same as bulk reads when written and read again features for %d features in tile.\n", len(feats1))
}

func TestM(t *testing.T) {
	bytevals := []byte{0x1a, 0xc3, 0x2, 0xa, 0x4, 0x54, 0x65, 0x73, 0x74, 0x12, 0x3c, 0x12, 0x1a, 0x0, 0x0, 0x1, 0x1, 0x2, 0x2, 0x3, 0x3, 0x4, 0x4, 0x5, 0x5, 0x6, 0x6, 0x7, 0x7, 0x8, 0x8, 0x9, 0x9, 0xa, 0xa, 0xb, 0xa, 0xc, 0xb, 0x18, 0x2, 0x22, 0x1c, 0x9, 0x80, 0x41, 0xde, 0x3, 0x42, 0x75, 0x8d, 0x1, 0xab, 0x1, 0x71, 0x5d, 0x5b, 0xa9, 0x1, 0x83, 0x1, 0x8f, 0x1, 0x57, 0xdb, 0x2, 0x69, 0x43, 0x21, 0x1d, 0x19, 0x1a, 0x7, 0x52, 0x4f, 0x55, 0x54, 0x45, 0x49, 0x44, 0x1a, 0x8, 0x53, 0x75, 0x62, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x1a, 0xa, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x79, 0x43, 0x6f, 0x64, 0x65, 0x1a, 0x8, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x63, 0x74, 0x1a, 0x6, 0x4f, 0x4e, 0x45, 0x57, 0x41, 0x59, 0x1a, 0x5, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x1a, 0x5, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x1a, 0xa, 0x53, 0x69, 0x67, 0x6e, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x1a, 0x3, 0x45, 0x4d, 0x50, 0x1a, 0xa, 0x53, 0x68, 0x61, 0x70, 0x65, 0x5f, 0x4c, 0x65, 0x6e, 0x67, 0x1a, 0x3, 0x42, 0x4d, 0x50, 0x1a, 0x8, 0x53, 0x75, 0x70, 0x70, 0x43, 0x6f, 0x64, 0x65, 0x1a, 0xa, 0x53, 0x74, 0x72, 0x65, 0x65, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0xf, 0xa, 0xd, 0x35, 0x30, 0x34, 0x30, 0x30, 0x35, 0x32, 0x36, 0x35, 0x30, 0x30, 0x30, 0x30, 0x22, 0x9, 0x19, 0x0, 0x0, 0x0, 0x0, 0x0, 0x40, 0x50, 0x40, 0x22, 0x9, 0x19, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x49, 0x40, 0x22, 0x9, 0x19, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x40, 0x22, 0x2, 0x38, 0x0, 0x22, 0x7, 0xa, 0x5, 0x35, 0x32, 0x2f, 0x36, 0x35, 0x22, 0x9, 0x19, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4a, 0x40, 0x22, 0x3, 0xa, 0x1, 0x34, 0x22, 0x9, 0x19, 0x40, 0x96, 0x4f, 0xa0, 0x99, 0x99, 0xd, 0x40, 0x22, 0x9, 0x19, 0x91, 0x80, 0xf2, 0xf3, 0x68, 0xbb, 0xb6, 0x40, 0x22, 0x9, 0x19, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x22, 0x14, 0xa, 0x12, 0x42, 0x49, 0x47, 0x20, 0x53, 0x41, 0x4e, 0x44, 0x59, 0x20, 0x52, 0x49, 0x56, 0x45, 0x52, 0x20, 0x52, 0x44, 0x78, 0x2}
	xyz := m.TileID{X: 1107, Y: 1578, Z: 12}

	tile, err := NewTile(bytevals, PROTO_MAPBOX)
	if err != nil {
		fmt.Println(err)
	}
	for layername, layer := range tile.LayerMap {
		fmt.Printf("LayerName: %s\n", layername)
		for layer.Next() {
			lazyfeature, err := layer.Feature()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("Lazy Feature: %+v\n", lazyfeature)

			geom, err := lazyfeature.LoadGeometry()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("Geometry: %+v\n", geom)

			feat, err := lazyfeature.ToGeoJSON(xyz)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("Geojson Feature: %+v\n", feat)

		}
	}
}
