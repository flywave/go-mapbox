package style

type Paint struct {
	BackgroundColor  *ColorType                   `json:"background-color"`
	FillColor        *ColorType                   `json:"fill-color"`
	FillOpacity      *NumberOrFunctionWrapperType `json:"fill-opacity"`
	FillAntialias    interface{}                  `json:"fill-antialias"`
	LineColor        *ColorType                   `json:"line-color"`
	LineDashArray    []float64                    `json:"line-dasharray"`
	LineGapWidth     *NumberOrFunctionWrapperType `json:"line-gap-width"`
	LineOpacity      *NumberOrFunctionWrapperType `json:"line-opacity"`
	LineOffset       *float64                     `json:"line-offset"`
	LineWidth        *NumberOrFunctionWrapperType `json:"line-width"`
	FillOutlineColor *ColorType                   `json:"fill-outline-color"`
	FillTranslate    interface{}                  `json:"fill-translate"`
	FillPattern      interface{}                  `json:"fill-pattern"`
	TextColor        *ColorType                   `json:"text-color"`
	TextHaloBlur     *float64                     `json:"text-halo-blur"`
	TextHaloColor    *ColorType                   `json:"text-halo-color"`
	TextHaloWidth    *float64                     `json:"text-halo-width"`
}
