package mbtiles

import (
	"testing"
)

// ─── TileFormat.String ─────────────────────────────────────────────────────

func TestTileFormatString(t *testing.T) {
	tests := []struct {
		fmt  TileFormat
		want string
	}{
		{UNKNOWN, ""},
		{GZIP, ""},   // not in switch
		{ZLIB, ""},   // not in switch
		{PNG, "png"},
		{JPG, "jpg"},
		{PBF, "pbf"},
		{WEBP, "webp"},
		{LERC, "lerc"},
		{TIFF, "tiff"},
		{TERRAIN, "terrain"},
		{TIF, "tif"},
	}
	for _, tc := range tests {
		if got := tc.fmt.String(); got != tc.want {
			t.Errorf("TileFormat(%d).String() = %q, want %q", tc.fmt, got, tc.want)
		}
	}
}

// ─── TileFormat.ContentType ────────────────────────────────────────────────

func TestTileFormatContentType(t *testing.T) {
	tests := []struct {
		fmt  TileFormat
		want string
	}{
		{UNKNOWN, ""},
		{PNG, "image/png"},
		{JPG, "image/jpeg"},
		{PBF, "application/x-protobuf"},
		{WEBP, "image/webp"},
		{LERC, "image/lerc"},
		{TIFF, "image/tiff"},
		{TIF, "image/tif"},
		{TERRAIN, "application/vnd.quantized-mesh"},
	}
	for _, tc := range tests {
		if got := tc.fmt.ContentType(); got != tc.want {
			t.Errorf("TileFormat(%d).ContentType() = %q, want %q", tc.fmt, got, tc.want)
		}
	}
}

// ─── formatStrings round-trip ──────────────────────────────────────────────

func TestTileFormatStringsRoundTrip(t *testing.T) {
	for i, s := range formatStrings {
		if s == "" {
			continue
		}
		parsed := StringToTileFormat(s)
		if parsed != TileFormat(i) {
			t.Errorf("StringToTileFormat(%q) = %d, want %d", s, parsed, i)
		}
		back := TileFormatToString(parsed)
		if back != s {
			t.Errorf("TileFormatToString(%d) = %q, want %q", parsed, back, s)
		}
	}
}

func TestStringToTileFormat_Unknown(t *testing.T) {
	if got := StringToTileFormat("nonexistent"); got != UNKNOWN {
		t.Fatalf("got %d", got)
	}
}

// ─── LayerType ─────────────────────────────────────────────────────────────

func TestLayerTypeString(t *testing.T) {
	if BaseLayer.String() != "baselayer" {
		t.Fatalf("BaseLayer.String() = %q", BaseLayer.String())
	}
	if Overlay.String() != "overlay" {
		t.Fatalf("Overlay.String() = %q", Overlay.String())
	}
}

func TestLayerTypeMarshalJSON(t *testing.T) {
	lt := Overlay
	b, err := lt.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `"overlay"` {
		t.Fatalf("got %s", string(b))
	}
}

func TestLayerTypeUnmarshalJSON(t *testing.T) {
	var lt LayerType
	if err := lt.UnmarshalJSON([]byte(`"baselayer"`)); err != nil {
		t.Fatal(err)
	}
	if lt != BaseLayer {
		t.Fatal("expected BaseLayer")
	}
}

func TestLayerTypeUnmarshalJSON_Invalid(t *testing.T) {
	var lt LayerType
	if err := lt.UnmarshalJSON([]byte(`"invalid"`)); err == nil {
		t.Fatal("expected error")
	}
}

func TestLayerTypeHelpers(t *testing.T) {
	if stringToLayerType("baselayer") != BaseLayer {
		t.Fatal("stringToLayerType baselayer")
	}
	if stringToLayerType("overlay") != Overlay {
		t.Fatal("stringToLayerType overlay")
	}
	if stringToLayerType("unknown") != BaseLayer {
		t.Fatal("stringToLayerType unknown should return BaseLayer")
	}
	if layerTypeToString(BaseLayer) != "baselayer" {
		t.Fatal("layerTypeToString baselayer")
	}
	if layerTypeToString(Overlay) != "overlay" {
		t.Fatal("layerTypeToString overlay")
	}
	if layerTypeToString(LayerType(99)) != "" {
		t.Fatal("layerTypeToString invalid")
	}
}

// ─── Metadata ──────────────────────────────────────────────────────────────

func TestMetadataToMap(t *testing.T) {
	md := &Metadata{
		Name:    "test",
		Format:  PNG,
		MinZoom: 0,
		MaxZoom: 10,
		Type:    BaseLayer,
		Center:  [3]float64{0, 0, 4},
		Bounds:  [4]float64{-180, -85, 180, 85},
	}
	m := md.ToMap()
	if m["name"] != "test" {
		t.Fatalf("name = %q", m["name"])
	}
	if m["format"] != "png" {
		t.Fatalf("format = %q", m["format"])
	}
	if m["type"] != "baselayer" {
		t.Fatalf("type = %q", m["type"])
	}
}

func TestMetadataToMapBoundsCenter(t *testing.T) {
	md := &Metadata{
		MinZoom: 0,
		MaxZoom: 5,
		Format:  PBF,
	}
	m := md.ToMap()
	if m["center"] == "" {
		t.Fatal("center should not be empty (zero value)")
	}
	if m["bounds"] == "" {
		t.Fatal("bounds should not be empty (zero value)")
	}
}

// ─── Bounds/Center conversion ──────────────────────────────────────────────

func TestBoundsRoundTrip(t *testing.T) {
	b := [4]float64{-180, -85.05, 180, 85.05}
	s, err := boundsToString(b)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := stringToBounds(s)
	if err != nil {
		t.Fatal(err)
	}
	if parsed != b {
		t.Fatalf("got %v, want %v", parsed, b)
	}
}

func TestCenterRoundTrip(t *testing.T) {
	c := [3]float64{12.5, 41.9, 8}
	s, err := centerToString(c)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := stringToCenter(s)
	if err != nil {
		t.Fatal(err)
	}
	if parsed != c {
		t.Fatalf("got %v, want %v", parsed, c)
	}
}

// ─── ResFactor ─────────────────────────────────────────────────────────────

func TestResFactorRoundTrip(t *testing.T) {
	tests := []interface{}{
		"sqrt2",
		float64(2.0),
		float64(1.5),
	}
	for _, v := range tests {
		s := resFactorToString(v)
		back := stringToResFactor(s)
		switch b := back.(type) {
		case string:
			if b != v {
				t.Fatalf("string: got %q, want %v", b, v)
			}
		case float64:
			if b != v {
				t.Fatalf("float64: got %f, want %v", b, v)
			}
		}
	}
}

func TestStringToResFactor_Invalid(t *testing.T) {
	if v := stringToResFactor("notanumber"); v != nil {
		t.Fatal("expected nil for invalid res_factor")
	}
}

// ─── TileSize ──────────────────────────────────────────────────────────────

func TestTileSizeRoundTrip(t *testing.T) {
	ts := [2]int{256, 256}
	s := tileSizeToString(ts)
	parsed, err := stringToTileSize(s)
	if err != nil {
		t.Fatal(err)
	}
	if *parsed != ts {
		t.Fatalf("got %v, want %v", *parsed, ts)
	}
}

func TestStringToTileSize_Invalid(t *testing.T) {
	_, err := stringToTileSize("notanumber,256")
	if err == nil {
		t.Fatal("expected error")
	}
}

// ─── detectTileFormat ──────────────────────────────────────────────────────

func TestDetectTileFormatPNG(t *testing.T) {
	data := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	fmt, err := detectTileFormat(&data)
	if err != nil {
		t.Fatal(err)
	}
	if fmt != PNG {
		t.Fatalf("got %d, want PNG(%d)", fmt, PNG)
	}
}

func TestDetectTileFormatJPG(t *testing.T) {
	data := []byte{0xFF, 0xD8, 0xFF}
	fmt, err := detectTileFormat(&data)
	if err != nil {
		t.Fatal(err)
	}
	if fmt != JPG {
		t.Fatalf("got %d, want JPG(%d)", fmt, JPG)
	}
}

func TestDetectTileFormatGZip(t *testing.T) {
	data := []byte{0x1f, 0x8b}
	fmt, err := detectTileFormat(&data)
	if err != nil {
		t.Fatal(err)
	}
	if fmt != GZIP {
		t.Fatalf("got %d, want GZIP(%d)", fmt, GZIP)
	}
}

func TestDetectTileFormatZLIB(t *testing.T) {
	data := []byte{0x78, 0x9c}
	fmt, err := detectTileFormat(&data)
	if err != nil {
		t.Fatal(err)
	}
	if fmt != ZLIB {
		t.Fatalf("got %d, want ZLIB(%d)", fmt, ZLIB)
	}
}

func TestDetectTileFormatUnknown(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02}
	_, err := detectTileFormat(&data)
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestDetectTileFormatWEBP(t *testing.T) {
	// WEBP signature with RIFF header
	data := []byte{0x52, 0x49, 0x46, 0x46, 0x00, 0x00, 0x00, 0x00, 0x57, 0x45, 0x42, 0x50, 0x56, 0x50}
	fmt, err := detectTileFormat(&data)
	if err != nil {
		t.Fatal(err)
	}
	if fmt != WEBP {
		t.Fatalf("got %d, want WEBP(%d)", fmt, WEBP)
	}
}

// ─── nil data edge case ────────────────────────────────────────────────────

func TestDetectTileFormatNil(t *testing.T) {
	_, err := detectTileFormat(nil)
	if err == nil {
		t.Fatal("expected error for nil data")
	}
}

func TestDetectTileFormatEmpty(t *testing.T) {
	data := []byte{}
	_, err := detectTileFormat(&data)
	if err == nil {
		t.Fatal("expected error for empty data")
	}
}

// ─── DB Create / Store / Read ──────────────────────────────────────────────

func TestCreateDBAndStoreTile(t *testing.T) {
	tmpFile := t.TempDir() + "/test.mbtiles"
	db, err := CreateDB(tmpFile, PNG, nil)
	if err != nil {
		t.Fatalf("CreateDB: %v", err)
	}
	defer db.Close()

	if err := db.StoreTile(0, 0, 0, []byte("tile data")); err != nil {
		t.Fatalf("StoreTile: %v", err)
	}

	if db.TileFormat() != PNG {
		t.Fatal("TileFormat mismatch")
	}
	if db.TileFormatString() != "png" {
		t.Fatal("TileFormatString mismatch")
	}
	if db.ContentType() != "image/png" {
		t.Fatal("ContentType mismatch")
	}
	if db.HasUTFGrid() {
		t.Fatal("should not have UTFGrid")
	}
}

func TestReadTileNonExistent(t *testing.T) {
	tmpFile := t.TempDir() + "/test.mbtiles"
	db, err := CreateDB(tmpFile, PBF, nil)
	if err != nil {
		t.Fatalf("CreateDB: %v", err)
	}
	defer db.Close()

	var data []byte
	if err := db.ReadTile(0, 0, 0, &data); err != nil {
		t.Fatalf("ReadTile: %v", err)
	}
	if data != nil {
		t.Fatal("expected nil data for non-existent tile")
	}
}

func TestStoreAndReadTile(t *testing.T) {
	tmpFile := t.TempDir() + "/test.mbtiles"
	db, err := CreateDB(tmpFile, PBF, nil)
	if err != nil {
		t.Fatalf("CreateDB: %v", err)
	}
	defer db.Close()

	expected := []byte("hello tile")
	if err := db.StoreTile(5, 100, 200, expected); err != nil {
		t.Fatalf("StoreTile: %v", err)
	}

	var got []byte
	if err := db.ReadTile(5, 100, 200, &got); err != nil {
		t.Fatalf("ReadTile: %v", err)
	}
	if string(got) != string(expected) {
		t.Fatalf("got %q, want %q", got, expected)
	}
}

func TestDBTimeStamp(t *testing.T) {
	tmpFile := t.TempDir() + "/test.mbtiles"
	db, err := CreateDB(tmpFile, PNG, nil)
	if err != nil {
		t.Fatalf("CreateDB: %v", err)
	}
	defer db.Close()

	if db.TimeStamp().IsZero() {
		t.Fatal("expected non-zero timestamp")
	}
}

func TestReadGridWithoutUTFGrid(t *testing.T) {
	tmpFile := t.TempDir() + "/test.mbtiles"
	db, err := CreateDB(tmpFile, PNG, nil)
	if err != nil {
		t.Fatalf("CreateDB: %v", err)
	}
	defer db.Close()

	var data []byte
	if err := db.ReadGrid(0, 0, 0, &data); err == nil {
		t.Fatal("expected error for tileset without UTFGrid")
	}
}

// ─── DB with Metadata ──────────────────────────────────────────────────────

func TestCreateDBWithMetadata(t *testing.T) {
	md := &Metadata{
		Name:    "test",
		Format:  PNG,
		MinZoom: 0,
		MaxZoom: 10,
		Type:    BaseLayer,
	}
	tmpFile := t.TempDir() + "/test.mbtiles"
	db, err := CreateDB(tmpFile, PNG, md)
	if err != nil {
		t.Fatalf("CreateDB: %v", err)
	}
	defer db.Close()

	got, err := db.GetMetadata()
	if err != nil {
		t.Fatalf("GetMetadata: %v", err)
	}
	if got.Name != "test" {
		t.Fatalf("Name = %q", got.Name)
	}
	if got.MinZoom != 0 || got.MaxZoom != 10 {
		t.Fatal("Zoom mismatch")
	}
}

func TestUpdateMetadata(t *testing.T) {
	tmpFile := t.TempDir() + "/test.mbtiles"
	db, err := CreateDB(tmpFile, PBF, nil)
	if err != nil {
		t.Fatalf("CreateDB: %v", err)
	}
	defer db.Close()

	md := &Metadata{Name: "updated", Format: PBF, MinZoom: 0, MaxZoom: 10}
	if err := db.UpdateMetadata(md); err != nil {
		t.Fatalf("UpdateMetadata: %v", err)
	}

	got, err := db.GetMetadata()
	if err != nil {
		t.Fatalf("GetMetadata: %v", err)
	}
	if got.Name != "updated" {
		t.Fatalf("Name = %q", got.Name)
	}
}

// ─── LayerData ─────────────────────────────────────────────────────────────

func TestLayerDataJSON(t *testing.T) {
	ld := &LayerData{
		VectorLayers: &[]VectorLayer{
			{ID: "water", Fields: map[string]interface{}{"depth": "number"}, MinZoom: 0, MaxZoom: 10},
		},
		TileStats: &TileStats{
			LayerCount: 1,
			Layers: []Layer{
				{Name: "water", Count: 100, Geometry: "Polygon", AttributeCount: 1},
			},
		},
	}
	md := &Metadata{Name: "test", Format: PBF, LayerData: ld}
	m := md.ToMap()
	if m["json"] == "" {
		t.Fatal("expected non-empty json metadata")
	}
}

// ─── VectorLayer ───────────────────────────────────────────────────────────

func TestVectorLayer(t *testing.T) {
	vl := VectorLayer{
		ID:          "roads",
		Fields:      map[string]interface{}{"type": "string", "lanes": "number"},
		Description: "Road network",
		MinZoom:     4,
		MaxZoom:     14,
	}
	if vl.ID != "roads" {
		t.Fatal("ID mismatch")
	}
}
