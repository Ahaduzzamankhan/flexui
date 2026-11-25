package utilities

import "github.com/Ahaduzzamankhan/flexui/internal/ast"

// insetValue resolves inset tokens.
func (c *Context) insetValue(v string) (string, bool) {
	switch v {
	case "auto":
		return "auto", true
	case "full":
		return "100%", true
	}
	if p, ok := fraction(v); ok {
		return p, true
	}
	return c.spacing(v)
}

func insetProp(prop string) resolveFunc {
	return func(c *Context, value string, neg bool) []ast.Block {
		v, ok := c.insetValue(value)
		if !ok {
			return nil
		}
		return []ast.Block{on(decl(prop, applyNeg(neg, v)))}
	}
}

func registerPosition(r *Registry) {
	r.prefix("inset", 30, insetProp("inset"))
	r.prefix("inset-x", 31, insetProp("inset-inline"))
	r.prefix("inset-y", 32, insetProp("inset-block"))
	r.prefix("top", 33, insetProp("top"))
	r.prefix("right", 34, insetProp("right"))
	r.prefix("bottom", 35, insetProp("bottom"))
	r.prefix("left", 36, insetProp("left"))
	r.prefix("start", 37, insetProp("inset-inline-start"))
	r.prefix("end", 38, insetProp("inset-inline-end"))
	r.prefix("z", 39, func(c *Context, value string, neg bool) []ast.Block {
		if value == "auto" {
			return []ast.Block{on(decl("z-index", "auto"))}
		}
		if v, ok := c.Theme.ZIndex[value]; ok {
			return []ast.Block{on(decl("z-index", v))}
		}
		if v, ok := numValue(value); ok {
			return []ast.Block{on(decl("z-index", applyNeg(neg, v)))}
		}
		return nil
	})
}
