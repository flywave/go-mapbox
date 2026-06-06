package style

type Import struct {
	ID         string                 `json:"id"`
	URL        string                 `json:"url"`
	ColorTheme *ColorTheme            `json:"color-theme,omitempty"`
	Config     map[string]interface{} `json:"config,omitempty"`
	Data       *Style                 `json:"data,omitempty"`
}

type ColorTheme struct {
	Data string `json:"data,omitempty"`
}
