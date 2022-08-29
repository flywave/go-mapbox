package style

type Layout struct {
	Visibility            string                       `json:"visibility,omitempty"`
	LineCap               string                       `json:"line-cap,omitempty"`
	LineJoin              string                       `json:"line-join,omitempty"`
	TextField             string                       `json:"text-field,omitempty"`
	TextFont              []string                     `json:"text-font,omitempty"`
	TextSize              *NumberOrFunctionWrapperType `json:"text-size,omitempty"`
	SymbolPlacement       interface{}                  `json:"symbol-placement,omitempty"`
	TextLetterSpacing     float64                      `json:"text-letter-spacing,omitempty"`
	TextRotationAlignment string                       `json:"text-rotation-alignment,omitempty"`
	TextTransform         string                       `json:"text-transform,omitempty"`
	IconSize              *NumberOrFunctionWrapperType `json:"icon-size,omitempty"`
	TextAnchor            string                       `json:"text-anchor,omitempty"`
	TextMaxWidth          float64                      `json:"text-max-width,omitempty"`
	TextOffset            []float64                    `json:"text-offset,omitempty"`
	SymbolSpacing         float64                      `json:"symbol-spacing,omitempty"`
	IconImage             string                       `json:"icon-image,omitempty"`
	TextPadding           float64                      `json:"text-padding,omitempty"`
}
