package sprite

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseStyleFile(t *testing.T) {
	for _, s := range []string{
		"sprite@2x.json",
		"sprite.json",
		"sprite@2x.png",
	} {
		texture, err := ParseStyleFile(s)
		require.NoError(t, err)
		require.Equal(t, s, texture.String())
	}
}
