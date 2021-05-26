package geojsonvt

import (
	"fmt"

	m "github.com/flywave/go-mapbox/tileid"

	"github.com/flywave/go-mapbox/mvt"

	"reflect"
)

type Tile struct {
	NumPoints     int
	NumSimplified int
	NumFeatures   int
	Source        []Feature
	TileID        m.TileID
	LayerWrite    *mvt.LayerWrite
	Transform     bool
	MinX          float64
	MinY          float64
	MaxX          float64
	MaxY          float64
	Options       Config
	Tolerance     float64
}

func NewTile() Tile {
	return Tile{
		MinX: 2,
		MinY: 1,
		MaxX: -1,
		MaxY: 0,
	}
}

func CreateTile(features []Feature, tileid m.TileID, options Config) Tile {
	tile := NewTile()
	tile.TileID = tileid

	for _, feature := range features {

		tile.NumFeatures++

		minX := feature.MinX
		minY := feature.MinY
		maxX := feature.MaxX
		maxY := feature.MaxY
		if minX < tile.MinX {
			tile.MinX = minX
		}
		if minY < tile.MinY {
			tile.MinY = minY
		}
		if maxX > tile.MaxX {
			tile.MaxX = maxX
		}
		if maxY > tile.MaxY {
			tile.MaxY = maxY
		}
	}
	tile.Source = features
	tile.Options = options
	return tile
}

func (tile Tile) Marshal() []byte {

	var indexval int
	if len(tile.Source) > tile.Options.MaxFeatures {
		indexval = tile.Options.MaxFeatures
	} else {
		indexval = len(tile.Source)
	}
	layerwrite := mvt.NewLayerConfig(mvt.Config{TileID: tile.TileID, Name: tile.Options.LayerName})
	var tolerance float64
	if tile.Options.MaxZoom == 0 {
		tolerance = 0.0
	} else {
		tolerance = float64(tile.Options.Tolerance) / (float64(int(1<<tile.TileID.Z)) * float64(tile.Options.Extent))
	}
	count := 0
	for _, feature := range tile.Source[:indexval] {
		layerwrite.RefreshCursor()
		z2 := float64(int(1 << tile.TileID.Z))
		tx := float64(tile.TileID.X)
		ty := float64(tile.TileID.Y)
		simplified := [][]int32{}
		extent := float64(tile.Options.Extent)
		var geomtype int
		if feature.Type == "Point" {
			simplified = append(simplified, transformPoint(feature.Geometry.Point[0], feature.Geometry.Point[1], extent, z2, tx, ty))
			tile.NumPoints++
			tile.NumSimplified++
			layerwrite.Cursor.MakePoint(simplified[0])
			geomtype = 1

		} else if feature.Type == "MultiPoint" {
			for i := 0; i < len(feature.Geometry.MultiPoint); i += 3 {
				simplified = append(simplified, transformPoint(feature.Geometry.MultiPoint[i], feature.Geometry.MultiPoint[i+1], extent, z2, tx, ty))
				tile.NumPoints++
				tile.NumSimplified++
			}
			if len(simplified) > 0 {
				layerwrite.Cursor.MakeMultiPoint(simplified)
			}

			geomtype = 1

		} else if feature.Type == "LineString" {
			simplified = tile.addLine(feature.Geometry.LineString, tolerance, false, false, extent, z2, tx, ty)
			if len(simplified) > 1 {
				layerwrite.Cursor.MakeLine(simplified)
			}
			geomtype = 2

		} else if feature.Type == "MultiLineString" {
			newlist := make([][][]int32, len(feature.Geometry.MultiLineString))
			for pos, i := range feature.Geometry.MultiLineString {
				newlist[pos] = tile.addLine(i, tolerance, false, false, extent, z2, tx, ty)
			}
			if len(newlist) > 0 {
				layerwrite.Cursor.MakeMultiLine(newlist)
			}
			geomtype = 2
		} else if feature.Type == "Polygon" {
			newlist := make([][][]int32, len(feature.Geometry.Polygon))
			for pos, i := range feature.Geometry.Polygon {
				newlist[pos] = tile.addLine(i, tolerance, true, pos == 0, extent, z2, tx, ty)
			}
			if len(newlist) > 0 {
				layerwrite.Cursor.MakePolygon(newlist)
			}
			geomtype = 3
		} else if feature.Type == "MultiPolygon" {
			newlist := make([][][][]int32, len(feature.Geometry.MultiPolygon))
			for i := range feature.Geometry.MultiPolygon {
				newlist[i] = make([][][]int32, len(feature.Geometry.MultiPolygon[i]))
				for j := range feature.Geometry.MultiPolygon[i] {
					newlist[i][j] = tile.addLine(feature.Geometry.MultiPolygon[i][j], tolerance, true, j == 0, extent, z2, tx, ty)
				}
			}
			if len(newlist) > 0 {
				layerwrite.Cursor.MakeMultiPolygon(newlist)
			}
			geomtype = 3
		}
		var id int
		if feature.ID != nil {
			vv := reflect.ValueOf(feature.ID)
			kd := vv.Kind()
			switch kd {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				id = int(vv.Int())
			}
		}
		if len(layerwrite.Cursor.Geometry) > 0 {
			count++
			layerwrite.AddFeatureRaw(id, geomtype, layerwrite.Cursor.Geometry, feature.Tags)
		}
	}

	return layerwrite.Flush()
}

func (tile Tile) AddFeature(feature Feature) {
	tolerance := tile.Tolerance
	tile.LayerWrite.RefreshCursor()
	z2 := float64(int(1 << tile.TileID.Z))
	tx := float64(tile.TileID.X)
	ty := float64(tile.TileID.Y)
	simplified := [][]int32{}
	extent := float64(tile.Options.Extent)
	var geomtype int
	if feature.Type == "Point" {
		simplified = append(simplified, transformPoint(feature.Geometry.Point[0], feature.Geometry.Point[1], extent, z2, tx, ty))
		tile.NumPoints++
		tile.NumSimplified++
		tile.LayerWrite.Cursor.MakePoint(simplified[0])
		geomtype = 1

	} else if feature.Type == "MultiPoint" {
		for i := 0; i < len(feature.Geometry.MultiPoint); i += 3 {
			simplified = append(simplified, transformPoint(feature.Geometry.MultiPoint[i], feature.Geometry.MultiPoint[i+1], extent, z2, tx, ty))
			tile.NumPoints++
			tile.NumSimplified++
		}
		if len(simplified) > 0 {
			tile.LayerWrite.Cursor.MakeMultiPoint(simplified)
		}

		geomtype = 1

	} else if feature.Type == "LineString" {
		simplified = tile.addLine(feature.Geometry.LineString, tolerance, false, false, extent, z2, tx, ty)
		if len(simplified) > 1 {
			tile.LayerWrite.Cursor.MakeLine(simplified)
		}
		geomtype = 2

	} else if feature.Type == "MultiLineString" {
		newlist := make([][][]int32, len(feature.Geometry.MultiLineString))
		for pos, i := range feature.Geometry.MultiLineString {
			newlist[pos] = tile.addLine(i, tolerance, false, false, extent, z2, tx, ty)
		}
		if len(newlist) > 0 {
			tile.LayerWrite.Cursor.MakeMultiLine(newlist)
		}
		geomtype = 2
	} else if feature.Type == "Polygon" {
		newlist := make([][][]int32, len(feature.Geometry.Polygon))
		for pos, i := range feature.Geometry.Polygon {
			newlist[pos] = tile.addLine(i, tolerance, true, pos == 0, extent, z2, tx, ty)
		}
		if len(newlist) > 0 {
			tile.LayerWrite.Cursor.MakePolygon(newlist)
		}
		geomtype = 3
	} else if feature.Type == "MultiPolygon" {
		fmt.Println("here")

		newlist := make([][][][]int32, len(feature.Geometry.MultiPolygon))
		for i := range feature.Geometry.MultiPolygon {
			for j := range feature.Geometry.MultiPolygon[i] {
				newlist[i][j] = tile.addLine(feature.Geometry.MultiPolygon[i][j], tolerance, true, j == 0, extent, z2, tx, ty)
			}
		}
		if len(newlist) > 0 {
			tile.LayerWrite.Cursor.MakeMultiPolygon(newlist)
		}
		geomtype = 3
	}
	var id int
	if feature.ID != nil {
		vv := reflect.ValueOf(feature.ID)
		kd := vv.Kind()
		switch kd {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			id = int(vv.Int())
		}
	}
	if len(tile.LayerWrite.Cursor.Geometry) > 0 {
		tile.LayerWrite.AddFeatureRaw(id, geomtype, tile.LayerWrite.Cursor.Geometry, feature.Tags)
	}
}

func (tile *Tile) addLine(geom []float64, tolerance float64, isPolygon bool, isOuter bool, extent, z2, x, y float64) [][]int32 {
	var sqTolerance = tolerance * tolerance
	if isPolygon {
		tolerance = sqTolerance
	}
	ring := [][]int32{}
	for i := 0; i < len(geom); i += 3 {
		if tolerance == 0 || geom[i+2] > sqTolerance {
			tile.NumSimplified++
			ring = append(ring, transformPoint(geom[i], geom[i+1], extent, z2, x, y))
		}
		tile.NumPoints++
	}

	if isPolygon {
		if isOuter {
			ring = assert_winding_order(ring, "counter")
		} else {
			ring = assert_winding_order(ring, "clockwise")
		}
	}

	return ring
}

func reverse(coord [][]int32) [][]int32 {
	current := len(coord) - 1
	newlist := [][]int32{}
	for current != -1 {
		newlist = append(newlist, coord[current])
		current = current - 1
	}
	return newlist
}

func assert_winding_order(coord [][]int32, exp_orient string) [][]int32 {
	if len(coord) == 0 {
		return coord
	}
	count := 0
	firstpt := coord[0]
	weight := 0.0
	var oldpt []int32
	for _, pt := range coord {
		if count == 0 {
			count = 1
		} else {
			weight += float64((pt[0] - oldpt[0]) * (pt[1] + oldpt[1]))
		}
		oldpt = pt
	}

	weight += float64((firstpt[0] - oldpt[0]) * (firstpt[1] + oldpt[1]))
	var orientation string
	if weight > 0 {
		orientation = "clockwise"
	} else {
		orientation = "counter"
	}

	if orientation != exp_orient {
		return reverse(coord)
	} else {
		return coord
	}
}

func rewind(ring []float64, clockwise bool) []float64 {
	var area = 0.0
	i := 0
	len1 := len(ring)
	for j := len1 - 2; i < len1; i += 2 {
		var newj int
		if j > len(ring) {
			newj = j - len(ring)
		} else {
			newj = j
		}
		area += (ring[i] - ring[newj]) * (ring[i+1] + ring[newj+1])
	}

	ringsize := len(ring)
	if area > 0 == clockwise {
		for i = 0; i < ringsize/2; i += 2 {
			x := ring[i]
			y := ring[i+1]
			ring[i] = ring[ringsize-2-i]
			ring[i+1] = ring[ringsize-1-i]
			ring[ringsize-2-i] = x
			ring[ringsize-1-i] = y
		}
	}
	return ring
}
