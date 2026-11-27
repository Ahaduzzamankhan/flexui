package config

import (
	"os"
	"strings"
	"testing"
)

func TestValidateComponentsAcceptsGoodDefinitions(t *testing.T) {
	components := map[string]ComponentConfig{
		"alert": {
			Base:     map[string]string{"display": "flex", "gap": "0.75rem"},
			Variants: map[string]map[string]string{"success": {"color": "green"}},
			Parts:    map[string]map[string]string{"title": {"font-weight": "600"}},
			States:   map[string]map[string]string{":hover": {"opacity": "0.9"}},
		},
		"hero-panel": {Extends: "card", Base: map[string]string{"padding": "2rem"}},
		"nav-item":   {Extends: "fui-btn"},
	}
	if err := validateComponents("flexui.config.json", components); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateComponentsRejectsBadDefinitions(t *testing.T) {
	cases := map[string]map[string]ComponentConfig{
		"uppercase name":   {"Alert": {Base: map[string]string{"display": "flex"}}},
		"spaced name":      {"alert box": {Base: map[string]string{"display": "flex"}}},
		"double dash":      {"alert--box": {Base: map[string]string{"display": "flex"}}},
		"no declarations":  {"alert": {}},
		"empty value":      {"alert": {Base: map[string]string{"display": " "}}},
		"injected value":   {"alert": {Base: map[string]string{"display": "flex;color:red"}}},
		"injected brace":   {"alert": {Base: map[string]string{"display": "flex}"}}},
		"bad property":     {"alert": {Base: map[string]string{"display;x": "flex"}}},
		"bad variant name": {"alert": {Variants: map[string]map[string]string{"Success": {"color": "green"}}}},
		"bad part name":    {"alert": {Parts: map[string]map[string]string{"_title": {"color": "green"}}}},
		"bad state":        {"alert": {States: map[string]map[string]string{"hover": {"color": "green"}}}},
		"bad extends":      {"alert": {Extends: "Not A Name", Base: map[string]string{"display": "flex"}}},
	}
	for name, components := range cases {
		err := validateComponents("flexui.config.json", components)
		if err == nil {
			t.Fatalf("%s should have failed", name)
		}
		if !strings.Contains(err.Error(), "flexui.config.json") {
			t.Fatalf("%s error should name the config file, got %v", name, err)
		}
	}
}

func TestLoadReadsComponents(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/flexui.config.json"
	body := `{
	  "components": {
	    "alert": {
	      "base": { "display": "flex" },
	      "variants": { "success": { "color": "green" } }
	    }
	  }
	}`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	alert, ok := cfg.Components["alert"]
	if !ok {
		t.Fatal("component was not loaded")
	}
	if alert.Base["display"] != "flex" || alert.Variants["success"]["color"] != "green" {
		t.Fatalf("component = %+v", alert)
	}
}
