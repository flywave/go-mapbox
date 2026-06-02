package mvt

import (
	"bytes"
	"os"
	"testing"

	"github.com/flywave/go-geom"
	gen "github.com/flywave/go-geom/general"
	"github.com/flywave/go-pbf"
	m "github.com/flywave/go-mapbox/tileid"
)

// === BUG #2: IsEmpty incorrectly reports Point as empty ===

func TestIsEmpty_PointWithCoords(t *testing.T) {
	gd := geom.NewPointGeometryData([]float64{100, 200})
	if IsEmpty(*gd) {
		t.Error("IsEmpty(Point{100,200}) = true, want false")
	}
}

func TestIsEmpty_PointEmpty(t *testing.T) {
	gd := geom.NewPointGeometryData(nil)
	if !IsEmpty(*gd) {
		t.Error("IsEmpty(Point{nil}) = false, want true")
	}
}

func TestIsEmpty_LineString(t *testing.T) {
	gd := geom.NewLineStringGeometryData([][]float64{{0, 0}, {1, 1}})
	if IsEmpty(*gd) {
		t.Error("IsEmpty(LineString{{0,0},{1,1}}) = true, want false")
	}
}

func TestIsEmpty_EmptyLineString(t *testing.T) {
	gd := geom.NewLineStringGeometryData(nil)
	if !IsEmpty(*gd) {
		t.Error("IsEmpty(LineString{nil}) = false, want true")
	}
}

func TestIsEmpty_Polygon(t *testing.T) {
	gd := geom.NewPolygonGeometryData([][][]float64{{{0, 0}, {1, 0}, {1, 1}, {0, 0}}})
	if IsEmpty(*gd) {
		t.Error("IsEmpty(Polygon{...}) = true, want false")
	}
}

// === BUG #3: PointClipAboutTile single-point branch ===

func TestPointClipAboutTile_SinglePoint(t *testing.T) {
	tileid := m.TileID{X: 0, Y: 0, Z: 0}

	// A single point in MultiPoint that falls in tile (0,0,0)
	feat := &geom.Feature{
		GeometryData: *geom.NewMultiPointGeometryData(
			[]float64{0, 0},
		),
		Properties: map[string]interface{}{},
	}

	result := PointClipAboutTile(feat, tileid)
	if len(result.Properties) == 0 && result.GeometryData.Type == "" {
		t.Fatal("PointClipAboutTile returned empty feature for point in tile")
	}
	if result.GeometryData.Type != "Point" {
		t.Errorf("expected single point to produce Type=Point, got %s", result.GeometryData.Type)
	}
}

func TestPointClipAboutTile_MultiPoint(t *testing.T) {
	tileid := m.TileID{X: 0, Y: 0, Z: 0}

	feat := &geom.Feature{
		GeometryData: *geom.NewMultiPointGeometryData(
			[]float64{0, 0},
			[]float64{1, 1},
		),
		Properties: map[string]interface{}{},
	}

	result := PointClipAboutTile(feat, tileid)
	if len(result.Properties) == 0 && result.GeometryData.Type == "" {
		t.Fatal("PointClipAboutTile returned empty feature for points in tile")
	}
	if result.GeometryData.Type != "MultiPoint" {
		t.Errorf("expected multi-point to produce Type=MultiPoint, got %s", result.GeometryData.Type)
	}
}

func TestPointClipAboutTile_PointOutOfTile(t *testing.T) {
	tileid := m.TileID{X: 0, Y: 0, Z: 0}

	feat := &geom.Feature{
		GeometryData: *geom.NewMultiPointGeometryData(
			[]float64{200, 200},
		),
		Properties: map[string]interface{}{},
	}

	result := PointClipAboutTile(feat, tileid)
	if result.GeometryData.Type != "" {
		t.Errorf("expected empty feature for point outside tile, got Type=%s", result.GeometryData.Type)
	}
}

// === BUG #3: PointClipAboutZoom single-point branch ===

func TestPointClipAboutZoom_SinglePoint(t *testing.T) {
	feat := &geom.Feature{
		GeometryData: *geom.NewMultiPointGeometryData(
			[]float64{0, 0},
		),
		Properties: map[string]interface{}{},
	}

	result := PointClipAboutZoom(feat, 0)
	if len(result) == 0 {
		t.Fatal("PointClipAboutZoom returned empty map")
	}
	for tid, f := range result {
		if f.GeometryData.Type != "Point" {
			t.Errorf("expected single point to produce Type=Point at tile %v, got %s", tid, f.GeometryData.Type)
		}
	}
}

func TestPointClipAboutZoom_MultiPoint(t *testing.T) {
	feat := &geom.Feature{
		GeometryData: *geom.NewMultiPointGeometryData(
			[]float64{0, 0},
			[]float64{1, 1},
		),
		Properties: map[string]interface{}{},
	}

	result := PointClipAboutZoom(feat, 0)
	if len(result) == 0 {
		t.Fatal("PointClipAboutZoom returned empty map")
	}
	for tid, f := range result {
		if f.GeometryData.Type != "MultiPoint" {
			t.Errorf("expected multi-point to produce Type=MultiPoint at tile %v, got %s", tid, f.GeometryData.Type)
		}
	}
}

// === BUG #4: Verify write-layer output structure (protobuf level) ===

func TestWriteLayer_ProducesValidProtobuf(t *testing.T) {
	tileid := m.TileID{X: 0, Y: 0, Z: 0}

	feats := []*geom.Feature{
		{
			Geometry:   gen.NewPoint([]float64{0, 0}),
			Properties: map[string]interface{}{"name": "origin"},
		},
	}

	conf := NewConfig("test", tileid, PROTO_MAPBOX)
	conf.ExtentBool = true
	data := WriteLayer(feats, conf)
	if len(data) == 0 {
		t.Fatal("WriteLayer returned empty data")
	}

	// Verify protobuf structure at the tag level
	r := pbf.NewReader(data)
	var foundLayer bool
	var layerVersion, layerExtent int
	var featureCount int

	for r.Pos < r.Length {
		k, v := r.ReadTag()
		if k == 3 && v == pbf.Bytes { // tile.layers
			foundLayer = true
			lLen := r.ReadVarint()
			lEnd := r.Pos + lLen

			for r.Pos < lEnd {
				fk, fv := r.ReadTag()
				switch {
				case fk == 1 && fv == pbf.Bytes: // layer.name
					r.ReadString()
				case fk == 2 && fv == pbf.Bytes: // layer.features
					fLen := r.ReadVarint()
					r.Pos += fLen
				case fk == 5 && fv == pbf.Varint: // layer.extent
					layerExtent = r.ReadVarint()
				case fk == 15 && fv == pbf.Varint: // layer.version
					layerVersion = r.ReadVarint()
				default:
					if fv == pbf.Bytes {
						skip := r.ReadVarint()
						r.Pos += skip
					} else if fv == pbf.Varint {
						r.ReadVarint()
					}
				}
			}
		}
	}

	if !foundLayer {
		t.Error("no layer found in tile")
	}
	if layerVersion != 2 {
		t.Errorf("expected version 2, got %d", layerVersion)
	}
	if layerExtent != 4096 {
		t.Errorf("expected extent 4096, got %d", layerExtent)
	}

	// Count features at the raw bytes level
	r2 := pbf.NewReader(data)
	for r2.Pos < r2.Length {
		k, v := r2.ReadTag()
		if k == 3 && v == pbf.Bytes {
			lLen := r2.ReadVarint()
			lEnd := r2.Pos + lLen
			for r2.Pos < lEnd {
				fk, fv := r2.ReadTag()
				if fk == 2 && fv == pbf.Bytes {
					featureCount++
					fLen := r2.ReadVarint()
					r2.Pos += fLen
				} else if fv == pbf.Bytes {
					skip := r2.ReadVarint()
					r2.Pos += skip
				} else if fv == pbf.Varint {
					r2.ReadVarint()
				}
			}
		}
	}
	if featureCount != 1 {
		t.Errorf("expected 1 feature, found %d", featureCount)
	}
}

func TestWriteLayer_MultiFeature(t *testing.T) {
	tileid := m.TileID{X: 0, Y: 0, Z: 0}

	// The +1 off-by-one bug in ReadRawTile would corrupt subsequent features.
	// Verify that 5 features round-trip correctly at the wire level.
	n := 5
	feats := make([]*geom.Feature, n)
	for i := range feats {
		feats[i] = &geom.Feature{
			Geometry:   gen.NewPoint([]float64{float64(i), float64(i)}),
			Properties: map[string]interface{}{"idx": i},
		}
	}

	conf := NewConfig("multi", tileid, PROTO_MAPBOX)
	conf.ExtentBool = true
	data := WriteLayer(feats, conf)

	r := pbf.NewReader(data)
	var count int
	for r.Pos < r.Length {
		k, v := r.ReadTag()
		if k == 3 && v == pbf.Bytes {
			lLen := r.ReadVarint()
			lEnd := r.Pos + lLen
			for r.Pos < lEnd {
				fk, fv := r.ReadTag()
				if fk == 2 && fv == pbf.Bytes {
					count++
					fLen := r.ReadVarint()
					r.Pos += fLen
				} else if fv == pbf.Bytes {
					skip := r.ReadVarint()
					r.Pos += skip
				} else if fv == pbf.Varint {
					r.ReadVarint()
				}
			}
		}
	}
	if count != n {
		t.Errorf("expected %d features on wire, got %d", n, count)
	}
}

// === BUG #1: assert_winding_order no debug stderr ===

func TestMakePolygonFloat_NoStderrOutput(t *testing.T) {
	// Capture stderr
	stderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stderr = w

	cur := NewCursorExtent(m.TileID{X: 0, Y: 0, Z: 0}, 4096)
	cur.MakePolygonFloat([][][]float64{
		{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}},
	})

	w.Close()
	os.Stderr = stderr

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if output != "" {
		t.Errorf("MakePolygonFloat wrote to stderr: %q", output)
	}
}

// === BUG #1: MakeMultiPolygonFloat also no debug stderr ===

func TestMakeMultiPolygonFloat_NoStderrOutput(t *testing.T) {
	stderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stderr = w

	cur := NewCursorExtent(m.TileID{X: 0, Y: 0, Z: 0}, 4096)
	cur.MakeMultiPolygonFloat([][][][]float64{
		{{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}},
	})

	w.Close()
	os.Stderr = stderr

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if output != "" {
		t.Errorf("MakeMultiPolygonFloat wrote to stderr: %q", output)
	}
}
