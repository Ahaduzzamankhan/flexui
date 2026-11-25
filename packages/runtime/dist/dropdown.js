import { closestFromEvent, createDisposer, emit, getElement, on, qsa, uid, } from "./types.js";
const controllers = new WeakMap();
const ITEM_SELECTOR = '[role="menuitem"], [data-fui-dropdown-item]';
/**
 * A menu button controller: click/Enter/Space/Arrow keys, roving focus,
 * Escape, outside dismissal and aria-expanded wiring.
 */
export class Dropdown {
    trigger;
    menu;
    options;
    disposer;
    openState = false;
    constructor(trigger, menu, options = {}) {
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
        this.disposer.add(on(trigger, "click", (event) => {
            event.preventDefault();
            this.toggle();
        }));
        this.disposer.add(on(trigger, "keydown", (event) => {
            this.handleTriggerKey(event);
        }));
        this.disposer.add(on(menu, "keydown", (event) => {
            this.handleMenuKey(event);
        }));
        this.disposer.add(on(menu, "click", (event) => {
            this.handleMenuClick(event);
        }));
        this.disposer.add(on(document, "pointerdown", (event) => {
            const target = event.target;
            if (this.openState &&
                target instanceof Node &&
                !this.menu.contains(target) &&
                !this.trigger.contains(target)) {
                this.close(false);
            }
        }, true));
    }
    get isOpen() {
        return this.openState;
    }
    open(focus = "first") {
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
        }
        else if (focus === "last") {
            items[items.length - 1]?.focus();
        }
        emit(this.trigger, "fui:open", { menu: this.menu, trigger: this.trigger });
    }
    close(restoreFocus = true) {
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
    toggle() {
        if (this.openState) {
            this.close();
        }
        else {
            this.open("first");
        }
    }
    dispose() {
        this.close(false);
        this.disposer.dispose();
        controllers.delete(this.trigger);
    }
    items() {
        return qsa(this.menu, ITEM_SELECTOR).filter((item) => item.getAttribute("aria-disabled") !== "true" && !item.hasAttribute("disabled"));
    }
    handleTriggerKey(event) {
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
    handleMenuKey(event) {
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
    handleMenuClick(event) {
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
/** Binds dropdowns declared with data-fui-dropdown="<target>". */
export function initDropdowns(root = document) {
    if (typeof document === "undefined") {
        return () => undefined;
    }
    const disposer = createDisposer();
    for (const trigger of qsa(root, "[data-fui-dropdown]")) {
        if (controllers.has(trigger)) {
            continue;
        }
        const value = trigger.getAttribute("data-fui-dropdown") ?? trigger.getAttribute("aria-controls");
        const menu = getElement(value);
        if (!menu) {
            continue;
        }
        const controller = new Dropdown(trigger, menu);
        controllers.set(trigger, controller);
        disposer.add(() => controller.dispose());
    }
    return () => disposer.dispose();
}
