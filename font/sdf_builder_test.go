package fonts

import (
	"os"
	"testing"

	"github.com/flywave/freetype/truetype"
	"github.com/flywave/go-pbf"
)

func fontPath(family string) string {
	return "../data/fonts/" + family + ".ttf"
}

func TestSDFBuilder_Creation(t *testing.T) {
	ttf, err := os.ReadFile(fontPath("NotoSans-Regular"))
	if err != nil {
		t.Skip("NotoSans-Regular.ttf not found, skipping")
	}
	font, err := truetype.Parse(ttf)
	if err != nil {
		t.Fatalf("parse font: %v", err)
	}
	b := NewSDFBuilder(font, SDFBuilderOpt{FontSize: 24, Buffer: 3})
	if b.FontSize != 24 {
		t.Fatalf("FontSize = %f", b.FontSize)
	}
	if b.Buffer != 3 {
		t.Fatalf("Buffer = %f", b.Buffer)
	}
	if b.Face == nil {
		t.Fatal("Face is nil")
	}
}

func TestSDFBuilder_Glyph(t *testing.T) {
	ttf, err := os.ReadFile(fontPath("NotoSans-Regular"))
	if err != nil {
		t.Skip("NotoSans-Regular.ttf not found, skipping")
	}
	font, err := truetype.Parse(ttf)
	if err != nil {
		t.Fatalf("parse font: %v", err)
	}
	b := NewSDFBuilder(font, SDFBuilderOpt{FontSize: 24, Buffer: 3})

	g := b.Glyph('A')
	if g == nil {
		t.Fatal("Glyph('A') returned nil")
	}
	if g.ID != 65 {
		t.Errorf("ID = %d, want 65", g.ID)
	}
	if g.Width == 0 || g.Height == 0 {
		t.Errorf("zero dimensions: %dx%d", g.Width, g.Height)
	}
	if len(g.Bitmap) == 0 {
		t.Error("empty SDF bitmap")
	}
}

func TestSDFBuilder_GlyphNil(t *testing.T) {
	ttf, err := os.ReadFile(fontPath("NotoSans-Regular"))
	if err != nil {
		t.Skip("NotoSans-Regular.ttf not found, skipping")
	}
	font, err := truetype.Parse(ttf)
	if err != nil {
		t.Fatalf("parse font: %v", err)
	}
	b := NewSDFBuilder(font)

	if g := b.Glyph(0); g != nil {
		t.Error("Glyph(0) should be nil")
	}
}

func TestSDFBuilder_GlyphsPBF(t *testing.T) {
	ttf, err := os.ReadFile(fontPath("NotoSans-Regular"))
	if err != nil {
		t.Skip("NotoSans-Regular.ttf not found, skipping")
	}
	font, err := truetype.Parse(ttf)
	if err != nil {
		t.Fatalf("parse font: %v", err)
	}
	b := NewSDFBuilder(font, SDFBuilderOpt{FontSize: 24, Buffer: 3})

	gs := b.Glyphs(65, 70)
	if gs == nil {
		t.Fatal("Glyphs returned nil")
	}
	if len(gs.Stacks) != 1 {
		t.Fatalf("expected 1 stack, got %d", len(gs.Stacks))
	}
	stack := gs.Stacks[0]
	if stack.Range != "65-69" {
		t.Errorf("Range = %q, want 65-69", stack.Range)
	}
	if len(stack.Glyphs) != 5 {
		t.Fatalf("expected 5 glyphs, got %d", len(stack.Glyphs))
	}

	w := pbf.NewWriter()
	if err := gs.WritePBF(w); err != nil {
		t.Fatalf("WritePBF: %v", err)
	}
	buf := w.Finish()

	r := pbf.NewReader(buf)
	read := &Glyphs{}
	if err := read.ReadPBF(r); err != nil {
		t.Fatalf("ReadPBF: %v", err)
	}
	if len(read.Stacks) != 1 {
		t.Fatal("round-trip: expected 1 stack")
	}
	if len(read.Stacks[0].Glyphs) != 5 {
		t.Fatalf("round-trip: expected 5 glyphs, got %d", len(read.Stacks[0].Glyphs))
	}
}
