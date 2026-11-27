# flexui-runtime

The optional browser runtime for FlexUI components. It adds keyboard
behaviour, focus management, ARIA state and DOM events on top of native
`<dialog>`, the popover API and your existing markup. The CSS works on its own;
this package is what makes the interactive components feel finished.

Zero dependencies, ESM, about as small as a runtime can be.

```bash
npm install flexui-runtime
```

```js
import { init } from "flexui-runtime";

const cleanup = init();
```

Call `init()` once per page. It scans the document for the data attributes below
and returns a cleanup function that unbinds everything it attached.

## What it handles

- **Dialogs and drawers** — backdrop dismissal, Escape, scroll locking, focus
  trap and focus restoration to the element that opened them.
- **Dropdown menus** — arrow, Home and End navigation, roving focus, dismissal
  on outside click.
- **Popovers** — native popover API with viewport-aware positioning.
- **Tabs** — roving tabindex, automatic activation, horizontal and vertical
  orientation.
- **Toasts** — polite live region, a four-toast cap, pause on hover or focus, and
  a sticky option.
- **Command palette** — Ctrl/Cmd+K, multi-term filtering, active descendant and
  events.
- **Theme switching** — light, dark and system, with `localStorage` and
  `prefers-color-scheme`.
- **Focus helpers** — `focusable`, `focusFirst` and `trapFocus` for your own
  widgets.

## Markup

| Feature | Markup |
| --- | --- |
| Dialog | `<dialog class="fui-dialog" data-fui-dialog>`, opened with `data-fui-open="#id"` and closed with `data-fui-close` |
| Drawer | `<dialog class="fui-drawer" data-fui-drawer data-fui-side="right">` |
| Dropdown | `data-fui-dropdown="#menu"` on the trigger, `role="menu"` on the items |
| Popover | `data-fui-popover="#panel"` on the trigger, `popover` on the panel |
| Tabs | `data-fui-tabs` with `role="tablist"`, `role="tab"` and `role="tabpanel"` |
| Toast | `toast({ title, description, variant, duration })` |
| Command palette | `data-fui-command` plus the input, list and item attributes |
| Theme | `data-fui-theme-set="light\|dark\|system"` or `data-fui-theme-toggle` |

## Per-feature entry points

Importing the whole runtime is convenient; importing one module keeps the rest
out of your bundle. Every feature has its own entry:

```js
import { initDialogs } from "flexui-runtime/dialog";
import { initDrawers } from "flexui-runtime/drawer";
import { initDropdowns } from "flexui-runtime/dropdown";
import { initPopovers } from "flexui-runtime/popover";
import { initTabs } from "flexui-runtime/tabs";
import { toast } from "flexui-runtime/toast";
import { initCommand } from "flexui-runtime/command";
import { setTheme, getTheme, initTheme } from "flexui-runtime/theme";
import { focusFirst, trapFocus } from "flexui-runtime/focus";
```

`flexui-runtime/auto` is a side-effect entry that initialises everything as soon
as the DOM is ready.

## Events

Every event bubbles and is cancelable.

| Event | Detail |
| --- | --- |
| `fui:open` / `fui:close` | `{ dialog, reason }` or `{ menu, trigger }` |
| `fui:select` | `{ item, menu, trigger }` |
| `fui:tabchange` | `{ tabs, tab, panel, index }` |
| `fui:popoveropen` / `fui:popoverclose` | `{ trigger, panel }` |
| `fui:toast` / `fui:toastdismiss` | `{ toast, variant }` |
| `fui:command` | `{ value, label, item, palette }` |
| `fui:themechange` | `{ mode, resolved }` |

## Theming

```js
import { setTheme, getTheme } from "flexui-runtime/theme";

setTheme("dark");          // "light" | "dark" | "system"
getTheme();                // { mode: "dark", resolved: "dark" }
```

Theme state is stored under a configurable key, so several sites on one origin
can keep independent preferences.

## Documentation

The full guide, including every option and the complete event list, lives in
[docs/runtime.md](https://github.com/Ahaduzzamankhan/flexui/blob/main/docs/runtime.md).

## License

MIT
