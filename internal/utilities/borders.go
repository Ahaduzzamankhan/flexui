package utilities

import (
	"strings"

	"github.com/Ahaduzzamankhan/flexui/internal/ast"
)

// borderAnyFn resolves the "border" family: width, style and colour.
func borderAnyFn(c *Context, value string, _ bool) []ast.Block {
	if value == "" {
		return []ast.Block{on(decl("border-width", "1px"))}
	}
	switch value {
	case "solid", "dashed", "dotted", "double", "hidden", "none":
		return []ast.Block{on(decl("border-style", value))}
	}
	if v, ok := borderWidthValue(value); ok {
		return []ast.Block{on(decl("border-width", v))}
	}
	if strings.HasPrefix(value, "[") {
		raw, kind := classifyArbitrary(value)
		if kind == "length" || kind == "" && isLengthLike(raw) {
			return []ast.Block{on(decl("border-width", raw))}
		}
		return []ast.Block{on(decl("border-color", raw))}
	}
	if v, ok := c.colorValue(value); ok {
		return []ast.Block{on(decl("border-color", v))}
	}
	return nil
}

// borderSideFn resolves border-t-*, border-x-* and friends.
func borderSideFn(prop string) resolveFunc {
	return func(c *Context, value string, _ bool) []ast.Block {
		if value == "" {
			return []ast.Block{on(decl(prop, "1px"))}
		}
		if v, ok := borderWidthValue(value); ok {
			return []ast.Block{on(decl(prop, v))}
		}
		if strings.HasPrefix(value, "[") {
			raw, kind := classifyArbitrary(value)
			if kind == "length" || kind == "" && isLengthLike(raw) {
				return []ast.Block{on(decl(prop, raw))}
			}
			return []ast.Block{on(decl(sideColorProp(prop), raw))}
		}
		if v, ok := c.colorValue(value); ok {
			return []ast.Block{on(decl(sideColorProp(prop), v))}
		}
		return nil
	}
}

// sideColorProp maps a side width property onto its colour counterpart.
func sideColorProp(widthProp string) string {
	return strings.Replace(widthProp, "-width", "-color", 1)
}

// radiusFn resolves the rounded family.
func radiusFn(props []string) resolveFunc {
	return func(c *Context, value string, _ bool) []ast.Block {
		var v string
		if value == "" {
			v = c.Theme.Radius[""]
		} else if r, ok := c.Theme.Radius[value]; ok {
			v = r
		} else if a, ok := arb(value); ok {
			v = a
		} else {
			return nil
		}
		decls := make([]ast.Decl, 0, len(props))
		for _, p := range props {
			decls = append(decls, decl(p, v))
		}
		return []ast.Block{on(decls...)}
	}
}

func registerBorders(r *Registry) {
	r.both("border", 90, borderAnyFn)
	r.both("border-x", 91, borderSideFn("border-inline-width"))
	r.both("border-y", 91, borderSideFn("border-block-width"))
	r.both("border-t", 93, borderSideFn("border-top-width"))
	r.both("border-r", 94, borderSideFn("border-right-width"))
	r.both("border-b", 95, borderSideFn("border-bottom-width"))
	r.both("border-l", 96, borderSideFn("border-left-width"))
	r.both("border-s", 97, borderSideFn("border-inline-start-width"))
	r.both("border-e", 98, borderSideFn("border-inline-end-width"))

	r.both("rounded", 140, radiusFn([]string{"border-radius"}))
	r.prefix("rounded-t", 141, radiusFn([]string{"border-top-left-radius", "border-top-right-radius"}))
	r.prefix("rounded-r", 142, radiusFn([]string{"border-top-right-radius", "border-bottom-right-radius"}))
	r.prefix("rounded-b", 143, radiusFn([]string{"border-bottom-right-radius", "border-bottom-left-radius"}))
	r.prefix("rounded-l", 144, radiusFn([]string{"border-top-left-radius", "border-bottom-left-radius"}))
	r.prefix("rounded-tl", 145, radiusFn([]string{"border-top-left-radius"}))
	r.prefix("rounded-tr", 145, radiusFn([]string{"border-top-right-radius"}))
	r.prefix("rounded-br", 146, radiusFn([]string{"border-bottom-right-radius"}))
	r.prefix("rounded-bl", 146, radiusFn([]string{"border-bottom-left-radius"}))
	r.prefix("rounded-ss", 147, radiusFn([]string{"border-start-start-radius"}))
	r.prefix("rounded-se", 147, radiusFn([]string{"border-start-end-radius"}))
	r.prefix("rounded-es", 148, radiusFn([]string{"border-end-start-radius"}))
	r.prefix("rounded-ee", 148, radiusFn([]string{"border-end-end-radius"}))

	r.both("outline", 150, func(c *Context, value string, _ bool) []ast.Block {
		if value == "" {
			return []ast.Block{on(decl("outline-style", "solid"), decl("outline-width", "1px"))}
		}
		switch value {
		case "none", "hidden":
			return []ast.Block{on(decl("outline-style", "none"))}
		case "solid", "dashed", "dotted", "double":
			return []ast.Block{on(decl("outline-style", value))}
		}
		if v, ok := borderWidthValue(value); ok {
			return []ast.Block{on(decl("outline-width", v))}
		}
		if v, ok := c.colorValue(value); ok {
			return []ast.Block{on(decl("outline-color", v))}
		}
		return nil
	})
	r.prefix("outline-offset", 151, func(_ *Context, value string, neg bool) []ast.Block {
		v, ok := numValue(value)
		if !ok {
			return nil
		}
		return []ast.Block{on(decl("outline-offset", applyNeg(neg, v+"px")))}
	})

	// Rings are composed into box-shadow with the shadow family so that a ring
	// and a shadow can coexist on one element.
	r.both("ring", 152, func(c *Context, value string, _ bool) []ast.Block {
		width := "3px"
		if value != "" {
			if v, ok := borderWidthValue(value); ok {
				width = v
			} else if v, ok := c.colorValue(value); ok {
				return []ast.Block{on(decl("--flexui-ring-color", v))}
			} else {
				return nil
			}
		}
		return []ast.Block{on(
			decl("--flexui-ring-width", width),
			decl("box-shadow", "var(--flexui-ring-width, 0) var(--flexui-ring-width, 0) 0 var(--flexui-ring-spread, 0) var(--flexui-ring-color, rgb(59 130 246 / 0.5)), var(--flexui-shadow, 0 0 #0000)"),
		)}
	})
	r.prefix("ring-offset", 153, func(c *Context, value string, _ bool) []ast.Block {
		if v, ok := c.colorValue(value); ok {
			return []ast.Block{on(decl("--flexui-ring-offset-color", v))}
		}
		return nil
	})
}
