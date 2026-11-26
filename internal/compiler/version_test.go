package compiler

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestVersionMatchesVersionFile keeps the injected release version honest: the
// default in compiler.Version must equal version.json at all times, so a plain
// `go build` produces the same version string as a release binary.
func TestVersionMatchesVersionFile(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "version.json"))
	if err != nil {
		t.Fatalf("read version.json: %v", err)
	}
	var info struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(data, &info); err != nil {
		t.Fatalf("parse version.json: %v", err)
	}
	if info.Version != Version {
		t.Fatalf("version.json says %q but compiler.Version is %q; run `bun run version:sync`", info.Version, Version)
	}
	if !strings.Contains(Banner, Version) {
		t.Fatalf("banner %q does not mention version %q", Banner, Version)
	}
}
