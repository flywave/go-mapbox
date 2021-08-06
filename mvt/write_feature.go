package mvt

import (
	"reflect"

	"github.com/flywave/go-geom"
	"github.com/flywave/go-geom/general"
	"github.com/flywave/go-pbf"
)

func (layer *LayerWrite) AddFeature(feature *geom.Feature) {
	layer.RefreshCursor()

	fwriter := pbf.NewWriter()

	if feature.Geometry == nil {
		feature.Geometry = general.GeometryDataAsGeometry(&feature.GeometryData)
	}

	if feature.ID != nil {
		vv := reflect.ValueOf(feature.ID)
		kd := vv.Kind()
		switch kd {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fwriter.WriteUInt64(layer.Proto.Feature.ID, uint64(vv.Int()))
		}
	}

	if len(feature.Properties) > 0 {
		tags := layer.GetTags(feature.Properties)
		fwriter.WritePackedUInt32(layer.Proto.Feature.Tags, tags)
	}
	if feature.Geometry != nil {
		var geomtype byte
		switch (feature.Geometry).GetType() {
		case "Point", "MultiPoint":
			geomtype = 1
		case "LineString", "MultiLineString":
			geomtype = 2
		case "Polygon", "MultiPolygon":
			geomtype = 3
		}
		fwriter.WriteVarint(layer.Proto.Feature.Type, int(geomtype))
	}
	if feature.Geometry != nil {
		switch (feature.Geometry).GetType() {
		case "Point":
			layer.Cursor.MakePointFloat((feature.Geometry).(geom.Point).Data())
			fwriter.WritePackedUInt32(layer.Proto.Feature.Geometry, layer.Cursor.Geometry)
		case "LineString":
			layer.Cursor.MakeLineFloat((feature.Geometry).(geom.LineString).Data())
			fwriter.WritePackedUInt32(layer.Proto.Feature.Geometry, layer.Cursor.Geometry)
		case "Polygon":
			layer.Cursor.MakePolygonFloat((feature.Geometry).(geom.Polygon).Data())
			fwriter.WritePackedUInt32(layer.Proto.Feature.Geometry, layer.Cursor.Geometry)
		case "MultiPoint":
			layer.Cursor.MakeMultiPointFloat((feature.Geometry).(geom.MultiPoint).Data())
			fwriter.WritePackedUInt32(layer.Proto.Feature.Geometry, layer.Cursor.Geometry)
		case "MultiLineString":
			layer.Cursor.MakeMultiLineFloat((feature.Geometry).(geom.MultiLine).Data())
			fwriter.WritePackedUInt32(layer.Proto.Feature.Geometry, layer.Cursor.Geometry)
		case "MultiPolygon":
			layer.Cursor.MakeMultiPolygonFloat((feature.Geometry).(geom.MultiPolygon).Data())
			fwriter.WritePackedUInt32(layer.Proto.Feature.Geometry, layer.Cursor.Geometry)
		}
	}

	allbyte := fwriter.Finish()
	tag := tagAndType(layer.Proto.Layer.Features, pbf.Bytes)
	lens := pbf.EncodeVarint(uint64(len(allbyte)))
	layer.Features = append(layer.Features, appendAll([]byte{tag}, lens, allbyte)...)
}

func (layer *LayerWrite) AddFeatureRaw(id int, geomtype int, geometry []uint32, properties map[string]interface{}) {
	layer.RefreshCursor()

	fwriter := pbf.NewWriter()

	if id > 0 {
		fwriter.WriteUInt64(layer.Proto.Feature.ID, uint64(id))
	}

	if len(properties) > 0 {
		tags := layer.GetTags(properties)
		fwriter.WritePackedUInt32(layer.Proto.Feature.Tags, tags)
	}
	if geomtype != 0 {
		fwriter.WriteVarint(layer.Proto.Feature.Type, int(geomtype))
	}
	if len(geometry) > 0 {
		fwriter.WritePackedUInt32(layer.Proto.Feature.Geometry, geometry)
	}

	allbyte := fwriter.Finish()
	tag := tagAndType(layer.Proto.Layer.Features, pbf.Bytes)
	lens := pbf.EncodeVarint(uint64(len(allbyte)))
	layer.Features = append(layer.Features, appendAll([]byte{tag}, lens, allbyte)...)
}
