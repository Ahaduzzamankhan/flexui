# Runtime

`@flexui/runtime` is the optional zero-dependency browser runtime for FlexUI's
interactive components. The component CSS works without it; the runtime adds
keyboard behaviour, focus management, ARIA state and events on top of native
platform features.

- No dependencies, no build step for consumers, ESM only.
- Import the whole runtime or a single module (`@flexui/runtime/tabs`).
- Progressive enhancement: dialogs use `<dialog>` + `showModal()`, popovers use
  the native popover API, and every controller has a fallback.
- Strict TypeScript types are shipped with every module.

## Install

```bash
npm install @flexui/runtime
```

Or copy `packages/runtime/dist` into your project. The runtime has no package
dependencies.

## Initialise

```ts
import { init } from "@flexui/runtime";

const dispose = init(); // scans document
// later, if the page is a SPA view being torn down:
dispose();
```

`init(root?)` initialises theme switching, dialogs/drawers, dropdowns, tabs,
popovers, command palettes and toasts found under `root`. It returns a cleanup
function and is safe to call once per page. Re-initialising a controller that is
already bound is a no-op.

You can also initialise individual features:

```ts
import { initDialogs } from "@flexui/runtime/dialog";
import { initTabs } from "@flexui/runtime/tabs";
```

## Markup contracts

| Feature | Markup |
| --- | --- |
| Dialog | `<dialog class="fui-dialog" data-fui-dialog>`; open with `data-fui-open="#id"`, close with `data-fui-close` |
| Drawer | `<dialog class="fui-drawer" data-fui-drawer data-fui-side="right">` |
| Dropdown | Trigger `data-fui-dropdown="#menu"`; menu `role="menu"` with `role="menuitem"` items |
| Popover | Trigger `data-fui-popover="#panel"`; panel `class="fui-popover" popover` |
| Tabs | Container `data-fui-tabs`; `role="tablist"`, `role="tab"`, `role="tabpanel"` |
| Toast | Call `toast({...})`, or mark static markup `data-fui-toast` with an optional `data-fui-toast-duration` |
| Command palette | Container `data-fui-command`; input `data-fui-command-input`; list `data-fui-command-list`; items `data-fui-command-item` |
| Theme | `data-fui-theme-set="light\|dark\|system"`, or `data-fui-theme-toggle` |

All triggers accept either an id (`#menu`) or an element reference.

## Dialogs and drawers

```html
<button class="fui-btn fui-btn--primary" data-fui-open="#invite">Invite</button>

<dialog id="invite" class="fui-dialog" data-fui-dialog aria-labelledby="invite-title">
  <div class="fui-dialog-header">
    <h2 class="fui-dialog-title" id="invite-title">Invite a teammate</h2>
    <button type="button" class="fui-dialog-close" data-fui-close aria-label="Close">×</button>
  </div>
  <div class="fui-dialog-body">
    <label class="fui-label" for="invite-email">Email</label>
    <input id="invite-email" class="fui-input" type="email" autofocus />
  </div>
  <div class="fui-dialog-footer">
    <button type="button" class="fui-btn fui-btn--ghost" data-fui-close>Cancel</button>
    <button type="submit" class="fui-btn fui-btn--primary">Send invite</button>
  </div>
</dialog>
```

`bindDialog(element, options)` returns a controller:

```ts
import { bindDialog } from "@flexui/runtime/dialog";

const dialog = bindDialog(document.querySelector("#invite")!, {
  closeOnBackdrop: true, // default
  restoreFocus: true,    // default
  onClose(el, reason) {
    console.log("closed because", reason);
  },
});

dialog.open();
dialog.close("saved");
dialog.isOpen;
dialog.dispose();
```

A drawer is a dialog with `class="fui-drawer"` and `data-fui-drawer`. Use
`data-fui-side="left"` for a left edge. Focus trapping, Escape, backdrop
dismissal, scroll locking (`html[data-fui-overlay-open]`) and focus restoration
are shared with dialogs.

## Dropdown menus

```html
<button class="fui-btn fui-btn--outline" data-fui-dropdown="#account-menu">Account</button>
<div id="account-menu" class="fui-menu" role="menu" hidden>
  <p class="fui-menu-label">Signed in</p>
  <button type="button" class="fui-menu-item" role="menuitem">Profile</button>
  <hr class="fui-menu-sep" />
  <button type="button" class="fui-menu-item" role="menuitem">Sign out</button>
</div>
```

Supported keys: Enter/Space/ArrowDown open at the first item, ArrowUp opens at
the last, ArrowUp/ArrowDown/Home/End move, Escape closes and restores focus, and
Tab closes without stealing focus. `aria-haspopup`, `aria-controls` and
`aria-expanded` are managed automatically.

## Popovers

```html
<button class="fui-btn fui-btn--outline" data-fui-popover="#shipping">Shipping</button>
<div id="shipping" class="fui-popover" popover>
  <p>Orders over $50 ship free.</p>
  <button type="button" class="fui-btn fui-btn--sm" data-fui-popover-close>Got it</button>
</div>
```

The controller calls `showPopover()`/`hidePopover()`, positions the panel against
the trigger with viewport flipping, and closes on Escape or an outside pointer
down. Browsers without the popover API use a `hidden` fallback with the same
behaviour.

## Tabs

```html
<div class="fui-tabs" data-fui-tabs>
  <div class="fui-tablist" role="tablist" aria-label="Settings">
    <button class="fui-tab" role="tab" id="tab-general" aria-controls="panel-general" aria-selected="true">General</button>
    <button class="fui-tab" role="tab" id="tab-billing" aria-controls="panel-billing" aria-selected="false">Billing</button>
  </div>
  <div class="fui-tabpanel" role="tabpanel" id="panel-general" aria-labelledby="tab-general">…</div>
  <div class="fui-tabpanel" role="tabpanel" id="panel-billing" aria-labelledby="tab-billing" hidden>…</div>
</div>
```

Roving tabindex, automatic activation, Arrow keys, Home/End and vertical
orientation (`aria-orientation="vertical"`) are handled. Panels are linked with
`aria-controls`/`aria-labelledby`; missing ids are generated.

## Toasts

```ts
import { toast } from "@flexui/runtime/toast";

const handle = toast({
  variant: "success", // "info" | "success" | "danger"
  title: "Workspace created",
  description: "The new workspace is ready.",
  duration: 5000,     // 0 keeps it open
});

handle.dismiss();
```

Toasts are appended to a polite live region (created on first use), capped at
four, and pause their timer while hovered or focused. A close button is added
unless `dismissible: false`. Static markup can opt in with `data-fui-toast` and
`data-fui-toast-duration`.

## Command palette

```html
<button class="fui-btn" data-fui-command-open="#palette">Commands</button>

<dialog id="palette" class="fui-dialog" data-fui-dialog data-fui-command aria-label="Command palette">
  <div class="fui-cmd">
    <input class="fui-cmd-input" data-fui-command-input placeholder="Search commands…" />
    <div class="fui-cmd-list" data-fui-command-list role="listbox">
      <button type="button" class="fui-cmd-item" role="option" data-fui-command-item
        data-value="new" data-keywords="create workspace">New workspace</button>
    </div>
    <p class="fui-cmd-empty" data-fui-command-empty hidden>No commands found.</p>
  </div>
</dialog>
```

Ctrl/Cmd+K toggles the first palette on the page. Filtering matches every
whitespace-separated term against the label, `data-value` and `data-keywords`.
Arrow keys, Home/End, Enter and Escape work; the active item is exposed through
`aria-activedescendant`. Choosing an item dispatches `fui:command`.

```ts
document.addEventListener("fui:command", (event) => {
  const { value } = (event as CustomEvent<{ value: string }>).detail;
  runCommand(value);
});
```

## Theme switching

```ts
import { getTheme, setTheme, toggleTheme, getResolvedTheme } from "@flexui/runtime/theme";

setTheme("dark");   // persists to localStorage by default
setTheme("system", { persist: false });
toggleTheme();      // light -> dark -> system -> light
getTheme();         // "light" | "dark" | "system"
getResolvedTheme(); // "light" | "dark"
```

`init()` applies the stored mode before first paint is expected to have happened
(use a tiny inline script if you need to avoid a flash entirely), toggles the
`dark` class on `<html>`, sets `data-fui-theme-mode` and `color-scheme`, and
keeps `aria-pressed` on theme controls in sync. A `prefers-color-scheme` change
reapplies the theme while the mode is `system`.

## Focus helpers

```ts
import { focusable, focusFirst, trapFocus, FOCUSABLE_SELECTOR } from "@flexui/runtime/focus";

focusFirst(panel);          // autofocus, first focusable, or the container
const release = trapFocus(panel); // Tab / Shift+Tab cycling
release();
```

`focusable()` excludes hidden, `inert`, `aria-hidden` and non-rendered nodes.

## Events

Every event bubbles and is cancelable unless noted.

| Event | Dispatched on | Detail |
| --- | --- | --- |
| `fui:open` | dialog | `{ dialog }` |
| `fui:close` | dialog | `{ dialog, reason }` |
| `fui:open` / `fui:close` | dropdown trigger | `{ menu, trigger }` |
| `fui:select` | menu item | `{ item, menu, trigger }` |
| `fui:tabchange` | tabs container | `{ tabs, tab, panel, index }` |
| `fui:popoveropen` / `fui:popoverclose` | popover trigger | `{ trigger, panel }` |
| `fui:toast` | toast element | `{ toast, variant }` |
| `fui:toastdismiss` | toast element | `{ toast }` |
| `fui:command` | command item | `{ value, label, item, palette }` |
| `fui:commandopen` / `fui:commandclose` | palette | `{ palette }` |
| `fui:themechange` | document | `{ mode, resolved }` |

Close reasons include `escape`, `backdrop`, `trigger`, `programmatic`, `submit`,
`save`, `command` and `dispose`.

## Accessibility notes

- Native `<dialog>` provides the top layer and inert background; the runtime
  adds scroll locking and focus restoration.
- Menus and tabs manage roving focus and expose state through ARIA attributes.
- Reduced-motion preferences are respected by the component CSS.
- Nothing depends on pointer input: every action has a keyboard path.

## Building from source

```bash
bun run build:runtime   # tsc -p packages/runtime/tsconfig.json
bun run typecheck
```
