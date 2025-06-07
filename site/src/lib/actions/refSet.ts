import type { Action } from 'svelte/action';

/**
 * Svelte action that adds the element to a `Set<HTMLElement>` when the
 * element is created, and removes it when the element is destroyed.
 */
export const refSet: Action<HTMLElement, Set<HTMLElement>> = (node, set) => {
  set.add(node);
  return {
    destroy() {
      set.delete(node);
    }
  };
};
