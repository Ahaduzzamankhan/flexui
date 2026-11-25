import { closestFromEvent, createDisposer, emit, on, qsa, } from "./types.js";
const active = [];
const MAX_TOASTS = 4;
const DEFAULT_DURATION = 5000;
let region = null;
function ensureRegion() {
    if (region?.isConnected) {
        return region;
    }
    const existing = document.querySelector("[data-fui-toast-region], .fui-toast-region");
    if (existing) {
        region = existing;
        if (!existing.hasAttribute("data-fui-toast-region")) {
            existing.setAttribute("data-fui-toast-region", "");
        }
        return existing;
    }
    const created = document.createElement("div");
    created.className = "fui-toast-region";
    created.setAttribute("data-fui-toast-region", "");
    created.setAttribute("role", "region");
    created.setAttribute("aria-label", "Notifications");
    created.setAttribute("aria-live", "polite");
    created.setAttribute("aria-atomic", "false");
    document.body.append(created);
    region = created;
    return created;
}
function dismissElement(element) {
    if (element.dataset["fuiToastState"] === "closed") {
        return;
    }
    element.dataset["fuiToastState"] = "closed";
    emit(element, "fui:toastdismiss", { toast: element });
    window.setTimeout(() => {
        element.remove();
    }, 160);
}
/** Creates and announces a toast, creating the live region on first use. */
export function toast(options) {
    const container = ensureRegion();
    const variant = options.variant ?? "info";
    const duration = options.duration ?? DEFAULT_DURATION;
    while (active.length >= MAX_TOASTS) {
        const oldest = active[0];
        if (!oldest) {
            break;
        }
        oldest.dismiss();
    }
    const element = document.createElement("div");
    element.className = `fui-toast fui-toast--${variant}`;
    element.setAttribute("data-fui-toast", "");
    element.setAttribute("role", variant === "danger" ? "alert" : "status");
    const content = document.createElement("div");
    content.className = "fui-toast-content";
    const title = document.createElement("p");
    title.className = "fui-toast-title";
    title.textContent = options.title;
    content.append(title);
    if (options.description) {
        const description = document.createElement("p");
        description.className = "fui-toast-description";
        description.textContent = options.description;
        content.append(description);
    }
    element.append(content);
    let timer = null;
    let remaining = Math.max(0, duration);
    let startedAt = 0;
    let dismissed = false;
    const dismiss = () => {
        if (dismissed) {
            return;
        }
        dismissed = true;
        if (timer !== null) {
            window.clearTimeout(timer);
            timer = null;
        }
        const index = active.findIndex((item) => item.element === element);
        if (index >= 0) {
            active.splice(index, 1);
        }
        dismissElement(element);
    };
    const start = () => {
        if (dismissed || remaining <= 0) {
            return;
        }
        startedAt = Date.now();
        timer = window.setTimeout(dismiss, remaining);
    };
    const pause = () => {
        if (timer === null) {
            return;
        }
        window.clearTimeout(timer);
        timer = null;
        remaining = Math.max(0, remaining - (Date.now() - startedAt));
    };
    const resume = () => {
        if (!dismissed) {
            start();
        }
    };
    if (options.dismissible !== false) {
        const close = document.createElement("button");
        close.type = "button";
        close.className = "fui-toast-close";
        close.setAttribute("data-fui-toast-dismiss", "");
        close.setAttribute("aria-label", "Dismiss notification");
        close.textContent = "\u00d7";
        close.addEventListener("click", dismiss);
        element.append(close);
    }
    element.addEventListener("mouseenter", pause);
    element.addEventListener("mouseleave", resume);
    element.addEventListener("focusin", pause);
    element.addEventListener("focusout", (event) => {
        const next = event.relatedTarget;
        if (!(next instanceof Node) || !element.contains(next)) {
            resume();
        }
    });
    container.append(element);
    active.push({ element, dismiss });
    start();
    emit(element, "fui:toast", { toast: element, variant });
    return { element, dismiss };
}
/** Binds static toast dismiss buttons and auto-dismiss timers. */
export function initToasts(root = document) {
    if (typeof document === "undefined") {
        return () => undefined;
    }
    ensureRegion();
    const disposer = createDisposer();
    const timers = [];
    disposer.add(on(document, "click", (event) => {
        const close = closestFromEvent(event, "[data-fui-toast-dismiss]");
        const element = close?.closest("[data-fui-toast]");
        if (element) {
            event.preventDefault();
            dismissElement(element);
        }
    }));
    for (const element of qsa(root, "[data-fui-toast]")) {
        if (element.dataset["fuiToastBound"] === "true") {
            continue;
        }
        element.dataset["fuiToastBound"] = "true";
        const duration = Number.parseInt(element.getAttribute("data-fui-toast-duration") ?? "", 10);
        if (Number.isFinite(duration) && duration > 0) {
            timers.push(window.setTimeout(() => {
                dismissElement(element);
            }, duration));
        }
    }
    return () => {
        disposer.dispose();
        for (const timer of timers) {
            window.clearTimeout(timer);
        }
    };
}
