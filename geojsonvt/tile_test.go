package geojsonvt

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/flywave/go-geom"
	"github.com/flywave/go-geom/general"
	m "github.com/flywave/go-mapbox/tileid"
)

func TestTile(t *testing.T) {
	bs, err := ioutil.ReadFile("../../tests/alaska.geojson")
	if err != nil {
		fmt.Println(err)
	}
	features, err := general.UnmarshalFeature(bs)
	if err != nil {
		fmt.Println(err)
	}

	options := NewConfig()

	tile := TileFromGeoJSON([]*geom.Feature{features}, m.TileID{X: 1, Y: 8, Z: 5}, options)

	bytevals := tile.Marshal()
	fmt.Println(len(bytevals))
}
