package utilities

import (
	"testing"

	"github.com/Ahaduzzamankhan/flexui/internal/themes"
)

func TestUtilityValues(t *testing.T) {
	reg := New(themes.Default())
	cases := []struct {
		token string
		neg   bool
		prop  string
		value string
	}{
		{"p-4", false, "padding", "1rem"},
		{"px-2", false, "padding-inline", "0.5rem"},
		{"py-0.5", false, "padding-block", "0.125rem"},
		{"pt-[3px]", false, "padding-top", "3px"},
		{"m-4", false, "margin", "1rem"},
		{"mt-4", true, "margin-top", "-1rem"},
		{"mx-auto", false, "margin-inline", "auto"},
		{"gap-6", false, "gap", "1.5rem"},
		{"gap-x-2", false, "column-gap", "0.5rem"},
		{"flex", false, "display", "flex"},
		{"grid", false, "display", "grid"},
		{"hidden", false, "display", "none"},
		{"inline-block", false, "display", "inline-block"},
		{"relative", false, "position", "relative"},
		{"items-center", false, "align-items", "center"},
		{"justify-between", false, "justify-content", "space-between"},
		{"flex-col", false, "flex-direction", "column"},
		{"grid-cols-3", false, "grid-template-columns", "repeat(3, minmax(0, 1fr))"},
		{"grid-cols-[250px_1fr]", false, "grid-template-columns", "250px 1fr"},
		{"col-span-2", false, "grid-column", "span 2 / span 2"},
		{"w-full", false, "width", "100%"},
		{"w-1/2", false, "width", "50%"},
		{"w-screen", false, "width", "100vw"},
		{"h-[calc(100vh-80px)]", false, "height", "calc(100vh-80px)"},
		{"max-w-prose", false, "max-width", "65ch"},
		{"aspect-video", false, "aspect-ratio", "16 / 9"},
		{"inset-0", false, "inset", "0"},
		{"top-1/2", false, "top", "50%"},
		{"left-4", true, "left", "-1rem"},
		{"z-10", false, "z-index", "10"},
		{"bg-[#121212]", false, "background-color", "#121212"},
		{"bg-red-500/50", false, "background-color", "color-mix(in srgb, #ef4444 50%, transparent)"},
		{"text-white", false, "color", "#ffffff"},
		{"text-sm", false, "font-size", "0.875rem"},
		{"text-[clamp(2rem,5vw,5rem)]", false, "font-size", "clamp(2rem,5vw,5rem)"},
		{"font-bold", false, "font-weight", "700"},
		{"leading-none", false, "line-height", "1"},
		{"tracking-tight", false, "letter-spacing", "-0.025em"},
		{"border", false, "border-width", "1px"},
		{"border-2", false, "border-width", "2px"},
		{"border-red-500", false, "border-color", "#ef4444"},
		{"border-t-red-500", false, "border-top-color", "#ef4444"},
		{"border-t", false, "border-top-width", "1px"},
		{"divide-y", false, "border-block-start-width", "1px"},
		{"divide-x", false, "border-inline-start-width", "1px"},
		{"rounded-lg", false, "border-radius", "0.75rem"},
		{"rounded-t-xl", false, "border-top-left-radius", "1rem"},
		{"opacity-50", false, "opacity", "0.5"},
		{"scale-105", false, "--flexui-sx", "1.05"},
		{"rotate-45", false, "--flexui-rotate", "45deg"},
		{"blur-sm", false, "filter", "blur(4px)"},
		{"duration-300", false, "transition-duration", "300ms"},
		{"ease-smooth", false, "transition-timing-function", "cubic-bezier(0.22, 1, 0.36, 1)"},
		{"cursor-pointer", false, "cursor", "pointer"},
		{"overflow-hidden", false, "overflow", "hidden"},
		{"object-cover", false, "object-fit", "cover"},
		{"uppercase", false, "text-transform", "uppercase"},
		{"sr-only", false, "position", "absolute"},
	}
	for _, tc := range cases {
		blocks, _, _, ok := reg.Resolve(tc.token, tc.neg)
		if !ok {
			t.Fatalf("%q did not resolve", tc.token)
		}
		if len(blocks) == 0 || len(blocks[0].Decls) == 0 {
			t.Fatalf("%q produced no declarations", tc.token)
		}
		got := blocks[0].Decls[0]
		if got.Prop != tc.prop || got.Value != tc.value {
			t.Fatalf("%q -> %s: %s, want %s: %s", tc.token, got.Prop, got.Value, tc.prop, tc.value)
		}
	}
}

func TestSpaceBetweenUsesChildSelector(t *testing.T) {
	reg := New(themes.Default())
	blocks, _, _, ok := reg.Resolve("space-x-4", false)
	if !ok {
		t.Fatal("space-x-4 did not resolve")
	}
	if blocks[0].Target != "& > * + *" {
		t.Fatalf("target = %q", blocks[0].Target)
	}
	if blocks[0].Decls[0].Prop != "margin-inline-start" || blocks[0].Decls[0].Value != "1rem" {
		t.Fatalf("decls = %+v", blocks[0].Decls)
	}
}

func TestGradientUtilities(t *testing.T) {
	reg := New(themes.Default())
	blocks, _, _, ok := reg.Resolve("bg-gradient-to-r", false)
	if !ok {
		t.Fatal("gradient direction did not resolve")
	}
	if blocks[0].Decls[0].Prop != "background-image" {
		t.Fatalf("decls = %+v", blocks[0].Decls)
	}
	stop, _, _, ok := reg.Resolve("from-indigo-500", false)
	if !ok || stop[0].Decls[0].Prop != "--flexui-gradient-from" {
		t.Fatalf("gradient stop = %+v", stop)
	}
}

func TestRingAndShadowCompose(t *testing.T) {
	reg := New(themes.Default())
	ring, _, _, ok := reg.Resolve("ring-2", false)
	if !ok {
		t.Fatal("ring did not resolve")
	}
	var shared bool
	for _, decl := range ring[0].Decls {
		if decl.Prop == "box-shadow" {
			shared = true
		}
	}
	if !shared {
		t.Fatalf("ring should emit a composed box-shadow: %+v", ring[0].Decls)
	}
	shadow, _, _, ok := reg.Resolve("shadow-md", false)
	if !ok {
		t.Fatal("shadow did not resolve")
	}
	if shadow[0].Decls[0].Prop != "--flexui-shadow" {
		t.Fatalf("shadow decls = %+v", shadow[0].Decls)
	}
}

func TestUnknownUtilities(t *testing.T) {
	reg := New(themes.Default())
	for _, token := range []string{"bogus-1", "p-", "text-notacolor", "w-", "grid-cols-0"} {
		if blocks, _, _, ok := reg.Resolve(token, false); ok {
			t.Fatalf("%q unexpectedly resolved to %+v", token, blocks)
		}
	}
}

func TestArbitraryPropertySupport(t *testing.T) {
	reg := New(themes.Default())
	blocks, _, _, ok := reg.Resolve("grid-cols-[repeat(auto-fill,_minmax(200px,1fr))]", false)
	if !ok {
		t.Fatal("arbitrary grid template did not resolve")
	}
	want := "repeat(auto-fill, minmax(200px,1fr))"
	if blocks[0].Decls[0].Value != want {
		t.Fatalf("value = %q, want %q", blocks[0].Decls[0].Value, want)
	}
}

func TestNegativeArbitrary(t *testing.T) {
	reg := New(themes.Default())
	blocks, _, _, ok := reg.Resolve("translate-x-[10px]", true)
	if !ok {
		t.Fatal("negative translate did not resolve")
	}
	if blocks[0].Decls[0].Value != "-10px" {
		t.Fatalf("value = %q", blocks[0].Decls[0].Value)
	}
}
