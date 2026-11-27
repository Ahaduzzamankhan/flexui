# flexui-react

React hooks for FlexUI's interactive components. They bind the
[flexui-runtime](https://www.npmjs.com/package/flexui-runtime) controllers to
your components, dispose them on unmount, and mirror state such as whether a
dialog is open into React so your UI re-renders when a dialog closes from the
keyboard or from a backdrop click.

```bash
npm install flexui-runtime flexui-react
```

`react` (19 or newer) is a peer dependency.

## Dialog

```tsx
import { useDialog } from "flexui-react";

export function Settings() {
  const { ref, open, close, isOpen } = useDialog({ closeOnBackdrop: true });

  return (
    <>
      <button className="fui-btn" onClick={open}>Settings</button>
      <dialog ref={ref} className="fui-dialog" data-fui-dialog>
        <h2>Settings</h2>
        <button className="fui-btn" onClick={close}>Done</button>
      </dialog>
    </>
  );
}
```

`useDialog(options)` returns `ref`, `open()`, `close()`, `toggle()` and `isOpen`.
Options are read when the dialog is first bound.

## Dropdown

```tsx
const { ref, isOpen, open, close } = useDropdown("#account-menu");

<button ref={ref} className="fui-btn" aria-haspopup="menu">Account</button>
<div id="account-menu" role="menu">…</div>
```

`useDropdown(target, options)` accepts a selector or the menu element itself and
returns the trigger `ref` plus the open state.

## Popover

```tsx
const { ref, isOpen } = usePopover("#filters");

<button ref={ref} className="fui-btn">Filters</button>
<div id="filters" popover>…</div>
```

## Tabs

```tsx
const { ref, activeIndex, select } = useTabs();

<div ref={ref} className="fui-tabs" data-fui-tabs>
  <div role="tablist">…</div>
</div>
```

The hook binds the tablist, keeps roving tabindex working and tracks the active
tab in state.

## Theme

```tsx
const { mode, resolved, setMode } = useTheme();

<button className="fui-btn" onClick={() => setMode(mode === "dark" ? "light" : "dark")}>
  {resolved === "dark" ? "Light mode" : "Dark mode"}
</button>
```

`useTheme()` reads the current mode, re-renders on `fui:themechange` and exposes
`setMode("light" | "dark" | "system")`.

## Styling

The hooks handle behaviour only. Class names come from the FlexUI compiler, so
the same markup works whether or not the hooks are mounted.

## Documentation

See the [FlexUI repository](https://github.com/Ahaduzzamankhan/flexui#readme) and
[docs/runtime.md](https://github.com/Ahaduzzamankhan/flexui/blob/main/docs/runtime.md)
for the underlying controllers, markup contracts and events.

## License

MIT
