import { createDisposer, emit, on, qsa, type Cleanup } from "./types.js";

export type ThemeMode = "light" | "dark" | "system";
export type ResolvedTheme = "light" | "dark";

export interface ThemeChangeDetail {
  mode: ThemeMode;
  resolved: ResolvedTheme;
}

export interface ThemeOptions {
  /** Root used for theme controls. Defaults to document. */
  root?: ParentNode;
  /** localStorage key. Defaults to "flexui-theme". */
  storageKey?: string;
}

const DEFAULT_STORAGE_KEY = "flexui-theme";
let currentMode: ThemeMode | null = null;

function isThemeMode(value: string | null | undefined): value is ThemeMode {
  return value === "light" || value === "dark" || value === "system";
}

function storage(): Storage | null {
  if (typeof window === "undefined") {
    return null;
  }
  try {
    return window.localStorage;
  } catch {
    return null;
  }
}

function mediaQuery(): MediaQueryList | null {
  if (typeof window === "undefined" || typeof window.matchMedia !== "function") {
    return null;
  }
  return window.matchMedia("(prefers-color-scheme: dark)");
}

function readStored(key: string): ThemeMode | null {
  try {
    const value = storage()?.getItem(key);
    return isThemeMode(value) ? value : null;
  } catch {
    return null;
  }
}

function writeStored(key: string, mode: ThemeMode): void {
  try {
    storage()?.setItem(key, mode);
  } catch {
    // Private mode and disabled storage are non-fatal.
  }
}

export function getResolvedTheme(mode: ThemeMode = getTheme()): ResolvedTheme {
  if (mode === "system") {
    return mediaQuery()?.matches ? "dark" : "light";
  }
  return mode;
}

export function getTheme(): ThemeMode {
  if (currentMode) {
    return currentMode;
  }
  currentMode = readStored(DEFAULT_STORAGE_KEY) ?? "system";
  return currentMode;
}

/** Applies a mode to <html> and optionally emits fui:themechange. */
export function applyTheme(mode: ThemeMode, silent = false): ResolvedTheme {
  currentMode = mode;
  const resolved = getResolvedTheme(mode);

  if (typeof document !== "undefined") {
    const root = document.documentElement;
    root.classList.toggle("dark", resolved === "dark");
    root.dataset["fuiThemeMode"] = mode;
    root.style.colorScheme = mode === "system" ? "light dark" : mode;
    if (!silent) {
      emit(document, "fui:themechange", {
        mode,
        resolved,
      } satisfies ThemeChangeDetail);
    }
  }

  return resolved;
}

export function nextTheme(mode: ThemeMode = getTheme()): ThemeMode {
  if (mode === "light") {
    return "dark";
  }
  if (mode === "dark") {
    return "system";
  }
  return "light";
}

/** Sets and (by default) persists a theme mode. */
export function setTheme(
  mode: ThemeMode,
  options: { persist?: boolean; storageKey?: string; silent?: boolean } = {},
): ThemeMode {
  applyTheme(mode, options.silent);
  if (options.persist !== false) {
    writeStored(options.storageKey ?? DEFAULT_STORAGE_KEY, mode);
  }
  syncThemeControls();
  return mode;
}

export function toggleTheme(): ThemeMode {
  return setTheme(nextTheme());
}

function syncThemeControls(root: ParentNode = document): void {
  if (typeof document === "undefined") {
    return;
  }
  const mode = getTheme();
  for (const button of qsa<HTMLElement>(root, "[data-fui-theme-set]")) {
    const value = button.getAttribute("data-fui-theme-set");
    button.setAttribute("aria-pressed", value === mode ? "true" : "false");
  }
}

/**
 * Applies the saved/system theme and binds theme controls.
 *
 * Controls:
 *   [data-fui-theme-set="light|dark|system"]
 *   [data-fui-theme-toggle]
 */
export function initTheme(options: ThemeOptions = {}): Cleanup {
  if (typeof document === "undefined") {
    return () => undefined;
  }

  const disposer = createDisposer();
  const root = options.root ?? document;
  const key = options.storageKey ?? DEFAULT_STORAGE_KEY;
  const saved = readStored(key);

  applyTheme(saved ?? "system", true);
  syncThemeControls(root);

  const query = mediaQuery();
  if (query) {
    const handleChange = (): void => {
      if (getTheme() === "system") {
        applyTheme("system");
      }
    };
    query.addEventListener("change", handleChange);
    disposer.add(() => query.removeEventListener("change", handleChange));
  }

  disposer.add(
    on(root, "click", (event) => {
      const target = event.target;
      if (!(target instanceof Element)) {
        return;
      }
      const control = target.closest<HTMLElement>("[data-fui-theme-set], [data-fui-theme-toggle]");
      if (!control || !root.contains(control)) {
        return;
      }

      const set = control.getAttribute("data-fui-theme-set");
      if (isThemeMode(set)) {
        setTheme(set, { storageKey: key });
      } else {
        setTheme(nextTheme(), { storageKey: key });
      }
      (event as MouseEvent).preventDefault();
    }),
  );

  disposer.add(on(document, "fui:themechange", () => syncThemeControls(root)));
  return () => disposer.dispose();
}
