package style

type Paint struct {
	BackgroundColor  *ColorType                   `json:"background-color,omitempty"`
	FillColor        *ColorType                   `json:"fill-color,omitempty"`
	FillOpacity      *NumberOrFunctionWrapperType `json:"fill-opacity,omitempty"`
	FillAntialias    interface{}                  `json:"fill-antialias,omitempty"`
	LineColor        *ColorType                   `json:"line-color,omitempty"`
	LineDashArray    []float64                    `json:"line-dasharray,omitempty"`
	LineGapWidth     *NumberOrFunctionWrapperType `json:"line-gap-width,omitempty"`
	LineOpacity      *NumberOrFunctionWrapperType `json:"line-opacity,omitempty"`
	LineOffset       *float64                     `json:"line-offset,omitempty"`
	LineWidth        *NumberOrFunctionWrapperType `json:"line-width,omitempty"`
	FillOutlineColor *ColorType                   `json:"fill-outline-color,omitempty"`
	FillTranslate    interface{}                  `json:"fill-translate,omitempty"`
	FillPattern      interface{}                  `json:"fill-pattern,omitempty"`
	TextColor        *ColorType                   `json:"text-color,omitempty"`
	TextHaloBlur     *float64                     `json:"text-halo-blur,omitempty"`
	TextHaloColor    *ColorType                   `json:"text-halo-color,omitempty"`
	TextHaloWidth    *float64                     `json:"text-halo-width,omitempty"`
}
