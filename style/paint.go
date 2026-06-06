package style

type Paint struct {
	// Background
	BackgroundColor          *ColorType `json:"background-color,omitempty"`
	BackgroundEmissiveStrength interface{} `json:"background-emissive-strength,omitempty"`
	BackgroundOpacity        interface{} `json:"background-opacity,omitempty"`
	BackgroundPattern        interface{} `json:"background-pattern,omitempty"`
	BackgroundPitchAlignment string      `json:"background-pitch-alignment,omitempty"`

	// Circle
	CircleBlur             interface{} `json:"circle-blur,omitempty"`
	CircleColor            *ColorType  `json:"circle-color,omitempty"`
	CircleEmissiveStrength interface{} `json:"circle-emissive-strength,omitempty"`
	CircleOpacity          interface{} `json:"circle-opacity,omitempty"`
	CirclePitchAlignment   string      `json:"circle-pitch-alignment,omitempty"`
	CirclePitchScale       string      `json:"circle-pitch-scale,omitempty"`
	CircleRadius           interface{} `json:"circle-radius,omitempty"`
	CircleStrokeColor      *ColorType  `json:"circle-stroke-color,omitempty"`
	CircleStrokeOpacity    interface{} `json:"circle-stroke-opacity,omitempty"`
	CircleStrokeWidth      interface{} `json:"circle-stroke-width,omitempty"`
	CircleTranslate        []float64   `json:"circle-translate,omitempty"`
	CircleTranslateAnchor  string      `json:"circle-translate-anchor,omitempty"`

	// Fill
	FillAntialias    interface{} `json:"fill-antialias,omitempty"`
	FillColor        *ColorType  `json:"fill-color,omitempty"`
	FillEmissiveStrength interface{} `json:"fill-emissive-strength,omitempty"`
	FillOpacity      interface{} `json:"fill-opacity,omitempty"`
	FillOutlineColor *ColorType  `json:"fill-outline-color,omitempty"`
	FillPattern      interface{} `json:"fill-pattern,omitempty"`
	FillPatternCrossFade interface{} `json:"fill-pattern-cross-fade,omitempty"`
	FillTranslate    []float64   `json:"fill-translate,omitempty"`
	FillTranslateAnchor string   `json:"fill-translate-anchor,omitempty"`
	FillZOffset      interface{} `json:"fill-z-offset,omitempty"`

	// Fill-Extrusion
	FillExtrusionAmbientOcclusionGroundAttenuation interface{} `json:"fill-extrusion-ambient-occlusion-ground-attenuation,omitempty"`
	FillExtrusionAmbientOcclusionGroundRadius      interface{} `json:"fill-extrusion-ambient-occlusion-ground-radius,omitempty"`
	FillExtrusionAmbientOcclusionIntensity         interface{} `json:"fill-extrusion-ambient-occlusion-intensity,omitempty"`
	FillExtrusionAmbientOcclusionRadius            interface{} `json:"fill-extrusion-ambient-occlusion-radius,omitempty"`
	FillExtrusionAmbientOcclusionWallRadius        interface{} `json:"fill-extrusion-ambient-occlusion-wall-radius,omitempty"`
	FillExtrusionBase                              interface{} `json:"fill-extrusion-base,omitempty"`
	FillExtrusionBaseAlignment                     string      `json:"fill-extrusion-base-alignment,omitempty"`
	FillExtrusionCastShadows                       *bool       `json:"fill-extrusion-cast-shadows,omitempty"`
	FillExtrusionColor                             *ColorType  `json:"fill-extrusion-color,omitempty"`
	FillExtrusionCutoffFadeRange                   interface{} `json:"fill-extrusion-cutoff-fade-range,omitempty"`
	FillExtrusionEmissiveStrength                  interface{} `json:"fill-extrusion-emissive-strength,omitempty"`
	FillExtrusionFloodLightColor                   *ColorType  `json:"fill-extrusion-flood-light-color,omitempty"`
	FillExtrusionFloodLightGroundAttenuation       interface{} `json:"fill-extrusion-flood-light-ground-attenuation,omitempty"`
	FillExtrusionFloodLightGroundRadius            interface{} `json:"fill-extrusion-flood-light-ground-radius,omitempty"`
	FillExtrusionFloodLightIntensity               interface{} `json:"fill-extrusion-flood-light-intensity,omitempty"`
	FillExtrusionFloodLightWallRadius              interface{} `json:"fill-extrusion-flood-light-wall-radius,omitempty"`
	FillExtrusionHeight                            interface{} `json:"fill-extrusion-height,omitempty"`
	FillExtrusionHeightAlignment                   string      `json:"fill-extrusion-height-alignment,omitempty"`
	FillExtrusionLineWidth                         interface{} `json:"fill-extrusion-line-width,omitempty"`
	FillExtrusionOpacity                           interface{} `json:"fill-extrusion-opacity,omitempty"`
	FillExtrusionPattern                           interface{} `json:"fill-extrusion-pattern,omitempty"`
	FillExtrusionPatternCrossFade                  interface{} `json:"fill-extrusion-pattern-cross-fade,omitempty"`
	FillExtrusionRoundedRoof                       *bool       `json:"fill-extrusion-rounded-roof,omitempty"`
	FillExtrusionTranslate                         []float64   `json:"fill-extrusion-translate,omitempty"`
	FillExtrusionTranslateAnchor                   string      `json:"fill-extrusion-translate-anchor,omitempty"`
	FillExtrusionVerticalGradient                  *bool       `json:"fill-extrusion-vertical-gradient,omitempty"`
	FillExtrusionVerticalScale                     interface{} `json:"fill-extrusion-vertical-scale,omitempty"`

	// Heatmap
	HeatmapColor    *ColorType  `json:"heatmap-color,omitempty"`
	HeatmapIntensity interface{} `json:"heatmap-intensity,omitempty"`
	HeatmapOpacity  interface{} `json:"heatmap-opacity,omitempty"`
	HeatmapRadius   interface{} `json:"heatmap-radius,omitempty"`
	HeatmapWeight   interface{} `json:"heatmap-weight,omitempty"`

	// Hillshade
	HillshadeAccentColor          *ColorType `json:"hillshade-accent-color,omitempty"`
	HillshadeEmissiveStrength     interface{} `json:"hillshade-emissive-strength,omitempty"`
	HillshadeExaggeration         interface{} `json:"hillshade-exaggeration,omitempty"`
	HillshadeHighlightColor       *ColorType `json:"hillshade-highlight-color,omitempty"`
	HillshadeIlluminationAnchor   string      `json:"hillshade-illumination-anchor,omitempty"`
	HillshadeIlluminationDirection interface{} `json:"hillshade-illumination-direction,omitempty"`
	HillshadeShadowColor          *ColorType `json:"hillshade-shadow-color,omitempty"`

	// Line
	LineBlur              interface{} `json:"line-blur,omitempty"`
	LineColor             *ColorType  `json:"line-color,omitempty"`
	LineDashArray         []float64   `json:"line-dasharray,omitempty"`
	LineEmissiveStrength  interface{} `json:"line-emissive-strength,omitempty"`
	LineGapWidth          interface{} `json:"line-gap-width,omitempty"`
	LineGradient          *ColorType  `json:"line-gradient,omitempty"`
	LineOcclusionOpacity  interface{} `json:"line-occlusion-opacity,omitempty"`
	LineOffset            interface{} `json:"line-offset,omitempty"`
	LineOpacity           interface{} `json:"line-opacity,omitempty"`
	LinePattern           interface{} `json:"line-pattern,omitempty"`
	LinePatternCrossFade  interface{} `json:"line-pattern-cross-fade,omitempty"`
	LineTranslate         []float64   `json:"line-translate,omitempty"`
	LineTranslateAnchor   string      `json:"line-translate-anchor,omitempty"`
	LineTrimColor         *ColorType  `json:"line-trim-color,omitempty"`
	LineTrimFadeRange     []float64   `json:"line-trim-fade-range,omitempty"`
	LineTrimOffset        []float64   `json:"line-trim-offset,omitempty"`
	LineWidth             interface{} `json:"line-width,omitempty"`

	// Model
	ModelAmbientOcclusionIntensity          interface{} `json:"model-ambient-occlusion-intensity,omitempty"`
	ModelCastShadows                        *bool       `json:"model-cast-shadows,omitempty"`
	ModelColor                              *ColorType  `json:"model-color,omitempty"`
	ModelColorMixIntensity                  interface{} `json:"model-color-mix-intensity,omitempty"`
	ModelCutoffFadeRange                    interface{} `json:"model-cutoff-fade-range,omitempty"`
	ModelElevationReference                 string      `json:"model-elevation-reference,omitempty"`
	ModelEmissiveStrength                   interface{} `json:"model-emissive-strength,omitempty"`
	ModelHeightBasedEmissiveStrengthMultiplier []float64 `json:"model-height-based-emissive-strength-multiplier,omitempty"`
	ModelOpacity                            interface{} `json:"model-opacity,omitempty"`
	ModelReceiveShadows                     *bool       `json:"model-receive-shadows,omitempty"`
	ModelRotation                           []float64   `json:"model-rotation,omitempty"`
	ModelRoughness                          interface{} `json:"model-roughness,omitempty"`
	ModelScale                              []float64   `json:"model-scale,omitempty"`
	ModelTranslation                        []float64   `json:"model-translation,omitempty"`
	ModelType                               string      `json:"model-type,omitempty"`

	// Raster
	RasterArrayBand              string      `json:"raster-array-band,omitempty"`
	RasterBrightnessMax          interface{} `json:"raster-brightness-max,omitempty"`
	RasterBrightnessMin          interface{} `json:"raster-brightness-min,omitempty"`
	RasterColorMix               []float64   `json:"raster-color-mix,omitempty"`
	RasterColorRange             []float64   `json:"raster-color-range,omitempty"`
	RasterContrast               interface{} `json:"raster-contrast,omitempty"`
	RasterElevation              interface{} `json:"raster-elevation,omitempty"`
	RasterEmissiveStrength       interface{} `json:"raster-emissive-strength,omitempty"`
	RasterHueRotate              interface{} `json:"raster-hue-rotate,omitempty"`
	RasterOpacity                interface{} `json:"raster-opacity,omitempty"`
	RasterParticleElevation      []float64   `json:"raster-particle-elevation,omitempty"`
	RasterParticleFadeOpacityFactor interface{} `json:"raster-particle-fade-opacity-factor,omitempty"`
	RasterParticleSpeedFactor    interface{} `json:"raster-particle-speed-factor,omitempty"`
	RasterResampling             string      `json:"raster-resampling,omitempty"`
	RasterSaturation             interface{} `json:"raster-saturation,omitempty"`

	// Sky
	SkyAtmosphereColor       *ColorType `json:"sky-atmosphere-color,omitempty"`
	SkyAtmosphereHaloColor   *ColorType `json:"sky-atmosphere-halo-color,omitempty"`
	SkyAtmosphereSun         []float64  `json:"sky-atmosphere-sun,omitempty"`
	SkyAtmosphereSunIntensity interface{} `json:"sky-atmosphere-sun-intensity,omitempty"`
	SkyGradientCenter        []float64  `json:"sky-gradient-center,omitempty"`
	SkyGradientRadius        interface{} `json:"sky-gradient-radius,omitempty"`
	SkyOpacity               interface{} `json:"sky-opacity,omitempty"`
	SkyType                  string     `json:"sky-type,omitempty"`

	// Symbol - Icon
	IconColor              *ColorType `json:"icon-color,omitempty"`
	IconCrossFade          interface{} `json:"icon-cross-fade,omitempty"`
	IconEmissiveStrength   interface{} `json:"icon-emissive-strength,omitempty"`
	IconHaloBlur           interface{} `json:"icon-halo-blur,omitempty"`
	IconHaloColor          *ColorType `json:"icon-halo-color,omitempty"`
	IconHaloWidth          interface{} `json:"icon-halo-width,omitempty"`
	IconOcclusionOpacity   interface{} `json:"icon-occlusion-opacity,omitempty"`
	IconOpacity            interface{} `json:"icon-opacity,omitempty"`
	IconTranslate          []float64  `json:"icon-translate,omitempty"`
	IconTranslateAnchor    string     `json:"icon-translate-anchor,omitempty"`

	// Symbol - Text
	TextColor            *ColorType `json:"text-color,omitempty"`
	TextEmissiveStrength interface{} `json:"text-emissive-strength,omitempty"`
	TextHaloBlur         interface{} `json:"text-halo-blur,omitempty"`
	TextHaloColor        *ColorType `json:"text-halo-color,omitempty"`
	TextHaloWidth        interface{} `json:"text-halo-width,omitempty"`
	TextOcclusionOpacity interface{} `json:"text-occlusion-opacity,omitempty"`
	TextOpacity          interface{} `json:"text-opacity,omitempty"`
	TextTranslate        []float64  `json:"text-translate,omitempty"`
	TextTranslateAnchor  string     `json:"text-translate-anchor,omitempty"`

	// Symbol - combined
	SymbolZOffset interface{} `json:"symbol-z-offset,omitempty"`
}
