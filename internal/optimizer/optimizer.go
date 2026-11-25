package optimizer

import (
	"bytes"

	"github.com/Ahaduzzamankhan/flexui/internal/ast"
	"github.com/Ahaduzzamankhan/flexui/internal/generator"
)

// Rules removes empty rules, duplicate rules and duplicate declarations.
func Rules(rules []ast.Rule) []ast.Rule {
	out := make([]ast.Rule, 0, len(rules))
	seen := make(map[string]bool, len(rules))
	for _, r := range rules {
		r = dedupeDecls(r)
		if r.Empty() {
			continue
		}
		key := ruleKey(r)
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, r)
	}
	return out
}

// dedupeDecls drops repeated declarations inside a block. When the same
// property appears more than once the last declaration wins, which matches CSS
// cascade semantics exactly, so removing the earlier ones cannot change meaning.
func dedupeDecls(r ast.Rule) ast.Rule {
	blocks := make([]ast.Block, 0, len(r.Blocks))
	for _, b := range r.Blocks {
		if len(b.Decls) == 0 {
			continue
		}
		last := make(map[string]int, len(b.Decls))
		for i, decl := range b.Decls {
			last[decl.Prop] = i
		}
		decls := make([]ast.Decl, 0, len(b.Decls))
		for i, decl := range b.Decls {
			if last[decl.Prop] != i {
				continue
			}
			decls = append(decls, decl)
		}
		b.Decls = decls
		blocks = append(blocks, b)
	}
	r.Blocks = blocks
	return r
}

// ruleKey is the identity of the CSS a rule produces.
func ruleKey(r ast.Rule) string {
	var b bytes.Buffer
	b.WriteString(r.AtRuleKey())
	b.WriteByte('\x01')
	for _, block := range r.Blocks {
		b.WriteString(generator.Selector(&r, block.Target))
		b.WriteByte('\x02')
		for _, decl := range block.Decls {
			b.WriteString(decl.Prop)
			b.WriteByte(':')
			b.WriteString(decl.Value)
			b.WriteByte(';')
		}
		b.WriteByte('\x03')
	}
	return b.String()
}

// Minify compacts third-party CSS such as a user's passthrough stylesheet.
//
// It is conservative: comments are removed and insignificant whitespace is
// collapsed, but descendant combinators are preserved because a single space
// can change which element a selector matches. FlexUI's own output is emitted
// compactly by the generator, so minification of generated CSS never depends on
// this function.
func Minify(css []byte) []byte {
	out := make([]byte, 0, len(css))
	inComment := false
	inBlock := false
	var quote byte
	lastSpace := false
	for i := 0; i < len(css); i++ {
		c := css[i]
		if inComment {
			if c == '*' && i+1 < len(css) && css[i+1] == '/' {
				inComment = false
				i++
			}
			continue
		}
		if quote != 0 {
			out = append(out, c)
			if c == '\\' && i+1 < len(css) {
				i++
				out = append(out, css[i])
				continue
			}
			if c == quote {
				quote = 0
			}
			continue
		}
		switch c {
		case '/':
			if i+1 < len(css) && css[i+1] == '*' {
				inComment = true
				i++
				continue
			}
			out = append(out, c)
			lastSpace = false
		case '\'', '"':
			quote = c
			out = append(out, c)
			lastSpace = false
		case ' ', '\t', '\n', '\r', '\f':
			if !lastSpace {
				out = append(out, ' ')
				lastSpace = true
			}
		case ',':
			out = trimSpaceEnd(out)
			out = append(out, ',')
			lastSpace = false
			i = skipSpaces(css, i)
		case '{':
			out = trimSpaceEnd(out)
			out = append(out, '{')
			inBlock = true
			lastSpace = false
			i = skipSpaces(css, i)
		case '}':
			out = trimSpaceEnd(out)
			// A semicolon directly before a closing brace is optional.
			for len(out) > 0 && out[len(out)-1] == ';' {
				out = out[:len(out)-1]
			}
			out = append(out, '}')
			inBlock = false
			lastSpace = false
			i = skipSpaces(css, i)
		case ';':
			out = trimSpaceEnd(out)
			out = append(out, ';')
			lastSpace = false
			i = skipSpaces(css, i)
		case ':':
			if inBlock {
				out = trimSpaceEnd(out)
				out = append(out, ':')
				lastSpace = false
				i = skipSpaces(css, i)
			} else {
				out = append(out, c)
				lastSpace = false
			}
		default:
			out = append(out, c)
			lastSpace = false
		}
	}
	return bytes.TrimSpace(out)
}

func isSpaceByte(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f'
}

func skipSpaces(b []byte, i int) int {
	for i+1 < len(b) && isSpaceByte(b[i+1]) {
		i++
	}
	return i
}

func trimSpaceEnd(b []byte) []byte {
	for len(b) > 0 && b[len(b)-1] == ' ' {
		b = b[:len(b)-1]
	}
	return b
}
