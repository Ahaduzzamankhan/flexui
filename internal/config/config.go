// Package config loads flexui.config.json.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// FileNames are searched in order when --config is not given.
var FileNames = []string{"flexui.config.json", "flexui.json"}

// DefaultContent is the default scan set relative to the project root.
var DefaultContent = []string{
	"./**/*.{html,htm,md,mdx,js,jsx,mjs,cjs,ts,tsx,mts,cts,vue,svelte,astro,php,twig,hbs,erb,ejs,njk,liquid,pug,gohtml,gotmpl}",
}

// Config mirrors flexui.config.json. Unknown fields are rejected so typos are
// reported instead of silently ignored.
type Config struct {
	Content  []string `json:"content"`
	Output   string   `json:"output"`
	Input    string   `json:"input"`
	Minify   bool     `json:"minify"`
	Safelist []string `json:"safelist"`

	// Base controls the preflight reset. Nil means enabled.
	Base *bool `json:"base"`

	// DarkMode selects how dark: variants and the dark token block are applied:
	// "class" (default, `.dark` ancestor), "media" (prefers-color-scheme) or
	// "both" (class selector plus a system default).
	DarkMode string `json:"darkMode"`

	// ReducedMotion is "safe" (default) or "off". Safe mode neutralises
	// animation durations inside prefers-reduced-motion.
	ReducedMotion string `json:"reducedMotion"`

	// Theme overrides individual design tokens.
	Theme *ThemeConfig `json:"theme"`

	// Variants defines custom variants: a name mapped to a selector template
	// using "&" (for example "&:hover, &:focus") or to an at-rule prefixed with
	// "@" (for example "@media (pointer: coarse)").
	Variants map[string]string `json:"variants"`

	// Animations defines custom keyframe animations.
	Animations map[string]AnimationConfig `json:"animations"`

	// Fonts declares local @font-face rules.
	Fonts []FontFaceConfig `json:"fonts"`
}

// ThemeConfig is the theme override block. Every scale is optional.
type ThemeConfig struct {
	Colors         map[string]string   `json:"colors"`
	Spacing        map[string]string   `json:"spacing"`
	FontSize       map[string]FontSize `json:"fontSize"`
	FontFamily     map[string]string   `json:"fontFamily"`
	FontWeight     map[string]string   `json:"fontWeight"`
	LineHeight     map[string]string   `json:"lineHeight"`
	LetterSpacing  map[string]string   `json:"letterSpacing"`
	Radius         map[string]string   `json:"radius"`
	Shadow         map[string]string   `json:"shadow"`
	Opacity        map[string]string   `json:"opacity"`
	ZIndex         map[string]string   `json:"zIndex"`
	Breakpoints    map[string]string   `json:"breakpoints"`
	Easing         map[string]string   `json:"easing"`
	Duration       map[string]string   `json:"duration"`
	AspectRatio    map[string]string   `json:"aspectRatio"`
	MaxWidth       map[string]string   `json:"maxWidth"`
	ContainerSizes map[string]string   `json:"containerSizes"`
	SemanticLight  map[string]string   `json:"semanticLight"`
	SemanticDark   map[string]string   `json:"semanticDark"`
}

// FontSize accepts either a plain string ("1.5rem") or an object with size,
// lineHeight and letterSpacing, mirroring the documented typography config.
type FontSize struct {
	Size          string `json:"fontSize"`
	LineHeight    string `json:"lineHeight"`
	LetterSpacing string `json:"letterSpacing"`
}

// UnmarshalJSON implements the string-or-object form.
func (f *FontSize) UnmarshalJSON(data []byte) error {
	trimmed := firstNonSpace(data)
	if trimmed == '"' {
		var size string
		if err := json.Unmarshal(data, &size); err != nil {
			return err
		}
		f.Size = size
		return nil
	}
	type alias FontSize
	var value alias
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*f = FontSize(value)
	return nil
}

func firstNonSpace(data []byte) byte {
	for _, c := range data {
		switch c {
		case ' ', '\t', '\n', '\r':
			continue
		default:
			return c
		}
	}
	return 0
}

// AnimationConfig defines a custom animation. Frames may be given as from/to or
// as explicit percentage steps.
type AnimationConfig struct {
	From       map[string]string            `json:"from"`
	To         map[string]string            `json:"to"`
	Frames     map[string]map[string]string `json:"frames"`
	Duration   string                       `json:"duration"`
	Easing     string                       `json:"easing"`
	Iterations string                       `json:"iterations"`
	Fill       string                       `json:"fillMode"`
	Direction  string                       `json:"direction"`
	Delay      string                       `json:"delay"`
}

// FontFaceConfig declares a local font.
type FontFaceConfig struct {
	Family  string `json:"family"`
	Src     string `json:"src"`
	Weight  string `json:"weight"`
	Style   string `json:"style"`
	Display string `json:"display"`
}

// Defaults returns the configuration used when no file exists.
func Defaults() Config {
	base := true
	return Config{
		Content:       append([]string(nil), DefaultContent...),
		Output:        "dist/flexui.css",
		Base:          &base,
		DarkMode:      "class",
		ReducedMotion: "safe",
	}
}

// BaseEnabled reports whether preflight should be emitted.
func (c Config) BaseEnabled() bool {
	if c.Base == nil {
		return true
	}
	return *c.Base
}

// Load reads a config file. An empty path returns the defaults.
func Load(path string) (Config, error) {
	cfg := Defaults()
	if path == "" {
		return cfg, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("read config %s: %w", path, err)
	}
	decoder := json.NewDecoder(newTrimReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cfg); err != nil {
		return cfg, fmt.Errorf("parse config %s: %w", path, err)
	}
	if len(cfg.Content) == 0 {
		cfg.Content = append([]string(nil), DefaultContent...)
	}
	if cfg.Output == "" {
		cfg.Output = "dist/flexui.css"
	}
	if cfg.DarkMode == "" {
		cfg.DarkMode = "class"
	}
	switch cfg.DarkMode {
	case "class", "media", "both":
	default:
		return cfg, fmt.Errorf("parse config %s: darkMode must be class, media or both, got %q", path, cfg.DarkMode)
	}
	if cfg.ReducedMotion == "" {
		cfg.ReducedMotion = "safe"
	}
	switch cfg.ReducedMotion {
	case "safe", "off":
	default:
		return cfg, fmt.Errorf("parse config %s: reducedMotion must be safe or off, got %q", path, cfg.ReducedMotion)
	}
	return cfg, nil
}

// Find returns the first config file present in root, or "" when there is none.
func Find(root string) string {
	for _, name := range FileNames {
		path := filepath.Join(root, name)
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return path
		}
	}
	return ""
}
