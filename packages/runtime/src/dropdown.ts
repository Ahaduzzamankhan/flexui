import {
  closestFromEvent,
  createDisposer,
  emit,
  getElement,
  on,
  qsa,
  uid,
  type Cleanup,
  type Disposer,
} from "./types.js";

export interface DropdownOptions {
  /** Close after a menu item is selected. Defaults to true. */
  closeOnSelect?: boolean;
}

const controllers = new WeakMap<HTMLElement, Dropdown>();
const ITEM_SELECTOR = '[role="menuitem"], [data-fui-dropdown-item]';

/**
 * A menu button controller: click/Enter/Space/Arrow keys, roving focus,
 * Escape, outside dismissal and aria-expanded wiring.
 */
export class Dropdown {
  readonly trigger: HTMLElement;
  readonly menu: HTMLElement;

  private readonly options: DropdownOptions;
  private readonly disposer: Disposer;
  private openState = false;

  constructor(trigger: HTMLElement, menu: HTMLElement, options: DropdownOptions = {}) {
    this.trigger = trigger;
    this.menu = menu;
    this.options = options;
    this.disposer = createDisposer();

    if (!menu.id) {
      menu.id = uid("fui-menu");
    }
    trigger.setAttribute("aria-haspopup", "menu");
    trigger.setAttribute("aria-controls", menu.id);
    trigger.setAttribute("aria-expanded", "false");
    if (!menu.hasAttribute("role")) {
      menu.setAttribute("role", "menu");
    }
    menu.hidden = true;

    this.disposer.add(
      on(trigger, "click", (event) => {
        (event as MouseEvent).preventDefault();
        this.toggle();
      }),
    );
    this.disposer.add(
      on(trigger, "keydown", (event) => {
        this.handleTriggerKey(event as KeyboardEvent);
      }),
    );
    this.disposer.add(
      on(menu, "keydown", (event) => {
        this.handleMenuKey(event as KeyboardEvent);
      }),
    );
    this.disposer.add(
      on(menu, "click", (event) => {
        this.handleMenuClick(event as MouseEvent);
      }),
    );
    this.disposer.add(
      on(
        document,
        "pointerdown",
        (event) => {
          const target = event.target;
          if (
            this.openState &&
            target instanceof Node &&
            !this.menu.contains(target) &&
            !this.trigger.contains(target)
          ) {
            // Focus stays where the user clicked; only Escape and selection
            // hand focus back to the trigger.
            this.close(false);
          }
        },
        true,
      ),
    );
  }

  get isOpen(): boolean {
    return this.openState;
  }

  open(focus: "first" | "last" | "none" = "first"): void {
    if (this.openState) {
      return;
    }
    this.openState = true;
    this.menu.hidden = false;
    this.trigger.setAttribute("aria-expanded", "true");
    this.menu.dataset["fuiState"] = "open";
    this.trigger.dataset["fuiState"] = "open";

    const items = this.items();
    if (focus === "first") {
      items[0]?.focus();
    } else if (focus === "last") {
      items[items.length - 1]?.focus();
    }
    emit(this.trigger, "fui:open", { menu: this.menu, trigger: this.trigger });
  }

  close(restoreFocus = true): void {
    if (!this.openState) {
      return;
    }
    this.openState = false;
    this.menu.hidden = true;
    this.trigger.setAttribute("aria-expanded", "false");
    delete this.menu.dataset["fuiState"];
    delete this.trigger.dataset["fuiState"];
    if (restoreFocus) {
      this.trigger.focus({ preventScroll: true });
    }
    emit(this.trigger, "fui:close", { menu: this.menu, trigger: this.trigger });
  }

  toggle(): void {
    if (this.openState) {
      this.close();
    } else {
      this.open("first");
    }
  }

  dispose(): void {
    this.close(false);
    this.disposer.dispose();
    controllers.delete(this.trigger);
  }

  private items(): HTMLElement[] {
    return qsa<HTMLElement>(this.menu, ITEM_SELECTOR).filter(
      (item) => item.getAttribute("aria-disabled") !== "true" && !item.hasAttribute("disabled"),
    );
  }

  private handleTriggerKey(event: KeyboardEvent): void {
    switch (event.key) {
      case "ArrowDown":
      case "Enter":
      case " ":
        event.preventDefault();
        this.open("first");
        break;
      case "ArrowUp":
        event.preventDefault();
        this.open("last");
        break;
      case "Escape":
        this.close();
        break;
      default:
        break;
    }
  }

  private handleMenuKey(event: KeyboardEvent): void {
    const items = this.items();
    if (items.length === 0) {
      return;
    }

    const current = document.activeElement;
    let index = items.findIndex((item) => item === current);

    switch (event.key) {
      case "ArrowDown":
        event.preventDefault();
        index = index < 0 ? 0 : (index + 1) % items.length;
        items[index]?.focus();
        break;
      case "ArrowUp":
        event.preventDefault();
        index = index <= 0 ? items.length - 1 : index - 1;
        items[index]?.focus();
        break;
      case "Home":
        event.preventDefault();
        items[0]?.focus();
        break;
      case "End":
        event.preventDefault();
        items[items.length - 1]?.focus();
        break;
      case "Escape":
        event.preventDefault();
        this.close();
        break;
      case "Tab":
        this.close(false);
        break;
      default:
        break;
    }
  }

  private handleMenuClick(event: MouseEvent): void {
    const item = closestFromEvent(event, ITEM_SELECTOR);
    if (!item || !this.menu.contains(item)) {
      return;
    }
    if (item.getAttribute("aria-disabled") === "true" || item.hasAttribute("disabled")) {
      event.preventDefault();
      return;
    }

    emit(item, "fui:select", { item, menu: this.menu, trigger: this.trigger });
    if (this.options.closeOnSelect !== false) {
      this.close();
    }
  }
}

/** Returns the controller bound to a trigger, if any. */
export function getDropdown(trigger: HTMLElement): Dropdown | undefined {
  return controllers.get(trigger);
}

/** Creates or reuses the controller for a trigger and its menu. */
export function bindDropdown(
  trigger: HTMLElement,
  menu: HTMLElement,
  options: DropdownOptions = {},
): Dropdown {
  const existing = controllers.get(trigger);
  if (existing) {
    return existing;
  }
  const controller = new Dropdown(trigger, menu, options);
  controllers.set(trigger, controller);
  return controller;
}

/** Binds dropdowns declared with data-fui-dropdown="<target>". */
export function initDropdowns(root: ParentNode = document): Cleanup {
  if (typeof document === "undefined") {
    return () => undefined;
  }

  const disposer = createDisposer();

  for (const trigger of qsa<HTMLElement>(root, "[data-fui-dropdown]")) {
    if (controllers.has(trigger)) {
      continue;
    }
    const value = trigger.getAttribute("data-fui-dropdown") ?? trigger.getAttribute("aria-controls");
    const menu = getElement(value);
    if (!menu) {
      continue;
    }
    const controller = bindDropdown(trigger, menu);
    disposer.add(() => controller.dispose());
  }

  return () => disposer.dispose();
}
