package compiler

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Ahaduzzamankhan/flexui/internal/cache"
)

func manifestPath(root string) string {
	return filepath.Join(root, ".flexui", "manifest.json")
}

func readManifest(t *testing.T, root string) *cache.Manifest {
	t.Helper()
	manifest := cache.OpenManifest(root).Load()
	if manifest == nil {
		t.Fatal("manifest was not written")
	}
	return manifest
}

// TestUnchangedRebuildHitsManifest proves a second build reads no source file
// and reports what it reused.
func TestUnchangedRebuildHitsManifest(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "a.html", `<div class="p-4 flex"></div>`)
	writeFile(t, root, "b.html", `<div class="m-2"></div>`)

	opts := Options{Root: root, Content: []string{"./**/*.html"}, Base: false}
	first, err := Build(opts)
	if err != nil {
		t.Fatalf("first build: %v", err)
	}
	if first.Stats.Parsed != 2 || first.Stats.Reused != 0 {
		t.Fatalf("first build parsed/reused = %d/%d, want 2/0", first.Stats.Parsed, first.Stats.Reused)
	}
	if first.Stats.CacheHit {
		t.Fatal("first build reported a whole-build cache hit")
	}

	second, err := Build(opts)
	if err != nil {
		t.Fatalf("second build: %v", err)
	}
	if !second.Stats.CacheHit {
		t.Fatal("second build did not use the build cache")
	}
	if second.Stats.Parsed != 0 || second.Stats.Reused != 2 {
		t.Fatalf("second build parsed/reused = %d/%d, want 0/2", second.Stats.Parsed, second.Stats.Reused)
	}
	if second.CSS != first.CSS {
		t.Fatal("cached output differs from fresh output")
	}
}

// TestIncrementalRebuildParsesOnlyTheEditedFile is the core Phase 8 promise.
func TestIncrementalRebuildParsesOnlyTheEditedFile(t *testing.T) {
	root := t.TempDir()
	for i := 0; i < 8; i++ {
		writeFile(t, root, fmt.Sprintf("page%d.html", i), `<div class="m-1"></div>`)
	}
	writeFile(t, root, "page3.html", `<div class="w-4"></div>`)

	opts := Options{Root: root, Content: []string{"./**/*.html"}, Base: false}
	if _, err := Build(opts); err != nil {
		t.Fatal(err)
	}

	writeFile(t, root, "page3.html", `<div class="w-6"></div>`)
	res, err := Build(opts)
	if err != nil {
		t.Fatal(err)
	}

	if res.Stats.CacheHit {
		t.Fatal("edited project should miss the whole-build cache")
	}
	if res.Stats.Parsed != 1 || res.Stats.Reused != 7 {
		t.Fatalf("parsed/reused = %d/%d, want 1/7", res.Stats.Parsed, res.Stats.Reused)
	}
	if !strings.Contains(res.CSS, ".w-6 { width: 1.5rem; }") {
		t.Fatalf("edited file's new utility missing:\n%s", res.CSS)
	}
	if strings.Contains(res.CSS, ".w-4") {
		t.Fatalf("edited file's old utility is still emitted:\n%s", res.CSS)
	}
	if !strings.Contains(res.CSS, ".m-1 { margin: 0.25rem; }") {
		t.Fatalf("reused files' utilities were lost:\n%s", res.CSS)
	}
}

// TestIncrementalDetectsSameSizeEdit guards the mtime half of the identity
// check: a same-size rewrite must still be re-parsed.
func TestIncrementalDetectsSameSizeEdit(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "a.html", `<div class="p-4"></div>`)
	opts := Options{Root: root, Content: []string{"./**/*.html"}, Base: false}
	if _, err := Build(opts); err != nil {
		t.Fatal(err)
	}

	writeFile(t, root, "a.html", `<div class="m-2"></div>`)
	res, err := Build(opts)
	if err != nil {
		t.Fatal(err)
	}
	if res.Stats.Parsed != 1 {
		t.Fatalf("parsed = %d, want 1", res.Stats.Parsed)
	}
	if !strings.Contains(res.CSS, ".m-2 { margin: 0.5rem; }") {
		t.Fatalf("new utility missing:\n%s", res.CSS)
	}
}

// TestIncrementalManifestPrunesDeletedFiles keeps the manifest proportional to
// the project instead of growing forever.
func TestIncrementalManifestPrunesDeletedFiles(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "a.html", `<div class="p-4"></div>`)
	writeFile(t, root, "b.html", `<div class="m-2"></div>`)
	opts := Options{Root: root, Content: []string{"./**/*.html"}, Base: false}
	if _, err := Build(opts); err != nil {
		t.Fatal(err)
	}
	if len(readManifest(t, root).Files) != 2 {
		t.Fatalf("manifest holds %d files, want 2", len(readManifest(t, root).Files))
	}

	if err := os.Remove(filepath.Join(root, "b.html")); err != nil {
		t.Fatal(err)
	}
	if _, err := Build(opts); err != nil {
		t.Fatal(err)
	}
	files := readManifest(t, root).Files
	if len(files) != 1 {
		t.Fatalf("manifest holds %d files after a deletion, want 1", len(files))
	}
	if _, ok := files["a.html"]; !ok {
		t.Fatalf("remaining manifest entries = %v", files)
	}
}

func TestManifestIsIgnoredWhenCorrupt(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "a.html", `<div class="p-4"></div>`)
	if err := os.MkdirAll(filepath.Join(root, ".flexui"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(manifestPath(root), []byte("{oops"), 0o644); err != nil {
		t.Fatal(err)
	}

	res, err := Build(Options{Root: root, Content: []string{"./**/*.html"}, Base: false})
	if err != nil {
		t.Fatalf("build with a corrupt manifest: %v", err)
	}
	if res.Stats.Parsed != 1 {
		t.Fatalf("parsed = %d, want 1", res.Stats.Parsed)
	}
}

func TestNoCacheSkipsManifest(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "a.html", `<div class="p-4"></div>`)
	res, err := Build(Options{Root: root, Content: []string{"./**/*.html"}, Base: false, NoCache: true})
	if err != nil {
		t.Fatal(err)
	}
	if res.Stats.Parsed != 1 || res.Stats.Reused != 0 {
		t.Fatalf("parsed/reused = %d/%d, want 1/0", res.Stats.Parsed, res.Stats.Reused)
	}
	if _, err := os.Stat(manifestPath(root)); !os.IsNotExist(err) {
		t.Fatalf("--no-cache wrote a manifest: %v", err)
	}
}

// TestManifestIsValidJSON documents the on-disk format other tools may read.
func TestManifestIsValidJSON(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "a.html", `<div class="p-4 flex"></div>`)
	if _, err := Build(Options{Root: root, Content: []string{"./**/*.html"}, Base: false}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(manifestPath(root))
	if err != nil {
		t.Fatal(err)
	}
	var decoded struct {
		Version int `json:"version"`
		Files   map[string]struct {
			Size   int64    `json:"size"`
			Hash   string   `json:"hash"`
			Tokens []string `json:"tokens"`
		} `json:"files"`
	}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("manifest is not valid JSON: %v", err)
	}
	entry, ok := decoded.Files["a.html"]
	if !ok {
		t.Fatalf("manifest files = %v", decoded.Files)
	}
	if entry.Hash == "" || entry.Size == 0 {
		t.Fatalf("entry = %+v", entry)
	}
	if len(entry.Tokens) != 2 {
		t.Fatalf("tokens = %v", entry.Tokens)
	}
}

// fixtureTree writes a project with the given number of HTML files, each using
// a realistic mix of utilities, for the build benchmarks.
func fixtureTree(tb testing.TB, files int) string {
	tb.Helper()
	root := tb.TempDir()
	classes := []string{
		"flex items-center justify-between gap-4 p-4 rounded-lg border border-slate-200",
		"grid grid-cols-1 md:grid-cols-3 gap-6 lg:gap-8 max-w-7xl mx-auto",
		"text-sm font-medium text-slate-600 hover:text-slate-900 transition-colors",
		"bg-white shadow-sm hover:shadow-md dark:bg-slate-900 dark:text-slate-100",
		"fui-btn fui-btn--primary fui-badge fui-card fui-input",
	}
	for i := 0; i < files; i++ {
		dir := filepath.Join(root, "src", fmt.Sprintf("group%d", i/20))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			tb.Fatal(err)
		}
		markup := fmt.Sprintf("<section class=\"%s\"><div class=\"%s\">card %d</div></section>\n",
			classes[i%len(classes)], classes[(i+1)%len(classes)], i)
		if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("page%d.html", i)), []byte(markup), 0o644); err != nil {
			tb.Fatal(err)
		}
	}
	return root
}

func BenchmarkBuildCold(b *testing.B) {
	root := fixtureTree(b, 200)
	opts := Options{Root: root, Content: []string{"./**/*.html"}, Base: true, NoCache: true}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Build(opts); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBuildWarmCache(b *testing.B) {
	root := fixtureTree(b, 200)
	opts := Options{Root: root, Content: []string{"./**/*.html"}, Base: true}
	if _, err := Build(opts); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Build(opts); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBuildIncremental reproduces a watch rebuild: one file changed, the
// other 199 reused from the manifest.
func BenchmarkBuildIncremental(b *testing.B) {
	root := fixtureTree(b, 200)
	opts := Options{Root: root, Content: []string{"./**/*.html"}, Base: true}
	if _, err := Build(opts); err != nil {
		b.Fatal(err)
	}
	target := filepath.Join(root, "src", "group0", "page0.html")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		body := fmt.Sprintf("<section class=\"p-%d flex\">edit %d</section>\n", i%8+1, i)
		if err := os.WriteFile(target, []byte(body), 0o644); err != nil {
			b.Fatal(err)
		}
		if _, err := Build(opts); err != nil {
			b.Fatal(err)
		}
	}
}
