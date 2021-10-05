package mapboxglstyle

type Paint struct {
	BackgroundColor *ColorType                   `json:"background-color"`
	FillColor       *ColorType                   `json:"fill-color"` // for example "hsl(47, 13%, 86%)"
	FillOpacity     *NumberOrFunctionWrapperType `json:"fill-opacity"`
	FillAntialias   interface{}                  `json:"fill-antialias"` // bool or e.g. {"base": 1, "stops": [[0, false], [9, true]]}
	LineColor       *ColorType                   `json:"line-color"`
	LineDashArray   []float64                    `json:"line-dasharray"`
	LineGapWidth    *NumberOrFunctionWrapperType `json:"line-gap-width"`
	LineOpacity     *NumberOrFunctionWrapperType `json:"line-opacity"` // either float64, or like {"base": 1, "stops": [[11, 0], [16, 1]]}
	LineOffset      *float64                     `json:"line-offset"`
	LineWidth       *NumberOrFunctionWrapperType `json:"line-width"`
	// FillOutlineColor *ColorStopsType              `json:"fill-outline-color"`
	FillOutlineColor *ColorType  `json:"fill-outline-color"`
	FillTranslate    interface{} `json:"fill-translate"`
	FillPattern      interface{} `json:"fill-pattern"`
	TextColor        *ColorType  `json:"text-color"`
	TextHaloBlur     *float64    `json:"text-halo-blur"`
	TextHaloColor    *ColorType  `json:"text-halo-color"`
	TextHaloWidth    *float64    `json:"text-halo-width"`
}
