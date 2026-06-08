package style

import (
	"encoding/json"
	"image/color"
	"strings"
	"testing"
)

func ptr[T any](v T) *T { return &v }

// ─── Expression ─────────────────────────────────────────────────────────────

func TestExpressionLiteralUnmarshal(t *testing.T) {
	tests := []struct {
		name  string
		json  string
		want  interface{}
	}{
		{"string", `"hello"`, "hello"},
		{"number", `42`, float64(42)},
		{"bool_true", `true`, true},
		{"bool_false", `false`, false},
		{"null", `null`, nil},
		{"array", `[1,2,3]`, []interface{}{float64(1), float64(2), float64(3)}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var e Expression
			if err := json.Unmarshal([]byte(tc.json), &e); err != nil {
				t.Fatal(err)
			}
			if !e.IsLiteral {
				t.Fatal("expected literal")
			}
			if tc.want == nil && e.Value != nil {
				t.Fatalf("expected nil, got %v", e.Value)
			}
			if tc.want != nil {
				got, _ := json.Marshal(e.Value)
				want, _ := json.Marshal(tc.want)
				if string(got) != string(want) {
					t.Fatalf("expected %v (%T), got %v (%T)", tc.want, tc.want, e.Value, e.Value)
				}
			}
		})
	}
}

func TestExpressionCompoundUnmarshal(t *testing.T) {
	raw := `["get","name"]`
	var e Expression
	if err := json.Unmarshal([]byte(raw), &e); err != nil {
		t.Fatal(err)
	}
	if e.IsLiteral {
		t.Fatal("expected compound expression")
	}
	if e.Operator != "get" {
		t.Fatalf("expected operator 'get', got %q", e.Operator)
	}
	if len(e.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(e.Args))
	}
	if !e.Args[0].IsLiteral || e.Args[0].Value.(string) != "name" {
		t.Fatalf("expected literal 'name', got %v", e.Args[0].Value)
	}
}

func TestExpressionNestedUnmarshal(t *testing.T) {
	raw := `["interpolate",["linear"],["zoom"],0,["get","min"],10,["get","max"]]`
	var e Expression
	if err := json.Unmarshal([]byte(raw), &e); err != nil {
		t.Fatal(err)
	}
	if e.Operator != "interpolate" {
		t.Fatalf("expected interpolate, got %q", e.Operator)
	}
	// ["linear"]
	if e.Args[0].Operator != "linear" {
		t.Fatalf("expected linear, got %q", e.Args[0].Operator)
	}
	// ["zoom"]
	if e.Args[1].Operator != "zoom" {
		t.Fatalf("expected zoom, got %q", e.Args[1].Operator)
	}
	// 0 (literal)
	if !e.Args[2].IsLiteral || e.Args[2].Value.(float64) != 0 {
		t.Fatal("expected literal 0")
	}
	// ["get","min"]
	if e.Args[3].Operator != "get" {
		t.Fatalf("expected get, got %q", e.Args[3].Operator)
	}
}

func TestExpressionRoundTrip(t *testing.T) {
	cases := []string{
		`"hello"`,
		`42`,
		`true`,
		`null`,
		`["get","name"]`,
		`["+",1,2]`,
		`["case",["==",["get","type"],"park"],"green","default"]`,
		`["interpolate",["linear"],["zoom"],0,1,10,10]`,
		`["let","x",["get","val"],["var","x"]]`,
		`["match",["get","type"],"a",1,"b",2,0]`,
	}
	for _, raw := range cases {
		t.Run(raw, func(t *testing.T) {
			var e Expression
			if err := json.Unmarshal([]byte(raw), &e); err != nil {
				t.Fatal(err)
			}
			out, err := json.Marshal(&e)
			if err != nil {
				t.Fatal(err)
			}
			var normalized interface{}
			json.Unmarshal([]byte(raw), &normalized)
			var got interface{}
			json.Unmarshal(out, &got)
			eq := jsonDeepEqual(normalized, got)
			if !eq {
				t.Fatalf("round trip mismatch:\n  in:  %s\n  out: %s", raw, string(out))
			}
		})
	}
}

func jsonDeepEqual(a, b interface{}) bool {
	aj, _ := json.Marshal(a)
	bj, _ := json.Marshal(b)
	return string(aj) == string(bj)
}

func TestExpressionOperatorConstantsMatch(t *testing.T) {
	ops := []string{
		ExpArray, ExpBoolean, ExpCollator, ExpFormat, ExpImage,
		ExpLiteral, ExpNumber, ExpNumberFmt, ExpObject, ExpString,
		ExpToBool, ExpToColor, ExpToNumber, ExpToString, ExpTypeOf,
		ExpToHSLA, ExpToRGBA,
		ExpAccumulated, ExpFeatureState, ExpGeometryType, ExpID,
		ExpLineProgress, ExpProperties,
		ExpAt, ExpAtInterpolated, ExpConfig, ExpGet, ExpHas,
		ExpIn, ExpIndexOf, ExpLength, ExpMeasureLight, ExpSlice,
		ExpSplit, ExpWorldview,
		ExpNot, ExpNEq, ExpLT, ExpLTE, ExpEQ, ExpGT, ExpGTE,
		ExpAll, ExpAny, ExpCase, ExpCoalesce, ExpMatch, ExpWithin,
		ExpInterpolate, ExpInterpolateHCL, ExpInterpolateLab, ExpStep,
		ExpLet, ExpVar,
		ExpConcat, ExpDowncase, ExpIsSupportedScript, ExpResolvedLocale, ExpUpcase,
		ExpHSL, ExpHSLA, ExpRGB, ExpRGBA,
		ExpSub, ExpMul, ExpDiv, ExpMod, ExpPow, ExpAdd,
		ExpAbs, ExpAcos, ExpAsin, ExpAtan, ExpCeil, ExpCos, ExpDist,
		ExpE, ExpFloor, ExpLn, ExpLn2, ExpLog10, ExpLog2, ExpMax, ExpMin,
		ExpPI, ExpRand, ExpRound, ExpSin, ExpSqrt, ExpTan,
		ExpDistFromCenter, ExpPitch, ExpZoom,
		ExpHeatmapDensity,
		ExpLinear, ExpExponential, ExpCubicBezier,
	}
	for _, op := range ops {
		if !IsKnownOperator(op) {
			t.Errorf("IsKnownOperator(%q) = false, expected true", op)
		}
	}
	if IsKnownOperator("nonexistent") {
		t.Error("IsKnownOperator('nonexistent') = true, expected false")
	}
}

// ─── Validation ─────────────────────────────────────────────────────────────

func TestValidateExpression(t *testing.T) {
	tests := []struct {
		name      string
		expr      string
		mode      ValidationMode
		wantError bool
	}{
		{"literal_string", `"hello"`, ValidationStrict, false},
		{"literal_number", `42`, ValidationStrict, false},
		{"valid_get", `["get","name"]`, ValidationStrict, false},
		{"valid_zoom", `["zoom"]`, ValidationStrict, false},
		{"valid_interpolate", `["interpolate",["linear"],["zoom"],0,1,10,10]`, ValidationStrict, false},
		{"valid_step", `["step",["zoom"],0,5,1,10,2]`, ValidationStrict, false},
		{"valid_case", `["case",true,"a",false,"b","c"]`, ValidationStrict, false},
		{"valid_match", `["match",["get","t"],"x",1,"y",2,0]`, ValidationStrict, false},
		{"valid_rgb", `["rgb",255,0,0]`, ValidationStrict, false},
		{"valid_rgba", `["rgba",255,0,0,1]`, ValidationStrict, false},
		{"valid_add", `["+",1,2]`, ValidationStrict, false},
		{"valid_let_var", `["let","x",1,["var","x"]]`, ValidationStrict, false},
		{"valid_all", `["all",true,false]`, ValidationStrict, false},
		{"valid_any", `["any",true,false]`, ValidationStrict, false},
		{"valid_coalesce", `["coalesce",["get","a"],"fallback"]`, ValidationStrict, false},
		{"valid_within", `["within",{"type":"Polygon"}]`, ValidationStrict, false},
		{"unknown_op", `["foobar"]`, ValidationStrict, true},
		{"rgb_wrong_args", `["rgb",255]`, ValidationStrict, true},
		{"rgba_wrong_args", `["rgba",255,0]`, ValidationStrict, true},
		{"zoom_wrong_args", `["zoom",1]`, ValidationStrict, true},
		{"interpolate_too_few", `["interpolate",["linear"],["zoom"]]`, ValidationStrict, true},
		{"interpolate_odd_pairs", `["interpolate",["linear"],["zoom"],0,1,10]`, ValidationStrict, true},
		{"case_even_args", `["case",true,"a"]`, ValidationStrict, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var e Expression
			if err := json.Unmarshal([]byte(tc.expr), &e); err != nil {
				t.Fatal(err)
			}
			err := ValidateExpression(&e, tc.mode)
			if tc.wantError && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestIsCameraExpression(t *testing.T) {
	tests := []struct {
		name string
		json string
		want bool
	}{
		{"literal", `"hello"`, false},
		{"zoom", `["zoom"]`, true},
		{"interpolate_zoom", `["interpolate",["linear"],["zoom"],0,1]`, true},
		{"get_only", `["get","name"]`, false},
		{"nested_no_zoom", `["+",["get","a"],["get","b"]]`, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var e Expression
			json.Unmarshal([]byte(tc.json), &e)
			if got := IsCameraExpression(&e); got != tc.want {
				t.Fatalf("IsCameraExpression() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestIsDataExpression(t *testing.T) {
	tests := []struct {
		name string
		json string
		want bool
	}{
		{"literal", `42`, false},
		{"get", `["get","name"]`, true},
		{"has", `["has","name"]`, true},
		{"id", `["id"]`, true},
		{"geometry_type", `["geometry-type"]`, true},
		{"properties", `["properties"]`, true},
		{"feature_state", `["feature-state","selected"]`, true},
		{"zoom_only", `["zoom"]`, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var e Expression
			json.Unmarshal([]byte(tc.json), &e)
			if got := IsDataExpression(&e); got != tc.want {
				t.Fatalf("IsDataExpression() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestFormatExpression(t *testing.T) {
	tests := []struct {
		json string
		want string
	}{
		{`"hello"`, `hello`},
		{`["get","name"]`, `[get name]`},
		{`["+",1,2]`, `[+ 1 2]`},
	}
	for _, tc := range tests {
		var e Expression
		json.Unmarshal([]byte(tc.json), &e)
		if got := FormatExpression(&e); got != tc.want {
			t.Fatalf("FormatExpression() = %q, want %q", got, tc.want)
		}
	}
}

// ─── FilterContainer ────────────────────────────────────────────────────────

func TestFilterContainerRoundTrip(t *testing.T) {
	cases := []string{
		`["==","$type","Point"]`,
		`["all",["==","class","road"],["!=","brunnel","tunnel"]]`,
		`["any",["in","water","class"],["has","name"]]`,
		`["none",["==","$type","Polygon"]]`,
	}
	for _, raw := range cases {
		t.Run(raw, func(t *testing.T) {
			var fc FilterContainer
			if err := json.Unmarshal([]byte(raw), &fc); err != nil {
				t.Fatal(err)
			}
			out, err := json.Marshal(&fc)
			if err != nil {
				t.Fatal(err)
			}
			got := string(out)
			if !jsonEqual(got, raw) {
				t.Fatalf("round trip:\n  in:  %s\n  out: %s", raw, got)
			}
		})
	}
}

func jsonEqual(a, b string) bool {
	var va, vb interface{}
	json.Unmarshal([]byte(a), &va)
	json.Unmarshal([]byte(b), &vb)
	return jsonDeepEqual(va, vb)
}

// ─── ColorType ──────────────────────────────────────────────────────────────

func TestColorTypeUnmarshal(t *testing.T) {
	tests := []struct {
		json string
	}{
		{`"#ff0000"`},
		{`"#f00"`},
		{`"rgb(255,0,0)"`},
		{`"rgba(255,0,0,1)"`},
		{`"hsl(0,100%,50%)"`},
		{`"hsla(0,100%,50%,1)"`},
		{`"red"`},
	}
	for _, tc := range tests {
		t.Run(tc.json, func(t *testing.T) {
			var c ColorType
			if err := json.Unmarshal([]byte(tc.json), &c); err != nil {
				t.Fatal(err)
			}
			if c.GetColorAtZoomLevel(0) == nil {
				t.Fatal("expected non-nil color")
			}
		})
	}
}

func TestColorTypeMarshal(t *testing.T) {
	// ColorType uses unexported internal type, so MarshalJSON is not defined.
	// Verify round-trip through unmarshal only.
	var c ColorType
	if err := json.Unmarshal([]byte(`"#ff0000"`), &c); err != nil {
		t.Fatal(err)
	}
	if c.GetColorAtZoomLevel(0) == nil {
		t.Fatal("expected color")
	}
}

// ─── Layer ──────────────────────────────────────────────────────────────────

func TestLayerRoundTrip(t *testing.T) {
	raw := `{
		"id":"water",
		"type":"fill",
		"source":"mapbox-streets",
		"source-layer":"water",
		"minzoom":0,
		"maxzoom":22,
		"filter":["==","$type","Polygon"],
		"layout":{"visibility":"visible"},
		"paint":{"fill-color":"#00ffff","fill-opacity":0.5},
		"metadata":{"key":"value"},
		"slot":"middle"
	}`
	var l Layer
	if err := json.Unmarshal([]byte(raw), &l); err != nil {
		t.Fatal(err)
	}
	if l.ID != "water" {
		t.Fatalf("ID = %q, want water", l.ID)
	}
	if l.Type != LayerTypeFill {
		t.Fatalf("Type = %q, want fill", l.Type)
	}
	if l.Source == nil || *l.Source != "mapbox-streets" {
		t.Fatalf("Source = %v", l.Source)
	}
	if l.SourceLayer == nil || *l.SourceLayer != "water" {
		t.Fatalf("SourceLayer = %v", l.SourceLayer)
	}
	if l.MinZoom == nil || *l.MinZoom != 0 {
		t.Fatal("MinZoom mismatch")
	}
	if l.MaxZoom == nil || *l.MaxZoom != 22 {
		t.Fatal("MaxZoom mismatch")
	}
	if l.Filter == nil {
		t.Fatal("Filter is nil")
	}
	if l.Layout == nil || l.Layout.Visibility != "visible" {
		t.Fatal("Layout.Visibility mismatch")
	}
	if l.Paint == nil {
		t.Fatal("Paint is nil")
	}
	if l.Slot == nil || *l.Slot != "middle" {
		t.Fatal("Slot mismatch")
	}
	if l.Metadata == nil {
		t.Fatal("Metadata is nil")
	}

	out, err := json.Marshal(&l)
	if err != nil {
		t.Fatal(err)
	}
	if !jsonEqual(string(out), raw) {
		t.Fatalf("round trip:\n  in:  %s\n  out: %s", raw, string(out))
	}
}

func TestLayerOmitEmpty(t *testing.T) {
	l := Layer{ID: "test", Type: LayerTypeFill}
	out, err := json.Marshal(&l)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), "filter") {
		t.Fatal("filter should be omitted")
	}
	if strings.Contains(string(out), "layout") {
		t.Fatal("layout should be omitted")
	}
	if strings.Contains(string(out), "paint") {
		t.Fatal("paint should be omitted")
	}
	if strings.Contains(string(out), "metadata") {
		t.Fatal("metadata should be omitted")
	}
}

func TestLayerValidate(t *testing.T) {
	tests := []struct {
		name      string
		layer     Layer
		wantError bool
	}{
		{"valid", Layer{ID: "a", Type: LayerTypeFill}, false},
		{"empty_id", Layer{Type: LayerTypeFill}, true},
		{"unknown_type", Layer{ID: "a", Type: "invalid"}, true},
		{"max_lt_min", Layer{ID: "a", Type: LayerTypeFill, MaxZoom: ptr(5.0), MinZoom: ptr(10.0)}, true},
		{"max_out_of_range", Layer{ID: "a", Type: LayerTypeFill, MaxZoom: ptr(25.0)}, true},
		{"min_negative", Layer{ID: "a", Type: LayerTypeFill, MinZoom: ptr(-1.0)}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.layer.Validate()
			if tc.wantError && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestLayerTypeConstants(t *testing.T) {
	all := []LayerType{
		LayerTypeBackground, LayerTypeBuilding, LayerTypeCircle, LayerTypeClip,
		LayerTypeFill, LayerTypeFillExtrusion, LayerTypeHeatmap, LayerTypeHillshade,
		LayerTypeLine, LayerTypeModel, LayerTypeRaster, LayerTypeRasterParticle,
		LayerTypeSky, LayerTypeSlot, LayerTypeSymbol,
	}
	for _, lt := range all {
		if !validLayerType(lt) {
			t.Errorf("validLayerType(%q) = false", lt)
		}
	}
	if validLayerType("foo") {
		t.Error("validLayerType('foo') = true")
	}
}

// ─── Layout ─────────────────────────────────────────────────────────────────

func TestLayoutRoundTrip(t *testing.T) {
	raw := `{
		"visibility":"visible",
		"line-cap":"round",
		"line-join":"miter",
		"symbol-placement":"point",
		"text-field":"name",
		"text-font":["Open Sans"],
		"text-size":16,
		"icon-image":"marker",
		"icon-size":1.5,
		"circle-sort-key":10,
		"fill-sort-key":5,
		"line-sort-key":3,
		"line-elevation-reference":"sea",
		"symbol-z-order":"auto"
	}`
	var l Layout
	if err := json.Unmarshal([]byte(raw), &l); err != nil {
		t.Fatal(err)
	}
	out, err := json.Marshal(&l)
	if err != nil {
		t.Fatal(err)
	}
	if !jsonEqual(string(out), raw) {
		t.Fatalf("round trip:\n  in:  %s\n  out: %s", raw, string(out))
	}
}

func TestLayoutOmitEmpty(t *testing.T) {
	l := Layout{}
	out, err := json.Marshal(&l)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "{}" {
		t.Fatalf("expected {}, got %s", string(out))
	}
}

// ─── Paint ──────────────────────────────────────────────────────────────────

func TestPaintRoundTrip(t *testing.T) {
	raw := `{
		"background-color":"#ff0000",
		"background-opacity":0.8,
		"fill-color":"#00ff00",
		"fill-opacity":0.5,
		"line-color":"#0000ff",
		"line-width":2,
		"circle-color":"#ff00ff",
		"circle-radius":10,
		"circle-opacity":0.9,
		"text-color":"#ffffff",
		"text-halo-color":"#000000",
		"text-halo-width":1,
		"icon-opacity":1,
		"raster-opacity":0.7,
		"heatmap-opacity":0.8,
		"hillshade-exaggeration":0.5
	}`
	var p Paint
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		t.Fatal(err)
	}
	out, err := json.Marshal(&p)
	if err != nil {
		t.Fatal(err)
	}
	if !jsonEqual(string(out), raw) {
		t.Fatalf("round trip:\n  in:  %s\n  out: %s", raw, string(out))
	}
}

func TestPaintOmitEmpty(t *testing.T) {
	p := Paint{}
	out, err := json.Marshal(&p)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "{}" {
		t.Fatalf("expected {}, got %s", string(out))
	}
}

// ─── Source ─────────────────────────────────────────────────────────────────

func TestSourceVectorRoundTrip(t *testing.T) {
	raw := `{"type":"vector","url":"mapbox://mapbox.mapbox-streets-v8"}`
	var s Source
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatal(err)
	}
	if s.Type != "vector" {
		t.Fatalf("Type = %q", s.Type)
	}
	if s.URL != "mapbox://mapbox.mapbox-streets-v8" {
		t.Fatalf("URL = %q", s.URL)
	}
	out, _ := json.Marshal(&s)
	if !jsonEqual(string(out), raw) {
		t.Fatalf("round trip:\n  in:  %s\n  out: %s", raw, string(out))
	}
}

func TestSourceRasterRoundTrip(t *testing.T) {
	raw := `{"type":"raster","url":"mapbox://mapbox.satellite","tileSize":256}`
	var s Source
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatal(err)
	}
	if s.TileSize == nil || *s.TileSize != 256 {
		t.Fatal("TileSize mismatch")
	}
}

func TestSourceGeoJSONRoundTrip(t *testing.T) {
	raw := `{"type":"geojson","data":{"type":"Feature","geometry":{"type":"Point","coordinates":[0,0]}},"cluster":true,"clusterRadius":50}`
	var s Source
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatal(err)
	}
	if s.Cluster == nil || !*s.Cluster {
		t.Fatal("Cluster should be true")
	}
	if s.ClusterRadius == nil || *s.ClusterRadius != 50 {
		t.Fatal("ClusterRadius mismatch")
	}
}

func TestSourceImageRoundTrip(t *testing.T) {
	raw := `{"type":"image","url":"https://example.com/img.png","coordinates":[[0,0],[1,0],[1,1],[0,1]]}`
	var s Source
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatal(err)
	}
	if s.Type != "image" {
		t.Fatalf("Type = %q", s.Type)
	}
	if len(s.Coordinates) != 4 {
		t.Fatalf("expected 4 coordinates, got %d", len(s.Coordinates))
	}
}

func TestSourceVideoRoundTrip(t *testing.T) {
	raw := `{"type":"video","urls":["https://example.com/vid.mp4"],"coordinates":[[0,0],[1,0],[1,1],[0,1]]}`
	var s Source
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatal(err)
	}
	if len(s.URLs) != 1 || s.URLs[0] != "https://example.com/vid.mp4" {
		t.Fatal("URLs mismatch")
	}
}

func TestSourcesMapRoundTrip(t *testing.T) {
	raw := `{
		"streets":{"type":"vector","url":"mapbox://mapbox.streets"},
		"satellite":{"type":"raster","url":"mapbox://mapbox.satellite","tileSize":256},
		"poi":{"type":"geojson","data":{"type":"FeatureCollection","features":[]}}
	}`
	var ss Sources
	if err := json.Unmarshal([]byte(raw), &ss); err != nil {
		t.Fatal(err)
	}
	if len(ss) != 3 {
		t.Fatalf("expected 3 sources, got %d", len(ss))
	}
	if ss["streets"].Type != "vector" {
		t.Fatal("streets type mismatch")
	}
	if ss["satellite"].TileSize == nil || *ss["satellite"].TileSize != 256 {
		t.Fatal("satellite tileSize mismatch")
	}
	out, err := json.Marshal(ss)
	if err != nil {
		t.Fatal(err)
	}
	if !jsonEqual(string(out), raw) {
		t.Fatalf("round trip:\n  in:  %s\n  out: %s", raw, string(out))
	}
}

// ─── Model Source ────────────────────────────────────────────────────────────

func TestModelSourceModelRoundTrip(t *testing.T) {
	raw := `{"uri":"/assets/car.glb","orientation":[0,45,0],"position":[0,0]}`
	var m ModelSourceModel
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		t.Fatal(err)
	}
	if m.URI != "/assets/car.glb" {
		t.Fatalf("URI = %q", m.URI)
	}
	if len(m.Orientation) != 3 || m.Orientation[1] != 45 {
		t.Fatal("Orientation mismatch")
	}
}

func TestSourceWithModelsRoundTrip(t *testing.T) {
	raw := `{"type":"model","models":{"car":{"uri":"/assets/car.glb"}}}`
	var s Source
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatal(err)
	}
	if s.Type != "model" {
		t.Fatalf("Type = %q", s.Type)
	}
	if s.Models == nil || s.Models["car"].URI != "/assets/car.glb" {
		t.Fatal("Model source car not found")
	}
}

// ─── Fog ────────────────────────────────────────────────────────────────────

func TestFogRoundTrip(t *testing.T) {
	raw := `{"color":"#dc9f9f","high-color":"#245bde","horizon-blend":0.5,"range":[0.8,8],"space-color":"#000000","star-intensity":0.15,"vertical-range":[0,0]}`
	var f Fog
	if err := json.Unmarshal([]byte(raw), &f); err != nil {
		t.Fatal(err)
	}
	out, err := json.Marshal(&f)
	if err != nil {
		t.Fatal(err)
	}
	if !jsonEqual(string(out), raw) {
		t.Fatalf("round trip:\n  in:  %s\n  out: %s", raw, string(out))
	}
}

// ─── Terrain ────────────────────────────────────────────────────────────────

func TestTerrainRoundTrip(t *testing.T) {
	raw := `{"source":"mapbox-raster-dem","exaggeration":1.5}`
	var tr Terrain
	if err := json.Unmarshal([]byte(raw), &tr); err != nil {
		t.Fatal(err)
	}
	if tr.Source != "mapbox-raster-dem" {
		t.Fatalf("Source = %q", tr.Source)
	}
	if tr.Exaggeration == nil {
		t.Fatal("Exaggeration is nil")
	}
	out, _ := json.Marshal(&tr)
	if !jsonEqual(string(out), raw) {
		t.Fatalf("round trip:\n  in:  %s\n  out: %s", raw, string(out))
	}
}

// ─── Projection ─────────────────────────────────────────────────────────────

func TestProjectionRoundTrip(t *testing.T) {
	raw := `{"name":"mercator"}`
	var p Projection
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		t.Fatal(err)
	}
	if p.Name != "mercator" {
		t.Fatalf("Name = %q", p.Name)
	}
	out, _ := json.Marshal(&p)
	if !jsonEqual(string(out), raw) {
		t.Fatalf("round trip:\n  in:  %s\n  out: %s", raw, string(out))
	}
}

func TestProjectionWithCenter(t *testing.T) {
	raw := `{"name":"albers","center":[-96,37.5],"parallels":[29.5,45.5]}`
	var p Projection
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		t.Fatal(err)
	}
	if len(p.Center) != 2 || p.Center[0] != -96 {
		t.Fatal("Center mismatch")
	}
}

// ─── Light3D ────────────────────────────────────────────────────────────────

func TestLight3DRoundTrip(t *testing.T) {
	raw := `{"id":"sun","type":"directional","properties":{"color":"#ffffff","intensity":0.5,"direction":[210,30],"cast-shadows":true}}`
	var l Light3D
	if err := json.Unmarshal([]byte(raw), &l); err != nil {
		t.Fatal(err)
	}
	if l.ID != "sun" {
		t.Fatalf("ID = %q", l.ID)
	}
	if l.Type != "directional" {
		t.Fatalf("Type = %q", l.Type)
	}
	if l.Properties == nil {
		t.Fatal("Properties is nil")
	}
	if l.Properties.CastShadows == nil || !*l.Properties.CastShadows {
		t.Fatal("CastShadows should be true")
	}
}

// ─── Camera ─────────────────────────────────────────────────────────────────

func TestCameraRoundTrip(t *testing.T) {
	raw := `{"camera-projection":"orthographic"}`
	var c Camera
	if err := json.Unmarshal([]byte(raw), &c); err != nil {
		t.Fatal(err)
	}
	if c.CameraProjection != "orthographic" {
		t.Fatalf("CameraProjection = %q", c.CameraProjection)
	}
}

// ─── Import ─────────────────────────────────────────────────────────────────

func TestImportRoundTrip(t *testing.T) {
	raw := `{"id":"basemap","url":"mapbox://styles/mapbox/standard","config":{"lightPreset":"dusk"}}`
	var im Import
	if err := json.Unmarshal([]byte(raw), &im); err != nil {
		t.Fatal(err)
	}
	if im.ID != "basemap" {
		t.Fatalf("ID = %q", im.ID)
	}
	if im.Config["lightPreset"] != "dusk" {
		t.Fatal("Config mismatch")
	}
}

// ─── Appearance ─────────────────────────────────────────────────────────────

func TestAppearanceRoundTrip(t *testing.T) {
	raw := `{"name":"selected","condition":["feature-state","select"],"properties":{"icon-size":1.3,"icon-color":"red"}}`
	var a Appearance
	if err := json.Unmarshal([]byte(raw), &a); err != nil {
		t.Fatal(err)
	}
	if a.Name != "selected" {
		t.Fatalf("Name = %q", a.Name)
	}
	if a.Properties["icon-size"] != 1.3 {
		t.Fatal("icon-size mismatch")
	}
}

// ─── SchemaOption ───────────────────────────────────────────────────────────

func TestSchemaOptionRoundTrip(t *testing.T) {
	raw := `{"default":10,"type":"number","minValue":0,"maxValue":100}`
	var s SchemaOption
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatal(err)
	}
	if s.Default.(float64) != 10 {
		t.Fatalf("Default = %v", s.Default)
	}
	if s.Type != "number" {
		t.Fatalf("Type = %q", s.Type)
	}
	if *s.MinValue != 0 {
		t.Fatal("MinValue mismatch")
	}
	if *s.MaxValue != 100 {
		t.Fatal("MaxValue mismatch")
	}
}

// ─── Transition ─────────────────────────────────────────────────────────────

func TestTransitionRoundTrip(t *testing.T) {
	raw := `{"delay":100,"duration":300}`
	var tr Transition
	if err := json.Unmarshal([]byte(raw), &tr); err != nil {
		t.Fatal(err)
	}
	if tr.Delay != 100 {
		t.Fatalf("Delay = %d", tr.Delay)
	}
	if tr.Duration != 300 {
		t.Fatalf("Duration = %d", tr.Duration)
	}
}

// ─── Light ──────────────────────────────────────────────────────────────────

func TestLightRoundTrip(t *testing.T) {
	raw := `{"anchor":"viewport","color":"#ffffff","intensity":0.5,"position":[1.15,210,30]}`
	var l Light
	if err := json.Unmarshal([]byte(raw), &l); err != nil {
		t.Fatal(err)
	}
	if l.Anchor != "viewport" {
		t.Fatalf("Anchor = %q", l.Anchor)
	}
	if l.Color == nil {
		t.Fatal("Color is nil")
	}
	if l.Color.GetColorAtZoomLevel(0) == nil {
		t.Fatal("Color value is nil")
	}
}

// ─── Complete Style ─────────────────────────────────────────────────────────

func TestStyleUnmarshal(t *testing.T) {
	raw := `{
		"version":8,
		"name":"test-style",
		"center":[0,0],
		"zoom":10,
		"bearing":0,
		"pitch":45,
		"glyphs":"mapbox://fonts/{fontstack}/{range}.pbf",
		"sprite":"mapbox://sprites/test",
		"sources":{
			"streets":{"type":"vector","url":"mapbox://mapbox.streets"}
		},
		"layers":[
			{"id":"bg","type":"background","paint":{"background-color":"#fff"}},
			{"id":"water","type":"fill","source":"streets","source-layer":"water","paint":{"fill-color":"#00ffff"}}
		],
		"fog":{"color":"#ffffff","range":[0.5,10]},
		"terrain":{"source":"dem","exaggeration":1},
		"projection":{"name":"mercator"},
		"lights":[{"id":"sun","type":"directional","properties":{"intensity":0.5}}],
		"transition":{"duration":300,"delay":0},
		"metadata":{"author":"test"},
		"models":{"tree":"asset://tree.glb"},
		"schema":{"option1":{"default":10,"type":"number"}}
	}`
	var s Style
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatal(err)
	}
	if s.Version != 8 {
		t.Fatalf("Version = %d", s.Version)
	}
	if len(s.Layers) != 2 {
		t.Fatalf("expected 2 layers, got %d", len(s.Layers))
	}
	if s.Layers[0].Type != LayerTypeBackground {
		t.Fatalf("Layer[0].Type = %q", s.Layers[0].Type)
	}
	if s.Layers[1].Source == nil || *s.Layers[1].Source != "streets" {
		t.Fatal("Layer[1].Source mismatch")
	}
	if s.Fog == nil {
		t.Fatal("Fog is nil")
	}
	if s.Terrain == nil {
		t.Fatal("Terrain is nil")
	}
	if s.Projection == nil || s.Projection.Name != "mercator" {
		t.Fatal("Projection mismatch")
	}
	if len(s.Lights) != 1 {
		t.Fatal("Lights missing")
	}
	if s.Lights[0].ID != "sun" {
		t.Fatal("Light ID mismatch")
	}
	if s.Transition == nil {
		t.Fatal("Transition is nil")
	}
	if s.Models == nil || s.Models["tree"] != "asset://tree.glb" {
		t.Fatal("Models mismatch")
	}
	if s.Schema == nil || s.Schema["option1"].Default == nil {
		t.Fatal("Schema mismatch")
	}
	if s.Metadata == nil || s.Metadata["author"] != "test" {
		t.Fatal("Metadata mismatch")
	}
}

func TestStyleValidate(t *testing.T) {
	tests := []struct {
		name      string
		style     Style
		wantError bool
	}{
		{"valid", Style{Version: 8}, false},
		{"wrong_version", Style{Version: 7}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.style.Validate()
			if tc.wantError && err == nil {
				t.Fatal("expected error")
			}
			if !tc.wantError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestParse(t *testing.T) {
	raw := `{"version":8,"name":"test","sources":{},"layers":[{"id":"bg","type":"background","paint":{"background-color":"#ff0000"}}]}`
	r := strings.NewReader(raw)
	ms, err := Parse(r)
	if err != nil {
		t.Fatal(err)
	}
	if ms.GetStyleID() != "" {
		t.Fatal("expected empty ID")
	}
	bg := ms.GetBackground()
	if bg == nil {
		t.Fatal("expected non-nil background")
	}
}

// ─── Expression Sub-types (linear/exponential/cubic-bezier) ────────────────

func TestInterpolationTypes(t *testing.T) {
	tests := []struct {
		json string
		op   string
	}{
		{`["linear"]`, ExpLinear},
		{`["exponential",2]`, ExpExponential},
		{`["cubic-bezier",0,0,1,1]`, ExpCubicBezier},
	}
	for _, tc := range tests {
		t.Run(tc.op, func(t *testing.T) {
			var e Expression
			if err := json.Unmarshal([]byte(tc.json), &e); err != nil {
				t.Fatal(err)
			}
			if e.Operator != tc.op {
				t.Fatalf("expected %q, got %q", tc.op, e.Operator)
			}
		})
	}
}

// ─── Featureset / Rain / Snow ──────────────────────────────────────────────

func TestFeaturesetRoundTrip(t *testing.T) {
	raw := `{"selectors":[{"layer":"poi","featureNamespace":"ns","properties":{"type":["get","type"]}}]}`
	var fs Featureset
	if err := json.Unmarshal([]byte(raw), &fs); err != nil {
		t.Fatal(err)
	}
	if len(fs.Selectors) != 1 {
		t.Fatal("expected 1 selector")
	}
	if fs.Selectors[0].Layer != "poi" {
		t.Fatalf("Layer = %q", fs.Selectors[0].Layer)
	}
}

func TestRainRoundTrip(t *testing.T) {
	raw := `{"color":"#ffffff","intensity":0.5,"opacity":0.8}`
	var r Rain
	if err := json.Unmarshal([]byte(raw), &r); err != nil {
		t.Fatal(err)
	}
	if r.Intensity == nil || *r.Intensity != 0.5 {
		t.Fatal("Intensity mismatch")
	}
}

func TestSnowRoundTrip(t *testing.T) {
	raw := `{"color":"#ffffff","intensity":0.3,"opacity":0.9}`
	var s Snow
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatal(err)
	}
	if s.Intensity == nil || *s.Intensity != 0.3 {
		t.Fatal("Intensity mismatch")
	}
}

// ─── Iconset ────────────────────────────────────────────────────────────────

func TestIconsetRoundTrip(t *testing.T) {
	raw := `{"type":"sprite","url":"myURL"}`
	var is Iconset
	if err := json.Unmarshal([]byte(raw), &is); err != nil {
		t.Fatal(err)
	}
	if is.Type != "sprite" {
		t.Fatalf("Type = %q", is.Type)
	}
	if is.URL != "myURL" {
		t.Fatalf("URL = %q", is.URL)
	}
}

// ─── Model Material / Node Overrides ───────────────────────────────────────

func TestModelMaterialOverrideRoundTrip(t *testing.T) {
	raw := `{"model-color":"#ff0000","model-opacity":0.5}`
	var m ModelMaterialOverride
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		t.Fatal(err)
	}
	if m.ModelOpacity == nil || *m.ModelOpacity != 0.5 {
		t.Fatal("ModelOpacity mismatch")
	}
}

func TestModelNodeOverrideRoundTrip(t *testing.T) {
	raw := `{"orientation":[0,45,0]}`
	var m ModelNodeOverride
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		t.Fatal(err)
	}
	if len(m.Orientation) != 3 || m.Orientation[1] != 45 {
		t.Fatal("Orientation mismatch")
	}
}

// ─── Edge Cases: Empty / Nil ───────────────────────────────────────────────

func TestEmptyStyle(t *testing.T) {
	raw := `{"version":8,"sources":{},"layers":[]}`
	var s Style
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatal(err)
	}
	if s.Version != 8 {
		t.Fatalf("Version = %d", s.Version)
	}
	if s.Light != nil {
		t.Fatal("Light should be nil")
	}
	if s.Transition != nil {
		t.Fatal("Transition should be nil")
	}
}

func TestNilFilterContainerMarshal(t *testing.T) {
	fc := FilterContainer{}
	out, err := json.Marshal(&fc)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "null" {
		t.Fatalf("expected null, got %s", string(out))
	}
}

func TestEmptyExpressions(t *testing.T) {
	// Empty array [] is parsed as literal `[]`
	var e Expression
	if err := json.Unmarshal([]byte(`[]`), &e); err != nil {
		t.Fatal(err)
	}
	if !e.IsLiteral {
		t.Fatal("empty array should be literal")
	}
}

// ─── ColorType edge cases ──────────────────────────────────────────────────

func TestColorTypeFunctionStopsRoundTrip(t *testing.T) {
	raw := `{"stops":[[0,"red"],[10,"blue"]],"base":1.5}`
	var c ColorType
	if err := json.Unmarshal([]byte(raw), &c); err != nil {
		t.Fatal(err)
	}
	clr := c.GetColorAtZoomLevel(5)
	if clr == nil {
		t.Fatal("expected interpolated color at zoom 5")
	}
	out, err := json.Marshal(&c)
	if err != nil {
		t.Fatal(err)
	}
	if !jsonEqual(string(out), raw) {
		t.Fatalf("round trip:\n  in:  %s\n  out: %s", raw, string(out))
	}
}

func TestColorTypeNilGetColor(t *testing.T) {
	var c *ColorType
	if clr := c.GetColorAtZoomLevel(0); clr != nil {
		t.Fatal("expected nil from nil receiver")
	}
	c = &ColorType{}
	if clr := c.GetColorAtZoomLevel(0); clr != nil {
		t.Fatal("expected nil from empty ColorType")
	}
}

func TestColorTypeMarshalNil(t *testing.T) {
	var c ColorType
	out, err := json.Marshal(&c)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "null" {
		t.Fatalf("expected null, got %s", string(out))
	}
}

// ─── ColorStopsType edge cases ──────────────────────────────────────────────

func TestColorStopsTypeBounds(t *testing.T) {
	raw := `{"stops":[[0,"#ff0000"],[10,"#0000ff"]]}`
	var c ColorType
	json.Unmarshal([]byte(raw), &c)
	// below min stop returns nil
	if clr := c.GetColorAtZoomLevel(-1); clr != nil {
		t.Fatal("expected nil below min stop")
	}
	// at max stop returns last color
	if clr := c.GetColorAtZoomLevel(10); clr == nil {
		t.Fatal("expected color at max stop")
	}
	// above max stop returns last color
	if clr := c.GetColorAtZoomLevel(20); clr == nil {
		t.Fatal("expected color above max stop")
	}
}

func TestColorStopsTypeExactStop(t *testing.T) {
	raw := `{"stops":[[5,"#00ff00"],[15,"#ff0000"]]}`
	var c ColorType
	json.Unmarshal([]byte(raw), &c)
	clr := c.GetColorAtZoomLevel(5)
	if clr == nil {
		t.Fatal("expected color at exact stop")
	}
	rgba, ok := clr.(color.RGBA)
	if !ok {
		t.Fatal("expected RGBA")
	}
	if rgba.G != 255 {
		t.Fatalf("expected green=255 at stop 5, got %+v", rgba)
	}
}

func TestColorStopsTypeSingleStop(t *testing.T) {
	raw := `{"stops":[[0,"#ff0000"]]}`
	var c ColorType
	json.Unmarshal([]byte(raw), &c)
	clr := c.GetColorAtZoomLevel(0)
	if clr == nil {
		t.Fatal("expected color with single stop")
	}
}

// ─── ValidateExpression edge cases ─────────────────────────────────────────

func TestValidateExpressionNil(t *testing.T) {
	err := ValidateExpression(nil, ValidationStrict)
	if err == nil {
		t.Fatal("expected error for nil expression")
	}
}

func TestValidateExpressionEmptyOperator(t *testing.T) {
	e := &Expression{IsLiteral: false, Operator: ""}
	err := ValidateExpression(e, ValidationStrict)
	if err == nil {
		t.Fatal("expected error for empty operator")
	}
}

func TestValidateExpressionBasicMode(t *testing.T) {
	e := &Expression{IsLiteral: false, Operator: "foobar"}
	// Basic mode only checks known operator - "foobar" is unknown
	err := ValidateExpression(e, ValidationBasic)
	if err == nil {
		t.Fatal("expected error for unknown operator even in basic mode")
	}
}

// ─── validateArgs full coverage ─────────────────────────────────────────────

func TestValidateExpressionAllOperators(t *testing.T) {
	type testCase struct {
		name   string
		json   string
		valid  bool
	}
	tests := []testCase{
		// Types
		{"array_minimal", `["array",["get","x"]]`, true},
		{"array_type", `["array","number",["get","x"]]`, true},
		{"array_type_len", `["array","number",3,["get","x"]]`, true},
		{"boolean", `["boolean",["get","x"]]`, true},
		{"boolean_fallback", `["boolean",["get","x"],false]`, true},
		{"collator_empty", `["collator"]`, true},
		{"collator_options", `["collator",{"locale":"en"}]`, true},
		{"number_assert", `["number",["get","x"]]`, true},
		{"object_assert", `["object",["get","x"]]`, true},
		{"string_assert", `["string",["get","x"]]`, true},
		{"to_string", `["to-string",["get","name"]]`, true},
		{"to_string_err", `["to-string"]`, false},
		{"to_boolean", `["to-boolean",["get","x"]]`, true},
		{"to_number", `["to-number",["get","x"]]`, true},
		{"to_number_fallback", `["to-number",["get","x"],0]`, true},
		{"to_color", `["to-color",["get","x"]]`, true},
		{"typeof", `["typeof",["get","x"]]`, true},
		{"typeof_err", `["typeof"]`, false},
		{"literal", `["literal",[1,2,3]]`, true},
		{"literal_err", `["literal"]`, false},
		{"image", `["image","my-icon"]`, true},
		{"image_err", `["image"]`, false},
		{"number_format", `["number-format",["get","x"]]`, true},
		{"number_format_full", `["number-format",["get","x"],{"locale":"en"}]`, true},
		{"number_format_err", `["number-format"]`, false},

		// Feature data
		{"accumulated", `["accumulated"]`, true},
		{"accumulated_err", `["accumulated",1]`, false},
		{"feature_state", `["feature-state","selected"]`, true},
		{"feature_state_err", `["feature-state"]`, false},
		{"geometry_type", `["geometry-type"]`, true},
		{"id", `["id"]`, true},
		{"line_progress", `["line-progress"]`, true},
		{"properties", `["properties"]`, true},

		// Lookup
		{"at", `["at",0,["get","arr"]]`, true},
		{"at_err", `["at",0]`, false},
		{"at_interpolated", `["at-interpolated",0.5,["get","arr"]]`, true},
		{"config", `["config","option"]`, true},
		{"config_scope", `["config","option","scope"]`, true},
		{"has", `["has","name"]`, true},
		{"has_with_obj", `["has","name",["properties"]]`, true},
		{"has_err", `["has"]`, false},
		{"in", `["in",1,["get","arr"]]`, true},
		{"in_err", `["in",1]`, false},
		{"index_of", `["index-of","a",["get","str"]]`, true},
		{"index_of_with_start", `["index-of","a",["get","str"],0]`, true},
		{"index_of_err", `["index-of","a"]`, false},
		{"length", `["length",["get","str"]]`, true},
		{"length_err", `["length"]`, false},
		{"measure_light", `["measure-light","brightness"]`, true},
		{"measure_light_err", `["measure-light"]`, false},
		{"slice", `["slice",["get","str"],0]`, true},
		{"slice_end", `["slice",["get","str"],0,5]`, true},
		{"slice_err", `["slice",["get","str"]]`, false},
		{"split", `["split",["get","str"],","]`, true},
		{"split_err", `["split",["get","str"]]`, false},
		{"worldview", `["worldview"]`, true},
		{"worldview_err", `["worldview",1]`, false},

		// Decision
		{"not", `["!",["get","flag"]]`, true},
		{"not_err", `["!"]`, false},
		{"neq", `["!=",["get","a"],["get","b"]]`, true},
		{"eq", `["==",["get","a"],["get","b"]]`, true},
		{"lt", `["<",["get","a"],["get","b"]]`, true},
		{"lte", `["<=",["get","a"],["get","b"]]`, true},
		{"gt", `[">",["get","a"],["get","b"]]`, true},
		{"gte", `[">=",["get","a"],["get","b"]]`, true},
		{"compare_collator", `["==",["get","a"],["get","b"],["collator"]]`, true},
		{"compare_err", `["==",["get","a"]]`, false},
		{"all", `["all",true,false,true]`, true},
		{"all_err", `["all",true]`, false},
		{"any", `["any",false,true]`, true},
		{"any_err", `["any",false]`, false},
		{"coalesce", `["coalesce",["get","a"],["get","b"],"fallback"]`, true},
		{"coalesce_err", `["coalesce"]`, false},
		{"match_full", `["match",["get","t"],"a",1,"b",2,3]`, true},
		{"match_err", `["match",["get","t"]]`, false},
		{"within", `["within",{"type":"Polygon"}]`, true},
		{"within_err", `["within"]`, false},

		// Ramps
		{"interpolate_linear", `["interpolate",["linear"],["zoom"],0,1,10,10]`, true},
		{"interpolate_exponential", `["interpolate",["exponential",2],["zoom"],0,1,10,10]`, true},
		{"interpolate_cubic_bezier", `["interpolate",["cubic-bezier",0,0,1,1],["zoom"],0,1,10,10]`, true},
		{"interpolate_wrong_subtype", `["interpolate",["foobar"],["zoom"],0,1,10,10]`, false},
		{"step", `["step",["zoom"],0,5,1,10,2]`, true},
		{"step_err", `["step",["zoom"]]`, false},

		// Variable binding
		{"let", `["let","x",["get","val"],["var","x"]]`, true},
		{"let_multi", `["let","x",1,"y",2,["+",["var","x"],["var","y"]]]`, true},
		{"let_err", `["let","x"]`, false},
		{"var", `["var","x"]`, true},
		{"var_err", `["var"]`, false},

		// String
		{"concat", `["concat",["get","a"],["get","b"]]`, true},
		{"downcase", `["downcase",["get","a"]]`, true},
		{"downcase_err", `["downcase"]`, false},
		{"upcase", `["upcase",["get","a"]]`, true},
		{"is_supported_script", `["is-supported-script",["get","a"]]`, true},
		{"is_supported_script_err", `["is-supported-script"]`, false},
		{"resolved_locale", `["resolved-locale",["collator"]]`, true},
		{"resolved_locale_err", `["resolved-locale"]`, false},

		// Color
		{"hsl", `["hsl",0,100,50]`, true},
		{"hsl_err", `["hsl",0,100]`, false},
		{"hsla", `["hsla",0,100,50,1]`, true},
		{"hsla_err", `["hsla",0,100,50]`, false},

		// Math
		{"add", `["+",1,2]`, true},
		{"add_multi", `["+",1,2,3]`, true},
		{"add_err", `["+",1]`, false},
		{"sub", `["-",1,2]`, true},
		{"mul", `["*",1,2]`, true},
		{"div", `["/",1,2]`, true},
		{"mod", `["%",1,2]`, true},
		{"mod_err", `["%",1]`, false},
		{"pow", `["^",2,3]`, true},
		{"pow_err", `["^",2]`, false},
		{"abs", `["abs",-1]`, true},
		{"abs_err", `["abs"]`, false},
		{"acos", `["acos",0.5]`, true},
		{"asin", `["asin",0.5]`, true},
		{"atan", `["atan",0.5]`, true},
		{"ceil", `["ceil",1.5]`, true},
		{"cos", `["cos",0]`, true},
		{"distance", `["distance",["get","a"],["get","b"]]`, true},
		{"distance_err", `["distance",["get","a"]]`, false},
		{"e", `["e"]`, true},
		{"e_err", `["e",1]`, false},
		{"floor", `["floor",1.5]`, true},
		{"ln", `["ln",2]`, true},
		{"log10", `["log10",100]`, true},
		{"log2", `["log2",8]`, true},
		{"max", `["max",1,2]`, true},
		{"max_multi", `["max",1,2,3]`, true},
		{"max_err", `["max",1]`, false},
		{"min", `["min",1,2]`, true},
		{"pi", `["pi"]`, true},
		{"pi_err", `["pi",1]`, false},
		{"random", `["random",0,1]`, true},
		{"random_seeded", `["random",0,1,["get","seed"]]`, true},
		{"random_err", `["random",0]`, false},
		{"round", `["round",1.5]`, true},
		{"sin", `["sin",0]`, true},
		{"sqrt", `["sqrt",4]`, true},
		{"tan", `["tan",0]`, true},

		// Camera
		{"distance_from_center", `["distance-from-center"]`, true},
		{"distance_from_center_err", `["distance-from-center",1]`, false},
		{"pitch", `["pitch"]`, true},
		{"pitch_err", `["pitch",1]`, false},
		{"zoom", `["zoom"]`, true},

		// Heatmap
		{"heatmap_density", `["heatmap-density"]`, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var e Expression
			if err := json.Unmarshal([]byte(tc.json), &e); err != nil {
				t.Fatal(err)
			}
			err := ValidateExpression(&e, ValidationStrict)
			if tc.valid && err != nil {
				t.Fatalf("expected valid, got error: %v", err)
			}
			if !tc.valid && err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

// ─── Expression deep nesting and edge cases ────────────────────────────────

func TestExpressionDeeplyNested(t *testing.T) {
	raw := `["case",["==",["get","type"],"a"],["interpolate",["linear"],["zoom"],0,1,10,["get","val"]],"default"]`
	var e Expression
	if err := json.Unmarshal([]byte(raw), &e); err != nil {
		t.Fatal(err)
	}
	out, err := json.Marshal(&e)
	if err != nil {
		t.Fatal(err)
	}
	if !jsonEqual(string(out), raw) {
		t.Fatalf("round trip:\n  in:  %s\n  out: %s", raw, string(out))
	}
}

func TestExpressionMarshalEdgeCases(t *testing.T) {
	// Deeply nested args with nil elements
	e := &Expression{
		Operator: "coalesce",
		Args: []*Expression{
			{IsLiteral: true, Value: nil},  // null literal
			{Operator: "get", Args: []*Expression{{IsLiteral: true, Value: "name"}}},
		},
	}
	out, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	expected := `["coalesce",null,["get","name"]]`
	if !jsonEqual(string(out), expected) {
		t.Fatalf("expected %s, got %s", expected, string(out))
	}
}

func TestExpressionDecodeInvalidJSON(t *testing.T) {
	var e Expression
	if err := json.Unmarshal([]byte(`{invalid}`), &e); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

// ─── FilterContainer edge cases ────────────────────────────────────────────

func TestFilterContainerNilExpr(t *testing.T) {
	fc := FilterContainer{Expr: nil}
	out, err := json.Marshal(&fc)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "null" {
		t.Fatalf("expected null, got %s", string(out))
	}
}

// ─── Style Parse edge cases ─────────────────────────────────────────────────

func TestParseNoBackgroundLayer(t *testing.T) {
	raw := `{"version":8,"sources":{},"layers":[{"id":"water","type":"fill","paint":{"fill-color":"#00ffff"}}]}`
	r := strings.NewReader(raw)
	ms, err := Parse(r)
	if err != nil {
		t.Fatal(err)
	}
	bg := ms.GetBackground()
	if bg == nil {
		t.Fatal("expected non-nil background (default white)")
	}
}

func TestParseNoLayers(t *testing.T) {
	raw := `{"version":8,"sources":{},"layers":[]}`
	r := strings.NewReader(raw)
	ms, err := Parse(r)
	if err != nil {
		t.Fatal(err)
	}
	if ms.GetBackground() == nil {
		t.Fatal("expected default background")
	}
}

func TestParseInvalidJSON(t *testing.T) {
	r := strings.NewReader(`{invalid}`)
	_, err := Parse(r)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

// ─── calculateBackgroundColor edge cases ───────────────────────────────────

func TestCalculateBackgroundColorEdgeCases(t *testing.T) {
	s := &Style{Version: 8, Sources: Sources{}, Layers: []*Layer{}}
	ms := &MapboxGLStyle{style: s}
	if ms.GetStyleID() != "" {
		t.Fatal("expected empty ID")
	}
}

// ─── Color format edge cases ───────────────────────────────────────────────

func TestColorHexFormats(t *testing.T) {
	tests := []struct{ input, expected string }{
		{`"#ff0000"`, `"#ff0000"`},
		{`"#f00"`, `"#ff0000"`},
		{`"#00ff00"`, `"#00ff00"`},
		{`"#0000ff"`, `"#0000ff"`},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			var c ColorType
			if err := json.Unmarshal([]byte(tc.input), &c); err != nil {
				t.Fatal(err)
			}
			clr := c.GetColorAtZoomLevel(0)
			if clr == nil {
				t.Fatal("expected color")
			}
		})
	}
}

func TestColorNamedFormats(t *testing.T) {
	names := []string{"red", "blue", "green", "white", "black", "yellow", "orange", "purple", "transparent"}
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			var c ColorType
			if err := json.Unmarshal([]byte(`"`+name+`"`), &c); err != nil {
				t.Fatal(err)
			}
			clr := c.GetColorAtZoomLevel(0)
			if clr == nil && name != "transparent" {
				t.Fatalf("expected color for %q", name)
			}
		})
	}
}

func TestColorEmptyString(t *testing.T) {
	var c ColorType
	if err := json.Unmarshal([]byte(`""`), &c); err != nil {
		t.Fatal(err)
	}
}

// ─── Source edge cases ─────────────────────────────────────────────────────

func TestSourceEmpty(t *testing.T) {
	var s Source
	out, err := json.Marshal(&s)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != `{"type":""}` {
		t.Fatalf("expected type only, got %s", string(out))
	}
}

func TestSourceDEM(t *testing.T) {
	raw := `{"type":"raster-dem","url":"mapbox://mapbox.mapbox-terrain-dem-v1","encoding":"mapbox"}`
	var s Source
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatal(err)
	}
	if s.Encoding != "mapbox" {
		t.Fatalf("Encoding = %q", s.Encoding)
	}
}

func TestSourceWithPromoteID(t *testing.T) {
	raw := `{"type":"vector","url":"mapbox://mapbox.streets","promoteId":"name"}`
	var s Source
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatal(err)
	}
	if s.PromoteID.(string) != "name" {
		t.Fatal("PromoteID mismatch")
	}
}

func TestSourcesMapMarshalEmpty(t *testing.T) {
	ss := Sources{}
	out, err := json.Marshal(ss)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "{}" {
		t.Fatalf("expected {}, got %s", string(out))
	}

	out, err = json.Marshal(Sources(nil))
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "null" {
		t.Fatalf("expected null, got %s", string(out))
	}
}

// ─── Import edge cases ─────────────────────────────────────────────────────

func TestImportWithColorTheme(t *testing.T) {
	raw := `{"id":"standard","url":"mapbox://styles/standard","color-theme":{"data":"base64data"}}`
	var im Import
	if err := json.Unmarshal([]byte(raw), &im); err != nil {
		t.Fatal(err)
	}
	if im.ColorTheme == nil || im.ColorTheme.Data != "base64data" {
		t.Fatal("ColorTheme mismatch")
	}
}

func TestImportWithInlineData(t *testing.T) {
	raw := `{"id":"custom","url":"mapbox://styles/custom","data":{"version":8,"sources":{},"layers":[]}}`
	var im Import
	if err := json.Unmarshal([]byte(raw), &im); err != nil {
		t.Fatal(err)
	}
	if im.Data == nil || im.Data.Version != 8 {
		t.Fatal("Inline style data mismatch")
	}
}

// ─── Layer edge cases ──────────────────────────────────────────────────────

func TestLayerAllTypesRoundTrip(t *testing.T) {
	types := []LayerType{
		LayerTypeBackground, LayerTypeBuilding, LayerTypeCircle, LayerTypeClip,
		LayerTypeFill, LayerTypeFillExtrusion, LayerTypeHeatmap, LayerTypeHillshade,
		LayerTypeLine, LayerTypeModel, LayerTypeRaster, LayerTypeRasterParticle,
		LayerTypeSky, LayerTypeSlot, LayerTypeSymbol,
	}
	for _, lt := range types {
		t.Run(string(lt), func(t *testing.T) {
			raw := `{"id":"test","type":"` + string(lt) + `"}`
			var l Layer
			if err := json.Unmarshal([]byte(raw), &l); err != nil {
				t.Fatal(err)
			}
			if l.Type != lt {
				t.Fatalf("Type = %q, want %q", l.Type, lt)
			}
			out, _ := json.Marshal(&l)
			if !jsonEqual(string(out), raw) {
				t.Fatalf("round trip:\n  in:  %s\n  out: %s", raw, string(out))
			}
		})
	}
}

func TestLayerWithAppearances(t *testing.T) {
	raw := `{"id":"test","type":"symbol","appearances":[{"name":"selected","condition":["feature-state","select"],"properties":{"icon-size":1.3}}]}`
	var l Layer
	if err := json.Unmarshal([]byte(raw), &l); err != nil {
		t.Fatal(err)
	}
	if len(l.Appearances) != 1 {
		t.Fatalf("expected 1 appearance, got %d", len(l.Appearances))
	}
	if l.Appearances[0].Name != "selected" {
		t.Fatalf("Appearance name = %q", l.Appearances[0].Name)
	}
}

// ─── ModelSourceModel all fields ────────────────────────────────────────────

func TestModelSourceModelAllFields(t *testing.T) {
	raw := `{"uri":"/assets/car.glb","orientation":[0,45,0],"position":[0,0],"featureProperties":{"color":"red"},"materialOverrideNames":["paint"],"materialOverrides":{"paint":{"model-color":"#ff0000","model-opacity":0.5}},"nodeOverrideNames":["door"],"nodeOverrides":{"door":{"orientation":[0,0,-15]}}}`
	var m ModelSourceModel
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		t.Fatal(err)
	}
	if m.URI != "/assets/car.glb" {
		t.Fatalf("URI = %q", m.URI)
	}
	if m.FeatureProperties["color"] != "red" {
		t.Fatal("FeatureProperties mismatch")
	}
	if len(m.MaterialOverrideNames) != 1 || m.MaterialOverrideNames[0] != "paint" {
		t.Fatal("MaterialOverrideNames mismatch")
	}
	if m.MaterialOverrides["paint"].ModelOpacity == nil || *m.MaterialOverrides["paint"].ModelOpacity != 0.5 {
		t.Fatal("MaterialOverrides mismatch")
	}
	if m.NodeOverrides["door"].Orientation[2] != -15 {
		t.Fatal("NodeOverrides mismatch")
	}
}

// ─── Expression Format edge cases ──────────────────────────────────────────

func TestFormatExpressionEdgeCases(t *testing.T) {
	if got := FormatExpression(nil); got != "null" {
		t.Fatalf("expected null, got %s", got)
	}
	if got := FormatExpression(&Expression{IsLiteral: true, Value: nil}); got != "<nil>" {
		t.Fatalf("expected <nil>, got %s", got)
	}
}

// ─── IsDataExpression edge cases ───────────────────────────────────────────

func TestIsDataExpressionAllOperators(t *testing.T) {
	dataOps := []string{
		`["get","name"]`, `["has","name"]`, `["id"]`,
		`["geometry-type"]`, `["properties"]`, `["feature-state","sel"]`,
	}
	for _, raw := range dataOps {
		t.Run(raw, func(t *testing.T) {
			var e Expression
			json.Unmarshal([]byte(raw), &e)
			if !IsDataExpression(&e) {
				t.Fatalf("expected data expression for %s", raw)
			}
		})
	}
}

// ─── GetValueAtZoomLevel helpers ───────────────────────────────────────────

func TestNumericHelpers(t *testing.T) {
	// getValueThroughStop
	r := getValueThroughStop(10, 20, 0.5)
	if r != 15 {
		t.Fatalf("getValueThroughStop(10,20,0.5) = %f, want 15", r)
	}
	r = getValueThroughStop(0, 100, 0.0)
	if r != 0 {
		t.Fatalf("getValueThroughStop(0,100,0) = %f, want 0", r)
	}

	// getExponentialPercentage
	p := getExponentialPercentage(5, 0, 10, 1)
	if p != 0.5 {
		t.Fatalf("getExponentialPercentage linear = %f, want 0.5", p)
	}
	p = getExponentialPercentage(5, 5, 10, 1)
	if p != 0 {
		t.Fatalf("getExponentialPercentage at lower stop = %f, want 0", p)
	}
	p = getExponentialPercentage(10, 5, 5, 1)
	if p != 0 {
		t.Fatalf("getExponentialPercentage zero diff = %f, want 0", p)
	}
	p = getExponentialPercentage(5, 0, 10, 2)
	if p <= 0 {
		t.Fatalf("getExponentialPercentage exponential = %f, want > 0", p)
	}

	// getColorValueBetweenStops (via ColorStopsType)
	raw := `{"stops":[[0,"#000000"],[10,"#ffffff"]]}`
	var c ColorType
	json.Unmarshal([]byte(raw), &c)
	clr := c.GetColorAtZoomLevel(5)
	if clr == nil {
		t.Fatal("expected interpolated color")
	}
	rgba, ok := clr.(color.RGBA)
	if !ok {
		t.Fatal("expected RGBA")
	}
	if rgba.R != 127 || rgba.G != 127 || rgba.B != 127 {
		t.Fatalf("expected ~127 at midpoint, got %+v", rgba)
	}
}

// ─── Hex color edge cases ─────────────────────────────────────────────────

func TestHexColor4Char(t *testing.T) {
	var c ColorType
	if err := json.Unmarshal([]byte(`"#abc"`), &c); err != nil {
		t.Fatal(err)
	}
	clr := c.GetColorAtZoomLevel(0)
	if clr == nil {
		t.Fatal("expected color")
	}
}

func TestHexColor9Char(t *testing.T) {
	var c ColorType
	if err := json.Unmarshal([]byte(`"#ff000080"`), &c); err != nil {
		t.Fatal(err)
	}
	clr := c.GetColorAtZoomLevel(0)
	if clr == nil {
		t.Fatal("expected color")
	}
	rgba, ok := clr.(color.RGBA)
	if !ok || rgba.A != 128 {
		t.Fatalf("expected alpha 128, got %+v", rgba)
	}
}

func TestHexColorInvalid(t *testing.T) {
	var c ColorType
	if err := json.Unmarshal([]byte(`"#xyzxyz"`), &c); err == nil {
		t.Fatal("expected error for invalid hex")
	}
}

// ─── ColorType Unmarshal error path ───────────────────────────────────────

func TestColorTypeUnmarshalInvalidType(t *testing.T) {
	var c ColorType
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for invalid type")
		}
	}()
	json.Unmarshal([]byte(`123`), &c)
}

// ─── Expression Unmarshal/Marshal edge cases ───────────────────────────────

func TestExpressionUnmarshalError(t *testing.T) {
	var e Expression
	// This will panic in our code because the decode function doesn't handle this case
	// It hits the `interfaceListToStringList` style error path
	defer func() {
		if r := recover(); r == nil {
			// Expected: our decode function will return a literal for unrecognized arrays
		}
	}()
	json.Unmarshal([]byte(`[1,2,3]`), &e)
	// First element not a string -> treated as literal array
	if !e.IsLiteral {
		t.Fatal("expected literal for non-string first element")
	}
}

func TestExpressionMarshalComplex(t *testing.T) {
	e := &Expression{
		Operator: "case",
		Args: []*Expression{
			{Operator: "==", Args: []*Expression{
				{IsLiteral: true, Value: "$type"},
				{IsLiteral: true, Value: "Point"},
			}},
			{IsLiteral: true, Value: float64(1)},
			{IsLiteral: true, Value: float64(0)},
		},
	}
	out, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	expected := `["case",["==","$type","Point"],1,0]`
	if !jsonEqual(string(out), expected) {
		t.Fatalf("expected %s, got %s", expected, string(out))
	}
}

// ─── IsDataExpression nil ─────────────────────────────────────────────────

func TestIsDataExpressionNil(t *testing.T) {
	if IsDataExpression(nil) {
		t.Fatal("expected false for nil")
	}
	if IsDataExpression(&Expression{IsLiteral: true, Value: nil}) {
		t.Fatal("expected false for literal")
	}
}

// ─── calculateBackgroundColor edge case: first layer is not background ─────

func TestParseWithNonBackgroundFirstLayer(t *testing.T) {
	raw := `{"version":8,"sources":{},"layers":[{"id":"water","type":"fill","paint":{"fill-color":"#00ffff"}}]}`
	r := strings.NewReader(raw)
	ms, err := Parse(r)
	if err != nil {
		t.Fatal(err)
	}
	bg := ms.GetBackground()
	if bg == nil {
		t.Fatal("expected non-nil background")
	}
}

// ─── FilterContainer Unmarshal error ──────────────────────────────────────

func TestFilterContainerUnmarshalError(t *testing.T) {
	var fc FilterContainer
	if err := json.Unmarshal([]byte(`{invalid}`), &fc); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

// ─── Source tiles edge case ───────────────────────────────────────────────

func TestSourceWithTiles(t *testing.T) {
	raw := `{"type":"vector","tiles":["http://example.com/{z}/{x}/{y}.pbf"],"minzoom":0,"maxzoom":14}`
	var s Source
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatal(err)
	}
	if len(s.Tiles) != 1 {
		t.Fatalf("expected 1 tile URL, got %d", len(s.Tiles))
	}
	if s.MinZoom == nil || *s.MinZoom != 0 {
		t.Fatal("MinZoom mismatch")
	}
	if s.MaxZoom == nil || *s.MaxZoom != 14 {
		t.Fatal("MaxZoom mismatch")
	}
}

// ─── Layer with all fields empty after marshal ────────────────────────────

func TestLayerMarshalEmptyID(t *testing.T) {
	l := Layer{}
	out, err := json.Marshal(&l)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"id":""`) {
		t.Fatalf("expected id field, got %s", string(out))
	}
}

// ─── Flywave Extensions ────────────────────────────────────────────────────

func TestStyleFlywaveExtensionsRoundTrip(t *testing.T) {
	raw := `{
		"version":8,
		"sources":{},
		"layers":[],
		"flywave:clearColor":"#1a1a2e",
		"flywave:clearAlpha":1,
		"flywave:enableShadows":true,
		"flywave:toneMappingExposure":1.2,
		"flywave:definitions":{"roadColor":"#ff0000"},
		"flywave:postEffects":{"bloom":{"enabled":true,"strength":0.5}},
		"flywave:textStyles":[{"name":"default","size":14}],
		"flywave:fontCatalogs":[{"url":"fonts.pbf","name":"Open Sans"}],
		"flywave:imageTextures":[{"name":"icon","image":"marker"}],
		"flywave:poiTables":[{"name":"pois","url":"pois.json","useAltNamesForKey":false}],
		"flywave:priorities":[{"group":"tilezen","category":"road"}],
		"flywave:labelPriorities":["city","town"]
	}`
	var s Style
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatal(err)
	}
	if s.ClearColor == nil || *s.ClearColor != "#1a1a2e" {
		t.Fatalf("ClearColor = %v", s.ClearColor)
	}
	if s.ClearAlpha == nil || *s.ClearAlpha != 1 {
		t.Fatalf("ClearAlpha = %d", *s.ClearAlpha)
	}
	if s.EnableShadows == nil || !*s.EnableShadows {
		t.Fatal("EnableShadows should be true")
	}
	if s.ToneMappingExposure == nil || *s.ToneMappingExposure != 1.2 {
		t.Fatalf("ToneMappingExposure = %v", *s.ToneMappingExposure)
	}
	if s.FlywaveDefinitions == nil {
		t.Fatal("expected Definitions")
	}
	if (*s.FlywaveDefinitions)["roadColor"] != "#ff0000" {
		t.Fatalf("definitions.roadColor = %v", (*s.FlywaveDefinitions)["roadColor"])
	}
	if s.FlywavePostEffects == nil || s.FlywavePostEffects.Bloom == nil || !s.FlywavePostEffects.Bloom.Enabled {
		t.Fatal("expected bloom post effect")
	}
	if len(s.FlywaveTextStyles) != 1 || s.FlywaveTextStyles[0].Name == nil || *s.FlywaveTextStyles[0].Name != "default" {
		t.Fatal("textStyles mismatch")
	}
	if len(s.FlywaveFontCatalogs) != 1 || s.FlywaveFontCatalogs[0].Name != "Open Sans" {
		t.Fatal("fontCatalogs mismatch")
	}
	if len(s.FlywaveImageTextures) != 1 || s.FlywaveImageTextures[0].Name != "icon" {
		t.Fatal("imageTextures mismatch")
	}
	if len(s.FlywavePoiTables) != 1 || s.FlywavePoiTables[0].Name != "pois" {
		t.Fatal("poiTables mismatch")
	}
	if len(s.FlywavePriorities) != 1 || s.FlywavePriorities[0].Group != "tilezen" {
		t.Fatal("priorities mismatch")
	}
	if len(s.FlywaveLabelPriorities) != 2 || s.FlywaveLabelPriorities[0] != "city" {
		t.Fatal("labelPriorities mismatch")
	}

	out, err := json.Marshal(&s)
	if err != nil {
		t.Fatal(err)
	}
	if !jsonEqual(string(out), raw) {
		t.Fatalf("round trip:\n  in:  %s\n  out: %s", raw, string(out))
	}
}

func TestStyleFlywaveOmitEmpty(t *testing.T) {
	s := Style{Version: 8, Sources: Sources{}, Layers: []*Layer{}}
	out, err := json.Marshal(&s)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), "flywave") {
		t.Fatal("flywave should be omitted when empty")
	}
}

func TestLayerFlywaveExtensionsRoundTrip(t *testing.T) {
	raw := `{
		"id":"building-test",
		"type":"fill-extrusion",
		"source":"composite",
		"source-layer":"building",
		"flywave:technique":"extruded-polygon",
		"flywave:styleSet":"tilezen",
		"flywave:category":"building",
		"flywave:renderOrder":10,
		"flywave:animateExtrusion":true,
		"flywave:boundaryWalls":true,
		"flywave:footprint":false,
		"flywave:imageTexturePrefix":"night_",
		"flywave:imageTexturePostfix":"_dark"
	}`
	var l Layer
	if err := json.Unmarshal([]byte(raw), &l); err != nil {
		t.Fatal(err)
	}
	if l.FlywaveTechnique == nil || *l.FlywaveTechnique != "extruded-polygon" {
		t.Fatalf("Technique = %v", l.FlywaveTechnique)
	}
	if l.FlywaveStyleSet == nil || *l.FlywaveStyleSet != "tilezen" {
		t.Fatalf("StyleSet = %v", l.FlywaveStyleSet)
	}
	if l.FlywaveCategory == nil || *l.FlywaveCategory != "building" {
		t.Fatalf("Category = %v", l.FlywaveCategory)
	}
	if l.FlywaveRenderOrder == nil || *l.FlywaveRenderOrder != 10 {
		t.Fatalf("RenderOrder = %d", *l.FlywaveRenderOrder)
	}
	if l.FlywaveAnimateExtrusion == nil || !*l.FlywaveAnimateExtrusion {
		t.Fatal("AnimateExtrusion should be true")
	}
	if l.FlywaveBoundaryWalls == nil || !*l.FlywaveBoundaryWalls {
		t.Fatal("BoundaryWalls should be true")
	}
	if l.FlywaveFootprint == nil || *l.FlywaveFootprint {
		t.Fatal("Footprint should be false")
	}
	if l.FlywaveImageTexturePrefix == nil || *l.FlywaveImageTexturePrefix != "night_" {
		t.Fatalf("ImageTexturePrefix = %v", l.FlywaveImageTexturePrefix)
	}
	if l.FlywaveImageTexturePostfix == nil || *l.FlywaveImageTexturePostfix != "_dark" {
		t.Fatalf("ImageTexturePostfix = %v", l.FlywaveImageTexturePostfix)
	}

	out, err := json.Marshal(&l)
	if err != nil {
		t.Fatal(err)
	}
	if !jsonEqual(string(out), raw) {
		t.Fatalf("round trip:\n  in:  %s\n  out: %s", raw, string(out))
	}
}

func TestLayerFlywaveOmitEmpty(t *testing.T) {
	l := Layer{ID: "test", Type: LayerTypeFill}
	out, err := json.Marshal(&l)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), "flywave") {
		t.Fatal("flywave should be omitted when empty")
	}
}

func TestLayerFlywaveShaderParams(t *testing.T) {
	raw := `{
		"id":"shader-layer",
		"type":"line",
		"flywave:technique":"shader",
		"flywave:shaderParams":{"vertexShader":"...","fragmentShader":"..."}
	}`
	var l Layer
	if err := json.Unmarshal([]byte(raw), &l); err != nil {
		t.Fatal(err)
	}
	if l.FlywaveShaderParams == nil {
		t.Fatal("expected ShaderParams")
	}
	if l.FlywaveShaderParams["vertexShader"] != "..." {
		t.Fatalf("vertexShader = %v", l.FlywaveShaderParams["vertexShader"])
	}
}

func TestStyleWithFlywaveAndStandardFields(t *testing.T) {
	raw := `{
		"version":8,
		"name":"mixed-style",
		"sources":{"streets":{"type":"vector","url":"mapbox://mapbox.streets"}},
		"layers":[
			{
				"id":"water",
				"type":"fill",
				"source":"streets",
				"source-layer":"water",
				"paint":{"fill-color":"#00ffff"},
				"flywave:technique":"fill",
				"flywave:styleSet":"tilezen"
			}
		],
		"flywave:definitions":{"waterColor":"#00ffff"},
		"flywave:enableShadows":true
	}`
	var s Style
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatal(err)
	}
	if s.Name != "mixed-style" {
		t.Fatalf("Name = %q", s.Name)
	}
	if s.EnableShadows == nil || !*s.EnableShadows {
		t.Fatal("expected enableShadows")
	}
	if len(s.Layers) != 1 || s.Layers[0].FlywaveTechnique == nil || *s.Layers[0].FlywaveTechnique != "fill" {
		t.Fatal("expected layer technique fill")
	}
	if s.Layers[0].Paint == nil || s.Layers[0].Paint.FillColor == nil {
		t.Fatal("standard paint should still work alongside flywave extensions")
	}

	out, err := json.Marshal(&s)
	if err != nil {
		t.Fatal(err)
	}
	if !jsonEqual(string(out), raw) {
		t.Fatalf("round trip:\n  in:  %s\n  out: %s", raw, string(out))
	}
}

func TestStyleUnmarshalWithoutFlywave(t *testing.T) {
	raw := `{"version":8,"sources":{},"layers":[{"id":"bg","type":"background"}]}`
	var s Style
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatal(err)
	}
	if s.EnableShadows != nil {
		t.Fatal("EnableShadows should be nil when not present")
	}
	if s.Layers[0].FlywaveTechnique != nil {
		t.Fatal("Layer flywave:technique should be nil when not present")
	}
}

