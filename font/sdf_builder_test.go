package fonts

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/flywave/freetype/truetype"
	"github.com/flywave/go-pbf"
	"github.com/stretchr/testify/require"
)

func builderFor(fontFamily string) *SDFBuilder {
	ttf, err := os.ReadFile("../data/fonts/" + fontFamily + ".ttf")
	if err != nil {
		panic(err)
	}

	font, err := truetype.Parse(ttf)

	if err != nil {
		panic(err)
	}

	return NewSDFBuilder(font, SDFBuilderOpt{FontSize: 24, Buffer: 3})
}

func TestSDFBuilder_Glyph(t *testing.T) {
	builder := builderFor("NotoSans-Regular")

	for i := 0; i < 255; i++ {
		g := builder.Glyph(rune(i))
		if g != nil {
			fmt.Printf("%s %d\n", strconv.Itoa(int(g.ID)), g.Top)
			img := DrawGlyph(g, true)
			os.MkdirAll("../data/fonts/data/NotoSans", os.ModePerm)
			SavePNG(fmt.Sprintf("../data/fonts/data/NotoSans/%d.png", i), img)
		}
	}
}

func TestSDFBuilder(t *testing.T) {
	t.Run("#Glyphs", func(t *testing.T) {
		builder := builderFor("NotoSans-Regular")

		for _, rng := range [][]int{
			{0, 255},
			{20224, 20479},
			{22784, 23039},
		} {
			s := builder.Glyphs(rng[0], rng[1])
			w := pbf.NewWriter()
			if err := s.WritePBF(w); err != nil {
				require.NoError(t, err)
			}
			bytes := w.Finish()
			os.WriteFile(fmt.Sprintf("../data/fonts/data/NotoSans/%d-%d.pbf", rng[0], rng[1]), bytes, os.ModePerm)
		}
	})
}
