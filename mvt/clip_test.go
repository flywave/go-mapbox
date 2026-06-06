package mvt

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/flywave/go-geom"
	"github.com/flywave/go-geom/general"
	m "github.com/flywave/go-mapbox/tileid"
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func TestClip(t *testing.T) {
	bs, err := ioutil.ReadFile("../../tests/alaska.geojson")
	if err != nil {
		t.Skip("skipping: alaska.geojson not found")
	}
	feature, err := general.UnmarshalFeature(bs)
	if err != nil {
		t.Skip("skipping: could not parse alaska.geojson")
	}

	tileid := m.TileID{X: 1, Y: 8, Z: 5}
	about_tile_feature := ClipTile(feature, tileid)
	if about_tile_feature == nil {
		t.Fatal("ClipTile returned nil")
	}
	about_tile_feature.Properties = map[string]interface{}{"COLORKEY": "purple", "TILEID": m.Tilestr(tileid)}

	keep_parents := false
	tilemap := ClipFeature(feature, int(tileid.Z), keep_parents)

	feats := []*geom.Feature{}
	for k, v := range tilemap {
		v.Properties = map[string]interface{}{"TILEID": m.Tilestr(k), "COLORKEY": "white"}
		feats = append(feats, v)
	}

	feats = append(feats, about_tile_feature)
	MakeFeatures(feats, "../../tests/a.geojson")

	if b, _ := exists("../../tests/a.geojson"); b {
		os.Remove("../../tests/a.geojson")
	}
}
