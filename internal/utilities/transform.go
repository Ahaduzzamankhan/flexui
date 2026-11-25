package utilities

import (
	"github.com/Ahaduzzamankhan/flexui/internal/ast"
)

// transformStack composes every transform utility through custom properties so
// that translate, rotate, scale and skew can be combined on one element instead
// of overwriting each other.
const transformStack = "translate(var(--flexui-tx, 0), var(--flexui-ty, 0)) rotate(var(--flexui-rotate, 0)) scale(var(--flexui-sx, 1), var(--flexui-sy, 1)) skew(var(--flexui-skx, 0), var(--flexui-sky, 0))"

// transformVar sets one custom property and re-emits the composed transform.
func transformVar(prop, value string) []ast.Block {
	return []ast.Block{on(decl(prop, value), decl("transform", transformStack))}
}

// scaleTokens maps the scale scale onto unitless multipliers.
var scaleTokens = map[string]string{
	"0": "0", "50": "0.5", "75": "0.75", "90": "0.9", "95": "0.95",
	"100": "1", "105": "1.05", "110": "1.1", "125": "1.25",
	"150": "1.5", "200": "2",
}

func scaleValue(v string) (string, bool) {
	if s, ok := scaleTokens[v]; ok {
		return s, true
	}
	return arb(v)
}

func registerTransforms(r *Registry) {
	r.prefix("translate-x", 190, func(c *Context, value string, neg bool) []ast.Block {
		if p, ok := fraction(value); ok {
			return transformVar("--flexui-tx", applyNeg(neg, p))
		}
		if v, ok := c.sizeValue(value, "100vw"); ok {
			return transformVar("--flexui-tx", applyNeg(neg, v))
		}
		return nil
	})
	r.prefix("translate-y", 191, func(c *Context, value string, neg bool) []ast.Block {
		if p, ok := fraction(value); ok {
			return transformVar("--flexui-ty", applyNeg(neg, p))
		}
		if v, ok := c.sizeValue(value, "100vh"); ok {
			return transformVar("--flexui-ty", applyNeg(neg, v))
		}
		return nil
	})
	r.prefix("rotate", 192, func(c *Context, value string, neg bool) []ast.Block {
		if v, ok := numValue(value); ok {
			return transformVar("--flexui-rotate", applyNeg(neg, v)+"deg")
		}
		return nil
	})
	r.both("scale", 193, func(c *Context, value string, neg bool) []ast.Block {
		if value == "" {
			return transformVar("--flexui-sx", "1")
		}
		if v, ok := scaleValue(value); ok {
			_ = neg
			return transformVar("--flexui-sx", v)
		}
		return nil
	})
	r.prefix("scale-x", 194, func(c *Context, value string, _ bool) []ast.Block {
		if v, ok := scaleValue(value); ok {
			return transformVar("--flexui-sx", v)
		}
		return nil
	})
	r.prefix("scale-y", 195, func(c *Context, value string, _ bool) []ast.Block {
		if v, ok := scaleValue(value); ok {
			return transformVar("--flexui-sy", v)
		}
		return nil
	})
	r.prefix("skew-x", 196, func(c *Context, value string, neg bool) []ast.Block {
		if v, ok := numValue(value); ok {
			return transformVar("--flexui-skx", applyNeg(neg, v)+"deg")
		}
		return nil
	})
	r.prefix("skew-y", 197, func(c *Context, value string, neg bool) []ast.Block {
		if v, ok := numValue(value); ok {
			return transformVar("--flexui-sky", applyNeg(neg, v)+"deg")
		}
		return nil
	})
	r.prefix("origin", 198, keyword("transform-origin", map[string]string{
		"center": "center", "top": "top", "top-right": "top right",
		"right": "right", "bottom-right": "bottom right", "bottom": "bottom",
		"bottom-left": "bottom left", "left": "left", "top-left": "top left",
	}))
	r.exactOnly("transform-none", 199, single("transform", "none"))
	r.prefix("perspective", 189, func(c *Context, value string, _ bool) []ast.Block {
		if v, ok := arb(value); ok {
			return []ast.Block{on(decl("perspective", v))}
		}
		if v, ok := c.spacing(value); ok {
			return []ast.Block{on(decl("perspective", v))}
		}
		return nil
	})
}
