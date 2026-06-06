package recipe

import (
	"encoding/json"
	"testing"
)

func ptr[T any](v T) *T { return &v }

// ─── Vector Recipe ──────────────────────────────────────────────────────────

func TestVectorRecipeMarshal(t *testing.T) {
	r := &Recipe{
		Version:     1,
		FillZoom:    ptr(7),
		Incremental: ptr(true),
		Layers: map[string]*Layer{
			"trees": {Source: "mapbox://tileset-source/user/trees", MinZoom: ptr(uint(4)), MaxZoom: ptr(uint(8))},
		},
	}
	out, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	var decoded Recipe
	if err := json.Unmarshal(out, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Version != 1 {
		t.Fatalf("Version = %d", decoded.Version)
	}
	if decoded.FillZoom == nil || *decoded.FillZoom != 7 {
		t.Fatal("FillZoom mismatch")
	}
	if decoded.Incremental == nil || !*decoded.Incremental {
		t.Fatal("Incremental mismatch")
	}
	if decoded.Layers["trees"] == nil {
		t.Fatal("trees layer missing")
	}
}

func TestVectorRecipeUnmarshal(t *testing.T) {
	raw := `{
		"version": 1,
		"layers": {
			"water": {
				"source": "mapbox://tileset-source/user/water",
				"minzoom": 0,
				"maxzoom": 12
			}
		}
	}`
	var r Recipe
	if err := json.Unmarshal([]byte(raw), &r); err != nil {
		t.Fatal(err)
	}
	if r.Version != 1 {
		t.Fatalf("Version = %d", r.Version)
	}
	layer := r.Layers["water"]
	if layer == nil {
		t.Fatal("water layer missing")
	}
	if layer.Source != "mapbox://tileset-source/user/water" {
		t.Fatalf("Source = %q", layer.Source)
	}
	if *layer.MinZoom != 0 {
		t.Fatalf("MinZoom = %d", *layer.MinZoom)
	}
	if *layer.MaxZoom != 12 {
		t.Fatalf("MaxZoom = %d", *layer.MaxZoom)
	}
}

func TestVectorRecipeWithFeatures(t *testing.T) {
	raw := `{
		"version": 1,
		"layers": {
			"buildings": {
				"source": "mapbox://tileset-source/user/buildings",
				"minzoom": 0,
				"maxzoom": 16,
				"features": {
					"id": ["get", "id"],
					"filter": ["all", [">=", ["zoom"], ["get", "minzoom"]], ["<=", ["zoom"], ["get", "maxzoom"]]],
					"simplification": 4,
					"attributes": {
						"zoom_element": ["name"],
						"set": {"height": ["get", "height"]},
						"allowed_output": ["name", "height"]
					}
				}
			}
		}
	}`
	var r Recipe
	if err := json.Unmarshal([]byte(raw), &r); err != nil {
		t.Fatal(err)
	}
	f := r.Layers["buildings"].Features
	if f == nil {
		t.Fatal("Features config missing")
	}
	if f.Simplification == nil {
		t.Fatal("Simplification missing")
	}
	if _, ok := f.Simplification.(float64); !ok || f.Simplification.(float64) != 4 {
		t.Fatal("Simplification expected 4")
	}
	if f.Attributes == nil || len(f.Attributes.ZoomElement) != 1 || f.Attributes.ZoomElement[0] != "name" {
		t.Fatal("ZoomElement mismatch")
	}
	if len(f.Attributes.AllowedOutput) != 2 {
		t.Fatal("AllowedOutput mismatch")
	}
}

func TestVectorRecipeWithTiles(t *testing.T) {
	raw := `{
		"version": 1,
		"layers": {
			"roads": {
				"source": "mapbox://tileset-source/user/roads",
				"minzoom": 4,
				"maxzoom": 14,
				"tiles": {
					"extent": 4096,
					"buffer_size": 4,
					"limit": [["highest_where", ["==", ["get", "class"], "motorway"], 5, "priority"]],
					"order": "priority",
					"layer_size": 500
				}
			}
		}
	}`
	var r Recipe
	if err := json.Unmarshal([]byte(raw), &r); err != nil {
		t.Fatal(err)
	}
	tc := r.Layers["roads"].Tiles
	if tc == nil {
		t.Fatal("Tiles config missing")
	}
	if tc.Order != "priority" {
		t.Fatalf("Order = %q", tc.Order)
	}
	if tc.LayerSize == nil || *tc.LayerSize != 500 {
		t.Fatal("LayerSize mismatch")
	}
	if len(tc.Limit) != 1 {
		t.Fatal("Limit mismatch")
	}
}

func TestVectorRecipeWithUnion(t *testing.T) {
	raw := `{
		"version": 1,
		"layers": {
			"blocks": {
				"source": "mapbox://tileset-source/user/blocks",
				"minzoom": 0,
				"maxzoom": 10,
				"tiles": {
					"union": [
						{
							"where": ["==", ["get", "type"], "building"],
							"group_by": ["height"],
							"aggregate": {"area": "sum"},
							"maintain_direction": true
						}
					]
				}
			}
		}
	}`
	var r Recipe
	if err := json.Unmarshal([]byte(raw), &r); err != nil {
		t.Fatal(err)
	}
	unions := r.Layers["blocks"].Tiles.Union
	if len(unions) != 1 {
		t.Fatal("expected 1 union")
	}
	u := unions[0]
	if u.Aggregate["area"] != "sum" {
		t.Fatalf("Aggregate area = %q", u.Aggregate["area"])
	}
	if u.MaintainDirection == nil || !*u.MaintainDirection {
		t.Fatal("MaintainDirection should be true")
	}
}

// ─── Raster Recipe ──────────────────────────────────────────────────────────

func TestRasterRecipeRoundTrip(t *testing.T) {
	r := &Recipe{
		Version: 1,
		Type:    RecipeRaster,
		Sources: []Source{{URI: "mapbox://tileset-source/user/source"}},
		MinZoom: ptr(0),
		MaxZoom: ptr(18),
		Layers: map[string]*Layer{
			"RGB": {
				SourceRules: &SourceRules{
					Filter: []interface{}{
						[]interface{}{"==", []interface{}{"get", "colorinterp"}, "red"},
						[]interface{}{"==", []interface{}{"get", "colorinterp"}, "green"},
						[]interface{}{"==", []interface{}{"get", "colorinterp"}, "blue"},
					},
				},
			},
		},
	}
	out, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	var decoded Recipe
	if err := json.Unmarshal(out, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Type != RecipeRaster {
		t.Fatalf("Type = %q", decoded.Type)
	}
	if len(decoded.Sources) != 1 {
		t.Fatal("Sources mismatch")
	}
	if *decoded.MinZoom != 0 || *decoded.MaxZoom != 18 {
		t.Fatal("Zoom range mismatch")
	}
}

func TestRasterRecipeUnmarshal(t *testing.T) {
	raw := `{
		"version": 1,
		"type": "raster",
		"sources": [{"uri": "mapbox://tileset-source/user/imagery"}],
		"minzoom": 10,
		"maxzoom": 16,
		"layers": {
			"RGBA": {}
		}
	}`
	var r Recipe
	if err := json.Unmarshal([]byte(raw), &r); err != nil {
		t.Fatal(err)
	}
	if !r.IsRaster() {
		t.Fatal("expected raster type")
	}
}

// ─── RasterArray Recipe ─────────────────────────────────────────────────────

func TestRasterArrayRecipeRoundTrip(t *testing.T) {
	r := &Recipe{
		Version: 1,
		Type:    RecipeRasterArray,
		Sources: []Source{{URI: "mapbox://tileset-source/user/weather"}},
		Layers: map[string]*Layer{
			"Pressure": {
				TileSize: ptr(512),
				Offset:   ptr(-100.0),
				Scale:    ptr(0.1),
				Units:    "Pa",
				SourceRules: &SourceRules{
					Name:  []interface{}{"to-number", []interface{}{"get", "GRIB_VALID_TIME"}},
					Order: "asc",
				},
			},
		},
	}
	out, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	var decoded Recipe
	if err := json.Unmarshal(out, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Type != RecipeRasterArray {
		t.Fatalf("Type = %q", decoded.Type)
	}
	l := decoded.Layers["Pressure"]
	if l == nil {
		t.Fatal("Pressure layer missing")
	}
	if l.TileSize == nil || *l.TileSize != 512 {
		t.Fatal("TileSize mismatch")
	}
	if l.Offset == nil || *l.Offset != -100 {
		t.Fatal("Offset mismatch")
	}
}

func TestRasterArrayWithOptionalZoom(t *testing.T) {
	raw := `{
		"version": 1,
		"type": "rasterarray",
		"sources": [{"uri": "mapbox://tileset-source/user/data"}],
		"minzoom": 0,
		"maxzoom": 6,
		"layers": {"band1": {}}
	}`
	var r Recipe
	if err := json.Unmarshal([]byte(raw), &r); err != nil {
		t.Fatal(err)
	}
	if r.MinZoom == nil || *r.MinZoom != 0 {
		t.Fatal("MinZoom mismatch")
	}
}

// ─── Validation ─────────────────────────────────────────────────────────────

func TestValidate_Vector(t *testing.T) {
	r := &Recipe{Version: 1, Layers: map[string]*Layer{"test": {Source: "src", MinZoom: ptr(uint(0)), MaxZoom: ptr(uint(10))}}}
	if err := r.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_WrongVersion(t *testing.T) {
	r := &Recipe{Version: 2}
	if err := r.Validate(); err == nil {
		t.Fatal("expected error for wrong version")
	}
}

func TestValidate_NoLayers(t *testing.T) {
	r := &Recipe{Version: 1}
	if err := r.Validate(); err == nil {
		t.Fatal("expected error for no layers")
	}
}

func TestValidate_TooManyLayers(t *testing.T) {
	layers := make(map[string]*Layer)
	for i := 0; i < 21; i++ {
		layers[string(rune('a'+i))] = &Layer{Source: "src", MinZoom: ptr(uint(0)), MaxZoom: ptr(uint(10))}
	}
	r := &Recipe{Version: 1, Layers: layers}
	if err := r.Validate(); err == nil {
		t.Fatal("expected error for >20 layers")
	}
}

func TestValidate_VectorMissingSource(t *testing.T) {
	r := &Recipe{Version: 1, Layers: map[string]*Layer{"test": {MinZoom: ptr(uint(0)), MaxZoom: ptr(uint(10))}}}
	if err := r.Validate(); err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestValidate_VectorMinGtMax(t *testing.T) {
	r := &Recipe{Version: 1, Layers: map[string]*Layer{"test": {Source: "src", MinZoom: ptr(uint(10)), MaxZoom: ptr(uint(5))}}}
	if err := r.Validate(); err == nil {
		t.Fatal("expected error for min > max")
	}
}

func TestValidate_Raster(t *testing.T) {
	r := &Recipe{Version: 1, Type: RecipeRaster, Sources: []Source{{URI: "src"}}, MinZoom: ptr(0), MaxZoom: ptr(10),
		Layers: map[string]*Layer{"RGB": {}}}
	if err := r.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_RasterMissingSources(t *testing.T) {
	r := &Recipe{Version: 1, Type: RecipeRaster, MinZoom: ptr(0), MaxZoom: ptr(10), Layers: map[string]*Layer{"RGB": {}}}
	if err := r.Validate(); err == nil {
		t.Fatal("expected error for missing sources")
	}
}

func TestValidate_RasterArray(t *testing.T) {
	r := &Recipe{Version: 1, Type: RecipeRasterArray, Sources: []Source{{URI: "src"}}, Layers: map[string]*Layer{"b": {}}}
	if err := r.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ─── Helpers ────────────────────────────────────────────────────────────────

func TestHelpers(t *testing.T) {
	r := &Recipe{Version: 1, Layers: map[string]*Layer{"a": {}, "b": {}}}
	names := r.LayerNames()
	if len(names) != 2 {
		t.Fatalf("expected 2 layer names, got %d", len(names))
	}
	if !r.IsVector() {
		t.Fatal("default recipe should be vector")
	}
	r.Type = RecipeRaster
	if !r.IsRaster() {
		t.Fatal("expected raster")
	}
	if r.IsVector() {
		t.Fatal("raster should not be vector")
	}
}

func TestAddLayer(t *testing.T) {
	r := &Recipe{Version: 1}
	if err := r.AddLayer("test", &Layer{Source: "src", MinZoom: ptr(uint(0)), MaxZoom: ptr(uint(10))}); err != nil {
		t.Fatal(err)
	}
	if len(r.Layers) != 1 {
		t.Fatal("expected 1 layer")
	}
	if err := r.AddLayer("test", nil); err == nil {
		t.Fatal("expected error for duplicate name")
	}
	if err := r.AddLayer("invalid name!", nil); err == nil {
		t.Fatal("expected error for invalid name")
	}
	if err := r.AddLayer("valid_name", &Layer{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRemoveLayer(t *testing.T) {
	r := &Recipe{Version: 1, Layers: map[string]*Layer{"a": {}}}
	r.RemoveLayer("a")
	if len(r.Layers) != 0 {
		t.Fatal("expected 0 layers after removal")
	}
	r.RemoveLayer("nonexistent") // should not panic
}

func TestNewHelpers(t *testing.T) {
	r := NewVectorRecipe(map[string]*Layer{"t": NewLayer("src", 0, 10)})
	if r.Version != 1 || !r.IsVector() || len(r.Layers) != 1 {
		t.Fatal("NewVectorRecipe failed")
	}

	r2 := NewRasterRecipe([]Source{{URI: "src"}}, 0, 10, map[string]*Layer{"RGB": {}})
	if !r2.IsRaster() || *r2.MinZoom != 0 || *r2.MaxZoom != 10 {
		t.Fatal("NewRasterRecipe failed")
	}

	r3 := NewRasterArrayRecipe([]Source{{URI: "src"}}, map[string]*Layer{"b": {}})
	if !r3.IsRasterArray() {
		t.Fatal("NewRasterArrayRecipe failed")
	}
}

func TestRasterLayerOptions(t *testing.T) {
	l := NewRasterLayer(
		WithTileSize(256),
		WithBuffer(1),
		WithUnits("C"),
		WithOffset(-100),
		WithScale(0.1),
		WithResampling("bilinear"),
		WithInputNoData(9999),
		WithSourceRules(&SourceRules{Order: "asc"}),
	)
	if l.TileSize == nil || *l.TileSize != 256 {
		t.Fatal("TileSize mismatch")
	}
	if l.Buffer == nil || *l.Buffer != 1 {
		t.Fatal("Buffer mismatch")
	}
	if l.Units != "C" {
		t.Fatal("Units mismatch")
	}
	if l.Offset == nil || *l.Offset != -100 {
		t.Fatal("Offset mismatch")
	}
	if l.Resampling != "bilinear" {
		t.Fatal("Resampling mismatch")
	}
	if l.SourceRules.Order != "asc" {
		t.Fatal("SourceRules mismatch")
	}
}

func TestString(t *testing.T) {
	r := &Recipe{Version: 1, Layers: map[string]*Layer{"a": {}}}
	s := r.String()
	if s == "" {
		t.Fatal("expected non-empty string")
	}
}

func TestValidLayerName(t *testing.T) {
	if !validLayerName("test") {
		t.Error("'test' should be valid")
	}
	if !validLayerName("test_123") {
		t.Error("'test_123' should be valid")
	}
	if validLayerName("") {
		t.Error("empty name should be invalid")
	}
	if validLayerName("has space") {
		t.Error("name with space should be invalid")
	}
	if validLayerName("has-hyphen") {
		t.Error("name with hyphen should be invalid")
	}
}

// ─── Json edge cases ────────────────────────────────────────────────────────

func TestRecipeMarshalUnmarshalNil(t *testing.T) {
	var r Recipe
	out, err := json.Marshal(&r)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(out, &r); err != nil {
		t.Fatal(err)
	}
}

func TestSourceCRS(t *testing.T) {
	raw := `{"uri":"mapbox://tileset-source/user/data","crs":"EPSG:4326"}`
	var s Source
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatal(err)
	}
	if s.CRS != "EPSG:4326" {
		t.Fatalf("CRS = %q", s.CRS)
	}
}
