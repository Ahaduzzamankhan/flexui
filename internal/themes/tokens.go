package themes

import (
	"sort"
	"strconv"
)

// FontFace is a local font declaration emitted into the base layer.
type FontFace struct {
	Family  string `json:"family"`
	Src     string `json:"src"`
	Weight  string `json:"weight"`
	Style   string `json:"style"`
	Display string `json:"display"`
}

// Overrides replaces individual theme tokens. Every field is optional; only the
// entries present are merged over the defaults, so a configuration file can
// extend a scale without restating it.
type Overrides struct {
	Colors         map[string]string
	Spacing        map[string]string
	FontSize       map[string]FontSize
	FontFamily     map[string]string
	FontWeight     map[string]string
	LineHeight     map[string]string
	LetterSpacing  map[string]string
	Radius         map[string]string
	Shadow         map[string]string
	Opacity        map[string]string
	ZIndex         map[string]string
	Breakpoints    map[string]string
	Easing         map[string]string
	Duration       map[string]string
	AspectRatio    map[string]string
	MaxWidth       map[string]string
	ContainerSizes map[string]string
	SemanticLight  map[string]string
	SemanticDark   map[string]string
}

// Apply merges overrides into the theme. Maps are mutated in place; Default
// returns fresh maps, so this never leaks between builds.
func (t *Theme) Apply(o Overrides) {
	merge := func(dst map[string]string, src map[string]string) {
		for k, v := range src {
			dst[k] = v
		}
	}
	merge(t.Colors, o.Colors)
	merge(t.Spacing, o.Spacing)
	merge(t.FontFamily, o.FontFamily)
	merge(t.FontWeight, o.FontWeight)
	merge(t.LineHeight, o.LineHeight)
	merge(t.LetterSpacing, o.LetterSpacing)
	merge(t.Radius, o.Radius)
	merge(t.Shadow, o.Shadow)
	merge(t.Opacity, o.Opacity)
	merge(t.ZIndex, o.ZIndex)
	merge(t.Breakpoints, o.Breakpoints)
	merge(t.Easing, o.Easing)
	merge(t.Duration, o.Duration)
	merge(t.AspectRatio, o.AspectRatio)
	merge(t.MaxWidth, o.MaxWidth)
	merge(t.ContainerSizes, o.ContainerSizes)
	merge(t.SemanticLight, o.SemanticLight)
	merge(t.SemanticDark, o.SemanticDark)
	for k, v := range o.FontSize {
		t.FontSize[k] = v
	}
}

// CSSVar is one custom property declaration.
type CSSVar struct {
	Name  string
	Value string
}

// TokenVars returns the custom properties emitted for :root. Ordering is fully
// deterministic so the same theme always produces the same stylesheet.
func (t *Theme) TokenVars() []CSSVar {
	vars := make([]CSSVar, 0, 128)
	for _, role := range sortedKeys(t.SemanticLight) {
		vars = append(vars, CSSVar{Name: "--flexui-" + role, Value: t.SemanticLight[role]})
	}
	for _, key := range numericKeys(t.Spacing) {
		vars = append(vars, CSSVar{Name: "--flexui-space-" + key, Value: t.Spacing[key]})
	}
	for _, key := range numericKeys(t.Radius) {
		if key == "" {
			vars = append(vars, CSSVar{Name: "--flexui-radius", Value: t.Radius[key]})
			continue
		}
		vars = append(vars, CSSVar{Name: "--flexui-radius-" + key, Value: t.Radius[key]})
	}
	for _, key := range sortedKeys(t.FontFamily) {
		vars = append(vars, CSSVar{Name: "--flexui-font-" + key, Value: t.FontFamily[key]})
	}
	textKeys := make([]string, 0, len(t.FontSize))
	for key := range t.FontSize {
		textKeys = append(textKeys, key)
	}
	sort.Strings(textKeys)
	for _, key := range textKeys {
		vars = append(vars, CSSVar{Name: "--flexui-text-" + key, Value: t.FontSize[key].Size})
	}
	return vars
}

// DarkTokenVars returns the role tokens overridden by the dark theme.
func (t *Theme) DarkTokenVars() []CSSVar {
	vars := make([]CSSVar, 0, len(t.SemanticDark))
	for _, role := range sortedKeys(t.SemanticDark) {
		vars = append(vars, CSSVar{Name: "--flexui-" + role, Value: t.SemanticDark[role]})
	}
	return vars
}

// SemanticVar returns the custom property reference for a role token such as
// "surface" or "fg-muted".
func (t *Theme) SemanticVar(role string) (string, bool) {
	if _, ok := t.SemanticLight[role]; !ok {
		return "", false
	}
	return "var(--flexui-" + role + ")", true
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// numericKeys sorts scale keys the way a designer reads them: 0, 0.5, 1, 2,
// 10, 96, px. Non-numeric keys such as "px" or "full" come last, sorted.
func numericKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		a, aok := parseNumber(keys[i])
		b, bok := parseNumber(keys[j])
		switch {
		case aok && bok:
			if a != b {
				return a < b
			}
			return keys[i] < keys[j]
		case aok:
			return true
		case bok:
			return false
		default:
			return keys[i] < keys[j]
		}
	})
	return keys
}

func parseNumber(s string) (float64, bool) {
	if s == "" {
		return 0, false
	}
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	return n, true
}
