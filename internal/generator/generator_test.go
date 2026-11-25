package generator

import (
	"strings"
	"testing"

	"github.com/Ahaduzzamankhan/flexui/internal/ast"
	"github.com/Ahaduzzamankhan/flexui/internal/resolver"
	"github.com/Ahaduzzamankhan/flexui/internal/themes"
)

func rule(t *testing.T, candidates ...string) []ast.Rule {
	t.Helper()
	r := resolver.New(themes.Default())
	out := make([]ast.Rule, 0, len(candidates))
	for _, c := range candidates {
		got, ok := r.Resolve(c)
		if !ok {
			t.Fatalf("%q did not resolve", c)
		}
		out = append(out, got)
	}
	return out
}

func TestEscapeClass(t *testing.T) {
	cases := map[string]string{
		"p-4":          "p-4",
		"w-1/2":        `w-1\/2`,
		"bg-[#121212]": `bg-\[\#121212\]`,
		"hover:p-4":    `hover\:p-4`,
		"2xl:p-4":      `\32 xl\:p-4`,
		"text-[10px]":  `text-\[10px\]`,
		"[&>*]:p-4":    `\[\&\>\*\]\:p-4`,
	}
	for in, want := range cases {
		if got := EscapeClass(in); got != want {
			t.Fatalf("EscapeClass(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestGenerateRuleAndMedia(t *testing.T) {
	rules := rule(t, "p-4", "md:hover:bg-red-500")
	css := string(Generate(rules, Options{}))
	if !strings.Contains(css, ".p-4 { padding: 1rem; }") {
		t.Fatalf("missing padding rule:\n%s", css)
	}
	if !strings.Contains(css, "@media (min-width: 768px) {") {
		t.Fatalf("missing media query:\n%s", css)
	}
	if !strings.Contains(css, `.md\:hover\:bg-red-500:hover`) {
		t.Fatalf("missing escaped responsive selector:\n%s", css)
	}
	if !strings.Contains(css, "@layer flexui.utilities") {
		t.Fatalf("missing utilities layer:\n%s", css)
	}
}

func TestGenerateGroupAndArbitrarySelectors(t *testing.T) {
	css := string(Generate(rule(t, "group-hover:opacity-100", "[&>*]:p-2"), Options{}))
	if !strings.Contains(css, ".group:hover .group-hover\\:opacity-100") {
		t.Fatalf("group selector missing:\n%s", css)
	}
	if !strings.Contains(css, `\[\&\>\*\]\:p-2>*`) {
		t.Fatalf("arbitrary child selector missing:\n%s", css)
	}
}

func TestGenerateOrderIsDeterministic(t *testing.T) {
	a := Generate(rule(t, "p-4", "pt-2", "md:p-4", "hover:p-1", "bg-red-500"), Options{})
	b := Generate(rule(t, "bg-red-500", "hover:p-1", "md:p-4", "pt-2", "p-4"), Options{})
	if string(a) != string(b) {
		t.Fatalf("output depends on input order:\n%s\n---\n%s", a, b)
	}
}

func TestShorthandSortsBeforeLonghand(t *testing.T) {
	css := string(Generate(rule(t, "pt-2", "p-4"), Options{}))
	if strings.Index(css, ".p-4 ") > strings.Index(css, ".pt-2 ") {
		t.Fatalf("padding longhand came before shorthand:\n%s", css)
	}
}

func TestResponsiveRulesComeAfterBase(t *testing.T) {
	css := string(Generate(rule(t, "md:p-4", "p-2"), Options{}))
	if strings.Index(css, ".p-2 ") > strings.Index(css, "@media") {
		t.Fatalf("base rule emitted after media query:\n%s", css)
	}
}

func TestMinifiedOutput(t *testing.T) {
	css := string(Generate(rule(t, "p-4", "md:p-4"), Options{Minify: true}))
	if strings.Contains(css, "\n\n") {
		t.Fatalf("minified output has blank lines:\n%q", css)
	}
	if strings.Contains(css, "    ") {
		t.Fatalf("minified output has indentation:\n%q", css)
	}
	if !strings.Contains(css, ".p-4{padding:1rem;}") {
		t.Fatalf("compact rule missing:\n%q", css)
	}
}

func TestBaseLayerToggle(t *testing.T) {
	withBase := string(Generate(rule(t, "p-4"), Options{Base: true}))
	if !strings.Contains(withBase, "@layer flexui.base, flexui.utilities;") {
		t.Fatalf("layer order statement missing:\n%s", withBase)
	}
	without := string(Generate(rule(t, "p-4"), Options{Base: false}))
	if strings.Contains(without, "@layer flexui.base {") {
		t.Fatalf("base layer should be off:\n%s", without)
	}
	// The order statement declares the layers that are actually emitted, so a
	// build without preflight must not claim a base layer.
	if !strings.Contains(without, "@layer flexui.utilities;") {
		t.Fatalf("layer order statement missing:\n%s", without)
	}
	if strings.Contains(without, "flexui.base") {
		t.Fatalf("unused base layer was declared:\n%s", without)
	}
}

func TestPassthroughIsAppended(t *testing.T) {
	css := string(Generate(rule(t, "p-4"), Options{Passthrough: ".custom { color: red; }"}))
	if !strings.Contains(css, ".custom { color: red; }") {
		t.Fatalf("passthrough missing:\n%s", css)
	}
	if strings.Index(css, ".p-4") > strings.Index(css, ".custom") {
		t.Fatalf("passthrough should come last:\n%s", css)
	}
}
