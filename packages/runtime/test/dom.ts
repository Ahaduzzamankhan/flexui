/**
 * Minimal DOM environment for the runtime tests.
 *
 * happy-dom is a development-only dependency: the published runtime still has
 * no dependencies. Installing the globals before a test body runs is enough,
 * because every runtime module reads `document` when a function is called.
 */

import { Window } from "happy-dom";

let current: Window | null = null;

export function setupDOM(): Window {
  const window = new Window({ url: "https://example.test/" });
  current = window;

  const globals: Record<string, unknown> = {
    window,
    document: window.document,
    navigator: window.navigator,
    localStorage: window.localStorage,
    Node: window.Node,
    Element: window.Element,
    HTMLElement: window.HTMLElement,
    HTMLDialogElement: window.HTMLDialogElement,
    HTMLInputElement: window.HTMLInputElement,
    Event: window.Event,
    CustomEvent: window.CustomEvent,
    MouseEvent: window.MouseEvent,
    KeyboardEvent: window.KeyboardEvent,
    FocusEvent: window.FocusEvent,
    getComputedStyle: window.getComputedStyle.bind(window),
  };
  for (const [key, value] of Object.entries(globals)) {
    Object.defineProperty(globalThis, key, { value, writable: true, configurable: true });
  }
  return window;
}

export function teardownDOM(): void {
  current?.document.body.replaceChildren();
  current = null;
}

/** Creates an element, appends it to the body and returns it. */
export function mount<T extends "dialog" | "div" | "button" | "input">(
  tag: T,
  attributes: Record<string, string> = {},
  html = "",
): HTMLElement {
  const element = document.createElement(tag);
  for (const [name, value] of Object.entries(attributes)) {
    element.setAttribute(name, value);
  }
  element.innerHTML = html;
  document.body.append(element);
  return element;
}

/** Dispatches a bubbling keydown event on an element. */
export function press(element: Element, key: string, init: KeyboardEventInit = {}): void {
  const KeyboardEventCtor = globalThis.KeyboardEvent as unknown as new (
    type: string,
    init: Record<string, unknown>,
  ) => Event;
  element.dispatchEvent(
    new KeyboardEventCtor("keydown", {
      key,
      bubbles: true,
      cancelable: true,
      ...init,
    }),
  );
}
