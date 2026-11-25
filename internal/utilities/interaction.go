package utilities

import (
	"github.com/Ahaduzzamankhan/flexui/internal/ast"
)

func registerInteraction(r *Registry) {
	r.prefix("cursor", 160, keyword("cursor", map[string]string{
		"auto": "auto", "default": "default", "pointer": "pointer",
		"wait": "wait", "text": "text", "move": "move", "help": "help",
		"not-allowed": "not-allowed", "none": "none", "grab": "grab",
		"grabbing": "grabbing", "zoom-in": "zoom-in", "zoom-out": "zoom-out",
		"crosshair": "crosshair", "progress": "progress", "cell": "cell",
		"copy": "copy", "alias": "alias", "all-scroll": "all-scroll",
		"col-resize": "col-resize", "row-resize": "row-resize",
		"e-resize": "e-resize", "w-resize": "w-resize",
		"n-resize": "n-resize", "s-resize": "s-resize",
		"ne-resize": "ne-resize", "nw-resize": "nw-resize",
		"se-resize": "se-resize", "sw-resize": "sw-resize",
	}))
	r.exactOnly("pointer-events-none", 161, single("pointer-events", "none"))
	r.exactOnly("pointer-events-auto", 161, single("pointer-events", "auto"))
	r.prefix("select", 162, keyword("user-select", map[string]string{
		"none": "none", "text": "text", "all": "all", "auto": "auto",
	}))
	r.prefix("object", 163, keyword("object-fit", map[string]string{
		"contain": "contain", "cover": "cover", "fill": "fill",
		"none": "none", "scale-down": "scale-down",
	}))
	r.prefix("object-position", 164, keyword("object-position", map[string]string{
		"center": "center", "top": "top", "bottom": "bottom",
		"left": "left", "right": "right", "top-left": "top left",
		"top-right": "top right", "bottom-left": "bottom left",
		"bottom-right": "bottom right",
	}))
	r.prefix("resize", 167, keyword("resize", map[string]string{
		"none": "none", "y": "vertical", "x": "horizontal", "both": "both",
	}))
	r.prefix("appearance", 168, keyword("appearance", map[string]string{
		"none": "none", "auto": "auto",
	}))

	// Transitions.
	r.both("transition", 210, func(_ *Context, value string, _ bool) []ast.Block {
		switch value {
		case "":
			return []ast.Block{on(
				decl("transition-property", "color, background-color, border-color, text-decoration-color, fill, stroke, opacity, box-shadow, transform, translate, scale, rotate, filter, backdrop-filter"),
				decl("transition-timing-function", "cubic-bezier(0.4, 0, 0.2, 1)"),
				decl("transition-duration", "150ms"),
			)}
		case "all":
			return []ast.Block{on(
				decl("transition-property", "all"),
				decl("transition-timing-function", "cubic-bezier(0.4, 0, 0.2, 1)"),
				decl("transition-duration", "150ms"),
			)}
		case "none":
			return []ast.Block{on(decl("transition-property", "none"))}
		case "colors", "color":
			return []ast.Block{on(
				decl("transition-property", "color, background-color, border-color, text-decoration-color, fill, stroke"),
				decl("transition-timing-function", "cubic-bezier(0.4, 0, 0.2, 1)"),
				decl("transition-duration", "150ms"),
			)}
		case "opacity":
			return []ast.Block{on(decl("transition-property", "opacity"), decl("transition-timing-function", "cubic-bezier(0.4, 0, 0.2, 1)"), decl("transition-duration", "150ms"))}
		case "shadow":
			return []ast.Block{on(decl("transition-property", "box-shadow"), decl("transition-timing-function", "cubic-bezier(0.4, 0, 0.2, 1)"), decl("transition-duration", "150ms"))}
		case "transform":
			return []ast.Block{on(decl("transition-property", "transform, translate, scale, rotate"), decl("transition-timing-function", "cubic-bezier(0.4, 0, 0.2, 1)"), decl("transition-duration", "150ms"))}
		case "size":
			return []ast.Block{on(decl("transition-property", "width, height, min-width, max-width, min-height, max-height"), decl("transition-timing-function", "cubic-bezier(0.4, 0, 0.2, 1)"), decl("transition-duration", "150ms"))}
		}
		if v, ok := arb(value); ok {
			return []ast.Block{on(decl("transition-property", v))}
		}
		return nil
	})
	r.prefix("duration", 211, func(c *Context, value string, _ bool) []ast.Block {
		if v, ok := c.Theme.Duration[value]; ok {
			return []ast.Block{on(decl("transition-duration", v))}
		}
		if v, ok := numValue(value); ok {
			return []ast.Block{on(decl("transition-duration", v+"ms"))}
		}
		return nil
	})
	r.prefix("delay", 212, func(c *Context, value string, _ bool) []ast.Block {
		if v, ok := c.Theme.Duration[value]; ok {
			return []ast.Block{on(decl("transition-delay", v))}
		}
		if v, ok := numValue(value); ok {
			return []ast.Block{on(decl("transition-delay", v+"ms"))}
		}
		return nil
	})
	r.prefix("ease", 213, func(c *Context, value string, _ bool) []ast.Block {
		if v, ok := c.Theme.Easing[value]; ok {
			return []ast.Block{on(decl("transition-timing-function", v))}
		}
		if v, ok := arb(value); ok {
			return []ast.Block{on(decl("transition-timing-function", v))}
		}
		return nil
	})

	// Container shell used by layout presets.
	r.exactOnly("container", 214, func(*Context, string, bool) []ast.Block {
		return []ast.Block{on(
			decl("width", "100%"),
			decl("max-width", "var(--flexui-container-width, 80rem)"),
			decl("margin-inline", "auto"),
			decl("padding-inline", "var(--flexui-container-gutter, 1rem)"),
		)}
	})
	r.exactOnly("will-change-transform", 215, single("will-change", "transform"))
	r.exactOnly("will-change-auto", 215, single("will-change", "auto"))
	r.exactOnly("will-change-opacity", 215, single("will-change", "opacity"))
}
