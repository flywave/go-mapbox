package mvt

import (
	"fmt"
	"testing"
)

func TestWriteLayer(t *testing.T) {
	feats1, _ := ReadTile(bytevals2, tileid2, PROTO_LK)
	for _, feat := range feats1 {
		for k, v := range feat.Properties {
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
}
