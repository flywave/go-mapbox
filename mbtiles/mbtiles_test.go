package mbtiles

import (
	"os"
	"testing"
)

func Test_TileFormat_String(t *testing.T) {
	var conditions = []struct {
		in  TileFormat
		out string
	}{
		{UNKNOWN, ""},
		{PNG, "png"},
		{JPG, "jpg"},
		{PNG, "png"},
		{PBF, "pbf"},
		{WEBP, "webp"},
	}

	for _, condition := range conditions {
		if condition.in.String() != condition.out {
			t.Errorf("%q.String() => %q, expected %q", condition.in, condition.in.String(), condition.out)
		}
	}
}

func Test_TileFormat_ContentType(t *testing.T) {
	var conditions = []struct {
		in  TileFormat
		out string
	}{
		{UNKNOWN, ""},
		{PNG, "image/png"},
		{JPG, "image/jpeg"},
		{PNG, "image/png"},
		{PBF, "application/x-protobuf"},
		{WEBP, "image/webp"},
	}

	for _, condition := range conditions {
		if condition.in.ContentType() != condition.out {
			t.Errorf("%q.ContentType() => %q, expected %q", condition.in, condition.in.ContentType(), condition.out)
		}
	}
}

func TestWriteMBTiles(t *testing.T) {
	db, err := CreateDB("./test.mbtils", PNG, "", "")

	if err != nil {
		t.FailNow()
	}

	db.StoreTile(0, 0, 0, []byte("test"))

	db.Close()

	os.Remove("./test.mbtils")
}
