package components

import (
	"strings"
	"testing"

	"github.com/Ahaduzzamankhan/flexui/internal/themes"
)

func declValues(t *testing.T, r *Registry, class string) map[string]string {
	t.Helper()
	entry, ok := r.Resolve(class)
	if !ok {
		t.Fatalf("%s is not registered", class)
	}
	values := make(map[string]string)
	for _, block := range entry.Blocks {
		if block.Target != "" || block.AtRule != "" {
			continue
		}
		for _, decl := range block.Decls {
			values[decl.Prop] = decl.Value
		}
	}
	return values
}

func TestCustomComponentRegistersClasses(t *testing.T) {
	registry, err := NewWithCustom(themes.Default(), []CustomComponent{{
		Name: "alert",
		Base: map[string]string{
			"display":          "flex",
			"gap":              "0.75rem",
			"padding":          "1rem",
			"border-radius":    "0.75rem",
			"background-color": "var(--flexui-surface-raised)",
		},
		Variants: map[string]map[string]string{
			"success": {"color": "var(--flexui-success)"},
		},
		Parts: map[string]map[string]string{
			"title": {"font-weight": "600"},
		},
		States: map[string]map[string]string{
			":hover": {"border-color": "var(--flexui-accent)"},
		},
	}})
	if err != nil {
		t.Fatalf("register: %v", err)
	}

	base := declValues(t, registry, "fui-alert")
	if base["display"] != "flex" || base["padding"] != "1rem" {
		t.Fatalf("base declarations = %v", base)
	}

	variant := declValues(t, registry, "fui-alert--success")
	if variant["color"] != "var(--flexui-success)" {
		t.Fatalf("variant declarations = %v", variant)
	}

	part := declValues(t, registry, "fui-alert-title")
	if part["font-weight"] != "600" {
		t.Fatalf("part declarations = %v", part)
	}

	entry, ok := registry.Resolve("fui-alert")
	if !ok {
		t.Fatal("base class missing")
	}
	foundState := false
	for _, block := range entry.Blocks {
		if block.Target == "&:hover" {
			foundState = true
		}
	}
	if !foundState {
		t.Fatalf("hover state block missing from %+v", entry.Blocks)
	}
}

func TestCustomComponentIsDeterministic(t *testing.T) {
	customs := []CustomComponent{{
		Name: "panel",
		Base: map[string]string{"z-index": "1", "align-items": "center", "display": "grid"},
	}}
	first, err := NewWithCustom(themes.Default(), customs)
	if err != nil {
		t.Fatal(err)
	}
	second, err := NewWithCustom(themes.Default(), customs)
	if err != nil {
		t.Fatal(err)
	}

	firstEntry, _ := first.Resolve("fui-panel")
	secondEntry, _ := second.Resolve("fui-panel")
	var firstProps, secondProps []string
	for _, decl := range firstEntry.Blocks[0].Decls {
		firstProps = append(firstProps, decl.Prop)
	}
	for _, decl := range secondEntry.Blocks[0].Decls {
		secondProps = append(secondProps, decl.Prop)
	}
	if strings.Join(firstProps, ",") != strings.Join(secondProps, ",") {
		t.Fatalf("declaration order differs: %v vs %v", firstProps, secondProps)
	}
	if strings.Join(firstProps, ",") != "align-items,display,z-index" {
		t.Fatalf("declarations are not sorted: %v", firstProps)
	}
}

func TestCustomComponentExtendsCustom(t *testing.T) {
	registry, err := NewWithCustom(themes.Default(), []CustomComponent{
		{
			Name: "surface",
			Base: map[string]string{"padding": "1rem", "border-radius": "1rem"},
			Parts: map[string]map[string]string{
				"title": {"font-size": "1rem"},
			},
		},
		{
			Name:    "callout",
			Extends: "surface",
			Base:    map[string]string{"padding": "1.5rem", "border": "1px solid red"},
		},
	})
	if err != nil {
		t.Fatalf("register: %v", err)
	}

	base := declValues(t, registry, "fui-callout")
	if base["padding"] != "1.5rem" {
		t.Fatalf("child should override the parent padding, got %v", base["padding"])
	}
	if base["border-radius"] != "1rem" || base["border"] == "" {
		t.Fatalf("inherited declarations missing: %v", base)
	}
	// Parts are inherited too.
	if part := declValues(t, registry, "fui-callout-title"); part["font-size"] != "1rem" {
		t.Fatalf("inherited part = %v", part)
	}
}

func TestCustomComponentExtendsBuiltin(t *testing.T) {
	registry, err := NewWithCustom(themes.Default(), []CustomComponent{{
		Name:    "banner",
		Extends: "card",
		Base:    map[string]string{"border-style": "dashed"},
	}})
	if err != nil {
		t.Fatalf("register: %v", err)
	}

	base := declValues(t, registry, "fui-banner")
	if base["border-style"] != "dashed" {
		t.Fatalf("own declaration missing: %v", base)
	}
	if base["display"] != "flex" || base["border-radius"] == "" {
		t.Fatalf("built-in declarations were not inherited: %v", base)
	}
}

func TestCustomComponentOverridesBuiltin(t *testing.T) {
	registry, err := NewWithCustom(themes.Default(), []CustomComponent{{
		Name: "btn",
		Base: map[string]string{"padding": "0.25rem"},
	}})
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	base := declValues(t, registry, "fui-btn")
	if len(base) != 1 || base["padding"] != "0.25rem" {
		t.Fatalf("custom btn should replace the built-in, got %v", base)
	}
}

func TestCustomComponentErrors(t *testing.T) {
	cases := map[string][]CustomComponent{
		"unknown parent": {{Name: "a", Extends: "missing", Base: map[string]string{"display": "block"}}},
		"cycle": {
			{Name: "a", Extends: "b", Base: map[string]string{"display": "block"}},
			{Name: "b", Extends: "a", Base: map[string]string{"display": "flex"}},
		},
		"self cycle": {{Name: "a", Extends: "a", Base: map[string]string{"display": "block"}}},
		"duplicate":  {{Name: "a", Base: map[string]string{"display": "block"}}, {Name: "a", Base: map[string]string{"display": "flex"}}},
		"empty name": {{Base: map[string]string{"display": "block"}}},
	}
	for name, customs := range cases {
		if _, err := NewWithCustom(themes.Default(), customs); err == nil {
			t.Fatalf("%s should have failed", name)
		}
	}
}

func TestCustomComponentNamesStaySorted(t *testing.T) {
	registry, err := NewWithCustom(themes.Default(), []CustomComponent{
		{Name: "zebra", Base: map[string]string{"display": "block"}},
		{Name: "alpha", Base: map[string]string{"display": "flex"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	names := registry.Names()
	index := func(target string) int {
		for i, name := range names {
			if name == target {
				return i
			}
		}
		return -1
	}
	if index("fui-alpha") < 0 || index("fui-zebra") < 0 {
		t.Fatal("custom classes missing from Names()")
	}
	if index("fui-alpha") > index("fui-zebra") {
		t.Fatalf("custom names are out of order: %d vs %d", index("fui-alpha"), index("fui-zebra"))
	}
	baseline := New(themes.Default()).Len()
	if registry.Len() != baseline+2 {
		t.Fatalf("registry holds %d classes, want %d built-in + 2 custom", registry.Len(), baseline)
	}
}
