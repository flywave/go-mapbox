package mvt

import (
	"github.com/flywave/go-pbf"
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
	Buf             *pbf.Reader
	Proto           Proto
}

func (tile *Tile) NewLayer(endpos int, pt ProtoType) {
	proto := getProto(pt)

	layer := &Layer{StartPos: tile.Buf.Pos, EndPos: endpos, Proto: proto}
	key, val := tile.Buf.ReadTag()
	for tile.Buf.Pos < layer.EndPos {
		if key == proto.Layer.Name && val == pbf.Bytes {
			layer.Name = tile.Buf.ReadString()
			tile.Layers = append(tile.Layers, layer.Name)
			key, val = tile.Buf.ReadTag()
		}
		for key == proto.Layer.Features && val == pbf.Bytes {
			layer.features = append(layer.features, tile.Buf.Pos)
			feat_size := tile.Buf.ReadVarint()

			tile.Buf.Pos += feat_size
			key, val = tile.Buf.ReadTag()
		}
		for key == proto.Layer.Keys && val == pbf.Bytes {
			layer.Keys = append(layer.Keys, tile.Buf.ReadString())
			key, val = tile.Buf.ReadTag()
		}
		for key == proto.Layer.Values && val == pbf.Bytes {
			tile.Buf.ReadVarint()
			newkey, _ := tile.Buf.ReadTag()
			switch newkey {
			case proto.Value.StringValue:
				layer.Values = append(layer.Values, tile.Buf.ReadString())
			case proto.Value.FloatValue:
				layer.Values = append(layer.Values, tile.Buf.ReadFloat())
			case proto.Value.DoubleValue:
				layer.Values = append(layer.Values, tile.Buf.ReadDouble())
			case proto.Value.IntValue:
				layer.Values = append(layer.Values, tile.Buf.ReadInt64())
			case proto.Value.UIntValue:
				layer.Values = append(layer.Values, tile.Buf.ReadUInt64())
			case proto.Value.SIntValue:
				layer.Values = append(layer.Values, tile.Buf.ReadUInt64())
			case proto.Value.BoolIntValue:
				layer.Values = append(layer.Values, tile.Buf.ReadBool())
			}
			key, val = tile.Buf.ReadTag()
		}
		if key == proto.Layer.Extent && val == pbf.Varint {
			layer.Extent = int(tile.Buf.ReadVarint())
			key, val = tile.Buf.ReadTag()
		}
		if key == proto.Layer.Version && val == pbf.Varint {
			layer.Version = int(tile.Buf.ReadVarint())
			key, val = tile.Buf.ReadTag()
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
