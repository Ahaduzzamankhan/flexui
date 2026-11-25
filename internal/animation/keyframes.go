package animation

import "github.com/Ahaduzzamankhan/flexui/internal/ast"

// builtinOrder lists the built-in animations in documentation order.
var builtinOrder = []string{
	"fade-in", "fade-out", "fade-up", "fade-down", "fade-left", "fade-right",
	"slide-in", "slide-out", "zoom-in", "zoom-out", "scale-up", "scale-down",
	"rotate-in", "rotate-out", "bounce", "pulse", "ping", "spin", "shake",
	"wobble", "float", "shimmer", "blur-in", "blur-out", "flip",
}

// builtins holds every built-in animation definition.
var builtins = map[string]Definition{
	"fade-in": {
		Steps: []ast.KeyframeStep{
			{Steps: "from", Decls: decl("opacity", "0")},
			{Steps: "to", Decls: decl("opacity", "1")},
		},
	},
	"fade-out": {
		Steps: []ast.KeyframeStep{
			{Steps: "from", Decls: decl("opacity", "1")},
			{Steps: "to", Decls: decl("opacity", "0")},
		},
	},
	"fade-up": {
		Steps: []ast.KeyframeStep{
			{Steps: "from", Decls: decl("opacity", "0", "transform", "translateY(16px)")},
			{Steps: "to", Decls: decl("opacity", "1", "transform", "translateY(0)")},
		},
	},
	"fade-down": {
		Steps: []ast.KeyframeStep{
			{Steps: "from", Decls: decl("opacity", "0", "transform", "translateY(-16px)")},
			{Steps: "to", Decls: decl("opacity", "1", "transform", "translateY(0)")},
		},
	},
	"fade-left": {
		Steps: []ast.KeyframeStep{
			{Steps: "from", Decls: decl("opacity", "0", "transform", "translateX(16px)")},
			{Steps: "to", Decls: decl("opacity", "1", "transform", "translateX(0)")},
		},
	},
	"fade-right": {
		Steps: []ast.KeyframeStep{
			{Steps: "from", Decls: decl("opacity", "0", "transform", "translateX(-16px)")},
			{Steps: "to", Decls: decl("opacity", "1", "transform", "translateX(0)")},
		},
	},
	"slide-in": {
		Steps: []ast.KeyframeStep{
			{Steps: "from", Decls: decl("opacity", "0", "transform", "translateY(24px)")},
			{Steps: "to", Decls: decl("opacity", "1", "transform", "translateY(0)")},
		},
	},
	"slide-out": {
		Steps: []ast.KeyframeStep{
			{Steps: "from", Decls: decl("opacity", "1", "transform", "translateY(0)")},
			{Steps: "to", Decls: decl("opacity", "0", "transform", "translateY(24px)")},
		},
	},
	"zoom-in": {
		Steps: []ast.KeyframeStep{
			{Steps: "from", Decls: decl("opacity", "0", "transform", "scale(0.92)")},
			{Steps: "to", Decls: decl("opacity", "1", "transform", "scale(1)")},
		},
	},
	"zoom-out": {
		Steps: []ast.KeyframeStep{
			{Steps: "from", Decls: decl("opacity", "0", "transform", "scale(1.08)")},
			{Steps: "to", Decls: decl("opacity", "1", "transform", "scale(1)")},
		},
	},
	"scale-up": {
		Steps: []ast.KeyframeStep{
			{Steps: "from", Decls: decl("transform", "scale(0.95)")},
			{Steps: "to", Decls: decl("transform", "scale(1)")},
		},
		Duration: "200ms",
	},
	"scale-down": {
		Steps: []ast.KeyframeStep{
			{Steps: "from", Decls: decl("transform", "scale(1.05)")},
			{Steps: "to", Decls: decl("transform", "scale(1)")},
		},
		Duration: "200ms",
	},
	"rotate-in": {
		Steps: []ast.KeyframeStep{
			{Steps: "from", Decls: decl("opacity", "0", "transform", "rotate(-8deg) scale(0.96)")},
			{Steps: "to", Decls: decl("opacity", "1", "transform", "rotate(0) scale(1)")},
		},
	},
	"rotate-out": {
		Steps: []ast.KeyframeStep{
			{Steps: "from", Decls: decl("opacity", "1", "transform", "rotate(0)")},
			{Steps: "to", Decls: decl("opacity", "0", "transform", "rotate(8deg)")},
		},
	},
	"bounce": {
		Steps: []ast.KeyframeStep{
			{Steps: "0%, 100%", Decls: decl("transform", "translateY(0)")},
			{Steps: "50%", Decls: decl("transform", "translateY(-10px)")},
		},
		Duration: "900ms", Easing: "ease-in-out", Iterations: "infinite", Fill: "none",
	},
	"pulse": {
		Steps: []ast.KeyframeStep{
			{Steps: "0%, 100%", Decls: decl("opacity", "1")},
			{Steps: "50%", Decls: decl("opacity", "0.55")},
		},
		Duration: "2s", Easing: "ease-in-out", Iterations: "infinite", Fill: "none",
	},
	"ping": {
		Steps: []ast.KeyframeStep{
			{Steps: "0%", Decls: decl("transform", "scale(1)", "opacity", "0.8")},
			{Steps: "75%, 100%", Decls: decl("transform", "scale(1.9)", "opacity", "0")},
		},
		Duration: "1.2s", Easing: "cubic-bezier(0, 0, 0.2, 1)", Iterations: "infinite", Fill: "none",
	},
	"spin": {
		Steps: []ast.KeyframeStep{
			{Steps: "to", Decls: decl("transform", "rotate(360deg)")},
		},
		Duration: "800ms", Easing: "linear", Iterations: "infinite", Fill: "none",
	},
	"shake": {
		Steps: []ast.KeyframeStep{
			{Steps: "0%, 100%", Decls: decl("transform", "translateX(0)")},
			{Steps: "25%", Decls: decl("transform", "translateX(-4px)")},
			{Steps: "75%", Decls: decl("transform", "translateX(4px)")},
		},
		Duration: "400ms",
	},
	"wobble": {
		Steps: []ast.KeyframeStep{
			{Steps: "0%, 100%", Decls: decl("transform", "rotate(0)")},
			{Steps: "25%", Decls: decl("transform", "rotate(-3deg)")},
			{Steps: "75%", Decls: decl("transform", "rotate(3deg)")},
		},
		Duration: "500ms",
	},
	"float": {
		Steps: []ast.KeyframeStep{
			{Steps: "0%, 100%", Decls: decl("transform", "translateY(0)")},
			{Steps: "50%", Decls: decl("transform", "translateY(-6px)")},
		},
		Duration: "3s", Easing: "ease-in-out", Iterations: "infinite", Fill: "none",
	},
	"shimmer": {
		Steps: []ast.KeyframeStep{
			{Steps: "from", Decls: decl("background-position", "-200% 0")},
			{Steps: "to", Decls: decl("background-position", "200% 0")},
		},
		Duration: "1.6s", Easing: "linear", Iterations: "infinite", Fill: "none",
	},
	"blur-in": {
		Steps: []ast.KeyframeStep{
			{Steps: "from", Decls: decl("opacity", "0", "filter", "blur(8px)")},
			{Steps: "to", Decls: decl("opacity", "1", "filter", "blur(0)")},
		},
	},
	"blur-out": {
		Steps: []ast.KeyframeStep{
			{Steps: "from", Decls: decl("opacity", "1", "filter", "blur(0)")},
			{Steps: "to", Decls: decl("opacity", "0", "filter", "blur(8px)")},
		},
	},
	"flip": {
		Steps: []ast.KeyframeStep{
			{Steps: "from", Decls: decl("transform", "perspective(600px) rotateY(-90deg)")},
			{Steps: "to", Decls: decl("transform", "perspective(600px) rotateY(0)")},
		},
		Duration: "400ms",
	},
}

// decl builds declarations from a prop, value, prop, value... sequence.
func decl(pairs ...string) []ast.Decl {
	out := make([]ast.Decl, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, ast.Decl{Prop: pairs[i], Value: pairs[i+1]})
	}
	return out
}
