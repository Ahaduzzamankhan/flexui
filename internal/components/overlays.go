package components

// registerOverlays adds Dialog, Drawer, Popover, Tooltip, Toast and Context
// menu. Everything here degrades gracefully: dialogs use the native <dialog>
// element, popovers use the native popover attribute, tooltips are CSS only and
// toasts are announced through a live region.
func registerOverlays(r *Registry) {
	r.Register("fui-dialog", orderOverlay, Self(
		Decl("padding", "0"),
		Decl("border", "1px solid var(--flexui-line)"),
		Decl("border-radius", "1rem"),
		Decl("background-color", "var(--flexui-surface-raised)"),
		Decl("color", "var(--flexui-fg)"),
		Decl("box-shadow", "0 30px 80px -24px rgb(0 0 0 / 0.45)"),
		Decl("width", "calc(100% - 2rem)"),
		Decl("max-width", "32rem"),
	), Part("&::backdrop", Decl("background-color", "var(--flexui-overlay)"), Decl("backdrop-filter", "blur(2px)")),
		Part("&[open]", Decl("animation", "fui-dialog-in 180ms cubic-bezier(0.22, 1, 0.36, 1)")))
	r.Register("fui-dialog-header", orderOverlay+1, Self(
		Decl("display", "flex"),
		Decl("align-items", "flex-start"),
		Decl("justify-content", "space-between"),
		Decl("gap", "1rem"),
		Decl("padding", "1.25rem 1.25rem 0"),
	))
	r.Register("fui-dialog-title", orderOverlay+2, Self(
		Decl("font-size", "1.0625rem"),
		Decl("font-weight", "600"),
		Decl("letter-spacing", "-0.01em"),
		Decl("margin", "0"),
	))
	r.Register("fui-dialog-body", orderOverlay+3, Self(
		Decl("padding", "0.75rem 1.25rem 1.25rem"),
		Decl("font-size", "0.9375rem"),
		Decl("line-height", "1.6"),
		Decl("color", "var(--flexui-fg-muted)"),
	))
	r.Register("fui-dialog-footer", orderOverlay+4, Self(
		Decl("display", "flex"),
		Decl("align-items", "center"),
		Decl("justify-content", "flex-end"),
		Decl("gap", "0.5rem"),
		Decl("padding", "0 1.25rem 1.25rem"),
	))
	r.Register("fui-dialog-close", orderOverlay+5, Self(
		Decl("display", "inline-flex"),
		Decl("align-items", "center"),
		Decl("justify-content", "center"),
		Decl("width", "2rem"),
		Decl("height", "2rem"),
		Decl("border-radius", "0.5rem"),
		Decl("border", "0"),
		Decl("background", "transparent"),
		Decl("color", "var(--flexui-fg-muted)"),
		Decl("cursor", "pointer"),
		Decl("flex-shrink", "0"),
	), Part("&:hover", Decl("background-color", "var(--flexui-surface-sunken)"), Decl("color", "var(--flexui-fg)")))

	// Drawers are <dialog> elements anchored to an edge.
	r.Register("fui-drawer", orderOverlay+10, Self(
		Decl("margin", "0 0 0 auto"),
		Decl("height", "100dvh"),
		Decl("max-height", "100dvh"),
		Decl("width", "min(24rem, 100%)"),
		Decl("max-width", "none"),
		Decl("padding", "0"),
		Decl("border", "0"),
		Decl("border-left", "1px solid var(--flexui-line)"),
		Decl("border-radius", "0"),
		Decl("background-color", "var(--flexui-surface-raised)"),
		Decl("color", "var(--flexui-fg)"),
		Decl("box-shadow", "-20px 0 60px -30px rgb(0 0 0 / 0.45)"),
	), Part("&::backdrop", Decl("background-color", "var(--flexui-overlay)")),
		Part("&[open]", Decl("animation", "fui-drawer-in 200ms cubic-bezier(0.22, 1, 0.36, 1)")))
	r.Register("fui-drawer--left", orderOverlay+11, Self(Decl("margin", "0 auto 0 0"), Decl("border-left", "0"), Decl("border-right", "1px solid var(--flexui-line)")))
	r.Register("fui-drawer-header", orderOverlay+12, Self(
		Decl("display", "flex"),
		Decl("align-items", "center"),
		Decl("justify-content", "space-between"),
		Decl("gap", "1rem"),
		Decl("padding", "1rem 1.25rem"),
		Decl("border-bottom", "1px solid var(--flexui-line)"),
	))
	r.Register("fui-drawer-body", orderOverlay+13, Self(
		Decl("padding", "1.25rem"),
		Decl("overflow-y", "auto"),
		Decl("height", "calc(100dvh - 8rem)"),
	))
	r.Register("fui-drawer-footer", orderOverlay+14, Self(
		Decl("display", "flex"),
		Decl("gap", "0.5rem"),
		Decl("padding", "1rem 1.25rem"),
		Decl("border-top", "1px solid var(--flexui-line)"),
	))

	// Popovers use the native popover attribute; the anchor is the trigger.
	r.Register("fui-popover", orderOverlay+20, Self(
		Decl("margin", "0"),
		Decl("padding", "0.875rem"),
		Decl("max-width", "20rem"),
		Decl("border", "1px solid var(--flexui-line)"),
		Decl("border-radius", "0.75rem"),
		Decl("background-color", "var(--flexui-surface-raised)"),
		Decl("color", "var(--flexui-fg)"),
		Decl("font-size", "0.875rem"),
		Decl("box-shadow", "0 16px 40px -18px rgb(0 0 0 / 0.35)"),
		Decl("inset", "auto"),
	))

	// CSS-only tooltip: content comes from the attribute.
	r.Register("fui-tooltip", orderOverlay+30, Self(Decl("position", "relative"), Decl("display", "inline-flex")),
		Part("&::after",
			Decl("content", "attr(data-fui-tooltip)"),
			Decl("position", "absolute"),
			Decl("left", "50%"),
			Decl("bottom", "calc(100% + 0.5rem)"),
			Decl("transform", "translateX(-50%)"),
			Decl("padding-inline", "0.5rem"),
			Decl("padding-block", "0.25rem"),
			Decl("border-radius", "0.375rem"),
			Decl("background-color", "var(--flexui-fg)"),
			Decl("color", "var(--flexui-surface)"),
			Decl("font-size", "0.75rem"),
			Decl("line-height", "1.25rem"),
			Decl("white-space", "nowrap"),
			Decl("opacity", "0"),
			Decl("pointer-events", "none"),
			Decl("transition", "opacity 120ms ease 60ms"),
			Decl("z-index", "60"),
		),
		Part("&:hover::after, &:focus-visible::after", Decl("opacity", "1")),
	)

	r.Register("fui-toast-region", orderOverlay+40, Self(
		Decl("position", "fixed"),
		Decl("right", "1rem"),
		Decl("bottom", "1rem"),
		Decl("z-index", "100"),
		Decl("display", "flex"),
		Decl("flex-direction", "column"),
		Decl("gap", "0.5rem"),
		Decl("width", "min(22rem, calc(100vw - 2rem))"),
		Decl("pointer-events", "none"),
	))
	r.Register("fui-toast", orderOverlay+41, Self(
		Decl("pointer-events", "auto"),
		Decl("display", "flex"),
		Decl("align-items", "flex-start"),
		Decl("gap", "0.625rem"),
		Decl("padding", "0.75rem 0.875rem"),
		Decl("border-radius", "0.75rem"),
		Decl("border", "1px solid var(--flexui-line)"),
		Decl("background-color", "var(--flexui-surface-raised)"),
		Decl("color", "var(--flexui-fg)"),
		Decl("font-size", "0.875rem"),
		Decl("box-shadow", "0 12px 32px -16px rgb(0 0 0 / 0.35)"),
		Decl("animation", "fui-toast-in 200ms cubic-bezier(0.22, 1, 0.36, 1)"),
	))
	r.Register("fui-toast--success", orderOverlay+42, Self(Decl("border-left", "3px solid var(--flexui-success)")))
	r.Register("fui-toast--danger", orderOverlay+43, Self(Decl("border-left", "3px solid var(--flexui-danger)")))
	r.Register("fui-toast--info", orderOverlay+44, Self(Decl("border-left", "3px solid var(--flexui-accent)")))
	r.Register("fui-toast-content", orderOverlay+45, Self(
		Decl("display", "flex"),
		Decl("flex-direction", "column"),
		Decl("gap", "0.125rem"),
		Decl("flex", "1"),
		Decl("min-width", "0"),
	))
	r.Register("fui-toast-title", orderOverlay+46, Self(
		Decl("margin", "0"),
		Decl("font-size", "0.875rem"),
		Decl("font-weight", "600"),
		Decl("line-height", "1.35"),
		Decl("color", "var(--flexui-fg)"),
	))
	r.Register("fui-toast-description", orderOverlay+47, Self(
		Decl("margin", "0"),
		Decl("font-size", "0.8125rem"),
		Decl("line-height", "1.45"),
		Decl("color", "var(--flexui-fg-muted)"),
	))
	r.Register("fui-toast-close", orderOverlay+48, Self(
		Decl("display", "inline-flex"),
		Decl("align-items", "center"),
		Decl("justify-content", "center"),
		Decl("width", "1.5rem"),
		Decl("height", "1.5rem"),
		Decl("flex-shrink", "0"),
		Decl("padding", "0"),
		Decl("border", "0"),
		Decl("border-radius", "0.375rem"),
		Decl("background", "transparent"),
		Decl("color", "var(--flexui-fg-subtle)"),
		Decl("font-size", "1rem"),
		Decl("line-height", "1"),
		Decl("cursor", "pointer"),
	), Part("&:hover", Decl("background-color", "var(--flexui-surface-sunken)"), Decl("color", "var(--flexui-fg)")),
		Part("&:focus-visible", Decl("outline", "2px solid var(--flexui-focus-ring)"), Decl("outline-offset", "1px")))
}
