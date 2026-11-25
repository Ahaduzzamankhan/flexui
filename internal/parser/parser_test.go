package parser

import "testing"

func TestParseVariants(t *testing.T) {
	c, ok := Parse("md:hover:bg-red-500")
	if !ok {
		t.Fatal("expected candidate to parse")
	}
	if len(c.Variants) != 2 || c.Variants[0] != "md" || c.Variants[1] != "hover" {
		t.Fatalf("variants = %v", c.Variants)
	}
	if c.Utility != "bg-red-500" {
		t.Fatalf("utility = %q", c.Utility)
	}
	if c.Important || c.Negative || c.ArbitraryProperty {
		t.Fatalf("unexpected flags: %+v", c)
	}
}

func TestParseArbitraryValueKeepsColons(t *testing.T) {
	c, ok := Parse("bg-[url(http://example.com/a.png)]")
	if !ok {
		t.Fatal("expected candidate to parse")
	}
	if len(c.Variants) != 0 {
		t.Fatalf("claimed variants %v", c.Variants)
	}
	if c.Utility != "bg-[url(http://example.com/a.png)]" {
		t.Fatalf("utility = %q", c.Utility)
	}
}

func TestParseArbitrarySelector(t *testing.T) {
	c, ok := Parse("[&>*]:p-2")
	if !ok {
		t.Fatal("expected candidate to parse")
	}
	if len(c.Variants) != 1 || c.Variants[0] != "[&>*]" {
		t.Fatalf("variants = %v", c.Variants)
	}
	if c.Utility != "p-2" {
		t.Fatalf("utility = %q", c.Utility)
	}
}

func TestParseDataVariant(t *testing.T) {
	c, ok := Parse("data-[state=open]:opacity-100")
	if !ok {
		t.Fatal("expected candidate to parse")
	}
	if len(c.Variants) != 1 || c.Variants[0] != "data-[state=open]" {
		t.Fatalf("variants = %v", c.Variants)
	}
	if c.Utility != "opacity-100" {
		t.Fatalf("utility = %q", c.Utility)
	}
}

func TestParseNegative(t *testing.T) {
	c, ok := Parse("-mt-4")
	if !ok || !c.Negative || c.Utility != "mt-4" {
		t.Fatalf("got %+v ok=%v", c, ok)
	}
}

func TestParseImportant(t *testing.T) {
	for _, raw := range []string{"!p-4", "p-4!"} {
		c, ok := Parse(raw)
		if !ok || !c.Important || c.Utility != "p-4" {
			t.Fatalf("%s -> %+v ok=%v", raw, c, ok)
		}
	}
}

func TestParseArbitraryProperty(t *testing.T) {
	c, ok := Parse("[mask-type:luminance]")
	if !ok || !c.ArbitraryProperty {
		t.Fatalf("got %+v ok=%v", c, ok)
	}
	if c.Utility != "[mask-type:luminance]" {
		t.Fatalf("utility = %q", c.Utility)
	}
}

func TestParseRejectsBrokenCandidates(t *testing.T) {
	for _, raw := range []string{"", "p-4:", ":p-4", "p-[4px", "md::p-4"} {
		if _, ok := Parse(raw); ok {
			t.Fatalf("expected %q to be rejected", raw)
		}
	}
}

func TestSplitUtility(t *testing.T) {
	cases := []struct{ token, prefix, value string }{
		{"p-4", "p", "4"},
		{"bg-red-500", "bg", "red-500"},
		{"grid-cols-[250px_1fr]", "grid", "cols-[250px_1fr]"},
		{"flex", "flex", ""},
	}
	for _, tc := range cases {
		prefix, value := SplitUtility(tc.token)
		if prefix != tc.prefix || value != tc.value {
			t.Fatalf("%s -> %q %q, want %q %q", tc.token, prefix, value, tc.prefix, tc.value)
		}
	}
}

func TestArbitraryHelpers(t *testing.T) {
	if v, ok := UnwrapArbitrary("[length:12px]"); !ok || v != "12px" {
		t.Fatalf("hint not stripped: %q %v", v, ok)
	}
	if v, ok := UnwrapArbitrary("[calc(100vh-80px)]"); !ok || v != "calc(100vh-80px)" {
		t.Fatalf("got %q %v", v, ok)
	}
	if got := NormalizeArbitrary("250px_1fr"); got != "250px 1fr" {
		t.Fatalf("normalize = %q", got)
	}
	if got := NormalizeArbitrary(`a\_b`); got != "a_b" {
		t.Fatalf("escaped underscore = %q", got)
	}
	if _, ok := UnwrapArbitrary("12px"); ok {
		t.Fatal("non-arbitrary value accepted")
	}
}
