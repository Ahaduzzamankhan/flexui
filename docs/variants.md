# Variants

Variants are prefixes on a utility and compose left to right, outermost first.

```html
<button class="md:dark:hover:bg-slate-800">…</button>
```

becomes

```css
@media (min-width: 768px) {
  .dark .md\:dark\:hover\:bg-slate-800:hover { background-color: …; }
}
```

## Responsive

`sm:` `md:` `lg:` `xl:` `2xl:` → `@media (min-width: …)`, using the theme
breakpoints (640, 768, 1024, 1280, 1536 px).

## Interaction and state

`hover:` `focus:` `focus-visible:` `focus-within:` `active:` `visited:`
`target:` `disabled:` `enabled:` `checked:` `required:` `optional:` `invalid:`
`valid:` `user-invalid:` `placeholder-shown:` `indeterminate:` `read-only:`
`read-write:` `autofill:` `default:` `open:` `fullscreen:`
`picture-in-picture:`

## Structural

`first:` `last:` `odd:` `even:` `only:` `empty:` `first-of-type:`
`last-of-type:` `only-of-type:` and their negations, for example
`not-first:` → `&:not(:first-child)`.

## Group and peer

Mark the parent with `group` and the sibling with `peer`. Both classes are
explicit markers: they emit no CSS themselves.

```html
<a class="group rounded-xl p-6 transition hover:shadow-lg">
  <span class="transition group-hover:scale-105">icon</span>
</a>

<input class="peer" type="checkbox" />
<label class="peer-checked:font-semibold">Enable</label>
```

Supported group states: `hover`, `focus`, `focus-visible`, `focus-within`,
`active`, `visited`, `disabled`, `checked`, `open`, `placeholder-shown`,
`invalid`, `valid`, `required`, `empty`, `target`.
Supported peer states: the same plus `enabled` and `indeterminate`.

Named groups are supported: `group/name` in markup and
`group-hover/name:` in the utility (emits `.group\/name:hover &`).

## Accessibility

`motion-safe:` → `@media (prefers-reduced-motion: no-preference)`
`motion-reduce:` → `@media (prefers-reduced-motion: reduce)`
`contrast-more:` → `@media (prefers-contrast: more)`
`contrast-less:` → `@media (prefers-contrast: less)`
`print:` → `@media print`
`portrait:`, `landscape:` → orientation media queries

FlexUI never disables motion globally; use `motion-reduce:` on the properties
that should stop animating.

## Theme

`dark:` prefixes an ancestor selector `.dark &`, so switching themes is a class
toggle on `<html>` and needs no recompilation:

```html
<html class="dark">
  <body class="bg-white dark:bg-slate-950">…</body>
</html>
```

`light:` re-asserts the base style with the lowest variant rank, which makes
explicit light-mode overrides predictable inside a dark subtree.

## Data, ARIA and has

```html
<div data-[state=open]:block hidden>…</div>
<button aria-[expanded=true]:bg-slate-100>Toggle</button>
<div class="has-[:checked]:ring-2">…</div>
```

## supports

```html
<div class="supports-[display:grid]:grid block">…</div>
```

## Container queries

```html
<div class="@container">
  <div class="@sm:grid-cols-2 @[40rem]:grid-cols-3 grid">…</div>
</div>
```

`@xs` … `@7xl` use the container size tokens (20rem–80rem); `@[<length>]`
accepts any length.

## Arbitrary selectors

```html
<ul class="[&>li]:py-2 [&>li:last-child]:border-b-0">
  <li>one</li>
</ul>
```

`&` marks where the element's own selector goes. A selector without `&` is
treated as an ancestor: `[.theme]` becomes `.theme &`.

## Ordering

Variant order is fixed so the cascade is predictable:

1. base utilities
2. state, arbitrary and theme variants
3. group and peer variants
4. accessibility media queries
5. breakpoints, ascending
6. container queries
7. `@supports`
8. `print`

Within one context, rules sort by utility order (shorthands before longhands),
variant rank, variant count and finally class name. The same input therefore
always produces the same stylesheet.
