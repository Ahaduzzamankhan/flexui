# FlexUI component lab

One page that exercises every feature of [`@flexui/runtime`](../../packages/runtime):
native dialogs, an edge drawer, a dropdown menu, a positioned popover, ARIA tabs,
toasts, a Ctrl/Cmd+K command palette and persistent light/dark/system theme
switching.

Everything is styled with the FlexUI component layer and utilities. The runtime
is optional; the page also logs runtime events into a live panel so you can see
what each interaction emits.

## Build

From the repository root:

```bash
# compile the example CSS with the native compiler
cd examples/components
../../bin/flexui build

# compile the TypeScript entry
cd ../..
bun run build:example
```

`bun run build:components` builds both the runtime and this example's TypeScript.

## Preview

```bash
bun scripts/serve.ts . 4173
```

Open <http://127.0.0.1:4173/examples/components/>.

The compiled `dist/flexui.css` and `dist/main.js` are committed so the page works
from a clean checkout. Rebuild them after changing the source.

## Files

| Path | Purpose |
| --- | --- |
| `index.html` | Markup for every component and the runtime log |
| `src/styles.css` | `@flexui` entry plus small lab-specific styles |
| `src/main.ts` | Wires toasts, commands, the event log and the theme status |
| `flexui.config.json` | Content globs, output and toast class safelist |
| `tsconfig.json` | Strict browser TypeScript build |

The TypeScript imports the built runtime from `packages/runtime/dist`, so run
`bun run build:runtime` before `build:example`.
