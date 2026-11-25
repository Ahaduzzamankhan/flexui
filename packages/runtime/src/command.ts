import { bindDialog, getDialog, type Dialog } from "./dialog.js";
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

export interface CommandPaletteOptions {
  /** Close the palette after a command is chosen. Defaults to true. */
  closeOnSelect?: boolean;
}

interface CommandEntry {
  element: HTMLElement;
  label: string;
  value: string;
  search: string;
}

const controllers = new WeakMap<HTMLElement, CommandPalette>();
const ITEM_SELECTOR = "[data-fui-command-item], [role='option']";

/**
 * Keyboard-first command palette over existing listbox markup.
 *
 * Item attributes: data-fui-command-item, data-value and optional
 * data-keywords. Selection dispatches a bubbling fui:command event.
 */
export class CommandPalette {
  readonly element: HTMLElement;

  private readonly input: HTMLInputElement;
  private readonly list: HTMLElement;
  private readonly empty: HTMLElement | null;
  private readonly dialog: Dialog | null;
  private readonly options: CommandPaletteOptions;
  private readonly disposer: Disposer;
  private entries: CommandEntry[] = [];
  private filtered: CommandEntry[] = [];
  private activeIndex = -1;

  constructor(element: HTMLElement, options: CommandPaletteOptions = {}) {
    this.element = element;
    this.options = options;
    this.disposer = createDisposer();

    const input = element.querySelector<HTMLInputElement>(
      "[data-fui-command-input], input[role='combobox'], input",
    );
    const list = element.querySelector<HTMLElement>("[data-fui-command-list], [role='listbox']");
    if (!input || !list) {
      throw new Error("@flexui/runtime: command palette needs an input and a list");
    }

    this.input = input;
    this.list = list;
    this.empty = element.querySelector<HTMLElement>("[data-fui-command-empty]");
    this.dialog = element instanceof HTMLDialogElement ? (getDialog(element) ?? bindDialog(element)) : null;

    if (!input.id) {
      input.id = uid("fui-command-input");
    }
    if (!list.id) {
      list.id = uid("fui-command-list");
    }
    input.setAttribute("role", "combobox");
    input.setAttribute("aria-controls", list.id);
    input.setAttribute("aria-expanded", "true");
    input.setAttribute("aria-autocomplete", "list");
    input.setAttribute("autocomplete", "off");
    input.spellcheck = false;
    list.setAttribute("role", "listbox");

    if (!this.dialog) {
      element.hidden = true;
    }

    this.disposer.add(
      on(input, "input", () => {
        this.filter(input.value);
      }),
    );
    this.disposer.add(
      on(input, "keydown", (event) => {
        this.handleKeydown(event as KeyboardEvent);
      }),
    );
    this.disposer.add(
      on(list, "click", (event) => {
        const mouse = event as MouseEvent;
        const item = closestFromEvent(mouse, ITEM_SELECTOR);
        if (!item || !list.contains(item)) {
          return;
        }
        mouse.preventDefault();
        const index = this.filtered.findIndex((entry) => entry.element === item);
        this.select(index);
      }),
    );
    this.disposer.add(
      on(list, "pointermove", (event) => {
        const item = closestFromEvent(event, ITEM_SELECTOR);
        if (!item) {
          return;
        }
        const index = this.filtered.findIndex((entry) => entry.element === item);
        if (index >= 0 && index !== this.activeIndex) {
          this.setActive(index, false);
        }
      }),
    );
    this.disposer.add(
      on(element, "close", () => {
        this.reset();
      }),
    );

    this.refresh();
  }

  get isOpen(): boolean {
    return this.dialog ? this.dialog.isOpen : !this.element.hidden;
  }

  /** Re-reads items from the DOM. Call after adding or removing commands. */
  refresh(): void {
    this.entries = qsa<HTMLElement>(this.list, ITEM_SELECTOR).map((element) => {
      if (!element.id) {
        element.id = uid("fui-command-item");
      }
      element.setAttribute("role", "option");
      element.tabIndex = -1;
      const label = (element.getAttribute("data-label") ?? element.textContent ?? "").trim();
      const keywords = element.getAttribute("data-keywords") ?? "";
      const value = element.getAttribute("data-value") ?? label;
      return {
        element,
        label,
        value,
        search: `${label} ${keywords} ${value}`.toLowerCase(),
      };
    });
    this.filter(this.input.value);
  }

  open(): void {
    if (this.dialog) {
      this.dialog.open();
    } else {
      this.element.hidden = false;
      this.element.dataset["fuiState"] = "open";
    }
    this.filter(this.input.value);
    this.input.focus({ preventScroll: true });
    this.input.select();
    emit(this.element, "fui:commandopen", { palette: this });
  }

  close(): void {
    if (this.dialog) {
      this.dialog.close("command");
    } else {
      this.element.hidden = true;
      delete this.element.dataset["fuiState"];
    }
    emit(this.element, "fui:commandclose", { palette: this });
  }

  toggle(): void {
    if (this.isOpen) {
      this.close();
    } else {
      this.open();
    }
  }

  dispose(): void {
    if (this.dialog?.isOpen) {
      this.dialog.close("dispose");
    }
    this.disposer.dispose();
    controllers.delete(this.element);
  }

  private reset(): void {
    this.input.value = "";
    this.filter("");
  }

  private filter(query: string): void {
    const terms = query.toLowerCase().trim().split(/\s+/).filter(Boolean);
    this.filtered = this.entries.filter((entry) => terms.every((term) => entry.search.includes(term)));

    const visible = new Set(this.filtered.map((entry) => entry.element));
    for (const entry of this.entries) {
      entry.element.hidden = !visible.has(entry.element);
      entry.element.setAttribute("aria-selected", "false");
    }
    if (this.empty) {
      this.empty.hidden = this.filtered.length > 0;
    }

    this.setActive(this.filtered.length > 0 ? 0 : -1, false);
  }

  private setActive(index: number, scroll: boolean): void {
    const previous = this.filtered[this.activeIndex];
    previous?.element.setAttribute("aria-selected", "false");

    this.activeIndex = index;
    const entry = this.filtered[index];
    if (!entry) {
      this.input.removeAttribute("aria-activedescendant");
      return;
    }

    entry.element.setAttribute("aria-selected", "true");
    this.input.setAttribute("aria-activedescendant", entry.element.id);
    if (scroll) {
      entry.element.scrollIntoView({ block: "nearest" });
    }
  }

  private handleKeydown(event: KeyboardEvent): void {
    switch (event.key) {
      case "ArrowDown":
        event.preventDefault();
        if (this.filtered.length > 0) {
          this.setActive((this.activeIndex + 1) % this.filtered.length, true);
        }
        break;
      case "ArrowUp":
        event.preventDefault();
        if (this.filtered.length > 0) {
          this.setActive(this.activeIndex <= 0 ? this.filtered.length - 1 : this.activeIndex - 1, true);
        }
        break;
      case "Home":
        event.preventDefault();
        this.setActive(0, true);
        break;
      case "End":
        event.preventDefault();
        this.setActive(this.filtered.length - 1, true);
        break;
      case "Enter":
        event.preventDefault();
        this.select(this.activeIndex);
        break;
      case "Escape":
        event.preventDefault();
        this.close();
        break;
      default:
        break;
    }
  }

  private select(index: number): void {
    const entry = this.filtered[index];
    if (!entry) {
      return;
    }

    if (this.options.closeOnSelect !== false) {
      this.close();
    }
    emit(entry.element, "fui:command", {
      value: entry.value,
      label: entry.label,
      item: entry.element,
      palette: this,
    });
  }
}

/** Binds [data-fui-command] palettes, Ctrl/Cmd+K and command openers. */
export function initCommandPalettes(root: ParentNode = document): Cleanup {
  if (typeof document === "undefined") {
    return () => undefined;
  }

  const disposer = createDisposer();
  const palettes: CommandPalette[] = [];

  for (const element of qsa<HTMLElement>(root, "[data-fui-command]")) {
    if (controllers.has(element)) {
      continue;
    }
    try {
      const controller = new CommandPalette(element);
      controllers.set(element, controller);
      disposer.add(() => controller.dispose());
      palettes.push(controller);
    } catch {
      // Misconfigured palettes are ignored rather than breaking the page.
    }
  }

  if (palettes.length === 0) {
    return () => disposer.dispose();
  }

  disposer.add(
    on(document, "keydown", (event) => {
      const keyboard = event as KeyboardEvent;
      if (
        (keyboard.metaKey || keyboard.ctrlKey) &&
        !keyboard.altKey &&
        keyboard.key.toLowerCase() === "k"
      ) {
        keyboard.preventDefault();
        palettes[0]?.toggle();
      }
    }),
  );

  disposer.add(
    on(document, "click", (event) => {
      const opener = closestFromEvent(event, "[data-fui-command-open]");
      if (!opener) {
        return;
      }
      (event as MouseEvent).preventDefault();
      const value = opener.getAttribute("data-fui-command-open");
      const target = getElement(value);
      const palette = target ? controllers.get(target) : palettes[0];
      palette?.open();
    }),
  );

  return () => disposer.dispose();
}
