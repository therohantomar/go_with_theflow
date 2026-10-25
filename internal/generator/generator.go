package generator

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"go-with-the-flow/internal/model"
)

//go:embed all:templates/static
var staticFS embed.FS

const staticRoot = "templates/static"

// Generate writes a complete Astro project into <parent>/<slug> and returns its path.
func Generate(c *model.Config, parent string) (string, error) {
	out := filepath.Join(parent, c.Slug)
	if _, err := os.Stat(out); err == nil {
		return "", fmt.Errorf("%s already exists", out)
	}

	// 1. static files (components, layout, config) copied verbatim
	err := fs.WalkDir(staticFS, staticRoot, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := staticFS.ReadFile(p)
		if err != nil {
			return err
		}
		return write(filepath.Join(out, strings.TrimPrefix(p, staticRoot+"/")), data)
	})
	if err != nil {
		return "", err
	}

	// 2. dynamic files rendered from the config
	type sec struct{ ID, Component, Layout string }
	var secs []sec
	for _, s := range c.Sections {
		secs = append(secs, sec{s.ID, strings.ToUpper(s.ID[:1]) + s.ID[1:], s.Layout})
	}
	// the faq section file is Faq.astro
	data := struct {
		*model.Config
		Sections []sec
		FontIDs  []string
	}{c, secs, uniqueFonts(c)}
	// shadow embedded Config.Sections with the richer list
	files := map[string]string{
		"package.json":          packageJSON,
		"src/styles/global.css": globalCSS,
		"src/styles/fonts.css":  fontsCSS,
		"src/pages/index.astro": indexAstro,
	}
	for path, tpl := range files {
		s, err := render(tpl, data)
		if err != nil {
			return "", err
		}
		if err := write(filepath.Join(out, path), []byte(s)); err != nil {
			return "", err
		}
	}

	// 3. section-wise content JSON + site.json
	selected := map[string]bool{}
	for _, s := range c.Sections {
		selected[s.ID] = true
	}
	for id, v := range contentFor(c) {
		if !selected[id] {
			continue
		}
		if err := writeJSON(filepath.Join(out, "src/content", id+".json"), v); err != nil {
			return "", err
		}
	}
	if err := writeJSON(filepath.Join(out, "src/content/site.json"), siteFor(c)); err != nil {
		return "", err
	}

	// 4. empty section-wise asset folders for photos
	for _, s := range c.Sections {
		if err := write(filepath.Join(out, "src/assets", s.ID, ".gitkeep"), nil); err != nil {
			return "", err
		}
	}

	// 5. flow.json (project state) inside the project, and a copy in ~/.go-with-the-flow/projects
	if err := writeJSON(filepath.Join(out, "flow.json"), c); err != nil {
		return "", err
	}
	if home, err := os.UserHomeDir(); err == nil {
		_ = writeJSON(filepath.Join(home, ".go-with-the-flow", "projects", c.Slug+".json"), c)
	}
	return out, nil
}

func uniqueFonts(c *model.Config) []string {
	ids := []string{c.Fonts.Heading.ID}
	if c.Fonts.Body.ID != c.Fonts.Heading.ID {
		ids = append(ids, c.Fonts.Body.ID)
	}
	return ids
}

func render(tpl string, data any) (string, error) {
	t, err := template.New("t").Parse(tpl)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	if err := t.Execute(&b, data); err != nil {
		return "", err
	}
	return b.String(), nil
}

func write(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func writeJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return write(path, append(b, '\n'))
}

const packageJSON = `{
  "name": "{{.Slug}}",
  "type": "module",
  "version": "0.1.0",
  "scripts": {
    "dev": "astro dev",
    "build": "astro build",
    "preview": "astro preview"
  },
  "dependencies": {
    "@tailwindcss/vite": "^4.0.0",
    "animejs": "^4.0.0",
    "astro": "^5.0.0",
    "motion": "^12.0.0",
    "tailwindcss": "^4.0.0"{{range .FontIDs}},
    "@fontsource/{{.}}": "^5.0.0"{{end}}
  }
}
`

const fontsCSS = `{{range .FontIDs}}@import "@fontsource/{{.}}/400.css";
@import "@fontsource/{{.}}/500.css";
@import "@fontsource/{{.}}/600.css";
@import "@fontsource/{{.}}/700.css";
{{end}}`

const globalCSS = `@import "tailwindcss";
@import "./fonts.css";

@theme {
  --color-primary: {{.Colors.Primary}};
  --color-secondary: {{.Colors.Secondary}};
  --color-tertiary: {{.Colors.Tertiary}};
  --font-heading: "{{.Fonts.Heading.Name}}", {{.Fonts.Heading.Fallback}};
  --font-body: "{{.Fonts.Body.Name}}", {{.Fonts.Body.Fallback}};
}

@layer base {
  html { scroll-behavior: smooth; }
  body { @apply bg-white font-body text-neutral-800; }
  h1, h2, h3, h4 { @apply font-heading text-primary; }
}

/* revealed by motion / anime.js on load or scroll */
[data-reveal], [data-hero-item] { opacity: 0; }
`

const indexAstro = `---
import Layout from "../layouts/Layout.astro";
import Header from "../components/layout/Header.astro";
import Footer from "../components/layout/Footer.astro";
import site from "../content/site.json";
{{range .Sections}}import {{.Component}} from "../components/sections/{{.Component}}.astro";
import {{.ID}}Data from "../content/{{.ID}}.json";
{{end}}---

<Layout title={site.name}>
  <Header />
  <main>
{{range .Sections}}    <{{.Component}} data={ {{.ID}}Data } layout="{{.Layout}}" />
{{end}}  </main>
  <Footer />
</Layout>
`
