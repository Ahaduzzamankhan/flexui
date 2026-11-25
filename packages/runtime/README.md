# @flexui/runtime

Optional zero-dependency browser runtime for [FlexUI](https://github.com/Ahaduzzamankhan/flexui) components. The CSS works without it; this package adds keyboard behaviour, focus management, ARIA state and DOM events on top of native `<dialog>` and the popover API.

```bash
npm install @flexui/runtime
```

```ts
import { init } from "@flexui/runtime";

const dispose = init();
```

## Features

- Dialogs and drawers: backdrop dismissal, Escape, scroll locking, focus trap and focus restoration.
- Dropdown menus: arrow/Home/End navigation, roving focus, outside dismissal.
- Popovers: native popover API with viewport-aware positioning and a fallback.
- Tabs: roving tabindex, automatic activation, horizontal/vertical orientation.
- Toasts: polite live region, four-toast cap, pause on hover/focus, sticky option.
- Command palette: Ctrl/Cmd+K, multi-term filtering, active descendant, events.
- Theme switching: light/dark/system, `localStorage`, `prefers-color-scheme`.
- Focus helpers: `focusable`, `focusFirst`, `trapFocus`.

## Markup contracts

| Feature | Markup |
| --- | --- |
| Dialog | `<dialog class="fui-dialog" data-fui-dialog>` with `data-fui-open="#id"` / `data-fui-close` |
| Drawer | `<dialog class="fui-drawer" data-fui-drawer data-fui-side="right">` |
| Dropdown | `data-fui-dropdown="#menu"` + `role="menu"` items |
| Popover | `data-fui-popover="#panel"` + `popover` panel |
| Tabs | `data-fui-tabs` + `role="tablist"` / `tab` / `tabpanel` |
| Toast | `toast({ title, description, variant, duration })` |
| Command palette | `data-fui-command` + input/list/item attributes |
| Theme | `data-fui-theme-set="light\|dark\|system"` or `data-fui-theme-toggle` |

`init(root?)` returns a cleanup function and is safe to call once per page. Every feature is also exported from its own entry point:

```ts
import { initDialogs } from "@flexui/runtime/dialog";
import { toast } from "@flexui/runtime/toast";
import { setTheme } from "@flexui/runtime/theme";
```

## Events

All events bubble and are cancelable.

| Event | Detail |
| --- | --- |
| `fui:open` / `fui:close` | `{ dialog, reason }` or `{ menu, trigger }` |
| `fui:select` | `{ item, menu, trigger }` |
| `fui:tabchange` | `{ tabs, tab, panel, index }` |
| `fui:popoveropen` / `fui:popoverclose` | `{ trigger, panel }` |
| `fui:toast` / `fui:toastdismiss` | `{ toast, variant }` |
| `fui:command` | `{ value, label, item, palette }` |
| `fui:themechange` | `{ mode, resolved }` |

## Docs

The full guide lives in [docs/runtime.md](https://github.com/Ahaduzzamankhan/flexui/blob/main/docs/runtime.md).

## Build

```bash
bun run build:runtime
```

Output is ESM plus declarations in `dist/`. MIT licensed.
