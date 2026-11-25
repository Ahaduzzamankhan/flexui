package components

// registerFoundation adds Button, Badge, Avatar, Card, Separator, Skeleton and
// Spinner.
func registerFoundation(r *Registry) {
	r.Register("fui-btn", orderFoundation, Self(
		Decl("display", "inline-flex"),
		Decl("align-items", "center"),
		Decl("justify-content", "center"),
		Decl("gap", "0.5rem"),
		Decl("padding-inline", "1rem"),
		Decl("padding-block", "0.625rem"),
		Decl("border-radius", "0.75rem"),
		Decl("border", "1px solid var(--flexui-line)"),
		Decl("background-color", "var(--flexui-surface-raised)"),
		Decl("color", "var(--flexui-fg)"),
		Decl("font-size", "0.875rem"),
		Decl("font-weight", "600"),
		Decl("line-height", "1.25rem"),
		Decl("cursor", "pointer"),
		Decl("text-decoration", "none"),
		Decl("white-space", "nowrap"),
		Decl("transition", "background-color 150ms ease, border-color 150ms ease, color 150ms ease, box-shadow 150ms ease"),
	), Part("&:hover", Decl("background-color", "var(--flexui-surface-sunken)")),
		Part("&:focus-visible", Decl("outline", "2px solid var(--flexui-focus-ring)"), Decl("outline-offset", "2px")),
		Part("&:disabled, &[aria-disabled=\"true\"]", Decl("opacity", "0.5"), Decl("pointer-events", "none")),
		Part("& > svg", Decl("width", "1rem"), Decl("height", "1rem"), Decl("flex-shrink", "0")),
	)
	r.Register("fui-btn--primary", orderFoundation+1, Self(
		Decl("background-color", "var(--flexui-accent)"),
		Decl("border-color", "transparent"),
		Decl("color", "var(--flexui-accent-fg)"),
	), Part("&:hover", Decl("background-color", "var(--flexui-accent-hover)")))
	r.Register("fui-btn--secondary", orderFoundation+2, Self(
		Decl("background-color", "var(--flexui-accent-muted)"),
		Decl("border-color", "transparent"),
		Decl("color", "var(--flexui-accent)"),
	))
	r.Register("fui-btn--outline", orderFoundation+3, Self(
		Decl("background-color", "transparent"),
		Decl("border-color", "var(--flexui-line-strong)"),
	))
	r.Register("fui-btn--ghost", orderFoundation+4, Self(
		Decl("background-color", "transparent"),
		Decl("border-color", "transparent"),
	), Part("&:hover", Decl("background-color", "var(--flexui-surface-sunken)")))
	r.Register("fui-btn--danger", orderFoundation+5, Self(
		Decl("background-color", "var(--flexui-danger)"),
		Decl("border-color", "transparent"),
		Decl("color", "var(--flexui-fg-inverted)"),
		Decl("background-color", "var(--flexui-danger)"),
		Decl("color", "#ffffff"),
	))
	r.Register("fui-btn--sm", orderFoundation+6, Self(
		Decl("padding-inline", "0.75rem"),
		Decl("padding-block", "0.375rem"),
		Decl("font-size", "0.8125rem"),
	))
	r.Register("fui-btn--lg", orderFoundation+7, Self(
		Decl("padding-inline", "1.25rem"),
		Decl("padding-block", "0.75rem"),
		Decl("font-size", "0.9375rem"),
	))
	r.Register("fui-btn--block", orderFoundation+8, Self(Decl("width", "100%")))
	r.Register("fui-btn--icon", orderFoundation+9, Self(
		Decl("padding", "0.5rem"),
		Decl("aspect-ratio", "1 / 1"),
	))

	r.Register("fui-badge", orderFoundation+10, Self(
		Decl("display", "inline-flex"),
		Decl("align-items", "center"),
		Decl("gap", "0.25rem"),
		Decl("padding-inline", "0.5rem"),
		Decl("padding-block", "0.125rem"),
		Decl("border-radius", "9999px"),
		Decl("border", "1px solid var(--flexui-line)"),
		Decl("background-color", "var(--flexui-surface-sunken)"),
		Decl("color", "var(--flexui-fg-muted)"),
		Decl("font-size", "0.75rem"),
		Decl("font-weight", "500"),
		Decl("line-height", "1.25rem"),
		Decl("white-space", "nowrap"),
	))
	r.Register("fui-badge--accent", orderFoundation+11, Self(
		Decl("background-color", "var(--flexui-accent-muted)"),
		Decl("border-color", "transparent"),
		Decl("color", "var(--flexui-accent)"),
	))
	r.Register("fui-badge--success", orderFoundation+12, Self(
		Decl("background-color", "var(--flexui-success-muted)"),
		Decl("border-color", "transparent"),
		Decl("color", "var(--flexui-success)"),
	))
	r.Register("fui-badge--warning", orderFoundation+13, Self(
		Decl("background-color", "var(--flexui-warning-muted)"),
		Decl("border-color", "transparent"),
		Decl("color", "var(--flexui-warning)"),
	))
	r.Register("fui-badge--danger", orderFoundation+14, Self(
		Decl("background-color", "var(--flexui-danger-muted)"),
		Decl("border-color", "transparent"),
		Decl("color", "var(--flexui-danger)"),
	))

	r.Register("fui-avatar", orderFoundation+20, Self(
		Decl("display", "inline-flex"),
		Decl("align-items", "center"),
		Decl("justify-content", "center"),
		Decl("width", "2.5rem"),
		Decl("height", "2.5rem"),
		Decl("border-radius", "9999px"),
		Decl("background-color", "var(--flexui-surface-sunken)"),
		Decl("color", "var(--flexui-fg-muted)"),
		Decl("font-size", "0.875rem"),
		Decl("font-weight", "600"),
		Decl("overflow", "hidden"),
		Decl("flex-shrink", "0"),
	), Part("& > img", Decl("width", "100%"), Decl("height", "100%"), Decl("object-fit", "cover")))
	r.Register("fui-avatar--sm", orderFoundation+21, Self(Decl("width", "2rem"), Decl("height", "2rem"), Decl("font-size", "0.75rem")))
	r.Register("fui-avatar--lg", orderFoundation+22, Self(Decl("width", "3.5rem"), Decl("height", "3.5rem"), Decl("font-size", "1rem")))

	r.Register("fui-card", orderFoundation+30, Self(
		Decl("display", "flex"),
		Decl("flex-direction", "column"),
		Decl("border-radius", "1rem"),
		Decl("border", "1px solid var(--flexui-line)"),
		Decl("background-color", "var(--flexui-surface-raised)"),
		Decl("color", "var(--flexui-fg)"),
		Decl("box-shadow", "0 1px 3px 0 rgb(0 0 0 / 0.06)"),
	))
	r.Register("fui-card-header", orderFoundation+31, Self(
		Decl("display", "flex"),
		Decl("align-items", "center"),
		Decl("justify-content", "space-between"),
		Decl("gap", "0.75rem"),
		Decl("padding", "1.25rem 1.25rem 0"),
	))
	r.Register("fui-card-title", orderFoundation+32, Self(
		Decl("font-size", "1rem"),
		Decl("font-weight", "600"),
		Decl("letter-spacing", "-0.01em"),
	))
	r.Register("fui-card-body", orderFoundation+33, Self(
		Decl("padding", "1rem 1.25rem"),
		Decl("font-size", "0.875rem"),
		Decl("line-height", "1.5"),
		Decl("color", "var(--flexui-fg-muted)"),
	))
	r.Register("fui-card-footer", orderFoundation+34, Self(
		Decl("display", "flex"),
		Decl("align-items", "center"),
		Decl("justify-content", "flex-end"),
		Decl("gap", "0.5rem"),
		Decl("padding", "0 1.25rem 1.25rem"),
	))

	r.Register("fui-separator", orderFoundation+40, Self(
		Decl("height", "1px"),
		Decl("width", "100%"),
		Decl("background-color", "var(--flexui-line)"),
		Decl("border", "0"),
	))
	r.Register("fui-separator--vertical", orderFoundation+41, Self(
		Decl("height", "auto"),
		Decl("width", "1px"),
		Decl("align-self", "stretch"),
	))

	r.Register("fui-skeleton", orderFoundation+50, Self(
		Decl("display", "block"),
		Decl("border-radius", "0.5rem"),
		Decl("background-color", "var(--flexui-surface-sunken)"),
		Decl("animation", "fui-skeleton-pulse 1.6s ease-in-out infinite"),
		Decl("min-height", "1rem"),
	))
	r.Register("fui-spinner", orderFoundation+51, Self(
		Decl("display", "inline-block"),
		Decl("width", "1.25rem"),
		Decl("height", "1.25rem"),
		Decl("border-radius", "9999px"),
		Decl("border", "2px solid var(--flexui-line)"),
		Decl("border-top-color", "var(--flexui-accent)"),
		Decl("animation", "fui-spin 700ms linear infinite"),
		Decl("flex-shrink", "0"),
	))
}
