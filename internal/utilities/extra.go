package utilities

import (
	"strconv"

	"github.com/Ahaduzzamankhan/flexui/internal/ast"
)

// registerExtra adds the remaining Phase 2 utilities: multi-column layout,
// place-self/place-content alignment and fragmentation controls.
func registerExtra(r *Registry) {
	r.prefix("columns", 60, func(c *Context, value string, _ bool) []ast.Block {
		if value == "auto" {
			return []ast.Block{on(decl("columns", "auto"))}
		}
		if size, ok := columnWidths[value]; ok {
			return []ast.Block{on(decl("columns", size))}
		}
		if n, err := strconv.Atoi(value); err == nil && n > 0 {
			return []ast.Block{on(decl("columns", strconv.Itoa(n)))}
		}
		if v, ok := arb(value); ok {
			return []ast.Block{on(decl("columns", v))}
		}
		return nil
	})
	r.prefix("place-self", 61, alignmentFn("place-self", map[string]string{
		"auto": "auto", "start": "start", "end": "end",
		"center": "center", "stretch": "stretch",
	}))
	r.prefix("place-content", 62, alignmentFn("place-content", map[string]string{
		"center": "center", "start": "start", "end": "end",
		"between": "space-between", "around": "space-around",
		"evenly": "space-evenly", "stretch": "stretch", "baseline": "baseline",
	}))
	r.both("break-inside", 63, func(_ *Context, value string, _ bool) []ast.Block {
		switch value {
		case "":
			return []ast.Block{on(decl("break-inside", "avoid"))}
		case "auto", "avoid", "avoid-page", "avoid-column", "avoid-region":
			return []ast.Block{on(decl("break-inside", value))}
		}
		return nil
	})
	r.prefix("break-before", 64, fragmentationFn("break-before"))
	r.prefix("break-after", 65, fragmentationFn("break-after"))
}

// columnWidths maps the named column widths onto CSS multi-column widths.
var columnWidths = map[string]string{
	"3xs": "16rem", "2xs": "18rem", "xs": "20rem", "sm": "24rem",
	"md": "28rem", "lg": "32rem", "xl": "36rem", "2xl": "42rem",
	"3xl": "48rem", "4xl": "56rem", "5xl": "64rem", "6xl": "72rem", "7xl": "80rem",
}

// alignmentFn builds a two-axis alignment utility from a value table.
func alignmentFn(prop string, values map[string]string) resolveFunc {
	return func(_ *Context, value string, _ bool) []ast.Block {
		v, ok := values[value]
		if !ok {
			if arbValue, ok := arb(value); ok {
				v, ok = arbValue, true
			}
		}
		if !ok || v == "" {
			return nil
		}
		if prop == "place-self" || prop == "place-content" {
			return []ast.Block{on(decl(prop, v))}
		}
		return []ast.Block{on(decl(prop, v))}
	}
}

// fragmentationFn builds break-before/break-after utilities.
func fragmentationFn(prop string) resolveFunc {
	return func(_ *Context, value string, _ bool) []ast.Block {
		switch value {
		case "auto", "avoid", "all", "avoid-page", "page", "left", "right",
			"column", "avoid-column", "region", "avoid-region":
			return []ast.Block{on(decl(prop, value))}
		}
		return nil
	}
}
