<script lang="ts">
  import { onMount } from "svelte";
  import { debounce } from "../utils";
  import { refSet } from "./actions/refSet";

  interface Props {
    solutions: Map<string, string[][]>;
    class: string | undefined;
  }

  const { solutions, class: cls = "" }: Props = $props();

  /** The scrollable element containing the solutions. */
  let solScroller: HTMLElement;
  /** The element that indicates the selected solution on mobile. */
  let solSelector: HTMLElement;
  /** The solutions that are fully visible in `solScroller`. */
  const visibleSolutions = new Set<HTMLElement>();
  /** The currently selected solution. */
  let selectedSol: HTMLElement;

  let pointerMediaQuery: MediaQueryList;
  let hasFinePointer = true;

  onMount(() => {
    pointerMediaQuery = window.matchMedia("(pointer: fine)");
    hasFinePointer = pointerMediaQuery.matches;
    pointerMediaQuery.addEventListener("change", (e) => (hasFinePointer = e.matches));

    const observer = new IntersectionObserver(
      (entries) => {
        if (hasFinePointer) return;
        updateVisibleSolutions(entries, visibleSolutions);
      },
      {
        root: solScroller,
        rootMargin: "0px",
        threshold: 1,
      }
    );

    for (const el of visibleSolutions) {
      observer.observe(el);
    }
  });

  const debouncedSnapSolScroller = debounce(() => snapSolScroller(solScroller, solSelector, selectedSol), 200);
  const handleScroll = () => {
    if (hasFinePointer) return;

    requestIdleCallback(() => (selectedSol = getSelectedSolution(visibleSolutions, solSelector) || selectedSol));
    debouncedSnapSolScroller();
  };

  /** Scroll the `solScroller` such that the `selectedSol` directly overlaps with `solSelector`. */
  function snapSolScroller(solScroller: HTMLElement, solSelector: HTMLElement, selectedSol: HTMLElement) {
    const scrollY =
      selectedSol.offsetTop + solScroller.getBoundingClientRect().top - solSelector.getBoundingClientRect().top;
    solScroller.scroll({ top: scrollY, behavior: "smooth" });
  }

  /** Update the set of visible solutions based on the entries provided by an `IntersectionObserver`. */
  function updateVisibleSolutions(entries: IntersectionObserverEntry[], visibleSolutions: Set<HTMLElement>) {
    for (const entry of entries) {
      if (entry.isIntersecting) visibleSolutions.add(entry.target as HTMLElement);
      else visibleSolutions.delete(entry.target as HTMLElement);
    }
  }

  /** Return the solution that overlaps the most with `solSelector`, or `null` if no solution overlaps. */
  function getSelectedSolution(visibleSolutions: Set<HTMLElement>, solSelector: HTMLElement): HTMLElement | null {
    let maxIntersectionRatio = 0;
    let selectedSol: HTMLElement | null = null;

    for (const sol of visibleSolutions) {
      const { bottom: solBottom, top: solTop } = sol.getBoundingClientRect();
      const { bottom: selectorBottom, top: selectorTop, height: selectorHeight } = solSelector.getBoundingClientRect();

      const intersectionRatio =
        Math.max(0, Math.min(solBottom, selectorBottom) - Math.max(solTop, selectorTop)) / selectorHeight;

      if (intersectionRatio > maxIntersectionRatio) {
        selectedSol = sol;
        maxIntersectionRatio = intersectionRatio;
      }
    }

    return selectedSol;
  }
</script>

<div class={["relative", cls]}>
  <div
    bind:this={solSelector}
    class="md:hidden pointer-fine:hidden -z-1 h-9 left-4 right-4 top-14 bg-rose-100 absolute"
  ></div>

  <div bind:this={solScroller} onscroll={handleScroll} class="overflow-y-scroll">
    {#each solutions.entries() as [date, sols]}
      <div class="pt-6 last:pb-[calc(100%-4rem)]">
        <span class="font-bold px-8 md:px-10 block mb-2">
          {new Date(date).toLocaleDateString("en-US", { year: "numeric", month: "long", day: "numeric" })}
        </span>
        <ul class="flex flex-col px-4 md:px-6">
          {#each sols as words}
            <li
              use:refSet={visibleSolutions}
              class="solution px-4 py-1.5 not-first:-mt-1.5 pointer-fine:hover:bg-rose-100 tracking-wider"
            >
              {words.join(" – ")}
            </li>
          {/each}
        </ul>
      </div>
    {/each}
  </div>
</div>
