package generator

import "github.com/Ahaduzzamankhan/flexui/internal/ast"

// baseRule is one preflight rule. Preflight is kept as structured data, not a
// CSS string, so minified and pretty output share one code path and stay in
// sync.
type baseRule struct {
	selector string
	decls    []ast.Decl
}

func d(prop, value string) ast.Decl { return ast.Decl{Prop: prop, Value: value} }

var baseRules = []baseRule{
	{
		selector: "*, ::before, ::after",
		decls: []ast.Decl{
			d("box-sizing", "border-box"),
			d("border", "0 solid currentColor"),
		},
	},
	{
		selector: "html",
		decls: []ast.Decl{
			d("-webkit-text-size-adjust", "100%"),
			d("text-size-adjust", "100%"),
			d("tab-size", "4"),
			d("font-family", `ui-sans-serif, system-ui, -apple-system, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif`),
			d("line-height", "1.5"),
		},
	},
	{
		selector: "body",
		decls: []ast.Decl{
			d("margin", "0"),
			d("min-height", "100vh"),
			d("text-rendering", "optimizeLegibility"),
			d("-webkit-font-smoothing", "antialiased"),
			d("-moz-osx-font-smoothing", "grayscale"),
		},
	},
	{
		selector: "h1, h2, h3, h4, h5, h6",
		decls: []ast.Decl{
			d("font-size", "inherit"),
			d("font-weight", "inherit"),
			d("margin", "0"),
		},
	},
	{
		selector: "p, figure, blockquote, dl, dd, pre, hr",
		decls: []ast.Decl{
			d("margin", "0"),
		},
	},
	{
		selector: "ul, ol",
		decls: []ast.Decl{
			d("margin", "0"),
			d("padding", "0"),
			d("list-style", "none"),
		},
	},
	{
		selector: "img, svg, video, canvas, audio, iframe, embed, object",
		decls: []ast.Decl{
			d("display", "block"),
			d("vertical-align", "middle"),
		},
	},
	{
		selector: "img, video",
		decls: []ast.Decl{
			d("max-width", "100%"),
			d("height", "auto"),
		},
	},
	{
		selector: "button, input, optgroup, select, textarea",
		decls: []ast.Decl{
			d("font", "inherit"),
			d("color", "inherit"),
			d("margin", "0"),
			d("padding", "0"),
			d("background", "transparent"),
			d("border-radius", "0"),
		},
	},
	{
		selector: "button, [type=\"button\"], [type=\"submit\"], [type=\"reset\"]",
		decls: []ast.Decl{
			d("appearance", "button"),
			d("cursor", "pointer"),
		},
	},
	{
		selector: "textarea",
		decls: []ast.Decl{
			d("resize", "vertical"),
		},
	},
	{
		selector: "a",
		decls: []ast.Decl{
			d("color", "inherit"),
			d("text-decoration", "inherit"),
		},
	},
	{
		selector: "table",
		decls: []ast.Decl{
			d("border-collapse", "collapse"),
			d("border-color", "inherit"),
			d("text-indent", "0"),
		},
	},
	{
		selector: "summary",
		decls: []ast.Decl{
			d("display", "list-item"),
		},
	},
	{
		selector: "[hidden]",
		decls: []ast.Decl{
			d("display", "none"),
		},
	},
	{
		selector: "::selection",
		decls: []ast.Decl{
			d("background-color", "#b4d7ff"),
			d("color", "#0b1729"),
		},
	},
}

// BaseRules exposes the preflight rule set for documentation and tests.
func BaseRuleCount() int { return len(baseRules) }
