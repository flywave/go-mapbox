package mvt

import (
	"github.com/flywave/go-pbf"
)

const (
	GeomTypeUnknown    = 0
	GeomTypePoint      = 1
	GeomTypeLineString = 2
	GeomTypePolygon    = 3
)

type ProtoValue struct {
	StringValue  pbf.TagType
	FloatValue   pbf.TagType
	DoubleValue  pbf.TagType
	IntValue     pbf.TagType
	UIntValue    pbf.TagType
	SIntValue    pbf.TagType
	BoolIntValue pbf.TagType
}

type ProtoFeature struct {
	ID       pbf.TagType
	Tags     pbf.TagType
	Type     pbf.TagType
	Geometry pbf.TagType
}

type ProtoLayer struct {
	Version  pbf.TagType
	Name     pbf.TagType
	Features pbf.TagType
	Keys     pbf.TagType
	Values   pbf.TagType
	Extent   pbf.TagType
}

type Proto struct {
	Layers pbf.TagType

	Layer   ProtoLayer
	Feature ProtoFeature
	Value   ProtoValue
}

var (
	MapboxProto = Proto{Layers: 3, Layer: ProtoLayer{Version: 15, Name: 1, Features: 2, Keys: 3, Values: 4, Extent: 5}, Feature: ProtoFeature{ID: 1, Tags: 2, Type: 3, Geometry: 4}, Value: ProtoValue{StringValue: 1, FloatValue: 2, DoubleValue: 3, IntValue: 4, UIntValue: 5, SIntValue: 6, BoolIntValue: 7}}
)

var (
	LKProto = Proto{Layers: 2, Layer: ProtoLayer{Version: 15, Name: 1, Features: 2, Keys: 4, Values: 5, Extent: 6}, Feature: ProtoFeature{ID: 1, Tags: 7, Type: 6, Geometry: 2}, Value: ProtoValue{StringValue: 1, FloatValue: 2, DoubleValue: 3, IntValue: 4, UIntValue: 5, SIntValue: 6, BoolIntValue: 7}}
)

type ProtoType int

const (
	PROTO_MAPBOX ProtoType = 0
	PROTO_LK     ProtoType = 1
)

func getProto(p ProtoType) Proto {
	switch p {
	case PROTO_MAPBOX:
		return MapboxProto
	case PROTO_LK:
		return LKProto
	}
	return MapboxProto
}
