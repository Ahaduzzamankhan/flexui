// Package cache stores compiled stylesheets keyed by their content hash.
package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
)

// Cache is a content-addressed directory of build outputs.
type Cache struct {
	dir string
}

// Open returns the cache rooted at .flexui/cache inside root.
func Open(root string) *Cache {
	return &Cache{dir: filepath.Join(root, ".flexui", "cache")}
}

// Dir returns the cache directory.
func (c *Cache) Dir() string { return c.dir }

// Hash returns the SHA-256 of one file's contents as hex, used to identify a
// file's content in the incremental manifest and in the build key.
func Hash(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// Key hashes a list of byte chunks into a file-safe key.
func Key(parts ...[]byte) string {
	h := sha256.New()
	for _, p := range parts {
		h.Write(p)
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))[:32]
}

// Get returns a cached stylesheet and its metadata.
func (c *Cache) Get(key string) (data, meta []byte, ok bool) {
	data, err := os.ReadFile(c.path(key, ".css"))
	if err != nil {
		return nil, nil, false
	}
	meta, _ = os.ReadFile(c.path(key, ".json"))
	return data, meta, true
}

// Put stores a stylesheet and its metadata atomically.
func (c *Cache) Put(key string, data, meta []byte) error {
	if err := os.MkdirAll(c.dir, 0o755); err != nil {
		return err
	}
	if err := writeAtomic(c.path(key, ".css"), data); err != nil {
		return err
	}
	return writeAtomic(c.path(key, ".json"), meta)
}

// Clear removes every cached entry.
func (c *Cache) Clear() error {
	return os.RemoveAll(filepath.Join(filepath.Dir(c.dir), ""))
}

func (c *Cache) path(key, ext string) string {
	return filepath.Join(c.dir, key+ext)
}

func writeAtomic(path string, data []byte) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
