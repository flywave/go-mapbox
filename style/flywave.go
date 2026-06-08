package style

// FlywavePostEffects holds post-processing effects configuration (flywave extension).
type FlywavePostEffects struct {
	Bloom    *FlywaveBloom    `json:"bloom,omitempty"`
	Outline  *FlywaveOutline  `json:"outline,omitempty"`
	Vignette *FlywaveVignette `json:"vignette,omitempty"`
	Sepia    *FlywaveSepia    `json:"sepia,omitempty"`
}

type FlywaveBloom struct {
	Enabled                bool     `json:"enabled"`
	Strength               *float64 `json:"strength,omitempty"`
	Radius                 *float64 `json:"radius,omitempty"`
	Levels                 *int     `json:"levels,omitempty"`
	Inverted               *bool    `json:"inverted,omitempty"`
	IgnoreBackground       *bool    `json:"ignoreBackground,omitempty"`
	LuminancePassEnabled   *bool    `json:"luminancePassEnabled,omitempty"`
	LuminancePassThreshold *float64 `json:"luminancePassThreshold,omitempty"`
	LuminancePassSmoothing *float64 `json:"luminancePassSmoothing,omitempty"`
}

type FlywaveOutline struct {
	Enabled               bool    `json:"enabled"`
	GhostExtrudedPolygons bool    `json:"ghostExtrudedPolygons"`
	Thickness             float64 `json:"thickness"`
	Color                 string  `json:"color"`
}

type FlywaveVignette struct {
	Enabled  bool    `json:"enabled"`
	Offset   float64 `json:"offset"`
	Darkness float64 `json:"darkness"`
}

type FlywaveSepia struct {
	Enabled bool    `json:"enabled"`
	Amount  float64 `json:"amount"`
}

type FlywaveTextStyle struct {
	Name              *string  `json:"name,omitempty"`
	FontCatalogName   *string  `json:"fontCatalogName,omitempty"`
	FontName          *string  `json:"fontName,omitempty"`
	Size              *int     `json:"size,omitempty"`
	Color             *string  `json:"color,omitempty"`
	BackgroundColor   *string  `json:"backgroundColor,omitempty"`
	Opacity           *float64 `json:"opacity,omitempty"`
	BackgroundOpacity *float64 `json:"backgroundOpacity,omitempty"`
}

type FlywaveFontCatalog struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

type FlywaveImageTexture struct {
	Name    string   `json:"name"`
	Image   string   `json:"image"`
	Origin  *string  `json:"origin,omitempty"`
	XOffset *int     `json:"xOffset,omitempty"`
	YOffset *int     `json:"yOffset,omitempty"`
	Width   *int     `json:"width,omitempty"`
	Height  *int     `json:"height,omitempty"`
	FlipH   *bool    `json:"flipH,omitempty"`
	FlipV   *bool    `json:"flipV,omitempty"`
	Opacity *float64 `json:"opacity,omitempty"`
}

type FlywavePoiTableRef struct {
	Name              string `json:"name"`
	URL               string `json:"url"`
	UseAltNamesForKey bool   `json:"useAltNamesForKey"`
}

type FlywavePriority struct {
	Group    string  `json:"group"`
	Category *string `json:"category,omitempty"`
}

type FlywaveDefinitions map[string]interface{}


