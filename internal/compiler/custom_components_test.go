package compiler

import (
	"strings"
	"testing"

	"github.com/Ahaduzzamankhan/flexui/internal/config"
)

func customConfig() config.Config {
	cfg := config.Defaults()
	cfg.Components = map[string]config.ComponentConfig{
		"alert": {
			Base: map[string]string{
				"display":          "flex",
				"gap":              "0.75rem",
				"padding":          "1rem",
				"border-radius":    "0.75rem",
				"background-color": "var(--flexui-surface-raised)",
			},
			Variants: map[string]map[string]string{
				"success": {"color": "var(--flexui-success)"},
				"danger":  {"color": "var(--flexui-danger)"},
			},
			Parts: map[string]map[string]string{
				"title": {"font-weight": "600"},
			},
			States: map[string]map[string]string{
				":hover": {"border-color": "var(--flexui-accent)"},
			},
		},
	}
	return cfg
}

func TestCustomComponentsAreEmittedWhenUsed(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "index.html", `<div class="fui-alert fui-alert--success"><p class="fui-alert-title">Saved</p></div>`)

	cfg := customConfig()
	res, err := Build(Options{Root: root, Content: []string{"./**/*.html"}, Config: &cfg, NoCache: true})
	if err != nil {
		t.Fatalf("build: %v", err)
	}

	for _, want := range []string{
		".fui-alert {",
		"padding: 1rem;",
		".fui-alert--success {",
		"color: var(--flexui-success);",
		".fui-alert-title {",
		"font-weight: 600;",
		".fui-alert:hover {",
		"border-color: var(--flexui-accent);",
	} {
		if !strings.Contains(res.CSS, want) {
			t.Fatalf("missing %q in output:\n%s", want, res.CSS)
		}
	}
	// Only used classes are emitted.
	if strings.Contains(res.CSS, ".fui-alert--danger") {
		t.Fatalf("unused variant was emitted:\n%s", res.CSS)
	}
}

func TestCustomComponentsLiveInTheComponentsLayer(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "index.html", `<div class="fui-alert flex"></div>`)
	cfg := customConfig()
	res, err := Build(Options{Root: root, Content: []string{"./**/*.html"}, Config: &cfg, NoCache: true})
	if err != nil {
		t.Fatal(err)
	}

	componentLayer := res.CSS[strings.Index(res.CSS, "@layer flexui.components {"):]
	utilityIndex := strings.Index(componentLayer, "@layer flexui.utilities {")
	if utilityIndex < 0 {
		t.Fatalf("utilities layer missing:\n%s", res.CSS)
	}
	componentLayer = componentLayer[:utilityIndex]
	if !strings.Contains(componentLayer, ".fui-alert") {
		t.Fatalf("custom component is not in the components layer:\n%s", componentLayer)
	}
	if strings.Contains(componentLayer, ".flex {") {
		t.Fatalf("utility leaked into the components layer:\n%s", componentLayer)
	}
}

func TestCustomComponentsRejectBrokenExtends(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "index.html", `<div class="fui-alert"></div>`)
	cfg := customConfig()
	cfg.Components["alert"] = config.ComponentConfig{Extends: "nope", Base: map[string]string{"display": "flex"}}

	if _, err := Build(Options{Root: root, Content: []string{"./**/*.html"}, Config: &cfg, NoCache: true}); err == nil {
		t.Fatal("build should fail for an unknown extends target")
	}
}

func TestCustomComponentScalesBeyondBuiltins(t *testing.T) {
	root := t.TempDir()
	cfg := config.Defaults()
	cfg.Components = make(map[string]config.ComponentConfig, 120)
	var markup strings.Builder
	for i := 0; i < 120; i++ {
		name := "widget" + string(rune('a'+i%26)) + string(rune('a'+i/26))
		cfg.Components[name] = config.ComponentConfig{
			Base: map[string]string{"display": "grid", "gap": "0.5rem"},
		}
		markup.WriteString(`<div class="fui-` + name + `"></div>`)
	}
	writeFile(t, root, "index.html", markup.String())

	res, err := Build(Options{Root: root, Content: []string{"./**/*.html"}, Config: &cfg, NoCache: true})
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if got := strings.Count(res.CSS, "display: grid;"); got != 120 {
		t.Fatalf("emitted %d custom component rules, want 120:\n%s", got, res.CSS)
	}
}
