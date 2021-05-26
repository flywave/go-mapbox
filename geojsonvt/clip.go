package geojsonvt

import (
	m "github.com/flywave/go-mapbox/tileid"
)

type ClipGeom struct {
	Geom      []float64
	NewGeom   [][]float64
	K1        float64
	K2        float64
	Axis      int
	IsPolygon bool
	SlicePos  int
}

type Slice struct {
	Pos   int
	Slice []float64
	Axis  int
}

func NewSlice(linesize int, axis int) *Slice {
	return &Slice{Slice: make([]float64, linesize), Axis: axis}
}

func (slice *Slice) IntersectX(ax, ay, bx, by, x float64) float64 {
	t := (x - ax) / (bx - ax)
	slice.Slice[slice.Pos] = x
	slice.Slice[slice.Pos+1] = ay + (by-ay)*t
	slice.Slice[slice.Pos+2] = 1
	slice.Pos += 3
	return t
}

func (slice *Slice) IntersectY(ax, ay, bx, by, y float64) float64 {
	t := (y - ay) / (by - ay)
	slice.Slice[slice.Pos] = ax + (bx-ax)*t
	slice.Slice[slice.Pos+1] = y
	slice.Slice[slice.Pos+2] = 1
	slice.Pos += 3
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

func (slice *Slice) AddPoint(x, y, z float64) {
	slice.Slice[slice.Pos] = x
	slice.Slice[slice.Pos+1] = y
	slice.Slice[slice.Pos+2] = z
	slice.Pos += 3
}

func (input *ClipGeom) clipLine() {
	slice := NewSlice(len(input.Geom)*3, input.Axis)
	var ax, ay, az, bx, by, a, b float64
	k1, k2 := input.K1, input.K2
	for i := 0; i < len(input.Geom)-3; i += 3 {
		ax = input.Geom[i]
		ay = input.Geom[i+1]
		az = input.Geom[i+2]
		bx = input.Geom[i+3]
		by = input.Geom[i+4]
		if input.Axis == 0 {
			a = ax
			b = bx
		} else if input.Axis == 1 {
			a = ay
			b = by
		}
		exited := false

		if a < k1 {
			if b >= k1 {
				slice.Intersect(ax, ay, bx, by, k1)
			}
		} else if a >= k2 {
			if b < k2 {
				slice.Intersect(ax, ay, bx, by, k2)
			}
		} else {
			slice.AddPoint(ax, ay, az)
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
			input.NewGeom = append(input.NewGeom, slice.Slice[:slice.Pos])
			slice = NewSlice(len(input.Geom)*3, input.Axis)
		}

	}

	last := len(input.Geom) - 3
	ax = input.Geom[last]
	ay = input.Geom[last+1]
	az = input.Geom[last+2]
	if input.Axis == 0 {
		a = ax
	} else if input.Axis == 1 {
		a = ay
	}
	if a >= k1 && a <= k2 {
		slice.AddPoint(ax, ay, az)
	}

	last = len(slice.Slice) - 3
	if input.IsPolygon && last >= 3 && (slice.Slice[last] != slice.Slice[0] || slice.Slice[last+1] != slice.Slice[1]) {
		slice.AddPoint(slice.Slice[0], slice.Slice[1], slice.Slice[2])
	}

	if slice.Pos > 0 {
		input.NewGeom = append(input.NewGeom, slice.Slice[:slice.Pos])
	}
}

func clipPoints(geometry []float64, k1, k2 float64, axis int) []float64 {
	slice := NewSlice(len(geometry), axis)
	for i := 0; i < len(geometry); i += 3 {
		a := geometry[i+axis]

		if a >= k2 && a <= k2 {
			slice.AddPoint(geometry[i], geometry[i+1], geometry[i+2])
		}
	}
	return slice.Slice[:slice.Pos]
}

func clipLines(geom [][]float64, k1, k2 float64, axis int, IsPolygon bool) *ClipGeom {
	clipthing := &ClipGeom{Geom: geom[0], K1: k1, K2: k2, Axis: axis, IsPolygon: IsPolygon}
	for pos := range geom {
		clipthing.Geom = geom[pos]
		clipthing.clipLine()

	}
	return clipthing
}

func (geometry *Geometry) Clip(k1 float64, k2 float64, axis int, ispolygon bool) (Geometry, bool) {
	switch geometry.Type {
	case "Point":
		geom := clipPoints(geometry.Point, k1, k2, axis)
		if len(geom) > 0 {
			return Geometry{Type: "Point", Point: geom}, true
		}
	case "MultiPoint":
		geom := clipPoints(geometry.MultiPoint, k1, k2, axis)
		if len(geom) > 0 {
			if len(geom) == 3 {
				return Geometry{Type: "Point", Point: geom}, true
			} else {
				return Geometry{Type: "MultiPoint", MultiPoint: geom}, true
			}
		}
	case "LineString":
		clipgeom := ClipGeom{K1: k1, K2: k2, Geom: geometry.LineString, Axis: axis, IsPolygon: false}
		clipgeom.clipLine()
		if len(clipgeom.NewGeom) > 0 {
			if len(clipgeom.NewGeom) == 1 {
				return Geometry{Type: "LineString", LineString: clipgeom.NewGeom[0]}, true
			} else {
				return Geometry{Type: "MultiLineString", MultiLineString: clipgeom.NewGeom}, true
			}
		}
	case "MultiLineString":
		clipgeom := clipLines(geometry.MultiLineString, k1, k2, axis, false)
		if len(clipgeom.NewGeom) > 0 {
			if len(clipgeom.NewGeom) == 1 {
				return Geometry{Type: "LineString", LineString: clipgeom.NewGeom[0]}, true
			} else {
				return Geometry{Type: "MultiLineString", MultiLineString: clipgeom.NewGeom}, true
			}
		}
	case "Polygon":
		clipgeom := clipLines(geometry.Polygon, k1, k2, axis, ispolygon)
		if len(clipgeom.NewGeom) > 0 {
			return Geometry{Type: "Polygon", Polygon: clipgeom.NewGeom}, true
		}
	case "MultiPolygon":
		multipolygon := [][][]float64{}
		for i := range geometry.MultiPolygon {
			clipgeom := clipLines(geometry.MultiPolygon[i], k1, k2, axis, ispolygon)
			if len(clipgeom.NewGeom) > 0 {
				multipolygon = append(multipolygon, clipgeom.NewGeom)
			}
		}
		if len(multipolygon) > 0 {
			if len(multipolygon) == 1 {
				return Geometry{Type: "Polygon", Polygon: multipolygon[0]}, true
			} else {
				return Geometry{Type: "MultiPolygon", MultiPolygon: multipolygon}, true
			}
		}

	}
	return Geometry{}, false
}

func clip(features []Feature, scale int, k1 float64, k2 float64, axis int, minAll float64, maxAll float64, options Config) []Feature {
	k1 = k1 / float64(scale)
	k2 = k2 / float64(scale)
	if minAll >= k1 && maxAll < k2 {
		return features
	} else if maxAll < k1 || minAll >= k2 {
		return []Feature{}
	}

	clipped := []Feature{}
	var min, max float64
	for _, feature := range features {
		if axis == 0 {
			min = feature.MinX
			max = feature.MaxX
		} else if axis == 1 {
			min = feature.MinY
			max = feature.MaxY

		}
		boolval := true
		if min >= k1 && max < k2 {
			//boolval = false
		} else if max < k1 || min >= k2 {
			//boolval = false
		}

		if boolval {
			clipgeom, _ := feature.Geometry.Clip(k1, k2, axis, feature.Type == "MultiPolygon" || feature.Type == "Polygon")
			boolval2 := (len(clipgeom.LineString) == 0) && (len(clipgeom.Point) == 0) && (len(clipgeom.Polygon) == 0) && (len(clipgeom.MultiLineString) == 0) && (len(clipgeom.MultiPoint) == 0) && (len(clipgeom.MultiPolygon) == 0)
			if !boolval2 {
				clipped = append(clipped, CreateFeature(feature.ID, clipgeom, feature.Tags))
			}
		}
	}
	return clipped
}

func clipcreate(features []Feature, scale int, k1 float64, k2 float64, axis int, minAll float64, maxAll float64, options Config, tileid m.TileID) Tile {
	tile := NewTile()
	tile.TileID = tileid

	k1 = k1 / float64(scale)
	k2 = k2 / float64(scale)
	if minAll >= k1 && maxAll < k2 {
		return CreateTile(features, tileid, options)
	} else if maxAll < k1 || minAll >= k2 {
		return Tile{}
	}

	clipped := []Feature{}
	var min, max float64
	for _, feature := range features {
		if axis == 0 {
			min = feature.MinX
			max = feature.MaxX
		} else if axis == 1 {
			min = feature.MinY
			max = feature.MaxY

		}
		boolval := true
		if min >= k1 && max < k2 {
			//boolval = false
		} else if max < k1 || min >= k2 {
			//boolval = false
		}

		if boolval {
			clipgeom, _ := feature.Geometry.Clip(k1, k2, axis, feature.Type == "MultiPolygon" || feature.Type == "Polygon")
			boolval2 := (len(clipgeom.LineString) == 0) && (len(clipgeom.Point) == 0) && (len(clipgeom.Polygon) == 0) && (len(clipgeom.MultiLineString) == 0) && (len(clipgeom.MultiPoint) == 0) && (len(clipgeom.MultiPolygon) == 0)
			if !boolval2 {
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
				clipped = append(clipped, CreateFeature(feature.ID, clipgeom, feature.Tags))
			}
		}
	}
	tile.Source = clipped
	tile.Options = options
	return tile
}
