<script lang="ts">
  import GithubIcon from "./GithubIcon.svelte";
  import BoxDiagram from "./BoxDiagram.svelte";
  import SolutionList from "./SolutionList.svelte";

  interface Props {
    allSols: Map<string, string[][]>;
    allSides: Map<string, string[]>;
  }

  let { allSols, allSides }: Props = $props();
  let selectedSol: { sol: string[]; date: string } | undefined = $state();
</script>

<div class="flex flex-col h-[100vh] md:grid md:grid-cols-[27rem_1fr] md:grid-rows-[min-content_1fr]">
  <header class="px-8 md:px-10 py-6 md:py-8 h-min">
    <div class="flex justify-between items-baseline mb-2">
      <h1 class="font-karnak text-2xl md:text-3xl">Letter unBoxed</h1>
      <a href="https://github.com/haziq21/letter-unboxed"><GithubIcon /></a>
    </div>
    <p class="text-sm md:text-base">
      A computed archive of every accepted 2-word solution for the New York Times'
      <a href="https://www.nytimes.com/puzzles/letter-boxed" class="underline">Letter Boxed</a>.
    </p>
  </header>

  <div class="bg-rose-300 min-h-90 md:row-span-2 md:h-full flex items-center justify-center">
    {#if selectedSol}
      <BoxDiagram sides={allSides.get(selectedSol.date)!} />
    {/if}
  </div>

  <main class="flex-1 min-h-0">
    <SolutionList solutions={allSols} bind:selected={selectedSol} class="flex flex-col max-h-full" />
  </main>
</div>
