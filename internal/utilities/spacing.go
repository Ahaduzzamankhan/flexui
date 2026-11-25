package utilities

import "github.com/Ahaduzzamankhan/flexui/internal/ast"

// spacingProps builds a resolver that sets one value on several properties.
func spacingProps(props ...string) resolveFunc {
	return func(c *Context, value string, neg bool) []ast.Block {
		v, ok := c.spacing(value)
		if !ok {
			return nil
		}
		v = applyNeg(neg, v)
		decls := make([]ast.Decl, 0, len(props))
		for _, p := range props {
			decls = append(decls, decl(p, v))
		}
		return []ast.Block{on(decls...)}
	}
}

// marginProps resolves margins, which additionally accept "auto".
func marginProps(props ...string) resolveFunc {
	return func(c *Context, value string, neg bool) []ast.Block {
		v := "auto"
		if value != "auto" {
			s, ok := c.spacing(value)
			if !ok {
				return nil
			}
			v = applyNeg(neg, s)
		}
		decls := make([]ast.Decl, 0, len(props))
		for _, p := range props {
			decls = append(decls, decl(p, v))
		}
		return []ast.Block{on(decls...)}
	}
}

// childSpacingProps applies a spacing value to following siblings, which is how
// space-x and space-y work.
func childSpacingProps(prop string) resolveFunc {
	return func(c *Context, value string, neg bool) []ast.Block {
		v, ok := c.spacing(value)
		if !ok {
			return nil
		}
		v = applyNeg(neg, v)
		return []ast.Block{under("& > * + *", decl(prop, v))}
	}
}

func registerSpacing(r *Registry) {
	r.prefix("p", 200, spacingProps("padding"))
	r.prefix("px", 201, spacingProps("padding-inline"))
	r.prefix("py", 202, spacingProps("padding-block"))
	r.prefix("pt", 203, spacingProps("padding-top"))
	r.prefix("pr", 204, spacingProps("padding-right"))
	r.prefix("pb", 205, spacingProps("padding-bottom"))
	r.prefix("pl", 206, spacingProps("padding-left"))
	r.prefix("ps", 207, spacingProps("padding-inline-start"))
	r.prefix("pe", 208, spacingProps("padding-inline-end"))

	r.prefix("m", 210, marginProps("margin"))
	r.prefix("mx", 211, marginProps("margin-inline"))
	r.prefix("my", 212, marginProps("margin-block"))
	r.prefix("mt", 213, marginProps("margin-top"))
	r.prefix("mr", 214, marginProps("margin-right"))
	r.prefix("mb", 215, marginProps("margin-bottom"))
	r.prefix("ml", 216, marginProps("margin-left"))
	r.prefix("ms", 217, marginProps("margin-inline-start"))
	r.prefix("me", 218, marginProps("margin-inline-end"))

	r.prefix("gap", 60, spacingProps("gap"))
	r.prefix("gap-x", 61, spacingProps("column-gap"))
	r.prefix("gap-y", 62, spacingProps("row-gap"))

	r.prefix("space-x", 63, childSpacingProps("margin-inline-start"))
	r.prefix("space-y", 64, childSpacingProps("margin-block-start"))
}
