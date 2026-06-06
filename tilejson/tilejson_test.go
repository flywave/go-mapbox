package tilejson

import (
	"encoding/json"
	"testing"
)

func ptr[T any](v T) *T { return &v }

// ─── Basic JSON round trip ─────────────────────────────────────────────────

func TestNewMinimal(t *testing.T) {
	tj := New([]string{"http://example.com/{z}/{x}/{y}.pbf"})
	if tj.TileJSON != SpecVersion {
		t.Fatalf("TileJSON = %q", tj.TileJSON)
	}
	if len(tj.Tiles) != 1 {
		t.Fatal("expected 1 tile")
	}
	if tj.Scheme != SchemeXYZ {
		t.Fatalf("Scheme = %q", tj.Scheme)
	}
	if tj.Version != DefaultVersion {
		t.Fatalf("Version = %q", tj.Version)
	}
}

func TestJSONRoundTrip(t *testing.T) {
	tj := &TileJSON{
		TileJSON: SpecVersion,
		Tiles:    []string{"http://example.com/{z}/{x}/{y}.pbf"},
		Name:     ptr("test tileset"),
		MinZoom:  ptr(0),
		MaxZoom:  ptr(14),
		Scheme:   SchemeXYZ,
		Version:  "1.2.0",
		Bounds:   NewTileBounds(-180, -85, 180, 85),
		Center:   NewTileCenter(0, 0, 4),
	}
	out, err := json.Marshal(tj)
	if err != nil {
		t.Fatal(err)
	}
	var decoded TileJSON
	if err := json.Unmarshal(out, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.TileJSON != SpecVersion {
		t.Fatal("TileJSON version mismatch")
	}
	if *decoded.Name != "test tileset" {
		t.Fatal("Name mismatch")
	}
	if *decoded.MinZoom != 0 || *decoded.MaxZoom != 14 {
		t.Fatal("Zoom mismatch")
	}
}

func TestJSONUnmarshalFull(t *testing.T) {
	raw := `{
		"tilejson": "3.0.0",
		"tiles": ["http://example.com/{z}/{x}/{y}.mvt"],
		"name": "Natural Earth",
		"description": "Admin 0 countries",
		"version": "1.0.0",
		"attribution": "<a href='https://openstreetmap.org'>OSM</a>",
		"minzoom": 0,
		"maxzoom": 6,
		"bounds": [-180, -85, 180, 85],
		"center": [0, 30, 3],
		"fillzoom": 4,
		"scheme": "xyz",
		"vector_layers": [{
			"id": "countries",
			"description": "Country boundaries",
			"minzoom": 0,
			"maxzoom": 6,
			"fields": {"name": "Country name", "iso_a2": "ISO 3166-1 code"}
		}],
		"legend": "Legend text",
		"template": "{{name}}",
		"grids": ["http://example.com/{z}/{x}/{y}.grid.json"],
		"data": ["http://example.com/data.geojson"]
	}`
	var tj TileJSON
	if err := json.Unmarshal([]byte(raw), &tj); err != nil {
		t.Fatal(err)
	}
	if tj.Name == nil || *tj.Name != "Natural Earth" {
		t.Fatal("Name mismatch")
	}
	if tj.Description == nil || *tj.Description != "Admin 0 countries" {
		t.Fatal("Description mismatch")
	}
	if tj.Attribution == nil || *tj.Attribution != "<a href='https://openstreetmap.org'>OSM</a>" {
		t.Fatal("Attribution mismatch")
	}
	if tj.FillZoom == nil || *tj.FillZoom != 4 {
		t.Fatal("FillZoom mismatch")
	}
	if tj.Legend == nil || *tj.Legend != "Legend text" {
		t.Fatal("Legend mismatch")
	}
	if tj.Template == nil || *tj.Template != "{{name}}" {
		t.Fatal("Template mismatch")
	}
	if len(tj.Grids) != 1 {
		t.Fatal("Grids mismatch")
	}
	if len(tj.Data) != 1 {
		t.Fatal("Data mismatch")
	}
	if len(tj.VectorLayers) != 1 {
		t.Fatal("VectorLayers mismatch")
	}
	vl := tj.VectorLayers[0]
	if vl.ID != "countries" {
		t.Fatalf("VectorLayer ID = %q", vl.ID)
	}
	if vl.Fields["name"] != "Country name" {
		t.Fatal("VectorLayer fields mismatch")
	}
	if vl.Description == nil || *vl.Description != "Country boundaries" {
		t.Fatal("VectorLayer description mismatch")
	}
}

// ─── VectorLayer ───────────────────────────────────────────────────────────

func TestVectorLayerMarshal(t *testing.T) {
	tj := New([]string{"http://example.com/{z}/{x}/{y}.mvt"})
	tj.AddVectorLayer(NewVectorLayer("roads", map[string]string{"type": "Road classification", "name": "Street name"}))
	tj.VectorLayers[0].MinZoom = ptr(4)
	tj.VectorLayers[0].MaxZoom = ptr(14)

	out, err := json.Marshal(tj)
	if err != nil {
		t.Fatal(err)
	}
	var decoded TileJSON
	json.Unmarshal(out, &decoded)
	if len(decoded.VectorLayers) != 1 {
		t.Fatal("expected 1 vector layer")
	}
	vl := decoded.VectorLayers[0]
	if vl.ID != "roads" {
		t.Fatalf("ID = %q", vl.ID)
	}
	if len(vl.Fields) != 2 {
		t.Fatal("expected 2 fields")
	}
	if vl.MinZoom == nil || *vl.MinZoom != 4 {
		t.Fatal("MinZoom mismatch")
	}
}

func TestVectorLayerEmptyFields(t *testing.T) {
	tj := New([]string{"http://example.com/{z}/{x}/{y}.mvt"})
	tj.AddVectorLayer(NewVectorLayer("buildings", map[string]string{}))
	out, _ := json.Marshal(tj)
	var decoded TileJSON
	json.Unmarshal(out, &decoded)
	vl := decoded.VectorLayers[0]
	if vl.Fields == nil {
		t.Fatal("fields should be empty object, not null")
	}
}

// ─── Helpers ───────────────────────────────────────────────────────────────

func TestNewTileBounds(t *testing.T) {
	b := NewTileBounds(-180, -85, 180, 85)
	if b[0] != -180 || b[1] != -85 || b[2] != 180 || b[3] != 85 {
		t.Fatal("bounds values mismatch")
	}
}

func TestNewTileCenter(t *testing.T) {
	c := NewTileCenter(-122.4, 37.8, 10)
	if c[0] != -122.4 || c[1] != 37.8 || c[2] != 10 {
		t.Fatal("center values mismatch")
	}
}

func TestAddTile(t *testing.T) {
	tj := New([]string{"http://a.com/{z}/{x}/{y}.pbf"})
	tj.AddTile("http://b.com/{z}/{x}/{y}.pbf")
	if len(tj.Tiles) != 2 {
		t.Fatal("expected 2 tiles")
	}
}

// ─── Validation ────────────────────────────────────────────────────────────

func TestValidateValid(t *testing.T) {
	tj := New([]string{"http://example.com/{z}/{x}/{y}.pbf"})
	if err := tj.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateSpecVersion(t *testing.T) {
	tj := New([]string{"http://example.com/{z}/{x}/{y}.pbf"})
	tj.TileJSON = "2.2.0"
	if err := tj.Validate(); err == nil {
		t.Fatal("expected error for wrong spec version")
	}
}

func TestValidateNoTiles(t *testing.T) {
	tj := &TileJSON{TileJSON: SpecVersion}
	if err := tj.Validate(); err == nil {
		t.Fatal("expected error for no tiles")
	}
}

func TestValidateZoomRange(t *testing.T) {
	tj := New([]string{"http://example.com/{z}/{x}/{y}.pbf"})
	tj.MinZoom = ptr(10)
	tj.MaxZoom = ptr(5)
	if err := tj.Validate(); err == nil {
		t.Fatal("expected error for min > max")
	}
}

func TestValidateZoomBounds(t *testing.T) {
	tj := New([]string{"http://example.com/{z}/{x}/{y}.pbf"})
	tj.MaxZoom = ptr(31)
	if err := tj.Validate(); err == nil {
		t.Fatal("expected error for zoom > 30")
	}
}

func TestValidateInvalidScheme(t *testing.T) {
	tj := New([]string{"http://example.com/{z}/{x}/{y}.pbf"})
	tj.Scheme = "invalid"
	if err := tj.Validate(); err == nil {
		t.Fatal("expected error for invalid scheme")
	}
}

func TestValidateVectorLayerMissingID(t *testing.T) {
	tj := New([]string{"http://example.com/{z}/{x}/{y}.pbf"})
	tj.AddVectorLayer(VectorLayer{Fields: map[string]string{}})
	if err := tj.Validate(); err == nil {
		t.Fatal("expected error for missing vector layer id")
	}
}

func TestValidateVectorLayerMissingFields(t *testing.T) {
	tj := New([]string{"http://example.com/{z}/{x}/{y}.pbf"})
	tj.AddVectorLayer(VectorLayer{ID: "test"})
	if err := tj.Validate(); err == nil {
		t.Fatal("expected error for missing vector layer fields")
	}
}

// ─── Marshal/Unmarshal edge cases ──────────────────────────────────────────

func TestMarshalOmitEmpty(t *testing.T) {
	tj := New([]string{"http://example.com/{z}/{x}/{y}.pbf"})
	out, err := json.Marshal(tj)
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]interface{}
	json.Unmarshal(out, &raw)
	if raw["attribution"] != nil {
		t.Fatal("attribution should be omitted")
	}
	if raw["name"] != nil {
		t.Fatal("name should be omitted")
	}
	if raw["bounds"] != nil {
		t.Fatal("bounds should be omitted")
	}
}

func TestMarshalDefaultValues(t *testing.T) {
	tj := New([]string{"http://example.com/{z}/{x}/{y}.pbf"})
	out, _ := json.Marshal(tj)
	var raw map[string]interface{}
	json.Unmarshal(out, &raw)
	if raw["tilejson"] != SpecVersion {
		t.Fatal("tilejson version should be set")
	}
	if raw["scheme"] != "xyz" {
		t.Fatal("scheme should default to xyz")
	}
}

func TestMarshalVectorLayer(t *testing.T) {
	tj := New([]string{"http://example.com/{z}/{x}/{y}.mvt"})
	tj.AddVectorLayer(VectorLayer{ID: "water", Fields: map[string]string{"depth": "Water depth in meters"}})
	out, _ := json.MarshalIndent(tj, "", "  ")
	var decoded TileJSON
	json.Unmarshal(out, &decoded)
	if len(decoded.VectorLayers) != 1 {
		t.Fatal("vector_layers round trip failed")
	}
}

func TestFillZoom(t *testing.T) {
	tj := New([]string{"http://example.com/{z}/{x}/{y}.pbf"})
	tj.FillZoom = ptr(7)
	out, _ := json.Marshal(tj)
	var decoded TileJSON
	json.Unmarshal(out, &decoded)
	if decoded.FillZoom == nil || *decoded.FillZoom != 7 {
		t.Fatal("FillZoom round trip failed")
	}
}

func TestEmptyTileset(t *testing.T) {
	tj := &TileJSON{}
	out, err := json.Marshal(tj)
	if err != nil {
		t.Fatal(err)
	}
	// Should still produce valid JSON
	var raw map[string]interface{}
	if err := json.Unmarshal(out, &raw); err != nil {
		t.Fatal(err)
	}
}
