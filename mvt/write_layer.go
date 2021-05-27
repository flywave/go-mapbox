package mvt

import (
	m "github.com/flywave/go-mapbox/tileid"

	"github.com/flywave/go-geom"

	"github.com/flywave/go-pbf"
)

type LayerWrite struct {
	TileID     m.TileID
	DeltaX     float64
	DeltaY     float64
	Name       string
	Extent     int
	Version    int
	Keys_Map   map[string]uint32
	Values_Map map[interface{}]uint32
	Cursor     *Cursor
	ReduceBool bool
	Buf        *pbf.Writer
	Proto      Proto
	Features   []byte
	Keys       []byte
	Values     []byte
}

type Config struct {
	TileID     m.TileID
	Name       string
	Extent     int32
	Version    int
	ReduceBool bool
	ExtentBool bool
	Tolerance  float64
	Proto      ProtoType
}

func NewLayer(tileid m.TileID, name string, pt ProtoType) LayerWrite {
	keys_map := map[string]uint32{}
	values_map := map[interface{}]uint32{}
	cur := NewCursor(tileid)
	proto := getProto(pt)
	return LayerWrite{TileID: tileid, Keys_Map: keys_map, Values_Map: values_map, Name: name, Cursor: cur, Buf: pbf.NewWriter(), Proto: proto}
}

func NewConfig(layername string, tileid m.TileID, pt ProtoType) Config {
	return Config{Name: layername, TileID: tileid, ExtentBool: true, Tolerance: 3, Proto: pt}
}

func NewLayerConfig(config Config) LayerWrite {
	keys_map := map[string]uint32{}
	values_map := map[interface{}]uint32{}
	if config.Extent == int32(0) {
		config.Extent = int32(4096)
	}
	if config.Version == 0 {
		config.Version = 2
	}
	proto := getProto(config.Proto)
	cur := NewCursorExtent(config.TileID, config.Extent)
	bds := m.Bounds(config.TileID)
	return LayerWrite{TileID: config.TileID,
		DeltaX:     bds.E - bds.W,
		DeltaY:     bds.N - bds.S,
		Keys_Map:   keys_map,
		Values_Map: values_map,
		Name:       config.Name,
		Cursor:     cur,
		Version:    config.Version,
		Extent:     int(config.Extent),
		ReduceBool: config.ReduceBool,
		Buf:        pbf.NewWriter(),
		Proto:      proto,
	}
}

func (layer *LayerWrite) AddKey(key string) uint32 {
	fwriter := pbf.NewWriter()
	fwriter.WriteString(layer.Proto.Layer.Keys, key)
	layer.Keys = append(layer.Keys, fwriter.Finish()...)
	myint := uint32(len(layer.Keys_Map))
	layer.Keys_Map[key] = myint
	return myint
}

func (layer *LayerWrite) AddValue(value interface{}) uint32 {
	fwriter := pbf.NewWriter()
	fwriter.WriteValue(layer.Proto.Layer.Values, value)
	layer.Values = append(layer.Values, fwriter.Finish()...)
	myint := uint32(len(layer.Values_Map))
	layer.Values_Map[value] = myint
	return myint
}

func (layer *LayerWrite) GetTags(properties map[string]interface{}) []uint32 {
	tags := make([]uint32, len(properties)*2)
	i := 0
	for k, v := range properties {
		keytag, keybool := layer.Keys_Map[k]
		if keybool == false {
			keytag = layer.AddKey(k)
		}
		valuetag, valuebool := layer.Values_Map[v]
		if valuebool == false {
			valuetag = layer.AddValue(v)
		}
		tags[i] = keytag
		tags[i+1] = valuetag
		i += 2
	}
	return tags
}

func (layer *LayerWrite) RefreshCursor() {
	layer.Cursor.Count = 0
	layer.Cursor.LastPoint = []int32{0, 0}
	layer.Cursor.Geometry = []uint32{}
	layer.Cursor.Bds = startbds
}

func WriteLayer(features []*geom.Feature, config Config) []byte {
	layer := NewLayerConfig(config)
	if config.ExtentBool {
		layer.Cursor.ExtentBool = true
	}

	for _, feat := range features {
		layer.AddFeature(feat)
	}
	if len(layer.Name) > 0 {
		layer.Buf.WriteString(layer.Proto.Layer.Name, layer.Name)
	}
	if len(layer.Features) > 0 {
		layer.Buf.WriteRaw(layer.Features)
	}
	if len(layer.Keys) > 0 {
		layer.Buf.WriteRaw(layer.Keys)
	}
	if len(layer.Values) > 0 {
		layer.Buf.WriteRaw(layer.Values)
	}

	layer.Buf.WriteUInt64(layer.Proto.Layer.Extent, uint64(layer.Extent))
	layer.Buf.WriteVarint(layer.Proto.Layer.Version, int(layer.Version))

	total_bytes := layer.Buf.Finish()

	tag := tagAndType(layer.Proto.Layers, pbf.Bytes)
	beg := append([]byte{tag}, pbf.EncodeVarint(uint64(len(total_bytes)))...)
	return append(beg, total_bytes...)
}

func (layer *LayerWrite) Flush() []byte {
	if len(layer.Name) > 0 {
		layer.Buf.WriteString(layer.Proto.Layer.Name, layer.Name)
	}
	if len(layer.Features) > 0 {
		layer.Buf.WriteRaw(layer.Features)
	}
	if len(layer.Keys) > 0 {
		layer.Buf.WriteRaw(layer.Keys)
	}
	if len(layer.Values) > 0 {
		layer.Buf.WriteRaw(layer.Values)
	}

	layer.Buf.WriteUInt64(layer.Proto.Layer.Extent, uint64(layer.Extent))
	layer.Buf.WriteVarint(layer.Proto.Layer.Version, int(layer.Version))

	total_bytes := layer.Buf.Finish()

	tag := tagAndType(layer.Proto.Layers, pbf.Bytes)
	beg := append([]byte{tag}, pbf.EncodeVarint(uint64(len(total_bytes)))...)
	return append(beg, total_bytes...)
}

func tagAndType(t pbf.TagType, w pbf.WireType) byte {
	return byte((uint32(t) << 3) | uint32(w))
}
