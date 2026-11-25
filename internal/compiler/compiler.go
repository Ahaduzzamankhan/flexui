// Package compiler orchestrates scanning, extraction, resolution, generation
// and caching.
package compiler

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Ahaduzzamankhan/flexui/internal/animation"
	"github.com/Ahaduzzamankhan/flexui/internal/ast"
	"github.com/Ahaduzzamankhan/flexui/internal/cache"
	"github.com/Ahaduzzamankhan/flexui/internal/config"
	"github.com/Ahaduzzamankhan/flexui/internal/extractor"
	"github.com/Ahaduzzamankhan/flexui/internal/generator"
	"github.com/Ahaduzzamankhan/flexui/internal/optimizer"
	"github.com/Ahaduzzamankhan/flexui/internal/resolver"
	"github.com/Ahaduzzamankhan/flexui/internal/scanner"
	"github.com/Ahaduzzamankhan/flexui/internal/themes"
)

// Version is the FlexUI release version reported by the CLI and the compiler.
const Version = "2.0.0-alpha.1"

// Banner is the comment written at the top of generated stylesheets.
const Banner = "FlexUI " + Version + " - generated stylesheet. Do not edit by hand."

// Options configures a build.
type Options struct {
	Root     string
	Content  []string
	Input    string
	Output   string
	Minify   bool
	Base     bool
	Safelist []string
	NoCache  bool
	Verbose  bool
	// Config supplies the theme overrides, custom variants, custom animations,
	// fonts and the dark-mode and reduced-motion strategies. Nil uses defaults.
	Config *config.Config
	Log    func(format string, args ...any)
}

// Stats describes a completed build.
type Stats struct {
	Files      int           `json:"files"`
	Candidates int           `json:"candidates"`
	Rules      int           `json:"rules"`
	Bytes      int           `json:"bytes"`
	Duration   time.Duration `json:"-"`
	CacheHit   bool          `json:"-"`
}

// Result is a compiled stylesheet plus statistics.
type Result struct {
	CSS      string
	Stats    Stats
	Files    []string
	Rules    []ast.Rule
	CacheKey string
}

// Build runs the full pipeline. It does not write the output file; callers do,
// so the compiler stays usable as a library.
func Build(opts Options) (Result, error) {
	start := time.Now()
	res := Result{}
	logf := opts.Log
	if logf == nil {
		logf = func(string, ...any) {}
	}
	root, err := filepath.Abs(opts.Root)
	if err != nil {
		return res, err
	}
	content := opts.Content
	if len(content) == 0 {
		content = []string{"./**/*.{html,js,jsx,ts,tsx,vue,svelte,astro,md,mdx}"}
	}
	files, err := scanner.Scan(scanner.Options{Root: root, Content: content})
	if err != nil {
		return res, fmt.Errorf("scan content: %w", err)
	}
	res.Files = files
	res.Stats.Files = len(files)
	logf("scanned %d files", len(files))

	// Read every source file concurrently, then hash them for the cache key.
	sources := make([][]byte, len(files))
	var readErr error
	var readMu sync.Mutex
	parallel(len(files), workers(), func(i int) {
		data, err := os.ReadFile(files[i])
		if err != nil {
			readMu.Lock()
			if readErr == nil {
				readErr = fmt.Errorf("read %s: %w", files[i], err)
			}
			readMu.Unlock()
			return
		}
		sources[i] = data
	})
	if readErr != nil {
		return res, readErr
	}

	baseEnabled, utilitiesEnabled, passthrough, err := readInput(root, opts.Input, opts.Base)
	if err != nil {
		return res, err
	}

	cfg := opts.Config
	if cfg == nil {
		defaults := config.Defaults()
		cfg = &defaults
	}
	darkMode := cfg.DarkMode
	if darkMode == "" {
		darkMode = "class"
	}
	theme := themes.Default()
	applyTheme(theme, cfg.Theme)
	cfgJSON, _ := json.Marshal(cfg)

	hashParts := [][]byte{
		[]byte(Version),
		[]byte(buildFingerprint()),
		[]byte(fmt.Sprintf("minify=%t;base=%t;utils=%t;utilsLayer=%t", opts.Minify, baseEnabled, utilitiesEnabled, utilitiesEnabled)),
		[]byte(passthrough),
		[]byte(strings.Join(opts.Safelist, ",")),
		cfgJSON,
	}
	for i, f := range files {
		rel, relErr := filepath.Rel(root, f)
		if relErr != nil {
			rel = f
		}
		hashParts = append(hashParts, []byte(filepath.ToSlash(rel)), sources[i])
	}
	key := cache.Key(hashParts...)
	res.CacheKey = key

	c := cache.Open(root)
	if !opts.NoCache {
		if data, meta, ok := c.Get(key); ok {
			res.CSS = string(data)
			res.Stats.Bytes = len(data)
			res.Stats.CacheHit = true
			res.Stats.Duration = time.Since(start)
			if len(meta) > 0 {
				_ = json.Unmarshal(meta, &res.Stats)
			}
			res.Stats.CacheHit = true
			res.Stats.Duration = time.Since(start)
			res.Stats.Bytes = len(data)
			logf("cache hit %s", key)
			return res, nil
		}
	}

	// Extract candidates concurrently. Order is irrelevant because the set is
	// sorted before resolution, which keeps output deterministic.
	candidates := make(map[string]bool, 1024)
	var candMu sync.Mutex
	parallel(len(files), workers(), func(i int) {
		tokens := extractor.Extract(files[i], sources[i])
		if len(tokens) == 0 {
			return
		}
		candMu.Lock()
		for _, t := range tokens {
			candidates[t] = true
		}
		candMu.Unlock()
	})
	for _, s := range opts.Safelist {
		for _, f := range strings.Fields(s) {
			if f != "" {
				candidates[f] = true
			}
		}
	}
	res.Stats.Candidates = len(candidates)

	rules := resolver.NewWith(resolver.Options{
		Theme:            theme,
		DarkMode:         darkMode,
		CustomVariants:   cfg.Variants,
		CustomAnimations: customAnimations(cfg.Animations),
	})
	names := make([]string, 0, len(candidates))
	for name := range candidates {
		names = append(names, name)
	}
	sort.Strings(names)

	// Resolution runs concurrently, but results are placed by index so the
	// output order is independent of scheduling.
	resolved := make([]ast.Rule, len(names))
	known := make([]bool, len(names))
	parallel(len(names), workers(), func(i int) {
		rule, ok := rules.Resolve(names[i])
		if ok {
			resolved[i] = rule
			known[i] = true
		}
	})
	kept := make([]ast.Rule, 0, len(resolved))
	for i, ok := range known {
		if ok {
			kept = append(kept, resolved[i])
		}
	}
	kept = optimizer.Rules(kept)
	res.Rules = kept
	res.Stats.Rules = len(kept)

	css := generator.Generate(kept, generator.Options{
		Minify:        opts.Minify,
		Base:          baseEnabled,
		Banner:        Banner,
		Passthrough:   passthrough,
		Theme:         theme,
		DarkMode:      darkMode,
		Fonts:         fontFaces(cfg.Fonts),
		ReducedMotion: cfg.ReducedMotion != "off",
	})
	res.CSS = string(css)
	res.Stats.Bytes = len(css)
	res.Stats.Duration = time.Since(start)

	meta, marshalErr := json.Marshal(res.Stats)
	if marshalErr == nil {
		if putErr := c.Put(key, css, meta); putErr != nil {
			logf("cache write failed: %v", putErr)
		}
	}
	return res, nil
}

// applyTheme merges the configuration overrides into the default theme.
func applyTheme(t *themes.Theme, cfg *config.ThemeConfig) {
	if cfg == nil {
		return
	}
	o := themes.Overrides{
		Colors:         cfg.Colors,
		Spacing:        cfg.Spacing,
		FontFamily:     cfg.FontFamily,
		FontWeight:     cfg.FontWeight,
		LineHeight:     cfg.LineHeight,
		LetterSpacing:  cfg.LetterSpacing,
		Radius:         cfg.Radius,
		Shadow:         cfg.Shadow,
		Opacity:        cfg.Opacity,
		ZIndex:         cfg.ZIndex,
		Breakpoints:    cfg.Breakpoints,
		Easing:         cfg.Easing,
		Duration:       cfg.Duration,
		AspectRatio:    cfg.AspectRatio,
		MaxWidth:       cfg.MaxWidth,
		ContainerSizes: cfg.ContainerSizes,
		SemanticLight:  cfg.SemanticLight,
		SemanticDark:   cfg.SemanticDark,
	}
	if len(cfg.FontSize) > 0 {
		o.FontSize = make(map[string]themes.FontSize, len(cfg.FontSize))
		for key, size := range cfg.FontSize {
			o.FontSize[key] = themes.FontSize{
				Size:          size.Size,
				LineHeight:    size.LineHeight,
				LetterSpacing: size.LetterSpacing,
			}
		}
	}
	t.Apply(o)
}

// customAnimations converts configuration animations into definitions.
func customAnimations(cfg map[string]config.AnimationConfig) map[string]animation.Definition {
	if len(cfg) == 0 {
		return nil
	}
	out := make(map[string]animation.Definition, len(cfg))
	for name, a := range cfg {
		def := animation.Definition{
			Duration:   a.Duration,
			Easing:     a.Easing,
			Iterations: a.Iterations,
			Fill:       a.Fill,
			Direction:  a.Direction,
			Delay:      a.Delay,
		}
		if len(a.Frames) > 0 {
			steps := make([]string, 0, len(a.Frames))
			for step := range a.Frames {
				steps = append(steps, step)
			}
			sort.Slice(steps, func(i, j int) bool { return frameOrder(steps[i]) < frameOrder(steps[j]) })
			for _, step := range steps {
				def.Steps = append(def.Steps, ast.KeyframeStep{Steps: step, Decls: sortedDecls(a.Frames[step])})
			}
		} else {
			if len(a.From) > 0 {
				def.Steps = append(def.Steps, ast.KeyframeStep{Steps: "from", Decls: sortedDecls(a.From)})
			}
			if len(a.To) > 0 {
				def.Steps = append(def.Steps, ast.KeyframeStep{Steps: "to", Decls: sortedDecls(a.To)})
			}
		}
		out[name] = def
	}
	return out
}

// sortedDecls converts a property map into declarations in property order.
func sortedDecls(props map[string]string) []ast.Decl {
	keys := make([]string, 0, len(props))
	for prop := range props {
		keys = append(keys, prop)
	}
	sort.Strings(keys)
	out := make([]ast.Decl, 0, len(keys))
	for _, prop := range keys {
		out = append(out, ast.Decl{Prop: prop, Value: props[prop]})
	}
	return out
}

// frameOrder sorts keyframe steps the way they read: from, percentages, to.
func frameOrder(step string) string {
	switch strings.TrimSpace(step) {
	case "from":
		return "0000"
	case "to":
		return "9999"
	}
	n := 0
	for i := 0; i < len(step) && step[i] >= '0' && step[i] <= '9'; i++ {
		n = n*10 + int(step[i]-'0')
	}
	return fmt.Sprintf("%04d", n)
}

// fontFaces converts configured @font-face rules.
func fontFaces(cfg []config.FontFaceConfig) []themes.FontFace {
	if len(cfg) == 0 {
		return nil
	}
	out := make([]themes.FontFace, 0, len(cfg))
	for _, f := range cfg {
		out = append(out, themes.FontFace{
			Family:  f.Family,
			Src:     f.Src,
			Weight:  f.Weight,
			Style:   f.Style,
			Display: f.Display,
		})
	}
	return out
}

// readInput resolves the entry CSS file and its @flexui directives.
func readInput(root, input string, baseDefault bool) (base, utilities bool, passthrough string, err error) {
	base, utilities = baseDefault, true
	if input == "" {
		return base, utilities, "", nil
	}
	path := input
	if !filepath.IsAbs(path) {
		path = filepath.Join(root, path)
	}
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		return false, false, "", fmt.Errorf("read input CSS %s: %w", input, readErr)
	}
	base, utilities, passthrough, err = ParseDirectives(string(data))
	if err != nil {
		return false, false, "", fmt.Errorf("%s: %w", input, err)
	}
	if !base {
		base = baseDefault
	}
	return base, utilities, passthrough, nil
}

// ParseDirectives extracts @flexui directives from an entry stylesheet and
// returns the remaining CSS, which is passed through untouched.
func ParseDirectives(src string) (base, utilities bool, passthrough string, err error) {
	var kept []string
	for _, line := range strings.Split(src, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "@flexui") {
			kept = append(kept, line)
			continue
		}
		directive := strings.TrimSpace(strings.TrimPrefix(trimmed, "@flexui"))
		directive = strings.TrimSpace(strings.TrimSuffix(directive, ";"))
		switch directive {
		case "base":
			base = true
		case "utilities":
			utilities = true
		default:
			if directive == "" {
				return false, false, "", fmt.Errorf("empty @flexui directive")
			}
			return false, false, "", fmt.Errorf("unknown @flexui directive %q (expected base or utilities)", directive)
		}
	}
	return base, utilities, strings.TrimSpace(strings.Join(kept, "\n")), nil
}

// WriteOutput writes the stylesheet, creating parent directories as needed.
func WriteOutput(path string, css []byte) error {
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create output directory: %w", err)
		}
	}
	if err := os.WriteFile(path, css, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

var (
	fingerprintOnce sync.Once
	fingerprint     string
)

// buildFingerprint identifies the compiler binary that produced a cached entry.
// Without it, editing the compiler during development would silently reuse a
// stylesheet generated by the previous build, because the cache key only
// otherwise covers inputs and options.
func buildFingerprint() string {
	fingerprintOnce.Do(func() {
		fingerprint = "unknown"
		exe, err := os.Executable()
		if err != nil {
			return
		}
		info, err := os.Stat(exe)
		if err != nil {
			return
		}
		fingerprint = fmt.Sprintf("%s-%d-%d", filepath.Base(exe), info.Size(), info.ModTime().UnixNano())
	})
	return fingerprint
}

func workers() int {
	n := runtime.NumCPU()
	if n < 1 {
		return 1
	}
	return n
}

// parallel runs fn for 0..n-1 on a bounded worker pool.
func parallel(n, limit int, fn func(int)) {
	if n == 0 {
		return
	}
	if limit < 1 {
		limit = 1
	}
	if limit > n {
		limit = n
	}
	var wg sync.WaitGroup
	ch := make(chan int)
	wg.Add(limit)
	for i := 0; i < limit; i++ {
		go func() {
			defer wg.Done()
			for idx := range ch {
				fn(idx)
			}
		}()
	}
	for i := 0; i < n; i++ {
		ch <- i
	}
	close(ch)
	wg.Wait()
}
