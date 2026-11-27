# Custom components

FlexUI ships 145 built-in component classes, but the component engine is open:
you define your own components in `flexui.config.json`, they resolve through the
same pipeline as utilities and built-ins, and only the classes you actually use
are emitted. There is no cap and no plugin to install.

## Define a component

```json
{
  "components": {
    "alert": {
      "base": {
        "display": "flex",
        "align-items": "flex-start",
        "gap": "0.75rem",
        "padding": "1rem",
        "border-radius": "0.75rem",
        "border": "1px solid var(--flexui-line)",
        "background-color": "var(--flexui-surface-raised)"
      },
      "variants": {
        "success": { "border-color": "var(--flexui-success)", "color": "var(--flexui-success)" },
        "danger": { "border-color": "var(--flexui-danger)", "color": "var(--flexui-danger)" }
      },
      "parts": {
        "title": { "font-weight": "600", "font-size": "0.9375rem" },
        "actions": { "display": "flex", "gap": "0.5rem", "margin-inline-start": "auto" }
      },
      "states": {
        ":hover": { "box-shadow": "0 4px 12px rgb(0 0 0 / 0.08)" },
        "[data-compact='true']": { "padding": "0.5rem 0.75rem" }
      }
    }
  }
}
```

Use it like any other component class:

```html
<div class="fui-alert fui-alert--success">
  <p class="fui-alert-title">Saved</p>
  <div class="fui-alert-actions">
    <button class="fui-btn fui-btn--sm">Undo</button>
  </div>
</div>
```

That produces classes in `@layer flexui.components`:

| Definition | Class | Notes |
| --- | --- | --- |
| `base` | `fui-<name>` | Element declarations, sorted by property |
| `variants.<v>` | `fui-<name>--<v>` | Emitted after the base class, so it wins |
| `parts.<p>` | `fui-<name>-<p>` | Independent class for a child element |
| `states.<selector>` | attached to `fui-<name>` | Selector suffix appended to `&`, e.g. `&:hover` |

Utilities still override components because `flexui.utilities` is declared after
`flexui.components`, so `class="fui-alert mt-6"` keeps the component styling and
adjusts the margin.

## Inheritance

`extends` copies declarations from another component before applying your own:

```json
{
  "components": {
    "callout": {
      "extends": "surface",
      "base": { "padding": "1.25rem" }
    },
    "banner": {
      "extends": "card",
      "base": { "border-style": "dashed" }
    }
  }
}
```

- Extending a **custom** component inherits its base declarations, variants,
  parts and states; your declarations win on conflict.
- Extending a **built-in** component (`btn`, `card`, `badge`, …) inherits the
  declarations that built-in applies to its own element. Built-in states and
  child parts are not inherited, so the result stays predictable.
- Chains are resolved in dependency order. An unknown parent or a cycle is a
  build error that names the config file.

## Overriding built-ins

A custom component with the same name as a built-in replaces it:

```json
{
  "components": {
    "btn": {
      "base": {
        "display": "inline-flex",
        "padding": "0.5rem 1rem",
        "border-radius": "9999px",
        "background-color": "var(--flexui-accent)",
        "color": "var(--flexui-accent-fg)"
      }
    }
  }
}
```

`fui-btn` keeps its name and its call sites; the CSS is yours. Combine this with
`extends: "btn"` in other components to build a design system on top of the
built-in foundation.

## Validation

Definitions are validated while the config loads, so a mistake is a build error
with the file name instead of a silently missing class:

- component names are lowercase letters, digits and single dashes (`hero-panel`);
- variant and part names follow the same rule;
- state selectors must start with `:`, `[`, `.`, `>`, `+`, `~` or whitespace;
- every declaration needs a non-empty property and value, and values cannot
  contain `{`, `}` or `;`;
- `extends` must name a custom component or a built-in one;
- cycles and duplicate definitions are rejected;
- a definition with no `base`, `variants`, `parts`, `states` or `extends` is
  rejected.

## Programmatic API (Go)

Library users can build a registry directly, without a config file:

```go
import (
	"github.com/Ahaduzzamankhan/flexui/internal/components"
	"github.com/Ahaduzzamankhan/flexui/internal/themes"
)

registry, err := components.NewWithCustom(themes.Default(), []components.CustomComponent{
	{
		Name:     "alert",
		Base:     map[string]string{"display": "flex", "gap": "0.75rem"},
		Variants: map[string]map[string]string{"success": {"color": "green"}},
	},
})
if err != nil {
	return err
}

rule, ok := registry.Resolve("fui-alert--success")
```

`ValidateCustom(definitions)` checks definitions without keeping the registry,
and `(*Registry).RegisterCustomComponents` adds them to an existing registry.

## Determinism

Custom components use fixed order buckets (`600` base, `601` variants, `602`
parts) and sort declarations by property, classes by name and variants/parts by
name. Two builds of the same sources and config therefore produce byte-identical
CSS, and the build cache key includes the whole config so a change invalidates it.
