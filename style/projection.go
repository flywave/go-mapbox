package style

type Projection struct {
	Name      string    `json:"name"`
	Center    []float64 `json:"center,omitempty"`
	Parallels []float64 `json:"parallels,omitempty"`
}
