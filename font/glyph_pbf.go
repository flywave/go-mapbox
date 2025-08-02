// Package fonts provides glyph and fontstack definitions using PBF format.
package fonts

import (
	"github.com/flywave/go-pbf"
)

// Proto definitions for font structures
const (
	// Glyph tags
	GlyphID      pbf.TagType = 1
	GlyphBitmap  pbf.TagType = 2
	GlyphWidth   pbf.TagType = 3
	GlyphHeight  pbf.TagType = 4
	GlyphLeft    pbf.TagType = 5
	GlyphTop     pbf.TagType = 6
	GlyphAdvance pbf.TagType = 7

	// Fontstack tags
	FontstackName   pbf.TagType = 1
	FontstackRange  pbf.TagType = 2
	FontstackGlyphs pbf.TagType = 3

	// Glyphs tags
	GlyphsStacks pbf.TagType = 1
)

// ProtoGlyph represents protobuf tags for Glyph
type ProtoGlyph struct {
	ID      pbf.TagType
	Bitmap  pbf.TagType
	Width   pbf.TagType
	Height  pbf.TagType
	Left    pbf.TagType
	Top     pbf.TagType
	Advance pbf.TagType
}

// ProtoFontstack represents protobuf tags for Fontstack
type ProtoFontstack struct {
	Name   pbf.TagType
	Range  pbf.TagType
	Glyphs pbf.TagType
}

// ProtoGlyphs represents protobuf tags for Glyphs
type ProtoGlyphs struct {
	Stacks pbf.TagType
}

// Proto represents all font protobuf definitions
type Proto struct {
	Glyph     ProtoGlyph
	Fontstack ProtoFontstack
	Glyphs    ProtoGlyphs
}

var (
	// Default font proto definition
	FontProto = Proto{
		Glyph: ProtoGlyph{
			ID:      GlyphID,
			Bitmap:  GlyphBitmap,
			Width:   GlyphWidth,
			Height:  GlyphHeight,
			Left:    GlyphLeft,
			Top:     GlyphTop,
			Advance: GlyphAdvance,
		},
		Fontstack: ProtoFontstack{
			Name:   FontstackName,
			Range:  FontstackRange,
			Glyphs: FontstackGlyphs,
		},
		Glyphs: ProtoGlyphs{
			Stacks: GlyphsStacks,
		},
	}
)

// Glyph represents a single glyph with metrics and optional SDF bitmap.
type Glyph struct {
	ID      uint32
	Bitmap  []byte
	Width   uint32
	Height  uint32
	Left    int32
	Top     int32
	Advance uint32
}

// Fontstack represents a font stack with name, range and associated glyphs.
type Fontstack struct {
	Name   string
	Range  string
	Glyphs []*Glyph
}

// Glyphs contains multiple font stacks.
type Glyphs struct {
	Stacks []*Fontstack
}

// ReadPBF reads a Glyph from PBF data.
func (g *Glyph) ReadPBF(r *pbf.Reader) error {
	// Track the start position for safety
	startPos := r.Pos
	maxBytesToRead := 1024 * 1024 // 1MB limit to prevent infinite loops

	for r.Pos < r.Length && r.Pos-startPos < maxBytesToRead {
		tag, typ := r.ReadTag()
		if tag == 0 {
			break
		}

		switch tag {
		case FontProto.Glyph.ID: // id
			if typ != pbf.Varint {
				skipField(r, typ)
			} else {
				g.ID = uint32(r.ReadUInt32())
			}
		case FontProto.Glyph.Bitmap: // bitmap
			if typ != pbf.Bytes {
				skipField(r, typ)
			} else {
				g.Bitmap = []byte(r.ReadString())
			}
		case FontProto.Glyph.Width: // width
			if typ != pbf.Varint {
				skipField(r, typ)
			} else {
				g.Width = uint32(r.ReadUInt32())
			}
		case FontProto.Glyph.Height: // height
			if typ != pbf.Varint {
				skipField(r, typ)
			} else {
				g.Height = uint32(r.ReadUInt32())
			}
		case FontProto.Glyph.Left: // left
			if typ != pbf.Fixed32 {
				skipField(r, typ)
			} else {
				g.Left = int32(r.ReadSFixed32())
			}
		case FontProto.Glyph.Top: // top
			if typ != pbf.Fixed32 {
				skipField(r, typ)
			} else {
				g.Top = int32(r.ReadSFixed32())
			}
		case FontProto.Glyph.Advance: // advance
			if typ != pbf.Varint {
				skipField(r, typ)
			} else {
				g.Advance = uint32(r.ReadUInt32())
			}
		default:
			return nil
		}
	}

	return nil
}

// WritePBF writes a Glyph to PBF data.
func (g *Glyph) WritePBF(w *pbf.Writer) error {
	w.WriteUInt32(pbf.TagType(FontProto.Glyph.ID), (g.ID))

	if len(g.Bitmap) > 0 {
		w.WriteString(pbf.TagType(FontProto.Glyph.Bitmap), string(g.Bitmap))
	} else {
		w.WriteString(pbf.TagType(FontProto.Glyph.Bitmap), "")
	}

	w.WriteUInt32(pbf.TagType(FontProto.Glyph.Width), (g.Width))
	w.WriteUInt32(pbf.TagType(FontProto.Glyph.Height), (g.Height))
	w.WriteSFixed32(pbf.TagType(FontProto.Glyph.Left), (g.Left))
	w.WriteSFixed32(pbf.TagType(FontProto.Glyph.Top), (g.Top))
	w.WriteUInt32(pbf.TagType(FontProto.Glyph.Advance), (g.Advance))
	return nil
}

// ReadPBF reads a Fontstack from PBF data.
func (f *Fontstack) ReadPBF(r *pbf.Reader) error {
	// Track the start position for safety
	startPos := r.Pos
	maxBytesToRead := 1024 * 1024 // 1MB limit to prevent infinite loops

	for r.Pos < r.Length && r.Pos-startPos < maxBytesToRead {
		tag, typ := r.ReadTag()
		if tag == 0 {
			break
		}

		switch tag {
		case FontProto.Fontstack.Name: // name
			if typ != pbf.Bytes {
				// Skip this field
				skipField(r, typ)
			} else {
				f.Name = r.ReadString()
			}
		case FontProto.Fontstack.Range: // range
			if typ != pbf.Bytes {
				skipField(r, typ)
			} else {
				f.Range = r.ReadString()
			}
		case FontProto.Fontstack.Glyphs: // glyphs (repeated)
			if typ != pbf.Bytes {
				skipField(r, typ)
			} else {
				// Read the glyph message
				size := r.ReadVarint()
				endpos := r.Pos + size

				// Ensure we don't read beyond the PBF data
				if endpos > r.Length {
					endpos = r.Length
				}

				// Create a slice of the PBF data for this glyph
				glyphData := r.Pbf[r.Pos:endpos]
				glyphReader := pbf.NewReader(glyphData)
				glyphReader.Length = len(glyphData)

				// Read the glyph
				g := &Glyph{}
				if err := g.ReadPBF(glyphReader); err != nil {
					return err
				}
				f.Glyphs = append(f.Glyphs, g)

				// Move the main reader to the end of this glyph data
				r.Pos = endpos
			}
		default:
			// Skip unknown fields
			skipField(r, typ)
		}
	}

	return nil
}

// Helper function to skip a field based on its type
func skipField(r *pbf.Reader, typ pbf.WireType) {
	switch typ {
	case pbf.Varint:
		r.ReadUInt64()
	case pbf.Bytes:
		length := r.ReadUInt32()
		// Ensure we don't read beyond the PBF data
		if r.Pos+int(length) > r.Length {
			length = uint32(r.Length - r.Pos)
		}
		r.Pos += int(length)
	case pbf.Fixed32:
		r.ReadUInt32()
	case pbf.Fixed64:
		r.ReadUInt64()
	default:
		// Handle unknown type by skipping
		// Try to skip 8 bytes as a default
		if r.Pos+8 <= r.Length {
			r.Pos += 8
		} else {
			r.Pos = r.Length // Move to end to exit loop
		}
	}
}

// WritePBF writes a Fontstack to PBF data.
func (f *Fontstack) WritePBF(w *pbf.Writer) error {
	w.WriteString(pbf.TagType(FontProto.Fontstack.Name), f.Name)
	w.WriteString(pbf.TagType(FontProto.Fontstack.Range), f.Range)

	// Write each glyph as a repeated field
	for _, glyph := range f.Glyphs {
		var err error
		w.WriteMessage(pbf.TagType(FontProto.Fontstack.Glyphs), func(w *pbf.Writer) {
			if e := glyph.WritePBF(w); e != nil {
				err = e
				return
			}
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// ReadPBF reads Glyphs from PBF data.
func (g *Glyphs) ReadPBF(r *pbf.Reader) error {
	// Track the start position for safety
	startPos := r.Pos
	maxBytesToRead := 1024 * 1024 // 1MB limit to prevent infinite loops

	for r.Pos < len(r.Pbf) && r.Pos-startPos < maxBytesToRead {
		tag, typ := r.ReadTag()
		if tag == 0 {
			break
		}

		switch tag {
		case FontProto.Glyphs.Stacks: // stacks (repeated)
			if typ != pbf.Bytes {
				skipField(r, typ)
			} else {
				// Read the fontstack message
				size := r.ReadVarint()
				endpos := r.Pos + size

				// Ensure we don't read beyond the PBF data
				if endpos > len(r.Pbf) {
					endpos = len(r.Pbf)
				}

				// Create a slice of the PBF data for this fontstack
				stackData := r.Pbf[r.Pos:endpos]
				stackReader := pbf.NewReader(stackData)
				stackReader.Length = len(stackData)

				// Read the fontstack
				stack := &Fontstack{}
				if err := stack.ReadPBF(stackReader); err != nil {
					return err
				}
				g.Stacks = append(g.Stacks, stack)

				// Move the main reader to the end of this fontstack data
				r.Pos = endpos
			}
		default:
			// Skip unknown fields
			skipField(r, typ)
		}
	}

	return nil
}

// WritePBF writes Glyphs to PBF data.
func (g *Glyphs) WritePBF(w *pbf.Writer) error {
	// Write each fontstack as a repeated field
	for _, stack := range g.Stacks {
		var err error
		w.WriteMessage(pbf.TagType(FontProto.Glyphs.Stacks), func(w *pbf.Writer) {
			if e := stack.WritePBF(w); e != nil {
				err = e
				return
			}
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// NewGlyph creates a new Glyph with the specified parameters.
func NewGlyph(id, width, height, advance uint32, left, top int32, bitmap []byte) *Glyph {
	return &Glyph{
		ID:      id,
		Bitmap:  bitmap,
		Width:   width,
		Height:  height,
		Left:    left,
		Top:     top,
		Advance: advance,
	}
}

// NewFontstack creates a new Fontstack with the specified name and range.
func NewFontstack(name, r string) *Fontstack {
	return &Fontstack{
		Name:   name,
		Range:  r,
		Glyphs: []*Glyph{},
	}
}

// NewGlyphs creates a new Glyphs container.
func NewGlyphs() *Glyphs {
	return &Glyphs{
		Stacks: []*Fontstack{},
	}
}
