// Package variants resolves FlexUI's variants (responsive, state, structural,
// group/peer, accessibility, theme, arbitrary selectors, supports, container
// queries, range breakpoints and user-defined variants) into selector templates
// and conditional group rules.
package variants

import (
	"sort"
	"strconv"
	"strings"

	"github.com/Ahaduzzamankhan/flexui/internal/ast"
	"github.com/Ahaduzzamankhan/flexui/internal/parser"
	"github.com/Ahaduzzamankhan/flexui/internal/themes"
)

// Ordering ranks. Responsive breakpoints sort after plain variants and before
// container queries, supports and print, which keeps the cascade predictable:
// a later breakpoint always wins over an earlier one.
const (
	RankBase          = 0
	RankLight         = 1
	RankPseudo        = 10
	RankArbitrary     = 12
	RankCustom        = 14
	RankTheme         = 15
	RankAncestor      = 20
	RankAccessibility = 30
	RankBreakpoint    = 100
	RankContainer     = 300
	RankSupports      = 400
	RankPrint         = 450
)

// Options configures variant resolution.
type Options struct {
	Theme *themes.Theme
	// DarkMode is "class", "media" or "both". Utilities follow "class" for
	// "both"; the token block covers the system preference separately.
	DarkMode string
	// Custom maps a variant name to a selector template using "&", or to an
	// at-rule when it starts with "@".
	Custom map[string]string
}

// Resolver resolves single variant tokens.
type Resolver struct {
	theme          *themes.Theme
	breakpoints    map[string]string
	breakpointRank map[string]int
	containers     map[string]string
	darkMode       string
	custom         map[string]string
	customNames    []string
}

// New builds a resolver with default options.
func New(t *themes.Theme) *Resolver {
	return NewWith(Options{Theme: t})
}

// NewWith builds a resolver with explicit options.
func NewWith(o Options) *Resolver {
	t := o.Theme
	if t == nil {
		t = themes.Default()
	}
	r := &Resolver{
		theme:          t,
		breakpoints:    t.Breakpoints,
		breakpointRank: make(map[string]int, len(t.Breakpoints)),
		containers:     t.ContainerSizes,
		darkMode:       o.DarkMode,
		custom:         o.Custom,
	}
	if r.darkMode == "" {
		r.darkMode = "class"
	}
	// BreakpointNames is ordered by width, so ranks follow the cascade order
	// (sm < md < lg < xl < 2xl) instead of alphabetical order.
	for i, name := range t.BreakpointNames() {
		r.breakpointRank[name] = RankBreakpoint + i*10
	}
	r.customNames = make([]string, 0, len(r.custom))
	for name := range r.custom {
		r.customNames = append(r.customNames, name)
	}
	sort.Slice(r.customNames, func(i, j int) bool {
		if len(r.customNames[i]) != len(r.customNames[j]) {
			return len(r.customNames[i]) > len(r.customNames[j])
		}
		return r.customNames[i] < r.customNames[j]
	})
	return r
}

// pseudoVariants map a variant to the selector suffix it appends.
var pseudoVariants = map[string]string{
	"hover":              ":hover",
	"focus":              ":focus",
	"focus-visible":      ":focus-visible",
	"focus-within":       ":focus-within",
	"active":             ":active",
	"visited":            ":visited",
	"target":             ":target",
	"disabled":           ":disabled",
	"enabled":            ":enabled",
	"checked":            ":checked",
	"required":           ":required",
	"optional":           ":optional",
	"invalid":            ":invalid",
	"valid":              ":valid",
	"user-invalid":       ":user-invalid",
	"placeholder-shown":  ":placeholder-shown",
	"indeterminate":      ":indeterminate",
	"read-only":          ":read-only",
	"read-write":         ":read-write",
	"autofill":           ":autofill",
	"first":              ":first-child",
	"last":               ":last-child",
	"odd":                ":nth-child(odd)",
	"even":               ":nth-child(even)",
	"only":               ":only-child",
	"empty":              ":empty",
	"first-of-type":      ":first-of-type",
	"last-of-type":       ":last-of-type",
	"only-of-type":       ":only-of-type",
	"open":               "[open]",
	"default":            ":default",
	"fullscreen":         ":fullscreen",
	"picture-in-picture": ":picture-in-picture",
	"inert":              ":is([inert], [inert] *)",
}

// groupVariants map a group-* suffix onto the state it observes.
var groupVariants = map[string]string{
	"hover":             ":hover",
	"focus":             ":focus",
	"focus-visible":     ":focus-visible",
	"focus-within":      ":focus-within",
	"active":            ":active",
	"visited":           ":visited",
	"disabled":          ":disabled",
	"checked":           ":checked",
	"open":              "[open]",
	"placeholder-shown": ":placeholder-shown",
	"invalid":           ":invalid",
	"valid":             ":valid",
	"required":          ":required",
	"empty":             ":empty",
	"target":            ":target",
}

// peerVariants map a peer-* suffix onto the sibling state it observes.
var peerVariants = map[string]string{
	"hover":             ":hover",
	"focus":             ":focus",
	"focus-visible":     ":focus-visible",
	"focus-within":      ":focus-within",
	"active":            ":active",
	"checked":           ":checked",
	"disabled":          ":disabled",
	"enabled":           ":enabled",
	"invalid":           ":invalid",
	"valid":             ":valid",
	"required":          ":required",
	"placeholder-shown": ":placeholder-shown",
	"open":              "[open]",
	"indeterminate":     ":indeterminate",
}

// mediaVariants map a variant onto a media query.
var mediaVariants = map[string]struct {
	query string
	rank  int
}{
	"motion-safe":    {"(prefers-reduced-motion: no-preference)", RankAccessibility},
	"motion-reduce":  {"(prefers-reduced-motion: reduce)", RankAccessibility + 1},
	"contrast-more":  {"(prefers-contrast: more)", RankAccessibility + 2},
	"contrast-less":  {"(prefers-contrast: less)", RankAccessibility + 3},
	"print":          {"print", RankPrint},
	"portrait":       {"(orientation: portrait)", RankAccessibility},
	"landscape":      {"(orientation: landscape)", RankAccessibility},
	"pointer-fine":   {"(pointer: fine)", RankAccessibility},
	"pointer-coarse": {"(pointer: coarse)", RankAccessibility},
}

// Resolve resolves one variant token.
func (v *Resolver) Resolve(name string) (ast.AppliedVariant, bool) {
	if name == "" {
		return ast.AppliedVariant{}, false
	}
	if suffix, ok := pseudoVariants[name]; ok {
		return ast.AppliedVariant{
			Name: name, Kind: ast.VariantSelector,
			Template: "&" + suffix, Rank: RankPseudo,
		}, true
	}
	if bp, ok := v.breakpoints[name]; ok {
		return ast.AppliedVariant{
			Name: name, Kind: ast.VariantAtRule,
			AtRule: "@media (min-width: " + bp + ")", Rank: v.rankOf(name),
		}, true
	}
	if strings.HasPrefix(name, "max-") {
		rest := strings.TrimPrefix(name, "max-")
		if bp, ok := v.breakpoints[rest]; ok {
			return ast.AppliedVariant{
				Name: name, Kind: ast.VariantAtRule,
				AtRule: "@media (max-width: " + below(bp) + ")", Rank: v.rankOf(rest) + 5,
			}, true
		}
		if size, ok := bracketValue(rest); ok {
			return ast.AppliedVariant{
				Name: name, Kind: ast.VariantAtRule,
				AtRule: "@media (max-width: " + size + ")", Rank: RankBreakpoint + 5,
			}, true
		}
	}
	if strings.HasPrefix(name, "min-[") {
		if size, ok := bracketValue(strings.TrimPrefix(name, "min-")); ok {
			return ast.AppliedVariant{
				Name: name, Kind: ast.VariantAtRule,
				AtRule: "@media (min-width: " + size + ")", Rank: RankBreakpoint + 1,
			}, true
		}
	}
	if strings.HasPrefix(name, "raw-") {
		// raw-[selector] escapes FlexUI entirely and is rejected here so a typo
		// does not silently emit an unsupported rule.
		if _, ok := bracketValue(strings.TrimPrefix(name, "raw-")); ok {
			return ast.AppliedVariant{}, false
		}
	}
	switch name {
	case "dark":
		if v.darkMode == "media" {
			return ast.AppliedVariant{
				Name: name, Kind: ast.VariantAtRule,
				AtRule: "@media (prefers-color-scheme: dark)", Rank: RankTheme,
			}, true
		}
		return ast.AppliedVariant{
			Name: name, Kind: ast.VariantAncestor,
			Template: ".dark &", Rank: RankTheme,
		}, true
	case "light":
		return ast.AppliedVariant{
			Name: name, Kind: ast.VariantAncestor,
			Template: "&", Rank: RankLight,
		}, true
	}
	if mv, ok := mediaVariants[name]; ok {
		return ast.AppliedVariant{
			Name: name, Kind: ast.VariantAtRule,
			AtRule: "@media " + mv.query, Rank: mv.rank,
		}, true
	}
	if strings.HasPrefix(name, "group-") {
		base, groupName := splitNamed(name[len("group-"):])
		suffix, ok := groupVariants[base]
		if !ok {
			return ast.AppliedVariant{}, false
		}
		selector := ".group"
		if groupName != "" {
			selector = ".group\\/" + groupName
		}
		return ast.AppliedVariant{
			Name: name, Kind: ast.VariantAncestor,
			Template: selector + suffix + " &", Rank: RankAncestor,
		}, true
	}
	if strings.HasPrefix(name, "peer-") {
		base, peerName := splitNamed(name[len("peer-"):])
		suffix, ok := peerVariants[base]
		if !ok {
			return ast.AppliedVariant{}, false
		}
		selector := ".peer"
		if peerName != "" {
			selector = ".peer\\/" + peerName
		}
		return ast.AppliedVariant{
			Name: name, Kind: ast.VariantAncestor,
			Template: selector + suffix + " ~ &", Rank: RankAncestor,
		}, true
	}
	if inner, ok := bracketed(name, "data-"); ok {
		return ast.AppliedVariant{
			Name: name, Kind: ast.VariantSelector,
			Template: "&[data-" + parser.NormalizeArbitrary(inner) + "]", Rank: RankPseudo,
		}, true
	}
	if inner, ok := bracketed(name, "aria-"); ok {
		return ast.AppliedVariant{
			Name: name, Kind: ast.VariantSelector,
			Template: "&[aria-" + parser.NormalizeArbitrary(inner) + "]", Rank: RankPseudo,
		}, true
	}
	if inner, ok := bracketed(name, "has-"); ok {
		return ast.AppliedVariant{
			Name: name, Kind: ast.VariantSelector,
			Template: "&:has(" + parser.NormalizeArbitrary(inner) + ")", Rank: RankPseudo,
		}, true
	}
	if inner, ok := bracketed(name, "supports-"); ok {
		return ast.AppliedVariant{
			Name: name, Kind: ast.VariantAtRule,
			AtRule: "@supports (" + parser.NormalizeArbitrary(inner) + ")", Rank: RankSupports,
		}, true
	}
	if inner, ok := bracketed(name, "not-"); ok {
		if suffix, found := pseudoVariants[inner]; found && suffix[0] == ':' {
			return ast.AppliedVariant{
				Name: name, Kind: ast.VariantSelector,
				Template: "&:not(" + suffix + ")", Rank: RankPseudo,
			}, true
		}
	}
	if strings.HasPrefix(name, "@[") {
		inner := name[2:]
		if strings.HasSuffix(inner, "]") {
			size := parser.NormalizeArbitrary(inner[:len(inner)-1])
			return ast.AppliedVariant{
				Name: name, Kind: ast.VariantAtRule,
				AtRule: "@container (min-width: " + size + ")", Rank: RankContainer,
			}, true
		}
		return ast.AppliedVariant{}, false
	}
	if strings.HasPrefix(name, "@") {
		preset := name[1:]
		if size, ok := v.containers[preset]; ok {
			return ast.AppliedVariant{
				Name: name, Kind: ast.VariantAtRule,
				AtRule: "@container (min-width: " + size + ")", Rank: RankContainer,
			}, true
		}
	}
	if v.custom != nil {
		for _, candidate := range v.customNames {
			if candidate != name {
				continue
			}
			template := strings.TrimSpace(v.custom[candidate])
			if template == "" {
				return ast.AppliedVariant{}, false
			}
			if strings.HasPrefix(template, "@") {
				return ast.AppliedVariant{
					Name: name, Kind: ast.VariantAtRule,
					AtRule: template, Rank: RankCustom,
				}, true
			}
			if !strings.Contains(template, "&") {
				template += " &"
			}
			return ast.AppliedVariant{
				Name: name, Kind: ast.VariantAncestor,
				Template: template, Rank: RankCustom,
			}, true
		}
	}
	if strings.HasPrefix(name, "[") && strings.HasSuffix(name, "]") {
		inner := parser.NormalizeArbitrary(name[1 : len(name)-1])
		if inner == "" {
			return ast.AppliedVariant{}, false
		}
		template := inner
		if !strings.Contains(inner, "&") {
			template = inner + " &"
		}
		return ast.AppliedVariant{
			Name: name, Kind: ast.VariantAncestor,
			Template: template, Rank: RankArbitrary,
		}, true
	}
	return ast.AppliedVariant{}, false
}

// splitNamed splits "hover/card" into "hover" and "card".
func splitNamed(s string) (string, string) {
	if i := strings.IndexByte(s, '/'); i >= 0 {
		return s[:i], s[i+1:]
	}
	return s, ""
}

// bracketed returns the inner text of "prefix-[...]".
func bracketed(name, prefix string) (string, bool) {
	if !strings.HasPrefix(name, prefix) || len(name) < len(prefix)+2 {
		return "", false
	}
	rest := name[len(prefix):]
	if rest[0] != '[' || rest[len(rest)-1] != ']' {
		return "", false
	}
	return rest[1 : len(rest)-1], true
}

// bracketValue returns the normalized inner text of "[...]".
func bracketValue(name string) (string, bool) {
	if len(name) < 3 || name[0] != '[' || name[len(name)-1] != ']' {
		return "", false
	}
	inner := parser.NormalizeArbitrary(name[1 : len(name)-1])
	if inner == "" {
		return "", false
	}
	return inner, true
}

// below turns a breakpoint width into the largest width that is still smaller,
// so max-* ranges never overlap the min-* breakpoint itself.
func below(px string) string {
	value := strings.TrimSuffix(strings.TrimSuffix(px, "px"), "rem")
	n, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return px
	}
	if strings.HasSuffix(px, "rem") {
		return strconv.FormatFloat(n-0.00125, 'f', -1, 64) + "rem"
	}
	if n <= 0 {
		return px
	}
	return strconv.FormatFloat(n-0.02, 'f', -1, 64) + "px"
}

func (v *Resolver) rankOf(name string) int {
	if rank, ok := v.breakpointRank[name]; ok {
		return rank
	}
	return RankBreakpoint
}

// IsBreakpoint reports whether a token is a theme breakpoint.
func (v *Resolver) IsBreakpoint(name string) bool {
	_, ok := v.breakpoints[name]
	return ok
}

// RankOf returns the sort rank for a breakpoint name.
func (v *Resolver) RankOf(name string) int { return v.rankOf(name) }

// ContainerSize returns a preset container size.
func (v *Resolver) ContainerSize(name string) (string, bool) {
	s, ok := v.containers[name]
	return s, ok
}

// CustomNames returns the configured custom variant names in longest-first
// order, for diagnostics.
func (v *Resolver) CustomNames() []string {
	out := make([]string, len(v.customNames))
	copy(out, v.customNames)
	return out
}
