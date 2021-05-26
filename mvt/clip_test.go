package mvt

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/flywave/go-geom"

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
		fmt.Println(err)
	}
	feature, err := geom.UnmarshalFeature(bs)
	if err != nil {
		fmt.Println(err)
	}

	tileid := m.TileID{1, 8, 5}
	about_tile_feature := ClipTile(feature, tileid)
	about_tile_feature.Properties = map[string]interface{}{"COLORKEY": "purple", "TILEID": m.Tilestr(tileid)}
	fmt.Printf("About Tile: %+v Feature: %+v\n", tileid, about_tile_feature)

	keep_parents := false
	tilemap := ClipFeature(feature, int(tileid.Z), keep_parents)

	feats := []*geom.Feature{}
	for k, v := range tilemap {
		v.Properties = map[string]interface{}{"TILEID": m.Tilestr(k), "COLORKEY": "white"}
		feats = append(feats, v)
		fmt.Printf("Tile: %+v Feature: %+v\n", k, v)
	}

	feats = append(feats, about_tile_feature)
	MakeFeatures(feats, "../../tests/a.geojson")

	if b, _ := exists("../../tests/a.geojson"); b {
		os.Remove("../../tests/a.geojson")
	}
}
