package mvt

import (
	"fmt"
	"io/ioutil"

	"github.com/flywave/go-geom"
	m "github.com/flywave/go-mapbox/tileid"

	"testing"
)

var bytevals, _ = ioutil.ReadFile("../../tests/701_1635_12.pbf")
var tileid = m.TileID{X: 701, Y: 1635, Z: 12}

func TestReads(t *testing.T) {
	feats1, _ := ReadTile(bytevals, tileid)
	m1, m2 := map[interface{}]*geom.Feature{}, map[interface{}]*geom.Feature{}
	for _, feat := range feats1 {
		delete(feat.Properties, "layer")
		m1[feat.Properties["@id"]] = feat
	}
	tile, _ := NewTile(bytevals)
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
	feats1, _ := ReadTile(bytevals, tileid)
	con := NewConfig("new", tileid)
	bs, _ := WriteLayer(feats1, con)
	feats1, _ = ReadTile(bs, tileid)

	m1, m2 := map[interface{}]*geom.Feature{}, map[interface{}]*geom.Feature{}
	for _, feat := range feats1 {
		delete(feat.Properties, "layer")
		m1[feat.Properties["@id"]] = feat
	}
	tile, _ := NewTile(bs)
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
