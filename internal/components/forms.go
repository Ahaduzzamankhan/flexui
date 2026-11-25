package components

import "github.com/Ahaduzzamankhan/flexui/internal/ast"

// registerForms adds field, input, textarea, select, checkbox, radio, switch,
// slider, date picker and validation styles.
func registerForms(r *Registry) {
	inputDecls := []ast.Decl{
		Decl("display", "block"),
		Decl("width", "100%"),
		Decl("padding-inline", "0.75rem"),
		Decl("padding-block", "0.5rem"),
		Decl("border-radius", "0.625rem"),
		Decl("border", "1px solid var(--flexui-line-strong)"),
		Decl("background-color", "var(--flexui-surface-raised)"),
		Decl("color", "var(--flexui-fg)"),
		Decl("font-size", "0.875rem"),
		Decl("line-height", "1.5rem"),
		Decl("transition", "border-color 150ms ease, box-shadow 150ms ease"),
	}
	invalid := Part("&[aria-invalid=\"true\"], &:user-invalid",
		Decl("border-color", "var(--flexui-danger)"),
		Decl("box-shadow", "0 0 0 1px var(--flexui-danger)"),
	)
	focus := Part("&:focus-visible",
		Decl("outline", "2px solid var(--flexui-focus-ring)"),
		Decl("outline-offset", "1px"),
		Decl("border-color", "var(--flexui-accent)"),
	)
	disabled := Part("&:disabled, &[aria-disabled=\"true\"]", Decl("opacity", "0.6"), Decl("cursor", "not-allowed"))

	r.Register("fui-field", orderForms, Self(
		Decl("display", "flex"),
		Decl("flex-direction", "column"),
		Decl("gap", "0.375rem"),
	))
	r.Register("fui-label", orderForms+1, Self(
		Decl("font-size", "0.8125rem"),
		Decl("font-weight", "500"),
		Decl("color", "var(--flexui-fg-muted)"),
	))
	r.Register("fui-help", orderForms+2, Self(
		Decl("font-size", "0.75rem"),
		Decl("color", "var(--flexui-fg-subtle)"),
	))
	r.Register("fui-error", orderForms+3, Self(
		Decl("font-size", "0.75rem"),
		Decl("font-weight", "500"),
		Decl("color", "var(--flexui-danger)"),
	))
	r.Register("fui-input", orderForms+4, Self(inputDecls...), focus, invalid, disabled,
		Part("&::placeholder", Decl("color", "var(--flexui-fg-subtle)")))
	r.Register("fui-textarea", orderForms+5, Self(append(inputDecls, Decl("min-height", "6rem"), Decl("resize", "vertical"))...), focus, invalid, disabled)
	r.Register("fui-select", orderForms+6, Self(append(inputDecls, Decl("appearance", "none"), Decl("padding-right", "2rem"))...), focus, invalid, disabled)
	r.Register("fui-date", orderForms+7, Self(inputDecls...), focus, invalid, disabled,
		Part("&::-webkit-calendar-picker-indicator", Decl("cursor", "pointer"), Decl("opacity", "0.6")))

	checkboxDecls := []ast.Decl{
		Decl("appearance", "none"),
		Decl("width", "1.125rem"),
		Decl("height", "1.125rem"),
		Decl("flex-shrink", "0"),
		Decl("margin", "0"),
		Decl("border-radius", "0.375rem"),
		Decl("border", "1px solid var(--flexui-line-strong)"),
		Decl("background-color", "var(--flexui-surface-raised)"),
		Decl("cursor", "pointer"),
		Decl("transition", "background-color 150ms ease, border-color 150ms ease"),
	}
	checked := Part("&:checked",
		Decl("background-color", "var(--flexui-accent)"),
		Decl("border-color", "var(--flexui-accent)"),
		Decl("background-image", "url(\"data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 16 16' fill='white'%3E%3Cpath d='M12.2 4.6 6.6 10.2 3.8 7.4l-1 1 3.8 3.8 6.6-6.6z'/%3E%3C/svg%3E\")"),
		Decl("background-size", "100% 100%"),
	)
	r.Register("fui-checkbox", orderForms+10, Self(checkboxDecls...), checked, focus, disabled)
	r.Register("fui-radio", orderForms+11, Self(append(checkboxDecls, Decl("border-radius", "9999px"))...),
		Part("&:checked",
			Decl("border-color", "var(--flexui-accent)"),
			Decl("background-color", "var(--flexui-surface-raised)"),
			Decl("background-image", "radial-gradient(circle, var(--flexui-accent) 0 42%, transparent 46%)"),
		), focus, disabled)

	r.Register("fui-check", orderForms+12, Self(
		Decl("display", "flex"),
		Decl("align-items", "center"),
		Decl("gap", "0.5rem"),
		Decl("font-size", "0.875rem"),
		Decl("color", "var(--flexui-fg)"),
	))

	r.Register("fui-switch", orderForms+20, Self(
		Decl("position", "relative"),
		Decl("display", "inline-flex"),
		Decl("align-items", "center"),
		Decl("width", "2.5rem"),
		Decl("height", "1.5rem"),
		Decl("flex-shrink", "0"),
		Decl("border-radius", "9999px"),
		Decl("border", "1px solid var(--flexui-line-strong)"),
		Decl("background-color", "var(--flexui-surface-sunken)"),
		Decl("cursor", "pointer"),
		Decl("appearance", "none"),
		Decl("transition", "background-color 150ms ease, border-color 150ms ease"),
	), Part("&::after",
		Decl("content", "\"\""),
		Decl("position", "absolute"),
		Decl("left", "0.125rem"),
		Decl("width", "1.125rem"),
		Decl("height", "1.125rem"),
		Decl("border-radius", "9999px"),
		Decl("background-color", "var(--flexui-surface-raised)"),
		Decl("box-shadow", "0 1px 2px 0 rgb(0 0 0 / 0.2)"),
		Decl("transition", "transform 150ms ease"),
	), Part("&:checked",
		Decl("background-color", "var(--flexui-accent)"),
		Decl("border-color", "var(--flexui-accent)"),
	), Part("&:checked::after", Decl("transform", "translateX(1rem)")), focus)

	r.Register("fui-slider", orderForms+21, Self(
		Decl("width", "100%"),
		Decl("height", "0.375rem"),
		Decl("border-radius", "9999px"),
		Decl("background-color", "var(--flexui-surface-sunken)"),
		Decl("appearance", "none"),
		Decl("cursor", "pointer"),
		Decl("border", "1px solid var(--flexui-line)"),
	), Part("&::-webkit-slider-thumb",
		Decl("appearance", "none"),
		Decl("width", "1rem"),
		Decl("height", "1rem"),
		Decl("border-radius", "9999px"),
		Decl("background-color", "var(--flexui-accent)"),
		Decl("border", "2px solid var(--flexui-surface-raised)"),
		Decl("box-shadow", "0 1px 3px 0 rgb(0 0 0 / 0.3)"),
	), Part("&::-moz-range-thumb",
		Decl("width", "1rem"),
		Decl("height", "1rem"),
		Decl("border-radius", "9999px"),
		Decl("background-color", "var(--flexui-accent)"),
		Decl("border", "2px solid var(--flexui-surface-raised)"),
	))

	r.Register("fui-fieldset", orderForms+30, Self(
		Decl("border", "1px solid var(--flexui-line)"),
		Decl("border-radius", "0.75rem"),
		Decl("padding", "1rem"),
		Decl("margin", "0"),
	))
	r.Register("fui-legend", orderForms+31, Self(
		Decl("padding-inline", "0.375rem"),
		Decl("font-size", "0.8125rem"),
		Decl("font-weight", "600"),
		Decl("color", "var(--flexui-fg-muted)"),
	))
}
