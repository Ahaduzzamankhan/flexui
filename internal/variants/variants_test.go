package variants

import (
	"testing"

	"github.com/Ahaduzzamankhan/flexui/internal/ast"
	"github.com/Ahaduzzamankhan/flexui/internal/themes"
)

func TestResolveVariants(t *testing.T) {
	r := New(themes.Default())
	cases := []struct {
		name     string
		kind     ast.VariantKind
		template string
		atRule   string
	}{
		{name: "hover", kind: ast.VariantSelector, template: "&:hover"},
		{name: "focus-visible", kind: ast.VariantSelector, template: "&:focus-visible"},
		{name: "disabled", kind: ast.VariantSelector, template: "&:disabled"},
		{name: "first", kind: ast.VariantSelector, template: "&:first-child"},
		{name: "odd", kind: ast.VariantSelector, template: "&:nth-child(odd)"},
		{name: "md", kind: ast.VariantAtRule, atRule: "@media (min-width: 768px)"},
		{name: "2xl", kind: ast.VariantAtRule, atRule: "@media (min-width: 1536px)"},
		{name: "dark", kind: ast.VariantAncestor, template: ".dark &"},
		{name: "group-hover", kind: ast.VariantAncestor, template: ".group:hover &"},
		{name: "peer-checked", kind: ast.VariantAncestor, template: ".peer:checked ~ &"},
		{name: "data-[state=open]", kind: ast.VariantSelector, template: "&[data-state=open]"},
		{name: "aria-[expanded=true]", kind: ast.VariantSelector, template: "&[aria-expanded=true]"},
		{name: "has-[:checked]", kind: ast.VariantSelector, template: "&:has(:checked)"},
		{name: "supports-[display:grid]", kind: ast.VariantAtRule, atRule: "@supports (display:grid)"},
		{name: "motion-reduce", kind: ast.VariantAtRule, atRule: "@media (prefers-reduced-motion: reduce)"},
		{name: "contrast-more", kind: ast.VariantAtRule, atRule: "@media (prefers-contrast: more)"},
		{name: "[&>*]", kind: ast.VariantAncestor, template: "&>*"},
		{name: "@[40rem]", kind: ast.VariantAtRule, atRule: "@container (min-width: 40rem)"},
		{name: "@md", kind: ast.VariantAtRule, atRule: "@container (min-width: 28rem)"},
	}
	for _, tc := range cases {
		got, ok := r.Resolve(tc.name)
		if !ok {
			t.Fatalf("%q did not resolve", tc.name)
		}
		if got.Kind != tc.kind {
			t.Fatalf("%q kind = %v, want %v", tc.name, got.Kind, tc.kind)
		}
		if got.Template != tc.template {
			t.Fatalf("%q template = %q, want %q", tc.name, got.Template, tc.template)
		}
		if got.AtRule != tc.atRule {
			t.Fatalf("%q at-rule = %q, want %q", tc.name, got.AtRule, tc.atRule)
		}
	}
}

func TestUnknownVariant(t *testing.T) {
	r := New(themes.Default())
	for _, name := range []string{"", "bogus", "group-bogus", "peer-bogus", "supports-["} {
		if got, ok := r.Resolve(name); ok {
			t.Fatalf("%q unexpectedly resolved to %+v", name, got)
		}
	}
}

func TestBreakpointOrdering(t *testing.T) {
	r := New(themes.Default())
	sm, _ := r.Resolve("sm")
	md, _ := r.Resolve("md")
	xl, _ := r.Resolve("xl")
	if !(sm.Rank < md.Rank && md.Rank < xl.Rank) {
		t.Fatalf("breakpoint ranks out of order: %d %d %d", sm.Rank, md.Rank, xl.Rank)
	}
	hover, _ := r.Resolve("hover")
	if hover.Rank >= sm.Rank {
		t.Fatalf("state variant should sort before breakpoints")
	}
}

func TestPatternVariantsWithValuesStillMatchPrefixes(t *testing.T) {
	r := New(themes.Default())
	got, ok := r.Resolve("[.parent_&]")
	if !ok {
		t.Fatal("arbitrary ancestor variant did not resolve")
	}
	if got.Template != ".parent &" {
		t.Fatalf("template = %q", got.Template)
	}
}
