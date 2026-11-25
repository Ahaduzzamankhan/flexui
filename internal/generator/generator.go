// Package generator turns resolved rules into CSS text.
//
// Generation is the only stage that produces CSS, and it sorts rules by a fully
// specified key so the same input always produces byte-identical output. Rules
// are grouped into the cascade layers flexui.theme, flexui.base,
// flexui.components and flexui.utilities; conditional group rules nest in the
// order the variants were written.
package generator

import (
	"bytes"
	"sort"
	"strings"

	"github.com/Ahaduzzamankhan/flexui/internal/ast"
	"github.com/Ahaduzzamankhan/flexui/internal/themes"
)

// Options controls output.
type Options struct {
	// Minify emits compact CSS with no insignificant whitespace.
	Minify bool
	// Base includes the preflight reset layer.
	Base bool
	// Banner is written as a leading comment when non-empty.
	Banner string
	// Passthrough CSS is appended verbatim after the generated layers.
	Passthrough string
	// Theme emits the design-token custom properties. Nil disables them.
	Theme *themes.Theme
	// DarkMode is "class", "media" or "both" and controls how the dark token
	// block is emitted.
	DarkMode string
	// Fonts are @font-face declarations emitted into the base layer.
	Fonts []themes.FontFace
	// ReducedMotion neutralises every animation under
	// prefers-reduced-motion: reduce when animations were generated.
	ReducedMotion bool
}

// Generate renders the stylesheet.
func Generate(rules []ast.Rule, opt Options) []byte {
	sorted := make([]ast.Rule, len(rules))
	copy(sorted, rules)
	Sort(sorted)

	var base, components, utilities []ast.Rule
	for _, r := range sorted {
		switch r.LayerName() {
		case ast.LayerComponents:
			components = append(components, r)
		default:
			utilities = append(utilities, r)
		}
	}
	_ = base

	keyframes := collectKeyframes(sorted)
	emitTheme := opt.Theme != nil && (len(sorted) > 0 || opt.Base)

	w := &writer{minify: opt.Minify}
	if opt.Banner != "" {
		w.comment(opt.Banner)
	}

	layers := make([]string, 0, 4)
	if emitTheme {
		layers = append(layers, "flexui.theme")
	}
	if opt.Base {
		layers = append(layers, "flexui.base")
	}
	if len(components) > 0 {
		layers = append(layers, "flexui.components")
	}
	if len(utilities) > 0 || len(keyframes) > 0 {
		layers = append(layers, "flexui.utilities")
	}
	if len(layers) > 0 {
		w.line("@layer " + strings.Join(layers, ", ") + ";")
	}

	if emitTheme {
		emitThemeLayer(w, opt)
	}
	if opt.Base {
		w.open("@layer flexui.base")
		if len(opt.Fonts) > 0 {
			emitFonts(w, opt.Fonts)
		}
		for _, br := range baseRules {
			w.rule(br.selector, br.decls)
		}
		w.close()
	}
	if len(components) > 0 {
		w.open("@layer flexui.components")
		emitRules(w, components)
		w.close()
	}
	if len(utilities) > 0 || len(keyframes) > 0 {
		w.open("@layer flexui.utilities")
		if len(keyframes) > 0 {
			emitKeyframes(w, keyframes)
			if opt.ReducedMotion {
				emitReducedMotion(w)
			}
		}
		emitRules(w, utilities)
		w.close()
	}
	if opt.Passthrough != "" {
		w.rawBlock(opt.Passthrough)
	}
	out := w.buf.Bytes()
	if opt.Minify && len(out) > 0 {
		out = append(out, '\n')
	}
	result := make([]byte, len(out))
	copy(result, out)
	return result
}

// emitThemeLayer writes the design tokens, the dark token override and the
// semantic roles every component reads.
func emitThemeLayer(w *writer, opt Options) {
	t := opt.Theme
	w.open("@layer flexui.theme")
	w.rule(":root", varsToDecls(t.TokenVars()))
	dark := varsToDecls(t.DarkTokenVars())
	if len(dark) > 0 {
		switch opt.DarkMode {
		case "media":
			w.open("@media (prefers-color-scheme: dark)")
			w.rule(":root", dark)
			w.close()
		case "both":
			w.rule(".dark", dark)
			w.open("@media (prefers-color-scheme: dark)")
			w.rule(":root:not(.light)", dark)
			w.close()
		default:
			w.rule(".dark", dark)
		}
	}
	w.close()
}

// varsToDecls converts custom properties into declarations.
func varsToDecls(vars []themes.CSSVar) []ast.Decl {
	decls := make([]ast.Decl, 0, len(vars))
	for _, v := range vars {
		decls = append(decls, ast.Decl{Prop: v.Name, Value: v.Value})
	}
	return decls
}

// emitFonts writes @font-face blocks for configured local fonts.
func emitFonts(w *writer, fonts []themes.FontFace) {
	for _, f := range fonts {
		if f.Family == "" || f.Src == "" {
			continue
		}
		decls := []ast.Decl{
			{Prop: "font-family", Value: quoteFamily(f.Family)},
			{Prop: "src", Value: f.Src},
		}
		if f.Weight != "" {
			decls = append(decls, ast.Decl{Prop: "font-weight", Value: f.Weight})
		}
		if f.Style != "" {
			decls = append(decls, ast.Decl{Prop: "font-style", Value: f.Style})
		}
		display := f.Display
		if display == "" {
			display = "swap"
		}
		decls = append(decls, ast.Decl{Prop: "font-display", Value: display})
		w.rule("@font-face", decls)
	}
}

// quoteFamily quotes a family name unless it is already a single identifier.
func quoteFamily(name string) string {
	if strings.HasPrefix(name, "\"") || strings.HasPrefix(name, "'") {
		return name
	}
	return "\"" + name + "\""
}

// emitReducedMotion collapses animation durations to a single frame. Because the
// animation utilities read custom properties, one override covers every
// animation in the build without duplicating rules.
func emitReducedMotion(w *writer) {
	w.open("@media (prefers-reduced-motion: reduce)")
	w.rule(":root", []ast.Decl{
		{Prop: "--flexui-animation-duration", Value: "1ms"},
		{Prop: "--flexui-animation-delay", Value: "0ms"},
		{Prop: "--flexui-animation-iteration", Value: "1"},
	})
	w.close()
}

// collectKeyframes gathers the keyframes required by the emitted rules, keeping
// the first definition of each name and returning them sorted.
func collectKeyframes(rules []ast.Rule) []ast.Keyframe {
	seen := make(map[string]ast.Keyframe, 8)
	for _, r := range rules {
		for _, kf := range r.Keyframes {
			if kf.Name == "" || len(kf.Steps) == 0 {
				continue
			}
			if _, exists := seen[kf.Name]; exists {
				continue
			}
			seen[kf.Name] = kf
		}
	}
	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	out := make([]ast.Keyframe, 0, len(names))
	for _, name := range names {
		out = append(out, seen[name])
	}
	return out
}

// emitKeyframes writes each keyframes block.
func emitKeyframes(w *writer, keyframes []ast.Keyframe) {
	for _, kf := range keyframes {
		w.open("@keyframes " + kf.Name)
		for _, step := range kf.Steps {
			w.rule(step.Steps, step.Decls)
		}
		w.close()
	}
}

// Sort orders rules deterministically. The key is: layer, group-rule rank,
// group-rule key, utility order, variant rank, variant count, class name,
// declaration signature.
func Sort(rules []ast.Rule) {
	sort.SliceStable(rules, func(i, j int) bool {
		a, b := rules[i], rules[j]
		if la, lb := layerRank(a.LayerName()), layerRank(b.LayerName()); la != lb {
			return la < lb
		}
		if ra, rb := a.AtRuleRank(), b.AtRuleRank(); ra != rb {
			return ra < rb
		}
		if ka, kb := a.AtRuleKey(), b.AtRuleKey(); ka != kb {
			return ka < kb
		}
		if a.UtilityOrder != b.UtilityOrder {
			return a.UtilityOrder < b.UtilityOrder
		}
		if ra, rb := a.VariantRank(), b.VariantRank(); ra != rb {
			return ra < rb
		}
		if len(a.Variants) != len(b.Variants) {
			return len(a.Variants) < len(b.Variants)
		}
		if a.Class != b.Class {
			return a.Class < b.Class
		}
		return a.Signature() < b.Signature()
	})
}

func layerRank(layer string) int {
	switch layer {
	case ast.LayerComponents:
		return 1
	case ast.LayerUtilities, "":
		return 2
	default:
		return 0
	}
}

// emitRules writes rules grouped by their at-rule context.
func emitRules(w *writer, rules []ast.Rule) {
	i := 0
	for i < len(rules) {
		atRules := rules[i].AtRules()
		key := rules[i].AtRuleKey()
		j := i
		for j < len(rules) && rules[j].AtRuleKey() == key {
			j++
		}
		for _, at := range atRules {
			w.open(at)
		}
		for _, r := range rules[i:j] {
			if r.Empty() {
				continue
			}
			for _, block := range r.Blocks {
				if len(block.Decls) == 0 {
					continue
				}
				selector := Selector(&r, block.Target)
				if block.AtRule != "" {
					w.open(block.AtRule)
					w.rule(selector, block.Decls)
					w.close()
					continue
				}
				w.rule(selector, block.Decls)
			}
		}
		for k := len(atRules) - 1; k >= 0; k-- {
			w.close()
		}
		i = j
	}
}

// Selector builds the final selector for a rule and a block target. Variant
// templates are applied in written order, then the escaped class name replaces
// the remaining "&".
func Selector(r *ast.Rule, target string) string {
	expr := "&"
	for _, v := range r.Variants {
		if v.Template == "" {
			continue
		}
		expr = strings.ReplaceAll(v.Template, "&", expr)
	}
	if target == "" {
		target = "&"
	}
	sel := strings.ReplaceAll(target, "&", expr)
	return strings.ReplaceAll(sel, "&", "."+EscapeClass(r.Class))
}

// EscapeClass escapes a class name for use in a CSS selector. It follows the
// CSS syntax rules for identifiers: punctuation is backslash-escaped, leading
// digits use a hexadecimal escape, and non-ASCII bytes are kept as-is.
func EscapeClass(s string) string {
	var b strings.Builder
	b.Grow(len(s) + 8)
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '-' || c == '_':
			b.WriteByte(c)
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z':
			b.WriteByte(c)
		case c >= '0' && c <= '9':
			if i == 0 || (i == 1 && s[0] == '-') {
				b.WriteString(`\3`)
				b.WriteByte(c)
				b.WriteByte(' ')
				continue
			}
			b.WriteByte(c)
		case c >= 0x80:
			b.WriteByte(c)
		case c == 0:
			b.WriteString(`\fffd `)
		default:
			b.WriteByte('\\')
			b.WriteByte(c)
		}
	}
	return b.String()
}

type writer struct {
	buf    bytes.Buffer
	minify bool
	depth  int
}

func (w *writer) indent() {
	if w.minify {
		return
	}
	for i := 0; i < w.depth; i++ {
		w.buf.WriteString("  ")
	}
}

func (w *writer) line(s string) {
	if w.minify {
		w.buf.WriteString(s)
		return
	}
	w.indent()
	w.buf.WriteString(s)
	w.buf.WriteByte('\n')
}

func (w *writer) open(s string) {
	if w.minify {
		w.buf.WriteString(s)
		w.buf.WriteByte('{')
	} else {
		w.line(s + " {")
	}
	w.depth++
}

func (w *writer) close() {
	if w.depth > 0 {
		w.depth--
	}
	if w.minify {
		w.buf.WriteByte('}')
		return
	}
	w.line("}")
}

func (w *writer) rule(selector string, decls []ast.Decl) {
	if len(decls) == 0 {
		return
	}
	if w.minify {
		w.buf.WriteString(selector)
		w.buf.WriteByte('{')
		for _, d := range decls {
			w.buf.WriteString(d.Prop)
			w.buf.WriteByte(':')
			w.buf.WriteString(d.Value)
			w.buf.WriteByte(';')
		}
		w.buf.WriteByte('}')
		return
	}
	w.indent()
	w.buf.WriteString(selector)
	w.buf.WriteString(" { ")
	for _, d := range decls {
		w.buf.WriteString(d.Prop)
		w.buf.WriteString(": ")
		w.buf.WriteString(d.Value)
		w.buf.WriteString("; ")
	}
	w.buf.WriteString("}\n")
}

func (w *writer) comment(s string) {
	for _, line := range strings.Split(s, "\n") {
		if w.minify {
			w.buf.WriteString("/*")
			w.buf.WriteString(strings.ReplaceAll(line, "*/", "* /"))
			w.buf.WriteString("*/")
			continue
		}
		w.buf.WriteString("/* ")
		w.buf.WriteString(strings.ReplaceAll(line, "*/", "* /"))
		w.buf.WriteString(" */\n")
	}
}

func (w *writer) rawBlock(s string) {
	if !w.minify {
		w.buf.WriteByte('\n')
	}
	w.buf.WriteString(s)
	if !strings.HasSuffix(s, "\n") && !w.minify {
		w.buf.WriteByte('\n')
	}
}
