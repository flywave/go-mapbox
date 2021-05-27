package recipe

type Recipe struct {
	Version int              `json:"version"`
	Layers  map[string]Layer `json:"layers"`
}

type Layer struct {
	Source  string `json:"source"`
	MinZoom uint   `json:"minzoom"`
	MaxZoom uint   `json:"maxzoom"`
}
