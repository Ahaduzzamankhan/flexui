package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMatch(t *testing.T) {
	cases := []struct {
		pattern, path string
		want          bool
	}{
		{"**/*.html", "index.html", true},
		{"**/*.html", "src/pages/index.html", true},
		{"./src/**/*.{ts,tsx}", "src/a/b.tsx", true},
		{"./src/**/*.{ts,tsx}", "src/a/b.js", false},
		{"src/**/*.ts", "other/a.ts", false},
		{"*.html", "src/index.html", false},
		{"**/*.html", "src/index.htmlx", false},
		{"src/*/index.html", "src/a/index.html", true},
		{"src/*/index.html", "src/a/b/index.html", false},
	}
	for _, tc := range cases {
		if got := Match(tc.pattern, tc.path); got != tc.want {
			t.Fatalf("Match(%q, %q) = %v, want %v", tc.pattern, tc.path, got, tc.want)
		}
	}
}

func TestScanFindsFilesAndSkipsIgnored(t *testing.T) {
	root := t.TempDir()
	write := func(rel, body string) {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("index.html", "<div class=\"p-4\"></div>")
	write("src/app.tsx", "export const x = 'flex';")
	write("src/nested/notes.md", "# notes")
	write("node_modules/pkg/index.html", "<div class=\"p-8\"></div>")

	files, err := Scan(Options{
		Root:    root,
		Content: []string{"./**/*.html", "./src/**/*.{ts,tsx}", "./src/nested"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 3 {
		t.Fatalf("got %d files: %v", len(files), files)
	}
	for _, f := range files {
		if filepath.Base(filepath.Dir(f)) == "node_modules" {
			t.Fatalf("node_modules was scanned: %v", files)
		}
	}
}

func TestScanIsSortedAndDeterministic(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"b.html", "a.html", "c.html"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte("<div></div>"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	first, err := Scan(Options{Root: root, Content: []string{"./**/*.html"}})
	if err != nil {
		t.Fatal(err)
	}
	second, err := Scan(Options{Root: root, Content: []string{"./**/*.html"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(first) != 3 {
		t.Fatalf("got %v", first)
	}
	for i := range first {
		if first[i] != second[i] {
			t.Fatalf("order changed between scans: %v vs %v", first, second)
		}
	}
}
