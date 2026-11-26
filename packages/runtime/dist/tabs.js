import { createDisposer, emit, on, qs, qsa, uid } from "./types.js";
const controllers = new WeakMap();
/** Accessible tabs with roving tabindex and automatic activation. */
export class Tabs {
    root;
    disposer;
    tabs;
    panels;
    activeIndex = 0;
    constructor(root) {
        this.root = root;
        this.disposer = createDisposer();
        const tablist = root.querySelector('[role="tablist"], [data-fui-tablist]') ??
            (root.getAttribute("role") === "tablist" ? root : null);
        this.tabs = tablist ? qsa(tablist, '[role="tab"], [data-fui-tab]') : [];
        this.panels = [];
        if (this.tabs.length === 0) {
            return;
        }
        for (const [index, tab] of this.tabs.entries()) {
            if (!tab.id) {
                tab.id = uid("fui-tab");
            }
            tab.setAttribute("role", "tab");
            tab.tabIndex = -1;
            const controls = tab.getAttribute("aria-controls");
            const panel = controls ? document.getElementById(controls) : null;
            if (panel) {
                panel.setAttribute("role", "tabpanel");
                panel.setAttribute("aria-labelledby", tab.id);
                if (!panel.id) {
                    panel.id = uid("fui-tabpanel");
                }
                if (!controls) {
                    tab.setAttribute("aria-controls", panel.id);
                }
            }
            this.panels[index] = panel;
            this.disposer.add(on(tab, "click", () => {
                this.activate(index, true);
            }));
            this.disposer.add(on(tab, "keydown", (event) => {
                this.handleKeydown(event, index, tablist);
            }));
        }
        const selected = this.tabs.findIndex((tab) => tab.getAttribute("aria-selected") === "true");
        this.activate(selected >= 0 ? selected : 0, false, false);
    }
    get active() {
        const tab = this.tabs[this.activeIndex];
        if (!tab) {
            return null;
        }
        return { tab, panel: this.panels[this.activeIndex] ?? null, index: this.activeIndex };
    }
    activate(index, focus = false, notify = true) {
        const next = this.tabs[index];
        if (!next) {
            return;
        }
        const changed = index !== this.activeIndex;
        this.activeIndex = index;
        for (const [itemIndex, tab] of this.tabs.entries()) {
            const selected = itemIndex === index;
            tab.setAttribute("aria-selected", selected ? "true" : "false");
            tab.tabIndex = selected ? 0 : -1;
            const panel = this.panels[itemIndex];
            if (panel) {
                panel.hidden = !selected;
            }
        }
        if (focus) {
            next.focus({ preventScroll: true });
        }
        if (notify && (changed || focus)) {
            emit(this.root, "fui:tabchange", {
                tabs: this,
                tab: next,
                panel: this.panels[index] ?? null,
                index,
            });
        }
    }
    dispose() {
        this.disposer.dispose();
        controllers.delete(this.root);
    }
    handleKeydown(event, index, tablist) {
        const vertical = tablist?.getAttribute("aria-orientation") === "vertical";
        let nextIndex = index;
        switch (event.key) {
            case "ArrowRight":
                if (vertical) {
                    return;
                }
                nextIndex = (index + 1) % this.tabs.length;
                break;
            case "ArrowLeft":
                if (vertical) {
                    return;
                }
                nextIndex = index === 0 ? this.tabs.length - 1 : index - 1;
                break;
            case "ArrowDown":
                if (!vertical) {
                    return;
                }
                nextIndex = (index + 1) % this.tabs.length;
                break;
            case "ArrowUp":
                if (!vertical) {
                    return;
                }
                nextIndex = index === 0 ? this.tabs.length - 1 : index - 1;
                break;
            case "Home":
                nextIndex = 0;
                break;
            case "End":
                nextIndex = this.tabs.length - 1;
                break;
            default:
                return;
        }
        event.preventDefault();
        this.activate(nextIndex, true);
    }
}
/** Returns the controller bound to a tabs container, if any. */
export function getTabs(root) {
    return controllers.get(root);
}
/** Creates or reuses the controller for a tabs container. */
export function bindTabs(root) {
    const existing = controllers.get(root);
    if (existing) {
        return existing;
    }
    const controller = new Tabs(root);
    controllers.set(root, controller);
    return controller;
}
/** Binds every [data-fui-tabs] container under root. */
export function initTabs(root = document) {
    if (typeof document === "undefined") {
        return () => undefined;
    }
    const disposer = createDisposer();
    for (const element of qsa(root, "[data-fui-tabs]")) {
        if (controllers.has(element)) {
            continue;
        }
        if (!qs(element, '[role="tablist"], [data-fui-tablist]') && element.getAttribute("role") !== "tablist") {
            continue;
        }
        const controller = bindTabs(element);
        disposer.add(() => controller.dispose());
    }
    return () => disposer.dispose();
}
