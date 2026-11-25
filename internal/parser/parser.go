// Package parser turns a raw class candidate into variants plus a utility
// token. It uses a bracket and quote aware state machine instead of regular
// expressions so that nested arbitrary values, function calls, URLs and
// arbitrary selectors all split correctly.
package parser

import "strings"

// Candidate is a parsed class candidate.
type Candidate struct {
	// Raw is the candidate exactly as found in the source.
	Raw string
	// Important marks declarations that must be emitted !important.
	Important bool
	// Negative marks a leading "-" before the utility.
	Negative bool
	// Variants are the raw variant tokens, outermost first.
	Variants []string
	// Utility is the final segment, for example "bg-red-500".
	Utility string
	// ArbitraryProperty marks a bracketed property/value utility such as
	// "[mask-type:luminance]".
	ArbitraryProperty bool
}

// Parse parses one class candidate. It reports false only when the candidate is
// structurally impossible (unbalanced brackets, empty segments).
func Parse(raw string) (Candidate, bool) {
	c := Candidate{Raw: raw}
	s := raw
	if s == "" {
		return c, false
	}
	if s[0] == '!' {
		c.Important = true
		s = s[1:]
	}
	segments, ok := Split(s)
	if !ok || len(segments) == 0 {
		return c, false
	}
	util := segments[len(segments)-1]
	c.Variants = segments[:len(segments)-1]
	if len(util) > 1 && util[0] == '!' {
		c.Important = true
		util = util[1:]
	}
	if len(util) > 2 && util[len(util)-1] == '!' {
		c.Important = true
		util = util[:len(util)-1]
	}
	if util == "" {
		return c, false
	}
	if isBracketed(util) && hasTopLevelColon(util[1:len(util)-1]) {
		c.ArbitraryProperty = true
		c.Utility = util
		return c, true
	}
	if len(util) > 1 && util[0] == '-' && util[1] != '-' {
		c.Negative = true
		util = util[1:]
	}
	if util == "" {
		return c, false
	}
	c.Utility = util
	return c, true
}

// Split splits s on colons that are not inside brackets, parentheses, braces,
// quotes or escapes.
func Split(s string) ([]string, bool) {
	var segments []string
	start := 0
	depth := 0
	var quote byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\\' && i+1 < len(s) {
			i++
			continue
		}
		if quote != 0 {
			if c == quote {
				quote = 0
			}
			continue
		}
		switch c {
		case '\'', '"':
			quote = c
		case '[', '(', '{':
			depth++
		case ']', ')', '}':
			depth--
			if depth < 0 {
				return nil, false
			}
		case ':':
			if depth == 0 {
				segments = append(segments, s[start:i])
				start = i + 1
			}
		}
	}
	if depth != 0 || quote != 0 {
		return nil, false
	}
	segments = append(segments, s[start:])
	for _, seg := range segments {
		if seg == "" {
			return nil, false
		}
	}
	return segments, true
}

// SplitUtility splits a utility token into its registry prefix and value at the
// first top-level "-". The utility registry performs its own longest-prefix
// matching; this helper is for diagnostics and tests.
func SplitUtility(token string) (prefix, value string) {
	depth := 0
	for i := 0; i < len(token); i++ {
		switch token[i] {
		case '[', '(', '{':
			depth++
		case ']', ')', '}':
			depth--
		case '-':
			if depth == 0 && i > 0 {
				return token[:i], token[i+1:]
			}
		}
	}
	return token, ""
}

// UnwrapArbitrary strips the surrounding brackets and any type hint from an
// arbitrary value such as "[length:12px]" or "[calc(100vh-80px)]".
func UnwrapArbitrary(v string) (string, bool) {
	if !isBracketed(v) {
		return "", false
	}
	inner := v[1 : len(v)-1]
	if i := strings.IndexByte(inner, ':'); i > 0 && isHint(inner[:i]) {
		inner = inner[i+1:]
	}
	return inner, true
}

// NormalizeArbitrary converts FlexUI's underscores into spaces, preserving
// escaped underscores, so "grid-cols-[250px_1fr]" becomes "250px 1fr".
func NormalizeArbitrary(v string) string {
	if !strings.ContainsRune(v, '_') {
		return v
	}
	var b strings.Builder
	b.Grow(len(v))
	for i := 0; i < len(v); i++ {
		switch {
		case v[i] == '\\' && i+1 < len(v) && v[i+1] == '_':
			b.WriteByte('_')
			i++
		case v[i] == '_':
			b.WriteByte(' ')
		default:
			b.WriteByte(v[i])
		}
	}
	return b.String()
}

func isBracketed(v string) bool {
	if len(v) < 2 || v[0] != '[' || v[len(v)-1] != ']' {
		return false
	}
	depth := 0
	var quote byte
	for i := 0; i < len(v); i++ {
		c := v[i]
		if quote != 0 {
			if c == quote {
				quote = 0
			}
			continue
		}
		switch c {
		case '\'', '"':
			quote = c
		case '[':
			depth++
		case ']':
			depth--
			if depth == 0 {
				return i == len(v)-1
			}
		}
	}
	return false
}

func hasTopLevelColon(s string) bool {
	depth := 0
	var quote byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		if quote != 0 {
			if c == quote {
				quote = 0
			}
			continue
		}
		switch c {
		case '\'', '"':
			quote = c
		case '[', '(', '{':
			depth++
		case ']', ')', '}':
			depth--
		case ':':
			if depth == 0 {
				return true
			}
		}
	}
	return false
}

func isHint(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '-') {
			return false
		}
	}
	return true
}
