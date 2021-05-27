package mvt

import (
	"github.com/flywave/go-geom"
)

/**
var array1, array2, array3, array4, array5, array6, array7, array8, array9 []byte
	layer.RefreshCursor()

	if feature.ID != nil {
		vv := reflect.ValueOf(feature.ID)
		kd := vv.Kind()
		switch kd {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			array3 = []byte{8}
			array4 = p.EncodeVarint(uint64(vv.Int()))
		}
	}

	if len(feature.Properties) > 0 {
		array5 = []byte{18}
		array6 = WritePackedUint32(layer.GetTags(feature.Properties))
	}
	if feature.Geometry != nil {
		var geomtype byte
		switch (*feature.Geometry).GetType() {
		case "Point", "MultiPoint":
			geomtype = 1
		case "LineString", "MultiLineString":
			geomtype = 2
		case "Polygon", "MultiPolygon":
			geomtype = 3
		}
		array7 = []byte{24, geomtype}
	}
	var abortBool bool
	if feature.Geometry != nil {
		switch (*feature.Geometry).GetType() {
		case "Point":
			array8 = []byte{34}
			layer.Cursor.MakePointFloat((*feature.Geometry).(geom.Point).Data())
			array9 = WritePackedUint32(layer.Cursor.Geometry)
		case "LineString":
			array8 = []byte{34}
			layer.Cursor.MakeLineFloat((*feature.Geometry).(geom.LineString).Data())
			if layer.Cursor.Count == 0 {
				abortBool = true
			}
			array9 = WritePackedUint32(layer.Cursor.Geometry)
		case "Polygon":
			array8 = []byte{34}
			layer.Cursor.MakePolygonFloat((*feature.Geometry).(geom.Polygon).Data())
			array9 = WritePackedUint32(layer.Cursor.Geometry)
		case "MultiPoint":
			array8 = []byte{34}
			layer.Cursor.MakeMultiPointFloat((*feature.Geometry).(geom.MultiPoint).Data())
			array9 = WritePackedUint32(layer.Cursor.Geometry)
		case "MultiLineString":
			array8 = []byte{34}
			layer.Cursor.MakeMultiLineFloat((*feature.Geometry).(geom.MultiLine).Data())
			array9 = WritePackedUint32(layer.Cursor.Geometry)
		case "MultiPolygon":
			array8 = []byte{34}
			layer.Cursor.MakeMultiPolygonFloat((*feature.Geometry).(geom.MultiPolygon).Data())
			array9 = WritePackedUint32(layer.Cursor.Geometry)
		}
	}

	if !abortBool {
		array1 = []byte{18}
		array2 = p.EncodeVarint(uint64(len(array3) + len(array4) + len(array5) + len(array6) + len(array7) + len(array8) + len(array9)))
		layer.Features = append(layer.Features, AppendAll(array1, array2, array3, array4, array5, array6, array7, array8, array9)...)
	}**/

func (layer *LayerWrite) AddFeature(feature *geom.Feature) {

}

/**
var array1, array2, array3, array4, array5, array6, array7, array8, array9 []byte

	layer.RefreshCursor()

	if id > 0 {
		array3 = []byte{8}
		array4 = p.EncodeVarint(uint64(id))
	}

	if len(properties) > 0 {
		array5 = []byte{18}
		array6 = WritePackedUint32(layer.GetTags(properties))
	}
	if geomtype != 0 {
		array7 = []byte{24, byte(geomtype)}
	}
	if len(geometry) > 0 {
		array8 = []byte{34}
		array9 = WritePackedUint32(geometry)
	}

	array1 = []byte{18}
	array2 = p.EncodeVarint(uint64(len(array3) + len(array4) + len(array5) + len(array6) + len(array7) + len(array8) + len(array9)))
	layer.Features = append(layer.Features, AppendAll(array1, array2, array3, array4, array5, array6, array7, array8, array9)...)
**/

func (layer *LayerWrite) AddFeatureRaw(id int, geomtype int, geometry []uint32, properties map[string]interface{}) {

}
