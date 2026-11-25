// Package themes holds FlexUI's design tokens: spacing, colour palettes,
// typography scales, radii, shadows, breakpoints and easing curves.
//
// Every token scale is a plain map so that a user configuration can extend or
// override individual entries without recompiling the compiler.
package themes

import (
	"sort"
	"strconv"
)

// FontSize is one entry of the typography scale.
type FontSize struct {
	Size          string
	LineHeight    string
	LetterSpacing string
}

// Theme is the resolved set of design tokens.
type Theme struct {
	Spacing        map[string]string
	Colors         map[string]string
	FontFamily     map[string]string
	FontSize       map[string]FontSize
	FontWeight     map[string]string
	LineHeight     map[string]string
	LetterSpacing  map[string]string
	Radius         map[string]string
	Shadow         map[string]string
	Opacity        map[string]string
	ZIndex         map[string]string
	Breakpoints    map[string]string
	ContainerSizes map[string]string
	Easing         map[string]string
	Duration       map[string]string
	AspectRatio    map[string]string
	MaxWidth       map[string]string

	// SemanticLight and SemanticDark hold role tokens (surface, fg, line,
	// accent, …). They are emitted as custom properties so themes can switch
	// at runtime without recompiling.
	SemanticLight map[string]string
	SemanticDark  map[string]string
}

// Default returns the built-in FlexUI theme.
func Default() *Theme {
	return &Theme{
		Spacing: map[string]string{
			"0":   "0",
			"px":  "1px",
			"0.5": "0.125rem",
			"1":   "0.25rem",
			"1.5": "0.375rem",
			"2":   "0.5rem",
			"2.5": "0.625rem",
			"3":   "0.75rem",
			"3.5": "0.875rem",
			"4":   "1rem",
			"5":   "1.25rem",
			"6":   "1.5rem",
			"7":   "1.75rem",
			"8":   "2rem",
			"9":   "2.25rem",
			"10":  "2.5rem",
			"11":  "2.75rem",
			"12":  "3rem",
			"14":  "3.5rem",
			"16":  "4rem",
			"20":  "5rem",
			"24":  "6rem",
			"28":  "7rem",
			"32":  "8rem",
			"36":  "9rem",
			"40":  "10rem",
			"44":  "11rem",
			"48":  "12rem",
			"52":  "13rem",
			"56":  "14rem",
			"60":  "15rem",
			"64":  "16rem",
			"72":  "18rem",
			"80":  "20rem",
			"96":  "24rem",
		},
		Colors: defaultColors(),
		FontFamily: map[string]string{
			"sans":   `ui-sans-serif, system-ui, -apple-system, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif`,
			"serif":  `ui-serif, Georgia, Cambria, "Times New Roman", Times, serif`,
			"mono":   `ui-monospace, SFMono-Regular, "SF Mono", Menlo, Consolas, "Liberation Mono", monospace`,
			"system": `system-ui, -apple-system, "Segoe UI", Roboto, sans-serif`,
		},
		FontSize: map[string]FontSize{
			"xs":         {Size: "0.75rem", LineHeight: "1rem"},
			"sm":         {Size: "0.875rem", LineHeight: "1.25rem"},
			"base":       {Size: "1rem", LineHeight: "1.5rem"},
			"lg":         {Size: "1.125rem", LineHeight: "1.75rem"},
			"xl":         {Size: "1.25rem", LineHeight: "1.75rem"},
			"2xl":        {Size: "1.5rem", LineHeight: "2rem"},
			"3xl":        {Size: "1.875rem", LineHeight: "2.25rem"},
			"4xl":        {Size: "2.25rem", LineHeight: "2.5rem"},
			"5xl":        {Size: "3rem", LineHeight: "1"},
			"6xl":        {Size: "3.75rem", LineHeight: "1"},
			"7xl":        {Size: "4.5rem", LineHeight: "1"},
			"8xl":        {Size: "6rem", LineHeight: "1"},
			"9xl":        {Size: "8rem", LineHeight: "1"},
			"display":    {Size: "clamp(2rem, 4vw, 3.5rem)", LineHeight: "1.05", LetterSpacing: "-0.02em"},
			"display-xl": {Size: "clamp(2.5rem, 6vw, 6rem)", LineHeight: "0.95", LetterSpacing: "-0.04em"},
		},
		FontWeight: map[string]string{
			"thin":       "100",
			"extralight": "200",
			"light":      "300",
			"normal":     "400",
			"medium":     "500",
			"semibold":   "600",
			"bold":       "700",
			"extrabold":  "800",
			"black":      "900",
		},
		LineHeight: map[string]string{
			"none":    "1",
			"tight":   "1.25",
			"snug":    "1.375",
			"normal":  "1.5",
			"relaxed": "1.625",
			"loose":   "2",
			"3":       "0.75rem",
			"4":       "1rem",
			"5":       "1.25rem",
			"6":       "1.5rem",
			"7":       "1.75rem",
			"8":       "2rem",
			"9":       "2.25rem",
			"10":      "2.5rem",
		},
		LetterSpacing: map[string]string{
			"tighter": "-0.05em",
			"tight":   "-0.025em",
			"normal":  "0em",
			"wide":    "0.025em",
			"wider":   "0.05em",
			"widest":  "0.1em",
		},
		Radius: map[string]string{
			"none": "0",
			"xs":   "0.125rem",
			"sm":   "0.25rem",
			"":     "0.375rem",
			"md":   "0.5rem",
			"lg":   "0.75rem",
			"xl":   "1rem",
			"2xl":  "1.5rem",
			"3xl":  "2rem",
			"full": "9999px",
		},
		Shadow: map[string]string{
			"none":  "none",
			"xs":    "0 1px 2px 0 rgb(0 0 0 / 0.05)",
			"sm":    "0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1)",
			"":      "0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1)",
			"md":    "0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)",
			"lg":    "0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1)",
			"xl":    "0 20px 25px -5px rgb(0 0 0 / 0.1), 0 8px 10px -6px rgb(0 0 0 / 0.1)",
			"2xl":   "0 25px 50px -12px rgb(0 0 0 / 0.25)",
			"inner": "inset 0 2px 4px 0 rgb(0 0 0 / 0.05)",
		},
		Opacity: map[string]string{
			"0": "0", "5": "0.05", "10": "0.1", "15": "0.15", "20": "0.2",
			"25": "0.25", "30": "0.3", "35": "0.35", "40": "0.4", "45": "0.45",
			"50": "0.5", "55": "0.55", "60": "0.6", "65": "0.65", "70": "0.7",
			"75": "0.75", "80": "0.8", "85": "0.85", "90": "0.9", "95": "0.95",
			"100": "1",
		},
		ZIndex: map[string]string{
			"0": "0", "10": "10", "20": "20", "30": "30", "40": "40", "50": "50",
			"auto":   "auto",
			"min":    "-1",
			"max":    "9999",
			"dialog": "100",
			"toast":  "200",
		},
		Breakpoints: map[string]string{
			"sm":  "640px",
			"md":  "768px",
			"lg":  "1024px",
			"xl":  "1280px",
			"2xl": "1536px",
		},
		ContainerSizes: map[string]string{
			"xs":  "20rem",
			"sm":  "24rem",
			"md":  "28rem",
			"lg":  "32rem",
			"xl":  "36rem",
			"2xl": "42rem",
			"3xl": "48rem",
			"4xl": "56rem",
			"5xl": "64rem",
			"6xl": "72rem",
			"7xl": "80rem",
		},
		Easing: map[string]string{
			"linear": "linear",
			"in":     "cubic-bezier(0.4, 0, 1, 1)",
			"out":    "cubic-bezier(0, 0, 0.2, 1)",
			"in-out": "cubic-bezier(0.4, 0, 0.2, 1)",
			"smooth": "cubic-bezier(0.22, 1, 0.36, 1)",
			"spring": "cubic-bezier(0.34, 1.56, 0.64, 1)",
			"bounce": "cubic-bezier(0.68, -0.55, 0.27, 1.55)",
			"snap":   "cubic-bezier(0.2, 0, 0, 1)",
		},
		Duration: map[string]string{
			"75":   "75ms",
			"100":  "100ms",
			"150":  "150ms",
			"200":  "200ms",
			"300":  "300ms",
			"400":  "400ms",
			"500":  "500ms",
			"700":  "700ms",
			"1000": "1000ms",
		},
		AspectRatio: map[string]string{
			"auto":     "auto",
			"square":   "1 / 1",
			"video":    "16 / 9",
			"photo":    "3 / 2",
			"portrait": "3 / 4",
			"wide":     "21 / 9",
		},
		MaxWidth: map[string]string{
			"none":       "none",
			"xs":         "20rem",
			"sm":         "24rem",
			"md":         "28rem",
			"lg":         "32rem",
			"xl":         "36rem",
			"2xl":        "42rem",
			"3xl":        "48rem",
			"4xl":        "56rem",
			"5xl":        "64rem",
			"6xl":        "72rem",
			"7xl":        "80rem",
			"prose":      "65ch",
			"screen-sm":  "640px",
			"screen-md":  "768px",
			"screen-lg":  "1024px",
			"screen-xl":  "1280px",
			"screen-2xl": "1536px",
		},
		SemanticLight: map[string]string{
			"surface":        "#ffffff",
			"surface-muted":  "#f8fafc",
			"surface-raised": "#ffffff",
			"surface-sunken": "#f1f5f9",
			"fg":             "#0f172a",
			"fg-muted":       "#475569",
			"fg-subtle":      "#64748b",
			"fg-inverted":    "#ffffff",
			"line":           "#e2e8f0",
			"line-strong":    "#cbd5e1",
			"accent":         "#4f46e5",
			"accent-hover":   "#4338ca",
			"accent-fg":      "#ffffff",
			"accent-muted":   "#eef2ff",
			"danger":         "#dc2626",
			"danger-muted":   "#fef2f2",
			"success":        "#16a34a",
			"success-muted":  "#f0fdf4",
			"warning":        "#d97706",
			"warning-muted":  "#fffbeb",
			"overlay":        "rgb(15 23 42 / 0.45)",
			"focus-ring":     "#6366f1",
		},
		SemanticDark: map[string]string{
			"surface":        "#0b1220",
			"surface-muted":  "#111a2e",
			"surface-raised": "#16213a",
			"surface-sunken": "#070d18",
			"fg":             "#e2e8f0",
			"fg-muted":       "#94a3b8",
			"fg-subtle":      "#7c8ba3",
			"fg-inverted":    "#0b1220",
			"line":           "#293548",
			"line-strong":    "#3b4a63",
			"accent":         "#818cf8",
			"accent-hover":   "#a5b4fc",
			"accent-fg":      "#0b1220",
			"accent-muted":   "#1e1b4b",
			"danger":         "#f87171",
			"danger-muted":   "#450a0a",
			"success":        "#4ade80",
			"success-muted":  "#052e16",
			"warning":        "#fbbf24",
			"warning-muted":  "#451a03",
			"overlay":        "rgb(2 6 23 / 0.7)",
			"focus-ring":     "#a5b4fc",
		},
	}
}

// Color resolves a colour token such as "slate-700", "white" or
// "transparent".
func (t *Theme) Color(token string) (string, bool) {
	v, ok := t.Colors[token]
	return v, ok
}

// BreakpointNames returns the breakpoint names ordered from smallest to
// largest. Ties (equal widths) are ordered by name so output is stable.
func (t *Theme) BreakpointNames() []string {
	names := make([]string, 0, len(t.Breakpoints))
	for name := range t.Breakpoints {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		a, b := t.Breakpoints[names[i]], t.Breakpoints[names[j]]
		if a != b {
			return pxValue(a) < pxValue(b)
		}
		return names[i] < names[j]
	})
	return names
}

func pxValue(v string) int {
	n, err := strconv.ParseFloat(trimUnit(v), 64)
	if err != nil {
		return 1 << 30
	}
	return int(n)
}

func trimUnit(v string) string {
	i := 0
	for i < len(v) && (v[i] >= '0' && v[i] <= '9' || v[i] == '.' || v[i] == '-') {
		i++
	}
	return v[:i]
}
