# Integrations

FlexUI's compiler is a native binary, and every integration runs that same
binary. None of them reimplement compilation, so a Vite build, a CLI build and a
framework helper produce byte-identical CSS.

## Vite

```bash
npm install -D flexui-vite flexui-compiler
```

```ts
// vite.config.ts
import { defineConfig } from "vite";
import { flexui } from "flexui-vite";

export default defineConfig({
  plugins: [
    flexui({
      config: "flexui.config.json",
      output: "dist/flexui.css",
      // minify defaults to Vite's build.minify setting
    }),
  ],
});
```

Link the generated stylesheet from your HTML entry:

```html
<link rel="stylesheet" href="/dist/flexui.css" />
```

| Option | Default | Meaning |
| --- | --- | --- |
| `root` | Vite's `root` | Project directory passed to the compiler |
| `config` | compiler default | Config file path |
| `input` | config value | Entry CSS with `@flexui` directives |
| `output` | config value | Stylesheet path |
| `content` | config value | Globs to scan |
| `safelist` | config value | Classes that must always be emitted |
| `minify` | Vite's `build.minify` | Force minification on or off |
| `watch` | `true` | Run `flexui watch` during `vite dev` |
| `command` | `FLEXUI_BINARY` → local `.bin` → `PATH` | Compiler binary |

- `vite build` compiles once in `buildStart`; a compiler failure is reported as a
  Vite error with the compiler's own diagnostic.
- `vite dev` starts `flexui watch` and forwards its output to the Vite logger;
  the watcher is stopped when the server closes.
- The plugin has no dependencies; `vite` is an optional peer dependency.

## React

```bash
npm install flexui-runtime flexui-react
```

```tsx
import { useDialog, useTheme } from "flexui-react";

export function Invite() {
  const dialog = useDialog({ closeOnBackdrop: true });
  const theme = useTheme();

  return (
    <>
      <button className="fui-btn fui-btn--primary" onClick={() => dialog.open()}>
        Invite
      </button>
      <dialog ref={dialog.ref} className="fui-dialog" data-fui-dialog>
        <div className="fui-dialog-body">…</div>
      </dialog>
      <button className="fui-btn fui-btn--ghost" onClick={theme.toggle}>
        {theme.resolved === "dark" ? "Light" : "Dark"}
      </button>
    </>
  );
}
```

Available hooks: `useDialog`, `useDropdown`, `usePopover`, `useTabs`, `useTheme`.
Each binds on mount, disposes on unmount, and mirrors controller state into React
state, so a dialog closed with Escape re-renders. `react` is a peer dependency.

## Vue

```vue
<script setup lang="ts">
import { useDialog } from "flexui-vue";
const invite = useDialog({ closeOnBackdrop: true });
</script>

<template>
  <button class="fui-btn fui-btn--primary" @click="invite.open()">Invite</button>
  <dialog ref="invite.element" class="fui-dialog" data-fui-dialog>
    <div class="fui-dialog-body">…</div>
  </dialog>
</template>
```

Composables: `useDialog`, `useDropdown`, `usePopover`, `useTabs`, `useTheme`,
each returning a template ref, the controller (shallow ref) and mirrored state.
`vue` is a peer dependency.

## Svelte

```svelte
<script>
  import { dialog, dropdown, tabs } from "flexui-svelte";
</script>

<dialog class="fui-dialog" use:dialog={{ closeOnBackdrop: true }}>…</dialog>
<button class="fui-btn" use:dropdown={"#account-menu"}>Account</button>
<div class="fui-tabs" data-fui-tabs use:tabs>…</div>
```

Actions: `dialog`, `dropdown`, `popover`, `tabs`. They import nothing from
Svelte, so they work in Svelte 4 and 5 and need no build step; parameters are
read when the action runs.

## Without a framework

```ts
import { init } from "flexui-runtime";
init();
```

or drop in the zero-config entry:

```html
<script type="module" src="/node_modules/flexui-runtime/dist/auto.js"></script>
```

That initialises every `data-fui-*` feature on the page and exposes
`window.flexui.start()` / `stop()`.

## Testing

The runtime has DOM tests (`bun test packages/runtime/test`) covering dialogs,
dropdowns, tabs, toasts, theme persistence, focus trapping and overlay scroll
locking. The Vite plugin has tests for argument construction, error formatting,
binary resolution and hook wiring (`bun test packages/vite-plugin/test`). The
framework packages are typechecked against their real framework type definitions;
they are thin wrappers over the tested runtime controllers.
