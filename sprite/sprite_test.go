package sprite

import (
	"image"
	"image/color"
	"testing"
)

type testTexture struct {
	id         string
	name       string
	width      int
	height     int
	pixelRatio int
	img        image.Image
}

func (t *testTexture) TextureId() *string        { return &t.id }
func (t *testTexture) TextureName() string        { return t.name }
func (t *testTexture) TextureWidth() int           { return t.width }
func (t *testTexture) TextureHeight() int          { return t.height }
func (t *testTexture) TexturePixelRatio() int      { return t.pixelRatio }
func (t *testTexture) TextureImage() image.Image   { return t.img }

func newTestTexture(name string, w, h, pr int) *testTexture {
	return &testTexture{
		id: name, name: name,
		width: w, height: h, pixelRatio: pr,
		img: image.NewUniform(color.RGBA{255, 0, 0, 255}),
	}
}

func TestGenerateSprite_NoRender(t *testing.T) {
	textures := []Texture{
		newTestTexture("a", 32, 32, 1),
		newTestTexture("b", 64, 64, 1),
	}
	set, img, err := GenerateSprite(textures, 1, false)
	if err != nil {
		t.Fatal(err)
	}
	if img != nil {
		t.Fatal("expected nil image when renderImage=false")
	}
	if len(set) != 2 {
		t.Fatalf("expected 2 sprites, got %d", len(set))
	}
	if set["a"] == nil || set["b"] == nil {
		t.Fatal("missing sprite entries")
	}
}

func TestGenerateSprite_WithRender(t *testing.T) {
	textures := []Texture{
		newTestTexture("a", 32, 32, 1),
		newTestTexture("b", 16, 16, 1),
	}
	set, img, err := GenerateSprite(textures, 1, true)
	if err != nil {
		t.Fatal(err)
	}
	if img == nil {
		t.Fatal("expected non-nil image when renderImage=true")
	}
	if len(set) != 2 {
		t.Fatalf("expected 2 sprites, got %d", len(set))
	}
	bounds := img.Bounds()
	if bounds.Dx() == 0 || bounds.Dy() == 0 {
		t.Fatal("sprite image has zero dimensions")
	}
}

func TestGenerateSprite_Empty(t *testing.T) {
	set, img, err := GenerateSprite(nil, 1, true)
	if err != nil {
		t.Fatal(err)
	}
	if set == nil {
		t.Fatal("expected non-nil set")
	}
	if img == nil {
		t.Fatal("expected non-nil image even for empty input")
	}
}

func TestGenerateSprite_PixelRatioScaling(t *testing.T) {
	// Texture at 2x pixel ratio, target 1x
	textures := []Texture{
		newTestTexture("a", 64, 64, 2),
	}
	set, _, err := GenerateSprite(textures, 1, false)
	if err != nil {
		t.Fatal(err)
	}
	ts := set["a"]
	if ts == nil {
		t.Fatal("missing sprite")
	}
	// Should be scaled down: 64/2 = 32
	if ts.Width != 32 || ts.Height != 32 {
		t.Fatalf("expected 32x32 after scale, got %dx%d", ts.Width, ts.Height)
	}
	if ts.PixelRatio != 1 {
		t.Fatalf("expected pixelRatio=1, got %d", ts.PixelRatio)
	}
}

func TestTextureMeta_ScaleTo(t *testing.T) {
	tm := &TextureMeta{Name: "test", Width: 200, Height: 100, PixelRatio: 2}
	scaled := tm.ScaleTo(1)
	if scaled.Width != 100 || scaled.Height != 50 {
		t.Fatalf("expected 100x50, got %dx%d", scaled.Width, scaled.Height)
	}
	if scaled.PixelRatio != 1 {
		t.Fatalf("PixelRatio = %d", scaled.PixelRatio)
	}
	// Scale up
	scaled2 := tm.ScaleTo(4)
	if scaled2.Width != 400 || scaled2.Height != 200 {
		t.Fatalf("expected 400x200, got %dx%d", scaled2.Width, scaled2.Height)
	}
}

func TestTextureSprite_String(t *testing.T) {
	ts := &TextureSprite{
		TextureMeta: &TextureMeta{Name: "test", Width: 32, Height: 32, PixelRatio: 1},
		X: 10, Y: 20,
	}
	s := ts.String()
	if s == "" {
		t.Fatal("expected non-empty string")
	}
}

func TestTextureMeta_String(t *testing.T) {
	tm := &TextureMeta{Name: "icon", Width: 64, Height: 64, PixelRatio: 2}
	s := tm.String()
	if s == "" {
		t.Fatal("expected non-empty string")
	}
}
