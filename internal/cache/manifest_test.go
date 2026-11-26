package cache

import (
	"os"
	"testing"
)

func TestManifestRoundTrip(t *testing.T) {
	root := t.TempDir()
	store := OpenManifest(root)

	if store.Load() != nil {
		t.Fatal("expected no manifest in an empty project")
	}

	manifest := NewManifest()
	manifest.Files["src/a.html"] = FileEntry{
		Size:    42,
		ModTime: 1700000000000000000,
		Hash:    "abc123",
		Tokens:  []string{"p-4", "flex"},
	}
	if err := store.Save(manifest); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded := store.Load()
	if loaded == nil {
		t.Fatal("manifest did not load back")
	}
	entry, ok := loaded.Files["src/a.html"]
	if !ok {
		t.Fatal("file entry missing")
	}
	if entry.Size != 42 || entry.Hash != "abc123" || len(entry.Tokens) != 2 {
		t.Fatalf("entry = %+v", entry)
	}
	if entry.Tokens[1] != "flex" {
		t.Fatalf("tokens = %v", entry.Tokens)
	}
}

func TestManifestRejectsWrongVersionAndCorruption(t *testing.T) {
	root := t.TempDir()
	store := OpenManifest(root)
	if err := os.MkdirAll(root+"/.flexui", 0o755); err != nil {
		t.Fatal(err)
	}

	for name, body := range map[string]string{
		"corrupt":   "{not json",
		"version":   `{"version":99,"files":{"a":{"size":1}}}`,
		"no-files":  `{"version":1}`,
		"empty-obj": `{}`,
	} {
		if err := os.WriteFile(store.Path(), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		if got := store.Load(); got != nil {
			t.Fatalf("%s manifest should be ignored, got %+v", name, got)
		}
	}
}

func TestManifestRemoveIsIdempotent(t *testing.T) {
	store := OpenManifest(t.TempDir())
	if err := store.Remove(); err != nil {
		t.Fatalf("remove on a missing manifest: %v", err)
	}
	if err := store.Save(NewManifest()); err != nil {
		t.Fatal(err)
	}
	if err := store.Remove(); err != nil {
		t.Fatalf("remove: %v", err)
	}
	if store.Load() != nil {
		t.Fatal("manifest still loads after remove")
	}
}
