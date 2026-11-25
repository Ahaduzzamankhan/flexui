package variants

import (
	"testing"

	"github.com/Ahaduzzamankhan/flexui/internal/ast"
	"github.com/Ahaduzzamankhan/flexui/internal/themes"
)

func TestMaxBreakpointsCoverBelowTheMin(t *testing.T) {
	r := New(themes.Default())
	cases := map[string]string{
		"max-sm":  "@media (max-width: 639.98px)",
		"max-md":  "@media (max-width: 767.98px)",
		"max-2xl": "@media (max-width: 1535.98px)",
	}
	for name, want := range cases {
		v, ok := r.Resolve(name)
		if !ok {
			t.Fatalf("%s did not resolve", name)
		}
		if v.Kind != ast.VariantAtRule || v.AtRule != want {
			t.Fatalf("%s = %q, want %q", name, v.AtRule, want)
		}
	}
}

func TestMaxBreakpointsSortAfterTheirBreakpoint(t *testing.T) {
	r := New(themes.Default())
	min, _ := r.Resolve("md")
	max, _ := r.Resolve("max-md")
	if max.Rank <= min.Rank {
		t.Fatalf("max-md rank %d should be after md rank %d", max.Rank, min.Rank)
	}
	lg, _ := r.Resolve("lg")
	if max.Rank >= lg.Rank {
		t.Fatalf("max-md must not sort after lg")
	}
}

func TestArbitraryRangeBreakpoints(t *testing.T) {
	r := New(themes.Default())
	min, ok := r.Resolve("min-[700px]")
	if !ok || min.AtRule != "@media (min-width: 700px)" {
		t.Fatalf("min-[700px] = %q, %v", min.AtRule, ok)
	}
	max, ok := r.Resolve("max-[40rem]")
	if !ok || max.AtRule != "@media (max-width: 40rem)" {
		t.Fatalf("max-[40rem] = %q, %v", max.AtRule, ok)
	}
	if _, ok := r.Resolve("max-nonsense"); ok {
		t.Fatalf("max-nonsense should not resolve")
	}
}

func TestCustomVariants(t *testing.T) {
	r := NewWith(Options{
		Theme: themes.Default(),
		Custom: map[string]string{
			"card":  ".card &",
			"touch": "@media (pointer: coarse)",
			"bare":  ".bare",
		},
	})
	card, ok := r.Resolve("card")
	if !ok || card.Kind != ast.VariantAncestor || card.Template != ".card &" {
		t.Fatalf("card variant = %+v, %v", card, ok)
	}
	touch, ok := r.Resolve("touch")
	if !ok || touch.Kind != ast.VariantAtRule || touch.AtRule != "@media (pointer: coarse)" {
		t.Fatalf("touch variant = %+v, %v", touch, ok)
	}
	bare, ok := r.Resolve("bare")
	if !ok || bare.Template != ".bare &" {
		t.Fatalf("bare variant = %+v, %v", bare, ok)
	}
	if names := r.CustomNames(); len(names) != 3 {
		t.Fatalf("custom names = %v", names)
	}
}

func TestDarkModeStrategies(t *testing.T) {
	class := NewWith(Options{Theme: themes.Default(), DarkMode: "class"})
	v, _ := class.Resolve("dark")
	if v.Kind != ast.VariantAncestor || v.Template != ".dark &" {
		t.Fatalf("class dark = %+v", v)
	}
	media := NewWith(Options{Theme: themes.Default(), DarkMode: "media"})
	v, _ = media.Resolve("dark")
	if v.Kind != ast.VariantAtRule || v.AtRule != "@media (prefers-color-scheme: dark)" {
		t.Fatalf("media dark = %+v", v)
	}
	both := NewWith(Options{Theme: themes.Default(), DarkMode: "both"})
	v, _ = both.Resolve("dark")
	if v.Kind != ast.VariantAncestor {
		t.Fatalf("both dark should use the class strategy for utilities: %+v", v)
	}
}

func TestMotionAndPrintVariants(t *testing.T) {
	r := New(themes.Default())
	cases := map[string]string{
		"motion-safe":   "@media (prefers-reduced-motion: no-preference)",
		"motion-reduce": "@media (prefers-reduced-motion: reduce)",
		"print":         "@media print",
		"portrait":      "@media (orientation: portrait)",
	}
	for name, want := range cases {
		v, ok := r.Resolve(name)
		if !ok || v.AtRule != want {
			t.Fatalf("%s = %q, %v", name, v.AtRule, ok)
		}
	}
}

func TestUnknownVariantsAreRejected(t *testing.T) {
	r := New(themes.Default())
	for _, name := range []string{"", "nope", "group-nope", "peer-nope", "group", "dark-nope", "[]"} {
		if _, ok := r.Resolve(name); ok {
			t.Fatalf("%q should not resolve", name)
		}
	}
}

func TestNamedGroupAndPeerVariants(t *testing.T) {
	r := New(themes.Default())
	group, ok := r.Resolve("group-hover/card")
	if !ok || group.Template != ".group\\/card:hover &" {
		t.Fatalf("named group = %+v, %v", group, ok)
	}
	peer, ok := r.Resolve("peer-checked/toggle")
	if !ok || peer.Template != ".peer\\/toggle:checked ~ &" {
		t.Fatalf("named peer = %+v, %v", peer, ok)
	}
}
