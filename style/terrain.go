package style

type Terrain struct {
	Source       string       `json:"source"`
	Exaggeration *Expression  `json:"exaggeration,omitempty"`
}
