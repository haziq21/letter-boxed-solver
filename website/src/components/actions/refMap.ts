import type { ActionReturn } from "svelte/action";

/**
 * Svelte action that adds the element to a `Map<HTMLElement, T>` when
 * the element is created, and removes it when the element is destroyed.
 */
export function refMap<T>(node: HTMLElement, param: { map: Map<HTMLElement, T>; value: T }): ActionReturn<HTMLElement> {
  param.map.set(node, param.value);
  return {
    destroy() {
      param.map.delete(node);
    },
  };
}
