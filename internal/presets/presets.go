package presets

import (
	_ "embed"
	"encoding/json"
)

//go:embed presets.json
var raw []byte

type Layout struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Section struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Layouts []Layout `json:"layouts"`
}

type Industry struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Tagline string   `json:"tagline"`
	Sub     []string `json:"sub"`
	Palette []string `json:"palette"` // primary, secondary, tertiary
}

type Font struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Fallback string `json:"fallback"`
}

type Presets struct {
	Industries []Industry `json:"industries"`
	Sections   []Section  `json:"sections"`
	Fonts      []Font     `json:"fonts"`
}

func Load() (*Presets, error) {
	var p Presets
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, err
	}
	return &p, nil
}
