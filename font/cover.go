package fonts

import (
	"bytes"
	"encoding/json"

	"github.com/flywave/freetype/truetype"
)

func difference(a, b []uint32) (diff []uint32) {
	m := make(map[uint32]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

type LanguageCoverage struct {
	Name      string `json:"name"`
	Id        string `json:"id"`
	Total     int    `json:"total"`
	Count     int    `json:"count"`
	Coverages []int  `json:"coverages"`
}

type CodePoints struct {
	ExemplarCharacters []uint32 `json:"exemplarCharacters"`
	Auxiliary          []uint32 `json:"auxiliary"`
	Index              []uint32 `json:"index"`
	Punctuation        []uint32 `json:"punctuation"`
}

type LanguageCodePoints struct {
	Id         string     `json:"id"`
	Name       string     `json:"name"`
	Codepoints CodePoints `json:"codepoints"`
}

type FontConverage struct {
	languages []*LanguageCodePoints
}

func (c *FontConverage) coverage(pointsArray [][]uint32) ([]*LanguageCoverage, error) {
	lang_coverages := make([]*LanguageCoverage, len(c.languages))

	for i, language := range c.languages {
		exemplarCharacters := language.Codepoints.ExemplarCharacters
		charactersLeft := exemplarCharacters

		coverages := make([]int, len(pointsArray))
		for i, points := range pointsArray {
			nextValue := difference(charactersLeft, points)
			covered := len(charactersLeft) - len(nextValue)
			charactersLeft = nextValue

			coverages[i] = covered
		}

		sum := 0

		for _, c := range coverages {
			sum += c
		}

		lang_coverages[i] = &LanguageCoverage{
			Name:      language.Name,
			Id:        language.Id,
			Total:     len(exemplarCharacters),
			Count:     sum,
			Coverages: coverages,
		}
	}

	return lang_coverages, nil
}

func (c *FontConverage) Coverage(ttfs ...[]byte) ([]*LanguageCoverage, error) {
	points := make([][]uint32, len(ttfs))
	for i, ttf := range ttfs {
		font, err := truetype.Parse(ttf)
		if err != nil {
			return nil, err
		}
		unicode := font.Unicodes()

		points[i] = unicode
	}

	return c.coverage(points)
}

func (c *FontConverage) CoverageTTFS(ttfs ...*truetype.Font) ([]*LanguageCoverage, error) {
	points := make([][]uint32, len(ttfs))
	for i, font := range ttfs {
		unicode := font.Unicodes()

		points[i] = unicode
	}

	return c.coverage(points)
}

func NewFontCoverage(data []byte) (*FontConverage, error) {
	c := FontConverage{}
	decode := json.NewDecoder(bytes.NewReader(data))

	err := decode.Decode(&c.languages)

	if err != nil {
		return nil, err
	}

	return &c, nil
}
