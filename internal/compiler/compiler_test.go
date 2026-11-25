package compiler

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFile(t *testing.T, root, rel, body string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestBuildEndToEnd(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "index.html", `<!-- mx-auto should not be compiled -->
<div class="p-4 flex md:hover:bg-red-500">hello</div>`)
	writeFile(t, root, "src/Card.tsx", `export const Card = () => <div className="rounded-lg shadow-md p-4" />;`)

	res, err := Build(Options{Root: root, Content: []string{"./**/*.{html,tsx}"}, Base: true})
	if err != nil {
		t.Fatalf("build failed: %v", err)
	}
	if res.Stats.Files != 2 {
		t.Fatalf("files = %d, want 2", res.Stats.Files)
	}
	want := []string{
		".p-4 { padding: 1rem; }",
		".flex { display: flex; }",
		`.md\:hover\:bg-red-500:hover`,
		".rounded-lg { border-radius: 0.75rem; }",
	}
	for _, w := range want {
		if !strings.Contains(res.CSS, w) {
			t.Fatalf("missing %q in output:\n%s", w, res.CSS)
		}
	}
	if strings.Contains(res.CSS, ".mx-auto") {
		t.Fatalf("utility from an HTML comment was compiled:\n%s", res.CSS)
	}
	if !strings.Contains(res.CSS, "@layer flexui.base") {
		t.Fatalf("preflight missing")
	}
	if res.Stats.Rules < 4 {
		t.Fatalf("rules = %d", res.Stats.Rules)
	}
}

func TestBuildCacheHit(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "index.html", `<div class="p-4"></div>`)
	opts := Options{Root: root, Content: []string{"./**/*.html"}, Base: true}
	first, err := Build(opts)
	if err != nil {
		t.Fatal(err)
	}
	if first.Stats.CacheHit {
		t.Fatal("first build reported a cache hit")
	}
	second, err := Build(opts)
	if err != nil {
		t.Fatal(err)
	}
	if !second.Stats.CacheHit {
		t.Fatal("second build did not use the cache")
	}
	if second.CSS != first.CSS {
		t.Fatal("cached output differs from fresh output")
	}
	if second.Stats.Rules != first.Stats.Rules {
		t.Fatalf("cached stats lost rule count: %d vs %d", second.Stats.Rules, first.Stats.Rules)
	}
}

func TestBuildCacheInvalidatedByContentChange(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "index.html", `<div class="p-4"></div>`)
	opts := Options{Root: root, Content: []string{"./**/*.html"}, Base: true}
	if _, err := Build(opts); err != nil {
		t.Fatal(err)
	}
	writeFile(t, root, "index.html", `<div class="p-6"></div>`)
	res, err := Build(opts)
	if err != nil {
		t.Fatal(err)
	}
	if res.Stats.CacheHit {
		t.Fatal("stale cache entry was reused")
	}
	if !strings.Contains(res.CSS, ".p-6 { padding: 1.5rem; }") {
		t.Fatalf("new utility missing:\n%s", res.CSS)
	}
}

func TestBuildIsDeterministic(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "a.html", `<div class="p-4 md:flex hover:bg-red-500"></div>`)
	writeFile(t, root, "b.html", `<div class="flex grid p-4"></div>`)
	opts := Options{Root: root, Content: []string{"./**/*.html"}, Base: true, NoCache: true}
	first, err := Build(opts)
	if err != nil {
		t.Fatal(err)
	}
	second, err := Build(opts)
	if err != nil {
		t.Fatal(err)
	}
	if first.CSS != second.CSS {
		t.Fatal("two identical builds produced different CSS")
	}
}

func TestBuildMinify(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "index.html", `<div class="p-4"></div>`)
	res, err := Build(Options{Root: root, Content: []string{"./**/*.html"}, Base: false, Minify: true, NoCache: true})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(res.CSS, "\n  ") {
		t.Fatalf("minified output is indented:\n%q", res.CSS)
	}
	if !strings.Contains(res.CSS, ".p-4{padding:1rem;}") {
		t.Fatalf("compact rule missing:\n%q", res.CSS)
	}
}

func TestSafelistKeepsClassWithoutMarkup(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "index.html", `<div class="p-4"></div>`)
	res, err := Build(Options{
		Root: root, Content: []string{"./**/*.html"}, Base: false, NoCache: true,
		Safelist: []string{"md:flex", "dialog-open"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.CSS, `.md\:flex`) {
		t.Fatalf("safelisted utility missing:\n%s", res.CSS)
	}
	if strings.Contains(res.CSS, "dialog-open") {
		t.Fatalf("unknown safelisted class should be ignored:\n%s", res.CSS)
	}
}

func TestInputDirectives(t *testing.T) {
	base, utilities, passthrough, err := ParseDirectives("@flexui base;\n@flexui utilities;\n.card { color: red; }")
	if err != nil {
		t.Fatal(err)
	}
	if !base || !utilities {
		t.Fatalf("directives not applied: %v %v", base, utilities)
	}
	if passthrough != ".card { color: red; }" {
		t.Fatalf("passthrough = %q", passthrough)
	}
	if _, _, _, err := ParseDirectives("@flexui bogus;"); err == nil {
		t.Fatal("unknown directive was accepted")
	}
}

func TestInputFileIsUsed(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "index.html", `<div class="p-4"></div>`)
	writeFile(t, root, "src/styles.css", "@flexui utilities;\n.hero { color: red; }\n")
	res, err := Build(Options{
		Root: root, Content: []string{"./**/*.html"}, Input: "src/styles.css", Base: true, NoCache: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.CSS, ".hero { color: red; }") {
		t.Fatalf("passthrough CSS missing:\n%s", res.CSS)
	}
	if !strings.Contains(res.CSS, "@layer flexui.base") {
		t.Fatalf("base layer missing when config enabled it")
	}
}

func TestWriteOutputCreatesDirectories(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "dist", "nested", "flexui.css")
	if err := WriteOutput(path, []byte(".a{}")); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != ".a{}" {
		t.Fatalf("content = %q", data)
	}
}
