package raster

import (
	"os"
	"testing"
)

func TestSlice(t *testing.T) {
	ts := []float64{0, 1, 2, 3}
	if l := ts[0:2]; len(l) != 2 {
		t.Error("error")
	} else {
		print(l[0])
	}
}

func TestDEM(t *testing.T) {
	var terrain *DEMData
	if fuji, err := LoadDEMData("../../tests/fuji.png", DEM_ENCODING_MAPBOX); err != nil {
		t.Error(err)
	} else {
		terrain = fuji
	}
	h := terrain.Get(0, 0)
	if h > 500 {
		t.Error("error")
	}
	data := terrain.GetData()
	if len(data) != 512*512 {
		t.Error("error")
	}
	if err := terrain.Save("../../tests/fuji1.png"); err != nil {
		t.Error(err)
	}
	if _, err := os.Stat("../../tests/fuji1.png"); err == nil {
		if err := os.Remove("../../tests/fuji1.png"); err != nil {
			t.Fail()
		}
	}
}
