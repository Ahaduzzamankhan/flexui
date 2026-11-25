// Package ast defines FlexUI's intermediate representation.
//
// The pipeline is: source files -> scanner -> extractor -> parser ->
// resolver -> ast (this package) -> generator -> optimizer.
//
// Rules are deliberately structural: they describe which selector transform
// and which conditional group rules apply, and the generator is the only
// place that turns them into CSS text. This keeps generation deterministic,
// makes rules cheap to sort, deduplicate and test, and avoids building CSS by
// concatenating strings during parsing.
package ast

import "strings"

// Cascade layers emitted by the generator.
const (
	LayerBase       = "base"
	LayerComponents = "components"
	LayerUtilities  = "utilities"
)

// Decl is a single CSS declaration.
type Decl struct {
	Prop  string
	Value string
}

// KeyframeStep is one step of a @keyframes block.
type KeyframeStep struct {
	// Steps is the step list, for example "from", "to" or "50%".
	Steps string
	Decls []Decl
}

// Keyframe is a keyframes block required by a rule. Identical keyframes are
// emitted once.
type Keyframe struct {
	Name  string
	Steps []KeyframeStep
}

// Block is a group of declarations that applies to a selector derived from the
// element the utility is placed on. Target is a selector template that uses
// "&" for the element; an empty Target means "&".
type Block struct {
	Target string
	Decls  []Decl
	// AtRule wraps this block in a conditional group rule, for example
	// "@media (min-width: 1024px)". Used by component parts that are
	// responsive without needing a variant on the class.
	AtRule string
}

// VariantKind classifies how a variant transforms a rule.
type VariantKind int

const (
	// VariantSelector rewrites the element selector, for example "&:hover".
	VariantSelector VariantKind = iota
	// VariantAncestor adds an ancestor or sibling condition, for example
	// ".group:hover &" or ".peer:checked ~ &".
	VariantAncestor
	// VariantAtRule wraps the rule in a conditional group rule.
	VariantAtRule
)

// AppliedVariant is a variant after resolution. Exactly one of Template and
// AtRule is meaningful, depending on Kind.
type AppliedVariant struct {
	Name     string
	Kind     VariantKind
	Template string
	AtRule   string
	Rank     int
}

// Rule is a resolved utility occurrence: one class candidate and the CSS it
// produces. Rules are the unit of sorting, deduplication and generation.
type Rule struct {
	// Class is the raw candidate, for example "md:hover:bg-red-500".
	Class string
	// Utility is the registry key that produced the declarations.
	Utility string
	// UtilityOrder is the stable cascade position of the utility.
	UtilityOrder int
	// Variants are ordered outermost first, exactly as written in the class.
	Variants []AppliedVariant
	// Blocks are the declaration groups produced by the utility.
	Blocks []Block
	// Layer is the cascade layer the rule belongs to. Empty means utilities.
	Layer string
	// Keyframes are the animation blocks this rule needs.
	Keyframes []Keyframe
}

// LayerName returns the rule's layer, defaulting to utilities.
func (r *Rule) LayerName() string {
	if r.Layer == "" {
		return LayerUtilities
	}
	return r.Layer
}

// AtRules returns the conditional group rules wrapping this rule, outermost
// first.
func (r *Rule) AtRules() []string {
	var out []string
	for _, v := range r.Variants {
		if v.AtRule != "" {
			out = append(out, v.AtRule)
		}
	}
	return out
}

// AtRuleKey returns a stable key for the rule's group-rule context. Rules with
// the same key can be emitted inside one nested block.
func (r *Rule) AtRuleKey() string {
	return strings.Join(r.AtRules(), "\x00")
}

// AtRuleRank returns the sort rank of the rule's group-rule context.
func (r *Rule) AtRuleRank() int {
	rank := 0
	for _, v := range r.Variants {
		if v.AtRule != "" && v.Rank > rank {
			rank = v.Rank
		}
	}
	return rank
}

// VariantRank returns the strongest variant rank in the rule.
func (r *Rule) VariantRank() int {
	rank := 0
	for _, v := range r.Variants {
		if v.Rank > rank {
			rank = v.Rank
		}
	}
	return rank
}

// Empty reports whether the rule would generate no CSS.
func (r *Rule) Empty() bool {
	for _, b := range r.Blocks {
		if len(b.Decls) > 0 {
			return false
		}
	}
	return true
}

// Signature returns a stable textual description of the rule's output, used to
// detect duplicate rules deterministically.
func (r *Rule) Signature() string {
	var b strings.Builder
	b.WriteString(r.LayerName())
	b.WriteByte('\x04')
	b.WriteString(r.AtRuleKey())
	b.WriteByte('\x01')
	for _, blk := range r.Blocks {
		b.WriteString(blk.Target)
		b.WriteByte('\x02')
		for _, d := range blk.Decls {
			b.WriteString(d.Prop)
			b.WriteByte(':')
			b.WriteString(d.Value)
			b.WriteByte(';')
		}
		b.WriteByte('\x03')
	}
	return b.String()
}
