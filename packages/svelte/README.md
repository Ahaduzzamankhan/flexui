# flexui-svelte

Svelte actions for FlexUI's interactive components. They wrap the
[flexui-runtime](https://www.npmjs.com/package/flexui-runtime) controllers, run
when the element mounts, and clean themselves up when it is destroyed. No Svelte
imports are needed, so the same actions work with Svelte 4 and 5.

```bash
npm install flexui-runtime flexui-svelte
```

## Dialog

```svelte
<script>
  import { dialog } from "flexui-svelte";
</script>

<button class="fui-btn" on:click={() => document.querySelector("#settings").showModal()}>
  Settings
</button>

<dialog id="settings" class="fui-dialog" data-fui-dialog use:dialog={{ closeOnBackdrop: true }}>
  <h2>Settings</h2>
</dialog>
```

## Dropdown

```svelte
<script>
  import { dropdown } from "flexui-svelte";
</script>

<button class="fui-btn" use:dropdown={"#account-menu"}>Account</button>
<div id="account-menu" role="menu">…</div>
```

## Popover

```svelte
<button class="fui-btn" use:popover={"#filters"}>Filters</button>
<div id="filters" popover>…</div>
```

## Tabs

```svelte
<div class="fui-tabs" data-fui-tabs use:tabs>
  <div role="tablist">…</div>
</div>
```

Parameters are read once, when the action runs. For reactive state, the runtime
events (`fui:open`, `fui:close`, `fui:tabchange`, `fui:themechange`) give you the
same information Svelte would otherwise duplicate.

## Lower-level binding

The package also exports `bindDialog`, `bindDropdown`, `bindPopover` and
`bindTabs`, which return the underlying controller when you want to drive it
yourself.

## Documentation

See the [FlexUI repository](https://github.com/Ahaduzzamankhan/flexui#readme) and
[docs/runtime.md](https://github.com/Ahaduzzamankhan/flexui/blob/main/docs/runtime.md)
for markup contracts, options and events.

## License

MIT
