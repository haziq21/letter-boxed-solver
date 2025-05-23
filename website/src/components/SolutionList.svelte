<script lang="ts">
  import { onMount } from "svelte";
  import { debounce } from "../utils";
  import { refSet, refMap } from "./actions";

  interface Props {
    solutions: Map<string, string[][]>;
    selected?: { sol: string[]; date: string };
    class: string | undefined;
  }

  let { solutions, selected = $bindable(), class: cls = "" }: Props = $props();

  /** The scrollable element containing the solutions. */
  let solScrollerElem: HTMLElement;
  /** The element that indicates the selected solution on mobile. */
  let solSelectorElem: HTMLElement;
  /** The solutions that are fully visible in `solScroller`. */
  const visibleSolElems = new Set<HTMLElement>();
  /** The currently selected solution. */
  let selectedSolElem: HTMLElement | undefined = $state();
  let hoveredSolElem: HTMLElement | undefined = $state();
  /** A map of solution elements to their corresponding data values. */
  const solElemData = new Map<HTMLElement, { sol: string[]; date: string }>();
  /** The result of the `(pointer: fine)` media query. */
  let hasFinePointer = $state(true);

  onMount(() => {
    const pointerMediaQuery = window.matchMedia("(pointer: fine)");
    hasFinePointer = pointerMediaQuery.matches;
    pointerMediaQuery.addEventListener("change", (e) => (hasFinePointer = e.matches));

    const observer = new IntersectionObserver(
      (entries) => {
        if (hasFinePointer) return;
        updateVisibleSolutions(entries, visibleSolElems);
      },
      {
        root: solScrollerElem,
        rootMargin: "0px",
        threshold: 1,
      }
    );

    for (const el of visibleSolElems) {
      observer.observe(el);
    }

    selectedSolElem = hasFinePointer
      ? solElemData.keys().next().value!
      : getSelectedSolution(visibleSolElems, solSelectorElem)!;
    selected = solElemData.get(selectedSolElem)!;
  });

  const debouncedSnapSolScroller = debounce(
    () => snapSolScroller(solScrollerElem, solSelectorElem, selectedSolElem!),
    200
  );
  const handleScroll = () => {
    if (hasFinePointer) return;

    requestIdleCallback(() => {
      selectedSolElem = getSelectedSolution(visibleSolElems, solSelectorElem) || selectedSolElem;
      selected = solElemData.get(selectedSolElem!)!;
    });
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

  /** Return the solution element the closest to `solSelector`, or `null` if `visibleSolutions` is empty. */
  function getSelectedSolution(visibleSolutions: Set<HTMLElement>, solSelector: HTMLElement): HTMLElement | null {
    let minDistance = Infinity;
    let selectedSol: HTMLElement | null = null;

    for (const sol of visibleSolutions) {
      const { top: solTop } = sol.getBoundingClientRect();
      const { top: selectorTop } = solSelector.getBoundingClientRect();
      const dist = Math.abs(solTop - selectorTop);

      if (dist < minDistance) {
        selectedSol = sol;
        minDistance = dist;
      }
    }

    return selectedSol;
  }
</script>

<div class={["relative", cls]}>
  <div
    bind:this={solSelectorElem}
    class="md:hidden pointer-fine:hidden -z-1 h-9 left-4 right-4 top-14 bg-rose-100 absolute"
  ></div>

  <div bind:this={solScrollerElem} onscroll={handleScroll} class="overflow-y-scroll">
    {#each solutions.entries() as [date, sols]}
      <div class="pt-6 last:pb-[calc(100%-3.5rem)]">
        <span class="font-bold px-8 md:px-10 block mb-2">
          {new Date(date).toLocaleDateString("en-US", { year: "numeric", month: "long", day: "numeric" })}
        </span>
        <ul class="flex flex-col px-4 md:px-6">
          {#each sols as words}
            <li class="not-first:-mt-1.5">
              <button
                use:refSet={visibleSolElems}
                use:refMap={{ map: solElemData, value: { sol: words, date } }}
                onmouseenter={(e) => {
                  if (!hasFinePointer) return;
                  hoveredSolElem = e.target as HTMLElement;
                  selected = solElemData.get(hoveredSolElem)!;
                }}
                onmouseleave={() => {
                  if (!hasFinePointer) return;
                  hoveredSolElem = undefined;
                  selected = solElemData.get(selectedSolElem!)!;
                }}
                onclick={(e) => {
                  selectedSolElem = e.target as HTMLElement;
                  selected = solElemData.get(selectedSolElem)!;
                  if (!hasFinePointer) snapSolScroller(solScrollerElem, solSelectorElem, selectedSolElem);
                }}
                class={[
                  "block w-full px-4 py-1.5 tracking-wider text-left",
                  hasFinePointer &&
                  selectedSolElem !== undefined &&
                  words.every((w, i) => w === solElemData.get(selectedSolElem!)!.sol[i])
                    ? "bg-rose-100 relative z-1"
                    : "pointer-fine:hover:bg-rose-50",
                ]}
              >
                {words.join(" – ")}
              </button>
            </li>
          {/each}
        </ul>
      </div>
    {/each}
  </div>
</div>
