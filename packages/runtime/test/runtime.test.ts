import { beforeEach, describe, expect, test } from "bun:test";

import { setupDOM, mount, press } from "./dom.js";
import { createDisposer } from "../src/types.js";
import { bindDialog, getDialog, initDialogs } from "../src/dialog.js";
import { bindDropdown, getDropdown } from "../src/dropdown.js";
import { bindTabs, getTabs } from "../src/tabs.js";
import { toast } from "../src/toast.js";
import { getTheme, initTheme, setTheme } from "../src/theme.js";

beforeEach(() => {
  setupDOM();
});

function dialogFixture(): HTMLDialogElement {
  return mount(
    "dialog",
    { id: "invite", "data-fui-dialog": "" },
    `<h2>Invite</h2><button type="button" id="invite-ok">Send</button>
     <button type="button" id="invite-cancel" data-fui-close>Cancel</button>`,
  ) as HTMLDialogElement;
}

describe("disposer", () => {
  test("runs cleanups in reverse order", () => {
    const order: number[] = [];
    const disposer = createDisposer();
    disposer.add(() => order.push(1));
    disposer.add(() => order.push(2));
    expect(disposer.size).toBe(2);
    disposer.dispose();
    expect(order).toEqual([2, 1]);
    expect(disposer.size).toBe(0);
  });

  test("runs a cleanup immediately when already disposed", () => {
    const disposer = createDisposer();
    disposer.dispose();
    let ran = false;
    disposer.add(() => {
      ran = true;
    });
    expect(ran).toBe(true);
  });
});

describe("dialog", () => {
  test("opens, emits fui:open and locks page scroll", () => {
    const element = dialogFixture();
    const controller = bindDialog(element);
    const events: string[] = [];
    element.addEventListener("fui:open", () => events.push("open"));

    controller.open();

    expect(controller.isOpen).toBe(true);
    expect(element.hasAttribute("open")).toBe(true);
    expect(document.documentElement.dataset["fuiOverlayOpen"]).toBe("");
    expect(events).toEqual(["open"]);
    controller.dispose();
  });

  test("keeps the scroll lock until the last overlay closes", () => {
    const first = bindDialog(dialogFixture());
    const second = bindDialog(dialogFixture());

    first.open();
    second.open();
    first.close("dismiss");
    expect(document.documentElement.hasAttribute("data-fui-overlay-open")).toBe(true);

    second.close("dismiss");
    expect(document.documentElement.hasAttribute("data-fui-overlay-open")).toBe(false);
    first.dispose();
    second.dispose();
  });

  test("closes with a reason, unlocks scroll and emits fui:close", () => {
    const element = dialogFixture();
    const controller = bindDialog(element);
    let reason = "";
    element.addEventListener("fui:close", (event) => {
      reason = (event as CustomEvent<{ reason: string }>).detail.reason;
    });

    controller.open();
    controller.close("saved");

    expect(controller.isOpen).toBe(false);
    expect(element.hasAttribute("open")).toBe(false);
    expect(document.documentElement.hasAttribute("data-fui-overlay-open")).toBe(false);
    expect(reason).toBe("saved");
    controller.dispose();
  });

  test("reuses one controller per element and forgets it on dispose", () => {
    const element = dialogFixture();
    const first = bindDialog(element);
    expect(bindDialog(element)).toBe(first);
    expect(getDialog(element)).toBe(first);
    first.dispose();
    expect(getDialog(element)).toBeUndefined();
  });

  test("closing an already closed dialog is a no-op", () => {
    const element = dialogFixture();
    const controller = bindDialog(element);
    let closed = 0;
    element.addEventListener("fui:close", () => {
      closed += 1;
    });
    controller.close("dismiss");
    expect(closed).toBe(0);
  });

  test("initDialogs binds marked dialogs and wires aria on triggers", () => {
    const element = dialogFixture();
    mount("button", { id: "open-invite", "data-fui-open": "#invite" }, "Open");
    const dispose = initDialogs(document);
    expect(getDialog(element)).toBeDefined();

    const trigger = document.getElementById("open-invite") as HTMLElement;
    expect(trigger.getAttribute("aria-haspopup")).toBe("dialog");
    expect(trigger.getAttribute("aria-controls")).toBe("invite");

    trigger.click();
    expect(getDialog(element)?.isOpen).toBe(true);
    dispose();
  });
});

describe("dropdown", () => {
  function dropdownFixture(): { trigger: HTMLElement; menu: HTMLElement } {
    const trigger = mount("button", { id: "account" }, "Account");
    const menu = mount(
      "div",
      { id: "account-menu", role: "menu", hidden: "" },
      `<button role="menuitem" id="item-a">Profile</button>
       <button role="menuitem" id="item-b">Sign out</button>`,
    );
    return { trigger, menu };
  }

  test("open sets aria-expanded and shows the menu", () => {
    const { trigger, menu } = dropdownFixture();
    const controller = bindDropdown(trigger, menu);

    expect(trigger.getAttribute("aria-haspopup")).toBe("menu");
    expect(trigger.getAttribute("aria-controls")).toBe("account-menu");
    expect(trigger.getAttribute("aria-expanded")).toBe("false");

    controller.open();
    expect(controller.isOpen).toBe(true);
    expect(menu.hidden).toBe(false);
    expect(trigger.getAttribute("aria-expanded")).toBe("true");
    controller.dispose();
    expect(getDropdown(trigger)).toBeUndefined();
  });

  test("ArrowDown moves to the next item and Escape closes", () => {
    const { trigger, menu } = dropdownFixture();
    const controller = bindDropdown(trigger, menu);
    controller.open();

    press(menu, "ArrowDown");
    expect(document.activeElement?.id).toBe("item-b");
    press(menu, "Home");
    expect(document.activeElement?.id).toBe("item-a");

    press(menu, "Escape");
    expect(controller.isOpen).toBe(false);
    expect(trigger.getAttribute("aria-expanded")).toBe("false");
    controller.dispose();
  });

  test("clicking an item emits fui:select and closes", () => {
    const { trigger, menu } = dropdownFixture();
    const controller = bindDropdown(trigger, menu);
    let value = "";
    menu.addEventListener("fui:select", (event) => {
      value = (event as CustomEvent<{ item: HTMLElement }>).detail.item.id;
    });

    controller.open();
    (document.getElementById("item-b") as HTMLElement).click();

    expect(value).toBe("item-b");
    expect(controller.isOpen).toBe(false);
    controller.dispose();
  });

  test("disabled items are skipped", () => {
    const trigger = mount("button", {});
    const menu = mount(
      "div",
      { hidden: "" },
      `<button role="menuitem" id="disabled" aria-disabled="true">Off</button>
       <button role="menuitem" id="enabled">On</button>`,
    );
    const controller = bindDropdown(trigger, menu);
    controller.open();
    expect(document.activeElement?.id).toBe("enabled");
    controller.dispose();
  });
});

describe("tabs", () => {
  function tabsFixture(): HTMLElement {
    return mount(
      "div",
      { "data-fui-tabs": "" },
      `<div role="tablist">
         <button role="tab" id="tab-a" aria-controls="panel-a" aria-selected="true">A</button>
         <button role="tab" id="tab-b" aria-controls="panel-b">B</button>
       </div>
       <div role="tabpanel" id="panel-a">A body</div>
       <div role="tabpanel" id="panel-b">B body</div>`,
    );
  }

  test("activates the selected tab and hides the other panel", () => {
    const root = tabsFixture();
    const controller = bindTabs(root);

    expect(controller.active?.index).toBe(0);
    expect(document.getElementById("panel-a")?.hidden).toBe(false);
    expect(document.getElementById("panel-b")?.hidden).toBe(true);
    expect(getTabs(root)).toBe(controller);
    controller.dispose();
  });

  test("activate updates ARIA and emits fui:tabchange", () => {
    const root = tabsFixture();
    const controller = bindTabs(root);
    let index = -1;
    root.addEventListener("fui:tabchange", (event) => {
      index = (event as CustomEvent<{ index: number }>).detail.index;
    });

    controller.activate(1);

    expect(index).toBe(1);
    expect(document.getElementById("tab-b")?.getAttribute("aria-selected")).toBe("true");
    expect(document.getElementById("tab-a")?.getAttribute("aria-selected")).toBe("false");
    expect(document.getElementById("panel-b")?.hidden).toBe(false);
    controller.dispose();
  });

  test("ArrowRight wraps from the last tab to the first", () => {
    const root = tabsFixture();
    const controller = bindTabs(root);
    const last = document.getElementById("tab-b") as HTMLElement;
    controller.activate(1);
    press(last, "ArrowRight");
    expect(controller.active?.index).toBe(0);
    controller.dispose();
  });

  test("generates ids when the markup omits them", () => {
    const root = mount(
      "div",
      { "data-fui-tabs": "" },
      `<div role="tablist"><button role="tab">One</button></div><div>Panel</div>`,
    );
    const controller = bindTabs(root);
    const tab = root.querySelector('[role="tab"]') as HTMLElement;
    expect(tab.id.startsWith("fui-tab")).toBe(true);
    controller.dispose();
  });
});

describe("toast", () => {
  test("creates a live region and announces a toast", () => {
    const handle = toast({ title: "Saved", description: "All good", variant: "success", duration: 0 });
    const region = document.querySelector("[data-fui-toast-region]") as HTMLElement;

    expect(region).toBeDefined();
    expect(region.getAttribute("aria-live")).toBe("polite");
    expect(handle.element.className).toContain("fui-toast--success");
    expect(handle.element.querySelector(".fui-toast-title")?.textContent).toBe("Saved");
    expect(handle.element.querySelector(".fui-toast-description")?.textContent).toBe("All good");

    handle.dismiss();
    expect(handle.element.dataset["fuiToastState"]).toBe("closed");
  });

  test("danger toasts are assertive alerts", () => {
    const handle = toast({ title: "Failed", variant: "danger", duration: 0 });
    expect(handle.element.getAttribute("role")).toBe("alert");
    handle.dismiss();
  });
});

describe("theme", () => {
  test("persists to the storage key configured by initTheme", () => {
    const dispose = initTheme({ storageKey: "custom-theme" });
    setTheme("dark");

    expect(window.localStorage.getItem("custom-theme")).toBe("dark");
    expect(getTheme()).toBe("dark");
    expect(document.documentElement.classList.contains("dark")).toBe(true);
    expect(document.documentElement.dataset["fuiThemeMode"]).toBe("dark");
    dispose();
  });

  test("toggling follows light -> dark -> system", () => {
    const dispose = initTheme();
    setTheme("light");
    const root = document.documentElement;
    expect(setTheme("dark")).toBe("dark");
    expect(root.classList.contains("dark")).toBe(true);
    expect(setTheme("system")).toBe("system");
    expect(root.style.colorScheme).toBe("light dark");
    dispose();
  });

  test("theme controls update aria-pressed", () => {
    mount("button", { "data-fui-theme-set": "dark" }, "Dark");
    const dispose = initTheme();
    setTheme("dark");
    const button = document.querySelector("[data-fui-theme-set='dark']") as HTMLElement;
    expect(button.getAttribute("aria-pressed")).toBe("true");
    dispose();
  });
});

describe("focus helpers", () => {
  test("trapFocus cycles Tab inside the container", async () => {
    const { trapFocus } = await import("../src/focus.js");
    const container = mount(
      "div",
      {},
      `<button id="first">First</button><button id="last">Last</button>`,
    );
    const release = trapFocus(container);
    (document.getElementById("last") as HTMLElement).focus();

    press(container, "Tab");
    expect(document.activeElement?.id).toBe("first");

    press(container, "Tab", { shiftKey: true });
    expect(document.activeElement?.id).toBe("last");
    release();
  });
});
