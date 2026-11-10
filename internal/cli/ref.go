package cli

import (
	"context"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/huh"

	"go-with-the-flow/internal/model"
	"go-with-the-flow/internal/presets"
)

var (
	titleTagRE = regexp.MustCompile(`(?is)<title[^>]*>(.*?)</title>`)
	metaRE     = regexp.MustCompile(`(?is)<meta\s+[^>]*(?:name|property)=["']([^"']+)["'][^>]*content=["']([^"']*)["'][^>]*>`)
	colorRE    = regexp.MustCompile(`(?i)#[0-9a-f]{6}\b`)
	tagRE      = regexp.MustCompile(`(?is)<(header|nav|main|section|footer)\b`)
	tagStripRE = regexp.MustCompile(`(?s)<[^>]*>`)
	spaceRE    = regexp.MustCompile(`\s+`)
)

// RunReference fetches public page markup and uses its available metadata, colors,
// and section hints to seed an editable Astro project.
func RunReference(rawURL string) (*model.Config, error) {
	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
		return nil, errors.New("provide a complete http:// or https:// website URL")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "GoWithTheFlow/1.0 (+website reference generator)")
	resp, err := (&http.Client{Timeout: 20 * time.Second}).Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch reference website: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, fmt.Errorf("reference website returned %s", resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	markup := string(body)

	name := ""
	if match := titleTagRE.FindStringSubmatch(markup); len(match) > 1 {
		name = cleanText(match[1])
	}
	if name == "" {
		name = strings.TrimPrefix(u.Hostname(), "www.")
	}
	// Use the first title fragment before common separators as the project name.
	for _, sep := range []string{" | ", " — ", " - ", " :: "} {
		if i := strings.Index(name, sep); i > 0 {
			name = strings.TrimSpace(name[:i])
			break
		}
	}
	c := &model.Config{Name: name, Slug: slug(name), PackageManager: "npm", Industry: "saas", IndustryName: "Website", SubIndustry: "Website reference", Style: "soft", Animation: "subtle", CTAs: []string{"Learn More", "Contact Us"}, CanonicalURL: u.String()}
	for _, m := range metaRE.FindAllStringSubmatch(markup, -1) {
		key, value := strings.ToLower(m[1]), cleanText(m[2])
		switch key {
		case "description", "og:description", "twitter:description":
			if c.Description == "" {
				c.Description = value
			}
		case "og:site_name":
			if value != "" {
				c.Name = value
				c.Slug = slug(value)
			}
		}
	}
	if c.Description == "" {
		c.Description = "A website inspired by the visual direction of " + c.Name + "."
	}

	p, err := presets.Load()
	if err != nil {
		return nil, err
	}
	c.Colors.Primary, c.Colors.Secondary, c.Colors.Tertiary = "#1e3a5f", "#c9a227", "#f4efe6"
	seen := map[string]bool{}
	for _, color := range colorRE.FindAllString(markup, -1) {
		color = strings.ToLower(color)
		if color == "#ffffff" || color == "#000000" || seen[color] {
			continue
		}
		seen[color] = true
		if c.Colors.Primary == "#1e3a5f" {
			c.Colors.Primary = color
		} else if c.Colors.Secondary == "#c9a227" {
			c.Colors.Secondary = color
		} else if c.Colors.Tertiary == "#f4efe6" {
			c.Colors.Tertiary = color
		}
	}
	// Preserve the default palette slots if the page exposes fewer colors.
	if len(seen) == 0 {
		c.Colors.Primary, c.Colors.Secondary, c.Colors.Tertiary = "#1e3a5f", "#c9a227", "#f4efe6"
	}
	c.Fonts.Heading = model.Font{ID: p.Fonts[0].ID, Name: p.Fonts[0].Name, Fallback: p.Fonts[0].Fallback}
	c.Fonts.Body = c.Fonts.Heading

	want := map[string]bool{"hero": true, "about": true, "features": true, "contact": true}
	// Section tags are useful structural cues, but HTML alone cannot identify their meaning.
	if len(tagRE.FindAllStringSubmatch(markup, -1)) >= 6 {
		want["testimonials"] = true
	}
	for _, section := range p.Sections {
		if !want[section.ID] {
			continue
		}
		layout := ""
		if len(section.Layouts) > 0 {
			layout = section.Layouts[0].ID
		}
		c.Sections = append(c.Sections, model.Section{ID: section.ID, Layout: layout})
	}
	return c, nil
}

func cleanText(s string) string {
	return strings.TrimSpace(spaceRE.ReplaceAllString(html.UnescapeString(tagStripRE.ReplaceAllString(s, " ")), " "))
}

func ConfirmReference(c *model.Config) (bool, error) {
	ok := true
	err := huh.NewConfirm().Title(fmt.Sprintf("Generate reference-based Astro project %q? (extracted content and colors are editable)", c.Name)).Value(&ok).Run()
	return ok, err
}
