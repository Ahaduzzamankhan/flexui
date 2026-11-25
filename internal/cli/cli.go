// Package cli implements the flexui command line interface.
package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Ahaduzzamankhan/flexui/internal/ast"
	"github.com/Ahaduzzamankhan/flexui/internal/compiler"
	"github.com/Ahaduzzamankhan/flexui/internal/config"
	"github.com/Ahaduzzamankhan/flexui/internal/generator"
	"github.com/Ahaduzzamankhan/flexui/internal/scanner"
	"github.com/Ahaduzzamankhan/flexui/internal/watcher"
)

// Main runs the CLI and returns the process exit code.
func Main(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		usage(stdout)
		return 2
	}
	cmd, rest := args[0], args[1:]
	switch cmd {
	case "build", "compile":
		return runBuild(rest, stdout, stderr, false)
	case "watch":
		return runBuild(rest, stdout, stderr, true)
	case "inspect":
		return runInspect(rest, stdout, stderr)
	case "init":
		return runInit(rest, stdout, stderr)
	case "clean":
		return runClean(rest, stdout, stderr)
	case "version", "--version", "-v":
		fmt.Fprintf(stdout, "flexui %s (%s/%s)\n", compiler.Version, runtime.GOOS, runtime.GOARCH)
		return 0
	case "help", "--help", "-h":
		usage(stdout)
		return 0
	default:
		fmt.Fprintf(stderr, "flexui: unknown command %q\n\n", cmd)
		usage(stderr)
		return 2
	}
}

func usage(w io.Writer) {
	fmt.Fprint(w, `flexui `+compiler.Version+` - utility-first CSS compiler

Usage:
  flexui <command> [options] [project-dir]

Commands:
  build      Compile the utilities used by your sources into one stylesheet
  watch      Rebuild whenever a scanned file changes
  inspect    Print resolved rules without writing a stylesheet
  init       Create flexui.config.json and an entry stylesheet
  clean      Remove the build cache (.flexui)
  version    Print the compiler version

Options for build, watch and inspect:
  --config <file>    Config file (default: flexui.config.json)
  --input <file>     Entry CSS file containing @flexui directives
  --output <file>    Output stylesheet
  --content <globs>  Comma separated glob patterns to scan
  --safelist <list>  Comma separated classes that must always be emitted
  --minify           Emit compact CSS
  --watch            Rebuild on change (same as the watch command)
  --no-cache         Ignore the build cache
  --verbose          Print build details
`)
}

type buildFlags struct {
	config   string
	input    string
	output   string
	content  string
	safelist string
	minify   bool
	verbose  bool
	watch    bool
	noCache  bool
}

func newBuildFlags(name string, stderr io.Writer) (*flag.FlagSet, *buildFlags) {
	f := &buildFlags{}
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.StringVar(&f.config, "config", "", "path to flexui.config.json")
	fs.StringVar(&f.input, "input", "", "entry CSS file containing @flexui directives")
	fs.StringVar(&f.output, "output", "", "output stylesheet")
	fs.StringVar(&f.content, "content", "", "comma separated glob patterns to scan")
	fs.StringVar(&f.safelist, "safelist", "", "comma separated classes to always emit")
	fs.BoolVar(&f.minify, "minify", false, "emit compact CSS")
	fs.BoolVar(&f.verbose, "verbose", false, "print build details")
	fs.BoolVar(&f.watch, "watch", false, "rebuild on change")
	fs.BoolVar(&f.noCache, "no-cache", false, "ignore the build cache")
	return fs, f
}

// resolveOptions merges CLI flags over the config file. Only flags that were
// actually passed override the config, so a config file's settings survive.
func resolveOptions(f *buildFlags, fs *flag.FlagSet, root string) (compiler.Options, string, error) {
	set := make(map[string]bool, 9)
	fs.Visit(func(fl *flag.Flag) { set[fl.Name] = true })

	cfgPath := f.config
	if cfgPath == "" && !set["config"] {
		cfgPath = config.Find(root)
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return compiler.Options{}, "", err
	}
	opts := compiler.Options{
		Root:     root,
		Config:   &cfg,
		Content:  cfg.Content,
		Input:    cfg.Input,
		Minify:   cfg.Minify,
		Base:     cfg.BaseEnabled(),
		Safelist: append([]string(nil), cfg.Safelist...),
		NoCache:  f.noCache,
		Verbose:  f.verbose,
	}
	output := cfg.Output
	if set["input"] {
		opts.Input = f.input
	}
	if set["output"] {
		output = f.output
	}
	if set["content"] {
		opts.Content = splitList(f.content)
	}
	if set["safelist"] {
		opts.Safelist = append(opts.Safelist, splitList(f.safelist)...)
	}
	if set["minify"] {
		opts.Minify = f.minify
	}
	if set["verbose"] {
		opts.Verbose = f.verbose
	}
	if output == "" {
		output = "dist/flexui.css"
	}
	return opts, output, nil
}

func runBuild(args []string, stdout, stderr io.Writer, forceWatch bool) int {
	fs, f := newBuildFlags("build", stderr)
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	root := "."
	if fs.NArg() > 0 {
		root = fs.Arg(0)
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		fmt.Fprintf(stderr, "flexui: %v\n", err)
		return 1
	}
	opts, output, err := resolveOptions(f, fs, absRoot)
	if err != nil {
		fmt.Fprintf(stderr, "flexui: %v\n", err)
		return 1
	}
	if forceWatch || f.watch {
		return watchLoop(opts, output, absRoot, stdout, stderr)
	}
	if _, ok := buildOnce(opts, output, absRoot, stdout, stderr); !ok {
		return 1
	}
	return 0
}

// buildOnce compiles and writes the stylesheet. It returns the scanned files so
// watch mode can start from the same set.
func buildOnce(opts compiler.Options, output, root string, stdout, stderr io.Writer) ([]string, bool) {
	if opts.Verbose {
		opts.Log = func(format string, args ...any) {
			fmt.Fprintf(stdout, "  · "+format+"\n", args...)
		}
	}
	res, err := compiler.Build(opts)
	if err != nil {
		fmt.Fprintf(stderr, "flexui: build failed: %v\n", err)
		return nil, false
	}
	outPath := output
	if !filepath.IsAbs(outPath) {
		outPath = filepath.Join(root, outPath)
	}
	if err := compiler.WriteOutput(outPath, []byte(res.CSS)); err != nil {
		fmt.Fprintf(stderr, "flexui: %v\n", err)
		return res.Files, false
	}
	printStats(stdout, res, displayPath(root, outPath))
	return res.Files, true
}

func watchLoop(opts compiler.Options, output, root string, stdout, stderr io.Writer) int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	files, ok := buildOnce(opts, output, root, stdout, stderr)
	if !ok {
		fmt.Fprintln(stderr, "flexui: build failed, waiting for changes")
	}
	prev := watcher.Take(files)
	fmt.Fprintln(stdout, "watching for changes (ctrl+c to stop)")

	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()
	lastScan := time.Now()

	for {
		select {
		case <-ctx.Done():
			fmt.Fprintln(stdout, "stopped")
			return 0
		case <-ticker.C:
		}
		if time.Since(lastScan) > 2*time.Second {
			lastScan = time.Now()
			if scanned, err := scanner.Scan(scanner.Options{Root: root, Content: opts.Content}); err == nil {
				files = scanned
			}
		}
		cur := watcher.Take(files)
		if len(watcher.Diff(prev, cur)) == 0 {
			continue
		}
		// Debounce: editors often write a file in several steps.
		time.Sleep(80 * time.Millisecond)
		cur = watcher.Take(files)
		changed := watcher.Diff(prev, cur)
		if len(changed) == 0 {
			continue
		}
		if opts.Verbose {
			for _, path := range changed {
				fmt.Fprintf(stdout, "  · changed %s\n", displayPath(root, path))
			}
		}
		if built, builtOK := buildOnce(opts, output, root, stdout, stderr); builtOK {
			files = built
		}
		prev = watcher.Take(files)
	}
}

func runInspect(args []string, stdout, stderr io.Writer) int {
	fs, f := newBuildFlags("inspect", stderr)
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	root := "."
	if fs.NArg() > 0 {
		root = fs.Arg(0)
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		fmt.Fprintf(stderr, "flexui: %v\n", err)
		return 1
	}
	opts, _, err := resolveOptions(f, fs, absRoot)
	if err != nil {
		fmt.Fprintf(stderr, "flexui: %v\n", err)
		return 1
	}
	opts.NoCache = true
	res, err := compiler.Build(opts)
	if err != nil {
		fmt.Fprintf(stderr, "flexui: %v\n", err)
		return 1
	}
	rules := append([]ast.Rule(nil), res.Rules...)
	generator.Sort(rules)
	for _, r := range rules {
		at := strings.Join(r.AtRules(), " ")
		for _, block := range r.Blocks {
			parts := make([]string, 0, len(block.Decls))
			for _, d := range block.Decls {
				parts = append(parts, d.Prop+": "+d.Value)
			}
			rule := generator.Selector(&r, block.Target) + " { " + strings.Join(parts, "; ") + " }"
			if at != "" {
				rule = at + " " + rule
			}
			fmt.Fprintln(stdout, rule)
		}
	}
	fmt.Fprintf(stdout, "\n%d rules from %d candidates across %d files\n", len(rules), res.Stats.Candidates, res.Stats.Files)
	return 0
}

func runInit(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(stderr)
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	root := "."
	if fs.NArg() > 0 {
		root = fs.Arg(0)
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		fmt.Fprintf(stderr, "flexui: %v\n", err)
		return 1
	}
	if err := os.MkdirAll(absRoot, 0o755); err != nil {
		fmt.Fprintf(stderr, "flexui: %v\n", err)
		return 1
	}

	initCfg := struct {
		Content  []string `json:"content"`
		Output   string   `json:"output"`
		Minify   bool     `json:"minify"`
		Safelist []string `json:"safelist"`
	}{
		Content: []string{
			"./index.html",
			"./src/**/*.{html,js,jsx,ts,tsx,vue,svelte,astro}",
		},
		Output:   "dist/flexui.css",
		Safelist: []string{},
	}
	data, err := json.MarshalIndent(initCfg, "", "  ")
	if err != nil {
		fmt.Fprintf(stderr, "flexui: %v\n", err)
		return 1
	}
	data = append(data, '\n')

	created := 0
	configPath := filepath.Join(absRoot, "flexui.config.json")
	if _, err := os.Stat(configPath); err == nil {
		fmt.Fprintf(stdout, "  exists  %s\n", displayPath(absRoot, configPath))
	} else if err := os.WriteFile(configPath, data, 0o644); err != nil {
		fmt.Fprintf(stderr, "flexui: %v\n", err)
		return 1
	} else {
		fmt.Fprintf(stdout, "  created %s\n", displayPath(absRoot, configPath))
		created++
	}

	entryPath := filepath.Join(absRoot, "src", "styles.css")
	if _, err := os.Stat(entryPath); err == nil {
		fmt.Fprintf(stdout, "  exists  %s\n", displayPath(absRoot, entryPath))
	} else {
		if err := os.MkdirAll(filepath.Dir(entryPath), 0o755); err != nil {
			fmt.Fprintf(stderr, "flexui: %v\n", err)
			return 1
		}
		entry := "@flexui base;\n@flexui utilities;\n"
		if err := os.WriteFile(entryPath, []byte(entry), 0o644); err != nil {
			fmt.Fprintf(stderr, "flexui: %v\n", err)
			return 1
		}
		fmt.Fprintf(stdout, "  created %s\n", displayPath(absRoot, entryPath))
		created++
	}

	fmt.Fprintf(stdout, "\nNext steps:\n")
	fmt.Fprintf(stdout, "  1. link the stylesheet in your HTML: <link rel=\"stylesheet\" href=\"dist/flexui.css\">\n")
	fmt.Fprintf(stdout, "  2. run: npx flexui build\n")
	if created == 0 {
		fmt.Fprintf(stdout, "\nNothing was overwritten.\n")
	}
	return 0
}

func runClean(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("clean", flag.ContinueOnError)
	fs.SetOutput(stderr)
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	root := "."
	if fs.NArg() > 0 {
		root = fs.Arg(0)
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		fmt.Fprintf(stderr, "flexui: %v\n", err)
		return 1
	}
	dir := filepath.Join(absRoot, ".flexui")
	if err := os.RemoveAll(dir); err != nil {
		fmt.Fprintf(stderr, "flexui: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "removed %s\n", displayPath(absRoot, dir))
	return 0
}

func printStats(w io.Writer, res compiler.Result, out string) {
	suffix := ""
	if res.Stats.CacheHit {
		suffix = " (cached)"
	}
	fmt.Fprintf(w, "flexui %s\n", compiler.Version)
	fmt.Fprintf(w, "  %s\n", out)
	fmt.Fprintf(w, "  %d files -> %d candidates -> %d rules -> %s%s\n",
		res.Stats.Files, res.Stats.Candidates, res.Stats.Rules, humanBytes(res.Stats.Bytes), suffix)
	fmt.Fprintf(w, "  %s\n", res.Stats.Duration.Round(time.Millisecond))
}

func splitList(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func displayPath(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	return filepath.ToSlash(rel)
}

func humanBytes(n int) string {
	switch {
	case n >= 1024*1024:
		return fmt.Sprintf("%.2f MB", float64(n)/(1024*1024))
	case n >= 1024:
		return fmt.Sprintf("%.1f kB", float64(n)/1024)
	default:
		return fmt.Sprintf("%d B", n)
	}
}
