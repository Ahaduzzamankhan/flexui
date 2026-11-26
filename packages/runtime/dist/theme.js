import { createDisposer, emit, on, qsa } from "./types.js";
const DEFAULT_STORAGE_KEY = "flexui-theme";
let currentMode = null;
// Set by initTheme so getTheme()/setTheme()/toggleTheme() read and write the
// same key as the controls, even when a custom storageKey was configured.
let activeStorageKey = DEFAULT_STORAGE_KEY;
function isThemeMode(value) {
    return value === "light" || value === "dark" || value === "system";
}
function storage() {
    if (typeof window === "undefined") {
        return null;
    }
    try {
        return window.localStorage;
    }
    catch {
        return null;
    }
}
function mediaQuery() {
    if (typeof window === "undefined" || typeof window.matchMedia !== "function") {
        return null;
    }
    return window.matchMedia("(prefers-color-scheme: dark)");
}
function readStored(key) {
    try {
        const value = storage()?.getItem(key);
        return isThemeMode(value) ? value : null;
    }
    catch {
        return null;
    }
}
function writeStored(key, mode) {
    try {
        storage()?.setItem(key, mode);
    }
    catch {
        // Private mode and disabled storage are non-fatal.
    }
}
export function getResolvedTheme(mode = getTheme()) {
    if (mode === "system") {
        return mediaQuery()?.matches ? "dark" : "light";
    }
    return mode;
}
export function getTheme() {
    if (currentMode) {
        return currentMode;
    }
    currentMode = readStored(activeStorageKey) ?? "system";
    return currentMode;
}
/** Applies a mode to <html> and optionally emits fui:themechange. */
export function applyTheme(mode, silent = false) {
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
            });
        }
    }
    return resolved;
}
export function nextTheme(mode = getTheme()) {
    if (mode === "light") {
        return "dark";
    }
    if (mode === "dark") {
        return "system";
    }
    return "light";
}
/** Sets and (by default) persists a theme mode. */
export function setTheme(mode, options = {}) {
    applyTheme(mode, options.silent);
    if (options.storageKey !== undefined) {
        activeStorageKey = options.storageKey;
    }
    if (options.persist !== false) {
        writeStored(activeStorageKey, mode);
    }
    syncThemeControls();
    return mode;
}
export function toggleTheme() {
    return setTheme(nextTheme());
}
function syncThemeControls(root = document) {
    if (typeof document === "undefined") {
        return;
    }
    const mode = getTheme();
    for (const button of qsa(root, "[data-fui-theme-set]")) {
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
export function initTheme(options = {}) {
    if (typeof document === "undefined") {
        return () => undefined;
    }
    const disposer = createDisposer();
    const root = options.root ?? document;
    const key = options.storageKey ?? DEFAULT_STORAGE_KEY;
    activeStorageKey = key;
    const saved = readStored(key);
    applyTheme(saved ?? "system", true);
    syncThemeControls(root);
    const query = mediaQuery();
    if (query) {
        const handleChange = () => {
            if (getTheme() === "system") {
                applyTheme("system");
            }
        };
        query.addEventListener("change", handleChange);
        disposer.add(() => query.removeEventListener("change", handleChange));
    }
    disposer.add(on(root, "click", (event) => {
        const target = event.target;
        if (!(target instanceof Element)) {
            return;
        }
        const control = target.closest("[data-fui-theme-set], [data-fui-theme-toggle]");
        if (!control || !root.contains(control)) {
            return;
        }
        const set = control.getAttribute("data-fui-theme-set");
        if (isThemeMode(set)) {
            setTheme(set, { storageKey: key });
        }
        else {
            setTheme(nextTheme(), { storageKey: key });
        }
        event.preventDefault();
    }));
    disposer.add(on(document, "fui:themechange", () => syncThemeControls(root)));
    return () => disposer.dispose();
}
