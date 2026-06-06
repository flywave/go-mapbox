package fonts

import (
	"bytes"
	"testing"

	"github.com/flywave/go-pbf"
)

func TestGlyphReadWritePBF(t *testing.T) {
	original := NewGlyph(65, 100, 120, 110, 10, -20, []byte{0x01, 0x02, 0x03, 0x04})

	w := pbf.NewWriter()
	if err := original.WritePBF(w); err != nil {
		t.Fatalf("WritePBF: %v", err)
	}
	buf := w.Finish()

	r := pbf.NewReader(buf)
	read := &Glyph{}
	if err := read.ReadPBF(r); err != nil {
		t.Fatalf("ReadPBF: %v", err)
	}

	if read.ID != original.ID {
		t.Errorf("ID: got %d, want %d", read.ID, original.ID)
	}
	if read.Width != original.Width {
		t.Errorf("Width: got %d, want %d", read.Width, original.Width)
	}
	if read.Height != original.Height {
		t.Errorf("Height: got %d, want %d", read.Height, original.Height)
	}
	if read.Advance != original.Advance {
		t.Errorf("Advance: got %d, want %d", read.Advance, original.Advance)
	}
	if read.Left != original.Left {
		t.Errorf("Left: got %d, want %d", read.Left, original.Left)
	}
	if read.Top != original.Top {
		t.Errorf("Top: got %d, want %d", read.Top, original.Top)
	}
	if !bytes.Equal(read.Bitmap, original.Bitmap) {
		t.Errorf("Bitmap: got %v, want %v", read.Bitmap, original.Bitmap)
	}
}

func TestGlyphWriteEmptyBitmap(t *testing.T) {
	g := NewGlyph(32, 50, 60, 40, 5, -10, nil)
	w := pbf.NewWriter()
	if err := g.WritePBF(w); err != nil {
		t.Fatal(err)
	}
	buf := w.Finish()
	r := pbf.NewReader(buf)
	read := &Glyph{}
	if err := read.ReadPBF(r); err != nil {
		t.Fatal(err)
	}
	if read.ID != 32 {
		t.Fatalf("ID = %d", read.ID)
	}
	if len(read.Bitmap) != 0 {
		t.Fatal("expected empty bitmap")
	}
}

func TestGlyphReadPBFUnknownTag(t *testing.T) {
	// Build a glyph PBF, then insert an extra field with an unused tag
	orig := NewGlyph(65, 100, 120, 110, 10, -20, []byte{0x01, 0x02})
	w := pbf.NewWriter()
	orig.WritePBF(w)
	// Append an unknown field after the normal glyph fields
	w.WriteUInt32(99, 99)
	buf := w.Finish()

	r := pbf.NewReader(buf)
	g := &Glyph{}
	if err := g.ReadPBF(r); err != nil {
		t.Fatalf("should skip unknown tag, got error: %v", err)
	}
	if g.ID != 65 {
		t.Fatalf("ID = %d, want 65", g.ID)
	}
}

func TestFontstackReadWritePBF(t *testing.T) {
	stack := NewFontstack("Arial", "0-255")
	stack.Glyphs = append(stack.Glyphs,
		NewGlyph(65, 100, 120, 110, 10, -20, []byte{0x01, 0x02}),
		NewGlyph(66, 90, 120, 100, 15, -22, []byte{0x03, 0x04}),
	)

	w := pbf.NewWriter()
	if err := stack.WritePBF(w); err != nil {
		t.Fatalf("WritePBF: %v", err)
	}
	buf := w.Finish()

	r := pbf.NewReader(buf)
	read := &Fontstack{}
	if err := read.ReadPBF(r); err != nil {
		t.Fatalf("ReadPBF: %v", err)
	}

	if read.Name != stack.Name {
		t.Errorf("Name: got %q, want %q", read.Name, stack.Name)
	}
	if read.Range != stack.Range {
		t.Errorf("Range: got %q, want %q", read.Range, stack.Range)
	}
	if len(read.Glyphs) != len(stack.Glyphs) {
		t.Fatalf("glyph count: got %d, want %d", len(read.Glyphs), len(stack.Glyphs))
	}
	for i := range stack.Glyphs {
		orig := stack.Glyphs[i]
		got := read.Glyphs[i]
		if got.ID != orig.ID {
			t.Errorf("glyph[%d].ID: got %d, want %d", i, got.ID, orig.ID)
		}
		if got.Width != orig.Width {
			t.Errorf("glyph[%d].Width: got %d, want %d", i, got.Width, orig.Width)
		}
		if got.Height != orig.Height {
			t.Errorf("glyph[%d].Height: got %d, want %d", i, got.Height, orig.Height)
		}
		if got.Left != orig.Left {
			t.Errorf("glyph[%d].Left: got %d, want %d", i, got.Left, orig.Left)
		}
		if got.Top != orig.Top {
			t.Errorf("glyph[%d].Top: got %d, want %d", i, got.Top, orig.Top)
		}
		if got.Advance != orig.Advance {
			t.Errorf("glyph[%d].Advance: got %d, want %d", i, got.Advance, orig.Advance)
		}
		if !bytes.Equal(got.Bitmap, orig.Bitmap) {
			t.Errorf("glyph[%d].Bitmap mismatch", i)
		}
	}
}

func TestGlyphsReadWritePBF(t *testing.T) {
	glyphs := NewGlyphs()
	glyphs.Stacks = append(glyphs.Stacks,
		NewFontstack("Arial", "0-255"),
		NewFontstack("Times New Roman", "0-127"),
	)
	glyphs.Stacks[0].Glyphs = append(glyphs.Stacks[0].Glyphs,
		NewGlyph(65, 100, 120, 110, 10, -20, []byte{0x01, 0x02}),
	)
	glyphs.Stacks[1].Glyphs = append(glyphs.Stacks[1].Glyphs,
		NewGlyph(66, 90, 120, 100, 15, -22, []byte{0x03, 0x04}),
	)

	w := pbf.NewWriter()
	if err := glyphs.WritePBF(w); err != nil {
		t.Fatalf("WritePBF: %v", err)
	}
	buf := w.Finish()

	r := pbf.NewReader(buf)
	read := &Glyphs{}
	if err := read.ReadPBF(r); err != nil {
		t.Fatalf("ReadPBF: %v", err)
	}

	if len(read.Stacks) != len(glyphs.Stacks) {
		t.Fatalf("stack count: got %d, want %d", len(read.Stacks), len(glyphs.Stacks))
	}
	for i := range glyphs.Stacks {
		orig := glyphs.Stacks[i]
		got := read.Stacks[i]
		if got.Name != orig.Name {
			t.Errorf("stack[%d].Name: got %q, want %q", i, got.Name, orig.Name)
		}
		if len(got.Glyphs) != len(orig.Glyphs) {
			t.Errorf("stack[%d] glyph count: got %d, want %d", i, len(got.Glyphs), len(orig.Glyphs))
			continue
		}
		for j := range orig.Glyphs {
			og := orig.Glyphs[j]
			g := got.Glyphs[j]
			if g.ID != og.ID {
				t.Errorf("stack[%d].glyph[%d].ID: got %d, want %d", i, j, g.ID, og.ID)
			}
			if !bytes.Equal(g.Bitmap, og.Bitmap) {
				t.Errorf("stack[%d].glyph[%d].Bitmap mismatch", i, j)
			}
		}
	}
}

func TestNewHelpers(t *testing.T) {
	g := NewGlyph(97, 50, 60, 45, 5, -10, []byte{0xff})
	if g.ID != 97 || g.Advance != 45 {
		t.Fatal("NewGlyph fields mismatch")
	}

	fs := NewFontstack("OpenSans", "0-127")
	if fs.Name != "OpenSans" || fs.Range != "0-127" || len(fs.Glyphs) != 0 {
		t.Fatal("NewFontstack fields mismatch")
	}

	gs := NewGlyphs()
	if len(gs.Stacks) != 0 {
		t.Fatal("NewGlyphs should have empty Stacks")
	}
}
