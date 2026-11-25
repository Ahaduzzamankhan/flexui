package generator

import (
	"strings"
	"testing"

	"github.com/Ahaduzzamankhan/flexui/internal/ast"
	"github.com/Ahaduzzamankhan/flexui/internal/resolver"
	"github.com/Ahaduzzamankhan/flexui/internal/themes"
)

// resolved resolves classes with the default theme, failing the test on an
// unknown class so a typo in a test is reported.
func resolved(t *testing.T, classes ...string) []ast.Rule {
	t.Helper()
	r := resolver.NewWith(resolver.Options{Theme: themes.Default()})
	out := make([]ast.Rule, 0, len(classes))
	for _, class := range classes {
		rule, ok := r.Resolve(class)
		if !ok {
			t.Fatalf("class %q did not resolve", class)
		}
		out = append(out, rule)
	}
	return out
}

func TestThemeLayerEmitted(t *testing.T) {
	css := string(Generate(resolved(t, "bg-surface", "text-fg"), Options{
		Theme: themes.Default(),
		Base:  true,
	}))
	if !strings.Contains(css, "@layer flexui.theme, flexui.base, flexui.utilities;") {
		t.Fatalf("layer order statement wrong:\n%s", css)
	}
	if !strings.Contains(css, "@layer flexui.theme {") {
		t.Fatalf("theme layer missing:\n%s", css)
	}
	if !strings.Contains(css, "--flexui-surface:") {
		t.Fatalf("surface token missing:\n%s", css)
	}
	if !strings.Contains(css, ".bg-surface { background-color: var(--flexui-surface); }") {
		t.Fatalf("semantic background did not use the token:\n%s", css)
	}
}

func TestThemeLayerOmittedWithoutTheme(t *testing.T) {
	css := string(Generate(resolved(t, "p-4"), Options{Base: true}))
	if strings.Contains(css, "flexui.theme") {
		t.Fatalf("theme layer emitted without a theme:\n%s", css)
	}
}

func TestDarkModeStrategies(t *testing.T) {
	theme := themes.Default()
	cases := []struct {
		mode  string
		want  []string
		avoid string
	}{
		{"class", []string{".dark {"}, "@media (prefers-color-scheme: dark)"},
		{"media", []string{"@media (prefers-color-scheme: dark)", ":root {"}, ".dark {"},
		{"both", []string{".dark {", "@media (prefers-color-scheme: dark)", ":root:not(.light) {"}, ""},
	}
	for _, tc := range cases {
		css := string(Generate(resolved(t, "p-4"), Options{Theme: theme, DarkMode: tc.mode, Base: true}))
		for _, want := range tc.want {
			if !strings.Contains(css, want) {
				t.Fatalf("%s mode missing %q:\n%s", tc.mode, want, css)
			}
		}
		if tc.avoid != "" && strings.Contains(css, tc.avoid) {
			t.Fatalf("%s mode should not emit %q:\n%s", tc.mode, tc.avoid, css)
		}
	}
}

func TestDarkVariantFollowsTokenStrategy(t *testing.T) {
	class := string(Generate(resolved(t, "dark:bg-surface"), Options{Theme: themes.Default(), Base: true}))
	if !strings.Contains(class, ".dark .dark\\:bg-surface") {
		t.Fatalf("dark variant did not use the class strategy:\n%s", class)
	}
	media := string(Generate(resolved(t, "dark:bg-surface"), Options{
		Theme: themes.Default(), DarkMode: "media", Base: true,
	}))
	if !strings.Contains(media, "@media (prefers-color-scheme: dark)") ||
		!strings.Contains(media, ".dark\\:bg-surface") {
		t.Fatalf("dark variant did not use the media strategy:\n%s", media)
	}
}

func TestComponentsLayerIsIndependent(t *testing.T) {
	css := string(Generate(resolved(t, "fui-btn", "p-4"), Options{Theme: themes.Default(), Base: true}))
	order := "@layer flexui.theme, flexui.base, flexui.components, flexui.utilities;"
	if !strings.Contains(css, order) {
		t.Fatalf("layer order statement wrong:\n%s", css)
	}
	componentAt := strings.Index(css, "@layer flexui.components {")
	utilitiesAt := strings.Index(css, "@layer flexui.utilities {")
	if componentAt < 0 || utilitiesAt < 0 || componentAt > utilitiesAt {
		t.Fatalf("components layer must precede utilities:\n%s", css)
	}
	if !strings.Contains(css, ".fui-btn {") {
		t.Fatalf("button component missing:\n%s", css)
	}
}

func TestComponentMediaPartsNest(t *testing.T) {
	css := string(Generate(resolved(t, "fui-app", "fui-app-main"), Options{Theme: themes.Default(), Base: true}))
	if !strings.Contains(css, "@media (min-width: 1024px) {") {
		t.Fatalf("responsive component part missing:\n%s", css)
	}
	if !strings.Contains(css, "grid-template-areas: \"header header\" \"sidebar main\";") {
		t.Fatalf("app shell layout missing:\n%s", css)
	}
	// The at-rule must be closed inside the components layer.
	open := strings.Count(css, "{")
	close := strings.Count(css, "}")
	if open != close {
		t.Fatalf("unbalanced braces (%d open, %d closed):\n%s", open, close, css)
	}
}

func TestKeyframesEmittedOnce(t *testing.T) {
	rules := resolved(t, "animate-spin", "animate-spin", "animate-fade-up")
	css := string(Generate(rules, Options{Theme: themes.Default(), Base: true}))
	if count := strings.Count(css, "@keyframes fui-spin {"); count != 1 {
		t.Fatalf("spin keyframes emitted %d times:\n%s", count, css)
	}
	if !strings.Contains(css, "@keyframes fui-fade-up {") {
		t.Fatalf("fade-up keyframes missing:\n%s", css)
	}
	if !strings.Contains(css, "animation-name: fui-spin;") {
		t.Fatalf("animation name missing:\n%s", css)
	}
}

func TestAnimationControlsAreSeparateRules(t *testing.T) {
	css := string(Generate(resolved(t, "animate-fade-in", "animate-duration-500"), Options{
		Theme: themes.Default(), Base: true,
	}))
	if !strings.Contains(css, "animation-duration: 500ms;") {
		t.Fatalf("duration override missing:\n%s", css)
	}
	// The named animation must not emit keyframes for the control utility.
	if strings.Contains(css, "@keyframes fui-duration-500") {
		t.Fatalf("control utility emitted keyframes:\n%s", css)
	}
}

func TestReducedMotionOverrideOnlyWithAnimations(t *testing.T) {
	withAnimation := string(Generate(resolved(t, "animate-fade-in"), Options{
		Theme: themes.Default(), Base: true, ReducedMotion: true,
	}))
	if !strings.Contains(withAnimation, "@media (prefers-reduced-motion: reduce)") {
		t.Fatalf("reduced-motion block missing:\n%s", withAnimation)
	}
	if !strings.Contains(withAnimation, "--flexui-animation-duration: 1ms;") {
		t.Fatalf("reduced-motion override missing:\n%s", withAnimation)
	}
	withoutAnimation := string(Generate(resolved(t, "p-4"), Options{
		Theme: themes.Default(), Base: true, ReducedMotion: true,
	}))
	if strings.Contains(withoutAnimation, "prefers-reduced-motion") {
		t.Fatalf("reduced-motion block emitted without animations:\n%s", withoutAnimation)
	}
	disabled := string(Generate(resolved(t, "animate-fade-in"), Options{
		Theme: themes.Default(), Base: true, ReducedMotion: false,
	}))
	if strings.Contains(disabled, "prefers-reduced-motion") {
		t.Fatalf("reduced-motion block emitted while disabled:\n%s", disabled)
	}
}

func TestFontFacesEmittedInBaseLayer(t *testing.T) {
	css := string(Generate(resolved(t, "p-4"), Options{
		Base:  true,
		Theme: themes.Default(),
		Fonts: []themes.FontFace{{Family: "Inter", Src: "url(/fonts/inter.woff2) format(\"woff2\")", Weight: "400"}},
	}))
	if !strings.Contains(css, "@font-face {") {
		t.Fatalf("@font-face missing:\n%s", css)
	}
	if !strings.Contains(css, "font-family: \"Inter\";") {
		t.Fatalf("font family not quoted:\n%s", css)
	}
	if !strings.Contains(css, "font-display: swap;") {
		t.Fatalf("font display default missing:\n%s", css)
	}
	baseAt := strings.Index(css, "@layer flexui.base {")
	faceAt := strings.Index(css, "@font-face")
	if baseAt < 0 || faceAt < baseAt {
		t.Fatalf("@font-face must be inside the base layer:\n%s", css)
	}
}

func TestVariantAtRuleAndBlockAtRuleBothClose(t *testing.T) {
	css := string(Generate(resolved(t, "md:fui-app"), Options{Theme: themes.Default(), Base: true}))
	open := strings.Count(css, "{")
	close := strings.Count(css, "}")
	if open != close {
		t.Fatalf("unbalanced braces (%d open, %d closed):\n%s", open, close, css)
	}
	if !strings.Contains(css, ".md\\:fui-app") {
		t.Fatalf("variant selector missing:\n%s", css)
	}
}

func TestMinifiedLayersStayCompact(t *testing.T) {
	css := string(Generate(resolved(t, "fui-btn", "p-4", "animate-spin"), Options{
		Theme: themes.Default(), Base: true, Minify: true,
	}))
	if strings.Contains(css, "\n\n") {
		t.Fatalf("minified output has blank lines:\n%q", css)
	}
	if !strings.Contains(css, "@layer flexui.theme{") {
		t.Fatalf("minified layer syntax wrong:\n%q", css)
	}
	if strings.Contains(css, "  ") {
		t.Fatalf("minified output has indentation:\n%q", css)
	}
}
