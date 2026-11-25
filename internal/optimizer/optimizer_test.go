package optimizer

import (
	"testing"

	"github.com/Ahaduzzamankhan/flexui/internal/ast"
)

func ruleWith(class string, decls ...ast.Decl) ast.Rule {
	return ast.Rule{
		Class:        class,
		Utility:      "test",
		UtilityOrder: 100,
		Blocks:       []ast.Block{{Decls: decls}},
	}
}

func TestRulesDeduplicateIdenticalOutput(t *testing.T) {
	rules := []ast.Rule{
		ruleWith("p-4", ast.Decl{Prop: "padding", Value: "1rem"}),
		ruleWith("p-4", ast.Decl{Prop: "padding", Value: "1rem"}),
		ruleWith("p-2", ast.Decl{Prop: "padding", Value: "0.5rem"}),
	}
	got := Rules(rules)
	if len(got) != 2 {
		t.Fatalf("got %d rules, want 2", len(got))
	}
}

func TestRulesDropEmpty(t *testing.T) {
	got := Rules([]ast.Rule{{Class: "group", Utility: "group", Blocks: []ast.Block{{}}}})
	if len(got) != 0 {
		t.Fatalf("empty rule survived: %+v", got)
	}
}

func TestDeclDeduplicationKeepsLast(t *testing.T) {
	rules := []ast.Rule{ruleWith("x",
		ast.Decl{Prop: "padding", Value: "1rem"},
		ast.Decl{Prop: "margin", Value: "0"},
		ast.Decl{Prop: "padding", Value: "2rem"},
	)}
	got := Rules(rules)
	decls := got[0].Blocks[0].Decls
	if len(decls) != 2 {
		t.Fatalf("decls = %+v", decls)
	}
	if decls[0].Prop != "margin" || decls[1].Prop != "padding" || decls[1].Value != "2rem" {
		t.Fatalf("wrong declarations kept: %+v", decls)
	}
}

func TestMinify(t *testing.T) {
	cases := []struct{ in, want string }{
		{"/* comment */ .a { color : red ; }", ".a{color:red}"},
		{".a > .b , .c { margin : 0 auto }", ".a > .b,.c{margin:0 auto}"},
		{"@media (min-width: 700px) { .a { color: red } }", "@media (min-width: 700px){.a{color:red}}"},
		{".a { content: \"a  b\" }", ".a{content:\"a  b\"}"},
	}
	for _, tc := range cases {
		if got := string(Minify([]byte(tc.in))); got != tc.want {
			t.Fatalf("Minify(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
