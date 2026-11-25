package components

// registerNavigation adds Navbar, Sidebar, Breadcrumb, Tabs, Pagination,
// Dropdown menu and Command palette.
func registerNavigation(r *Registry) {
	r.Register("fui-navbar", orderNav, Self(
		Decl("display", "flex"),
		Decl("align-items", "center"),
		Decl("justify-content", "space-between"),
		Decl("gap", "1.5rem"),
		Decl("padding-inline", "1.5rem"),
		Decl("padding-block", "0.75rem"),
		Decl("border-bottom", "1px solid var(--flexui-line)"),
		Decl("background-color", "var(--flexui-surface)"),
	))
	r.Register("fui-navbar-brand", orderNav+1, Self(
		Decl("display", "flex"),
		Decl("align-items", "center"),
		Decl("gap", "0.5rem"),
		Decl("font-weight", "600"),
		Decl("letter-spacing", "-0.01em"),
		Decl("color", "var(--flexui-fg)"),
		Decl("text-decoration", "none"),
	))
	r.Register("fui-navbar-links", orderNav+2, Self(
		Decl("display", "flex"),
		Decl("align-items", "center"),
		Decl("gap", "1.5rem"),
		Decl("font-size", "0.875rem"),
	))
	r.Register("fui-navbar-link", orderNav+3, Self(
		Decl("color", "var(--flexui-fg-muted)"),
		Decl("text-decoration", "none"),
		Decl("transition", "color 150ms ease"),
	), Part("&:hover, &[aria-current=\"page\"]", Decl("color", "var(--flexui-fg)")))

	r.Register("fui-sidebar", orderNav+10, Self(
		Decl("display", "flex"),
		Decl("flex-direction", "column"),
		Decl("gap", "0.25rem"),
		Decl("width", "16rem"),
		Decl("padding", "1rem"),
		Decl("border-right", "1px solid var(--flexui-line)"),
		Decl("background-color", "var(--flexui-surface-muted)"),
		Decl("overflow-y", "auto"),
	))
	r.Register("fui-sidebar-section", orderNav+11, Self(
		Decl("padding-inline", "0.5rem"),
		Decl("padding-block", "0.375rem"),
		Decl("font-size", "0.6875rem"),
		Decl("font-weight", "600"),
		Decl("letter-spacing", "0.08em"),
		Decl("text-transform", "uppercase"),
		Decl("color", "var(--flexui-fg-subtle)"),
	))
	r.Register("fui-sidebar-link", orderNav+12, Self(
		Decl("display", "flex"),
		Decl("align-items", "center"),
		Decl("gap", "0.5rem"),
		Decl("padding-inline", "0.5rem"),
		Decl("padding-block", "0.4375rem"),
		Decl("border-radius", "0.5rem"),
		Decl("font-size", "0.875rem"),
		Decl("color", "var(--flexui-fg-muted)"),
		Decl("text-decoration", "none"),
		Decl("transition", "background-color 150ms ease, color 150ms ease"),
	), Part("&:hover", Decl("background-color", "var(--flexui-surface-sunken)"), Decl("color", "var(--flexui-fg)")),
		Part("&[aria-current=\"page\"]", Decl("background-color", "var(--flexui-accent-muted)"), Decl("color", "var(--flexui-accent)"), Decl("font-weight", "500")))

	r.Register("fui-breadcrumb", orderNav+20, Self(
		Decl("display", "flex"),
		Decl("flex-wrap", "wrap"),
		Decl("align-items", "center"),
		Decl("gap", "0.5rem"),
		Decl("font-size", "0.8125rem"),
		Decl("color", "var(--flexui-fg-subtle)"),
	))
	r.Register("fui-breadcrumb-item", orderNav+21, Self(
		Decl("color", "inherit"),
		Decl("text-decoration", "none"),
	), Part("&:hover", Decl("color", "var(--flexui-fg)")))
	r.Register("fui-breadcrumb-sep", orderNav+22, Self(
		Decl("color", "var(--flexui-fg-subtle)"),
		Decl("user-select", "none"),
	))

	r.Register("fui-tabs", orderNav+30, Self(Decl("display", "flex"), Decl("flex-direction", "column"), Decl("gap", "0.75rem")))
	r.Register("fui-tablist", orderNav+31, Self(
		Decl("display", "flex"),
		Decl("gap", "0.25rem"),
		Decl("border-bottom", "1px solid var(--flexui-line)"),
	))
	r.Register("fui-tab", orderNav+32, Self(
		Decl("padding-inline", "0.75rem"),
		Decl("padding-block", "0.5rem"),
		Decl("margin-bottom", "-1px"),
		Decl("border-bottom", "2px solid transparent"),
		Decl("background", "transparent"),
		Decl("border-top", "0"),
		Decl("border-left", "0"),
		Decl("border-right", "0"),
		Decl("font-size", "0.875rem"),
		Decl("font-weight", "500"),
		Decl("color", "var(--flexui-fg-muted)"),
		Decl("cursor", "pointer"),
	), Part("&:hover", Decl("color", "var(--flexui-fg)")),
		Part("&[aria-selected=\"true\"]", Decl("color", "var(--flexui-accent)"), Decl("border-bottom-color", "var(--flexui-accent)")),
		Part("&:focus-visible", Decl("outline", "2px solid var(--flexui-focus-ring)"), Decl("outline-offset", "-2px")))
	r.Register("fui-tabpanel", orderNav+33, Self(
		Decl("font-size", "0.9375rem"),
		Decl("line-height", "1.6"),
		Decl("color", "var(--flexui-fg-muted)"),
	), Part("&[hidden]", Decl("display", "none")))

	r.Register("fui-pagination", orderNav+40, Self(
		Decl("display", "flex"),
		Decl("align-items", "center"),
		Decl("gap", "0.25rem"),
	))
	r.Register("fui-page", orderNav+41, Self(
		Decl("display", "inline-flex"),
		Decl("align-items", "center"),
		Decl("justify-content", "center"),
		Decl("min-width", "2.25rem"),
		Decl("height", "2.25rem"),
		Decl("padding-inline", "0.5rem"),
		Decl("border-radius", "0.5rem"),
		Decl("border", "1px solid transparent"),
		Decl("font-size", "0.8125rem"),
		Decl("font-weight", "500"),
		Decl("color", "var(--flexui-fg-muted)"),
		Decl("text-decoration", "none"),
	), Part("&:hover", Decl("background-color", "var(--flexui-surface-sunken)")),
		Part("&[aria-current=\"page\"]", Decl("background-color", "var(--flexui-accent)"), Decl("color", "var(--flexui-accent-fg)")),
		Part("&[aria-disabled=\"true\"]", Decl("opacity", "0.45"), Decl("pointer-events", "none")))

	r.Register("fui-menu", orderNav+50, Self(
		Decl("min-width", "12rem"),
		Decl("padding", "0.375rem"),
		Decl("border-radius", "0.75rem"),
		Decl("border", "1px solid var(--flexui-line)"),
		Decl("background-color", "var(--flexui-surface-raised)"),
		Decl("box-shadow", "0 10px 30px -12px rgb(0 0 0 / 0.25)"),
		Decl("z-index", "50"),
	))
	r.Register("fui-menu-item", orderNav+51, Self(
		Decl("display", "flex"),
		Decl("align-items", "center"),
		Decl("gap", "0.5rem"),
		Decl("width", "100%"),
		Decl("padding-inline", "0.5rem"),
		Decl("padding-block", "0.4375rem"),
		Decl("border-radius", "0.5rem"),
		Decl("font-size", "0.875rem"),
		Decl("color", "var(--flexui-fg)"),
		Decl("text-align", "left"),
		Decl("text-decoration", "none"),
		Decl("background", "transparent"),
		Decl("border", "0"),
		Decl("cursor", "pointer"),
	), Part("&:hover, &[aria-selected=\"true\"], &:focus-visible",
		Decl("background-color", "var(--flexui-surface-sunken)"),
		Decl("outline", "none")),
		Part("&[aria-disabled=\"true\"]", Decl("opacity", "0.5"), Decl("pointer-events", "none")))
	r.Register("fui-menu-sep", orderNav+52, Self(
		Decl("height", "1px"),
		Decl("margin-block", "0.375rem"),
		Decl("background-color", "var(--flexui-line)"),
		Decl("border", "0"),
	))
	r.Register("fui-menu-label", orderNav+53, Self(
		Decl("padding-inline", "0.5rem"),
		Decl("padding-block", "0.25rem"),
		Decl("font-size", "0.6875rem"),
		Decl("font-weight", "600"),
		Decl("letter-spacing", "0.06em"),
		Decl("text-transform", "uppercase"),
		Decl("color", "var(--flexui-fg-subtle)"),
	))

	r.Register("fui-cmd", orderNav+60, Self(
		Decl("width", "100%"),
		Decl("max-width", "32rem"),
		Decl("border-radius", "1rem"),
		Decl("border", "1px solid var(--flexui-line)"),
		Decl("background-color", "var(--flexui-surface-raised)"),
		Decl("box-shadow", "0 24px 60px -20px rgb(0 0 0 / 0.45)"),
		Decl("overflow", "hidden"),
	))
	r.Register("fui-cmd-input", orderNav+61, Self(
		Decl("width", "100%"),
		Decl("padding", "0.875rem 1rem"),
		Decl("border", "0"),
		Decl("border-bottom", "1px solid var(--flexui-line)"),
		Decl("background", "transparent"),
		Decl("font-size", "0.9375rem"),
		Decl("color", "var(--flexui-fg)"),
	), Part("&:focus-visible", Decl("outline", "none")))
	r.Register("fui-cmd-list", orderNav+62, Self(
		Decl("max-height", "20rem"),
		Decl("overflow-y", "auto"),
		Decl("padding", "0.375rem"),
	))
	r.Register("fui-cmd-item", orderNav+63, Self(
		Decl("display", "flex"),
		Decl("align-items", "center"),
		Decl("gap", "0.625rem"),
		Decl("padding-inline", "0.625rem"),
		Decl("padding-block", "0.5rem"),
		Decl("border-radius", "0.5rem"),
		Decl("font-size", "0.875rem"),
		Decl("color", "var(--flexui-fg)"),
		Decl("cursor", "pointer"),
	), Part("&[aria-selected=\"true\"]", Decl("background-color", "var(--flexui-accent-muted)"), Decl("color", "var(--flexui-accent)")),
		Part("&[hidden]", Decl("display", "none")))
	r.Register("fui-cmd-empty", orderNav+64, Self(
		Decl("padding", "1rem"),
		Decl("text-align", "center"),
		Decl("font-size", "0.8125rem"),
		Decl("color", "var(--flexui-fg-subtle)"),
	))
}
