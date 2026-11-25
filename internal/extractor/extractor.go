// Package extractor pulls raw class candidates out of source files.
//
// Extraction is structural: markup is read with a tag/attribute state machine
// and scripts with a string/comment/template-literal state machine. Regular
// expressions are never used to locate markup structure; the only pattern-like
// check is a per-character shape test on an already-split token.
package extractor

import (
	"path/filepath"
	"strings"
)

// Kind classifies how a file should be read.
type Kind int

const (
	KindMarkup Kind = iota
	KindScript
)

// KindOf returns the extraction strategy for a file extension.
func KindOf(path string) Kind {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".js", ".mjs", ".cjs", ".ts", ".mts", ".cts":
		return KindScript
	default:
		return KindMarkup
	}
}

// Extract returns the class candidates found in src, in source order.
// Duplicates are left in; the compiler de-duplicates globally.
func Extract(path string, src []byte) []string {
	var out []string
	add := func(s string) { out = append(out, s) }
	text := string(src)
	if KindOf(path) == KindScript {
		scanScript(text, add)
	} else {
		scanMarkup(text, add)
	}
	return out
}

// Tokens splits a class attribute value into candidate tokens.
func Tokens(value string, add func(string)) {
	for _, f := range strings.Fields(value) {
		f = strings.Trim(f, "'\"`,;")
		if ShapeOK(f) {
			add(f)
		}
	}
}

// ShapeOK reports whether a token could plausibly be a FlexUI class. The test
// is cheap and permissive: the resolver decides whether the candidate is real.
// Its job is to discard prose, identifiers and JS syntax so that arbitrary
// strings in application code do not bloat the candidate set.
func ShapeOK(t string) bool {
	if len(t) == 0 || len(t) > 200 {
		return false
	}
	depth := 0
	for i := 0; i < len(t); i++ {
		c := t[i]
		switch c {
		case '[':
			depth++
			continue
		case ']':
			depth--
			if depth < 0 {
				return false
			}
			continue
		}
		if depth > 0 {
			switch c {
			case '\n', '\r', '\t', ' ':
				return false
			}
			continue
		}
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9':
		case strings.IndexByte("-_:/.%#@!+*&~>", c) >= 0:
		default:
			return false
		}
	}
	if depth != 0 {
		return false
	}
	if len(t) == 1 {
		return false
	}
	c := t[0]
	switch {
	case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z':
		return true
	case c == '[' || c == '-' || c == '!' || c == '@' || c == '*': // ':' is deliberately rejected as a leading char
		return c != ':'
	}
	return false
}
