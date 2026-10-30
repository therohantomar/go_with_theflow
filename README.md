# Go With The Flow

A Go CLI that asks a few questions and generates a complete **Astro** landing page, styled with **Tailwind CSS**, animated with **anime.js** and **motion**, and using **Fontsource** fonts. All page copy lives in section-wise JSON files, and photo slots are left blank.

> Status: `flow new` works. `flow list`, `flow edit`, `flow generate` and `flow preview` are not built yet.

## What it asks

1. Package manager: npm, pnpm, yarn or bun
2. Project name
3. Industry and sub-industry (Packaging & Boxes, Food Packaging, Fashion & Apparel, Creative Studio, Real Estate, Pharmaceutical, SaaS, Fintech, Healthcare, Agency, E-commerce, Education)
4. Template style, with the best fit for the industry first (see below)
5. Sections: hero, marquee strip, trusted-by logos, about, stats counters, features, products, process, video, gallery, industries served, pricing, testimonials, team, blog, FAQ, instant quote (WhatsApp), CTA banner, newsletter, contact
6. A layout for each chosen section (carousel, grid, header pointers with image, and more)
7. Colors: primary, secondary, tertiary
8. Fonts: heading font and paragraph font (Inter, Plus Jakarta Sans, Urbanist, Outfit, Sora, Poppins, DM Sans, Montserrat, Manrope, Space Grotesk, Figtree, Lexend, Bricolage Grotesque, Red Hat Display, Rubik, Nunito, Open Sans, Oswald, Jost, Archivo, Syne, Playfair Display, Lora, Fraunces, Cormorant Garamond)
9. Number of CTAs and their labels
10. Animation level: none, subtle (motion) or expressive (anime.js + motion). The animation libraries are only installed if you choose them

## Images

Photos go in `src/assets/<section>/` (never `public/`). Astro optimises them with **sharp** at build time: AVIF and WebP, several sizes, lazy loading (hero images load eagerly). Add a file, then set its name in the section JSON, for example `"image": "hero-1.jpg"`.

## Animation

The wizard asks how much animation you want. Choosing none leaves out the libraries and the animation script. Reveals respect `prefers-reduced-motion`.

## Template styles

Each style changes type, shape, spacing, header, hero layout, default colors and fonts, so industries no longer look alike.

| Style | Look | Inspired by |
|---|---|---|
| Bold Sport | huge uppercase type, square corners, marquee strip, high contrast | Nike, Adidas, Myntra |
| Editorial Minimal | airy, hairline details, centered header, tall lookbook hero | Zara, Uniqlo |
| Industrial Catalog | solid color banner hero with key points, structured product catalogue | Huhtamaki, Dart, Pactiv, Genpak, Sabert, Novolex |
| Creative Dark | dark theme, oversized type, glow accents, pill buttons | Awwwards, Godly, Land-book, Dribbble |
| Soft Classic | rounded cards, friendly and balanced | - |

These are original designs that borrow general patterns from the reference sites. Nothing is copied.

## What it generates

```
your-project/
├── src/
│   ├── components/
│   │   ├── layout/        Header, Footer
│   │   ├── sections/      one component per section
│   │   └── Button, Carousel, Img
│   ├── content/           one JSON file per section, plus site.json
│   ├── assets/<section>/  empty photo folders
│   ├── layouts/  lib/  pages/  styles/
├── flow.json              your saved wizard choices
├── astro.config.mjs
└── package.json
```

## Project layout (this repo)

```
cmd/flow/                  entry point
internal/cli/              interactive wizard
internal/model/            config types
internal/presets/          presets.json (industries, sections, layouts, fonts)
internal/generator/        generator, placeholder content, embedded Astro templates
```

To add an industry, section option or font, edit `internal/presets/presets.json`. New section components go in `internal/generator/templates/static/src/components/sections/`.

## Requirements

- Go 1.22+
- Node.js 18+ and your chosen package manager

See [USAGE.md](USAGE.md) for copy-paste commands.
