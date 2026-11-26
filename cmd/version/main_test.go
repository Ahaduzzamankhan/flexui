package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeVersionFile(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "version.json")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestReadValidFile(t *testing.T) {
	path := writeVersionFile(t, `{"name":"flexui","version":"2.0.0-alpha.1","channel":"alpha"}`)
	info, err := Read(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if info.Name != "flexui" || info.Version != "2.0.0-alpha.1" || info.Channel != "alpha" {
		t.Fatalf("info = %+v", info)
	}
	if info.Tag() != "v2.0.0-alpha.1" {
		t.Fatalf("tag = %q", info.Tag())
	}
}

func TestReadRejectsInvalidFiles(t *testing.T) {
	cases := map[string]string{
		"corrupt":       `{`,
		"no name":       `{"version":"1.0.0"}`,
		"no version":    `{"name":"flexui"}`,
		"not semver":    `{"name":"flexui","version":"alpha"}`,
		"blank version": `{"name":"flexui","version":"  "}`,
		"bad channel":   `{"name":"flexui","version":"1.0.0","channel":"a b"}`,
	}
	for name, body := range cases {
		if _, err := Read(writeVersionFile(t, body)); err == nil {
			t.Fatalf("%s should have failed", name)
		}
	}
}

func TestReadMissingFile(t *testing.T) {
	if _, err := Read(filepath.Join(t.TempDir(), "absent.json")); err == nil {
		t.Fatal("missing file should fail")
	}
}

func TestField(t *testing.T) {
	info := Info{Name: "flexui", Version: "1.2.3", Channel: "stable"}
	for field, want := range map[string]string{
		"name": "flexui", "version": "1.2.3", "channel": "stable", "tag": "v1.2.3",
	} {
		got, err := info.Field(field)
		if err != nil || got != want {
			t.Fatalf("Field(%q) = (%q, %v), want %q", field, got, err, want)
		}
	}
	if _, err := info.Field("bogus"); err == nil {
		t.Fatal("unknown field should fail")
	}
}
