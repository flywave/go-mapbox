package mvt

import (
	"fmt"
	"io/ioutil"
	"math"

	m "github.com/flywave/go-mapbox/tileid"

	"github.com/flywave/go-geom"
)

type ClipGeom struct {
	Geom      [][]float64
	NewGeom   [][][]float64
	K1        float64
	K2        float64
	Axis      int
	IsPolygon bool
	SlicePos  int
}

type Slice struct {
	Pos   int
	Slice [][]float64
	Axis  int
}

func NewSlice(axis int) *Slice {
	return &Slice{Slice: [][]float64{}, Axis: axis}
}

func (slice *Slice) IntersectX(ax, ay, bx, by, x float64) float64 {
	t := (x - ax) / (bx - ax)
	if (bx - ax) == 0 {
		slice.Slice = append(slice.Slice, []float64{x, ay})
		slice.Pos += 1
		return 0.0
	} else {
		slice.Slice = append(slice.Slice, []float64{x, ay + (by-ay)*t})
		slice.Pos += 1
	}
	return t
}

func (slice *Slice) IntersectY(ax, ay, bx, by, y float64) float64 {
	t := (y - ay) / (by - ay)
	if (by - ay) == 0 {
		slice.Slice = append(slice.Slice, []float64{ax, y})
		slice.Pos += 1
		return 0.0
	} else {
		slice.Slice = append(slice.Slice, []float64{ax + (bx-ax)*t, y})
		slice.Pos += 1
	}
	return t
}

func (slice *Slice) Intersect(ax, ay, bx, by, val float64) float64 {
	if slice.Axis == 0 {
		return slice.IntersectX(ax, ay, bx, by, val)
	} else if slice.Axis == 1 {
		return slice.IntersectY(ax, ay, bx, by, val)
	}
	return 0.0
}

func (slice *Slice) AddPoint(x, y float64) {
	slice.Slice = append(slice.Slice, []float64{x, y})
	slice.Pos += 1
}

func (input *ClipGeom) clipLine() {
	slice := NewSlice(input.Axis)
	var ax, ay, bx, by, a, b float64
	k1, k2 := input.K1, input.K2
	for i := 0; i < len(input.Geom)-1; i += 1 {
		ax = input.Geom[i][0]
		ay = input.Geom[i][1]
		bx = input.Geom[i+1][0]
		by = input.Geom[i+1][1]
		if input.Axis == 0 {
			a = ax
			b = bx
		} else if input.Axis == 1 {
			a = ay
			b = by
		}
		exited := false

		if a <= k1 {
			if b >= k1 {
				slice.Intersect(ax, ay, bx, by, k1)
			}
		} else if a >= k2 {
			if b <= k2 {
				slice.Intersect(ax, ay, bx, by, k2)
			}
		} else {
			slice.AddPoint(ax, ay)
		}
		if b < k1 && a >= k1 {
			slice.Intersect(ax, ay, bx, by, k1)
			exited = true
		}
		if b > k2 && a <= k2 {
			slice.Intersect(ax, ay, bx, by, k2)
			exited = true
		}

		if !input.IsPolygon && exited {
			input.NewGeom = append(input.NewGeom, slice.Slice)
			slice.Slice = [][]float64{}
			slice.Pos = 0
		}
	}

	last := len(input.Geom) - 1
	ax = input.Geom[last][0]
	ay = input.Geom[last][1]
	if input.Axis == 0 {
		a = ax
	} else if input.Axis == 1 {
		a = ay
	}
	if a >= k1 && a <= k2 {
		slice.AddPoint(ax, ay)
	}

	last = len(slice.Slice) - 1
	if input.IsPolygon && last >= 3 && (slice.Slice[last][0] != slice.Slice[0][0] || slice.Slice[last][1] != slice.Slice[0][1]) {
		slice.AddPoint(slice.Slice[0][0], slice.Slice[0][1])
	}

	if slice.Pos > 0 {
		input.NewGeom = append(input.NewGeom, slice.Slice)
	}
}

func clipLine(geom [][]float64, k1, k2 float64, axis int, IsPolygon bool) [][][]float64 {
	clipthing := &ClipGeom{Geom: geom, K1: k1, K2: k2, Axis: axis, IsPolygon: IsPolygon}
	clipthing.clipLine()
	return clipthing.NewGeom
}

func clipLines(geom [][][]float64, k1, k2 float64, axis int, IsPolygon bool) [][][]float64 {
	if len(geom) == 0 {
		return [][][]float64{}
	}
	clipthing := &ClipGeom{Geom: geom[0], K1: k1, K2: k2, Axis: axis, IsPolygon: IsPolygon}
	for pos := range geom {
		clipthing.Geom = geom[pos]
		clipthing.clipLine()

	}
	return clipthing.NewGeom
}

func clipMultiPolygon(geom [][][][]float64, k1, k2 float64, axis int, IsPolygon bool) [][][][]float64 {
	newgeom := [][][][]float64{}
	for i := range geom {
		tmp := clipLines(geom[i], k1, k2, axis, IsPolygon)
		if len(tmp) > 0 {
			newgeom = append(newgeom, tmp)
		}
	}
	return newgeom
}

var Power7 = math.Pow(10, -7)

func DeltaFirstLast(firstpt, lastpt []float64) bool {
	dx, dy := math.Abs(firstpt[0]-lastpt[0]), math.Abs(firstpt[1]-lastpt[1])
	return dx > Power7 || dy > Power7
}

func LintPolygon(polygon [][][]float64) [][][]float64 {
	for pos, line := range polygon {
		fpt, lpt := line[0], line[len(line)-1]
		if DeltaFirstLast(fpt, lpt) {
			polygon[pos] = append(line, fpt)
		} else {
			line[len(line)-1] = fpt
			polygon[pos] = line
		}
	}
	return polygon
}

func LintMultiPolygon(multipolygon [][][][]float64) [][][][]float64 {
	for pos := range multipolygon {
		multipolygon[pos] = LintPolygon(multipolygon[pos])
	}
	return multipolygon
}

func clip(geom_ geom.GeometryData, k1, k2 float64, axis int) geom.GeometryData {
	var newgeom geom.GeometryData
	switch geom_.Type {
	case "LineString":
		newgeom.Type = "LineString"
		lines := clipLine(geom_.LineString, k1, k2, axis, false)
		if len(lines) == 1 {
			newgeom.LineString = lines[0]
		} else {
			newgeom.LineString = [][]float64{}
			newgeom.Type = "MultiLineString"
			newgeom.MultiLineString = lines
		}
	case "MultiLine":
		newgeom.Type = "MultiLineString"
		lines := clipLines(geom_.MultiLineString, k1, k2, axis, false)
		if len(lines) == 1 {
			newgeom.LineString = lines[0]
			newgeom.Type = "LineString"
			newgeom.MultiLineString = [][][]float64{}
		} else {
			newgeom.MultiLineString = lines
		}
	case "Polygon":
		newgeom.Type = "Polygon"
		geom_.Polygon = LintPolygon(geom_.Polygon)
		lines := clipLines(geom_.Polygon, k1, k2, axis, true)
		if len(lines) > 0 {
			newgeom.Polygon = lines
		}
	case "MultiPolygon":
		newgeom.Type = "MultiPolygon"
		geom_.MultiPolygon = LintMultiPolygon(geom_.MultiPolygon)
		multilines := clipMultiPolygon(geom_.MultiPolygon, k1, k2, axis, true)
		if len(multilines) == 1 {
			newgeom.Polygon = multilines[0]
			newgeom.Type = "Polygon"
			newgeom.MultiPolygon = [][][][]float64{}
		} else if len(multilines) > 0 {
			newgeom.MultiPolygon = multilines
		}
	}
	return newgeom
}

func IsEmpty(geom geom.GeometryData) bool {
	switch geom.Type {
	case "Point":
		return true
	case "MultiPoint":
		return len(geom.MultiPoint) == 0
	case "LineString":
		return len(geom.LineString) == 0
	case "MultiLineString":
		return len(geom.MultiLineString) == 0
	case "Polygon":
		return len(geom.Polygon) == 0
	case "MultiPolygon":
		return len(geom.MultiPolygon) == 0
	}
	return false
}

func PointClipAboutTile(feature *geom.Feature, tileid m.TileID) *geom.Feature {
	if feature.GeometryData.Type == "Point" {
		checktileid := m.Tile(feature.GeometryData.Point[0], feature.GeometryData.Point[1], int(tileid.Z))
		if m.IsEqual(checktileid, tileid) {
			feature.Properties["TILEID"] = tileid
			return feature
		}

		return &geom.Feature{}
	} else if feature.GeometryData.Type == "MultiPoint" {
		newpoints := [][]float64{}
		for _, pt := range feature.GeometryData.MultiPoint {
			checktileid := m.Tile(pt[0], pt[1], int(tileid.Z))
			if m.IsEqual(checktileid, tileid) {
				newpoints = append(newpoints, pt)
			}
		}
		if len(newpoints) > 0 {
			if len(newpoints) == 0 {
				newfeature := geom.NewPointFeature(newpoints[0])
				newfeature.Properties = feature.Properties
				newfeature.Properties["TILEID"] = tileid
				return newfeature
			} else {
				newfeature := geom.NewMultiPointFeature(newpoints...)
				newfeature.Properties = feature.Properties
				newfeature.Properties["TILEID"] = tileid
				return newfeature
			}
		}
	}
	return &geom.Feature{}
}

func PointClipAboutZoom(feature *geom.Feature, zoom int) map[m.TileID]*geom.Feature {
	if feature.GeometryData.Type == "Point" {
		checktileid := m.Tile(feature.GeometryData.Point[0], feature.GeometryData.Point[1], zoom)
		feature.Properties["TILEID"] = checktileid
		return map[m.TileID]*geom.Feature{checktileid: feature}
	} else if feature.GeometryData.Type == "MultiPoint" {
		newpoints := map[m.TileID][][]float64{}
		for _, pt := range feature.GeometryData.MultiPoint {
			checktileid := m.Tile(pt[0], pt[1], zoom)
			newpoints[checktileid] = append(newpoints[checktileid], pt)
		}
		totalmap := map[m.TileID]*geom.Feature{}
		for k, newpoints2 := range newpoints {
			if len(newpoints2) > 0 {
				if len(newpoints2) == 0 {
					newfeature := geom.NewPointFeature(newpoints2[0])
					newfeature.Properties = feature.Properties
					newfeature.Properties["TILEID"] = k
					totalmap[k] = newfeature
				} else {
					newfeature := geom.NewMultiPointFeature(newpoints2...)
					newfeature.Properties = feature.Properties
					newfeature.Properties["TILEID"] = k
					totalmap[k] = newfeature
				}
			}
		}
		return totalmap
	}
	return map[m.TileID]*geom.Feature{}
}

func makefeature(addgeom geom.GeometryData, prop map[string]interface{}, id interface{}) *geom.Feature {
	feat2 := &geom.Feature{GeometryData: geom.GeometryData{}}
	feat2.GeometryData.Type = addgeom.Type
	switch feat2.GeometryData.Type {
	case "Point":
		feat2.GeometryData.Point = addgeom.Point
	case "MultiPoint":
		feat2.GeometryData.MultiPoint = addgeom.MultiPoint
	case "LineString":
		feat2.GeometryData.LineString = addgeom.LineString
	case "MultiLineString":
		feat2.GeometryData.MultiLineString = addgeom.MultiLineString
	case "Polygon":
		if len(addgeom.Polygon[0][0]) == 8 && len(addgeom.BoundingBox) == 4 {
			bb := addgeom.BoundingBox
			w, s, e, n := bb[0], bb[1], bb[2], bb[3]
			poly := [][][]float64{{{e, n}, {w, n}, {w, s}, {e, s}, {e, n}}}
			addgeom.Polygon = poly
		}
		feat2.GeometryData.Polygon = addgeom.Polygon

	case "MultiPolygon":
		feat2.GeometryData.MultiPolygon = addgeom.MultiPolygon
	}
	feat2.Properties = prop
	feat2.ID = id
	feat2.BoundingBox = geom.BoundingBoxFromGeometryData(&addgeom)
	return feat2
}

func ClipTile(feature *geom.Feature, tileid m.TileID) *geom.Feature {
	gtype := string(feature.GeometryData.Type)
	if gtype == "Point" || gtype == "MultiPoint" {
		return PointClipAboutTile(feature, tileid)
	}
	addgeom := feature.GeometryData
	bds := m.Bounds(tileid)
	addgeom = clip(addgeom, bds.W, bds.E, 0)
	addgeom = clip(addgeom, bds.S, bds.N, 1)
	return makefeature(addgeom, feature.Properties, feature.ID)
}

func getbounds(tileid m.TileID) []float64 {
	bds := m.Bounds(tileid)
	return []float64{bds.W, bds.S, bds.E, bds.N}
}

var squaregeom = geom.GeometryData{Type: "Polygon", Polygon: [][][]float64{{{100.0, 100.0, 100.0, 100.0, 100.0, 100.0, 100.0, 100.0}}}}

func getgeomsquaretile(tileid m.TileID) geom.GeometryData {
	var val geom.GeometryData
	val.Type = squaregeom.Type
	val.Polygon = squaregeom.Polygon
	val.BoundingBox = getbounds(tileid)
	return val
}

func ClipDownTile(geom_ geom.GeometryData, tileid m.TileID) map[m.TileID]geom.GeometryData {
	bds := m.Bounds(tileid)
	cs := m.Children(tileid)
	cbds := m.Bounds(cs[0])
	if geom_.Type == "Polygon" {
		if len(geom_.Polygon[0][0]) == 8 {
			return map[m.TileID]geom.GeometryData{
				cs[0]: getgeomsquaretile(cs[0]),
				cs[1]: getgeomsquaretile(cs[1]),
				cs[2]: getgeomsquaretile(cs[2]),
				cs[3]: getgeomsquaretile(cs[3]),
			}
		}
	}
	if geom_.Type == "Polygon" {
		if len(geom_.Polygon) == 1 {
			if len(geom_.Polygon[0]) == 4 || len(geom_.Polygon[0]) == 5 {
				bbb := geom.BoundingBoxFromGeometryData(&geom_)
				bdsref := m.Extrema{W: bbb[0], S: bbb[1], E: bbb[2], N: bbb[3]}

				if DeltaBounds(bdsref, bds) {
					return map[m.TileID]geom.GeometryData{
						cs[0]: getgeomsquaretile(cs[0]),
						cs[1]: getgeomsquaretile(cs[1]),
						cs[2]: getgeomsquaretile(cs[2]),
						cs[3]: getgeomsquaretile(cs[3]),
					}
				}

			}
		}
	}

	mx, my := cbds.E, cbds.S

	l, r := clip(geom_, bds.W, mx, 0), clip(geom_, mx, bds.E, 0)
	ld, lu := clip(l, bds.S, my, 1), clip(l, my, bds.N, 1)
	rd, ru := clip(r, bds.S, my, 1), clip(r, my, bds.N, 1)

	lut, rut, rdt, ldt := cs[0], cs[1], cs[2], cs[3]

	mymap := map[m.TileID]geom.GeometryData{
		lut: lu,
		rut: ru,
		rdt: rd,
		ldt: ld,
	}

	for k, v := range mymap {
		if IsEmpty(v) {
			delete(mymap, k)
		}
	}

	return mymap
}

func GetFirstZoom(bb m.Extrema) (int, m.TileID) {
	corners := [][]float64{{bb.E, bb.N}, {bb.E, bb.S}, {bb.W, bb.N}, {bb.W, bb.S}}
	for i := 0; i < 30; i++ {
		mymap := map[m.TileID]string{}
		for _, corner := range corners {
			mymap[m.Tile(corner[0], corner[1], i)] = ""
		}

		if len(mymap) > 1 {
			return i - 1, m.Tile(corners[0][0], corners[0][1], i-1)
		}
	}
	return 30, m.TileID{X: 0, Y: 0, Z: 30}
}

func DeltaBounds(bds1, bds2 m.Extrema) bool {
	de, dw, dn, ds := math.Abs(bds1.E-bds2.E), math.Abs(bds1.W-bds2.W), math.Abs(bds1.N-bds2.N), math.Abs(bds1.S-bds2.S)
	return dw < Power7 && de < Power7 && dn < Power7 && ds < Power7
}

func ClipFeature(feature *geom.Feature, endzoom int, keep_parents bool) map[m.TileID]*geom.Feature {
	gtype := string(feature.GeometryData.Type)
	if gtype == "Point" || gtype == "MultiPoint" {
		return PointClipAboutZoom(feature, endzoom)
	}
	geom_ := feature.GeometryData
	bb := geom.BoundingBoxFromGeometryData(&geom_)
	firstzoom, tileid := GetFirstZoom(m.Extrema{W: bb[0], S: bb[1], E: bb[2], N: bb[3]})
	currentzoom := firstzoom
	mymap := map[m.TileID]*geom.Feature{tileid: feature}

	if currentzoom >= endzoom {
		for int(tileid.Z) != endzoom {
			tileid = m.Parent(tileid)
		}
		mymap[tileid] = makefeature(feature.GeometryData, feature.Properties, feature.ID)
		currentzoom = endzoom
	}

	for currentzoom < endzoom {
		var lastk m.TileID
		for k, tempgeom := range mymap {
			if int(k.Z) == currentzoom {
				tmap := ClipDownTile(tempgeom.GeometryData, k)
				for myk, addgeom := range tmap {
					if (myk.Z) != 0 {
						lastk = myk
					}
					mymap[myk] = makefeature(addgeom, feature.Properties, feature.ID)
				}

				if !keep_parents {
					delete(mymap, k)
				}
			}
		}
		currentzoom = int(lastk.Z)
	}

	return mymap
}

func ReadFeatures(filename string) []*geom.Feature {
	bs, _ := ioutil.ReadFile(filename)
	fc, _ := geom.UnmarshalFeatureCollection(bs)
	return fc.Features
}

func MakeFeatures(feats []*geom.Feature, filename string) {
	fc := geom.NewFeatureCollection()
	fc.Features = feats
	s, err := fc.MarshalJSON()
	fmt.Println(err)
	ioutil.WriteFile(filename, s, 0677)
}

func NewFeature(geom_ geom.GeometryData, props map[string]interface{}) *geom.Feature {
	feat2 := geom.NewFeatureFromGeometryData(&geom_)
	feat2.Properties = map[string]interface{}{}
	for k, v := range props {
		feat2.Properties[k] = v
	}
	feat2.Properties["COLORKEY"] = "white"
	return feat2
}

func NewFeatures(mymap map[m.TileID]geom.GeometryData, props map[string]interface{}) []*geom.Feature {
	feats := make([]*geom.Feature, len(mymap))
	i := 0
	for k, geom_ := range mymap {
		feat2 := &geom.Feature{GeometryData: geom.GeometryData{}}
		feat2.GeometryData.Type = geom_.Type
		switch feat2.GeometryData.Type {
		case "Point":
			feat2.GeometryData.Point = geom_.Point
		case "MultiPoint":
			feat2.GeometryData.MultiPoint = geom_.MultiPoint
		case "LineString":
			feat2.GeometryData.LineString = geom_.LineString
		case "MultiLineString":
			feat2.GeometryData.MultiLineString = geom_.MultiLineString
		case "Polygon":
			feat2.GeometryData.Polygon = geom_.Polygon
		case "MultiPolygon":
			feat2.GeometryData.MultiPolygon = geom_.MultiPolygon
		}
		feat2.Properties = map[string]interface{}{}
		for k, v := range props {
			feat2.Properties[k] = v
		}
		feat2.Properties["TILEID"] = m.Tilestr(k)
		feat2.Properties["COLORKEY"] = "white"
		feats[i] = feat2
		i++
	}
	return feats
}
