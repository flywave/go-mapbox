package style

type SchemaOption struct {
	Default    interface{} `json:"default"`
	Array      *bool       `json:"array,omitempty"`
	MaxValue   *float64    `json:"maxValue,omitempty"`
	MinValue   *float64    `json:"minValue,omitempty"`
	StepValue  *float64    `json:"stepValue,omitempty"`
	Metadata   Metadata    `json:"metadata,omitempty"`
	Type       string      `json:"type,omitempty"`
	Values     []interface{} `json:"values,omitempty"`
}
