package mapboxglstyle

type Layout struct {
	Visibility            string                       `json:"visibility"`
	LineCap               string                       `json:"line-cap"`
	LineJoin              string                       `json:"line-join"`
	TextField             string                       `json:"text-field"`
	TextFont              []string                     `json:"text-font"`
	TextSize              *NumberOrFunctionWrapperType `json:"text-size"`        // float64 or {"base": 1.4, "stops": [[10, 8], [20, 14]]}
	SymbolPlacement       interface{}                  `json:"symbol-placement"` // string or {"base": 1, "stops": [[10, "point"], [11, "line"]]}
	TextLetterSpacing     float64                      `json:"text-letter-spacing"`
	TextRotationAlignment string                       `json:"text-rotation-alignment"`
	TextTransform         string                       `json:"text-transform"`
	IconSize              *NumberOrFunctionWrapperType `json:"icon-size"`
	TextAnchor            string                       `json:"text-anchor"`
	TextMaxWidth          float64                      `json:"text-max-width"`
	TextOffset            []float64                    `json:"text-offset"`
	SymbolSpacing         float64                      `json:"symbol-spacing"`
	IconImage             string                       `json:"icon-image"`
	TextPadding           float64                      `json:"text-padding"`
}
