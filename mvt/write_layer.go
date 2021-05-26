package mvt

import (
	"errors"

	m "github.com/flywave/go-mapbox/tileid"

	"github.com/flywave/go-geom"

	"github.com/murphy214/pbf"
)

type LayerWrite struct {
	TileID       m.TileID
	DeltaX       float64
	DeltaY       float64
	Name         string
	Extent       int
	Version      int
	Keys_Map     map[string]uint32
	Keys_Bytes   []byte
	Values_Map   map[interface{}]uint32
	Values_Bytes []byte
	Features     []byte
	Cursor       *Cursor
	ReduceBool   bool
}

type Config struct {
	TileID     m.TileID
	Name       string
	Extent     int32
	Version    int
	ReduceBool bool
	ExtentBool bool
	Tolerance  float64
}

func NewLayer(tileid m.TileID, name string) LayerWrite {
	keys_map := map[string]uint32{}
	values_map := map[interface{}]uint32{}
	cur := NewCursor(tileid)
	return LayerWrite{TileID: tileid, Keys_Map: keys_map, Values_Map: values_map, Name: name, Cursor: cur}
}

func NewConfig(layername string, tileid m.TileID) Config {
	return Config{Name: layername, TileID: tileid, ExtentBool: true, Tolerance: 3}
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
	}
}

func (layer *LayerWrite) AddKey(key string) uint32 {
	layer.Keys_Bytes = append(layer.Keys_Bytes, 26)
	layer.Keys_Bytes = append(layer.Keys_Bytes, pbf.EncodeVarint(uint64(len(key)))...)
	layer.Keys_Bytes = append(layer.Keys_Bytes, []byte(key)...)
	myint := uint32(len(layer.Keys_Map))
	layer.Keys_Map[key] = myint
	return myint
}

func (layer *LayerWrite) AddValue(value interface{}) uint32 {
	layer.Values_Bytes = append(layer.Values_Bytes, WriteValue(value)...)
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

func WriteLayer(features []*geom.Feature, config Config) (total_bytes []byte, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("Error in WriteLayer().")
		}
	}()

	mylayer := NewLayerConfig(config)
	if config.ExtentBool {
		mylayer.Cursor.ExtentBool = true
	}

	for _, feat := range features {
		mylayer.AddFeature(feat)
	}

	if len(mylayer.Name) > 0 {
		total_bytes = append(total_bytes, 10)
		total_bytes = append(total_bytes, pbf.EncodeVarint(uint64(len(mylayer.Name)))...)
		total_bytes = append(total_bytes, []byte(mylayer.Name)...)
	}

	total_bytes = append(total_bytes, mylayer.Features...)
	total_bytes = append(total_bytes, mylayer.Keys_Bytes...)
	total_bytes = append(total_bytes, mylayer.Values_Bytes...)

	if true {
		total_bytes = append(total_bytes, 40)
		total_bytes = append(total_bytes, pbf.EncodeVarint(uint64(mylayer.Extent))...)
	}

	total_bytes = append(total_bytes, 120)
	total_bytes = append(total_bytes, byte(mylayer.Version))

	beg := append([]byte{26}, pbf.EncodeVarint(uint64(len(total_bytes)))...)
	total_bytes = append(beg, total_bytes...)
	return total_bytes, err
}

func (mylayer *LayerWrite) Flush() []byte {
	total_bytes := []byte{}

	if len(mylayer.Name) > 0 {
		total_bytes = append(total_bytes, 10)
		total_bytes = append(total_bytes, pbf.EncodeVarint(uint64(len(mylayer.Name)))...)
		total_bytes = append(total_bytes, []byte(mylayer.Name)...)
	}

	total_bytes = append(total_bytes, mylayer.Features...)
	total_bytes = append(total_bytes, mylayer.Keys_Bytes...)
	total_bytes = append(total_bytes, mylayer.Values_Bytes...)

	if true {
		total_bytes = append(total_bytes, 40)
		total_bytes = append(total_bytes, pbf.EncodeVarint(uint64(mylayer.Extent))...)
	}

	total_bytes = append(total_bytes, 120)
	total_bytes = append(total_bytes, byte(mylayer.Version))

	beg := append([]byte{26}, pbf.EncodeVarint(uint64(len(total_bytes)))...)
	return append(beg, total_bytes...)
}
