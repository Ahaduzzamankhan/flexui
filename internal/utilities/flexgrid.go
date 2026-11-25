package utilities

import (
	"strconv"

	"github.com/Ahaduzzamankhan/flexui/internal/ast"
)

// flexFn resolves the "flex" family: display, direction, wrapping and the
// flex shorthand itself.
func flexFn(c *Context, value string, neg bool) []ast.Block {
	if value == "" {
		return []ast.Block{on(decl("display", "flex"))}
	}
	switch value {
	case "row", "row-reverse", "column", "column-reverse":
		return []ast.Block{on(decl("flex-direction", value))}
	case "col":
		return []ast.Block{on(decl("flex-direction", "column"))}
	case "col-reverse":
		return []ast.Block{on(decl("flex-direction", "column-reverse"))}
	case "wrap", "wrap-reverse", "nowrap":
		return []ast.Block{on(decl("flex-wrap", value))}
	case "1":
		return []ast.Block{on(decl("flex", "1 1 0%"))}
	case "auto":
		return []ast.Block{on(decl("flex", "1 1 auto"))}
	case "initial":
		return []ast.Block{on(decl("flex", "0 1 auto"))}
	case "none":
		return []ast.Block{on(decl("flex", "none"))}
	}
	if v, ok := arb(value); ok {
		if isBareNumber(v) {
			return []ast.Block{on(decl("flex", v))}
		}
		return []ast.Block{on(decl("flex", v))}
	}
	if v, ok := c.spacing(value); ok {
		return []ast.Block{on(decl("flex-basis", applyNeg(neg, v)))}
	}
	return nil
}

func isBareNumber(v string) bool {
	if v == "" {
		return false
	}
	for i := 0; i < len(v); i++ {
		if !(v[i] >= '0' && v[i] <= '9' || v[i] == '.') {
			return false
		}
	}
	return true
}

// gridTemplate resolves grid-cols-*/grid-rows-* counts and arbitrary templates.
func gridTemplate(prop string) resolveFunc {
	return func(c *Context, value string, _ bool) []ast.Block {
		switch value {
		case "none":
			return []ast.Block{on(decl(prop, "none"))}
		case "subgrid":
			return []ast.Block{on(decl(prop, "subgrid"))}
		}
		if n, err := strconv.Atoi(value); err == nil && n > 0 {
			return []ast.Block{on(decl(prop, "repeat("+strconv.Itoa(n)+", minmax(0, 1fr))"))}
		}
		if v, ok := arb(value); ok {
			return []ast.Block{on(decl(prop, v))}
		}
		return nil
	}
}

// gridPlacement resolves col-span-*/row-span-* style utilities.
func gridPlacement(prop, spanProp string) resolveFunc {
	return func(_ *Context, value string, _ bool) []ast.Block {
		switch value {
		case "auto":
			return []ast.Block{on(decl(prop, "auto"))}
		case "full":
			return []ast.Block{on(decl(spanProp, "1 / -1"))}
		}
		if v, ok := arb(value); ok {
			return []ast.Block{on(decl(prop, v))}
		}
		if isBareNumber(value) {
			return []ast.Block{on(decl(spanProp, "span "+value+" / span "+value))}
		}
		return nil
	}
}

func registerFlexGrid(r *Registry) {
	r.both("flex", 12, flexFn)
	r.exactOnly("inline-flex", 12, single("display", "inline-flex"))
	r.exactOnly("inline-grid", 12, single("display", "inline-grid"))
	r.exactOnly("grid", 12, single("display", "grid"))

	r.prefix("basis", 45, func(c *Context, value string, neg bool) []ast.Block {
		v, ok := c.sizeValue(value, "auto")
		if !ok {
			return nil
		}
		return []ast.Block{on(decl("flex-basis", applyNeg(neg, v)))}
	})

	r.both("grow", 46, func(_ *Context, value string, _ bool) []ast.Block {
		if value == "" {
			return []ast.Block{on(decl("flex-grow", "1"))}
		}
		v, ok := numValue(value)
		if !ok {
			return nil
		}
		return []ast.Block{on(decl("flex-grow", v))}
	})
	r.both("shrink", 47, func(_ *Context, value string, _ bool) []ast.Block {
		if value == "" {
			return []ast.Block{on(decl("flex-shrink", "1"))}
		}
		v, ok := numValue(value)
		if !ok {
			return nil
		}
		return []ast.Block{on(decl("flex-shrink", v))}
	})

	r.prefix("grid-cols", 50, gridTemplate("grid-template-columns"))
	r.prefix("grid-rows", 51, gridTemplate("grid-template-rows"))
	r.prefix("grid-flow", 52, keyword("grid-auto-flow", map[string]string{
		"row": "row", "col": "column", "dense": "dense",
		"row-dense": "row dense", "col-dense": "column dense",
	}))
	r.prefix("auto-cols", 53, keyword("grid-auto-columns", map[string]string{
		"auto": "auto", "min": "min-content", "max": "max-content", "fr": "minmax(0, 1fr)",
	}))
	r.prefix("auto-rows", 54, keyword("grid-auto-rows", map[string]string{
		"auto": "auto", "min": "min-content", "max": "max-content", "fr": "minmax(0, 1fr)",
	}))
	r.prefix("col-span", 55, gridPlacement("grid-column", "grid-column"))
	r.prefix("col-start", 56, gridPlacement("grid-column-start", "grid-column-start"))
	r.prefix("col-end", 56, gridPlacement("grid-column-end", "grid-column-end"))
	r.prefix("row-span", 57, gridPlacement("grid-row", "grid-row"))
	r.prefix("row-start", 58, gridPlacement("grid-row-start", "grid-row-start"))
	r.prefix("row-end", 58, gridPlacement("grid-row-end", "grid-row-end"))

	r.prefix("items", 40, keyword("align-items", map[string]string{
		"start": "flex-start", "end": "flex-end", "center": "center",
		"baseline": "baseline", "stretch": "stretch",
	}))
	r.prefix("justify", 41, keyword("justify-content", map[string]string{
		"start": "flex-start", "end": "flex-end", "center": "center",
		"between": "space-between", "around": "space-around", "evenly": "space-evenly",
		"stretch": "stretch", "normal": "normal",
	}))
	r.exactOnly("items-baseline", 40, single("align-items", "baseline"))
	r.prefix("content", 42, keyword("align-content", map[string]string{
		"start": "flex-start", "end": "flex-end", "center": "center",
		"between": "space-between", "around": "space-around", "evenly": "space-evenly",
		"baseline": "baseline", "stretch": "stretch", "normal": "normal",
	}))
	r.prefix("self", 43, keyword("align-self", map[string]string{
		"auto": "auto", "start": "flex-start", "end": "flex-end",
		"center": "center", "baseline": "baseline", "stretch": "stretch",
	}))
	r.prefix("justify-self", 43, keyword("justify-self", map[string]string{
		"auto": "auto", "start": "start", "end": "end",
		"center": "center", "stretch": "stretch",
	}))
	r.prefix("justify-items", 44, keyword("justify-items", map[string]string{
		"start": "start", "end": "end", "center": "center", "stretch": "stretch",
	}))
	r.prefix("place-items", 44, keyword("place-items", map[string]string{
		"start": "start", "end": "end", "center": "center", "stretch": "stretch",
	}))
	r.prefix("place-content", 44, keyword("place-content", map[string]string{
		"start": "start", "end": "end", "center": "center",
		"between": "space-between", "around": "space-around", "evenly": "space-evenly",
		"stretch": "stretch",
	}))
	r.prefix("order", 48, func(_ *Context, value string, neg bool) []ast.Block {
		if value == "first" {
			return []ast.Block{on(decl("order", "-9999"))}
		}
		if value == "last" {
			return []ast.Block{on(decl("order", "9999"))}
		}
		if value == "none" {
			return []ast.Block{on(decl("order", "0"))}
		}
		v, ok := numValue(value)
		if !ok {
			return nil
		}
		return []ast.Block{on(decl("order", applyNeg(neg, v)))}
	})
}
