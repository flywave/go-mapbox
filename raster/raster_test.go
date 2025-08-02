package raster

import (
	"fmt"
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
	if fuji, err := LoadDEMData("../data/fuji.png", DEM_ENCODING_MAPBOX); err != nil {
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
	if err := terrain.Save("../data/fuji1.png"); err != nil {
		t.Error(err)
	}
	if _, err := os.Stat("../data/fuji1.png"); err == nil {
		if err := os.Remove("../data/fuji1.png"); err != nil {
			t.Fail()
		}
	}
}

func TestPack(t *testing.T) {
	mp := &MapboxPacker{}
	b := mp.Pack(200)
	res := float64(b[0])*UNPACK_MAPBOX[0] + float64(b[1])*UNPACK_MAPBOX[1] + float64(b[2])*UNPACK_MAPBOX[2] - float64(b[0])*UNPACK_MAPBOX[3]
	fmt.Println(res)
}
