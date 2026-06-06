package style

type Fog struct {
	Color         *ColorType `json:"color,omitempty"`
	HighColor     *ColorType `json:"high-color,omitempty"`
	HorizonBlend  interface{} `json:"horizon-blend,omitempty"`
	Range         []float64   `json:"range,omitempty"`
	SpaceColor    *ColorType  `json:"space-color,omitempty"`
	StarIntensity interface{} `json:"star-intensity,omitempty"`
	VerticalRange []float64   `json:"vertical-range,omitempty"`
}
