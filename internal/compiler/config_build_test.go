package compiler

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Ahaduzzamankhan/flexui/internal/config"
)

func writeProject(t *testing.T, html, cfgJSON string) (string, config.Config) {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte(html), 0o644); err != nil {
		t.Fatalf("write html: %v", err)
	}
	cfgPath := filepath.Join(dir, "flexui.config.json")
	if err := os.WriteFile(cfgPath, []byte(cfgJSON), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	return dir, cfg
}

func TestBuildWiresThemeVariantsAnimationsAndFonts(t *testing.T) {
	dir, cfg := writeProject(t,
		`<!doctype html><div class="fui-btn bg-brand animate-wiggle card:p-4 dark:bg-surface"></div>`,
		`{
		  "content": ["./**/*.html"],
		  "darkMode": "both",
		  "theme": {
		    "colors": {"brand": "#123456"},
		    "semanticLight": {"accent": "#ff0000"},
		    "spacing": {"18": "4.5rem"}
		  },
		  "variants": {"card": ".card &"},
		  "animations": {"wiggle": {
		    "from": {"transform": "rotate(0deg)"},
		    "to": {"transform": "rotate(4deg)"},
		    "duration": "500ms"
		  }},
		  "fonts": [{"family": "Inter", "src": "url(/inter.woff2)", "weight": "400"}]
		}`)

	res, err := Build(Options{Root: dir, Content: cfg.Content, Base: true, NoCache: true, Config: &cfg})
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	css := res.CSS
	for _, want := range []string{
		"@layer flexui.theme, flexui.base, flexui.components, flexui.utilities;",
		"--flexui-accent: #ff0000;",
		"--flexui-space-18: 4.5rem;",
		".fui-btn {",
		"background-color: #123456;",
		"@keyframes fui-wiggle {",
		"animation-name: fui-wiggle;",
		"font-family: \"Inter\";",
		".card .card\\:p-4",
		".dark {",
		":root:not(.light) {",
	} {
		if !strings.Contains(css, want) {
			t.Fatalf("missing %q in:\n%s", want, css)
		}
	}
	if open, close := strings.Count(css, "{"), strings.Count(css, "}"); open != close {
		t.Fatalf("unbalanced braces: %d open, %d close", open, close)
	}
}

func TestBuildIsDeterministicWithConfig(t *testing.T) {
	dir, cfg := writeProject(t, `<div class="p-4 md:p-6 bg-surface"></div>`,
		`{"content": ["./**/*.html"], "theme": {"semanticLight": {"surface": "#ffffff"}}, "variants": {"card": ".card &"}}`)
	opts := Options{Root: dir, Content: cfg.Content, Base: true, NoCache: true, Config: &cfg}
	first, err := Build(opts)
	if err != nil {
		t.Fatalf("first build: %v", err)
	}
	second, err := Build(opts)
	if err != nil {
		t.Fatalf("second build: %v", err)
	}
	if first.CSS != second.CSS {
		t.Fatalf("build is not deterministic")
	}
}

func TestBuildIgnoresUnknownClassesAndKeepsComponents(t *testing.T) {
	dir, cfg := writeProject(t, `<div class="not-a-flexui-class fui-card p-2"></div>`, `{"content": ["./**/*.html"]}`)
	res, err := Build(Options{Root: dir, Content: cfg.Content, Base: true, NoCache: true, Config: &cfg})
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if strings.Contains(res.CSS, "not-a-flexui-class") {
		t.Fatalf("unknown class leaked into the stylesheet")
	}
	if !strings.Contains(res.CSS, ".fui-card {") || !strings.Contains(res.CSS, ".p-2 {") {
		t.Fatalf("known classes missing:\n%s", res.CSS)
	}
}
