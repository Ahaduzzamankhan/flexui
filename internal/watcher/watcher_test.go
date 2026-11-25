package watcher

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTakeAndDiff(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "a.html")
	if err := os.WriteFile(file, []byte("<div></div>"), 0o644); err != nil {
		t.Fatal(err)
	}
	prev := Take([]string{file})
	if changed := Diff(prev, prev); len(changed) != 0 {
		t.Fatalf("identical snapshots reported %v", changed)
	}

	if err := os.WriteFile(file, []byte("<div class=\"p-4\"></div>"), 0o644); err != nil {
		t.Fatal(err)
	}
	next := Take([]string{file})
	changed := Diff(prev, next)
	if len(changed) != 1 || changed[0] != file {
		t.Fatalf("changed = %v", changed)
	}
}

func TestDiffDetectsNewAndRemovedFiles(t *testing.T) {
	root := t.TempDir()
	existing := filepath.Join(root, "a.html")
	added := filepath.Join(root, "b.html")
	if err := os.WriteFile(existing, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	prev := Take([]string{existing})
	if err := os.WriteFile(added, []byte("y"), 0o644); err != nil {
		t.Fatal(err)
	}
	next := Take([]string{existing, added})
	if got := Diff(prev, next); len(got) != 1 || got[0] != added {
		t.Fatalf("added file not detected: %v", got)
	}
	if got := Diff(next, Take([]string{existing})); len(got) != 1 || got[0] != added {
		t.Fatalf("removed file not detected: %v", got)
	}
}

func TestMissingFileIsRecorded(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope.html")
	snap := Take([]string{missing})
	if !snap[missing].Missing {
		t.Fatal("missing file not recorded")
	}
}
