package model

type Section struct {
	ID     string `json:"id"`
	Layout string `json:"layout"`
}

type Font struct {
	ID       string `json:"id"` // fontsource id, e.g. "playfair-display"
	Name     string `json:"name"`
	Fallback string `json:"fallback"`
}

type Config struct {
	Name           string    `json:"name"`
	Slug           string    `json:"slug"`
	PackageManager string    `json:"packageManager"`
	Industry       string    `json:"industry"`
	IndustryName   string    `json:"industryName"`
	Tagline        string    `json:"tagline"`
	SubIndustry    string    `json:"subIndustry"`
	Style          string    `json:"style"`
	Sections       []Section `json:"sections"`
	Colors         struct {
		Primary   string `json:"primary"`
		Secondary string `json:"secondary"`
		Tertiary  string `json:"tertiary"`
	} `json:"colors"`
	Fonts struct {
		Heading Font `json:"heading"`
		Body    Font `json:"body"`
	} `json:"fonts"`
	CTAs []string `json:"ctas"`
}
