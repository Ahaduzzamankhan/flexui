import { type Cleanup } from "./types.js";
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
export declare function getResolvedTheme(mode?: ThemeMode): ResolvedTheme;
export declare function getTheme(): ThemeMode;
/** Applies a mode to <html> and optionally emits fui:themechange. */
export declare function applyTheme(mode: ThemeMode, silent?: boolean): ResolvedTheme;
export declare function nextTheme(mode?: ThemeMode): ThemeMode;
/** Sets and (by default) persists a theme mode. */
export declare function setTheme(mode: ThemeMode, options?: {
    persist?: boolean;
    storageKey?: string;
    silent?: boolean;
}): ThemeMode;
export declare function toggleTheme(): ThemeMode;
/**
 * Applies the saved/system theme and binds theme controls.
 *
 * Controls:
 *   [data-fui-theme-set="light|dark|system"]
 *   [data-fui-theme-toggle]
 */
export declare function initTheme(options?: ThemeOptions): Cleanup;
