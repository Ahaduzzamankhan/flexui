// Package components implements FlexUI's component layer.
//
// Components are plain CSS classes in the `components` cascade layer, so
// utilities still override them and no JavaScript is required to render one.
// Only classes that appear in the scanned sources are emitted, because each
// component is resolved from the same candidate set as utilities.
//
// Interactive behaviour (dialogs, menus, tabs, toasts) lives in the optional
// TypeScript runtime and is expressed through ARIA attributes and data
// attributes that these rules style. A component without the runtime is still
// usable and accessible wherever the platform provides the behaviour (native
// <dialog>, <details>, form controls).
package components

import (
	"sort"

	"github.com/Ahaduzzamankhan/flexui/internal/ast"
	"github.com/Ahaduzzamankhan/flexui/internal/themes"
)

// Result is the resolved CSS for one component class.
type Result struct {
	Name      string
	Order     int
	Blocks    []ast.Block
	Keyframes []ast.Keyframe
}

// Registry maps component class names to their declarations.
type Registry struct {
	theme   *themes.Theme
	entries map[string]Result
	names   []string
}

// New builds the component registry for a theme.
func New(t *themes.Theme) *Registry {
	r := &Registry{theme: t, entries: make(map[string]Result, 128)}
	registerFoundation(r)
	registerForms(r)
	registerNavigation(r)
	registerOverlays(r)
	registerData(r)
	registerLayouts(r)
	// Keyframes that the component styles reference directly.
	r.attach("fui-spinner", keyframesSpin)
	r.attach("fui-skeleton", keyframesSkeleton)
	r.attach("fui-dialog", keyframesDialog)
	r.attach("fui-drawer", keyframesDrawer)
	r.attach("fui-toast", keyframesToast)
	r.names = make([]string, 0, len(r.entries))
	for name := range r.entries {
		r.names = append(r.names, name)
	}
	sort.Strings(r.names)
	return r
}

// Register adds a component class. Declarations are applied to the element
// itself; use Target templates ("& > svg", "&[aria-selected='true']") for
// child parts and states.
func (r *Registry) Register(name string, order int, blocks ...ast.Block) {
	r.entries[name] = Result{Name: name, Order: order, Blocks: blocks}
}

// Resolve returns a component by exact class name.
func (r *Registry) Resolve(name string) (Result, bool) {
	result, ok := r.entries[name]
	return result, ok
}

// Names returns every registered class name in sorted order.
func (r *Registry) Names() []string {
	out := make([]string, len(r.names))
	copy(out, r.names)
	return out
}

// Len returns the number of registered component classes.
func (r *Registry) Len() int { return len(r.entries) }

// Decl builds a declaration.
func Decl(prop, value string) ast.Decl { return ast.Decl{Prop: prop, Value: value} }

// At builds a block wrapped in a conditional group rule.
func At(condition string, decls ...ast.Decl) ast.Block {
	return ast.Block{Target: "&", AtRule: condition, Decls: decls}
}

// attach adds keyframes to a registered class.
func (r *Registry) attach(name string, kfs ...ast.Keyframe) {
	entry, ok := r.entries[name]
	if !ok {
		return
	}
	entry.Keyframes = kfs
	r.entries[name] = entry
}

// Self returns a block for the component element itself.
func Self(decls ...ast.Decl) ast.Block { return ast.Block{Decls: decls} }

// Part returns a block for a descendant or state, using "&" for the element.
func Part(target string, decls ...ast.Decl) ast.Block {
	return ast.Block{Target: target, Decls: decls}
}

var (
	keyframesSpin = ast.Keyframe{Name: "fui-spin", Steps: []ast.KeyframeStep{
		{Steps: "to", Decls: []ast.Decl{{Prop: "transform", Value: "rotate(360deg)"}}},
	}}
	keyframesSkeleton = ast.Keyframe{Name: "fui-skeleton-pulse", Steps: []ast.KeyframeStep{
		{Steps: "0%, 100%", Decls: []ast.Decl{{Prop: "opacity", Value: "1"}}},
		{Steps: "50%", Decls: []ast.Decl{{Prop: "opacity", Value: "0.55"}}},
	}}
	keyframesDialog = ast.Keyframe{Name: "fui-dialog-in", Steps: []ast.KeyframeStep{
		{Steps: "from", Decls: []ast.Decl{{Prop: "opacity", Value: "0"}, {Prop: "transform", Value: "translateY(8px) scale(0.98)"}}},
		{Steps: "to", Decls: []ast.Decl{{Prop: "opacity", Value: "1"}, {Prop: "transform", Value: "translateY(0) scale(1)"}}},
	}}
	keyframesDrawer = ast.Keyframe{Name: "fui-drawer-in", Steps: []ast.KeyframeStep{
		{Steps: "from", Decls: []ast.Decl{{Prop: "transform", Value: "translateX(100%)"}}},
		{Steps: "to", Decls: []ast.Decl{{Prop: "transform", Value: "translateX(0)"}}},
	}}
	keyframesToast = ast.Keyframe{Name: "fui-toast-in", Steps: []ast.KeyframeStep{
		{Steps: "from", Decls: []ast.Decl{{Prop: "opacity", Value: "0"}, {Prop: "transform", Value: "translateY(8px)"}}},
		{Steps: "to", Decls: []ast.Decl{{Prop: "opacity", Value: "1"}, {Prop: "transform", Value: "translateY(0)"}}},
	}}
)

// Order buckets keep component layers stable and readable: foundation, forms,
// navigation, overlays, data display, layouts.
const (
	orderFoundation = 10
	orderForms      = 100
	orderNav        = 200
	orderOverlay    = 300
	orderData       = 400
	orderLayout     = 500
)
