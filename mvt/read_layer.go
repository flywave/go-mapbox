package mvt

import (
	"github.com/murphy214/pbf"
)

type Layer struct {
	Name            string
	Extent          int
	Version         int
	Keys            []string
	Values          []interface{}
	Number_Features int
	features        []int
	StartPos        int
	EndPos          int
	featurePosition int
	Buf             *pbf.PBF
}

func (tile *Tile) NewLayer(endpos int) {
	layer := &Layer{StartPos: tile.Buf.Pos, EndPos: endpos}
	key, val := tile.Buf.ReadKey()
	for tile.Buf.Pos < layer.EndPos {
		if key == 1 && val == 2 {
			layer.Name = tile.Buf.ReadString()
			tile.Layers = append(tile.Layers, layer.Name)
			key, val = tile.Buf.ReadKey()
		}
		for key == 2 && val == 2 {
			layer.features = append(layer.features, tile.Buf.Pos)
			feat_size := tile.Buf.ReadVarint()

			tile.Buf.Pos += feat_size
			key, val = tile.Buf.ReadKey()
		}
		for key == 3 && val == 2 {
			layer.Keys = append(layer.Keys, tile.Buf.ReadString())
			key, val = tile.Buf.ReadKey()
		}
		for key == 4 && val == 2 {
			tile.Buf.ReadVarint()
			newkey, _ := tile.Buf.ReadKey()
			switch newkey {
			case 1:
				layer.Values = append(layer.Values, tile.Buf.ReadString())
			case 2:
				layer.Values = append(layer.Values, tile.Buf.ReadFloat())
			case 3:
				layer.Values = append(layer.Values, tile.Buf.ReadDouble())
			case 4:
				layer.Values = append(layer.Values, tile.Buf.ReadInt64())
			case 5:
				layer.Values = append(layer.Values, tile.Buf.ReadUInt64())
			case 6:
				layer.Values = append(layer.Values, tile.Buf.ReadUInt64())
			case 7:
				layer.Values = append(layer.Values, tile.Buf.ReadBool())
			}
			key, val = tile.Buf.ReadKey()
		}
		if key == 5 && val == 0 {
			layer.Extent = int(tile.Buf.ReadVarint())
			key, val = tile.Buf.ReadKey()
		}
		if key == 15 && val == 0 {
			layer.Version = int(tile.Buf.ReadVarint())
			key, val = tile.Buf.ReadKey()

		}
	}

	if layer.Extent == 0 {
		layer.Extent = 4096
	}
	layer.Number_Features = len(layer.features)
	tile.LayerMap[layer.Name] = layer
	tile.Buf.Pos = endpos
	layer.Buf = tile.Buf
}

func (layer *Layer) Next() bool {
	return layer.featurePosition < layer.Number_Features
}

func (layer *Layer) Reset() {
	layer.featurePosition = 0
}
