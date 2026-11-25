package utilities

import (
	"math"
	"strconv"
	"strings"

	"github.com/Ahaduzzamankhan/flexui/internal/ast"
	"github.com/Ahaduzzamankhan/flexui/internal/parser"
)

// decl builds a declaration.
func decl(prop, value string) ast.Decl { return ast.Decl{Prop: prop, Value: value} }

// on returns a block that targets the utility's own element.
func on(decls ...ast.Decl) ast.Block { return ast.Block{Decls: decls} }

// under returns a block that targets a descendant or pseudo element via a
// selector template containing "&".
func under(target string, decls ...ast.Decl) ast.Block {
	return ast.Block{Target: target, Decls: decls}
}

// arb unwraps an arbitrary value and normalizes underscores to spaces.
func arb(v string) (string, bool) {
	inner, ok := parser.UnwrapArbitrary(v)
	if !ok {
		return "", false
	}
	out := parser.NormalizeArbitrary(inner)
	if out == "" {
		return "", false
	}
	return out, true
}

// spacing resolves a spacing-scale token or arbitrary length.
func (c *Context) spacing(v string) (string, bool) {
	if s, ok := c.Theme.Spacing[v]; ok {
		return s, true
	}
	return arb(v)
}

// fraction converts "1/2" into "50%".
func fraction(v string) (string, bool) {
	i := strings.IndexByte(v, '/')
	if i <= 0 || i == len(v)-1 {
		return "", false
	}
	num, err1 := strconv.Atoi(v[:i])
	den, err2 := strconv.Atoi(v[i+1:])
	if err1 != nil || err2 != nil || den == 0 {
		return "", false
	}
	pct := math.Round(float64(num)*100/float64(den)*10000) / 10000
	return strconv.FormatFloat(pct, 'f', -1, 64) + "%", true
}

// negate produces the negative form of a length while staying valid CSS.
func negate(v string) string {
	switch v {
	case "0", "0px", "0rem", "0em", "":
		return v
	}
	if strings.HasPrefix(v, "-") {
		return strings.TrimPrefix(v, "-")
	}
	if isSimpleLength(v) {
		return "-" + v
	}
	return "calc(" + v + " * -1)"
}

func isSimpleLength(v string) bool {
	i := 0
	if i < len(v) && (v[i] == '+' || v[i] == '-') {
		i++
	}
	digits := 0
	for i < len(v) && v[i] >= '0' && v[i] <= '9' {
		i++
		digits++
	}
	if i < len(v) && v[i] == '.' {
		i++
		for i < len(v) && v[i] >= '0' && v[i] <= '9' {
			i++
			digits++
		}
	}
	if digits == 0 {
		return false
	}
	for ; i < len(v); i++ {
		c := v[i]
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '%') {
			return false
		}
	}
	return true
}

// applyNeg applies a negative modifier when neg is set.
func applyNeg(neg bool, v string) string {
	if neg {
		return negate(v)
	}
	return v
}

// colorValue resolves a colour token including an optional "/opacity"
// modifier, using color-mix so the result stays a real CSS colour.
func (c *Context) colorValue(token string) (string, bool) {
	base, opacity := token, ""
	if i := strings.LastIndexByte(token, '/'); i > 0 {
		base, opacity = token[:i], token[i+1:]
	}
	var value string
	if strings.HasPrefix(base, "[") {
		v, ok := arb(base)
		if !ok {
			return "", false
		}
		value = v
	} else if v, ok := c.Theme.SemanticVar(base); ok {
		// Semantic roles such as surface, fg-muted, line or accent resolve to
		// custom properties so a theme or dark mode can switch them at runtime.
		value = v
	} else {
		v, ok := c.Theme.Color(base)
		if !ok {
			return "", false
		}
		value = v
	}
	if opacity == "" {
		return value, true
	}
	return mixOpacity(value, opacity), true
}

func mixOpacity(color, opacity string) string {
	pct := opacity
	if v, err := strconv.ParseFloat(opacity, 64); err == nil {
		pct = strconv.FormatFloat(v, 'f', -1, 64) + "%"
	} else if inner, ok := arb("[" + opacity + "]"); ok {
		pct = inner
	}
	return "color-mix(in srgb, " + color + " " + pct + ", transparent)"
}

// numValue resolves a bare number (for opacity, z-index and similar) or an
// arbitrary value.
func numValue(v string) (string, bool) {
	if v != "" {
		ok := true
		for i := 0; i < len(v); i++ {
			c := v[i]
			if !(c >= '0' && c <= '9' || c == '.') {
				ok = false
				break
			}
		}
		if ok {
			return v, true
		}
	}
	return arb(v)
}

// percentValue resolves a percentage token such as "50" into "50%".
func percentValue(v string) (string, bool) {
	if v != "" {
		ok := true
		for i := 0; i < len(v); i++ {
			if !(v[i] >= '0' && v[i] <= '9' || v[i] == '.') {
				ok = false
				break
			}
		}
		if ok {
			return v + "%", true
		}
	}
	if inner, ok := arb(v); ok {
		return inner, true
	}
	return "", false
}
