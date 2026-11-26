import { createDisposer, emit, getElement, on, qsa, uid, } from "./types.js";
const controllers = new WeakMap();
/**
 * Positions a native popover against its trigger and keeps aria-expanded in
 * sync. Browsers without the popover API get a hidden-attribute fallback with
 * the same keyboard and outside-click behaviour.
 */
export class PopoverController {
    trigger;
    panel;
    options;
    disposer;
    native;
    openState = false;
    constructor(trigger, panel, options = {}) {
        this.trigger = trigger;
        this.panel = panel;
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
        this.disposer.add(on(trigger, "click", (event) => {
            event.preventDefault();
            this.toggle();
        }));
        this.disposer.add(on(this.panel, "toggle", (event) => {
            const state = event.newState;
            this.syncOpen(state === "open");
        }));
        this.disposer.add(on(document, "pointerdown", (event) => {
            const target = event.target;
            if (this.openState &&
                target instanceof Node &&
                !this.panel.contains(target) &&
                !this.trigger.contains(target)) {
                this.close();
            }
        }));
        this.disposer.add(on(document, "keydown", (event) => {
            if (this.openState && event.key === "Escape") {
                this.close();
            }
        }));
        this.disposer.add(on(window, "resize", () => this.position()));
        this.disposer.add(on(window, "scroll", () => this.position(), true));
        for (const closer of qsa(this.panel, "[data-fui-popover-close]")) {
            this.disposer.add(on(closer, "click", (event) => {
                event.preventDefault();
                this.close();
            }));
        }
    }
    get isOpen() {
        return this.openState;
    }
    open() {
        if (this.native) {
            try {
                this.panel.showPopover?.();
            }
            catch {
                // InvalidStateError means it is already open.
            }
        }
        else {
            this.panel.hidden = false;
        }
        this.syncOpen(true);
    }
    close() {
        if (this.native) {
            try {
                this.panel.hidePopover?.();
            }
            catch {
                // Nothing to hide.
            }
        }
        else {
            this.panel.hidden = true;
        }
        this.syncOpen(false);
    }
    toggle() {
        if (this.openState) {
            this.close();
        }
        else {
            this.open();
        }
    }
    dispose() {
        this.close();
        this.disposer.dispose();
        controllers.delete(this.trigger);
    }
    syncOpen(open) {
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
        }
        else {
            delete this.panel.dataset["fuiState"];
            delete this.panel.dataset["fuiSide"];
            emit(this.trigger, "fui:popoverclose", { trigger: this.trigger, panel: this.panel });
        }
    }
    position() {
        if (!this.openState) {
            return;
        }
        const offset = this.options.offset ?? 8;
        const margin = 8;
        const triggerRect = this.trigger.getBoundingClientRect();
        const panelRect = this.panel.getBoundingClientRect();
        let side = "bottom";
        let top = triggerRect.bottom + offset;
        if (top + panelRect.height > window.innerHeight - margin && triggerRect.top - panelRect.height - offset >= margin) {
            top = triggerRect.top - panelRect.height - offset;
            side = "top";
        }
        const left = Math.min(Math.max(margin, triggerRect.left), Math.max(margin, window.innerWidth - panelRect.width - margin));
        this.panel.style.position = "fixed";
        this.panel.style.left = `${Math.round(left)}px`;
        this.panel.style.top = `${Math.round(top)}px`;
        this.panel.style.right = "auto";
        this.panel.style.bottom = "auto";
        this.panel.style.margin = "0";
        this.panel.dataset["fuiSide"] = side;
    }
}
/** Returns the controller bound to a trigger, if any. */
export function getPopover(trigger) {
    return controllers.get(trigger);
}
/** Creates or reuses the controller for a trigger and its panel. */
export function bindPopover(trigger, panel, options = {}) {
    const existing = controllers.get(trigger);
    if (existing) {
        return existing;
    }
    const controller = new PopoverController(trigger, panel, options);
    controllers.set(trigger, controller);
    return controller;
}
/** Binds popovers declared with data-fui-popover="<target>". */
export function initPopovers(root = document) {
    if (typeof document === "undefined") {
        return () => undefined;
    }
    const disposer = createDisposer();
    for (const trigger of qsa(root, "[data-fui-popover]")) {
        if (controllers.has(trigger)) {
            continue;
        }
        const value = trigger.getAttribute("data-fui-popover") ??
            trigger.getAttribute("aria-controls") ??
            trigger.getAttribute("popovertarget");
        const panel = getElement(value);
        if (!panel) {
            continue;
        }
        const controller = bindPopover(trigger, panel);
        disposer.add(() => controller.dispose());
    }
    return () => disposer.dispose();
}
