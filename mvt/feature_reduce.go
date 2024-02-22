package mvt

import (
	m "github.com/flywave/go-mapbox/tileid"

	"github.com/flywave/go-geom"
)

var default_steps = 8
var default_percentage = .005

type Reduce_Config struct {
	Point_Steps int
	Percent     float64
	TileID      m.TileID
	TileMap     map[m.TileID]string
	DeltaX      float64
	DeltaY      float64
	Zoom        int
	OperateBool bool
}

func NewReduceConfig(tileid m.TileID) *Reduce_Config {
	bds := m.Bounds(tileid)
	deltax := bds.E - bds.W
	deltay := bds.N - bds.S
	zoom := int(tileid.Z) + default_steps
	return &Reduce_Config{
		Point_Steps: default_steps,
		Percent:     default_percentage,
		TileID:      tileid,
		TileMap:     map[m.TileID]string{},
		DeltaX:      deltax,
		DeltaY:      deltay,
		Zoom:        zoom,
		OperateBool: true,
	}
}

func Filter(feature *geom.Feature, config *Reduce_Config) bool {
	if !config.OperateBool {
		return true
	}
	switch feature.GeometryData.Type {
	case "Point":

		tile := m.Tile(feature.GeometryData.Point[0], feature.GeometryData.Point[1], config.Zoom)
		_, boolval := config.TileMap[tile]
		if !boolval {
			config.TileMap[tile] = ""
			return true
		} else {
			return false
		}
	case "MultiPoint":
		total_x, total_y := 0.0, 0.0
		for _, point := range feature.GeometryData.MultiPoint {
			total_x += point[0]
			total_y += point[1]
		}
		size := float64(len(feature.GeometryData.MultiPoint))
		avg_x, avg_y := total_x/size, total_y/size
		tile := m.Tile(avg_x, avg_y, config.Zoom)
		_, boolval := config.TileMap[tile]
		if !boolval {
			config.TileMap[tile] = ""
			feature.GeometryData.Type = "Point"
			feature.GeometryData.MultiPoint = [][]float64{}
			feature.GeometryData.Point = []float64{avg_x, avg_y}
			return true
		} else {
			return false
		}
	case "LineString", "MultiLineString", "MultiPolygon", "Polygon":
		bbox := geom.BoundingBoxFromGeometryData(&feature.GeometryData)
		percentx := (bbox[1][0] - bbox[0][0]) / config.DeltaX
		percenty := (bbox[1][1] - bbox[0][1]) / config.DeltaY
		return percentx > config.Percent || percenty > config.Percent
	}
	return false
}
