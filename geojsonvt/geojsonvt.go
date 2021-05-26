package geojsonvt

import (
	"math"

	"github.com/flywave/go-geom"
	m "github.com/flywave/go-mapbox/tileid"
)

type Config struct {
	MaxZoom        int
	MinZoom        int
	IndexMaxZoom   uint64
	IndexMaxPoints int
	Tolerance      int
	LayerName      string
	Extent         int
	Buffer         int
	Debug          int
	MaxFeatures    int
	PropertiesBool bool
}

func NewConfig() Config {
	return Config{
		MaxZoom:        14,
		IndexMaxZoom:   5,
		IndexMaxPoints: 100000,
		Tolerance:      3,
		Extent:         4096,
		Buffer:         64,
		Debug:          1,
		LayerName:      "a",
		MaxFeatures:    50000,
		PropertiesBool: true,
	}
}

type Geometry struct {
	Type            string
	Point           []float64
	MultiPoint      []float64
	LineString      []float64
	MultiLineString [][]float64
	Polygon         [][]float64
	MultiPolygon    [][][]float64
}

type Feature struct {
	ID       interface{}
	Type     string
	Geometry Geometry
	Tags     map[string]interface{}
	MinX     float64
	MaxX     float64
	MinY     float64
	MaxY     float64
}

func CreateFeature(id interface{}, geom Geometry, tags map[string]interface{}) Feature {
	feature := Feature{
		ID:       id,
		Geometry: geom,
		Type:     geom.Type,
		Tags:     tags,
	}
	feature.calcBBox()
	return feature
}

func (feature *Feature) calcBBox() {
	switch feature.Type {
	case "Point":
		feature.calcLineBBox(feature.Geometry.Point)
	case "MultiPoint":
		feature.calcLineBBox(feature.Geometry.MultiPoint)
	case "LineString":
		feature.calcLineBBox(feature.Geometry.LineString)
	case "MultiLineString":
		for i := range feature.Geometry.MultiLineString {
			feature.calcLineBBox(feature.Geometry.MultiLineString[i])
		}
	case "Polygon":
		for i := range feature.Geometry.Polygon {
			feature.calcLineBBox(feature.Geometry.Polygon[i])
		}
	case "MultiPolygon":
		for ii := range feature.Geometry.MultiPolygon {
			for i := range feature.Geometry.MultiPolygon[ii] {
				feature.calcLineBBox(feature.Geometry.MultiPolygon[ii][i])
			}
		}
	}
}

func (feature *Feature) calcLineBBox(line []float64) {
	for i := 0; i < len(line); i += 3 {
		feature.MinX = math.Min(feature.MinX, line[i])
		feature.MinY = math.Min(feature.MinY, line[i+1])
		feature.MaxX = math.Max(feature.MaxX, line[i])
		feature.MaxY = math.Max(feature.MaxY, line[i+1])
	}
}

func (tile Tile) SplitTileChildren() map[m.TileID]Tile {
	var minx, maxx, miny, maxy float64
	minx, maxx, miny, maxy = tile.MinX, tile.MaxX, tile.MinY, tile.MaxY

	tileid := tile.TileID

	k1 := 0.5 * float64(tile.Options.Buffer) / float64(tile.Options.Extent)
	k2 := 0.5 - k1
	k3 := 0.5 + k1
	k4 := 1 + k1
	z2 := 1 << tileid.Z

	chanl := make(chan []Feature, 1)
	chanr := make(chan []Feature, 1)

	clipcon := func(x chan []Feature, features []Feature, k1, k2 float64, min, max float64, yint int) {
		x <- clip(features, z2, k1, k2, yint, min, max, tile.Options)
	}
	clipcontile := func(x chan Tile, features []Feature, k1, k2 float64, min, max float64, yint int, tileid m.TileID) {
		x <- clipcreate(features, z2, k1, k2, yint, min, max, tile.Options, tileid)
	}

	tilemap := map[m.TileID]Tile{}
	tl_tile, bl_tile, tr_tile, br_tile := m.TileID{X: tileid.X * 2, Y: tileid.Y * 2, Z: tileid.Z + 1}, m.TileID{X: tileid.X * 2, Y: tileid.Y * 2, Z: tileid.Z + 1}, m.TileID{X: tileid.X * 2, Y: tileid.Y * 2, Z: tileid.Z + 1}, m.TileID{X: tileid.X * 2, Y: tileid.Y * 2, Z: tileid.Z + 1}
	bl_tile.Y = bl_tile.Y + 1
	tr_tile.X = tr_tile.X + 1
	br_tile.X, br_tile.Y = br_tile.X+1, br_tile.Y+1

	go clipcon(chanl, tile.Source, float64(tileid.X)-k1, float64(tileid.X)+k3, minx, maxx, 0)
	go clipcon(chanr, tile.Source, float64(tileid.X)+k2, float64(tileid.X)+k4, minx, maxx, 0)

	var chantr, chanbr, chantl, chanbl chan Tile
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-chanl:
			chantl = make(chan Tile, 1)
			chanbl = make(chan Tile, 1)
			if len(msg1) > 0 {
				go clipcontile(chantl, msg1, float64(tileid.Y)-k1, float64(tileid.Y)+k3, miny, maxy, 1, tl_tile)
				go clipcontile(chanbl, msg1, float64(tileid.Y)+k2, float64(tileid.Y)+k4, miny, maxy, 1, bl_tile)

			} else {
				chantl <- Tile{}
				chanbl <- Tile{}
			}
		case msg2 := <-chanr:
			chantr = make(chan Tile, 1)
			chanbr = make(chan Tile, 1)
			if len(msg2) > 0 {
				go clipcontile(chantr, msg2, float64(tileid.Y)-k1, float64(tileid.Y)+k3, miny, maxy, 1, tr_tile)
				go clipcontile(chanbr, msg2, float64(tileid.Y)+k2, float64(tileid.Y)+k4, miny, maxy, 1, br_tile)
			} else {
				chantr <- Tile{}
				chanbr <- Tile{}
			}
		}
	}

	for i := 0; i < 4; i++ {
		select {
		case msg1 := <-chantl:
			if len(msg1.Source) > 0 {
				tilemap[tl_tile] = msg1
			}
		case msg2 := <-chanbl:
			if len(msg2.Source) > 0 {
				tilemap[bl_tile] = msg2
			}
		case msg3 := <-chantr:
			if len(msg3.Source) > 0 {
				tilemap[tr_tile] = msg3
			}

		case msg4 := <-chanbr:
			if len(msg4.Source) > 0 {
				tilemap[br_tile] = msg4
			}
		}
	}
	return tilemap
}

func TileFromGeoJSON(geojsonfeatures []*geom.Feature, tileid m.TileID, options Config) Tile {
	tile := NewTile()
	tile.TileID = m.Parent(tileid)
	features := ConvertFeatures(geojsonfeatures, options)
	features = wrap(features, options)
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
	return tile.SplitTileChildren()[tileid]
}
