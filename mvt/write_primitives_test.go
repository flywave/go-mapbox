package mvt

import (
	"testing"
	//"fmt"
	"github.com/murphy214/pbf"
)

var values = []interface{}{"string", float32(100.23), float64(100.23), int64(10201203912), uint64(10201203912), true}
var values_bytes = [][]byte{{0x22, 0x8, 0xa, 0x6, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67}, {0x22, 0x5, 0x15, 0xc3, 0x75, 0xc8, 0x42}, {0x22, 0x9, 0x19, 0x1f, 0x85, 0xeb, 0x51, 0xb8, 0xe, 0x59, 0x40}, {0x22, 0x6, 0x20, 0xc8, 0x89, 0xa8, 0x80, 0x26}, {0x22, 0x6, 0x28, 0xc8, 0x89, 0xa8, 0x80, 0x26}, {0x22, 0x2, 0x38, 0x1}}

func TestWriteValue(t *testing.T) {
	for pos, i := range values {
		expected_bytevals := values_bytes[pos]
		bytevals := WriteValue(i)
		for i := range bytevals {
			if bytevals[i] != expected_bytevals[i] {
				t.Errorf("WriteValue Error w/ value %v", i)
			}
		}
	}
}

func TestEncodeVarint(t *testing.T) {
	expected_bytevals := []byte{0xe8, 0x1}
	bytevals := pbf.EncodeVarint(232)
	for i := range bytevals {
		if bytevals[i] != expected_bytevals[i] {
			t.Errorf("EncodeVarint Error")
		}
	}
}

func TestEncodeVarint32(t *testing.T) {
	expected_bytevals := []byte{0xe8, 0x1}
	bytevals := EncodeVarint32(232)
	for i := range bytevals {
		if bytevals[i] != expected_bytevals[i] {
			t.Errorf("EncodeVarint32 Error")
		}
	}
}

func TestWritePackedUint32(t *testing.T) {
	expected_bytevals := []byte{0x9, 0xa, 0x90, 0x3, 0xf4, 0x3, 0xd8, 0x4, 0xbc, 0x5}
	bytevals := WritePackedUint32([]uint32{10, 400, 500, 600, 700})
	for i := range bytevals {
		if bytevals[i] != expected_bytevals[i] {
			t.Errorf("WritePackedUint32 Error")
		}
	}
}
