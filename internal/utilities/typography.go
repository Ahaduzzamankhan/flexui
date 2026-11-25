package utilities

import (
	"strings"

	"github.com/Ahaduzzamankhan/flexui/internal/ast"
)

// textFn resolves the "text" family: alignment, wrapping, overflow, colour and
// the typography scale.
func textFn(c *Context, value string, _ bool) []ast.Block {
	if value == "" {
		return nil
	}
	switch value {
	case "left", "center", "right", "justify", "start", "end":
		return []ast.Block{on(decl("text-align", value))}
	case "wrap", "nowrap", "balance", "pretty", "stable":
		return []ast.Block{on(decl("text-wrap", value))}
	case "ellipsis", "clip":
		return []ast.Block{on(decl("text-overflow", value), decl("white-space", "nowrap"), decl("overflow", "hidden"))}
	}
	if strings.HasPrefix(value, "[") {
		raw, kind := classifyArbitrary(value)
		switch kind {
		case "color":
			return []ast.Block{on(decl("color", raw))}
		case "length", "size":
			return []ast.Block{on(decl("font-size", raw))}
		default:
			if isLengthLike(raw) {
				return []ast.Block{on(decl("font-size", raw))}
			}
			return []ast.Block{on(decl("color", raw))}
		}
	}
	if fs, ok := c.Theme.FontSize[value]; ok {
		decls := []ast.Decl{decl("font-size", fs.Size)}
		if fs.LineHeight != "" {
			decls = append(decls, decl("line-height", fs.LineHeight))
		}
		if fs.LetterSpacing != "" {
			decls = append(decls, decl("letter-spacing", fs.LetterSpacing))
		}
		return []ast.Block{on(decls...)}
	}
	if v, ok := c.colorValue(value); ok {
		return []ast.Block{on(decl("color", v))}
	}
	return nil
}

// fontFn resolves font-family and font-weight utilities.
func fontFn(c *Context, value string, _ bool) []ast.Block {
	if value == "" {
		return nil
	}
	if family, ok := c.Theme.FontFamily[value]; ok {
		return []ast.Block{on(decl("font-family", family))}
	}
	if weight, ok := c.Theme.FontWeight[value]; ok {
		return []ast.Block{on(decl("font-weight", weight))}
	}
	if value == "stretch" {
		return nil
	}
	if v, ok := arb(value); ok {
		return []ast.Block{on(decl("font-family", v))}
	}
	return nil
}

// typeScale maps semantic typography classes onto theme scales.
var typeScale = map[string]struct {
	size   string
	weight string
	mono   bool
}{
	"display-xl": {size: "display-xl", weight: "bold"},
	"display":    {size: "display", weight: "bold"},
	"heading-xl": {size: "4xl", weight: "semibold"},
	"heading-lg": {size: "3xl", weight: "semibold"},
	"heading":    {size: "2xl", weight: "semibold"},
	"subheading": {size: "lg", weight: "medium"},
	"body-lg":    {size: "lg", weight: "normal"},
	"body":       {size: "base", weight: "normal"},
	"body-sm":    {size: "sm", weight: "normal"},
	"caption":    {size: "xs", weight: "normal"},
	"label":      {size: "sm", weight: "medium"},
	"code":       {size: "sm", weight: "normal", mono: true},
}

func registerTypography(r *Registry) {
	r.prefix("text", 110, textFn)
	r.prefix("font", 111, fontFn)

	r.prefix("leading", 112, func(c *Context, value string, _ bool) []ast.Block {
		if v, ok := c.Theme.LineHeight[value]; ok {
			return []ast.Block{on(decl("line-height", v))}
		}
		if v, ok := arb(value); ok {
			return []ast.Block{on(decl("line-height", v))}
		}
		return nil
	})
	r.prefix("tracking", 113, func(c *Context, value string, neg bool) []ast.Block {
		if v, ok := c.Theme.LetterSpacing[value]; ok {
			return []ast.Block{on(decl("letter-spacing", v))}
		}
		if v, ok := arb(value); ok {
			return []ast.Block{on(decl("letter-spacing", applyNeg(neg, v)))}
		}
		return nil
	})
	r.prefix("indent", 114, func(c *Context, value string, neg bool) []ast.Block {
		if v, ok := c.spacing(value); ok {
			return []ast.Block{on(decl("text-indent", applyNeg(neg, v)))}
		}
		return nil
	})
	r.prefix("align", 115, keyword("vertical-align", map[string]string{
		"baseline": "baseline", "top": "top", "middle": "middle",
		"bottom": "bottom", "text-top": "text-top", "text-bottom": "text-bottom",
		"sub": "sub", "super": "super",
	}))
	r.prefix("whitespace", 116, keyword("white-space", map[string]string{
		"normal": "normal", "nowrap": "nowrap", "pre": "pre",
		"pre-line": "pre-line", "pre-wrap": "pre-wrap", "break-spaces": "break-spaces",
	}))
	r.prefix("break", 117, func(_ *Context, value string, _ bool) []ast.Block {
		switch value {
		case "normal":
			return []ast.Block{on(decl("overflow-wrap", "normal"), decl("word-break", "normal"))}
		case "words":
			return []ast.Block{on(decl("overflow-wrap", "break-word"))}
		case "all":
			return []ast.Block{on(decl("word-break", "break-all"))}
		case "keep":
			return []ast.Block{on(decl("word-break", "keep-all"))}
		}
		return nil
	})
	r.prefix("hyphens", 118, keyword("hyphens", map[string]string{
		"none": "none", "manual": "manual", "auto": "auto",
	}))
	r.prefix("line-clamp", 119, func(_ *Context, value string, _ bool) []ast.Block {
		v, ok := numValue(value)
		if !ok {
			return nil
		}
		if value == "none" {
			return []ast.Block{on(decl("-webkit-line-clamp", "none"), decl("line-clamp", "none"))}
		}
		return []ast.Block{on(
			decl("overflow", "hidden"),
			decl("display", "-webkit-box"),
			decl("-webkit-box-orient", "vertical"),
			decl("-webkit-line-clamp", v),
			decl("line-clamp", v),
		)}
	})
	r.prefix("list", 120, keyword("list-style-type", map[string]string{
		"none": "none", "disc": "disc", "decimal": "decimal",
		"circle": "circle", "square": "square", "roman": "upper-roman",
		"alpha": "lower-alpha",
	}))
	r.prefix("decoration", 121, func(c *Context, value string, _ bool) []ast.Block {
		switch value {
		case "solid", "double", "dotted", "dashed", "wavy":
			return []ast.Block{on(decl("text-decoration-style", value))}
		}
		if v, ok := c.colorValue(value); ok {
			return []ast.Block{on(decl("text-decoration-color", v))}
		}
		return nil
	})
	r.prefix("underline-offset", 122, func(_ *Context, value string, _ bool) []ast.Block {
		v, ok := numValue(value)
		if !ok {
			return nil
		}
		return []ast.Block{on(decl("text-underline-offset", v+"px"))}
	})
	r.prefix("type", 123, func(c *Context, value string, _ bool) []ast.Block {
		spec, ok := typeScale[value]
		if !ok {
			return nil
		}
		fs, ok := c.Theme.FontSize[spec.size]
		if !ok {
			return nil
		}
		decls := []ast.Decl{decl("font-size", fs.Size)}
		if fs.LineHeight != "" {
			decls = append(decls, decl("line-height", fs.LineHeight))
		}
		if fs.LetterSpacing != "" {
			decls = append(decls, decl("letter-spacing", fs.LetterSpacing))
		}
		if w, ok := c.Theme.FontWeight[spec.weight]; ok {
			decls = append(decls, decl("font-weight", w))
		}
		if spec.mono {
			decls = append(decls, decl("font-family", c.Theme.FontFamily["mono"]))
		}
		return []ast.Block{on(decls...)}
	})

	for name, value := range map[string]string{
		"italic":       "italic",
		"not-italic":   "normal",
		"underline":    "underline",
		"overline":     "overline",
		"line-through": "line-through",
		"no-underline": "none",
		"uppercase":    "uppercase",
		"lowercase":    "lowercase",
		"capitalize":   "capitalize",
		"normal-case":  "none",
	} {
		prop := "font-style"
		switch name {
		case "underline", "overline", "line-through", "no-underline":
			prop = "text-decoration-line"
		case "uppercase", "lowercase", "capitalize", "normal-case":
			prop = "text-transform"
		}
		r.exactOnly(name, 124, single(prop, value))
	}
	r.exactOnly("truncate", 125, func(*Context, string, bool) []ast.Block {
		return []ast.Block{on(
			decl("overflow", "hidden"),
			decl("text-overflow", "ellipsis"),
			decl("white-space", "nowrap"),
		)}
	})
	r.exactOnly("antialiased", 126, func(*Context, string, bool) []ast.Block {
		return []ast.Block{on(decl("-webkit-font-smoothing", "antialiased"), decl("-moz-osx-font-smoothing", "grayscale"))}
	})
	r.exactOnly("subpixel-antialiased", 126, func(*Context, string, bool) []ast.Block {
		return []ast.Block{on(decl("-webkit-font-smoothing", "auto"), decl("-moz-osx-font-smoothing", "auto"))}
	})
	r.exactOnly("tabular-nums", 127, single("font-variant-numeric", "tabular-nums"))
	r.exactOnly("proportional-nums", 127, single("font-variant-numeric", "proportional-nums"))
	r.exactOnly("oldstyle-nums", 127, single("font-variant-numeric", "oldstyle-nums"))
	r.exactOnly("lining-nums", 127, single("font-variant-numeric", "lining-nums"))
	r.exactOnly("normal-nums", 127, single("font-variant-numeric", "normal"))
	r.exactOnly("font-feature-normal", 128, single("font-feature-settings", "normal"))
	r.prefix("feat", 128, func(_ *Context, value string, _ bool) []ast.Block {
		if v, ok := arb(value); ok {
			return []ast.Block{on(decl("font-feature-settings", v))}
		}
		return nil
	})
}
