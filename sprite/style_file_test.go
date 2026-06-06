package sprite

import (
	"testing"
)

func TestParseStyleFile(t *testing.T) {
	tests := []struct {
		input    string
		wantType string
		wantExt  string
		wantPR   int
		out      string
	}{
		{"sprite@2x.json", "sprite", ".json", 2, "sprite@2x.json"},
		{"sprite.json", "sprite", ".json", 1, "sprite.json"},
		{"sprite@2x.png", "sprite", ".png", 2, "sprite@2x.png"},
		{"icons@3x.png", "icons", ".png", 3, "icons@3x.png"},
		{"marker@1x.png", "marker", ".png", 1, "marker.png"},
		{"MARKER.PNG", "marker", ".png", 1, "marker.png"},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			sf, err := ParseStyleFile(tc.input)
			if err != nil {
				t.Fatalf("ParseStyleFile(%q): %v", tc.input, err)
			}
			if sf.Type != tc.wantType {
				t.Errorf("Type: got %q, want %q", sf.Type, tc.wantType)
			}
			if sf.Ext != tc.wantExt {
				t.Errorf("Ext: got %q, want %q", sf.Ext, tc.wantExt)
			}
			if sf.PixelRatio != tc.wantPR {
				t.Errorf("PixelRatio: got %d, want %d", sf.PixelRatio, tc.wantPR)
			}
			if got := sf.String(); got != tc.out {
				t.Errorf("String(): got %q, want %q", got, tc.out)
			}
		})
	}
}

func TestParseStyleFile_Errors(t *testing.T) {
	inputs := []string{
		"",
		"noextension",
		".hidden",
	}
	for _, s := range inputs {
		t.Run(s, func(t *testing.T) {
			_, err := ParseStyleFile(s)
			if err == nil {
				t.Errorf("expected error for %q", s)
			}
		})
	}
}

func TestStyleFileMarshalText(t *testing.T) {
	sf := &StyleFile{Type: "sprite", Ext: ".json", PixelRatio: 2}
	text, err := sf.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	if string(text) != "sprite@2x.json" {
		t.Fatalf("got %q", string(text))
	}
}

func TestStyleFileUnmarshalText(t *testing.T) {
	var sf StyleFile
	if err := sf.UnmarshalText([]byte("marker@3x.png")); err != nil {
		t.Fatal(err)
	}
	if sf.Type != "marker" || sf.Ext != ".png" || sf.PixelRatio != 3 {
		t.Fatalf("got %+v", sf)
	}
}

func TestStyleFileString(t *testing.T) {
	tests := []struct {
		sf   StyleFile
		want string
	}{
		{StyleFile{Type: "test", Ext: ".json", PixelRatio: 1}, "test.json"},
		{StyleFile{Type: "test", Ext: ".json", PixelRatio: 0}, "test.json"},
		{StyleFile{Type: "test", Ext: ".png", PixelRatio: 2}, "test@2x.png"},
	}
	for _, tc := range tests {
		if got := tc.sf.String(); got != tc.want {
			t.Errorf("String(%+v) = %q, want %q", tc.sf, got, tc.want)
		}
	}
}
