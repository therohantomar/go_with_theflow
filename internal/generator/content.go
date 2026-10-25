package generator

import (
	"fmt"

	"go-with-the-flow/internal/model"
)

type kv = map[string]any

func siteFor(c *model.Config) kv {
	href := "#"
	for _, s := range c.Sections {
		if s.ID == "contact" {
			href = "#contact"
		}
	}
	var ctas []kv
	for i, l := range c.CTAs {
		v := "primary"
		if i > 0 {
			v = "secondary"
		}
		ctas = append(ctas, kv{"label": l, "href": href, "variant": v})
	}
	nav := []kv{}
	for _, s := range c.Sections {
		if s.ID != "hero" {
			nav = append(nav, kv{"id": s.ID, "label": s.ID})
		}
	}
	return kv{"name": c.Name, "industry": c.IndustryName, "style": c.Style, "ctas": ctas, "nav": nav}
}

// contentFor returns placeholder copy per section. Images are left blank ("") —
// drop files into src/assets/<section>/ and put the filename in the JSON.
func contentFor(c *model.Config) map[string]any {
	ind := c.IndustryName
	items := func(n int, f func(i int) kv) []kv {
		var r []kv
		for i := 1; i <= n; i++ {
			r = append(r, f(i))
		}
		return r
	}
	m := map[string]any{
		"marquee": kv{"items": []string{"Quality", "Design", "Craft", "Delivery"}},
		"products": kv{"title": "Our products", "subtitle": "The right format for every need.",
			"items": items(3, func(i int) kv { return kv{"tag": "Category", "title": fmt.Sprintf("Product %d", i), "text": "Short product description.", "image": ""} })},
		"process": kv{"title": "How we work",
			"items": items(4, func(i int) kv { return kv{"title": fmt.Sprintf("Step %d", i), "text": "What happens in this step."} })},
		"industries": kv{"title": "Industries we serve",
			"items": items(6, func(i int) kv { return kv{"icon": "●", "name": fmt.Sprintf("Industry %d", i)} })},
		"quote": kv{"title": "Get an instant estimate", "text": "Tell us what you need and we will confirm a quote.",
			"whatsapp": "919999999999", "styles": []string{"Option A", "Option B"}, "boards": []string{"Standard", "Premium"},
			"minQty": 100, "maxQty": 10000, "step": 100},
		"hero": kv{"eyebrow": c.SubIndustry, "headline": c.Tagline,
			"subheadline": fmt.Sprintf("%s is your partner in %s. Replace this copy in src/content/hero.json.", c.Name, ind), "image": "", "points": []string{"Quality you can see", "Fast turnaround", "Support that listens"}},
		"about": kv{"title": "About " + c.Name, "text": fmt.Sprintf("Tell your %s story here.", ind), "image": "",
			"stats": items(3, func(i int) kv { return kv{"value": fmt.Sprintf("%d00+", i), "label": fmt.Sprintf("Metric %d", i)} })},
		"features": kv{"title": "Why choose us", "subtitle": "What sets " + c.Name + " apart.", "image": "",
			"items": items(4, func(i int) kv { return kv{"title": fmt.Sprintf("Feature %d", i), "text": "Short description of this feature."} })},
		"gallery": kv{"title": "Gallery",
			"items": items(6, func(i int) kv { return kv{"caption": fmt.Sprintf("Image %d", i), "image": ""} })},
		"testimonials": kv{"title": "What people say",
			"items": items(3, func(i int) kv {
				return kv{"quote": "A short, glowing testimonial goes here.", "name": fmt.Sprintf("Customer %d", i), "role": "Role, Company", "image": ""}
			})},
		"faq": kv{"title": "Frequently asked questions",
			"items": items(4, func(i int) kv { return kv{"q": fmt.Sprintf("Question %d?", i), "a": "Answer goes here."} })},
		"contact": kv{"title": "Get in touch", "text": "We'd love to hear from you.", "email": "hello@example.com", "phone": "+00 000 000 000", "address": "Your address", "image": ""},
	}
	switch c.Industry {
	case "packaging":
		packagingOverrides(c, m)
	case "food-packaging":
		foodPackagingOverrides(c, m)
	case "fashion":
		fashionOverrides(c, m)
	}
	return m
}

// packagingOverrides swaps the generic copy for box / packaging copy.
func packagingOverrides(c *model.Config, m map[string]any) {
	m["hero"] = kv{"eyebrow": "Design · Printing · Packaging", "headline": c.Tagline,
		"subheadline": "Corrugated boxes, printed cartons and custom packaging, made to your size, board and brand.", "image": "", "points": []string{"Custom sizes", "In-house printing", "Bulk and sample orders"}}
	m["about"] = kv{"title": "A packaging partner you can build on", "image": "",
		"text": "Tell your story here: when you started, who runs the business, and why packaging is more than protection for your customers.",
		"stats": []kv{{"value": "20+", "label": "Years in business"}, {"value": "3", "label": "Manufacturing units"}, {"value": "50k+", "label": "Boxes a month"}}}
	m["features"] = kv{"title": "Why " + c.Name, "subtitle": "One stop for design, printing and packaging.", "image": "",
		"items": []kv{
			{"title": "Quality focused", "text": "Consistent quality and attention to detail at every stage."},
			{"title": "Precision printing", "text": "Sharp, brand-accurate offset printing."},
			{"title": "Integrated manufacturing", "text": "Corrugation, lamination, die-cutting and gluing under one roof."},
			{"title": "Customer first", "text": "Customised solutions and long-term partnerships."}}}
	m["products"] = kv{"title": "The right format. The right board. The right finish.", "subtitle": "Our product range",
		"items": []kv{
			{"tag": "Corrugated", "title": "Corrugated Shipping Boxes", "text": "Strong cartons built for stacking strength and safe transit.", "image": ""},
			{"tag": "Printed", "title": "Printed Corrugated Boxes", "text": "Brand-printed boxes that look as good as they protect.", "image": ""},
			{"tag": "Retail ready", "title": "Printed Cartons", "text": "Offset-printed cartons for shelves and product launches.", "image": ""},
			{"tag": "E-commerce", "title": "Mailer Boxes", "text": "Easy-open mailers for direct-to-consumer shipping.", "image": ""},
			{"tag": "Gifting", "title": "Gift & Rigid Boxes", "text": "Premium boxes for gifting and luxury unboxing.", "image": ""},
			{"tag": "Custom", "title": "Die-cut Trays", "text": "Custom die-cut trays and inserts cut to your product.", "image": ""}}}
	m["process"] = kv{"title": "From idea to doorstep", "items": []kv{
		{"title": "Design", "text": "Structure and artwork shaped around your product and brand."},
		{"title": "Print", "text": "Precision offset printing for high-impact presentation."},
		{"title": "Corrugate & finish", "text": "Corrugation, lamination, die-cutting, gluing and stitching."},
		{"title": "Deliver", "text": "Dependable packaging delivered with consistent quality."}}}
	m["industries"] = kv{"title": "Packaging for diverse industries", "items": []kv{
		{"icon": "🛒", "name": "E-commerce & D2C"}, {"icon": "🍎", "name": "FMCG & Food"}, {"icon": "🔌", "name": "Electricals & Electronics"},
		{"icon": "🚗", "name": "Automotive & Engineering"}, {"icon": "👕", "name": "Textiles & Apparel"}, {"icon": "🎁", "name": "Gifting & Retail"}}}
	m["quote"] = kv{"title": "Tell us what you need — we'll quote it",
		"text": "Pick a style, enter dimensions and quantity, and send your requirement straight to our team on WhatsApp.",
		"whatsapp": "919999999999", "styles": []string{"RSC", "Die-cut", "Tray", "Mailer", "Gift box"},
		"boards": []string{"3-ply", "5-ply", "7-ply", "Duplex"}, "minQty": 100, "maxQty": 10000, "step": 100}
	m["contact"] = kv{"title": "Ready to pack your next big idea?", "text": "Tell us your product, size and quantity and we will reply with the right packaging solution.",
		"email": "hello@example.com", "phone": "+00 000 000 000", "address": "Unit address", "image": ""}
}

func foodPackagingOverrides(c *model.Config, m map[string]any) {
	m["hero"] = kv{"eyebrow": "Food packaging & disposables", "headline": c.Tagline, "image": "",
		"subheadline": "Containers, cups, tableware and bags for restaurants, cloud kitchens and food brands, in food-safe materials.",
		"points": []string{"Food-safe materials", "Custom printing", "Bulk supply"}}
	m["marquee"] = kv{"items": []string{"Food safe", "Leak resistant", "Recyclable options", "Custom branded"}}
	m["products"] = kv{"title": "Packaging for every menu", "subtitle": "Product catalogue",
		"items": []kv{
			{"tag": "Containers", "title": "Food Containers", "text": "Round and rectangular containers with secure lids.", "image": ""},
			{"tag": "Hot & cold", "title": "Paper Cups", "text": "Single and double wall cups for beverages and soups.", "image": ""},
			{"tag": "Tableware", "title": "Plates & Bowls", "text": "Disposable plates and bowls in paper and moulded fibre.", "image": ""},
			{"tag": "Takeaway", "title": "Bags & Wraps", "text": "Paper bags, wraps and sleeves for takeaway.", "image": ""},
			{"tag": "Bakery", "title": "Clamshells & Trays", "text": "Clear and paper clamshells for bakery and desserts.", "image": ""},
			{"tag": "Cutlery", "title": "Cutlery & Straws", "text": "Spoons, forks and straws in wood and paper.", "image": ""}}}
	m["industries"] = kv{"title": "Who we supply", "items": []kv{
		{"icon": "🍔", "name": "Restaurants & QSR"}, {"icon": "☁️", "name": "Cloud kitchens"}, {"icon": "🥐", "name": "Bakeries & cafes"},
		{"icon": "🍽️", "name": "Caterers"}, {"icon": "🛍️", "name": "Retail & grocery"}, {"icon": "🚚", "name": "Delivery brands"}}}
	m["quote"] = kv{"title": "Request a bulk quote", "text": "Tell us what you serve and how much you need.",
		"whatsapp": "919999999999", "styles": []string{"Container", "Cup", "Plate", "Bag", "Cutlery"},
		"boards": []string{"Paper", "PET", "Bagasse", "PP"}, "minQty": 500, "maxQty": 50000, "step": 500}
}

func fashionOverrides(c *model.Config, m map[string]any) {
	m["hero"] = kv{"eyebrow": "New season", "headline": c.Tagline, "image": "",
		"subheadline": "The new collection is here. Designed to be worn hard and styled your way.",
		"points": []string{"Free shipping", "Easy returns", "New drops weekly"}}
	m["marquee"] = kv{"items": []string{"New season", "Free shipping", "Limited drops", "Made to move"}}
	m["products"] = kv{"title": "Shop the collections", "subtitle": "Featured",
		"items": []kv{
			{"tag": "Men", "title": "Men", "text": "Essentials and statement pieces.", "image": ""},
			{"tag": "Women", "title": "Women", "text": "Everyday staples and seasonal edits.", "image": ""},
			{"tag": "Kids", "title": "Kids", "text": "Comfort that keeps up with them.", "image": ""},
			{"tag": "Sport", "title": "Performance", "text": "Built for training and movement.", "image": ""},
			{"tag": "Accessories", "title": "Accessories", "text": "Bags, caps and finishing touches.", "image": ""},
			{"tag": "Sale", "title": "Last chance", "text": "Final reductions on past seasons.", "image": ""}}}
	m["gallery"] = kv{"title": "Lookbook", "items": []kv{
		{"caption": "Look 01", "image": ""}, {"caption": "Look 02", "image": ""}, {"caption": "Look 03", "image": ""},
		{"caption": "Look 04", "image": ""}, {"caption": "Look 05", "image": ""}, {"caption": "Look 06", "image": ""}}}
}
