package style

type Layout struct {
	// Common
	Visibility string `json:"visibility,omitempty"`

	// Circle
	CircleElevationReference string `json:"circle-elevation-reference,omitempty"`
	CircleSortKey            interface{} `json:"circle-sort-key,omitempty"`

	// Clip
	ClipLayerScope  []string `json:"clip-layer-scope,omitempty"`
	ClipLayerTypes  []string `json:"clip-layer-types,omitempty"`

	// Fill
	FillSortKey interface{} `json:"fill-sort-key,omitempty"`

	// Fill-Extrusion
	FillExtrusionEdgeRadius interface{} `json:"fill-extrusion-edge-radius,omitempty"`

	// Line
	LineCap                 string      `json:"line-cap,omitempty"`
	LineCrossSlope          interface{} `json:"line-cross-slope,omitempty"`
	LineElevationGroundScale interface{} `json:"line-elevation-ground-scale,omitempty"`
	LineElevationReference  string      `json:"line-elevation-reference,omitempty"`
	LineJoin                string      `json:"line-join,omitempty"`
	LineMiterLimit          interface{} `json:"line-miter-limit,omitempty"`
	LineRoundLimit          interface{} `json:"line-round-limit,omitempty"`
	LineSortKey             interface{} `json:"line-sort-key,omitempty"`
	LineZOffset             interface{} `json:"line-z-offset,omitempty"`

	// Model
	ModelAllowDensityReduction *bool   `json:"model-allow-density-reduction,omitempty"`
	ModelID                    string  `json:"model-id,omitempty"`

	// Symbol - general
	SymbolAvoidEdges     *bool       `json:"symbol-avoid-edges,omitempty"`
	SymbolPlacement      interface{} `json:"symbol-placement,omitempty"`
	SymbolSpacing        interface{} `json:"symbol-spacing,omitempty"`
	SymbolZElevate       *bool       `json:"symbol-z-elevate,omitempty"`
	SymbolZOrder         string      `json:"symbol-z-order,omitempty"`

	// Symbol - Icon
	IconAllowOverlap       *bool       `json:"icon-allow-overlap,omitempty"`
	IconAnchor             string      `json:"icon-anchor,omitempty"`
	IconIgnorePlacement    *bool       `json:"icon-ignore-placement,omitempty"`
	IconImage              interface{} `json:"icon-image,omitempty"`
	IconKeepUpright        *bool       `json:"icon-keep-upright,omitempty"`
	IconOffset             []float64   `json:"icon-offset,omitempty"`
	IconOptional           *bool       `json:"icon-optional,omitempty"`
	IconPadding            interface{} `json:"icon-padding,omitempty"`
	IconPitchAlignment     string      `json:"icon-pitch-alignment,omitempty"`
	IconRotate             interface{} `json:"icon-rotate,omitempty"`
	IconRotationAlignment  string      `json:"icon-rotation-alignment,omitempty"`
	IconSize               interface{} `json:"icon-size,omitempty"`
	IconTextFit            string      `json:"icon-text-fit,omitempty"`
	IconTextFitPadding     []float64   `json:"icon-text-fit-padding,omitempty"`

	// Symbol - Text
	TextAllowOverlap       *bool       `json:"text-allow-overlap,omitempty"`
	TextAnchor             string      `json:"text-anchor,omitempty"`
	TextField              interface{} `json:"text-field,omitempty"`
	TextFont               []string    `json:"text-font,omitempty"`
	TextIgnorePlacement    *bool       `json:"text-ignore-placement,omitempty"`
	TextJustify            string      `json:"text-justify,omitempty"`
	TextKeepUpright        *bool       `json:"text-keep-upright,omitempty"`
	TextLetterSpacing      interface{} `json:"text-letter-spacing,omitempty"`
	TextLineHeight         interface{} `json:"text-line-height,omitempty"`
	TextMaxAngle           interface{} `json:"text-max-angle,omitempty"`
	TextMaxWidth           interface{} `json:"text-max-width,omitempty"`
	TextOffset             []float64   `json:"text-offset,omitempty"`
	TextOptional           *bool       `json:"text-optional,omitempty"`
	TextPadding            interface{} `json:"text-padding,omitempty"`
	TextPitchAlignment     string      `json:"text-pitch-alignment,omitempty"`
	TextRadialOffset       interface{} `json:"text-radial-offset,omitempty"`
	TextRotate             interface{} `json:"text-rotate,omitempty"`
	TextRotationAlignment  string      `json:"text-rotation-alignment,omitempty"`
	TextSize               interface{} `json:"text-size,omitempty"`
	TextTransform          string      `json:"text-transform,omitempty"`
	TextVariableAnchor     []string    `json:"text-variable-anchor,omitempty"`
	TextWritingMode        []string    `json:"text-writing-mode,omitempty"`
}
