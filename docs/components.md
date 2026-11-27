# Components

FlexUI ships a component layer written as plain CSS classes. Components are not
a separate framework: they use semantic HTML, the same design tokens as
utilities, and live in `@layer flexui.components`, so a utility class always wins
when you need an override.

A component class is emitted only when it appears in your scanned content, just
like a utility. There is no import or extra directive:

```html
<button class="fui-btn fui-btn--primary">Save changes</button>
```

```css
/* src/styles.css */
@flexui base;
@flexui utilities;
```

The generated stylesheet contains `flexui.theme`, `flexui.base`,
`flexui.components` and `flexui.utilities` in that order.

## Buttons and badges

| Class | Purpose |
| --- | --- |
| `fui-btn` | Base button, link or label |
| `fui-btn--primary`, `--secondary`, `--outline`, `--ghost`, `--danger` | Variants |
| `fui-btn--sm`, `--lg`, `--icon`, `--block` | Size and layout |
| `fui-badge`, `fui-badge--accent`, `--success`, `--warning`, `--danger` | Status pills |

```html
<button class="fui-btn fui-btn--primary">Primary</button>
<button class="fui-btn fui-btn--outline fui-btn--sm">Small outline</button>
<button class="fui-btn fui-btn--icon" aria-label="Close">×</button>
<span class="fui-badge fui-badge--success">Active</span>
```

## Cards, stats and data display

`fui-card` + `fui-card-header` / `fui-card-title` / `fui-card-body` /
`fui-card-footer` / `fui-card-grid`; `fui-stat` with `fui-stat-label`,
`fui-stat-value` and `fui-stat-delta` (add `data-trend="up"` or `down`).

```html
<article class="fui-card">
  <div class="fui-card-header"><h2 class="fui-card-title">Revenue</h2></div>
  <div class="fui-card-body">
    <div class="fui-stat">
      <span class="fui-stat-label">This month</span>
      <span class="fui-stat-value">$24,800</span>
      <span class="fui-stat-delta" data-trend="up">12.4%</span>
    </div>
  </div>
</article>
```

Also available: `fui-avatar` (`--sm`, `--lg`), `fui-separator`
(`--vertical`), `fui-skeleton`, `fui-spinner`, `fui-progress` +
`fui-progress-bar`, `fui-table-wrap` / `fui-table` / `fui-tr` / `fui-th` /
`fui-td` (`--num`), `fui-timeline` (`fui-timeline-item`, `fui-timeline-dot`,
`fui-timeline-body`), `fui-tree`, `fui-calendar`, `fui-chart` and `fui-price`.

## Forms

`fui-field`, `fui-label`, `fui-help`, `fui-error`, `fui-input`, `fui-select`,
`fui-textarea`, `fui-checkbox` with `fui-check`, `fui-radio`, `fui-switch`,
`fui-slider`, `fui-search` + `fui-search-input`, `fui-fieldset` and
`fui-legend`.

```html
<div class="fui-field">
  <label class="fui-label" for="email">Email</label>
  <input class="fui-input" id="email" type="email" aria-describedby="email-help" />
  <p class="fui-help" id="email-help">We only use this for receipts.</p>
</div>

<label class="fui-check">
  <input class="fui-checkbox" type="checkbox" checked />
  Email me product updates
</label>
```

Pair `fui-error` with `aria-invalid="true"` and `aria-describedby` on the input.

## Navigation

`fui-navbar` + `fui-navbar-brand` / `fui-navbar-links` / `fui-navbar-link`;
`fui-breadcrumb` + `fui-breadcrumb-item` / `fui-breadcrumb-sep`;
`fui-pagination`; `fui-tabs` + `fui-tablist` / `fui-tab` / `fui-tabpanel`;
`fui-sidebar` + `fui-sidebar-section` / `fui-sidebar-link`; `fui-steps` +
`fui-step`; `fui-menu` + `fui-menu-label` / `fui-menu-item` / `fui-menu-sep`.

Add `data-fui-tabs`, `data-fui-dropdown` or `data-fui-command` to opt into the
[runtime](runtime.md). Without it, the markup still renders and native controls
still work.

## Overlays

`fui-dialog` with `fui-dialog-header` / `fui-dialog-title` / `fui-dialog-body` /
`fui-dialog-footer` / `fui-dialog-close`; `fui-drawer` with `fui-drawer-header`
/ `fui-drawer-body` / `fui-drawer-footer` and `fui-drawer--left`; `fui-popover`;
`fui-toast-region` / `fui-toast` (`--info`, `--success`, `--danger`) with
`fui-toast-content` / `fui-toast-title` / `fui-toast-description` /
`fui-toast-close`; `fui-tooltip`.

```html
<dialog class="fui-dialog" data-fui-dialog aria-labelledby="title">
  <div class="fui-dialog-header">
    <h2 class="fui-dialog-title" id="title">Edit profile</h2>
    <button class="fui-dialog-close" data-fui-close aria-label="Close">×</button>
  </div>
  <div class="fui-dialog-body">…</div>
  <div class="fui-dialog-footer">
    <button class="fui-btn fui-btn--ghost" data-fui-close>Cancel</button>
    <button class="fui-btn fui-btn--primary">Save</button>
  </div>
</dialog>
```

## Command palette

`fui-cmd`, `fui-cmd-input`, `fui-cmd-list`, `fui-cmd-item`, `fui-cmd-empty`.
Use them inside a dialog with `data-fui-command`; see the
[runtime guide](runtime.md#command-palette).

## Layout presets

Optional shells built from the same tokens: `fui-app` (`fui-app-header`,
`fui-app-sidebar`, `fui-app-main`), `fui-dashboard` + `fui-dashboard-grid` /
`fui-dashboard-span-2`, `fui-docs` (`fui-docs-nav`, `fui-docs-article`,
`fui-docs-toc`), `fui-settings` (`fui-settings-nav`, `fui-settings-panel`),
`fui-shop-grid`, `fui-product-media`, `fui-bento` (`fui-bento-wide`,
`fui-bento-tall`), `fui-page` and `fui-search-results`.

## Theming and overrides

Components read the shared `--flexui-*` custom properties for colour, radius,
shadow, spacing and focus rings. Override a token once, or add utilities on the
element:

```css
:root {
  --flexui-accent: oklch(0.62 0.19 265);
}
```

```html
<button class="fui-btn fui-btn--primary rounded-full px-6">Pill button</button>
```

Because utilities are in a later cascade layer, `px-6` and `rounded-full` win
over the component defaults without `!important`.

## Custom components

The registry is not limited to the classes listed above. Define your own
components in `flexui.config.json` with base declarations, variants, parts,
state selectors and `extends`, then use them as `fui-<name>` classes. A custom
component can also replace a built-in of the same name.

See [custom components](custom-components.md) for the configuration format, the
generated class names, inheritance rules, validation and the Go registry API.
