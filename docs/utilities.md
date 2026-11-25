# Utility reference

Every utility is resolved at build time. Values come from the theme scales
below or from an arbitrary value in brackets.

## Theme scales

| Scale | Tokens |
| --- | --- |
| Spacing | `0`, `px`, `0.5`, `1`, `1.5`, `2`, `2.5`, `3`, `3.5`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `14`, `16`, `20`, `24`, `28`, `32`, `36`, `40`, `44`, `48`, `52`, `56`, `60`, `64`, `72`, `80`, `96` |
| Colours | `slate`, `gray`, `zinc`, `neutral`, `stone`, `red`, `orange`, `amber`, `yellow`, `lime`, `green`, `emerald`, `teal`, `cyan`, `sky`, `blue`, `indigo`, `violet`, `purple`, `fuchsia`, `pink`, `rose` × `50`–`950`, plus `white`, `black`, `transparent`, `current`, `inherit` |
| Font size | `xs`, `sm`, `base`, `lg`, `xl`, `2xl`, `3xl`, `4xl`, `5xl`, `6xl`, `7xl`, `8xl`, `9xl`, `display`, `display-xl` |
| Font weight | `thin`, `extralight`, `light`, `normal`, `medium`, `semibold`, `bold`, `extrabold`, `black` |
| Line height | `none`, `tight`, `snug`, `normal`, `relaxed`, `loose`, `3`–`10` |
| Letter spacing | `tighter`, `tight`, `normal`, `wide`, `wider`, `widest` |
| Radius | `none`, `xs`, `sm`, `md`, `lg`, `xl`, `2xl`, `3xl`, `full` |
| Shadow | `xs`, `sm`, `md`, `lg`, `xl`, `2xl`, `inner`, `none` |
| Opacity | `0`–`100` in steps of 5 |
| z-index | `0`, `10`, `20`, `30`, `40`, `50`, `auto`, `min`, `max`, `dialog`, `toast` |
| Breakpoints | `sm` 640px, `md` 768px, `lg` 1024px, `xl` 1280px, `2xl` 1536px |

## Spacing

| Utility | CSS |
| --- | --- |
| `p-4` | `padding: 1rem` |
| `px-6` | `padding-inline: 1.5rem` |
| `py-2` | `padding-block: 0.5rem` |
| `pt-*`, `pr-*`, `pb-*`, `pl-*` | physical side padding |
| `ps-*`, `pe-*` | logical start/end padding |
| `m-4`, `mx-auto`, `my-0`, `mt-*`, `mr-*`, `mb-*`, `ml-*`, `ms-*`, `me-*` | margin equivalents |
| `gap-4`, `gap-x-2`, `gap-y-2` | `gap`, `column-gap`, `row-gap` |
| `space-x-4`, `space-y-4` | margin between siblings (`& > * + *`) |
| `-mt-4` | negative values work on every spacing utility |

## Layout

`block`, `inline-block`, `inline`, `flex`, `inline-flex`, `grid`, `inline-grid`,
`table`, `table-row`, `table-cell`, `contents`, `hidden`, `flow-root`,
`list-item`, `static`, `relative`, `absolute`, `fixed`, `sticky`, `visible`,
`invisible`, `collapse`, `isolate`, `isolation-auto`, `box-border`,
`box-content`, `sr-only`, `not-sr-only`, `container`.

Overflow: `overflow-auto`, `overflow-hidden`, `overflow-clip`, `overflow-visible`,
`overflow-scroll`, and the same set for `overflow-x-*` and `overflow-y-*`.

## Flexbox

`flex-row`, `flex-row-reverse`, `flex-col`, `flex-col-reverse`, `flex-wrap`,
`flex-wrap-reverse`, `flex-nowrap`, `flex-1`, `flex-auto`, `flex-initial`,
`flex-none`, `flex-[2_2_0%]`, `basis-*`, `grow`, `grow-0`, `shrink`, `shrink-0`,
`order-*`, `order-first`, `order-last`, `items-start|end|center|baseline|stretch`,
`items-baseline`, `justify-start|end|center|between|around|evenly`,
`content-*`, `self-*`, `justify-self-*`, `justify-items-*`, `place-items-*`,
`place-content-*`.

## Grid

`grid-cols-<n>` → `repeat(n, minmax(0, 1fr))`, `grid-cols-none`,
`grid-cols-subgrid`, `grid-cols-[250px_1fr]`, `grid-rows-<n>`, `grid-flow-row`,
`grid-flow-col`, `grid-flow-dense`, `auto-cols-fr`, `auto-rows-*`,
`col-span-<n>`, `col-span-full`, `col-start-<n>`, `col-end-<n>`, `row-span-<n>`,
`row-start-<n>`, `row-end-<n>`.

```html
<div class="grid grid-cols-1 gap-4 md:grid-cols-3">
  <div class="md:col-span-2 md:row-span-2">hero</div>
</div>
```

## Sizing

`w-*`, `h-*`, `min-w-*`, `min-h-*`, `max-w-*`, `max-h-*`, `size-*`.
Tokens: spacing scale, fractions (`w-1/2`), `auto`, `full`, `screen`, `svw`,
`lvw`, `dvw`, `svh`, `lvh`, `dvh`, `min`, `max`, `fit`, arbitrary values.
`max-w-*` also accepts `prose`, `screen-sm` … `screen-2xl` and the container
sizes `xs`–`7xl`. `aspect-auto|square|video|photo|portrait|wide`,
`aspect-[16/9]`.

## Position

`inset-*`, `inset-x-*`, `inset-y-*`, `top-*`, `right-*`, `bottom-*`, `left-*`,
`start-*`, `end-*`, `z-*`, plus negative variants.

## Colour

`bg-*` (colour, gradient direction, position, size, repeat, attachment, clip,
origin), `text-*` (colour), `border-*` (colour), `divide-*` (colour),
`outline-*`, `ring-*`, `accent-*`, `caret-*`, `fill-*`, `stroke-*`,
`placeholder:text-*`.

Opacity modifiers use `color-mix`, so they stay real colours:

```html
<div class="bg-slate-900/70">…</div>
<!-- background-color: color-mix(in srgb, #0f172a 70%, transparent) -->
```

Gradients:

```html
<div class="bg-gradient-to-r from-indigo-500 via-purple-500 to-sky-400">…</div>
```

`bg-gradient-to-{t,tr,r,br,b,bl,l,tl}` sets the image; `from-*`, `via-*` and
`to-*` set the stops through `--flexui-gradient-*` custom properties.

## Typography

`text-<size>` (sets size, line height and letter spacing from the scale),
`text-<colour>`, `text-left|center|right|justify|start|end`,
`text-wrap|nowrap|balance|pretty|stable`, `text-ellipsis`, `text-clip`,
`font-<family>`, `font-<weight>`, `leading-*`, `tracking-*`, `indent-*`,
`align-*` (vertical), `whitespace-*`, `break-normal|words|all|keep`,
`hyphens-*`, `line-clamp-<n>`, `list-<style>`, `decoration-*`,
`underline-offset-<n>`, `italic`, `not-italic`, `underline`, `overline`,
`line-through`, `no-underline`, `uppercase`, `lowercase`, `capitalize`,
`normal-case`, `truncate`, `antialiased`, `subpixel-antialiased`,
`tabular-nums`, `proportional-nums`, `font-feature-normal`, `feat-[...]`.

Semantic typography:

| Class | Scale |
| --- | --- |
| `type-display-xl` | `display-xl`, bold |
| `type-display` | `display`, bold |
| `type-heading-xl` | `4xl`, semibold |
| `type-heading-lg` | `3xl`, semibold |
| `type-heading` | `2xl`, semibold |
| `type-subheading` | `lg`, medium |
| `type-body-lg` | `lg`, normal |
| `type-body` | `base`, normal |
| `type-body-sm` | `sm`, normal |
| `type-caption` | `xs`, normal |
| `type-label` | `sm`, medium |
| `type-code` | `sm`, normal, monospace |

## Effects

`shadow-*`, `opacity-*`, `blur-*`, `backdrop-blur-*`, `brightness-*`,
`contrast-*`, `grayscale`, `grayscale-0`, `invert`, `saturate-*`, `sepia`,
`hue-rotate-*`, `drop-shadow-*`, `filter-none`, `backdrop-filter-none`,
`mix-blend-*`, `bg-blend-*`.

Shadows and rings compose: both write a `box-shadow` that reads
`--flexui-ring-width`, `--flexui-ring-color` and `--flexui-shadow`, so
`shadow-md ring-2` produces one combined shadow.

Filter utilities (`blur-*`, `grayscale-*`, …) each emit `filter` directly and
therefore replace one another; they do not compose yet.

## Borders

`border` (1px), `border-0`, `border-2`, `border-4`, `border-8`,
`border-<colour>`, `border-solid|dashed|dotted|double|hidden|none`,
`border-{x,y,t,r,b,l,s,e}-*`, `divide-x-*`, `divide-y-*`, `divide-<colour>`,
`rounded`, `rounded-{t,r,b,l,tl,tr,br,bl,ss,se,es,ee}-*`,
`outline`, `outline-<width>`, `outline-<colour>`, `outline-none`,
`outline-offset-*`, `ring`, `ring-<width>`, `ring-<colour>`, `ring-offset-*`.

## Transforms

`translate-x-*`, `translate-y-*`, `rotate-*`, `scale-*`, `scale-x-*`,
`scale-y-*`, `skew-x-*`, `skew-y-*`, `origin-*`, `transform-none`,
`perspective-*`.

Transforms compose through custom properties:

```css
.rotate-45 {
  --flexui-rotate: 45deg;
  transform: translate(var(--flexui-tx, 0), var(--flexui-ty, 0)) rotate(var(--flexui-rotate, 0)) scale(var(--flexui-sx, 1), var(--flexui-sy, 1)) skew(var(--flexui-skx, 0), var(--flexui-sky, 0));
}
```

## Interaction and transitions

`cursor-*`, `pointer-events-none|auto`, `select-none|text|all|auto`,
`resize-none|y|x|both`, `appearance-none|auto`, `object-contain|cover|fill|none|scale-down`,
`object-{position}-*`, `transition`, `transition-all|colors|opacity|shadow|transform|size|none`,
`duration-*`, `delay-*`, `ease-*`, `will-change-{transform,opacity,auto}`.

Easing presets: `linear`, `in`, `out`, `in-out`, `smooth`, `spring`, `bounce`,
`snap`.

## Arbitrary values

```html
<div class="w-[350px] h-[calc(100vh-80px)] bg-[#121212] text-[clamp(2rem,5vw,5rem)] grid-cols-[250px_1fr] shadow-[0_10px_40px_rgba(0,0,0,0.1)]"></div>
```

Underscores become spaces. Escape an underscore with `\_` to keep it literal.
A type hint resolves ambiguity between length and colour:

```html
<span class="text-[length:12px] text-[color:oklch(0.6_0.2_250)]">…</span>
```

Arbitrary properties work without a registry entry:

```html
<div class="[mask-type:luminance]">…</div>
```

`!important` is supported as a prefix (`!p-4`) or suffix (`p-4!`).
