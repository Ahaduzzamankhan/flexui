package config

import (
	"os"
	"path/filepath"
	"testing"
)

func write(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "flexui.config.json")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

func TestLoadThemeOverrides(t *testing.T) {
	cfg, err := Load(write(t, `{
	  "theme": {
	    "colors": {"brand": "#123456"},
	    "spacing": {"18": "4.5rem"},
	    "semanticLight": {"accent": "#ff0000"},
	    "semanticDark": {"accent": "#00ff00"},
	    "breakpoints": {"3xl": "1920px"},
	    "radius": {"card": "1.5rem"}
	  }
	}`))
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.Theme == nil {
		t.Fatalf("theme block not parsed")
	}
	if cfg.Theme.Colors["brand"] != "#123456" {
		t.Fatalf("brand = %q", cfg.Theme.Colors["brand"])
	}
	if cfg.Theme.SemanticDark["accent"] != "#00ff00" {
		t.Fatalf("dark accent = %q", cfg.Theme.SemanticDark["accent"])
	}
	if cfg.Theme.Breakpoints["3xl"] != "1920px" {
		t.Fatalf("breakpoint = %q", cfg.Theme.Breakpoints["3xl"])
	}
}

func TestFontSizeAcceptsStringAndObject(t *testing.T) {
	cfg, err := Load(write(t, `{
	  "theme": {"fontSize": {
	    "huge": "5rem",
	    "hero": {"fontSize": "4rem", "lineHeight": "1", "letterSpacing": "-0.02em"}
	  }}
	}`))
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.Theme.FontSize["huge"].Size != "5rem" {
		t.Fatalf("string form = %+v", cfg.Theme.FontSize["huge"])
	}
	hero := cfg.Theme.FontSize["hero"]
	if hero.Size != "4rem" || hero.LineHeight != "1" || hero.LetterSpacing != "-0.02em" {
		t.Fatalf("object form = %+v", hero)
	}
}

func TestVariantsAnimationsAndFonts(t *testing.T) {
	cfg, err := Load(write(t, `{
	  "variants": {"card": ".card &", "touch": "@media (pointer: coarse)"},
	  "animations": {"wiggle": {
	    "from": {"transform": "rotate(0deg)"},
	    "to": {"transform": "rotate(4deg)"},
	    "duration": "500ms"
	  }},
	  "fonts": [{"family": "Inter", "src": "url(/inter.woff2)", "weight": "400"}],
	  "darkMode": "both",
	  "reducedMotion": "off"
	}`))
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.Variants["card"] != ".card &" {
		t.Fatalf("card variant = %q", cfg.Variants["card"])
	}
	if cfg.Animations["wiggle"].Duration != "500ms" {
		t.Fatalf("animation duration = %q", cfg.Animations["wiggle"].Duration)
	}
	if len(cfg.Fonts) != 1 || cfg.Fonts[0].Family != "Inter" {
		t.Fatalf("fonts = %+v", cfg.Fonts)
	}
	if cfg.DarkMode != "both" {
		t.Fatalf("darkMode = %q", cfg.DarkMode)
	}
	if cfg.ReducedMotion != "off" {
		t.Fatalf("reducedMotion = %q", cfg.ReducedMotion)
	}
}

func TestAnimationFramesForm(t *testing.T) {
	cfg, err := Load(write(t, `{
	  "animations": {"tipper": {"frames": {
	    "0%": {"transform": "rotate(0deg)"},
	    "100%": {"transform": "rotate(10deg)"}
	  }}}
	}`))
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	frames := cfg.Animations["tipper"].Frames
	if len(frames) != 2 || frames["100%"]["transform"] != "rotate(10deg)" {
		t.Fatalf("frames = %+v", frames)
	}
}

func TestDefaultsAreApplied(t *testing.T) {
	if Defaults().DarkMode != "class" || Defaults().ReducedMotion != "safe" {
		t.Fatalf("defaults = %+v", Defaults())
	}
	if !Defaults().BaseEnabled() {
		t.Fatalf("base should default to enabled")
	}
	if _, err := Load(""); err != nil {
		t.Fatalf("empty path: %v", err)
	}
}

func TestInvalidDarkModeIsRejected(t *testing.T) {
	if _, err := Load(write(t, `{"darkMode": "dim"}`)); err == nil {
		t.Fatalf("invalid darkMode accepted")
	}
	if _, err := Load(write(t, `{"reducedMotion": "maybe"}`)); err == nil {
		t.Fatalf("invalid reducedMotion accepted")
	}
}

func TestUnknownFieldsAreRejected(t *testing.T) {
	if _, err := Load(write(t, `{"them": {}}`)); err == nil {
		t.Fatalf("typo in a top-level field was accepted")
	}
}
