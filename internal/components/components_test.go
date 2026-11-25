package components

import (
	"strings"
	"testing"

	"github.com/Ahaduzzamankhan/flexui/internal/ast"
	"github.com/Ahaduzzamankhan/flexui/internal/themes"
)

func TestRequiredComponentsExist(t *testing.T) {
	reg := New(themes.Default())
	required := []string{
		// Foundation
		"fui-btn", "fui-btn--primary", "fui-badge", "fui-avatar", "fui-card",
		"fui-card-header", "fui-separator", "fui-skeleton", "fui-spinner",
		// Forms
		"fui-field", "fui-label", "fui-input", "fui-textarea", "fui-select",
		"fui-checkbox", "fui-radio", "fui-switch", "fui-slider", "fui-error",
		// Navigation
		"fui-navbar", "fui-sidebar", "fui-breadcrumb", "fui-tabs", "fui-tab",
		"fui-pagination", "fui-menu", "fui-menu-item", "fui-cmd", "fui-cmd-item",
		// Overlays
		"fui-dialog", "fui-dialog-header", "fui-drawer", "fui-popover",
		"fui-tooltip", "fui-toast-region", "fui-toast",
		// Data display
		"fui-table", "fui-th", "fui-td", "fui-calendar", "fui-calendar-day",
		"fui-timeline", "fui-tree", "fui-progress", "fui-progress-bar",
		"fui-stat", "fui-chart",
		// Layout presets
		"fui-app", "fui-app-sidebar", "fui-dashboard", "fui-docs",
		"fui-shop-grid", "fui-search", "fui-settings", "fui-steps", "fui-bento",
	}
	for _, name := range required {
		if _, ok := reg.Resolve(name); !ok {
			t.Fatalf("%s is not registered", name)
		}
	}
}

func TestEveryComponentIsPrefixedAndDeclaresSomething(t *testing.T) {
	reg := New(themes.Default())
	for _, name := range reg.Names() {
		if !strings.HasPrefix(name, "fui-") {
			t.Fatalf("%s is not in the fui- namespace", name)
		}
		result, _ := reg.Resolve(name)
		decls := 0
		for _, block := range result.Blocks {
			decls += len(block.Decls)
		}
		if decls == 0 {
			t.Fatalf("%s declares nothing", name)
		}
	}
}

func TestComponentsUseSemanticTokens(t *testing.T) {
	reg := New(themes.Default())
	button, _ := reg.Resolve("fui-btn--primary")
	found := false
	for _, block := range button.Blocks {
		for _, decl := range block.Decls {
			if strings.Contains(decl.Value, "var(--flexui-accent") {
				found = true
			}
		}
	}
	if !found {
		t.Fatalf("primary button does not use the accent token: %+v", button.Blocks)
	}
	if _, ok := themes.Default().SemanticVar("accent"); !ok {
		t.Fatalf("accent is not a semantic token")
	}
}

func TestInteractiveStatesAreStyled(t *testing.T) {
	reg := New(themes.Default())
	tab, _ := reg.Resolve("fui-tab")
	targets := make([]string, 0, len(tab.Blocks))
	for _, block := range tab.Blocks {
		targets = append(targets, block.Target)
	}
	joined := strings.Join(targets, " ")
	for _, want := range []string{"aria-selected", "hover", "focus-visible"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("tab is missing the %s state: %v", want, targets)
		}
	}
}

func TestKeyframesAreAttached(t *testing.T) {
	reg := New(themes.Default())
	for name, want := range map[string]string{
		"fui-spinner":  "fui-spin",
		"fui-skeleton": "fui-skeleton-pulse",
		"fui-dialog":   "fui-dialog-in",
		"fui-toast":    "fui-toast-in",
	} {
		result, _ := reg.Resolve(name)
		if len(result.Keyframes) != 1 || result.Keyframes[0].Name != want {
			t.Fatalf("%s keyframes = %+v, want %s", name, result.Keyframes, want)
		}
	}
}

func TestResponsivePartsCarryAtRules(t *testing.T) {
	reg := New(themes.Default())
	app, _ := reg.Resolve("fui-app")
	found := false
	for _, block := range app.Blocks {
		if block.AtRule != "" {
			found = true
		}
	}
	if !found {
		t.Fatalf("fui-app has no responsive part: %+v", app.Blocks)
	}
}

func TestUnknownComponentIsNotFound(t *testing.T) {
	reg := New(themes.Default())
	for _, name := range []string{"fui-nope", "btn", "fui-", "fui-btn--"} {
		if _, ok := reg.Resolve(name); ok {
			t.Fatalf("%q should not resolve", name)
		}
	}
}

func TestNamesAreSortedAndStable(t *testing.T) {
	reg := New(themes.Default())
	names := reg.Names()
	for i := 1; i < len(names); i++ {
		if names[i-1] > names[i] {
			t.Fatalf("names are not sorted at %d: %q > %q", i, names[i-1], names[i])
		}
	}
	if reg.Len() != len(names) {
		t.Fatalf("Len() = %d, names = %d", reg.Len(), len(names))
	}
	// Mutating the returned slice must not affect the registry.
	names[0] = "mutated"
	if New(themes.Default()).Names()[0] == "mutated" {
		t.Fatalf("registry names leak")
	}
}

func TestPartsTargetTheirElement(t *testing.T) {
	reg := New(themes.Default())
	body, _ := reg.Resolve("fui-card-body")
	for _, block := range body.Blocks {
		if block.Target != "" && block.Target != "&" {
			continue
		}
		if len(block.Decls) == 0 {
			t.Fatalf("card body has an empty block")
		}
	}
	if _, ok := reg.Resolve("fui-card-body"); !ok {
		t.Fatalf("card body missing")
	}
	var _ ast.Block
}
