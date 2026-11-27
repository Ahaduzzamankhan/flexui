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
	"github.com/Ahaduzzamankhan/flexui/internal/components"
	"github.com/Ahaduzzamankhan/flexui/internal/config"
	"github.com/Ahaduzzamankhan/flexui/internal/extractor"
	"github.com/Ahaduzzamankhan/flexui/internal/generator"
	"github.com/Ahaduzzamankhan/flexui/internal/optimizer"
	"github.com/Ahaduzzamankhan/flexui/internal/resolver"
	"github.com/Ahaduzzamankhan/flexui/internal/scanner"
	"github.com/Ahaduzzamankhan/flexui/internal/themes"
)

// Version is the FlexUI release version reported by the CLI and the compiler.
//
// version.json is the single source of truth. Release builds inject the value
// with -ldflags "-X .../internal/compiler.Version=<version>"; the default below
// keeps `go run` and `go build` working, and TestVersionMatchesVersionFile fails
// if the two ever drift apart.
var Version = "2.0.0"

// Banner is the comment written at the top of generated stylesheets.
var Banner = "FlexUI " + Version + " - generated stylesheet. Do not edit by hand."

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
	Files      int `json:"files"`
	Candidates int `json:"candidates"`
	Rules      int `json:"rules"`
	Bytes      int `json:"bytes"`
	// Reused counts files whose candidates came from the incremental manifest
	// because their size and modification time were unchanged. Parsed counts
	// files that had to be read and extracted again.
	Reused   int           `json:"reused,omitempty"`
	Parsed   int           `json:"parsed,omitempty"`
	Duration time.Duration `json:"-"`
	CacheHit bool          `json:"-"`
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

	// Reuse the manifest's per-file entries when a file's size and modification
	// time are unchanged. An unchanged file is never opened, so a rebuild after
	// editing one file reads that file instead of the whole project. `--no-cache`
	// skips the manifest entirely and re-reads everything.
	manifests := cache.OpenManifest(root)
	var manifest *cache.Manifest
	if !opts.NoCache {
		manifest = manifests.Load()
	}
	parsed := 0
	entries := make([]cache.FileEntry, len(files))
	var readErr error
	var readMu sync.Mutex
	parallel(len(files), workers(), func(i int) {
		path := files[i]
		info, statErr := os.Stat(path)
		if statErr != nil {
			setReadError(&readErr, &readMu, fmt.Errorf("read %s: %w", path, statErr))
			return
		}
		if manifest != nil {
			if entry, ok := manifest.Files[relPath(root, path)]; ok &&
				entry.Size == info.Size() && entry.ModTime == info.ModTime().UnixNano() {
				entries[i] = entry
				return
			}
		}

		data, err := os.ReadFile(path)
		if err != nil {
			setReadError(&readErr, &readMu, fmt.Errorf("read %s: %w", path, err))
			return
		}
		entries[i] = cache.FileEntry{
			Size:    info.Size(),
			ModTime: info.ModTime().UnixNano(),
			Hash:    cache.Hash(data),
			Tokens:  extractor.Extract(path, data),
		}
		readMu.Lock()
		parsed++
		readMu.Unlock()
	})
	if readErr != nil {
		return res, readErr
	}
	res.Stats.Parsed = parsed
	res.Stats.Reused = len(files) - parsed

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
		hashParts = append(hashParts, []byte(relPath(root, f)), []byte(entries[i].Hash))
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
			// The whole build was cached, so no file was read in this run;
			// reporting the previous run's counters would lie.
			res.Stats.Reused = len(files)
			res.Stats.Parsed = 0
			res.Stats.CacheHit = true
			res.Stats.Duration = time.Since(start)
			res.Stats.Bytes = len(data)
			if manifest == nil || parsed > 0 {
				if saveErr := manifests.Save(manifestOf(root, files, entries)); saveErr != nil {
					logf("manifest write failed: %v", saveErr)
				}
			}
			logf("cache hit %s", key)
			return res, nil
		}
	}

	// Candidates are collected from the per-file entries, so a reused file is
	// never parsed again. Order is irrelevant: the set is sorted before
	// resolution, which keeps output deterministic.
	candidates := make(map[string]bool, 1024)
	for i := range files {
		for _, t := range entries[i].Tokens {
			candidates[t] = true
		}
	}
	for _, s := range opts.Safelist {
		for _, f := range strings.Fields(s) {
			if f != "" {
				candidates[f] = true
			}
		}
	}
	res.Stats.Candidates = len(candidates)

	customs := customComponents(cfg.Components)
	if err := components.ValidateCustom(customs); err != nil {
		return res, fmt.Errorf("config: %w", err)
	}

	rules := resolver.NewWith(resolver.Options{
		Theme:            theme,
		DarkMode:         darkMode,
		CustomVariants:   cfg.Variants,
		CustomAnimations: customAnimations(cfg.Animations),
		CustomComponents: customs,
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
	if !opts.NoCache {
		if saveErr := manifests.Save(manifestOf(root, files, entries)); saveErr != nil {
			logf("manifest write failed: %v", saveErr)
		}
	}
	return res, nil
}

// setReadError records the first read failure without a data race.
func setReadError(target *error, mu *sync.Mutex, err error) {
	mu.Lock()
	defer mu.Unlock()
	if *target == nil {
		*target = err
	}
}

// relPath returns a slash-separated path relative to the project root, so the
// manifest is portable between platforms and working directories.
func relPath(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return filepath.ToSlash(path)
	}
	return filepath.ToSlash(rel)
}

// manifestOf builds the manifest for the current scan. Files that disappeared
// are dropped, so the manifest prunes itself as the project changes.
func manifestOf(root string, files []string, entries []cache.FileEntry) *cache.Manifest {
	manifest := cache.NewManifest()
	for i, f := range files {
		manifest.Files[relPath(root, f)] = entries[i]
	}
	return manifest
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

// customComponents converts configuration components into registry definitions
// in sorted order, so registration and output are deterministic.
func customComponents(cfg map[string]config.ComponentConfig) []components.CustomComponent {
	if len(cfg) == 0 {
		return nil
	}
	names := make([]string, 0, len(cfg))
	for name := range cfg {
		names = append(names, name)
	}
	sort.Strings(names)

	out := make([]components.CustomComponent, 0, len(names))
	for _, name := range names {
		c := cfg[name]
		out = append(out, components.CustomComponent{
			Name:     name,
			Extends:  c.Extends,
			Base:     c.Base,
			Variants: c.Variants,
			Parts:    c.Parts,
			States:   c.States,
		})
	}
	return out
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
