package utilities

import "github.com/Ahaduzzamankhan/flexui/internal/ast"

// single builds a resolver that always emits one declaration.
func single(prop, value string) resolveFunc {
	return func(*Context, string, bool) []ast.Block {
		return []ast.Block{on(decl(prop, value))}
	}
}

// keyword builds a resolver mapping a value keyword onto a property, with an
// optional default for the bare form.
func keyword(prop string, values map[string]string) resolveFunc {
	return func(_ *Context, value string, _ bool) []ast.Block {
		if value == "" {
			return nil
		}
		if v, ok := values[value]; ok {
			return []ast.Block{on(decl(prop, v))}
		}
		if v, ok := arb(value); ok {
			return []ast.Block{on(decl(prop, v))}
		}
		return nil
	}
}

func registerLayout(r *Registry) {
	// Display and box generation.
	displays := map[string]string{
		"block": "block", "inline-block": "inline-block", "inline": "inline",
		"flow-root": "flow-root", "list-item": "list-item",
		"inline-table": "inline-table", "table": "table",
		"table-row": "table-row", "table-cell": "table-cell",
		"table-caption": "table-caption", "table-column": "table-column",
		"contents": "contents", "hidden": "none",
	}
	for name, value := range displays {
		r.exactOnly(name, 10, single("display", value))
	}

	// Positioning schemes.
	positions := map[string]string{
		"static": "static", "relative": "relative", "absolute": "absolute",
		"fixed": "fixed", "sticky": "sticky",
	}
	for name, value := range positions {
		r.exactOnly(name, 20, single("position", value))
	}

	// Visibility and isolation.
	r.exactOnly("visible", 15, single("visibility", "visible"))
	r.exactOnly("invisible", 15, single("visibility", "hidden"))
	r.exactOnly("collapse", 15, single("visibility", "collapse"))
	r.exactOnly("isolate", 16, single("isolation", "isolate"))
	r.exactOnly("isolation-auto", 16, single("isolation", "auto"))

	// Screen-reader-only helpers.
	r.exactOnly("sr-only", 160, func(*Context, string, bool) []ast.Block {
		return []ast.Block{on(
			decl("position", "absolute"),
			decl("width", "1px"),
			decl("height", "1px"),
			decl("padding", "0"),
			decl("margin", "-1px"),
			decl("overflow", "hidden"),
			decl("clip-path", "inset(50%)"),
			decl("white-space", "nowrap"),
			decl("border-width", "0"),
		)}
	})
	r.exactOnly("not-sr-only", 161, func(*Context, string, bool) []ast.Block {
		return []ast.Block{on(
			decl("position", "static"),
			decl("width", "auto"),
			decl("height", "auto"),
			decl("padding", "0"),
			decl("margin", "0"),
			decl("overflow", "visible"),
			decl("clip-path", "none"),
			decl("white-space", "normal"),
		)}
	})

	// Overflow.
	r.prefix("overflow", 165, keyword("overflow", map[string]string{
		"auto": "auto", "hidden": "hidden", "clip": "clip",
		"visible": "visible", "scroll": "scroll",
	}))
	r.prefix("overflow-x", 166, keyword("overflow-x", map[string]string{
		"auto": "auto", "hidden": "hidden", "clip": "clip",
		"visible": "visible", "scroll": "scroll",
	}))
	r.prefix("overflow-y", 166, keyword("overflow-y", map[string]string{
		"auto": "auto", "hidden": "hidden", "clip": "clip",
		"visible": "visible", "scroll": "scroll",
	}))

	// Box sizing.
	r.prefix("box", 17, keyword("box-sizing", map[string]string{
		"border": "border-box", "content": "content-box",
	}))
}
