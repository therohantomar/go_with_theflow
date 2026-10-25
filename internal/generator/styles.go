package generator

import "fmt"

// tokens are the per-style design values written into global.css.
type tokens struct {
	CardR, FieldR, BtnR                string
	HTransform, HWeight, HTrack, HLead string
	H1, H2, SectionY                   string
	Surface, Ink, Muted                string
}

var light = tokens{Surface: "#ffffff", Ink: "#262626", Muted: "#525252"}

func tokensFor(style string) tokens {
	t := light
	switch style {
	case "bold":
		t.CardR, t.FieldR, t.BtnR = "0px", "0px", "0px"
		t.HTransform, t.HWeight, t.HTrack, t.HLead = "uppercase", "700", "-0.02em", "0.95"
		t.H1, t.H2, t.SectionY = "clamp(3rem, 9vw, 7.5rem)", "clamp(2rem, 5vw, 3.75rem)", "6rem"
	case "editorial":
		t.CardR, t.FieldR, t.BtnR = "0px", "0px", "0px"
		t.HTransform, t.HWeight, t.HTrack, t.HLead = "none", "400", "0.02em", "1.15"
		t.H1, t.H2, t.SectionY = "clamp(2rem, 4.5vw, 3.5rem)", "clamp(1.75rem, 3vw, 2.5rem)", "7rem"
		t.Surface, t.Muted = "#fbfaf8", "#6b6560"
	case "industrial":
		t.CardR, t.FieldR, t.BtnR = "6px", "4px", "4px"
		t.HTransform, t.HWeight, t.HTrack, t.HLead = "none", "700", "-0.01em", "1.1"
		t.H1, t.H2, t.SectionY = "clamp(2.25rem, 5vw, 3.75rem)", "clamp(1.75rem, 3.5vw, 2.5rem)", "4.5rem"
	case "creative":
		t.CardR, t.FieldR, t.BtnR = "24px", "12px", "9999px"
		t.HTransform, t.HWeight, t.HTrack, t.HLead = "none", "600", "-0.04em", "0.95"
		t.H1, t.H2, t.SectionY = "clamp(3rem, 8vw, 7rem)", "clamp(2rem, 5vw, 3.5rem)", "7rem"
		t.Surface, t.Ink, t.Muted = "#0b0b0d", "#ececec", "#a1a1aa"
	default: // soft
		t.CardR, t.FieldR, t.BtnR = "1rem", "0.5rem", "9999px"
		t.HTransform, t.HWeight, t.HTrack, t.HLead = "none", "700", "-0.01em", "1.15"
		t.H1, t.H2, t.SectionY = "clamp(2.25rem, 5vw, 3.75rem)", "clamp(1.75rem, 3.5vw, 2.5rem)", "5rem"
	}
	return t
}

// onColor picks readable text (white or near-black) for a hex background.
func onColor(hex string) string {
	var r, g, b int
	if len(hex) != 7 {
		return "#ffffff"
	}
	if _, err := fmt.Sscanf(hex[1:], "%02x%02x%02x", &r, &g, &b); err != nil {
		return "#ffffff"
	}
	if (299*r+587*g+114*b)/1000 > 150 {
		return "#111111"
	}
	return "#ffffff"
}
