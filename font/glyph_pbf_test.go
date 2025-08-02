// Package fonts provides tests for glyph_pbf.go
package fonts

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/flywave/go-pbf"
)

func TestGlyphReadWritePBF(t *testing.T) {
	// Create a new glyph
	original := NewGlyph(65, 100, 120, 110, 10, -20, []byte{0x01, 0x02, 0x03, 0x04})

	// Write the glyph to PBF
	w := pbf.NewWriter()
	if err := original.WritePBF(w); err != nil {
		t.Fatalf("Failed to write glyph to PBF: %v", err)
	}
	buf := w.Finish()
	fmt.Printf("Glyph PBF buffer length: %d bytes\n", len(buf))

	// Read the glyph from PBF with timeout
	r := pbf.NewReader(buf)
	read := &Glyph{}

	// Set up timeout channel
	timeout := time.After(5 * time.Second)
	done := make(chan bool)
	errChan := make(chan error)

	go func() {
		fmt.Printf("Starting to read glyph PBF data...\n")
		startPos := r.Pos
		startTime := time.Now()

		err := read.ReadPBF(r)
		if err != nil {
			errChan <- err
			return
		}

		duration := time.Since(startTime)
		bytesRead := r.Pos - startPos
		fmt.Printf("Finished reading glyph PBF in %v, read %d bytes\n", duration, bytesRead)
		done <- true
	}()

	// Wait for completion or timeout
	select {
	case <-done:
		// Reading completed successfully
	case err := <-errChan:
		t.Fatalf("Failed to read glyph from PBF: %v", err)
	case <-timeout:
		t.Fatalf("Glyph ReadPBF timed out after 5 seconds - possible infinite loop")
	}

	// Verify the glyph data
	if read.ID != original.ID {
		t.Errorf("Expected ID %d, got %d", original.ID, read.ID)
	}
	if read.Width != original.Width {
		t.Errorf("Expected Width %d, got %d", original.Width, read.Width)
	}
	if read.Height != original.Height {
		t.Errorf("Expected Height %d, got %d", original.Height, read.Height)
	}
	if read.Advance != original.Advance {
		t.Errorf("Expected Advance %d, got %d", original.Advance, read.Advance)
	}
	if read.Left != original.Left {
		t.Errorf("Expected Left %d, got %d", original.Left, read.Left)
	}
	if read.Top != original.Top {
		t.Errorf("Expected Top %d, got %d", original.Top, read.Top)
	}
	if !bytes.Equal(read.Bitmap, original.Bitmap) {
		t.Errorf("Expected Bitmap %v, got %v", original.Bitmap, read.Bitmap)
	}
}

func TestFontstackReadWritePBF(t *testing.T) {
	// Create a fontstack with glyphs
	stack := NewFontstack("Arial", "0-255")
	stack.Glyphs = append(stack.Glyphs,
		NewGlyph(65, 100, 120, 110, 10, -20, []byte{0x01, 0x02}),
		NewGlyph(66, 90, 120, 100, 15, -22, []byte{0x03, 0x04}),
	)

	// Write the fontstack to PBF
	w := pbf.NewWriter()
	if err := stack.WritePBF(w); err != nil {
		t.Fatalf("Failed to write fontstack to PBF: %v", err)
	}
	buf := w.Finish()
	fmt.Printf("Fontstack PBF buffer length: %d bytes\n", len(buf))

	// Read the fontstack from PBF with timeout
	r := pbf.NewReader(buf)
	read := &Fontstack{}

	// Set up timeout channel
	timeout := time.After(5 * time.Second)
	done := make(chan bool)
	errChan := make(chan error)

	go func() {
		fmt.Printf("Starting to read fontstack PBF data...\n")
		startPos := r.Pos
		startTime := time.Now()

		err := read.ReadPBF(r)
		if err != nil {
			errChan <- err
			return
		}

		duration := time.Since(startTime)
		bytesRead := r.Pos - startPos
		fmt.Printf("Finished reading fontstack PBF in %v, read %d bytes\n", duration, bytesRead)
		done <- true
	}()

	// Wait for completion or timeout
	select {
	case <-done:
		// Reading completed successfully
	case err := <-errChan:
		t.Fatalf("Failed to read fontstack from PBF: %v", err)
	case <-timeout:
		t.Fatalf("Fontstack ReadPBF timed out after 5 seconds - possible infinite loop")
	}

	// Verify the fontstack data
	if read.Name != stack.Name {
		t.Errorf("Expected Name %s, got %s", stack.Name, read.Name)
	}
	if read.Range != stack.Range {
		t.Errorf("Expected Range %s, got %s", stack.Range, read.Range)
	}
	if len(read.Glyphs) != len(stack.Glyphs) {
		t.Errorf("Expected %d glyphs, got %d", len(stack.Glyphs), len(read.Glyphs))
	} else {
		for i := range stack.Glyphs {
			origGlyph := stack.Glyphs[i]
			readGlyph := read.Glyphs[i]

			if readGlyph.ID != origGlyph.ID {
				t.Errorf("Glyph %d: Expected ID %d, got %d", i, origGlyph.ID, readGlyph.ID)
			}
			if readGlyph.Width != origGlyph.Width {
				t.Errorf("Glyph %d: Expected Width %d, got %d", i, origGlyph.Width, readGlyph.Width)
			}
			// Additional checks for other fields...
		}
	}
}

func TestGlyphsReadWritePBF(t *testing.T) {
	// Create glyphs with multiple fontstacks
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

	// Write the glyphs to PBF
	w := pbf.NewWriter()
	if err := glyphs.WritePBF(w); err != nil {
		t.Fatalf("Failed to write glyphs to PBF: %v", err)
	}
	buf := w.Finish()
	fmt.Printf("PBF buffer length: %d bytes\n", len(buf))

	// Read the glyphs from PBF with timeout
	r := pbf.NewReader(buf)
	read := &Glyphs{}

	// Set up timeout channel
	timeout := time.After(5 * time.Second)
	done := make(chan bool)
	errChan := make(chan error)

	go func() {
		fmt.Printf("Starting to read PBF data...\n")
		startPos := r.Pos
		startTime := time.Now()

		err := read.ReadPBF(r)
		if err != nil {
			errChan <- err
			return
		}

		duration := time.Since(startTime)
		bytesRead := r.Pos - startPos
		fmt.Printf("Finished reading PBF in %v, read %d bytes\n", duration, bytesRead)
		done <- true
	}()

	// Wait for completion or timeout
	select {
	case <-done:
		// Reading completed successfully
	case err := <-errChan:
		t.Fatalf("Failed to read glyphs from PBF: %v", err)
	case <-timeout:
		t.Fatalf("ReadPBF timed out after 5 seconds - possible infinite loop")
	}

	// Verify the glyphs data
	if len(read.Stacks) != len(glyphs.Stacks) {
		t.Errorf("Expected %d fontstacks, got %d", len(glyphs.Stacks), len(read.Stacks))
	} else {
		for i := range glyphs.Stacks {
			origStack := glyphs.Stacks[i]
			readStack := read.Stacks[i]

			if readStack.Name != origStack.Name {
				t.Errorf("Fontstack %d: Expected Name %s, got %s", i, origStack.Name, readStack.Name)
			}
			if readStack.Range != origStack.Range {
				t.Errorf("Fontstack %d: Expected Range %s, got %s", i, origStack.Range, readStack.Range)
			}
			if len(readStack.Glyphs) != len(origStack.Glyphs) {
				t.Errorf("Fontstack %d: Expected %d glyphs, got %d", i, len(origStack.Glyphs), len(readStack.Glyphs))
			} else {
				for j := range origStack.Glyphs {
					origGlyph := origStack.Glyphs[j]
					readGlyph := readStack.Glyphs[j]

					if readGlyph.ID != origGlyph.ID {
						t.Errorf("Fontstack %d, Glyph %d: Expected ID %d, got %d", i, j, origGlyph.ID, readGlyph.ID)
					}
					if readGlyph.Width != origGlyph.Width {
						t.Errorf("Fontstack %d, Glyph %d: Expected Width %d, got %d", i, j, origGlyph.Width, readGlyph.Width)
					}
					// Additional checks for other fields...
				}
			}
		}
	}
}
