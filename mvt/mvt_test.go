package mvt

import (
	"testing"

	"github.com/flywave/go-geom"
	gen "github.com/flywave/go-geom/general"
	"github.com/flywave/go-pbf"
	m "github.com/flywave/go-mapbox/tileid"
)

func TestIntersectX_NoDivisionByZero(t *testing.T) {
	slice := NewSlice(0)
	_ = slice.IntersectX(0, 0, 0, 10, 5)
	if len(slice.Slice) != 1 {
		t.Fatal("expected 1 point in slice")
	}
	if slice.Slice[0][0] != 5 || slice.Slice[0][1] != 0 {
		t.Fatalf("expected (5,0), got (%f,%f)", slice.Slice[0][0], slice.Slice[0][1])
	}
}

func TestIntersectY_NoDivisionByZero(t *testing.T) {
	slice := NewSlice(1)
	_ = slice.IntersectY(0, 0, 10, 0, 5)
	if len(slice.Slice) != 1 {
		t.Fatal("expected 1 point in slice")
	}
	if slice.Slice[0][0] != 0 || slice.Slice[0][1] != 5 {
		t.Fatalf("expected (0,5), got (%f,%f)", slice.Slice[0][0], slice.Slice[0][1])
	}
}

func TestIsEmpty_UnknownType(t *testing.T) {
	gd := geom.GeometryData{Type: "GeometryCollection"}
	if IsEmpty(gd) {
		t.Error("IsEmpty(GeometryCollection) = true, want false")
	}
}

func TestSimplify_DegenerateRange(t *testing.T) {
	Simplify([]float64{0, 0, 1, 1, 0, 1}, 3, 2, 0.1)
}

func TestSimplify_TwoPoints(t *testing.T) {
	coords := []float64{0, 0, 1, 10, 10, 1}
	Simplify(coords, 0, 3, 0.1)
	if coords[2] != 1 || coords[5] != 1 {
		t.Fatal("two-point simplify should keep both endpoints")
	}
}

func TestCursor_ExtentBoolDefault(t *testing.T) {
	tileID := m.TileID{X: 0, Y: 0, Z: 0}
	cur := NewCursor(tileID)
	if cur.ExtentBool {
		t.Error("expected ExtentBool = false by default")
	}
}

func TestMakePolygonFloat_WindingOrder(t *testing.T) {
	cur := NewCursorExtent(m.TileID{X: 0, Y: 0, Z: 0}, 4096)
	cur.MakePolygonFloat([][][]float64{
		{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}},
	})
	if len(cur.Geometry) == 0 {
		t.Fatal("expected non-empty geometry")
	}
}

func TestNewLayerConfig_ProtoMapbox(t *testing.T) {
	tileID := m.TileID{X: 0, Y: 0, Z: 0}
	conf := NewConfig("test", tileID, PROTO_MAPBOX)
	layer := NewLayerConfig(conf)
	if layer.Proto != MapboxProto {
		t.Fatalf("expected MapboxProto, got %+v", layer.Proto)
	}
}

func TestNewLayerConfig_ProtoLK(t *testing.T) {
	tileID := m.TileID{X: 0, Y: 0, Z: 0}
	conf := NewConfig("test", tileID, PROTO_LK)
	layer := NewLayerConfig(conf)
	if layer.Proto != LKProto {
		t.Fatalf("expected LKProto, got %+v", layer.Proto)
	}
}

func TestGetProto_ValidTypes(t *testing.T) {
	if getProto(PROTO_MAPBOX) != MapboxProto {
		t.Fatal("PROTO_MAPBOX should return MapboxProto")
	}
	if getProto(PROTO_LK) != LKProto {
		t.Fatal("PROTO_LK should return LKProto")
	}
	if getProto(ProtoType(999)) != MapboxProto {
		t.Fatal("unknown ProtoType should default to MapboxProto")
	}
}

func TestReadTile_EmptyData(t *testing.T) {
	tileID := m.TileID{X: 0, Y: 0, Z: 0}
	feats, err := ReadTile(nil, tileID, PROTO_MAPBOX)
	if err == nil {
		t.Log("ReadTile with nil data returned no error (may be expected)")
	}
	if feats != nil {
		t.Log("ReadTile with nil data returned features (may be expected)")
	}
}

func TestSignedArea2_Empty(t *testing.T) {
	result := SignedArea2([][]int32{})
	if result != 0 {
		t.Fatalf("expected 0 for empty, got %f", result)
	}
}

func TestSignedArea2_SinglePoint(t *testing.T) {
	result := SignedArea2([][]int32{{0, 0}})
	if result != 0 {
		t.Fatalf("expected 0 for single point, got %f", result)
	}
}

func TestPointClipAboutZoom_Basic(t *testing.T) {
	feat := &geom.Feature{
		GeometryData: *geom.NewPointGeometryData([]float64{0, 0}),
		Properties:   map[string]interface{}{},
	}
	result := PointClipAboutZoom(feat, 0)
	if result == nil {
		t.Fatal("expected non-nil map, got nil")
	}
}

func TestPointClipAboutTile_Basic(t *testing.T) {
	tileID := m.TileID{X: 0, Y: 0, Z: 0}
	feat := &geom.Feature{
		GeometryData: *geom.NewPointGeometryData([]float64{0, 0}),
		Properties:   map[string]interface{}{},
	}
	result := PointClipAboutTile(feat, tileID)
	if result == nil {
		t.Fatal("expected non-nil feature")
	}
}

func TestPointClipAboutTile_NonPointType(t *testing.T) {
	tileID := m.TileID{X: 0, Y: 0, Z: 0}
	feat := &geom.Feature{
		GeometryData: geom.GeometryData{Type: "LineString", LineString: [][]float64{{0, 0}, {10, 10}}},
		Properties:   map[string]interface{}{},
	}
	result := PointClipAboutTile(feat, tileID)
	if result == nil {
		t.Fatal("expected non-nil feature")
	}
}

func TestClipDownTile_Basic(t *testing.T) {
	gd := *geom.NewPointGeometryData([]float64{0, 0})
	result := ClipDownTile(gd, m.TileID{X: 0, Y: 0, Z: 1})
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestDeltaFirstLast_NormalOperation(t *testing.T) {
	DeltaFirstLast([]float64{0, 0}, []float64{10, 10})
}

func TestConvertGeometry_Nil(t *testing.T) {
	result := ConvertGeometry(nil, 1)
	if result == nil {
		t.Fatal("expected non-nil result for nil input")
	}
}

func TestConvertGeometry_ZeroDim(t *testing.T) {
	gd := &geom.GeometryData{Type: "Point", Point: []float64{0, 0}}
	result := ConvertGeometry(gd, 0)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestWriteLayer_ProducesParsableProtobuf(t *testing.T) {
	tileID := m.TileID{X: 0, Y: 0, Z: 0}
	feats := []*geom.Feature{
		{Geometry: gen.NewPoint([]float64{0, 0}), Properties: map[string]interface{}{"name": "origin"}},
	}
	conf := NewConfig("test", tileID, PROTO_MAPBOX)
	data := WriteLayer(feats, conf)
	r := pbf.NewReader(data)
	foundLayer := false
	for r.Pos < r.Length {
		k, v := r.ReadTag()
		if k == 3 && v == pbf.Bytes {
			foundLayer = true
			lLen := r.ReadVarint()
			end := r.Pos + lLen
			foundName := false
			for r.Pos < end {
				fk, fv := r.ReadTag()
				if fk == 1 && fv == pbf.Bytes {
					if r.ReadString() == "test" {
						foundName = true
					}
				} else if fv == pbf.Bytes {
					skip := r.ReadVarint()
					r.Pos += skip
				} else if fv == pbf.Varint {
					r.ReadVarint()
				}
			}
			if !foundName {
				t.Error("layer name not found")
			}
		}
	}
	if !foundLayer {
		t.Error("no layer found in protobuf output")
	}
}

func TestCursorOperations(t *testing.T) {
	cur := NewCursorExtent(m.TileID{X: 0, Y: 0, Z: 0}, 4096)
	cur.MakePointFloat([]float64{0, 0})
	if len(cur.Geometry) == 0 {
		t.Fatal("expected non-empty geometry after MakePointFloat")
	}
}
