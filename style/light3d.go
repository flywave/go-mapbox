package style

type Light3D struct {
	ID         string           `json:"id"`
	Type       string           `json:"type,omitempty"`
	Properties *Light3DProperties `json:"properties,omitempty"`
}

type Light3DProperties struct {
	Color                *ColorType `json:"color,omitempty"`
	Intensity            interface{} `json:"intensity,omitempty"`
	Direction            []float64   `json:"direction,omitempty"`
	CastShadows          *bool      `json:"cast-shadows,omitempty"`
	ShadowIntensity      interface{} `json:"shadow-intensity,omitempty"`
	ShadowDrawBeforeLayer *string   `json:"shadow-draw-before-layer,omitempty"`
	ShadowQuality        *float64   `json:"shadow-quality,omitempty"`
}
