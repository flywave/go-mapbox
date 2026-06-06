package style

type Appearance struct {
	Name       string                 `json:"name,omitempty"`
	Condition  interface{}            `json:"condition,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}
