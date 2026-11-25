package utilities

import (
	"strings"

	"github.com/Ahaduzzamankhan/flexui/internal/ast"
)

// sizeValue resolves width/height tokens.
func (c *Context) sizeValue(v, screen string) (string, bool) {
	switch v {
	case "auto":
		return "auto", true
	case "full":
		return "100%", true
	case "screen":
		return screen, true
	case "svw", "lvw", "dvw", "svh", "lvh", "dvh":
		return "100" + v, true
	case "min":
		return "min-content", true
	case "max":
		return "max-content", true
	case "fit":
		return "fit-content", true
	}
	if p, ok := fraction(v); ok {
		return p, true
	}
	return c.spacing(v)
}

// sizeProp builds a width/height resolver.
func sizeProp(prop, screen string) resolveFunc {
	return func(c *Context, value string, _ bool) []ast.Block {
		v, ok := c.sizeValue(value, screen)
		if !ok {
			return nil
		}
		return []ast.Block{on(decl(prop, v))}
	}
}

func registerSizing(r *Registry) {
	r.prefix("w", 70, sizeProp("width", "100vw"))
	r.prefix("h", 71, sizeProp("height", "100vh"))
	r.prefix("min-w", 72, sizeProp("min-width", "100vw"))
	r.prefix("min-h", 73, sizeProp("min-height", "100vh"))
	r.prefix("max-w", 74, func(c *Context, value string, _ bool) []ast.Block {
		if v, ok := c.Theme.MaxWidth[value]; ok {
			return []ast.Block{on(decl("max-width", v))}
		}
		if v, ok := c.sizeValue(value, "100vw"); ok {
			return []ast.Block{on(decl("max-width", v))}
		}
		return nil
	})
	r.prefix("max-h", 75, func(c *Context, value string, _ bool) []ast.Block {
		if v, ok := c.Theme.MaxWidth[value]; ok {
			return []ast.Block{on(decl("max-height", v))}
		}
		if v, ok := c.sizeValue(value, "100vh"); ok {
			return []ast.Block{on(decl("max-height", v))}
		}
		return nil
	})
	r.prefix("size", 76, func(c *Context, value string, _ bool) []ast.Block {
		v, ok := c.sizeValue(value, "100vh")
		if !ok {
			return nil
		}
		return []ast.Block{on(decl("width", v), decl("height", v))}
	})
	r.prefix("aspect", 77, func(c *Context, value string, _ bool) []ast.Block {
		if v, ok := c.Theme.AspectRatio[value]; ok {
			return []ast.Block{on(decl("aspect-ratio", v))}
		}
		if strings.Contains(value, "/") {
			if p, ok := fraction(value); ok && !strings.HasSuffix(p, "%") {
				return []ast.Block{on(decl("aspect-ratio", strings.Replace(value, "/", " / ", 1)))}
			}
			return []ast.Block{on(decl("aspect-ratio", strings.Replace(value, "/", " / ", 1)))}
		}
		if v, ok := arb(value); ok {
			return []ast.Block{on(decl("aspect-ratio", v))}
		}
		return nil
	})
}
