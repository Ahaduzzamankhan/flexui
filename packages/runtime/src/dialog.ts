import { focusFirst, trapFocus } from "./focus.js";
import {
  closestFromEvent,
  createDisposer,
  emit,
  getElement,
  on,
  qsa,
  type Cleanup,
  type Disposer,
} from "./types.js";

export interface DialogOptions {
  /** Close when the backdrop is clicked. Defaults to true. */
  closeOnBackdrop?: boolean;
  /** Return focus to the opener when the dialog closes. Defaults to true. */
  restoreFocus?: boolean;
  onOpen?: (dialog: HTMLDialogElement) => void;
  onClose?: (dialog: HTMLDialogElement, reason: string) => void;
}

const controllers = new WeakMap<HTMLDialogElement, Dialog>();
let openOverlays = 0;

function syncOverlayLock(): void {
  if (typeof document === "undefined") {
    return;
  }
  const root = document.documentElement;
  if (openOverlays > 0) {
    root.dataset["fuiOverlayOpen"] = "";
  } else {
    delete root.dataset["fuiOverlayOpen"];
  }
}

function eventReason(element: HTMLDialogElement, fallback: string): string {
  return element.returnValue !== "" ? element.returnValue : fallback;
}

/**
 * Accessible controller for a native <dialog>.
 *
 * showModal() supplies the top layer and inert background. This class adds
 * backdrop dismissal, scroll locking, focus restoration and DOM events.
 */
export class Dialog {
  readonly element: HTMLDialogElement;

  private readonly options: DialogOptions;
  private readonly disposer: Disposer;
  private lastFocus: HTMLElement | null = null;
  private openState = false;
  private pendingReason = "dismiss";
  private releaseTrap: Cleanup | null = null;

  constructor(element: HTMLDialogElement, options: DialogOptions = {}) {
    this.element = element;
    this.options = options;
    this.disposer = createDisposer();

    this.disposer.add(
      on(element, "click", (event) => {
        this.handleClick(event as MouseEvent);
      }),
    );
    this.disposer.add(
      on(element, "cancel", () => {
        this.pendingReason = "escape";
      }),
    );
    this.disposer.add(
      on(element, "close", () => {
        this.handleClosed();
      }),
    );
  }

  get isOpen(): boolean {
    return this.openState;
  }

  open(): void {
    if (this.openState) {
      return;
    }

    if (this.element.open) {
      this.adoptOpen();
      return;
    }

    this.lastFocus = document.activeElement instanceof HTMLElement ? document.activeElement : null;
    this.pendingReason = "programmatic";
    this.element.returnValue = "";

    try {
      if (typeof this.element.showModal === "function") {
        this.element.showModal();
      } else {
        this.element.setAttribute("open", "");
      }
    } catch {
      this.element.setAttribute("open", "");
    }

    this.markOpen();
  }

  close(reason = "dismiss"): void {
    if (!this.openState && !this.element.open) {
      return;
    }

    this.pendingReason = reason;
    try {
      if (typeof this.element.close === "function") {
        this.element.close(reason);
        this.handleClosed();
      } else {
        this.element.removeAttribute("open");
        this.handleClosed();
      }
    } catch {
      this.element.removeAttribute("open");
      this.handleClosed();
    }
  }

  dispose(): void {
    if (this.openState) {
      this.close("dispose");
    }
    this.disposer.dispose();
    controllers.delete(this.element);
  }

  private adoptOpen(): void {
    this.lastFocus = document.activeElement instanceof HTMLElement ? document.activeElement : null;
    this.pendingReason = "programmatic";
    this.markOpen();
  }

  private markOpen(): void {
    if (this.openState) {
      return;
    }
    this.openState = true;
    openOverlays += 1;
    this.element.dataset["fuiState"] = "open";
    syncOverlayLock();
    this.releaseTrap = trapFocus(this.element);
    focusFirst(this.element);
    emit(this.element, "fui:open", { dialog: this.element });
    this.options.onOpen?.(this.element);
  }

  private handleClosed(): void {
    if (!this.openState) {
      return;
    }

    this.openState = false;
    openOverlays = Math.max(0, openOverlays - 1);
    delete this.element.dataset["fuiState"];
    syncOverlayLock();
    this.releaseTrap?.();
    this.releaseTrap = null;

    const reason = eventReason(this.element, this.pendingReason);
    this.pendingReason = "dismiss";
    this.element.returnValue = "";

    const previous = this.lastFocus;
    this.lastFocus = null;
    if (this.options.restoreFocus !== false && previous?.isConnected) {
      previous.focus({ preventScroll: true });
    }

    emit(this.element, "fui:close", { dialog: this.element, reason });
    this.options.onClose?.(this.element, reason);
  }

  private handleClick(event: MouseEvent): void {
    if (!this.openState || this.options.closeOnBackdrop === false) {
      return;
    }
    if (event.target !== this.element) {
      return;
    }

    const rect = this.element.getBoundingClientRect();
    const inside =
      event.clientX >= rect.left &&
      event.clientX <= rect.right &&
      event.clientY >= rect.top &&
      event.clientY <= rect.bottom;
    if (!inside) {
      this.close("backdrop");
    }
  }
}

/** Returns an existing controller, if one was bound. */
export function getDialog(element: HTMLDialogElement): Dialog | undefined {
  return controllers.get(element);
}

/** Creates or reuses the controller for a dialog element. */
export function bindDialog(element: HTMLDialogElement, options: DialogOptions = {}): Dialog {
  const existing = controllers.get(element);
  if (existing) {
    return existing;
  }
  const controller = new Dialog(element, options);
  controllers.set(element, controller);
  return controller;
}

function resolveDialog(value: string | null, fallback: Element | null): HTMLDialogElement | null {
  if (value) {
    const target = getElement(value);
    if (target instanceof HTMLDialogElement) {
      return target;
    }
  }
  const dialog = fallback?.closest("dialog");
  return dialog instanceof HTMLDialogElement ? dialog : null;
}

/**
 * Delegated open/close listeners for data-fui-open and data-fui-close.
 * This is exported so custom overlays can share the same trigger protocol.
 */
export function bindOverlayTriggers(): Cleanup {
  if (typeof document === "undefined") {
    return () => undefined;
  }

  const disposer = createDisposer();

  disposer.add(
    on(document, "click", (event) => {
      const mouse = event as MouseEvent;
      const opener = closestFromEvent(mouse, "[data-fui-open], [data-fui-dialog-open]");
      if (opener) {
        if (opener.getAttribute("aria-disabled") === "true") {
          return;
        }
        const value = opener.getAttribute("data-fui-open") ?? opener.getAttribute("data-fui-dialog-open");
        const target = resolveDialog(value, opener);
        if (target) {
          mouse.preventDefault();
          bindDialog(target).open();
        }
        return;
      }

      const closer = closestFromEvent(mouse, "[data-fui-close], [data-fui-dialog-close]");
      if (closer) {
        const value = closer.getAttribute("data-fui-close") ?? closer.getAttribute("data-fui-dialog-close");
        const target = resolveDialog(value, closer);
        if (target) {
          mouse.preventDefault();
          bindDialog(target).close("trigger");
        }
      }
    }),
  );

  return () => disposer.dispose();
}

/** Binds every dialog/drawer under root and the shared trigger delegation. */
export function initDialogs(root: ParentNode = document): Cleanup {
  if (typeof document === "undefined") {
    return () => undefined;
  }

  const disposer = createDisposer();

  for (const element of qsa<HTMLDialogElement>(root, "dialog[data-fui-dialog], dialog[data-fui-drawer]")) {
    if (!controllers.has(element)) {
      disposer.add(() => bindDialog(element).dispose());
    }
  }

  for (const trigger of qsa<HTMLElement>(root, "[data-fui-open], [data-fui-dialog-open]")) {
    const value = trigger.getAttribute("data-fui-open") ?? trigger.getAttribute("data-fui-dialog-open");
    const target = resolveDialog(value, trigger);
    if (target) {
      trigger.setAttribute("aria-haspopup", "dialog");
      if (target.id) {
        trigger.setAttribute("aria-controls", target.id);
      }
    }
  }

  disposer.add(bindOverlayTriggers());
  return () => disposer.dispose();
}
