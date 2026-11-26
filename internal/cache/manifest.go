package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// ManifestVersion guards the on-disk shape. Bumping it invalidates every
// manifest written by an older compiler, so a changed extraction semantics or
// manifest layout can never be read back as if it were current.
const ManifestVersion = 1

// FileEntry is everything the compiler needs to skip re-reading one file: its
// identity (size and modification time), the content hash used in the build
// key, and the candidates extracted from it.
type FileEntry struct {
	Size    int64    `json:"size"`
	ModTime int64    `json:"modTime"`
	Hash    string   `json:"hash"`
	Tokens  []string `json:"tokens"`
}

// Manifest is the persisted per-file state of one project.
//
// It is a single file on purpose. Measured on Windows, loading 200 separate
// per-file cache entries took ~26 ms while extracting all 200 files took
// 0.66 ms, so a file-per-file cache made rebuilds slower. One JSON document
// keeps the reuse win without the small-file I/O cost.
type Manifest struct {
	Version int                  `json:"version"`
	Files   map[string]FileEntry `json:"files"`
}

// NewManifest returns an empty manifest at the current version.
func NewManifest() *Manifest {
	return &Manifest{Version: ManifestVersion, Files: make(map[string]FileEntry)}
}

// ManifestStore reads and writes .flexui/manifest.json.
type ManifestStore struct {
	path string
}

// OpenManifest returns the manifest store for a project root.
func OpenManifest(root string) *ManifestStore {
	return &ManifestStore{path: filepath.Join(root, ".flexui", "manifest.json")}
}

// Path returns the manifest file path.
func (s *ManifestStore) Path() string { return s.path }

// Load returns the stored manifest, or nil when it is absent, corrupt, or was
// written by a different manifest version. A nil result simply means every file
// is read and parsed again.
func (s *ManifestStore) Load() *Manifest {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return nil
	}
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil
	}
	if manifest.Version != ManifestVersion || manifest.Files == nil {
		return nil
	}
	return &manifest
}

// Save writes the manifest atomically. Only files in the current scan are
// stored, so removed files are pruned simply by saving.
func (s *ManifestStore) Save(manifest *Manifest) error {
	if manifest == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	data, err := json.Marshal(manifest)
	if err != nil {
		return err
	}
	return writeAtomic(s.path, data)
}

// Remove deletes the manifest. A missing file is not an error.
func (s *ManifestStore) Remove() error {
	err := os.Remove(s.path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
