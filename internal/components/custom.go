package components

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Ahaduzzamankhan/flexui/internal/ast"
	"github.com/Ahaduzzamankhan/flexui/internal/themes"
)

// CustomComponent is a user-defined component from flexui.config.json.
//
// A component produces:
//
//	fui-<name>              base declarations and states
//	fui-<name>--<variant>   variant declarations
//	fui-<name>-<part>       part declarations
//
// This is the extension point that keeps the component system open: projects
// are not limited to the built-in classes, and a custom component with the same
// name replaces a built-in one so the built-in set stays customisable.
type CustomComponent struct {
	Name     string
	Extends  string
	Base     map[string]string
	Variants map[string]map[string]string
	Parts    map[string]map[string]string
	States   map[string]map[string]string
}

// Order buckets for custom components. All custom base classes share one order
// and sort by class name, so output stays deterministic across runs and
// platforms.
const (
	orderCustomBase    = 600
	orderCustomVariant = 601
	orderCustomPart    = 602
)

type resolvedCustom struct {
	base     map[string]string
	variants map[string]map[string]string
	parts    map[string]map[string]string
	states   map[string]map[string]string
}

// NewWithCustom builds the built-in registry and registers user components.
// Invalid definitions (unknown parent, extends cycle) are returned as errors.
func NewWithCustom(t *themes.Theme, customs []CustomComponent) (*Registry, error) {
	r := New(t)
	if err := r.RegisterCustomComponents(customs); err != nil {
		return nil, err
	}
	return r, nil
}

// ValidateCustom reports whether the definitions can be resolved against the
// built-in component set, without keeping the registry.
func ValidateCustom(customs []CustomComponent) error {
	_, err := NewWithCustom(themes.Default(), customs)
	return err
}

// RegisterCustomComponents resolves and registers every definition.
func (r *Registry) RegisterCustomComponents(customs []CustomComponent) error {
	resolved, err := buildCustom(r, customs)
	if err != nil {
		return err
	}

	names := sortedKeys(resolved)
	for _, name := range names {
		component := resolved[name]
		r.Register("fui-"+name, orderCustomBase, customBlocks(component)...)

		for _, variant := range sortedKeys(component.variants) {
			r.Register("fui-"+name+"--"+variant, orderCustomVariant, Self(declsOf(component.variants[variant])...))
		}
		for _, part := range sortedKeys(component.parts) {
			r.Register("fui-"+name+"-"+part, orderCustomPart, Self(declsOf(component.parts[part])...))
		}
	}
	r.refreshNames()
	return nil
}

// customBlocks returns the base element block plus one block per state.
func customBlocks(component resolvedCustom) []ast.Block {
	blocks := make([]ast.Block, 0, 1+len(component.states))
	blocks = append(blocks, Self(declsOf(component.base)...))
	for _, state := range sortedKeys(component.states) {
		blocks = append(blocks, Part("&"+state, declsOf(component.states[state])...))
	}
	return blocks
}

// buildCustom resolves extends chains into fully merged declarations.
func buildCustom(registry *Registry, customs []CustomComponent) (map[string]resolvedCustom, error) {
	byName := make(map[string]CustomComponent, len(customs))
	for _, component := range customs {
		if component.Name == "" {
			return nil, fmt.Errorf("component with no name")
		}
		if _, exists := byName[component.Name]; exists {
			return nil, fmt.Errorf("component %q is defined more than once", component.Name)
		}
		byName[component.Name] = component
	}

	resolved := make(map[string]resolvedCustom, len(byName))
	resolving := make(map[string]bool, len(byName))

	var resolve func(name string) error
	resolve = func(name string) error {
		if _, done := resolved[name]; done {
			return nil
		}
		if resolving[name] {
			return fmt.Errorf("component %q has a circular extends chain", name)
		}
		component, ok := byName[name]
		if !ok {
			return fmt.Errorf("component %q is referenced but not defined", name)
		}
		resolving[name] = true
		defer delete(resolving, name)

		result := resolvedCustom{
			base:     cloneStrings(component.Base),
			variants: cloneGroups(component.Variants),
			parts:    cloneGroups(component.Parts),
			states:   cloneGroups(component.States),
		}

		if parent := component.Extends; parent != "" {
			if _, isCustom := byName[parent]; isCustom {
				if err := resolve(parent); err != nil {
					return err
				}
				inherited := resolved[parent]
				result.base = mergeStrings(inherited.base, result.base)
				result.variants = mergeGroups(inherited.variants, result.variants)
				result.parts = mergeGroups(inherited.parts, result.parts)
				result.states = mergeGroups(inherited.states, result.states)
			} else {
				decls, found := builtinElementDecls(registry, parent)
				if !found {
					return fmt.Errorf("component %q extends %q, which is not a custom or built-in component", name, parent)
				}
				result.base = mergeStrings(decls, result.base)
			}
		}

		resolved[name] = result
		return nil
	}

	names := sortedKeys(byName)
	for _, name := range names {
		if err := resolve(name); err != nil {
			return nil, err
		}
	}
	return resolved, nil
}

// builtinElementDecls returns the declarations a built-in component applies to
// its own element, so `extends` can reuse a built-in's foundation without
// inheriting its states or child parts.
func builtinElementDecls(registry *Registry, name string) (map[string]string, bool) {
	if registry == nil {
		return nil, false
	}
	class := name
	if !strings.HasPrefix(class, "fui-") {
		class = "fui-" + name
	}
	entry, ok := registry.Resolve(class)
	if !ok {
		return nil, false
	}
	decls := make(map[string]string)
	for _, block := range entry.Blocks {
		if block.Target != "" || block.AtRule != "" {
			continue
		}
		for _, decl := range block.Decls {
			decls[decl.Prop] = decl.Value
		}
	}
	return decls, len(decls) > 0
}

// declsOf converts a property map into sorted declarations.
func declsOf(values map[string]string) []ast.Decl {
	props := sortedKeys(values)
	out := make([]ast.Decl, 0, len(props))
	for _, prop := range props {
		out = append(out, Decl(prop, values[prop]))
	}
	return out
}

func cloneStrings(values map[string]string) map[string]string {
	out := make(map[string]string, len(values))
	for key, value := range values {
		out[key] = value
	}
	return out
}

func cloneGroups(values map[string]map[string]string) map[string]map[string]string {
	out := make(map[string]map[string]string, len(values))
	for key, group := range values {
		out[key] = cloneStrings(group)
	}
	return out
}

// mergeStrings returns parent values overridden by child values.
func mergeStrings(parent, child map[string]string) map[string]string {
	out := cloneStrings(parent)
	for key, value := range child {
		out[key] = value
	}
	return out
}

// mergeGroups merges each variant/part/state, letting child declarations win.
func mergeGroups(parent, child map[string]map[string]string) map[string]map[string]string {
	out := cloneGroups(parent)
	for key, group := range child {
		out[key] = mergeStrings(out[key], group)
	}
	return out
}

func sortedKeys[V any](values map[string]V) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// refreshNames rebuilds the sorted class-name index after dynamic registration.
func (r *Registry) refreshNames() {
	r.names = sortedKeys(r.entries)
}
