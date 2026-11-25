import {
  createDisposer,
  emit,
  getElement,
  on,
  qsa,
  uid,
  type Cleanup,
  type Disposer,
} from "./types.js";

interface NativePopover extends HTMLElement {
  showPopover: () => void;
  hidePopover: () => void;
  togglePopover: (force?: boolean) => boolean;
}

export interface PopoverOptions {
  /** Pixels between trigger and panel. Defaults to 8. */
  offset?: number;
}

const controllers = new WeakMap<HTMLElement, PopoverController>();

/**
 * Positions a native popover against its trigger and keeps aria-expanded in
 * sync. Browsers without the popover API get a hidden-attribute fallback with
 * the same keyboard and outside-click behaviour.
 */
export class PopoverController {
  readonly trigger: HTMLElement;
  readonly panel: NativePopover;

  private readonly options: PopoverOptions;
  private readonly disposer: Disposer;
  private readonly native: boolean;
  private openState = false;

  constructor(trigger: HTMLElement, panel: HTMLElement, options: PopoverOptions = {}) {
    this.trigger = trigger;
    this.panel = panel as NativePopover;
    this.options = options;
    this.disposer = createDisposer();
    this.native = typeof this.panel.showPopover === "function";

    if (!this.panel.id) {
      this.panel.id = uid("fui-popover");
    }
    if (!this.panel.hasAttribute("popover")) {
      this.panel.setAttribute("popover", "auto");
    }
    this.trigger.setAttribute("aria-controls", this.panel.id);
    this.trigger.setAttribute("aria-expanded", "false");

    if (!this.native) {
      this.panel.hidden = true;
    }

    this.disposer.add(
      on(trigger, "click", (event) => {
        (event as MouseEvent).preventDefault();
        this.toggle();
      }),
    );
    this.disposer.add(
      on(this.panel, "toggle", (event) => {
        const state = (event as Event & { newState?: string }).newState;
        this.syncOpen(state === "open");
      }),
    );
    this.disposer.add(
      on(document, "pointerdown", (event) => {
        const target = event.target;
        if (
          this.openState &&
          target instanceof Node &&
          !this.panel.contains(target) &&
          !this.trigger.contains(target)
        ) {
          this.close();
        }
      }),
    );
    this.disposer.add(
      on(document, "keydown", (event) => {
        if (this.openState && (event as KeyboardEvent).key === "Escape") {
          this.close();
        }
      }),
    );
    this.disposer.add(on(window, "resize", () => this.position()));
    this.disposer.add(on(window, "scroll", () => this.position(), true));

    for (const closer of qsa<HTMLElement>(this.panel, "[data-fui-popover-close]")) {
      this.disposer.add(
        on(closer, "click", (event) => {
          (event as MouseEvent).preventDefault();
          this.close();
        }),
      );
    }
  }

  get isOpen(): boolean {
    return this.openState;
  }

  open(): void {
    if (this.native) {
      try {
        this.panel.showPopover?.();
      } catch {
        // InvalidStateError means it is already open.
      }
    } else {
      this.panel.hidden = false;
    }
    this.syncOpen(true);
  }

  close(): void {
    if (this.native) {
      try {
        this.panel.hidePopover?.();
      } catch {
        // Nothing to hide.
      }
    } else {
      this.panel.hidden = true;
    }
    this.syncOpen(false);
  }

  toggle(): void {
    if (this.openState) {
      this.close();
    } else {
      this.open();
    }
  }

  dispose(): void {
    this.close();
    this.disposer.dispose();
    controllers.delete(this.trigger);
  }

  private syncOpen(open: boolean): void {
    if (this.openState === open) {
      if (open) {
        this.position();
      }
      return;
    }

    this.openState = open;
    this.trigger.setAttribute("aria-expanded", open ? "true" : "false");
    if (open) {
      this.panel.dataset["fuiState"] = "open";
      this.position();
      emit(this.trigger, "fui:popoveropen", { trigger: this.trigger, panel: this.panel });
    } else {
      delete this.panel.dataset["fuiState"];
      delete this.panel.dataset["fuiSide"];
      emit(this.trigger, "fui:popoverclose", { trigger: this.trigger, panel: this.panel });
    }
  }

  private position(): void {
    if (!this.openState) {
      return;
    }

    const offset = this.options.offset ?? 8;
    const margin = 8;
    const triggerRect = this.trigger.getBoundingClientRect();
    const panelRect = this.panel.getBoundingClientRect();

    let side: "top" | "bottom" = "bottom";
    let top = triggerRect.bottom + offset;
    if (top + panelRect.height > window.innerHeight - margin && triggerRect.top - panelRect.height - offset >= margin) {
      top = triggerRect.top - panelRect.height - offset;
      side = "top";
    }
    const left = Math.min(
      Math.max(margin, triggerRect.left),
      Math.max(margin, window.innerWidth - panelRect.width - margin),
    );

    this.panel.style.position = "fixed";
    this.panel.style.left = `${Math.round(left)}px`;
    this.panel.style.top = `${Math.round(top)}px`;
    this.panel.style.right = "auto";
    this.panel.style.bottom = "auto";
    this.panel.style.margin = "0";
    this.panel.dataset["fuiSide"] = side;
  }
}

/** Binds popovers declared with data-fui-popover="<target>". */
export function initPopovers(root: ParentNode = document): Cleanup {
  if (typeof document === "undefined") {
    return () => undefined;
  }

  const disposer = createDisposer();

  for (const trigger of qsa<HTMLElement>(root, "[data-fui-popover]")) {
    if (controllers.has(trigger)) {
      continue;
    }
    const value =
      trigger.getAttribute("data-fui-popover") ??
      trigger.getAttribute("aria-controls") ??
      trigger.getAttribute("popovertarget");
    const panel = getElement(value);
    if (!panel) {
      continue;
    }
    const controller = new PopoverController(trigger, panel);
    controllers.set(trigger, controller);
    disposer.add(() => controller.dispose());
  }

  return () => disposer.dispose();
}
