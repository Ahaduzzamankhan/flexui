package animation

import (
	"strings"
	"testing"

	"github.com/Ahaduzzamankhan/flexui/internal/ast"
)

func TestEveryBuiltinResolvesWithKeyframes(t *testing.T) {
	reg := New(nil)
	for _, name := range Builtin() {
		blocks, keyframes, key, _, ok := reg.Resolve("animate-" + name)
		if !ok {
			t.Fatalf("animate-%s did not resolve", name)
		}
		if key != "animate-"+name {
			t.Fatalf("animate-%s key = %q", name, key)
		}
		if len(keyframes) != 1 || keyframes[0].Name != "fui-"+name {
			t.Fatalf("animate-%s keyframes = %+v", name, keyframes)
		}
		if len(keyframes[0].Steps) == 0 {
			t.Fatalf("animate-%s has no steps", name)
		}
		if len(blocks) != 1 || len(blocks[0].Decls) == 0 {
			t.Fatalf("animate-%s has no declarations", name)
		}
		for _, decl := range blocks[0].Decls {
			if strings.Contains(decl.Value, "var(--flexui-animation") && !strings.Contains(decl.Value, ",") {
				t.Fatalf("animate-%s emitted an unfallback'd var: %s", name, decl.Value)
			}
		}
	}
}

func TestUnknownAnimationIsRejected(t *testing.T) {
	reg := New(nil)
	for _, token := range []string{"animate-nope", "animate-", "animate", "fade-in", "animate-duration-fast"} {
		if _, _, _, _, ok := reg.Resolve(token); ok {
			t.Fatalf("%q should not resolve", token)
		}
	}
}

func TestAnimationControlUtilities(t *testing.T) {
	reg := New(nil)
	cases := map[string]string{
		"animate-duration-500":      "animation-duration: 500ms",
		"animate-duration-1s":       "animation-duration: 1s",
		"animate-delay-[250ms]":     "animation-delay: 250ms",
		"animate-ease-spring":       "animation-timing-function: cubic-bezier(0.22, 1, 0.36, 1)",
		"animate-iteration-3":       "animation-iteration-count: 3",
		"animate-fill-forwards":     "animation-fill-mode: forwards",
		"animate-direction-reverse": "animation-direction: reverse",
		"animate-play-paused":       "animation-play-state: paused",
	}
	for token, want := range cases {
		blocks, _, _, _, ok := reg.Resolve(token)
		if !ok {
			t.Fatalf("%s did not resolve", token)
		}
		got := blocks[0].Decls[0].Prop + ": " + blocks[0].Decls[0].Value
		if got != want {
			t.Fatalf("%s = %q, want %q", token, got, want)
		}
	}
}

func TestInvalidControlValuesAreRejected(t *testing.T) {
	reg := New(nil)
	for _, token := range []string{
		"animate-duration-abc",
		"animate-ease-bouncy",
		"animate-iteration--1",
		"animate-fill-sometimes",
	} {
		if _, _, _, _, ok := reg.Resolve(token); ok {
			t.Fatalf("%q should not resolve", token)
		}
	}
}

func TestStaggerTargetsChildren(t *testing.T) {
	reg := New(nil)
	blocks, _, _, _, ok := reg.Resolve("animate-stagger")
	if !ok {
		t.Fatalf("animate-stagger did not resolve")
	}
	if len(blocks) != 12 {
		t.Fatalf("stagger blocks = %d, want 12", len(blocks))
	}
	if blocks[0].Target != "& > :nth-child(1)" || blocks[0].Decls[0].Value != "0ms" {
		t.Fatalf("first stagger block = %+v", blocks[0])
	}
	if blocks[11].Decls[0].Value != "660ms" {
		t.Fatalf("last stagger delay = %q", blocks[11].Decls[0].Value)
	}
}

func TestCustomAnimationUsesDefaults(t *testing.T) {
	reg := New(map[string]Definition{
		"wiggle": {Steps: []ast.KeyframeStep{
			{Steps: "from", Decls: []ast.Decl{{Prop: "rotate", Value: "0deg"}}},
			{Steps: "to", Decls: []ast.Decl{{Prop: "rotate", Value: "4deg"}}},
		}},
	})
	blocks, keyframes, _, _, ok := reg.Resolve("animate-wiggle")
	if !ok {
		t.Fatalf("custom animation did not resolve")
	}
	if keyframes[0].Name != "fui-wiggle" {
		t.Fatalf("keyframe name = %q", keyframes[0].Name)
	}
	values := make(map[string]string, len(blocks[0].Decls))
	for _, d := range blocks[0].Decls {
		values[d.Prop] = d.Value
	}
	if !strings.Contains(values["animation-duration"], "300ms") {
		t.Fatalf("custom duration default missing: %q", values["animation-duration"])
	}
	if !strings.Contains(values["animation-fill-mode"], "both") {
		t.Fatalf("custom fill default missing: %q", values["animation-fill-mode"])
	}
}

func TestCustomAnimationWithoutStepsIsIgnored(t *testing.T) {
	reg := New(map[string]Definition{"empty": {Duration: "1s"}})
	if _, _, _, _, ok := reg.Resolve("animate-empty"); ok {
		t.Fatalf("animation without steps should not resolve")
	}
}

func TestTimelineUtilities(t *testing.T) {
	reg := New(nil)
	for _, token := range []string{"animate-timeline-view", "animate-timeline-scroll"} {
		blocks, _, _, _, ok := reg.Resolve(token)
		if !ok {
			t.Fatalf("%s did not resolve", token)
		}
		if blocks[0].Decls[0].Prop != "animation-timeline" {
			t.Fatalf("%s did not set animation-timeline", token)
		}
	}
}
