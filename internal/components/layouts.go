package components

// registerLayouts adds the configurable layout presets: application shell,
// dashboard, documentation, e-commerce, search, settings, multi-step form,
// bento and card grids. They are ordinary CSS grid/flex compositions, not a
// proprietary rendering model, and every part can be overridden with utilities.
func registerLayouts(r *Registry) {
	r.Register("fui-app", orderLayout, Self(
		Decl("display", "grid"),
		Decl("grid-template-columns", "minmax(0, 1fr)"),
		Decl("grid-template-rows", "auto auto 1fr"),
		Decl("grid-template-areas", "\"header\" \"sidebar\" \"main\""),
		Decl("min-height", "100dvh"),
	), At("@media (min-width: 1024px)",
		Decl("grid-template-columns", "auto minmax(0, 1fr)"),
		Decl("grid-template-rows", "auto 1fr"),
		Decl("grid-template-areas", "\"header header\" \"sidebar main\""),
	))
	r.Register("fui-app-header", orderLayout+1, Self(
		Decl("grid-area", "header"),
		Decl("position", "sticky"),
		Decl("top", "0"),
		Decl("z-index", "40"),
		Decl("border-bottom", "1px solid var(--flexui-line)"),
		Decl("background-color", "var(--flexui-surface)"),
	))
	r.Register("fui-app-sidebar", orderLayout+2, Self(
		Decl("grid-area", "sidebar"),
		Decl("border-right", "1px solid var(--flexui-line)"),
		Decl("background-color", "var(--flexui-surface-muted)"),
		Decl("overflow-y", "auto"),
	))
	r.Register("fui-app-main", orderLayout+3, Self(
		Decl("grid-area", "main"),
		Decl("min-width", "0"),
		Decl("padding", "1.5rem"),
	))

	r.Register("fui-dashboard", orderLayout+10, Self(
		Decl("display", "flex"),
		Decl("flex-direction", "column"),
		Decl("gap", "1.25rem"),
	))
	r.Register("fui-dashboard-grid", orderLayout+11, Self(
		Decl("display", "grid"),
		Decl("gap", "1rem"),
		Decl("grid-template-columns", "repeat(auto-fit, minmax(min(100%, 15rem), 1fr))"),
	))
	r.Register("fui-dashboard-span-2", orderLayout+12, Self(Decl("grid-column", "span 2")))

	r.Register("fui-docs", orderLayout+20, Self(
		Decl("display", "grid"),
		Decl("gap", "2rem"),
		Decl("grid-template-columns", "minmax(0, 1fr)"),
		Decl("max-width", "90rem"),
		Decl("margin-inline", "auto"),
		Decl("padding-inline", "1.5rem"),
		Decl("padding-block", "2rem"),
	), At("@media (min-width: 1024px)",
		Decl("grid-template-columns", "15rem minmax(0, 1fr)"),
	), At("@media (min-width: 1280px)",
		Decl("grid-template-columns", "15rem minmax(0, 1fr) 14rem"),
	))
	r.Register("fui-docs-nav", orderLayout+21, Self(
		Decl("position", "sticky"),
		Decl("top", "5rem"),
		Decl("align-self", "start"),
		Decl("display", "flex"),
		Decl("flex-direction", "column"),
		Decl("gap", "0.125rem"),
		Decl("font-size", "0.875rem"),
	))
	r.Register("fui-docs-article", orderLayout+22, Self(
		Decl("min-width", "0"),
		Decl("font-size", "1rem"),
		Decl("line-height", "1.7"),
		Decl("color", "var(--flexui-fg-muted)"),
	))
	r.Register("fui-docs-toc", orderLayout+23, Self(
		Decl("position", "sticky"),
		Decl("top", "5rem"),
		Decl("align-self", "start"),
		Decl("font-size", "0.8125rem"),
		Decl("color", "var(--flexui-fg-subtle)"),
	), At("@media (max-width: 1279.98px)", Decl("display", "none")))

	r.Register("fui-shop-grid", orderLayout+30, Self(
		Decl("display", "grid"),
		Decl("gap", "1.25rem"),
		Decl("grid-template-columns", "repeat(auto-fill, minmax(min(100%, 15rem), 1fr))"),
	))
	r.Register("fui-product-media", orderLayout+31, Self(
		Decl("aspect-ratio", "4 / 3"),
		Decl("overflow", "hidden"),
		Decl("border-radius", "0.75rem 0.75rem 0 0"),
		Decl("background-color", "var(--flexui-surface-sunken)"),
	), Part("& > img", Decl("width", "100%"), Decl("height", "100%"), Decl("object-fit", "cover")))
	r.Register("fui-price", orderLayout+32, Self(
		Decl("font-size", "1.125rem"),
		Decl("font-weight", "600"),
		Decl("font-variant-numeric", "tabular-nums"),
		Decl("color", "var(--flexui-fg)"),
	))

	r.Register("fui-search", orderLayout+40, Self(
		Decl("display", "flex"),
		Decl("align-items", "center"),
		Decl("gap", "0.5rem"),
		Decl("padding", "0.5rem 0.75rem"),
		Decl("border-radius", "0.75rem"),
		Decl("border", "1px solid var(--flexui-line-strong)"),
		Decl("background-color", "var(--flexui-surface-raised)"),
	))
	r.Register("fui-search-input", orderLayout+41, Self(
		Decl("flex", "1"),
		Decl("min-width", "0"),
		Decl("border", "0"),
		Decl("background", "transparent"),
		Decl("font-size", "0.875rem"),
		Decl("color", "var(--flexui-fg)"),
	), Part("&:focus-visible", Decl("outline", "none")))
	r.Register("fui-search-results", orderLayout+42, Self(
		Decl("display", "flex"),
		Decl("flex-direction", "column"),
		Decl("gap", "0.125rem"),
		Decl("margin-top", "0.5rem"),
		Decl("max-height", "24rem"),
		Decl("overflow-y", "auto"),
	))

	r.Register("fui-settings", orderLayout+50, Self(
		Decl("display", "grid"),
		Decl("gap", "1.5rem"),
		Decl("grid-template-columns", "minmax(0, 1fr)"),
	), At("@media (min-width: 48rem)",
		Decl("grid-template-columns", "12rem minmax(0, 1fr)"),
	))
	r.Register("fui-settings-nav", orderLayout+51, Self(
		Decl("display", "flex"),
		Decl("flex-direction", "column"),
		Decl("gap", "0.125rem"),
		Decl("font-size", "0.875rem"),
	))
	r.Register("fui-settings-panel", orderLayout+52, Self(
		Decl("display", "flex"),
		Decl("flex-direction", "column"),
		Decl("gap", "1rem"),
		Decl("min-width", "0"),
	))

	r.Register("fui-steps", orderLayout+60, Self(
		Decl("display", "flex"),
		Decl("flex-wrap", "wrap"),
		Decl("align-items", "center"),
		Decl("gap", "0.5rem"),
		Decl("margin", "0"),
		Decl("padding", "0"),
		Decl("list-style", "none"),
		Decl("counter-reset", "fui-step"),
	))
	r.Register("fui-step", orderLayout+61, Self(
		Decl("display", "flex"),
		Decl("align-items", "center"),
		Decl("gap", "0.5rem"),
		Decl("font-size", "0.8125rem"),
		Decl("font-weight", "500"),
		Decl("color", "var(--flexui-fg-subtle)"),
	), Part("&::before",
		Decl("content", "counter(fui-step)"),
		Decl("counter-increment", "fui-step"),
		Decl("display", "grid"),
		Decl("place-items", "center"),
		Decl("width", "1.5rem"),
		Decl("height", "1.5rem"),
		Decl("border-radius", "9999px"),
		Decl("border", "1px solid var(--flexui-line-strong)"),
		Decl("font-size", "0.75rem"),
	), Part("&[aria-current=\"step\"]", Decl("color", "var(--flexui-fg)")),
		Part("&[aria-current=\"step\"]::before",
			Decl("background-color", "var(--flexui-accent)"),
			Decl("border-color", "var(--flexui-accent)"),
			Decl("color", "var(--flexui-accent-fg)"),
		),
		Part("&[data-state=\"done\"]::before",
			Decl("background-color", "var(--flexui-success-muted)"),
			Decl("border-color", "var(--flexui-success)"),
			Decl("color", "var(--flexui-success)"),
		))

	r.Register("fui-bento", orderLayout+70, Self(
		Decl("display", "grid"),
		Decl("gap", "1rem"),
		Decl("grid-template-columns", "repeat(auto-fit, minmax(min(100%, 16rem), 1fr))"),
	))
	r.Register("fui-bento-wide", orderLayout+71, Self(Decl("grid-column", "span 2")), At("@media (max-width: 47.98rem)", Decl("grid-column", "auto")))
	r.Register("fui-bento-tall", orderLayout+72, Self(Decl("grid-row", "span 2")))
	r.Register("fui-card-grid", orderLayout+73, Self(
		Decl("display", "grid"),
		Decl("gap", "1rem"),
		Decl("grid-template-columns", "repeat(auto-fill, minmax(min(100%, 18rem), 1fr))"),
	))
}
