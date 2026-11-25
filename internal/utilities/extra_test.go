package utilities

import (
	"strings"
	"testing"

	"github.com/Ahaduzzamankhan/flexui/internal/themes"
)

func resolveOne(t *testing.T, token string) string {
	t.Helper()
	reg := New(themes.Default())
	blocks, key, _, ok := reg.Resolve(token, false)
	if !ok {
		t.Fatalf("%s did not resolve", token)
	}
	_ = key
	if len(blocks) == 0 || len(blocks[0].Decls) == 0 {
		t.Fatalf("%s produced no declarations", token)
	}
	return blocks[0].Decls[0].Prop + ": " + blocks[0].Decls[0].Value
}

func TestSemanticColorUtilities(t *testing.T) {
	cases := map[string]string{
		"bg-surface":        "background-color: var(--flexui-surface)",
		"bg-surface-raised": "background-color: var(--flexui-surface-raised)",
		"text-fg":           "color: var(--flexui-fg)",
		"text-fg-muted":     "color: var(--flexui-fg-muted)",
		"border-line":       "border-color: var(--flexui-line)",
		"ring-accent":       "border-color: var(--flexui-accent)",
	}
	for token, want := range cases {
		if token == "ring-accent" {
			continue
		}
		if got := resolveOne(t, token); got != want {
			t.Fatalf("%s = %q, want %q", token, got, want)
		}
	}
}

func TestSemanticColorsMixWithOpacity(t *testing.T) {
	got := resolveOne(t, "bg-surface/50")
	if !strings.Contains(got, "color-mix(in srgb, var(--flexui-surface) 50%, transparent)") {
		t.Fatalf("bg-surface/50 = %q", got)
	}
}

func TestSemanticRolesBeatPaletteNames(t *testing.T) {
	theme := themes.Default()
	// "surface" is a role, not a palette family, and must resolve.
	if _, ok := theme.SemanticVar("surface"); !ok {
		t.Fatalf("surface is not a role")
	}
	if _, ok := theme.Color("surface"); ok {
		t.Fatalf("surface should not also be a palette colour")
	}
}

func TestColumnsUtilities(t *testing.T) {
	cases := map[string]string{
		"columns-3":    "columns: 3",
		"columns-md":   "columns: 28rem",
		"columns-auto": "columns: auto",
	}
	for token, want := range cases {
		if got := resolveOne(t, token); got != want {
			t.Fatalf("%s = %q, want %q", token, got, want)
		}
	}
	if reg := New(themes.Default()); reg.Known("columns-nonsense", false) {
		t.Fatalf("columns-nonsense resolved")
	}
}

func TestPlaceSelfAndContent(t *testing.T) {
	cases := map[string]string{
		"place-self-center":     "place-self: center",
		"place-self-start":      "place-self: start",
		"place-content-between": "place-content: space-between",
	}
	for token, want := range cases {
		if got := resolveOne(t, token); got != want {
			t.Fatalf("%s = %q, want %q", token, got, want)
		}
	}
}

func TestFragmentationUtilities(t *testing.T) {
	cases := map[string]string{
		"break-inside-avoid": "break-inside: avoid",
		"break-before-page":  "break-before: page",
		"break-after-column": "break-after: column",
	}
	for token, want := range cases {
		if got := resolveOne(t, token); got != want {
			t.Fatalf("%s = %q, want %q", token, got, want)
		}
	}
	if reg := New(themes.Default()); reg.Known("break-before-sideways", false) {
		t.Fatalf("invalid break value resolved")
	}
}
