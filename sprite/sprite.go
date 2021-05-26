package sprite

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"sort"

	"golang.org/x/image/draw"
)

func GenerateSprite(textures []Texture, pixelRatio int, renderImage bool) (map[string]*TextureSprite, image.Image, error) {
	set := map[string]*TextureSprite{}
	chDone := make(chan struct{})
	defer close(chDone)

	blocks := make([]*Node, 0)

	for i := range textures {
		t := textures[i]

		img := t.TextureImage()

		tm := &TextureMeta{
			Name:       t.TextureName(),
			PixelRatio: t.TexturePixelRatio(),
			Width:      t.TextureWidth(),
			Height:     t.TextureHeight(),
		}

		ts := &TextureSprite{
			TextureMeta:       tm.ScaleTo(pixelRatio),
			OriginTextureMeta: tm,
			Image:             img,
		}

		blocks = append(blocks, NewNode(ts.Name, ts.Width, ts.Height))

		set[t.TextureName()] = ts
	}

	sort.Slice(blocks, func(i, j int) bool {
		return blocks[j].Width*blocks[j].Height < blocks[i].Width*blocks[i].Height
	})

	(&GrowingPacker{}).Fit(blocks)

	w := 0
	h := 0

	for _, b := range blocks {
		set[b.Key].X = b.X
		set[b.Key].Y = b.Y

		if renderImage {
			mayMaxX := b.Width + b.X
			mayMaxY := b.Height + b.Y

			if mayMaxX > w {
				w = mayMaxX
			}
			if mayMaxY > h {
				h = mayMaxY
			}
		}
	}

	if renderImage {
		if w == 0 {
			w = 1
		}

		if h == 0 {
			h = 1
		}

		dst := image.NewRGBA(image.Rect(0, 0, w, h))

		for _, ts := range set {
			img := ts.Image

			if ts.OriginTextureMeta != nil && ts.OriginTextureMeta.PixelRatio != pixelRatio {
				d := image.NewRGBA(image.Rect(0, 0, ts.Width, ts.Height))
				draw.BiLinear.Scale(d, d.Bounds(), ts.Image, ts.Image.Bounds(), draw.Over, nil)
				img = d
			}

			draw.Draw(dst, dst.Bounds(), img, image.Pt(-ts.X, -ts.Y), draw.Src)
		}

		return set, dst, nil
	}

	return set, nil, nil
}

type Texture interface {
	TextureName() string
	TextureWidth() int
	TextureHeight() int
	TexturePixelRatio() int
	TextureImage() image.Image
}

type TextureSprite struct {
	*TextureMeta
	X                 int          `json:"x"`
	Y                 int          `json:"y"`
	OriginTextureMeta *TextureMeta `json:"-"`
	Image             image.Image  `json:"-"`
}

func (s TextureSprite) String() string {
	return fmt.Sprintf("(%d,%d) %s", s.X, s.Y, s.TextureMeta)
}

type TextureMeta struct {
	Name       string `json:"-"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	PixelRatio int    `json:"pixelRatio"`
}

func (t TextureMeta) ScaleTo(pixelRatio int) *TextureMeta {
	scale := float64(pixelRatio) / float64(t.PixelRatio)

	return &TextureMeta{
		Name:       t.Name,
		Width:      int(float64(t.Width) * scale),
		Height:     int(float64(t.Height) * scale),
		PixelRatio: pixelRatio,
	}
}

func (t TextureMeta) String() string {
	return fmt.Sprintf("%s#%d,%d@%dx", t.Name, t.Width, t.Height, t.PixelRatio)
}
