<script lang="ts">
  import GithubIcon from './GithubIcon.svelte';
  import BoxDiagram from './BoxDiagram.svelte';
  import SolutionList from './SolutionList.svelte';
  import type { PageProps } from './$types';

  const { data }: PageProps = $props();
  const { solsByDate, sidesByDate } = data;

  let selectedSol: { sol: string[]; date: string } | undefined = $state();
</script>

<div
  class="flex h-[100vh] flex-col md:grid md:grid-cols-[27rem_1fr] md:grid-rows-[min-content_1fr]"
>
  <header class="h-min px-8 py-6 md:px-10 md:py-8">
    <div class="mb-2 flex items-baseline justify-between">
      <h1 class="font-roboto-slab text-2xl font-black tracking-tight md:text-2xl">
        Letter unBoxed
      </h1>
      <a href="https://github.com/haziq21/letter-unboxed"><GithubIcon /></a>
    </div>
    <p class="text-sm md:text-base">
      A computed archive of every accepted 2-word solution for the New York Times'
      <a href="https://www.nytimes.com/puzzles/letter-boxed" class="underline">Letter Boxed</a>.
    </p>
  </header>

  <div class="min-h-90 flex items-center justify-center bg-rose-300 md:row-span-2 md:h-full">
    <BoxDiagram
      sides={selectedSol ? sidesByDate.get(selectedSol.date)! : Array(4).fill('')}
      letterSeq={selectedSol?.sol.join('')}
    />
  </div>

  <main class="min-h-0 flex-1">
    <SolutionList
      solutions={solsByDate}
      bind:selected={selectedSol}
      class="flex max-h-full flex-col"
    />
  </main>
</div>
