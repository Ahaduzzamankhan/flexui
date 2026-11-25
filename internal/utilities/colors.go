package utilities

import (
	"strings"

	"github.com/Ahaduzzamankhan/flexui/internal/ast"
	"github.com/Ahaduzzamankhan/flexui/internal/parser"
)

// classifyArbitrary interprets a bracketed value. It returns the normalized
// value plus a kind: "color", "length", "image" or "" when it cannot tell.
// Type hints such as "[length:12px]" win over the heuristics.
func classifyArbitrary(v string) (string, string) {
	inner, ok := parser.UnwrapArbitrary(v)
	if !ok {
		return "", ""
	}
	kind := ""
	raw := inner
	if i := strings.IndexByte(raw, ':'); i > 0 && isHintWord(raw[:i]) {
		kind = strings.ToLower(raw[:i])
		raw = raw[i+1:]
	}
	raw = parser.NormalizeArbitrary(raw)
	if raw == "" {
		return "", ""
	}
	if kind == "" {
		switch {
		case hasAnyPrefix(raw, "url(", "linear-gradient(", "radial-gradient(", "conic-gradient(", "image-set(", "repeating-linear-gradient(", "repeating-radial-gradient("):
			kind = "image"
		case hasAnyPrefix(raw, "#", "rgb", "hsl", "oklch", "oklab", "lab(", "lch(", "color(", "color-mix(", "light-dark("):
			kind = "color"
		case isLengthLike(raw):
			kind = "length"
		}
	}
	return raw, kind
}

func hasAnyPrefix(s string, prefixes ...string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

func isHintWord(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '-') {
			return false
		}
	}
	return true
}

// isLengthLike reports whether v is a plain length, percentage or a sizing
// function such as clamp() or calc().
func isLengthLike(v string) bool {
	if hasAnyPrefix(v, "clamp(", "calc(", "min(", "max(", "var(", "env(") {
		return true
	}
	i := 0
	if i < len(v) && (v[i] == '-' || v[i] == '+') {
		i++
	}
	digits := 0
	for i < len(v) && v[i] >= '0' && v[i] <= '9' {
		i++
		digits++
	}
	if i < len(v) && v[i] == '.' {
		i++
		for i < len(v) && v[i] >= '0' && v[i] <= '9' {
			i++
			digits++
		}
	}
	if digits == 0 {
		return false
	}
	// A unit must follow, otherwise the value is a bare number such as 0.
	return i < len(v)
}

// registerColors implements background, gradient, border, text and SVG colour
// utilities.
func registerColors(r *Registry) {
	r.prefix("bg", 100, bgFn)
	r.prefix("from", 101, gradientStop("--flexui-gradient-from"))
	r.prefix("to", 102, gradientStop("--flexui-gradient-to"))
	r.prefix("via", 103, gradientStop("--flexui-gradient-via"))
	r.prefix("accent", 104, func(c *Context, value string, _ bool) []ast.Block {
		if v, ok := c.colorValue(value); ok {
			return []ast.Block{on(decl("accent-color", v))}
		}
		return nil
	})
	r.prefix("caret", 105, func(c *Context, value string, _ bool) []ast.Block {
		if v, ok := c.colorValue(value); ok {
			return []ast.Block{on(decl("caret-color", v))}
		}
		return nil
	})
	r.prefix("fill", 106, func(c *Context, value string, _ bool) []ast.Block {
		if value == "none" {
			return []ast.Block{on(decl("fill", "none"))}
		}
		if v, ok := c.colorValue(value); ok {
			return []ast.Block{on(decl("fill", v))}
		}
		return nil
	})
	r.prefix("stroke", 107, func(c *Context, value string, _ bool) []ast.Block {
		if value == "none" {
			return []ast.Block{on(decl("stroke", "none"))}
		}
		if v, ok := numValue(value); ok && !strings.HasPrefix(v, "#") {
			if w, ok2 := c.Theme.Spacing[value]; ok2 {
				return []ast.Block{on(decl("stroke-width", w))}
			}
			_ = v
		}
		if v, ok := c.colorValue(value); ok {
			return []ast.Block{on(decl("stroke", v))}
		}
		return nil
	})
	r.both("divide-x", 92, borderWidthFor("border-inline-start-width", "& > * + *"))
	r.both("divide-y", 93, borderWidthFor("border-block-start-width", "& > * + *"))
	r.prefix("divide", 94, func(c *Context, value string, _ bool) []ast.Block {
		switch value {
		case "solid", "dashed", "dotted", "double", "none":
			return []ast.Block{under("& > * + *", decl("border-style", value))}
		}
		if v, ok := c.colorValue(value); ok {
			return []ast.Block{under("& > * + *", decl("border-color", v))}
		}
		return nil
	})
}

// bgFn resolves background utilities: colours, gradients, images, positions,
// sizes, repeats, attachments and clipping.
func bgFn(c *Context, value string, _ bool) []ast.Block {
	if value == "" {
		return nil
	}
	switch value {
	case "fixed", "local", "scroll":
		return []ast.Block{on(decl("background-attachment", value))}
	case "cover", "contain", "auto":
		return []ast.Block{on(decl("background-size", value))}
	case "center", "top", "bottom", "left", "right", "left-top", "left-bottom", "right-top", "right-bottom", "top-left", "top-right", "bottom-left", "bottom-right":
		return []ast.Block{on(decl("background-position", strings.ReplaceAll(value, "-", " ")))}
	case "repeat", "no-repeat", "repeat-x", "repeat-y", "repeat-round", "repeat-space":
		return []ast.Block{on(decl("background-repeat", value))}
	case "none":
		return []ast.Block{on(decl("background-image", "none"))}
	case "clip-border", "clip-padding", "clip-content", "clip-text":
		return []ast.Block{on(decl("background-clip", strings.TrimPrefix(value, "clip-")))}
	case "origin-border", "origin-padding", "origin-content":
		return []ast.Block{on(decl("background-origin", strings.TrimPrefix(value, "origin-")))}
	}
	if dir, ok := gradientDirections[value]; ok {
		return []ast.Block{on(decl("background-image",
			"linear-gradient("+dir+", var(--flexui-gradient-from, transparent), var(--flexui-gradient-via, var(--flexui-gradient-from, transparent)), var(--flexui-gradient-to, transparent))"))}
	}
	if strings.HasPrefix(value, "[") {
		raw, kind := classifyArbitrary(value)
		switch kind {
		case "image":
			return []ast.Block{on(decl("background-image", raw))}
		case "length":
			return []ast.Block{on(decl("background-size", raw))}
		case "position":
			return []ast.Block{on(decl("background-position", raw))}
		default:
			return []ast.Block{on(decl("background-color", raw))}
		}
	}
	if v, ok := c.colorValue(value); ok {
		return []ast.Block{on(decl("background-color", v))}
	}
	return nil
}

var gradientDirections = map[string]string{
	"gradient-to-t":  "to top",
	"gradient-to-tr": "to top right",
	"gradient-to-r":  "to right",
	"gradient-to-br": "to bottom right",
	"gradient-to-b":  "to bottom",
	"gradient-to-bl": "to bottom left",
	"gradient-to-l":  "to left",
	"gradient-to-tl": "to top left",
}

func gradientStop(prop string) resolveFunc {
	return func(c *Context, value string, _ bool) []ast.Block {
		if v, ok := c.colorValue(value); ok {
			return []ast.Block{on(decl(prop, v))}
		}
		return nil
	}
}

// borderWidthFor resolves divide-x/divide-y widths.
func borderWidthFor(prop, target string) resolveFunc {
	return func(_ *Context, value string, _ bool) []ast.Block {
		v, ok := borderWidthValue(value)
		if !ok {
			return nil
		}
		return []ast.Block{under(target, decl(prop, v))}
	}
}

// borderWidthValue resolves border width tokens.
func borderWidthValue(value string) (string, bool) {
	switch value {
	case "":
		return "1px", true
	case "0":
		return "0px", true
	case "2", "4", "8":
		return value + "px", true
	}
	if v, ok := arb(value); ok {
		return v, true
	}
	return "", false
}
