package cli

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"

	"go-with-the-flow/internal/model"
	"go-with-the-flow/internal/presets"
)

var (
	hexRe  = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)
	slugRe = regexp.MustCompile(`[^a-z0-9]+`)
)

func slug(s string) string {
	return strings.Trim(slugRe.ReplaceAllString(strings.ToLower(s), "-"), "-")
}

func hexCheck(s string) error {
	if !hexRe.MatchString(s) {
		return errors.New("use a 6-digit hex colour, e.g. #1e3a5f")
	}
	return nil
}

// RunWizard walks the user through every choice. Returns nil config if the user declines at review.
func RunWizard() (*model.Config, error) {
	p, err := presets.Load()
	if err != nil {
		return nil, err
	}
	c := &model.Config{PackageManager: "npm"}

	// 1. package manager + project name
	if err := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().Title("Which package manager?").
			Options(huh.NewOptions("npm", "pnpm", "yarn", "bun")...).Value(&c.PackageManager),
		huh.NewInput().Title("Project name").Value(&c.Name).Validate(func(s string) error {
			if slug(s) == "" {
				return errors.New("name is required")
			}
			return nil
		}),
	)).Run(); err != nil {
		return nil, err
	}
	c.Slug = slug(c.Name)

	// 2. industry
	var indOpts []huh.Option[string]
	for _, i := range p.Industries {
		indOpts = append(indOpts, huh.NewOption(i.Name, i.ID))
	}
	if err := huh.NewSelect[string]().Title("What are we building?").Options(indOpts...).Value(&c.Industry).Run(); err != nil {
		return nil, err
	}
	var ind presets.Industry
	for _, i := range p.Industries {
		if i.ID == c.Industry {
			ind = i
		}
	}
	c.IndustryName, c.Tagline = ind.Name, ind.Tagline

	var subOpts []huh.Option[string]
	for _, s := range ind.Sub {
		subOpts = append(subOpts, huh.NewOption(s, s))
	}
	if err := huh.NewSelect[string]().Title("Sub-industry").Options(subOpts...).Value(&c.SubIndustry).Run(); err != nil {
		return nil, err
	}

	// 3. sections
	var secOpts []huh.Option[string]
	for _, s := range p.Sections {
		secOpts = append(secOpts, huh.NewOption(s.Name, s.ID).Selected(s.ID == "hero"))
	}
	var picked []string
	if err := huh.NewMultiSelect[string]().Title("Which sections? (space to toggle)").Options(secOpts...).Value(&picked).
		Validate(func(s []string) error {
			if len(s) == 0 {
				return errors.New("pick at least one section")
			}
			return nil
		}).Run(); err != nil {
		return nil, err
	}

	// 4. layout per section (keep preset order)
	var chosen []presets.Section
	for _, s := range p.Sections {
		for _, id := range picked {
			if id == s.ID {
				chosen = append(chosen, s)
			}
		}
	}
	layouts := make([]string, len(chosen))
	var fields []huh.Field
	for i, s := range chosen {
		var opts []huh.Option[string]
		for _, l := range s.Layouts {
			opts = append(opts, huh.NewOption(l.Name, l.ID))
		}
		layouts[i] = s.Layouts[0].ID
		fields = append(fields, huh.NewSelect[string]().Title(s.Name+" layout").Options(opts...).Value(&layouts[i]))
	}
	if err := huh.NewForm(huh.NewGroup(fields...)).Run(); err != nil {
		return nil, err
	}
	for i, s := range chosen {
		c.Sections = append(c.Sections, model.Section{ID: s.ID, Layout: layouts[i]})
	}

	// 5. colours
	c.Colors.Primary, c.Colors.Secondary, c.Colors.Tertiary = ind.Palette[0], ind.Palette[1], ind.Palette[2]
	if err := huh.NewForm(huh.NewGroup(
		huh.NewInput().Title("Primary colour").Value(&c.Colors.Primary).Validate(hexCheck),
		huh.NewInput().Title("Secondary colour").Value(&c.Colors.Secondary).Validate(hexCheck),
		huh.NewInput().Title("Tertiary colour").Value(&c.Colors.Tertiary).Validate(hexCheck),
	)).Run(); err != nil {
		return nil, err
	}

	// 6. fonts
	var fontOpts []huh.Option[string]
	for _, f := range p.Fonts {
		fontOpts = append(fontOpts, huh.NewOption(f.Name, f.ID))
	}
	headID, bodyID := "playfair-display", "inter"
	if err := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().Title("Primary font (headings)").Options(fontOpts...).Value(&headID),
		huh.NewSelect[string]().Title("Secondary font (paragraphs)").Options(fontOpts...).Value(&bodyID),
	)).Run(); err != nil {
		return nil, err
	}
	for _, f := range p.Fonts {
		if f.ID == headID {
			c.Fonts.Heading = model.Font{ID: f.ID, Name: f.Name, Fallback: f.Fallback}
		}
		if f.ID == bodyID {
			c.Fonts.Body = model.Font{ID: f.ID, Name: f.Name, Fallback: f.Fallback}
		}
	}

	// 7. CTAs
	count := "2"
	if err := huh.NewInput().Title("How many CTAs? (1-4)").Value(&count).Validate(func(s string) error {
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 || n > 4 {
			return errors.New("enter a number from 1 to 4")
		}
		return nil
	}).Run(); err != nil {
		return nil, err
	}
	n, _ := strconv.Atoi(count)
	c.CTAs = make([]string, n)
	var ctaFields []huh.Field
	for i := range c.CTAs {
		c.CTAs[i] = []string{"Get Started", "Learn More", "Contact Us", "Book a Call"}[i]
		ctaFields = append(ctaFields, huh.NewInput().Title(fmt.Sprintf("CTA %d label", i+1)).Value(&c.CTAs[i]))
	}
	if err := huh.NewForm(huh.NewGroup(ctaFields...)).Run(); err != nil {
		return nil, err
	}

	// review
	ok := true
	if err := huh.NewConfirm().Title(fmt.Sprintf("Generate %s (%s / %s, %d sections, %d CTAs)?",
		c.Name, c.IndustryName, c.SubIndustry, len(c.Sections), len(c.CTAs))).Value(&ok).Run(); err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return c, nil
}
