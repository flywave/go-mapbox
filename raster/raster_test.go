package raster

import (
	"image"
	"image/png"
	"math"
	"os"
	"testing"
)

// ─── NewDEMData ────────────────────────────────────────────────────────────

func TestNewDEMData_Nil(t *testing.T) {
	d := NewDEMData(nil, DEM_ENCODING_MAPBOX)
	if d != nil {
		t.Fatal("expected nil for nil data")
	}
}

func TestNewDEMData_Empty(t *testing.T) {
	d := NewDEMData([][4]byte{}, DEM_ENCODING_MAPBOX)
	if d != nil {
		t.Fatal("expected nil for empty data")
	}
}

func TestNewDEMData_OddLength(t *testing.T) {
	d := NewDEMData([][4]byte{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}}, DEM_ENCODING_MAPBOX)
	if d != nil {
		t.Fatal("expected nil for odd-length data (not a perfect square)")
	}
}

func TestNewDEMData_SquareGrid(t *testing.T) {
	n := 4
	data := make([][4]byte, n*n)
	for i := range data {
		data[i] = [4]byte{byte(i), byte(i + 1), byte(i + 2), 255}
	}
	d := NewDEMData(data, DEM_ENCODING_MAPBOX)
	if d == nil {
		t.Fatal("expected non-nil DEMData")
	}
	if d.Dim != n {
		t.Fatalf("Dim = %d, want %d", d.Dim, n)
	}
	if d.Stride != n+2 {
		t.Fatalf("Stride = %d, want %d", d.Stride, n+2)
	}
	if d.Encoding != DEM_ENCODING_MAPBOX {
		t.Fatal("Encoding mismatch")
	}
	// Should have (n+2)*(n+2) elements (border padding)
	expectedLen := (n + 2) * (n + 2)
	if len(d.Data) != expectedLen {
		t.Fatalf("Data len = %d, want %d", len(d.Data), expectedLen)
	}
}

// ─── Get / GetData ─────────────────────────────────────────────────────────

func TestGetAndGetData(t *testing.T) {
	n := 4
	data := make([][4]byte, n*n)
	for y := 0; y < n; y++ {
		for x := 0; x < n; x++ {
			idx := y*n + x
			data[idx] = [4]byte{0, 0, byte(idx), 255}
		}
	}
	d := NewDEMData(data, DEM_ENCODING_MAPBOX)
	if d == nil {
		t.Fatal("NewDEMData failed")
	}

	// Get should return a value (not panic)
	for x := 0; x < n; x++ {
		for y := 0; y < n; y++ {
			v := d.Get(x, y)
			if math.IsNaN(v) {
				t.Fatalf("Get(%d,%d) = NaN", x, y)
			}
		}
	}

	// GetData should return correct length
	gd := d.GetData()
	if len(gd) != n*n {
		t.Fatalf("GetData len = %d, want %d", len(gd), n*n)
	}
}

// ─── BackfillBorder ────────────────────────────────────────────────────────

func TestBackfillBorder_SameDim(t *testing.T) {
	n := 4
	data := make([][4]byte, n*n)
	for i := range data {
		data[i] = [4]byte{1, 2, 3, 255}
	}
	d := NewDEMData(data, DEM_ENCODING_MAPBOX)
	other := NewDEMData(data, DEM_ENCODING_MAPBOX)
	// Same dim: should be no-op
	d.BackfillBorder(*other, 0, 0)
}

func TestBackfillBorder_DifferentDim(t *testing.T) {
	n := 4
	data := make([][4]byte, n*n)
	for i := range data {
		data[i] = [4]byte{10, 20, 30, 255}
	}
	d := NewDEMData(data, DEM_ENCODING_MAPBOX)
	other := NewDEMData(data, DEM_ENCODING_MAPBOX)
	// dx=1, dy=0
	d.BackfillBorder(*other, 1, 0)
}

// ─── Save ──────────────────────────────────────────────────────────────────

func TestSaveAndLoad(t *testing.T) {
	n := 4
	data := make([][4]byte, n*n)
	for y := 0; y < n; y++ {
		for x := 0; x < n; x++ {
			idx := y*n + x
			data[idx] = [4]byte{byte(x * 50), byte(y * 50), byte(idx * 10), 255}
		}
	}
	d := NewDEMData(data, DEM_ENCODING_MAPBOX)
	if d == nil {
		t.Fatal("NewDEMData failed")
	}

	tmpFile := t.TempDir() + "/test_dem.png"
	if err := d.Save(tmpFile); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Fatal("saved file does not exist")
	}

	// Load it back
	loaded, err := LoadDEMData(tmpFile, DEM_ENCODING_MAPBOX)
	if err != nil {
		t.Fatalf("LoadDEMData: %v", err)
	}
	if loaded.Dim != n {
		t.Fatalf("loaded Dim = %d, want %d", loaded.Dim, n)
	}
}

func TestLoadDEMData_NotSquare(t *testing.T) {
	// Create a non-square image
	img := image.NewRGBA(image.Rect(0, 0, 4, 8))
	tmpFile := t.TempDir() + "/notsquare.png"
	f, _ := os.Create(tmpFile)
	png.Encode(f, img)
	f.Close()

	_, err := LoadDEMData(tmpFile, DEM_ENCODING_MAPBOX)
	if err == nil {
		t.Fatal("expected error for non-square image")
	}
}

// ─── MapboxPacker ──────────────────────────────────────────────────────────

func TestMapboxPacker_RoundTrip(t *testing.T) {
	mp := &MapboxPacker{}
	heights := []float64{0, 100, 500, 1000, 5000}
	for _, h := range heights {
		b := mp.Pack(h)
		if b[3] != 255 {
			t.Fatalf("alpha = %d, want 255", b[3])
		}
		unpacked := float64(b[0])*UNPACK_MAPBOX[0] +
			float64(b[1])*UNPACK_MAPBOX[1] +
			float64(b[2])*UNPACK_MAPBOX[2] -
			UNPACK_MAPBOX[3]
		// Should approximately match (some precision loss is expected)
		diff := math.Abs(unpacked - h)
		if diff > 100 {
			t.Fatalf("height %f: packed/unpacked diff = %f", h, diff)
		}
	}
}

func TestMapboxPacker_Negative(t *testing.T) {
	mp := &MapboxPacker{}
	b := mp.Pack(-500)
	if b[3] != 255 {
		t.Fatalf("alpha = %d", b[3])
	}
}

// ─── TerrariumPacker ───────────────────────────────────────────────────────

func TestTerrariumPacker_RoundTrip(t *testing.T) {
	tp := &TerrariumPacker{}
	heights := []float64{0, 100, 500, 1000}
	for _, h := range heights {
		b := tp.Pack(h)
		if b[3] != 255 {
			t.Fatalf("alpha = %d, want 255", b[3])
		}
		unpacked := float64(b[0])*UNPACK_TERRARIUM[0] +
			float64(b[1])*UNPACK_TERRARIUM[1] +
			float64(b[2])*UNPACK_TERRARIUM[2] -
			UNPACK_TERRARIUM[3]
		diff := math.Abs(unpacked - h)
		if diff > 10 {
			t.Fatalf("height %f: packed/unpacked diff = %f", h, diff)
		}
	}
}

func TestTerrariumPacker_Negative(t *testing.T) {
	tp := &TerrariumPacker{}
	b := tp.Pack(-500)
	if b[3] != 255 {
		t.Fatalf("alpha = %d", b[3])
	}
}

// ─── getUnpackVector ───────────────────────────────────────────────────────

func TestGetUnpackVector_Mapbox(t *testing.T) {
	n := 4
	data := make([][4]byte, n*n)
	d := NewDEMData(data, DEM_ENCODING_MAPBOX)
	v := d.getUnpackVector()
	if v != UNPACK_MAPBOX {
		t.Fatal("expected UNPACK_MAPBOX")
	}
}

func TestGetUnpackVector_Terrarium(t *testing.T) {
	n := 4
	data := make([][4]byte, n*n)
	d := NewDEMData(data, DEM_ENCODING_TERRARIUM)
	v := d.getUnpackVector()
	if v != UNPACK_TERRARIUM {
		t.Fatal("expected UNPACK_TERRARIUM")
	}
}

// ─── DemPacker interface ───────────────────────────────────────────────────

func TestDemPackerInterface(t *testing.T) {
	var mp DemPacker = &MapboxPacker{}
	var tp DemPacker = &TerrariumPacker{}
	_ = mp.Pack(100)
	_ = tp.Pack(100)
}

// ─── LoadDEMData with stream ───────────────────────────────────────────────

func TestLoadDEMDataWithStream_InvalidData(t *testing.T) {
	_, err := LoadDEMDataWithStream(nil, DEM_ENCODING_MAPBOX)
	if err == nil {
		t.Fatal("expected error for nil reader")
	}
}

// ─── edge cases: NewDEMData border mirror ──────────────────────────────────

func TestNewDEMData_BorderMirror(t *testing.T) {
	n := 2
	data := make([][4]byte, n*n)
	for i := range data {
		data[i] = [4]byte{byte(i), byte(i), byte(i), 255}
	}
	d := NewDEMData(data, DEM_ENCODING_MAPBOX)
	if d == nil {
		t.Fatal("NewDEMData failed")
	}
	// Corner border check: Data[0] (top-left of stride) should mirror Data[stride+1]
	cornerOutside := d.Data[0]
	cornerInside := d.Data[d.Stride+1]
	if cornerOutside != cornerInside {
		t.Fatal("border mirror failed for top-left corner")
	}
}

func TestNewDEMData_PerfectSquareLength(t *testing.T) {
	// 9 elements = 3x3
	data := make([][4]byte, 9)
	for i := range data {
		data[i] = [4]byte{1, 2, 3, 255}
	}
	d := NewDEMData(data, DEM_ENCODING_MAPBOX)
	if d == nil {
		t.Fatal("expected non-nil for 9 elements (3x3)")
	}
	if d.Dim != 3 {
		t.Fatalf("Dim = %d, want 3", d.Dim)
	}
}

func TestNewDEMData_NotPerfectSquare(t *testing.T) {
	// 10 elements is not a perfect square
	data := make([][4]byte, 10)
	d := NewDEMData(data, DEM_ENCODING_MAPBOX)
	if d != nil {
		t.Fatal("expected nil for non-perfect-square length")
	}
}
