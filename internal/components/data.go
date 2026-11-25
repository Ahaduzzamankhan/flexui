package components

// registerData adds data table, calendar, timeline, tree view, progress,
// statistics and chart primitives.
func registerData(r *Registry) {
	r.Register("fui-table-wrap", orderData, Self(
		Decl("width", "100%"),
		Decl("overflow-x", "auto"),
		Decl("border", "1px solid var(--flexui-line)"),
		Decl("border-radius", "0.75rem"),
		Decl("background-color", "var(--flexui-surface-raised)"),
	))
	r.Register("fui-table", orderData+1, Self(
		Decl("width", "100%"),
		Decl("border-collapse", "collapse"),
		Decl("font-size", "0.875rem"),
		Decl("color", "var(--flexui-fg)"),
	))
	r.Register("fui-caption", orderData+2, Self(
		Decl("padding", "0.75rem 1rem"),
		Decl("text-align", "left"),
		Decl("font-size", "0.75rem"),
		Decl("color", "var(--flexui-fg-subtle)"),
	))
	r.Register("fui-tr", orderData+3, Self(), Part("& + &", Decl("border-top", "1px solid var(--flexui-line)")),
		Part("&:hover", Decl("background-color", "var(--flexui-surface-muted)")))
	r.Register("fui-th", orderData+4, Self(
		Decl("padding", "0.625rem 1rem"),
		Decl("text-align", "left"),
		Decl("font-size", "0.75rem"),
		Decl("font-weight", "600"),
		Decl("letter-spacing", "0.04em"),
		Decl("text-transform", "uppercase"),
		Decl("color", "var(--flexui-fg-subtle)"),
		Decl("white-space", "nowrap"),
		Decl("border-bottom", "1px solid var(--flexui-line)"),
	))
	r.Register("fui-td", orderData+5, Self(
		Decl("padding", "0.75rem 1rem"),
		Decl("vertical-align", "middle"),
	))
	r.Register("fui-td--num", orderData+6, Self(Decl("text-align", "right"), Decl("font-variant-numeric", "tabular-nums")))

	r.Register("fui-calendar", orderData+10, Self(
		Decl("width", "100%"),
		Decl("border", "1px solid var(--flexui-line)"),
		Decl("border-radius", "0.875rem"),
		Decl("padding", "0.875rem"),
		Decl("background-color", "var(--flexui-surface-raised)"),
	))
	r.Register("fui-calendar-head", orderData+11, Self(
		Decl("display", "flex"),
		Decl("align-items", "center"),
		Decl("justify-content", "space-between"),
		Decl("margin-bottom", "0.75rem"),
		Decl("font-size", "0.875rem"),
		Decl("font-weight", "600"),
	))
	r.Register("fui-calendar-grid", orderData+12, Self(
		Decl("display", "grid"),
		Decl("grid-template-columns", "repeat(7, minmax(0, 1fr))"),
		Decl("gap", "0.125rem"),
	))
	r.Register("fui-calendar-day", orderData+13, Self(
		Decl("display", "flex"),
		Decl("align-items", "center"),
		Decl("justify-content", "center"),
		Decl("aspect-ratio", "1 / 1"),
		Decl("border-radius", "0.5rem"),
		Decl("font-size", "0.8125rem"),
		Decl("color", "var(--flexui-fg)"),
		Decl("cursor", "pointer"),
	), Part("&:hover", Decl("background-color", "var(--flexui-surface-sunken)")),
		Part("&[aria-current=\"date\"]", Decl("background-color", "var(--flexui-accent)"), Decl("color", "var(--flexui-accent-fg)")),
		Part("&[data-muted=\"true\"]", Decl("color", "var(--flexui-fg-subtle)"), Decl("opacity", "0.6")))

	r.Register("fui-timeline", orderData+20, Self(
		Decl("display", "flex"),
		Decl("flex-direction", "column"),
		Decl("gap", "0"),
		Decl("margin", "0"),
		Decl("padding", "0"),
		Decl("list-style", "none"),
	))
	r.Register("fui-timeline-item", orderData+21, Self(
		Decl("position", "relative"),
		Decl("padding-left", "1.75rem"),
		Decl("padding-bottom", "1.25rem"),
	), Part("&::before",
		Decl("content", "\"\""),
		Decl("position", "absolute"),
		Decl("left", "0.4375rem"),
		Decl("top", "1.125rem"),
		Decl("bottom", "0"),
		Decl("width", "1px"),
		Decl("background-color", "var(--flexui-line)"),
	), Part("&:last-child", Decl("padding-bottom", "0")), Part("&:last-child::before", Decl("display", "none")))
	r.Register("fui-timeline-dot", orderData+22, Self(
		Decl("position", "absolute"),
		Decl("left", "0"),
		Decl("top", "0.375rem"),
		Decl("width", "0.875rem"),
		Decl("height", "0.875rem"),
		Decl("border-radius", "9999px"),
		Decl("border", "2px solid var(--flexui-accent)"),
		Decl("background-color", "var(--flexui-surface-raised)"),
	))
	r.Register("fui-timeline-body", orderData+23, Self(
		Decl("font-size", "0.875rem"),
		Decl("line-height", "1.55"),
		Decl("color", "var(--flexui-fg-muted)"),
	))

	r.Register("fui-tree", orderData+30, Self(Decl("font-size", "0.875rem"), Decl("color", "var(--flexui-fg)")))
	r.Register("fui-tree-item", orderData+31, Self(
		Decl("display", "flex"),
		Decl("align-items", "center"),
		Decl("gap", "0.375rem"),
		Decl("padding-inline", "0.5rem"),
		Decl("padding-block", "0.3125rem"),
		Decl("border-radius", "0.375rem"),
		Decl("cursor", "pointer"),
	), Part("&:hover", Decl("background-color", "var(--flexui-surface-sunken)")),
		Part("&[aria-selected=\"true\"]", Decl("background-color", "var(--flexui-accent-muted)"), Decl("color", "var(--flexui-accent)")))
	r.Register("fui-tree-group", orderData+32, Self(Decl("padding-left", "1rem"), Decl("display", "flex"), Decl("flex-direction", "column"), Decl("gap", "0.125rem")))

	r.Register("fui-progress", orderData+40, Self(
		Decl("position", "relative"),
		Decl("width", "100%"),
		Decl("height", "0.5rem"),
		Decl("border-radius", "9999px"),
		Decl("background-color", "var(--flexui-surface-sunken)"),
		Decl("overflow", "hidden"),
	))
	r.Register("fui-progress-bar", orderData+41, Self(
		Decl("width", "var(--fui-progress, 0%)"),
		Decl("height", "100%"),
		Decl("border-radius", "inherit"),
		Decl("background-color", "var(--flexui-accent)"),
		Decl("transition", "width 300ms cubic-bezier(0.22, 1, 0.36, 1)"),
	))

	r.Register("fui-stat", orderData+50, Self(
		Decl("display", "flex"),
		Decl("flex-direction", "column"),
		Decl("gap", "0.25rem"),
		Decl("padding", "1rem"),
		Decl("border-radius", "0.875rem"),
		Decl("border", "1px solid var(--flexui-line)"),
		Decl("background-color", "var(--flexui-surface-raised)"),
	))
	r.Register("fui-stat-label", orderData+51, Self(
		Decl("font-size", "0.75rem"),
		Decl("font-weight", "500"),
		Decl("letter-spacing", "0.04em"),
		Decl("text-transform", "uppercase"),
		Decl("color", "var(--flexui-fg-subtle)"),
	))
	r.Register("fui-stat-value", orderData+52, Self(
		Decl("font-size", "1.5rem"),
		Decl("font-weight", "600"),
		Decl("letter-spacing", "-0.02em"),
		Decl("font-variant-numeric", "tabular-nums"),
		Decl("color", "var(--flexui-fg)"),
	))
	r.Register("fui-stat-delta", orderData+53, Self(
		Decl("font-size", "0.8125rem"),
		Decl("font-weight", "500"),
		Decl("color", "var(--flexui-fg-muted)"),
	), Part("&[data-trend=\"up\"]", Decl("color", "var(--flexui-success)")),
		Part("&[data-trend=\"down\"]", Decl("color", "var(--flexui-danger)")))

	r.Register("fui-chart", orderData+60, Self(
		Decl("position", "relative"),
		Decl("width", "100%"),
		Decl("border-radius", "0.875rem"),
		Decl("border", "1px solid var(--flexui-line)"),
		Decl("padding", "1rem"),
		Decl("background-color", "var(--flexui-surface-raised)"),
	))
	r.Register("fui-chart-bars", orderData+61, Self(
		Decl("display", "flex"),
		Decl("align-items", "flex-end"),
		Decl("gap", "0.375rem"),
		Decl("height", "8rem"),
	))
	r.Register("fui-chart-bar", orderData+62, Self(
		Decl("flex", "1"),
		Decl("height", "var(--fui-value, 50%)"),
		Decl("border-radius", "0.25rem 0.25rem 0 0"),
		Decl("background-color", "var(--flexui-accent)"),
		Decl("opacity", "0.85"),
		Decl("transition", "opacity 150ms ease"),
	), Part("&:hover", Decl("opacity", "1")))
	r.Register("fui-chart-legend", orderData+63, Self(
		Decl("display", "flex"),
		Decl("flex-wrap", "wrap"),
		Decl("gap", "0.75rem"),
		Decl("margin-top", "0.75rem"),
		Decl("font-size", "0.75rem"),
		Decl("color", "var(--flexui-fg-muted)"),
	))
}
