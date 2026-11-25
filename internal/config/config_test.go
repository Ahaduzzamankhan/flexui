package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Content) == 0 {
		t.Fatal("default content patterns missing")
	}
	if cfg.Output != "dist/flexui.css" {
		t.Fatalf("output = %q", cfg.Output)
	}
	if !cfg.BaseEnabled() {
		t.Fatal("base should default to enabled")
	}
}

func TestLoadFileWithBOM(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "flexui.config.json")
	body := "\xEF\xBB\xBF" + `{"content":["./src/**/*.tsx"],"output":"build/app.css","minify":true,"safelist":["md:flex"]}`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Content) != 1 || cfg.Content[0] != "./src/**/*.tsx" {
		t.Fatalf("content = %v", cfg.Content)
	}
	if cfg.Output != "build/app.css" || !cfg.Minify {
		t.Fatalf("cfg = %+v", cfg)
	}
	if len(cfg.Safelist) != 1 || cfg.Safelist[0] != "md:flex" {
		t.Fatalf("safelist = %v", cfg.Safelist)
	}
	if !cfg.BaseEnabled() {
		t.Fatal("base should stay enabled when omitted")
	}
}

func TestLoadRejectsUnknownFields(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "flexui.config.json")
	if err := os.WriteFile(path, []byte(`{"contents":["./a.html"]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("typo in config was accepted")
	}
}

func TestLoadBaseOverride(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "flexui.config.json")
	if err := os.WriteFile(path, []byte(`{"base":false}`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.BaseEnabled() {
		t.Fatal("base override ignored")
	}
}

func TestFind(t *testing.T) {
	root := t.TempDir()
	if got := Find(root); got != "" {
		t.Fatalf("found %q in an empty directory", got)
	}
	path := filepath.Join(root, "flexui.json")
	if err := os.WriteFile(path, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := Find(root); got != path {
		t.Fatalf("Find = %q, want %q", got, path)
	}
}
