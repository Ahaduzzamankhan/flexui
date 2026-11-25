import { getDialog, getResolvedTheme, getTheme, init, toast, toggleTheme, } from "../../../packages/runtime/dist/index.js";
const log = document.querySelector("[data-event-log]");
const themeStatus = document.querySelector("[data-theme-status]");
function record(message) {
    if (!log) {
        return;
    }
    const now = new Date();
    const item = document.createElement("li");
    item.className = "fui-log-item";
    const time = document.createElement("time");
    time.className = "fui-log-time";
    time.dateTime = now.toISOString();
    time.textContent = now.toLocaleTimeString(undefined, {
        hour: "2-digit",
        minute: "2-digit",
        second: "2-digit",
        hour12: false,
    });
    const text = document.createElement("span");
    text.textContent = message;
    item.append(time, text);
    log.prepend(item);
    while (log.children.length > 8) {
        log.lastElementChild?.remove();
    }
}
function updateThemeStatus() {
    if (themeStatus) {
        themeStatus.textContent = `${getTheme()} (${getResolvedTheme()})`;
    }
}
function openDialog(selector) {
    const element = document.querySelector(selector);
    if (element instanceof HTMLDialogElement) {
        getDialog(element)?.open();
    }
}
document.addEventListener("fui:open", (event) => {
    const detail = event.detail;
    record(`dialog opened: #${detail.dialog.id || "(unnamed)"}`);
});
document.addEventListener("fui:close", (event) => {
    const detail = event.detail;
    record(`dialog closed: #${detail.dialog.id || "(unnamed)"} (${detail.reason})`);
});
document.addEventListener("fui:tabchange", (event) => {
    const detail = event.detail;
    record(`tab changed: ${detail.tab.textContent?.trim() ?? detail.index}`);
});
document.addEventListener("fui:themechange", (event) => {
    const detail = event.detail;
    updateThemeStatus();
    record(`theme: ${detail.mode} -> ${detail.resolved}`);
});
document.addEventListener("fui:select", (event) => {
    const detail = event.detail;
    record(`menu selection: ${detail.item.textContent?.trim() ?? "item"}`);
});
document.addEventListener("fui:popoveropen", () => record("popover opened"));
document.addEventListener("fui:popoverclose", () => record("popover closed"));
document.addEventListener("fui:toast", (event) => {
    const detail = event.detail;
    record(`toast shown: ${detail.variant}`);
});
document.addEventListener("fui:command", (event) => {
    const detail = event.detail;
    record(`command: ${detail.label}`);
    switch (detail.value) {
        case "dialog":
            openDialog("#demo-dialog");
            break;
        case "drawer":
            openDialog("#settings-drawer");
            break;
        case "theme":
            toggleTheme();
            break;
        case "toast":
            toast({
                variant: "success",
                title: "Workspace created",
                description: "The new workspace is ready for your team.",
            });
            break;
        case "tabs":
            document.getElementById("tabs-card")?.scrollIntoView({ behavior: "smooth", block: "start" });
            document.getElementById("tab-behaviour")?.focus({ preventScroll: true });
            break;
        default:
            break;
    }
});
document.addEventListener("click", (event) => {
    const target = event.target;
    if (!(target instanceof Element)) {
        return;
    }
    const demo = target.closest("[data-toast-demo]");
    if (demo) {
        switch (demo.getAttribute("data-toast-demo")) {
            case "success":
                toast({
                    variant: "success",
                    title: "Changes saved",
                    description: "Your workspace preferences were updated.",
                });
                break;
            case "danger":
                toast({
                    variant: "danger",
                    title: "Deployment failed",
                    description: "The build exited with code 1. Check the logs and retry.",
                });
                break;
            case "sticky":
                toast({
                    variant: "info",
                    title: "Sticky notification",
                    description: "This toast stays until you dismiss it.",
                    duration: 0,
                });
                break;
            default:
                toast({
                    variant: "info",
                    title: "Heads up",
                    description: "Toasts pause while hovered or focused.",
                });
                break;
        }
        return;
    }
    if (target.closest("[data-drawer-save]")) {
        const drawer = document.getElementById("settings-drawer");
        if (drawer instanceof HTMLDialogElement) {
            getDialog(drawer)?.close("save");
        }
        toast({
            variant: "success",
            title: "Preferences saved",
            description: "Your drawer settings are up to date.",
        });
    }
    if (target.closest("[data-clear-log]")) {
        log?.replaceChildren();
        record("event log cleared");
    }
});
document.querySelector("[data-demo-form]")?.addEventListener("submit", (event) => {
    event.preventDefault();
    const form = event.currentTarget;
    if (!(form instanceof HTMLFormElement)) {
        return;
    }
    const dialog = form.closest("dialog");
    if (dialog) {
        getDialog(dialog)?.close("submit");
    }
    const name = new FormData(form).get("name");
    toast({
        variant: "success",
        title: "Workspace created",
        description: `${typeof name === "string" && name.trim() ? name.trim() : "Your workspace"} is ready.`,
    });
});
updateThemeStatus();
init(document);
record("runtime initialized");
