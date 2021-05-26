package mvt

import (
	"errors"
	"math"

	"github.com/flywave/go-geom"
	m "github.com/flywave/go-mapbox/tileid"

	"github.com/murphy214/pbf"
)

type Feature struct {
	ID          int
	Type        string
	Properties  map[string]interface{}
	GeometryPos int
	extent      int
	GeomInt     int
	Buf         *pbf.PBF
}

func DeltaDim(num int) float64 {
	if num%2 == 1 {
		return float64((num + 1) / -2)
	} else {
		return float64(num / 2)
	}
	return float64(0)
}

func SignedArea(ring [][]float64) float64 {
	sum := 0.0
	i := 0
	lenn := len(ring)
	j := lenn - 1
	var p1, p2 []float64

	for i < lenn {
		if i != 0 {
			j = i - 1
		}
		p1 = ring[i]
		p2 = ring[j]
		sum += (p2[0] - p1[0]) * (p1[1] + p2[1])
		i++
	}
	return sum
}

func Project(line [][]float64, x0 float64, y0 float64, size float64) [][]float64 {
	for j := range line {
		p := line[j]
		y2 := 180.0 - (float64(p[1])+y0)*360.0/size
		line[j] = []float64{
			(float64(p[0])+x0)*360.0/size - 180.0,
			360.0/math.Pi*math.Atan(math.Exp(y2*math.Pi/180.0)) - 90.0}
	}
	return line
}

func (layer *Layer) Feature() (feature *Feature, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("Error in Feature()")
			layer.featurePosition++
		}
	}()

	layer.Buf.Pos = layer.features[layer.featurePosition]
	endpos := layer.Buf.Pos + layer.Buf.ReadVarint()

	feature = &Feature{Properties: map[string]interface{}{}}

	for layer.Buf.Pos < endpos {
		key, val := layer.Buf.ReadKey()

		if key == 1 && val == 0 {
			feature.ID = int(layer.Buf.ReadUInt64())
		}
		if key == 2 && val == 2 {
			tags := layer.Buf.ReadPackedUInt32()
			i := 0
			for i < len(tags) {
				var key string
				if len(layer.Keys) <= int(tags[i]) {
					key = ""
				} else {
					key = layer.Keys[tags[i]]
				}
				var val interface{}
				if len(layer.Values) <= int(tags[i+1]) {
					val = ""
				} else {
					val = layer.Values[tags[i+1]]
				}
				feature.Properties[key] = val
				i += 2
			}
		}
		if key == 3 && val == 0 {
			geomType := int(layer.Buf.Varint()[0])
			feature.GeomInt = geomType
			switch geomType {
			case 1:
				feature.Type = "Point"
			case 2:
				feature.Type = "LineString"
			case 3:
				feature.Type = "Polygon"
			}
		}
		if key == 4 && val == 2 {
			feature.GeometryPos = layer.Buf.Pos
			size := layer.Buf.ReadVarint()
			layer.Buf.Pos += size + 1
		}
	}
	feature.extent = layer.Extent
	feature.Buf = layer.Buf
	layer.featurePosition += 1
	return feature, err
}

func (feature *Feature) LoadGeometryRaw() (geom []uint32, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("Error in feature.LoadGeometry()")
		}
	}()

	feature.Buf.Pos = feature.GeometryPos
	geom = feature.Buf.ReadPackedUInt32()
	return geom, err
}

func (feature *Feature) LoadGeometry() (geomm *geom.GeometryData, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("Error in feature.LoadGeometry()")
		}
	}()

	feature.Buf.Pos = feature.GeometryPos
	geom_ := feature.Buf.ReadPackedUInt32()

	pos := 0
	var lines [][][]float64
	var polygons [][][][]float64
	var firstpt []float64
	geomType := feature.GeomInt

	for pos < len(geom_) {
		if geom_[pos] == 9 {
			pos += 1
			if pos != 1 && geomType == 2 {
				firstpt = []float64{firstpt[0] + DeltaDim(int(geom_[pos])), firstpt[1] + DeltaDim(int(geom_[pos+1]))}
			} else {
				firstpt = []float64{DeltaDim(int(geom_[pos])), DeltaDim(int(geom_[pos+1]))}
			}
			pos += 2
			if len(geom_) == 3 {
				lines = [][][]float64{{firstpt}}
			}
			if pos < len(geom_) {
				cmdLen := geom_[pos]
				length := int(cmdLen >> 3)

				line := make([][]float64, length+1)
				pos += 1
				endpos := pos + length*2
				line[0] = firstpt
				i := 1
				for pos < endpos && pos+1 < len(geom_) {
					firstpt = []float64{firstpt[0] + DeltaDim(int(geom_[pos])), firstpt[1] + DeltaDim(int(geom_[pos+1]))}
					line[i] = firstpt
					i++
					pos += 2
				}
				lines = append(lines, line[:i])
				line = [][]float64{firstpt}

			} else {
				pos += 1
			}

		} else if pos < len(geom_) {
			if geom_[pos] == 15 {
				pos += 1
			} else {
				pos += 1
			}
		} else {
			pos += 1
		}
	}
	if geomType == 3 {
		for pos, line := range lines {
			f, l := line[0], line[len(line)-1]
			if !(f[0] == l[0] && l[1] == f[1]) {
				line = append(line, line[0])
			}
			lines[pos] = line
		}

		if len(lines) == 1 {
			polygons = append(polygons, lines)
		} else {
			for _, line := range lines {
				if len(line) > 0 {
					val := SignedArea(line)
					if val < 0 {
						polygons = append(polygons, [][][]float64{line})
					} else {
						if len(polygons) == 0 {
							polygons = append(polygons, [][][]float64{line})

						} else {
							polygons[len(polygons)-1] = append(polygons[len(polygons)-1], line)

						}
					}
				}
			}
		}
	} else {
		polygons = append(polygons, lines)
	}

	switch geomType {
	case 1:
		if len(polygons[0][0]) == 1 {
			geomm = geom.NewPointGeometryData(polygons[0][0][0])
		} else {
			geomm = geom.NewMultiPointGeometryData(polygons[0][0]...)
		}
	case 2:
		if len(polygons[0]) == 1 {
			geomm = geom.NewLineStringGeometryData(polygons[0][0])
		} else {
			geomm = geom.NewMultiLineStringGeometryData(polygons[0]...)
		}
	case 3:
		if len(polygons) == 1 {
			geomm = geom.NewPolygonGeometryData(polygons[0])
		} else {
			geomm = geom.NewMultiPolygonGeometryData(polygons...)
		}
	}

	return geomm, err
}

func (feature *Feature) ToGeoJSON(tile m.TileID) (*geom.Feature, error) {
	var err error
	defer func() {
		if recover() != nil {
			err = errors.New("Error in feature.ToGeoJSON()")
		}
	}()
	extent := feature.extent
	size := float64(extent) * float64(math.Pow(2, float64(tile.Z)))
	x0 := float64(extent) * float64(tile.X)
	y0 := float64(extent) * float64(tile.Y)
	geometry, err := feature.LoadGeometry()
	if err != nil {
		return &geom.Feature{}, err
	}
	switch geometry.Type {
	case "Point":
		geometry.Point = Project([][]float64{geometry.Point}, x0, y0, size)[0]
	case "MultiPoint":
		geometry.MultiPoint = Project(geometry.MultiPoint, x0, y0, size)
	case "LineString":
		geometry.LineString = Project(geometry.LineString, x0, y0, size)
	case "MultiLineString":
		for i := range geometry.MultiLineString {
			geometry.MultiLineString[i] = Project(geometry.MultiLineString[i], x0, y0, size)
		}
	case "Polygon":
		for i := range geometry.Polygon {
			geometry.Polygon[i] = Project(geometry.Polygon[i], x0, y0, size)
		}
	case "MultiPolygon":
		for i := range geometry.MultiPolygon {
			for j := range geometry.MultiPolygon[i] {
				geometry.MultiPolygon[i][j] = Project(geometry.MultiPolygon[i][j], x0, y0, size)
			}
		}
	}

	newFeature := geom.NewFeatureFromGeometryData(geometry)
	newFeature.Properties = feature.Properties
	newFeature.ID = feature.ID

	return newFeature, err
}

func convertpt(pt []float64, dim float64) []float64 {
	if pt[0] < 0 {
		//pt[0] = 0
	}
	if pt[1] < 0 {
		//pt[1] = 0
	}
	return []float64{pbf.Round(pt[0]/dim, .5, 0), pbf.Round(pt[1]/dim, .5, 0)}
}

func convertln(ln [][]float64, dim float64) [][]float64 {
	for i := range ln {
		ln[i] = convertpt(ln[i], dim)
	}
	return ln
}

func convertlns(lns [][][]float64, dim float64) [][][]float64 {
	for i := range lns {
		lns[i] = convertln(lns[i], dim)
	}
	return lns
}

func ConvertGeometry(geom_ *geom.GeometryData, dimf float64) *geom.GeometryData {
	if geom_ == nil {
		return &geom.GeometryData{}
	}
	switch geom_.Type {
	case "Point":
		geom_.Point = convertpt(geom_.Point, dimf)
	case "MultiPoint":
		geom_.MultiPoint = convertln(geom_.MultiPoint, dimf)
	case "LineString":
		geom_.LineString = convertln(geom_.LineString, dimf)
	case "MultiLineString":
		geom_.MultiLineString = convertlns(geom_.MultiLineString, dimf)
	case "Polygon":
		geom_.Polygon = convertlns(geom_.Polygon, dimf)
	case "MultiPolygon":
		for i := range geom_.MultiPolygon {
			geom_.MultiPolygon[i] = convertlns(geom_.MultiPolygon[i], dimf)
		}
	}
	return geom_
}

func (feature *Feature) LoadGeometryScaled(dim float64) (geomm *geom.GeometryData, err error) {
	geom, err := feature.LoadGeometry()
	geom2 := ConvertGeometry(geom, dim)
	return geom2, err
}
