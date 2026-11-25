# Getting started

## 1. Install

```bash
npm install flexui
```

FlexUI is a native compiler. npm installs the binary for your platform through
an optional dependency package (`@flexui/cli-<os>-<arch>`), so no Go toolchain
and no install scripts are required. If optional dependencies are disabled in
your environment, install the platform package directly:

```bash
npm install @flexui/cli-linux-x64      # or win32-x64, darwin-arm64, ...
```

You can always bypass resolution with `FLEXUI_BINARY=/path/to/flexui`.

## 2. Initialise

```bash
npx flexui init
```

This creates `flexui.config.json` and `src/styles.css`:

```json
{
  "content": ["./index.html", "./src/**/*.{html,js,jsx,ts,tsx,vue,svelte,astro}"],
  "output": "dist/flexui.css",
  "minify": false,
  "safelist": []
}
```

```css
@flexui base;
@flexui utilities;
```

`@flexui base` emits the preflight reset. `@flexui utilities` emits the resolved
utilities. Anything else in the entry file is copied into the output after the
generated layers, so your own CSS always wins.

Existing files are never overwritten.

## 3. Write markup

```html
<link rel="stylesheet" href="dist/flexui.css">

<div class="mx-auto max-w-3xl px-6 py-12">
  <h1 class="type-display text-balance">Ship interfaces, not stylesheets.</h1>
  <p class="mt-4 text-lg leading-relaxed text-slate-600">
    Utilities are resolved at build time and emitted once.
  </p>
  <div class="mt-8 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
    <article class="rounded-2xl border border-slate-200 p-6 transition hover:shadow-md">
      ...
    </article>
  </div>
</div>
```

## 4. Build

```bash
npx flexui build            # development build
npx flexui build --minify   # production build
npx flexui watch            # rebuild on change
```

The build prints what it did:

```
flexui 2.0.0-alpha.1
  dist/flexui.css
  1 files -> 42 candidates -> 39 rules -> 18.6 kB
  24ms
```

## 5. Inspect

`npx flexui inspect` prints every resolved rule, ordered exactly as it will be
emitted, without writing a file. This is the fastest way to understand why a
class behaves the way it does:

```
.p-4 { padding: 1rem; }
@media (min-width: 768px) .md\:flex { display: flex; }
```

## 6. Components and the optional runtime

Component classes work like utilities: use them in markup and FlexUI emits the
CSS you used.

```html
<button class="fui-btn fui-btn--primary" data-fui-open="#invite">Invite</button>

<dialog id="invite" class="fui-dialog" data-fui-dialog aria-labelledby="invite-title">
  <div class="fui-dialog-header">
    <h2 class="fui-dialog-title" id="invite-title">Invite a teammate</h2>
    <button class="fui-dialog-close" data-fui-close aria-label="Close">×</button>
  </div>
  <div class="fui-dialog-body">…</div>
</dialog>
```

For keyboard behaviour, focus restoration and theme switching, add the optional
runtime:

```bash
npm install @flexui/runtime
```

```ts
import { init } from "@flexui/runtime";
init();
```

Nothing breaks without it: dialogs are native `<dialog>` elements, menus are
lists of buttons and tabs are standard ARIA markup. See
[components](components.md) and [runtime](runtime.md).

## Where to go next

- [Components](components.md)
- [Runtime](runtime.md)
- [CLI reference](cli.md)
- [Utility reference](utilities.md)
- [Variants](variants.md)
- [Architecture](architecture.md)
