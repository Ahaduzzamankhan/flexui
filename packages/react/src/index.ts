/**
 * React bindings for `flexui-runtime`.
 *
 * The hooks own the controller lifecycle: they bind on mount, dispose on
 * unmount, and mirror the controller's open state into React state so a
 * component re-renders when a dialog closes from the keyboard or a click
 * outside it. Component CSS still works without any of this.
 */

import { useCallback, useEffect, useRef, useState, type RefObject } from "react";

import { bindDialog, type Dialog, type DialogOptions } from "flexui-runtime/dialog";
import { bindDropdown, type Dropdown, type DropdownOptions } from "flexui-runtime/dropdown";
import { bindPopover, type PopoverController, type PopoverOptions } from "flexui-runtime/popover";
import { bindTabs, type Tabs } from "flexui-runtime/tabs";
import {
  getResolvedTheme,
  getTheme,
  setTheme,
  toggleTheme,
  type ResolvedTheme,
  type ThemeChangeDetail,
  type ThemeMode,
} from "flexui-runtime/theme";

/**
 * Mirrors a controller's open state into React state. The ref is read inside
 * the effect, after React has attached it, which is why the element itself is
 * not passed in: on the first render it is still null.
 */
function useOpenState(
  ref: RefObject<Element | null>,
  openEvent: "fui:open" | "fui:popoveropen",
): boolean {
  const [isOpen, setIsOpen] = useState(false);
  const closeEvent = openEvent === "fui:open" ? "fui:close" : "fui:popoverclose";

  useEffect(() => {
    const element = ref.current;
    if (element === null) {
      return;
    }
    const open = (): void => setIsOpen(true);
    const close = (): void => setIsOpen(false);
    element.addEventListener(openEvent, open);
    element.addEventListener(closeEvent, close);
    return () => {
      element.removeEventListener(openEvent, open);
      element.removeEventListener(closeEvent, close);
    };
  }, [ref, openEvent, closeEvent]);

  return isOpen;
}

export interface UseDialogResult {
  /** Attach to the `<dialog>` element. */
  ref: RefObject<HTMLDialogElement | null>;
  open(): void;
  close(reason?: string): void;
  /** Mirrors the dialog's state, including closes from Escape or the backdrop. */
  isOpen: boolean;
}

/**
 * Binds a native `<dialog class="fui-dialog" data-fui-dialog>`.
 * Options are read when the dialog is first bound.
 */
export function useDialog(options: DialogOptions = {}): UseDialogResult {
  const ref = useRef<HTMLDialogElement | null>(null);
  const controller = useRef<Dialog | null>(null);
  const isOpen = useOpenState(ref, "fui:open");

  useEffect(() => {
    const element = ref.current;
    if (element === null) {
      return;
    }
    const instance = bindDialog(element, options);
    controller.current = instance;
    return () => {
      controller.current = null;
      instance.dispose();
    };
  }, []);

  const open = useCallback((): void => {
    controller.current?.open();
  }, []);
  const close = useCallback((reason?: string): void => {
    controller.current?.close(reason);
  }, []);

  return { ref, open, close, isOpen };
}

export interface UseDropdownResult {
  /** Attach to the trigger button. */
  triggerRef: RefObject<HTMLElement | null>;
  open(focus?: "first" | "last" | "none"): void;
  close(restoreFocus?: boolean): void;
  toggle(): void;
  isOpen: boolean;
}

/** Binds a trigger to the menu returned by `menu`, an id or an element. */
export function useDropdown(
  menu: string | HTMLElement | null,
  options: DropdownOptions = {},
): UseDropdownResult {
  const triggerRef = useRef<HTMLElement | null>(null);
  const controller = useRef<Dropdown | null>(null);
  const isOpen = useOpenState(triggerRef, "fui:open");

  useEffect(() => {
    const trigger = triggerRef.current;
    if (trigger === null) {
      return;
    }
    const target = typeof menu === "string" ? document.getElementById(menu.replace(/^#/, "")) : menu;
    if (target === null) {
      return;
    }
    const instance = bindDropdown(trigger, target, options);
    controller.current = instance;
    return () => {
      controller.current = null;
      instance.dispose();
    };
  }, [menu]);

  const open = useCallback((focus?: "first" | "last" | "none"): void => {
    controller.current?.open(focus);
  }, []);
  const close = useCallback((restoreFocus?: boolean): void => {
    controller.current?.close(restoreFocus);
  }, []);
  const toggle = useCallback((): void => {
    controller.current?.toggle();
  }, []);

  return { triggerRef, open, close, toggle, isOpen };
}

export interface UsePopoverResult {
  /** Attach to the trigger button. */
  triggerRef: RefObject<HTMLElement | null>;
  open(): void;
  close(): void;
  toggle(): void;
  isOpen: boolean;
}

/** Binds a trigger to a `popover` panel returned by `panel`, an id or an element. */
export function usePopover(
  panel: string | HTMLElement | null,
  options: PopoverOptions = {},
): UsePopoverResult {
  const triggerRef = useRef<HTMLElement | null>(null);
  const controller = useRef<PopoverController | null>(null);
  const isOpen = useOpenState(triggerRef, "fui:popoveropen");

  useEffect(() => {
    const trigger = triggerRef.current;
    if (trigger === null) {
      return;
    }
    const target = typeof panel === "string" ? document.getElementById(panel.replace(/^#/, "")) : panel;
    if (target === null) {
      return;
    }
    const instance = bindPopover(trigger, target, options);
    controller.current = instance;
    return () => {
      controller.current = null;
      instance.dispose();
    };
  }, [panel]);

  const open = useCallback((): void => {
    controller.current?.open();
  }, []);
  const close = useCallback((): void => {
    controller.current?.close();
  }, []);
  const toggle = useCallback((): void => {
    controller.current?.toggle();
  }, []);

  return { triggerRef, open, close, toggle, isOpen };
}

export interface UseTabsResult {
  /** Attach to the `[data-fui-tabs]` container. */
  ref: RefObject<HTMLElement | null>;
  activeIndex: number;
  activate(index: number, focus?: boolean): void;
}

/** Binds an ARIA tablist and tracks its active tab in React state. */
export function useTabs(): UseTabsResult {
  const ref = useRef<HTMLElement | null>(null);
  const controller = useRef<Tabs | null>(null);
  const [activeIndex, setActiveIndex] = useState(0);

  useEffect(() => {
    const element = ref.current;
    if (element === null) {
      return;
    }
    const instance = bindTabs(element);
    controller.current = instance;
    setActiveIndex(instance.active?.index ?? 0);

    const handle = (event: Event): void => {
      const detail = (event as CustomEvent<{ index: number }>).detail;
      setActiveIndex(detail.index);
    };
    element.addEventListener("fui:tabchange", handle);
    return () => {
      element.removeEventListener("fui:tabchange", handle);
      controller.current = null;
      instance.dispose();
    };
  }, []);

  const activate = useCallback((index: number, focus?: boolean): void => {
    controller.current?.activate(index, focus);
  }, []);

  return { ref, activeIndex, activate };
}

export interface UseThemeResult {
  mode: ThemeMode;
  resolved: ResolvedTheme;
  setMode(mode: ThemeMode): void;
  toggle(): void;
}

/** Reads and changes the FlexUI theme mode, re-rendering on `fui:themechange`. */
export function useTheme(): UseThemeResult {
  const [mode, setModeState] = useState<ThemeMode>(() => getTheme());
  const [resolved, setResolved] = useState<ResolvedTheme>(() => getResolvedTheme(getTheme()));

  useEffect(() => {
    const handle = (event: Event): void => {
      const detail = (event as CustomEvent<ThemeChangeDetail>).detail;
      setModeState(detail.mode);
      setResolved(detail.resolved);
    };
    document.addEventListener("fui:themechange", handle);
    return () => {
      document.removeEventListener("fui:themechange", handle);
    };
  }, []);

  const setMode = useCallback((next: ThemeMode): void => {
    setTheme(next);
  }, []);
  const toggle = useCallback((): void => {
    toggleTheme();
  }, []);

  return { mode, resolved, setMode, toggle };
}
