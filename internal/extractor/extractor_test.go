package extractor

import (
	"reflect"
	"testing"
)

func extract(t *testing.T, path, src string) []string {
	t.Helper()
	return Extract(path, []byte(src))
}

func contains(list []string, want string) bool {
	for _, v := range list {
		if v == want {
			return true
		}
	}
	return false
}

func TestExtractHTMLClassAttribute(t *testing.T) {
	got := extract(t, "index.html", `<div class="flex p-4 md:hover:bg-red-500" data-note="not-a-class">x</div>`)
	for _, want := range []string{"flex", "p-4", "md:hover:bg-red-500"} {
		if !contains(got, want) {
			t.Fatalf("missing %q in %v", want, got)
		}
	}
	if contains(got, "not-a-class") {
		t.Fatalf("non-class attribute was scanned: %v", got)
	}
}

func TestExtractSkipsCommentsAndStyles(t *testing.T) {
	src := `<!-- <div class="p-8"> -->
<style>.p-12 { padding: 3rem; }</style>
<div class="p-1"></div>`
	got := extract(t, "index.html", src)
	if contains(got, "p-8") || contains(got, "p-12") {
		t.Fatalf("comment or style leaked: %v", got)
	}
	if !contains(got, "p-1") {
		t.Fatalf("real class missing: %v", got)
	}
}

func TestExtractInlineScript(t *testing.T) {
	src := `<script>const cls = "items-center";</script><div class="flex"></div>`
	got := extract(t, "index.html", src)
	if !contains(got, "items-center") || !contains(got, "flex") {
		t.Fatalf("got %v", got)
	}
}

func TestExtractJSX(t *testing.T) {
	src := `<div className={clsx("p-4", active && "bg-blue-500", { "opacity-50": dim })} />`
	got := extract(t, "App.jsx", src)
	for _, want := range []string{"p-4", "bg-blue-500", "opacity-50"} {
		if !contains(got, want) {
			t.Fatalf("missing %q in %v", want, got)
		}
	}
}

func TestExtractVueBinding(t *testing.T) {
	src := `<div :class="['p-2', { 'text-sm': active }]"></div>`
	got := extract(t, "Card.vue", src)
	if !contains(got, "p-2") || !contains(got, "text-sm") {
		t.Fatalf("got %v", got)
	}
}

func TestExtractSvelteClassDirective(t *testing.T) {
	src := `<div class="flex" class:active></div>`
	got := extract(t, "Card.svelte", src)
	if !contains(got, "active") || !contains(got, "flex") {
		t.Fatalf("got %v", got)
	}
}

func TestExtractScriptStrings(t *testing.T) {
	src := `
// "hidden" is in a comment and must be ignored
const base = "flex items-center";
const padded = ` + "`p-4 ${cond ? 'm-2' : 'm-0'} gap-2`" + `;
const ignored = 42;
`
	got := extract(t, "lib.js", src)
	for _, want := range []string{"flex", "items-center", "p-4", "m-2", "m-0", "gap-2"} {
		if !contains(got, want) {
			t.Fatalf("missing %q in %v", want, got)
		}
	}
	if contains(got, "hidden") {
		t.Fatalf("comment was scanned: %v", got)
	}
}

func TestExtractTemplateLiteralAcrossLines(t *testing.T) {
	src := "const s = `\n  flex\n  p-4\n`;"
	got := extract(t, "lib.ts", src)
	if !contains(got, "flex") || !contains(got, "p-4") {
		t.Fatalf("got %v", got)
	}
}

func TestShapeOK(t *testing.T) {
	valid := []string{"flex", "p-4", "md:hover:bg-red-500", "bg-[#121212]", "[&>*]:p-2", "-mt-4"}
	for _, v := range valid {
		if !ShapeOK(v) {
			t.Fatalf("expected %q to be a candidate", v)
		}
	}
	invalid := []string{"", "a=b", "foo();", "x\ny", "[unbalanced", ":", "\"quoted\""}
	for _, v := range invalid {
		if ShapeOK(v) {
			t.Fatalf("expected %q to be rejected", v)
		}
	}
}

func TestExtractIsDeterministic(t *testing.T) {
	src := `<div class="p-4 flex p-4"></div>`
	first := extract(t, "index.html", src)
	second := extract(t, "index.html", src)
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("extraction is not deterministic: %v vs %v", first, second)
	}
}
