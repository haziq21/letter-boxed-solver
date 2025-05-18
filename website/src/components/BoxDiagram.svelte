<script lang="ts">
  interface Props {
    sides: string[];
  }

  const { sides }: Props = $props();
</script>

<svg viewBox="-150 -150 300 300" fill="none" xmlns="http://www.w3.org/2000/svg" class="w-75 h-75">
  <rect x="-90" y="-90" width="180" height="180" fill="white" stroke="black" stroke-width="2"></rect>

  {#each [-60, 0, 60] as x, i}
    <text {x} y="-110" text-anchor="middle" class="fill-black text-2xl font-semibold">{sides[0][i]}</text>
    <circle cx={x} cy="-90" r="6" fill="white" stroke="black" stroke-width="2"></circle>
  {/each}
  {#each [-60, 0, 60] as y, i}
    <text x="110" {y} dominant-baseline="middle" class="fill-black text-2xl font-semibold">{sides[1][i]}</text>
    <circle cx="90" cy={y} r="6" fill="white" stroke="black" stroke-width="2"></circle>
  {/each}
  {#each [-60, 0, 60] as x, i}
    <text {x} y="110" text-anchor="middle" dominant-baseline="hanging" class="fill-black text-2xl font-semibold">
      {sides[2][i]}
    </text>
    <circle cx={x} cy="90" r="6" fill="white" stroke="black" stroke-width="2"></circle>
  {/each}
  {#each [-60, 0, 60] as y, i}
    <text x="-110" {y} text-anchor="end" dominant-baseline="middle" class="fill-black text-2xl font-semibold">
      {sides[3][i]}
    </text>
    <circle cx="-90" cy={y} r="6" fill="white" stroke="black" stroke-width="2"></circle>
  {/each}
</svg>
