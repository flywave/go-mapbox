package mvt

import (
	m "github.com/flywave/go-mapbox/tileid"

	"github.com/flywave/go-geom"

	"github.com/murphy214/pbf"
)

func RemoveLayer(bytevals []byte, layername string) ([]byte, error) {
	tile, err := NewTile(bytevals)
	if err != nil {
		return []byte{}, err
	}
	layer, boolval := tile.LayerMap[layername]
	size := layer.EndPos - layer.StartPos
	bsize := len(pbf.EncodeVarint(uint64(size)))
	if boolval {
		return append(bytevals[:layer.StartPos-bsize-1], bytevals[layer.EndPos:]...), nil
	}
	return bytevals, nil
}

func CleanLayer(tileid m.TileID, layer *Layer) LayerWrite {
	keys_map := map[string]uint32{}
	values_map := map[interface{}]uint32{}
	layerwrite := LayerWrite{
		Name:       layer.Name,
		Extent:     layer.Extent,
		Version:    layer.Version,
		Cursor:     NewCursor(tileid),
		Keys_Map:   keys_map,
		Values_Map: values_map,
	}
	for _, key := range layer.Keys {
		layerwrite.AddKey(key)
	}
	for _, value := range layer.Values {
		layerwrite.AddValue(value)
	}

	firstfeature := layer.features[0] - 1
	lastfeature := layer.features[len(layer.features)-1]
	layer.Buf.Pos = lastfeature
	feat_size := layer.Buf.ReadVarint()
	layer.Buf.Pos += feat_size

	endpos := layer.Buf.Pos
	layerwrite.Features = layer.Buf.Pbf[firstfeature:endpos]
	return layerwrite
}

func AddFeaturesToLayer(bytevals []byte, config Config, features []*geom.Feature) ([]byte, error) {
	tileid := config.TileID
	layername := config.Name

	tile, err := NewTile(bytevals)
	if err != nil {
		return []byte{}, err
	}
	layer, boolval := tile.LayerMap[layername]
	if boolval {
		layerwrite := CleanLayer(tileid, layer)

		size := layer.EndPos - layer.StartPos
		bsize := len(pbf.EncodeVarint(uint64(size)))
		if boolval {
			bytevals = append(bytevals[:layer.StartPos-bsize-1], bytevals[layer.EndPos:]...)
		}

		for _, feature := range features {
			layerwrite.AddFeature(feature)
		}
		return append(bytevals, layerwrite.Flush()...), nil

	} else {
		layerwrite := NewLayerConfig(config)
		for _, feature := range features {
			layerwrite.AddFeature(feature)
		}
		return append(bytevals, layerwrite.Flush()...), nil
	}
}
