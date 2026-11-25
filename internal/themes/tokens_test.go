package themes

import (
	"reflect"
	"testing"
)

func TestApplyOverridesTokens(t *testing.T) {
	theme := Default()
	theme.Apply(Overrides{
		Spacing:       map[string]string{"18": "4.5rem", "4": "1.125rem"},
		Colors:        map[string]string{"brand": "#123456"},
		FontSize:      map[string]FontSize{"hero": {Size: "4rem", LineHeight: "1"}},
		SemanticLight: map[string]string{"accent": "#ff0000"},
		SemanticDark:  map[string]string{"accent": "#00ff00"},
		Radius:        map[string]string{"card": "1.5rem"},
	})
	if got := theme.Spacing["18"]; got != "4.5rem" {
		t.Fatalf("spacing 18 = %q", got)
	}
	if got := theme.Spacing["4"]; got != "1.125rem" {
		t.Fatalf("spacing 4 override = %q", got)
	}
	if got := theme.Colors["brand"]; got != "#123456" {
		t.Fatalf("brand colour = %q", got)
	}
	if got := theme.FontSize["hero"].Size; got != "4rem" {
		t.Fatalf("hero size = %q", got)
	}
	if got := theme.SemanticLight["accent"]; got != "#ff0000" {
		t.Fatalf("light accent = %q", got)
	}
	if got := theme.SemanticDark["accent"]; got != "#00ff00" {
		t.Fatalf("dark accent = %q", got)
	}
	if _, ok := theme.Color("brand"); !ok {
		t.Fatalf("brand colour did not become resolvable")
	}
}

func TestApplyLeavesUnrelatedTokensAlone(t *testing.T) {
	theme := Default()
	before := theme.Spacing["16"]
	theme.Apply(Overrides{Spacing: map[string]string{"18": "4.5rem"}})
	if theme.Spacing["16"] != before {
		t.Fatalf("unrelated spacing changed: %q -> %q", before, theme.Spacing["16"])
	}
}

func TestTwoThemesDoNotShareState(t *testing.T) {
	first := Default()
	first.Apply(Overrides{SemanticLight: map[string]string{"accent": "#ff0000"}})
	second := Default()
	if second.SemanticLight["accent"] == "#ff0000" {
		t.Fatalf("theme overrides leaked into a new default theme")
	}
}

func TestTokenVarsAreDeterministic(t *testing.T) {
	first := Default().TokenVars()
	second := Default().TokenVars()
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("token vars are not deterministic")
	}
	if len(first) == 0 {
		t.Fatalf("no token vars emitted")
	}
	seen := make(map[string]bool, len(first))
	for _, v := range first {
		if seen[v.Name] {
			t.Fatalf("duplicate token %s", v.Name)
		}
		seen[v.Name] = true
		if v.Value == "" {
			t.Fatalf("token %s has an empty value", v.Name)
		}
	}
}

func TestTokenVarsIncludeSemanticAndScales(t *testing.T) {
	seen := map[string]bool{}
	for _, v := range Default().TokenVars() {
		seen[v.Name] = true
	}
	for _, want := range []string{
		"--flexui-surface", "--flexui-fg", "--flexui-fg-muted", "--flexui-line",
		"--flexui-accent", "--flexui-space-4", "--flexui-radius-lg",
		"--flexui-font-sans", "--flexui-text-base",
	} {
		if !seen[want] {
			t.Fatalf("%s missing from :root tokens", want)
		}
	}
}

func TestSemanticVarOnlyForKnownRoles(t *testing.T) {
	theme := Default()
	value, ok := theme.SemanticVar("surface")
	if !ok || value != "var(--flexui-surface)" {
		t.Fatalf("surface = %q, %v", value, ok)
	}
	if _, ok := theme.SemanticVar("not-a-role"); ok {
		t.Fatalf("unknown role resolved")
	}
}

func TestDarkTokenVarsOverrideRoles(t *testing.T) {
	theme := Default()
	dark := theme.DarkTokenVars()
	if len(dark) == 0 {
		t.Fatalf("no dark tokens")
	}
	for _, v := range dark {
		if _, ok := theme.SemanticLight[v.Name[len("--flexui-"):]]; !ok {
			t.Fatalf("dark token %s has no light counterpart", v.Name)
		}
	}
}
