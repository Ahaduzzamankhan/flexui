package utilities

import (
	"github.com/Ahaduzzamankhan/flexui/internal/ast"
)

// shadowFn resolves box-shadow utilities. The value is also stored in a custom
// property so rings and shadows can be combined by the cascade.
func shadowFn(c *Context, value string, _ bool) []ast.Block {
	var v string
	if value == "" {
		v = c.Theme.Shadow[""]
	} else if s, ok := c.Theme.Shadow[value]; ok {
		v = s
	} else if a, ok := arb(value); ok {
		v = a
	} else {
		return nil
	}
	return []ast.Block{on(
		decl("--flexui-shadow", v),
		decl("box-shadow", "var(--flexui-ring-width, 0) var(--flexui-ring-width, 0) 0 var(--flexui-ring-spread, 0) var(--flexui-ring-color, transparent), var(--flexui-shadow)"),
	)}
}

// filterFn resolves a single filter function such as blur(4px).
func filterFn(prop, fn, unit string, defaults map[string]string) resolveFunc {
	return func(c *Context, value string, neg bool) []ast.Block {
		amount := ""
		if value == "" {
			amount = defaults[""]
		} else if v, ok := defaults[value]; ok {
			amount = v
		} else if v, ok := arb(value); ok {
			amount = v
		} else {
			return nil
		}
		_ = neg
		return []ast.Block{on(decl(prop, fn+"("+amount+unit+")"))}
	}
}

func registerEffects(r *Registry) {
	r.both("shadow", 170, shadowFn)
	r.prefix("shadow", 170, shadowFn)
	r.prefix("opacity", 171, func(c *Context, value string, _ bool) []ast.Block {
		if v, ok := c.Theme.Opacity[value]; ok {
			return []ast.Block{on(decl("opacity", v))}
		}
		if v, ok := numValue(value); ok {
			return []ast.Block{on(decl("opacity", v))}
		}
		return nil
	})

	blurs := map[string]string{
		"": "8px", "none": "0", "xs": "4px", "sm": "4px",
		"md": "12px", "lg": "16px", "xl": "24px", "2xl": "40px", "3xl": "64px",
	}
	r.both("blur", 172, filterFn("filter", "blur", "", blurs))
	r.both("backdrop-blur", 173, filterFn("backdrop-filter", "blur", "", blurs))
	r.prefix("brightness", 174, filterFn("filter", "brightness", "%", map[string]string{"": "100", "50": "50", "75": "75", "90": "90", "95": "95", "105": "105", "110": "110", "125": "125", "150": "150", "200": "200"}))
	r.prefix("contrast", 175, filterFn("filter", "contrast", "%", map[string]string{"": "100", "0": "0", "50": "50", "75": "75", "100": "100", "125": "125", "150": "150", "200": "200"}))
	r.prefix("grayscale", 176, filterFn("filter", "grayscale", "%", map[string]string{"": "100", "0": "0"}))
	r.prefix("invert", 177, filterFn("filter", "invert", "%", map[string]string{"": "100", "0": "0"}))
	r.prefix("saturate", 178, filterFn("filter", "saturate", "%", map[string]string{"": "150", "0": "0", "50": "50", "100": "100", "150": "150", "200": "200"}))
	r.prefix("sepia", 179, filterFn("filter", "sepia", "%", map[string]string{"": "100", "0": "0"}))
	r.prefix("hue-rotate", 180, func(c *Context, value string, neg bool) []ast.Block {
		if v, ok := numValue(value); ok {
			return []ast.Block{on(decl("filter", "hue-rotate("+applyNeg(neg, v)+"deg)"))}
		}
		return nil
	})
	r.prefix("drop-shadow", 181, func(c *Context, value string, _ bool) []ast.Block {
		if v, ok := c.Theme.Shadow[value]; ok {
			return []ast.Block{on(decl("filter", "drop-shadow("+v+")"))}
		}
		if v, ok := arb(value); ok {
			return []ast.Block{on(decl("filter", "drop-shadow("+v+")"))}
		}
		return nil
	})
	r.exactOnly("filter-none", 182, single("filter", "none"))
	r.exactOnly("backdrop-filter-none", 183, single("backdrop-filter", "none"))

	r.prefix("mix-blend", 184, keyword("mix-blend-mode", map[string]string{
		"normal": "normal", "multiply": "multiply", "screen": "screen",
		"overlay": "overlay", "darken": "darken", "lighten": "lighten",
		"color-dodge": "color-dodge", "color-burn": "color-burn",
		"hard-light": "hard-light", "soft-light": "soft-light",
		"difference": "difference", "exclusion": "exclusion",
		"hue": "hue", "saturation": "saturation", "color": "color",
		"luminosity": "luminosity", "plus-lighter": "plus-lighter",
		"plus-darker": "plus-darker",
	}))
	r.prefix("bg-blend", 185, keyword("background-blend-mode", map[string]string{
		"normal": "normal", "multiply": "multiply", "screen": "screen",
		"overlay": "overlay", "darken": "darken", "lighten": "lighten",
		"color-dodge": "color-dodge", "color-burn": "color-burn",
		"hard-light": "hard-light", "soft-light": "soft-light",
		"difference": "difference", "exclusion": "exclusion",
		"hue": "hue", "saturation": "saturation", "color": "color",
		"luminosity": "luminosity",
	}))
}
