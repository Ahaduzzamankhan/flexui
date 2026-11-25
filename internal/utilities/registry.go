// Package utilities implements FlexUI's utility registry.
//
// A utility is either an exact name ("flex", "hidden") or a prefix with a
// value ("p-4", "bg-red-500"). Prefixes are matched longest-first so that
// "border-t-2" resolves through "border-t" rather than "border". Every entry
// carries a stable cascade order number; the generator sorts rules by it, which
// is what makes output deterministic and keeps shorthand utilities such as
// "p-4" ahead of their longhand counterparts such as "pt-2".
package utilities

import (
	"sort"
	"strings"

	"github.com/Ahaduzzamankhan/flexui/internal/ast"
	"github.com/Ahaduzzamankhan/flexui/internal/themes"
)

// Context carries the theme for resolution.
type Context struct {
	Theme *themes.Theme
}

// resolveFunc produces declaration blocks for a utility value. An empty value
// means the utility was written without one ("flex", "border").
type resolveFunc func(c *Context, value string, neg bool) []ast.Block

type entry struct {
	prefix string
	order  int
	fn     resolveFunc
}

// Registry holds every known utility.
type Registry struct {
	ctx      *Context
	prefixes []string
	byPrefix map[string]entry
	exact    map[string]entry
}

// New builds the default registry for a theme.
func New(t *themes.Theme) *Registry {
	r := &Registry{
		ctx:      &Context{Theme: t},
		byPrefix: make(map[string]entry, 256),
		exact:    make(map[string]entry, 128),
	}
	registerSpacing(r)
	registerLayout(r)
	registerFlexGrid(r)
	registerSizing(r)
	registerPosition(r)
	registerColors(r)
	registerTypography(r)
	registerBorders(r)
	registerEffects(r)
	registerTransforms(r)
	registerInteraction(r)
	registerExtra(r)
	r.finish()
	return r
}

// prefix registers a value-carrying utility.
func (r *Registry) prefix(name string, order int, fn resolveFunc) {
	r.byPrefix[name] = entry{prefix: name, order: order, fn: fn}
}

// exact registers an exact-name utility.
func (r *Registry) exactOnly(name string, order int, fn resolveFunc) {
	r.exact[name] = entry{prefix: name, order: order, fn: fn}
}

// both registers a utility that is valid with and without a value.
func (r *Registry) both(name string, order int, fn resolveFunc) {
	r.prefix(name, order, fn)
	r.exactOnly(name, order, fn)
}

func (r *Registry) finish() {
	r.prefixes = make([]string, 0, len(r.byPrefix))
	for p := range r.byPrefix {
		r.prefixes = append(r.prefixes, p)
	}
	sort.Slice(r.prefixes, func(i, j int) bool {
		if len(r.prefixes[i]) != len(r.prefixes[j]) {
			return len(r.prefixes[i]) > len(r.prefixes[j])
		}
		return r.prefixes[i] < r.prefixes[j]
	})
}

// Resolve looks up a utility token. It returns the declaration blocks, the
// canonical utility key and its cascade order.
func (r *Registry) Resolve(token string, neg bool) (blocks []ast.Block, key string, order int, ok bool) {
	if e, found := r.exact[token]; found {
		if b := e.fn(r.ctx, "", neg); len(b) > 0 {
			return b, e.prefix, e.order, true
		}
		return nil, "", 0, false
	}
	for _, p := range r.prefixes {
		if len(token) <= len(p)+1 || !strings.HasPrefix(token, p) || token[len(p)] != '-' {
			continue
		}
		e := r.byPrefix[p]
		if b := e.fn(r.ctx, token[len(p)+1:], neg); len(b) > 0 {
			return b, e.prefix, e.order, true
		}
	}
	return nil, "", 0, false
}

// Known reports whether the token resolves to at least one declaration.
func (r *Registry) Known(token string, neg bool) bool {
	_, _, _, ok := r.Resolve(token, neg)
	return ok
}

// Theme exposes the theme used by this registry.
func (r *Registry) Theme() *themes.Theme { return r.ctx.Theme }

// Context exposes the resolution context.
func (r *Registry) Context() *Context { return r.ctx }
