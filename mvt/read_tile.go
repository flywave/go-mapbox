package mvt

import (
	"errors"
	"math"

	m "github.com/flywave/go-mapbox/tileid"

	"github.com/flywave/go-geom"
	"github.com/flywave/go-pbf"
)

type Tile struct {
	LayerMap map[string]*Layer
	Buf      *pbf.Reader
	TileID   m.TileID
	Layers   []string
	Proto    Proto
}

func NewTile(bytevals []byte, pt ProtoType) (tile *Tile, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("Error in NewTile.")
		}
	}()

	proto := getProto(pt)

	tile = &Tile{
		LayerMap: map[string]*Layer{},
		Buf:      &pbf.Reader{Pbf: bytevals, Length: len(bytevals)},
		Proto:    proto,
	}
	for tile.Buf.Pos < tile.Buf.Length {
		key, val := tile.Buf.ReadTag()
		if key == proto.Layers && val == pbf.Bytes {
			size := tile.Buf.ReadVarint()
			if size != 0 {
				tile.NewLayer(tile.Buf.Pos+size, pt)
			}

		}
	}
	return tile, err
}

func (tile *Tile) Render() []byte {
	totalbs := []byte{}
	for _, v := range tile.LayerMap {
		bs := v.Buf.Pbf[v.StartPos:v.EndPos]
		prefix := pbf.EncodeVarint(uint64(len(bs)))
		layerbs := append(append([]byte{26}, prefix...), bs...)
		totalbs = append(totalbs, layerbs...)
	}
	return totalbs
}

func ReadTile(bytevals []byte, tileid m.TileID, pt ProtoType) (totalfeautures []*geom.Feature, err error) {

	defer func() {
		if recover() != nil {
			err = errors.New("Error in ReadTile")
		}
	}()

	proto := getProto(pt)

	tile := &Tile{
		Buf:    pbf.NewReader(bytevals),
		TileID: tileid,
	}
	totalfeautures = []*geom.Feature{}
	for tile.Buf.Pos < tile.Buf.Length {
		key, val := tile.Buf.ReadTag()
		if key == proto.Layers && val == 2 {
			sizex := tile.Buf.ReadVarint()
			endpos := tile.Buf.Pos + sizex
			var extent, number_features int
			var layername string
			var features []int
			var keys []string
			var values []interface{}
			if sizex != 0 {
				key, val := tile.Buf.ReadTag()
				for tile.Buf.Pos < endpos {
					if key == proto.Layer.Name && val == pbf.Bytes {
						layername = tile.Buf.ReadString()
						key, val = tile.Buf.ReadTag()
					}
					for key == proto.Layer.Features && val == pbf.Bytes {
						features = append(features, tile.Buf.Pos)
						feat_size := tile.Buf.ReadVarint()
						tile.Buf.Pos += feat_size
						key, val = tile.Buf.ReadTag()
					}
					for key == proto.Layer.Keys && val == pbf.Bytes {
						keys = append(keys, tile.Buf.ReadString())
						key, val = tile.Buf.ReadTag()
					}
					for key == proto.Layer.Values && val == pbf.Bytes {
						tile.Buf.ReadVarint()
						newkey, _ := tile.Buf.ReadTag()
						switch newkey {
						case proto.Value.StringValue:
							values = append(values, tile.Buf.ReadString())
						case proto.Value.FloatValue:
							values = append(values, tile.Buf.ReadFloat())
						case proto.Value.DoubleValue:
							values = append(values, tile.Buf.ReadDouble())
						case proto.Value.IntValue:
							values = append(values, tile.Buf.ReadInt64())
						case proto.Value.UIntValue:
							values = append(values, tile.Buf.ReadUInt64())
						case proto.Value.SIntValue:
							values = append(values, tile.Buf.ReadUInt64())
						case proto.Value.BoolIntValue:
							values = append(values, tile.Buf.ReadBool())
						}
						key, val = tile.Buf.ReadTag()
					}
					if key == proto.Layer.Extent && val == pbf.Varint {
						extent = int(tile.Buf.ReadVarint())
						key, val = tile.Buf.ReadTag()
					}
					if key == proto.Layer.Version && val == pbf.Varint {
						_ = int(tile.Buf.ReadVarint())
						key, val = tile.Buf.ReadTag()
					}
				}
				if extent == 0 {
					extent = 4096
				}
				number_features = len(features)
				tile.Buf.Pos = endpos
			}
			feats := make([]*geom.Feature, number_features)
			size := float64(extent) * float64(math.Pow(2, float64(tile.TileID.Z)))
			x0 := float64(extent) * float64(tile.TileID.X)
			y0 := float64(extent) * float64(tile.TileID.Y)
			var feature_geometry, id, geom_type int
			if extent == 0 {
				extent = 4096
			}
			for i, pos := range features {
				tile.Buf.Pos = pos
				endpos := tile.Buf.Pos + tile.Buf.ReadVarint()

				feature := &geom.Feature{Properties: map[string]interface{}{}}

				for tile.Buf.Pos < endpos {
					key, val := tile.Buf.ReadTag()

					if key == proto.Feature.ID && val == pbf.Varint {
						id = int(tile.Buf.ReadUInt64())
					}

					if key == proto.Feature.Tags && val == pbf.Bytes {
						tags := tile.Buf.ReadPackedUInt32()
						i := 0
						for i < len(tags) {
							var key string
							if len(keys) <= int(tags[i]) {
								key = ""
							} else {
								key = keys[tags[i]]
							}
							var val interface{}
							if len(values) <= int(tags[i+1]) {
								val = ""
							} else {
								val = values[tags[i+1]]
							}
							feature.Properties[key] = val
							i += 2
						}
					}
					if key == proto.Feature.Type && val == pbf.Varint {
						geom_type = int(tile.Buf.Varint()[0])
					}
					if key == proto.Feature.Geometry && val == pbf.Bytes {
						feature_geometry = tile.Buf.Pos
						size := tile.Buf.ReadVarint()
						tile.Buf.Pos += size
					}
				}

				tile.Buf.Pos = feature_geometry
				geom_ := tile.Buf.ReadPackedUInt32()
				pos := 0
				var lines [][][]float64
				var polygons [][][][]float64
				var firstpt []float64
				for pos < len(geom_) {
					if geom_[pos] == 9 {
						pos += 1
						if pos != 1 && geom_type == 2 {
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
				if geom_type == 3 {
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

				for i := range polygons {
					for j := range polygons[i] {
						polygons[i][j] = Project(polygons[i][j], x0, y0, size)
					}
				}

				switch geom_type {
				case GeomTypePoint:
					if len(polygons[0][0]) == 1 {
						feature.GeometryData = *geom.NewPointGeometryData(polygons[0][0][0])
					} else {
						feature.GeometryData = *geom.NewMultiPointGeometryData(polygons[0][0]...)
					}
				case GeomTypeLineString:
					if len(polygons[0]) == 1 {
						feature.GeometryData = *geom.NewLineStringGeometryData(polygons[0][0])
					} else {
						feature.GeometryData = *geom.NewMultiLineStringGeometryData(polygons[0]...)
					}
				case GeomTypePolygon:
					if len(polygons) == 1 {
						feature.GeometryData = *geom.NewPolygonGeometryData(polygons[0])
					} else {
						feature.GeometryData = *geom.NewMultiPolygonGeometryData(polygons...)
					}
				}
				feature.GeometryData.EPSG = 4326
				if id != 0 {
					feature.ID = id
				}
				feature.Properties[`layer`] = layername
				feats[i] = feature
			}

			totalfeautures = append(totalfeautures, feats...)
			tile.Buf.Pos = endpos

		}
	}
	if len(totalfeautures) == 0 {
		err = errors.New("No features read from given tile.")
	}
	return totalfeautures, err
}

func ReadRawTile(bytevals []byte, tileId m.TileID, pt ProtoType) ([]*geom.Feature, [][]float64, error) {
	var err error
	defer func() {
		if recover() != nil {
			err = errors.New("Error in ReadTile")
		}
	}()

	proto := getProto(pt)

	tile := &Tile{Buf: pbf.NewReader(bytevals), TileID: tileId}
	totalfeautures := []*geom.Feature{}
	var extent int
	for tile.Buf.Pos < tile.Buf.Length {
		key, val := tile.Buf.ReadTag()
		if key == proto.Layers && val == 2 {
			sizex := tile.Buf.ReadVarint()
			endpos := tile.Buf.Pos + sizex
			var number_features int
			var layername string
			var features []int
			var keys []string
			var values []interface{}
			if sizex != 0 {
				key, val := tile.Buf.ReadTag()
				for tile.Buf.Pos < endpos {
					if key == proto.Layer.Name && val == pbf.Bytes {
						layername = tile.Buf.ReadString()
						key, val = tile.Buf.ReadTag()
					}
					for key == proto.Layer.Features && val == pbf.Bytes {
						features = append(features, tile.Buf.Pos)
						feat_size := tile.Buf.ReadVarint()
						tile.Buf.Pos += feat_size
						key, val = tile.Buf.ReadTag()
					}
					for key == proto.Layer.Keys && val == pbf.Bytes {
						keys = append(keys, tile.Buf.ReadString())
						key, val = tile.Buf.ReadTag()
					}
					for key == proto.Layer.Values && val == pbf.Bytes {
						tile.Buf.ReadVarint()
						newkey, _ := tile.Buf.ReadTag()
						switch newkey {
						case proto.Value.StringValue:
							values = append(values, tile.Buf.ReadString())
						case proto.Value.FloatValue:
							values = append(values, tile.Buf.ReadFloat())
						case proto.Value.DoubleValue:
							values = append(values, tile.Buf.ReadDouble())
						case proto.Value.IntValue:
							values = append(values, tile.Buf.ReadInt64())
						case proto.Value.UIntValue:
							values = append(values, tile.Buf.ReadUInt64())
						case proto.Value.SIntValue:
							values = append(values, tile.Buf.ReadUInt64())
						case proto.Value.BoolIntValue:
							values = append(values, tile.Buf.ReadBool())
						}
						key, val = tile.Buf.ReadTag()
					}
					if key == proto.Layer.Extent && val == pbf.Varint {
						extent = int(tile.Buf.ReadVarint())
						key, val = tile.Buf.ReadTag()
					}
					if key == proto.Layer.Version && val == pbf.Varint {
						_ = int(tile.Buf.ReadVarint())
						key, val = tile.Buf.ReadTag()
					}
				}
				if extent == 0 {
					extent = 4096
				}
				number_features = len(features)
				tile.Buf.Pos = endpos
			}
			feats := make([]*geom.Feature, number_features)
			var feature_geometry, id, geom_type int
			if extent == 0 {
				extent = 4096
			}
			for i, pos := range features {
				tile.Buf.Pos = pos
				endpos := tile.Buf.Pos + tile.Buf.ReadVarint()

				feature := &geom.Feature{Properties: map[string]interface{}{}}

				for tile.Buf.Pos < endpos {
					key, val := tile.Buf.ReadTag()

					if key == proto.Feature.ID && val == pbf.Varint {
						id = int(tile.Buf.ReadUInt64())
					}

					if key == proto.Feature.Tags && val == pbf.Bytes {
						tags := tile.Buf.ReadPackedUInt32()
						i := 0
						for i < len(tags) {
							var key string
							if len(keys) <= int(tags[i]) {
								key = ""
							} else {
								key = keys[tags[i]]
							}
							var val interface{}
							if len(values) <= int(tags[i+1]) {
								val = ""
							} else {
								val = values[tags[i+1]]
							}
							feature.Properties[key] = val
							i += 2
						}
					}
					if key == proto.Feature.Type && val == pbf.Varint {
						geom_type = int(tile.Buf.Varint()[0])
					}
					if key == proto.Feature.Geometry && val == pbf.Bytes {
						feature_geometry = tile.Buf.Pos
						size := tile.Buf.ReadVarint()
						tile.Buf.Pos += size + 1
					}
				}

				tile.Buf.Pos = feature_geometry
				geom_ := tile.Buf.ReadPackedUInt32()
				pos := 0
				var lines [][][]float64
				var polygons [][][][]float64
				var firstpt []float64
				for pos < len(geom_) {
					if geom_[pos] == 9 {
						pos += 1
						if pos != 1 && geom_type == 2 {
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
				if geom_type == 3 {
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

				switch geom_type {
				case GeomTypePoint:
					if len(polygons[0][0]) == 1 {
						feature.GeometryData = *geom.NewPointGeometryData(polygons[0][0][0])
					} else {
						feature.GeometryData = *geom.NewMultiPointGeometryData(polygons[0][0]...)
					}
				case GeomTypeLineString:
					if len(polygons[0]) == 1 {
						feature.GeometryData = *geom.NewLineStringGeometryData(polygons[0][0])
					} else {
						feature.GeometryData = *geom.NewMultiLineStringGeometryData(polygons[0]...)
					}
				case GeomTypePolygon:
					if len(polygons) == 1 {
						feature.GeometryData = *geom.NewPolygonGeometryData(polygons[0])
					} else {
						feature.GeometryData = *geom.NewMultiPolygonGeometryData(polygons...)
					}
				}
				if id != 0 {
					feature.ID = id
				}
				feature.Properties[`layer`] = layername
				feats[i] = feature
			}

			totalfeautures = append(totalfeautures, feats...)
			tile.Buf.Pos = endpos

		}
	}

	if extent == 0 {
		extent = 4096
	}
	size := float64(extent) * float64(math.Pow(2, float64(tile.TileID.Z)))
	x0 := float64(extent) * float64(tile.TileID.X)
	y0 := float64(extent) * float64(tile.TileID.Y)
	pts := [][]float64{{0, 0, 0}, {float64(extent), float64(extent), 0}}
	Project(pts, x0, y0, size)
	if len(totalfeautures) == 0 {
		err = errors.New("No features read from given tile.")
	}
	return totalfeautures, pts, err
}
