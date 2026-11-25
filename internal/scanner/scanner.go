package scanner

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// SupportedExtensions are the file types FlexUI scans by default.
var SupportedExtensions = map[string]bool{
	".html": true, ".htm": true, ".md": true, ".mdx": true,
	".js": true, ".jsx": true, ".mjs": true, ".cjs": true,
	".ts": true, ".tsx": true, ".mts": true, ".cts": true,
	".vue": true, ".svelte": true, ".astro": true,
	".php": true, ".twig": true, ".hbs": true, ".erb": true,
	".ejs": true, ".njk": true, ".liquid": true, ".pug": true,
	".gohtml": true, ".gotmpl": true,
}

// IgnoredDirs are never traversed.
var IgnoredDirs = map[string]bool{
	".git": true, ".hg": true, ".svn": true, "node_modules": true,
	"vendor": true, "dist": true, "build": true, "out": true,
	".next": true, ".nuxt": true, ".svelte-kit": true, ".output": true,
	"coverage": true, ".cache": true, ".flexui": true, ".toolchain": true,
	".turbo": true, ".vercel": true, "target": true, "tmp": true,
}

// Options configures a scan.
type Options struct {
	// Root is the project directory. Patterns are relative to it.
	Root string
	// Content holds glob patterns or directories, relative to Root.
	Content []string
	// Extensions overrides SupportedExtensions when non-empty.
	Extensions map[string]bool
}

// Scan returns the sorted, de-duplicated list of matching files as absolute
// paths, so later stages do not depend on the working directory.
func Scan(opts Options) ([]string, error) {
	root, err := filepath.Abs(opts.Root)
	if err != nil {
		return nil, err
	}
	exts := opts.Extensions
	if len(exts) == 0 {
		exts = SupportedExtensions
	}
	var patterns, dirs []string
	for _, c := range opts.Content {
		c = normalize(c)
		if c == "" {
			continue
		}
		if info, statErr := os.Stat(filepath.Join(root, filepath.FromSlash(c))); statErr == nil && info.IsDir() {
			dirs = append(dirs, strings.TrimSuffix(c, "/")+"/")
			continue
		}
		patterns = append(patterns, c)
	}
	seen := make(map[string]bool, 256)
	files := make([]string, 0, 256)
	add := func(abs string) {
		if !seen[abs] {
			seen[abs] = true
			files = append(files, abs)
		}
	}

	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			if d != nil && d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			if path != root && IgnoredDirs[d.Name()] {
				return fs.SkipDir
			}
			return nil
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return nil
		}
		rel = filepath.ToSlash(rel)
		if !exts[strings.ToLower(filepath.Ext(rel))] {
			return nil
		}
		for _, dir := range dirs {
			if strings.HasPrefix(rel, dir) {
				add(path)
				return nil
			}
		}
		for _, p := range patterns {
			if Match(p, rel) {
				add(path)
				return nil
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

// IsContentFile reports whether a path has a scannable extension.
func IsContentFile(path string) bool {
	return SupportedExtensions[strings.ToLower(filepath.Ext(path))]
}
