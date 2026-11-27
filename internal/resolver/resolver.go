// Package resolver turns a parsed candidate into an AST rule by combining the
// utility registry, the component registry and the variant resolver.
package resolver

import (
	"strings"

	"github.com/Ahaduzzamankhan/flexui/internal/animation"
	"github.com/Ahaduzzamankhan/flexui/internal/ast"
	"github.com/Ahaduzzamankhan/flexui/internal/components"
	"github.com/Ahaduzzamankhan/flexui/internal/parser"
	"github.com/Ahaduzzamankhan/flexui/internal/themes"
	"github.com/Ahaduzzamankhan/flexui/internal/utilities"
	"github.com/Ahaduzzamankhan/flexui/internal/variants"
)

// Options configures a resolver.
type Options struct {
	Theme *themes.Theme
	// DarkMode is "class", "media" or "both".
	DarkMode string
	// CustomVariants are user-defined variants from configuration.
	CustomVariants map[string]string
	// CustomAnimations are user-defined @keyframes from configuration.
	CustomAnimations map[string]animation.Definition
	// CustomComponents are user-defined components from configuration. Their
	// classes resolve in the components layer exactly like built-in ones.
	CustomComponents []components.CustomComponent
}

// Resolver resolves candidates.
type Resolver struct {
	utils *utilities.Registry
	comps *components.Registry
	anims *animation.Registry
	vars  *variants.Resolver
	theme *themes.Theme
}

// New builds a resolver with default options.
func New(t *themes.Theme) *Resolver {
	return NewWith(Options{Theme: t})
}

// NewWith builds a resolver with explicit options.
func NewWith(o Options) *Resolver {
	theme := o.Theme
	if theme == nil {
		theme = themes.Default()
	}
	comps, err := components.NewWithCustom(theme, o.CustomComponents)
	if err != nil {
		// Invalid definitions are rejected by the compiler before resolution, so
		// this only guards direct library use; the built-in set stays available.
		comps = components.New(theme)
	}
	return &Resolver{
		utils: utilities.New(theme),
		comps: comps,
		anims: animation.New(o.CustomAnimations),
		vars: variants.NewWith(variants.Options{
			Theme:    theme,
			DarkMode: o.DarkMode,
			Custom:   o.CustomVariants,
		}),
		theme: theme,
	}
}

// Utilities exposes the utility registry.
func (r *Resolver) Utilities() *utilities.Registry { return r.utils }

// Components exposes the component registry.
func (r *Resolver) Components() *components.Registry { return r.comps }

// Resolve resolves one raw candidate. The boolean reports whether the candidate
// is a known FlexUI class; unknown candidates are silently skipped by the
// compiler, exactly like unknown words in markup.
func (r *Resolver) Resolve(raw string) (ast.Rule, bool) {
	c, ok := parser.Parse(raw)
	if !ok {
		return ast.Rule{}, false
	}
	var blocks []ast.Block
	var key, layer string
	var order int
	var keyframes []ast.Keyframe
	switch {
	case c.ArbitraryProperty:
		prop, value, found := splitProperty(c.Utility)
		if !found {
			return ast.Rule{}, false
		}
		blocks = []ast.Block{{Decls: []ast.Decl{{Prop: prop, Value: value}}}}
		key, order = "arbitrary-property", 500
	default:
		if comp, found := r.comps.Resolve(c.Utility); found {
			blocks = comp.Blocks
			keyframes = comp.Keyframes
			key = comp.Name
			order = comp.Order
			layer = ast.LayerComponents
		} else if animBlocks, animKF, animKey, animOrder, found := r.anims.Resolve(c.Utility); found {
			blocks = animBlocks
			keyframes = animKF
			key = animKey
			order = animOrder
		} else {
			var found bool
			blocks, key, order, found = r.utils.Resolve(c.Utility, c.Negative)
			if !found {
				return ast.Rule{}, false
			}
		}
	}
	applied := make([]ast.AppliedVariant, 0, len(c.Variants))
	for _, v := range c.Variants {
		av, ok := r.vars.Resolve(v)
		if !ok {
			return ast.Rule{}, false
		}
		applied = append(applied, av)
	}
	if c.Important {
		for bi := range blocks {
			for di := range blocks[bi].Decls {
				if !strings.HasSuffix(blocks[bi].Decls[di].Value, "!important") {
					blocks[bi].Decls[di].Value += " !important"
				}
			}
		}
	}
	return ast.Rule{
		Class:        c.Raw,
		Utility:      key,
		UtilityOrder: order,
		Layer:        layer,
		Variants:     applied,
		Blocks:       blocks,
		Keyframes:    keyframes,
	}, true
}

// splitProperty splits an arbitrary property utility body into property and
// value at the first top-level colon.
func splitProperty(utility string) (string, string, bool) {
	if len(utility) < 3 || utility[0] != '[' || utility[len(utility)-1] != ']' {
		return "", "", false
	}
	inner := utility[1 : len(utility)-1]
	depth := 0
	for i := 0; i < len(inner); i++ {
		switch inner[i] {
		case '[':
			depth++
		case ']':
			depth--
		case ':':
			if depth == 0 {
				prop := strings.TrimSpace(inner[:i])
				value := parser.NormalizeArbitrary(strings.TrimSpace(inner[i+1:]))
				if prop == "" || value == "" {
					return "", "", false
				}
				return prop, value, true
			}
		}
	}
	return "", "", false
}

// Theme exposes the theme.
func (r *Resolver) Theme() *themes.Theme { return r.theme }
