// Package scanner discovers the source files that may contain FlexUI classes.
package scanner

import "strings"

// Match reports whether a slash-separated relative path matches a glob
// pattern. Supported syntax: "**" (any number of path segments), "*" (any run
// of characters inside a segment), "?" (one character) and "{a,b}" alternates.
func Match(pattern, path string) bool {
	pat := strings.Split(normalize(pattern), "/")
	name := strings.Split(normalize(path), "/")
	return matchSegments(pat, name)
}

func normalize(p string) string {
	p = strings.ReplaceAll(p, "\\", "/")
	for strings.HasPrefix(p, "./") {
		p = p[2:]
	}
	p = strings.TrimPrefix(p, "/")
	if p == "" {
		return ""
	}
	return p
}

func matchSegments(pat, name []string) bool {
	for len(pat) > 0 {
		if pat[0] == "**" {
			if len(pat) == 1 {
				return true
			}
			for i := 0; i <= len(name); i++ {
				if matchSegments(pat[1:], name[i:]) {
					return true
				}
			}
			return false
		}
		if len(name) == 0 {
			return false
		}
		if !matchSegment(pat[0], name[0]) {
			return false
		}
		pat, name = pat[1:], name[1:]
	}
	return len(name) == 0
}

func matchSegment(pattern, segment string) bool {
	for _, alt := range expandBraces(pattern) {
		if matchFlat(alt, segment) {
			return true
		}
	}
	return false
}

// expandBraces turns "*.{ts,tsx}" into ["*.ts", "*.tsx"].
func expandBraces(p string) []string {
	open := strings.IndexByte(p, '{')
	if open < 0 {
		return []string{p}
	}
	closeIdx := strings.IndexByte(p[open:], '}')
	if closeIdx < 0 {
		return []string{p}
	}
	closeIdx += open
	prefix, body, suffix := p[:open], p[open+1:closeIdx], p[closeIdx+1:]
	var out []string
	for _, part := range strings.Split(body, ",") {
		for _, rest := range expandBraces(suffix) {
			out = append(out, prefix+part+rest)
		}
	}
	return out
}

func matchFlat(pattern, s string) bool {
	pi, si := 0, 0
	star, mark := -1, 0
	for si < len(s) {
		switch {
		case pi < len(pattern) && (pattern[pi] == '?' || pattern[pi] == s[si]):
			pi++
			si++
		case pi < len(pattern) && pattern[pi] == '*':
			star = pi
			mark = si
			pi++
		case star >= 0:
			pi = star + 1
			mark++
			si = mark
		default:
			return false
		}
	}
	for pi < len(pattern) && pattern[pi] == '*' {
		pi++
	}
	return pi == len(pattern)
}
