// Package animation implements FlexUI's animation engine.
//
// Animations are ordinary CSS: a named @keyframes block plus an `animation`
// declaration whose parts read custom properties. That makes animation
// utilities override one another with utility specificity instead of being
// baked into the keyframes, and lets a single reduced-motion override neutralise
// every animation in the build.
package animation

import (
	"strconv"
	"strings"

	"github.com/Ahaduzzamankhan/flexui/internal/ast"
	"github.com/Ahaduzzamankhan/flexui/internal/themes"
)

// Definition describes one named animation.
type Definition struct {
	Steps      []ast.KeyframeStep
	Duration   string
	Easing     string
	Iterations string
	Fill       string
	Direction  string
	Delay      string
}

// Order is the cascade order of the animation family.
const Order = 480

// Registry holds the built-in animations plus user-defined ones.
type Registry struct {
	defs   map[string]Definition
	custom []string
}

// New builds the registry with the built-in animations and custom definitions.
func New(custom map[string]Definition) *Registry {
	r := &Registry{defs: make(map[string]Definition, len(builtins)+len(custom))}
	for name, def := range builtins {
		// Defaults are applied here, not at resolution time, so every animation
		// has a complete declaration set and never emits a var() reference with
		// no fallback and no value.
		r.defs[name] = withDefaults(def)
	}
	for name, def := range custom {
		name = strings.TrimSpace(name)
		if name == "" || len(def.Steps) == 0 {
			continue
		}
		r.defs[name] = withDefaults(def)
		r.custom = append(r.custom, name)
	}
	return r
}

// Names returns every animation name in the registry.
func (r *Registry) Names() []string {
	out := make([]string, 0, len(r.defs))
	for name := range r.defs {
		out = append(out, name)
	}
	return out
}

// Has reports whether a named animation exists.
func (r *Registry) Has(name string) bool {
	_, ok := r.defs[name]
	return ok
}

// Resolve resolves an `animate-*` utility. It returns the declarations, the
// keyframes the rule needs, a stable key and an order offset so that the
// animation family sorts as a group.
func (r *Registry) Resolve(token string) (blocks []ast.Block, keyframes []ast.Keyframe, key string, order int, ok bool) {
	rest, found := strings.CutPrefix(token, "animate-")
	if !found || rest == "" {
		return nil, nil, "", 0, false
	}
	if rest == "none" {
		return []ast.Block{{Decls: []ast.Decl{
			{Prop: "animation", Value: "none"},
		}}}, nil, "animate-none", Order, true
	}
	if def, ok := r.defs[rest]; ok {
		return animationBlocks(rest, def),
			[]ast.Keyframe{{Name: "fui-" + rest, Steps: def.Steps}},
			"animate-" + rest, Order, true
	}
	// Controls: animate-duration-*, animate-delay-*, animate-ease-* and friends.
	if prop, value, remaining, found := control(rest); found {
		if value == "" {
			return nil, nil, "", 0, false
		}
		return []ast.Block{{Decls: []ast.Decl{{Prop: prop, Value: value}}}},
			nil, "animate-" + remaining, Order + 1, true
	}
	if rest == "stagger" {
		blocks = staggerBlocks()
		return blocks, nil, "animate-stagger", Order + 2, true
	}
	if rest == "timeline-view" {
		return []ast.Block{{Decls: []ast.Decl{
			{Prop: "animation-timeline", Value: "view()"},
			{Prop: "animation-range", Value: "entry 10% cover 60%"},
		}}}, nil, "animate-timeline-view", Order + 3, true
	}
	if rest == "timeline-scroll" {
		return []ast.Block{{Decls: []ast.Decl{
			{Prop: "animation-timeline", Value: "scroll()"},
		}}}, nil, "animate-timeline-scroll", Order + 3, true
	}
	return nil, nil, "", 0, false
}

// control maps an animate-<property>-[value] token onto a declaration.
func control(rest string) (prop, value, key string, ok bool) {
	prefixes := []struct {
		name string
		prop string
	}{
		{"duration-", "animation-duration"},
		{"delay-", "animation-delay"},
		{"ease-", "animation-timing-function"},
		{"iteration-", "animation-iteration-count"},
		{"direction-", "animation-direction"},
		{"fill-", "animation-fill-mode"},
		{"play-", "animation-play-state"},
		{"composition-", "animation-composition"},
	}
	for _, prefix := range prefixes {
		remainder, found := strings.CutPrefix(rest, prefix.name)
		if !found {
			continue
		}
		v, ok := controlValue(prefix.name, remainder)
		if !ok {
			return "", "", "", false
		}
		return prefix.prop, v, strings.TrimSuffix(rest, "-"+remainder), true
	}
	return "", "", "", false
}

// controlValue validates and expands a control value.
func controlValue(kind, raw string) (string, bool) {
	if raw == "" {
		return "", false
	}
	if kind == "ease-" {
		switch raw {
		case "linear":
			return "linear", true
		case "in":
			return "cubic-bezier(0.4, 0, 1, 1)", true
		case "out":
			return "cubic-bezier(0, 0, 0.2, 1)", true
		case "in-out":
			return "cubic-bezier(0.4, 0, 0.2, 1)", true
		case "spring":
			return "cubic-bezier(0.22, 1, 0.36, 1)", true
		}
		return "", false
	}
	if kind == "direction-" {
		switch raw {
		case "normal", "reverse", "alternate", "alternate-reverse":
			return raw, true
		}
		return "", false
	}
	if kind == "fill-" {
		switch raw {
		case "none", "forwards", "backwards", "both":
			return raw, true
		}
		return "", false
	}
	if kind == "play-" {
		switch raw {
		case "running", "paused":
			return raw, true
		}
		return "", false
	}
	if kind == "composition-" {
		switch raw {
		case "replace", "add", "accumulate":
			return raw, true
		}
		return "", false
	}
	if kind == "iteration-" {
		if raw == "infinite" {
			return "infinite", true
		}
		if n, err := strconv.Atoi(raw); err == nil && n >= 0 {
			return strconv.Itoa(n), true
		}
		return "", false
	}
	// duration-* / delay-* take bare milliseconds or an arbitrary value.
	if value, ok := milliseconds(raw); ok {
		return value, true
	}
	if strings.HasPrefix(raw, "[") && strings.HasSuffix(raw, "]") && len(raw) > 2 {
		inner := strings.ReplaceAll(raw[1:len(raw)-1], "_", " ")
		if inner == "" {
			return "", false
		}
		return inner, true
	}
	return "", false
}

// milliseconds turns "150" / "150ms" / "1s" into a CSS time.
func milliseconds(raw string) (string, bool) {
	switch {
	case strings.HasSuffix(raw, "ms"):
		if _, err := strconv.Atoi(strings.TrimSuffix(raw, "ms")); err != nil {
			return "", false
		}
		return raw, true
	case strings.HasSuffix(raw, "s"):
		if _, err := strconv.ParseFloat(strings.TrimSuffix(raw, "s"), 64); err != nil {
			return "", false
		}
		return raw, true
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 0 {
		return "", false
	}
	return strconv.Itoa(n) + "ms", true
}

// animationBlocks builds the declaration set for a named animation. Every part
// reads a custom property with the definition's value as the fallback, so
// `animate-duration-500` can override a single property afterwards.
func animationBlocks(name string, def Definition) []ast.Block {
	return []ast.Block{{Decls: []ast.Decl{
		{Prop: "animation-name", Value: "fui-" + name},
		{Prop: "animation-duration", Value: varOr("--flexui-animation-duration", def.Duration)},
		{Prop: "animation-timing-function", Value: varOr("--flexui-animation-easing", def.Easing)},
		{Prop: "animation-iteration-count", Value: varOr("--flexui-animation-iteration", def.Iterations)},
		{Prop: "animation-fill-mode", Value: varOr("--flexui-animation-fill", def.Fill)},
		{Prop: "animation-direction", Value: varOr("--flexui-animation-direction", def.Direction)},
		{Prop: "animation-delay", Value: varOr("--flexui-animation-delay", def.Delay)},
		{Prop: "animation-play-state", Value: "var(--flexui-animation-play-state, running)"},
	}}}
}

// varOr renders a var() reference with a definition fallback.
func varOr(name, fallback string) string {
	if fallback == "" {
		return "var(" + name + ")"
	}
	return "var(" + name + ", " + fallback + ")"
}

// staggerBlocks delays children 1..12 by an increasing amount.
func staggerBlocks() []ast.Block {
	blocks := make([]ast.Block, 0, 12)
	for i := 1; i <= 12; i++ {
		blocks = append(blocks, ast.Block{
			Target: "& > :nth-child(" + strconv.Itoa(i) + ")",
			Decls:  []ast.Decl{{Prop: "animation-delay", Value: strconv.Itoa((i-1)*60) + "ms"}},
		})
	}
	return blocks
}

// withDefaults fills in the shared defaults for custom animations.
func withDefaults(def Definition) Definition {
	if def.Duration == "" {
		def.Duration = "300ms"
	}
	if def.Easing == "" {
		def.Easing = "cubic-bezier(0.22, 1, 0.36, 1)"
	}
	if def.Iterations == "" {
		def.Iterations = "1"
	}
	if def.Fill == "" {
		def.Fill = "both"
	}
	if def.Direction == "" {
		def.Direction = "normal"
	}
	if def.Delay == "" {
		def.Delay = "0ms"
	}
	return def
}

// Builtin returns the built-in animation names, ordered for docs.
func Builtin() []string {
	return builtinOrder
}

// Theme is unused today but keeps the registry signature symmetric with the
// other registries; custom animations may later read theme tokens.
var _ = themes.Default
