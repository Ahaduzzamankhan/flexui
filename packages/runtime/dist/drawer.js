import { bindDialog, bindOverlayTriggers, getDialog } from "./dialog.js";
import { createDisposer, qsa } from "./types.js";
/** Reads the side declared by data-fui-side, defaulting to right. */
export function drawerSide(element) {
    return element.getAttribute("data-fui-side") === "left" ? "left" : "right";
}
/** Binds every drawer under root, including its open/close triggers. */
export function initDrawers(root = document) {
    if (typeof document === "undefined") {
        return () => undefined;
    }
    const disposer = createDisposer();
    for (const element of qsa(root, "dialog[data-fui-drawer]")) {
        element.dataset["fuiSide"] = drawerSide(element);
        if (getDialog(element) !== undefined) {
            continue;
        }
        const controller = bindDialog(element);
        disposer.add(() => controller.dispose());
    }
    disposer.add(bindOverlayTriggers());
    return () => disposer.dispose();
}
