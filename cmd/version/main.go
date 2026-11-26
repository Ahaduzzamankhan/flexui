// Command version resolves the repository's version.json for build scripts.
//
// The version lives in exactly one file. Go build flags and shell scripts read
// it through this command, and the TypeScript scripts (scripts/version.ts)
// resolve the same file, so a release cannot ship mismatched versions.
//
// Usage:
//
//	go run ./cmd/version                  # 2.0.0-alpha.1
//	go run ./cmd/version -field tag       # v2.0.0-alpha.1
//	go run ./cmd/version -field name      # flexui
//	go run ./cmd/version -field channel   # alpha
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Info is the parsed contents of version.json.
type Info struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Channel string `json:"channel"`
}

// Tag returns the git tag for this version.
func (i Info) Tag() string { return "v" + i.Version }

var versionPattern = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.]+)?(\+[0-9A-Za-z.]+)?$`)

// Read loads and validates a version.json file.
func Read(path string) (Info, error) {
	var info Info
	data, err := os.ReadFile(path)
	if err != nil {
		return info, fmt.Errorf("read %s: %w", path, err)
	}
	if err := json.Unmarshal(data, &info); err != nil {
		return info, fmt.Errorf("parse %s: %w", path, err)
	}
	if info.Name == "" {
		return info, errors.New("version.json is missing \"name\"")
	}
	if strings.TrimSpace(info.Version) == "" {
		return info, errors.New("version.json is missing \"version\"")
	}
	if !versionPattern.MatchString(info.Version) {
		return info, fmt.Errorf("version.json version %q is not a semantic version", info.Version)
	}
	if info.Channel != "" && strings.ContainsAny(info.Channel, " \t") {
		return info, fmt.Errorf("version.json channel %q must not contain whitespace", info.Channel)
	}
	return info, nil
}

// Field returns one resolved value.
func (i Info) Field(name string) (string, error) {
	switch name {
	case "name":
		return i.Name, nil
	case "version":
		return i.Version, nil
	case "channel":
		return i.Channel, nil
	case "tag":
		return i.Tag(), nil
	default:
		return "", fmt.Errorf("unknown field %q (expected name, version, channel or tag)", name)
	}
}

func main() {
	file := flag.String("file", "version.json", "path to version.json")
	field := flag.String("field", "version", "name, version, channel or tag")
	flag.Parse()

	info, err := Read(*file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "flexui version:", err)
		os.Exit(1)
	}
	value, err := info.Field(*field)
	if err != nil {
		fmt.Fprintln(os.Stderr, "flexui version:", err)
		os.Exit(1)
	}
	fmt.Println(value)
}
