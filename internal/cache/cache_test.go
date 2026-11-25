package cache

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPutGetClear(t *testing.T) {
	root := t.TempDir()
	c := Open(root)
	key := Key([]byte("inputs"))
	if _, _, ok := c.Get(key); ok {
		t.Fatal("empty cache reported a hit")
	}
	if err := c.Put(key, []byte(".a{}"), []byte(`{"rules":1}`)); err != nil {
		t.Fatal(err)
	}
	data, meta, ok := c.Get(key)
	if !ok {
		t.Fatal("cache entry not found")
	}
	if string(data) != ".a{}" || string(meta) != `{"rules":1}` {
		t.Fatalf("data=%q meta=%q", data, meta)
	}
	if err := c.Clear(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, ".flexui")); !os.IsNotExist(err) {
		t.Fatalf(".flexui still exists: %v", err)
	}
}

func TestKeyDependsOnAllParts(t *testing.T) {
	a := Key([]byte("a"), []byte("b"))
	b := Key([]byte("a"), []byte("c"))
	if a == b {
		t.Fatal("different inputs produced the same key")
	}
	if a != Key([]byte("a"), []byte("b")) {
		t.Fatal("key is not stable")
	}
}
