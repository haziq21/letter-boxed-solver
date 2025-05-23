<script lang="ts">
  interface Props {
    sides: string[];
    letterSeq?: string;
  }

  let { sides, letterSeq = "" }: Props = $props();
  let points = $derived.by(() => {
    type SideDirection = "top" | "left" | "bottom" | "right";
    const map = new Map<string, { x: number; y: number; side: SideDirection }>();

    for (let i = 0; i < sides.length; i++) {
      for (let j = 0; j < sides[i].length; j++) {
        map.set(sides[i][j], {
          x: [(j - 1) * 60, 90, (j - 1) * 60, -90][i],
          y: [-90, (j - 1) * 60, 90, (1 - j) * 60][i],
          side: ["top", "right", "bottom", "left"][i] as SideDirection,
        });
      }
    }

    return map;
  });
</script>

<svg viewBox="-150 -150 300 300" fill="none" xmlns="http://www.w3.org/2000/svg" class="w-75 h-75 lg:w-100 lg:h-100">
  <rect x="-90" y="-90" width="180" height="180" fill="white" stroke="black" stroke-width="2"></rect>
  <polyline
    points={letterSeq
      .split("")
      .map((letter) => {
        const { x, y } = points.get(letter)!;
        return `${x},${y}`;
      })
      .join(" ")}
    class="fill-none stroke-rose-200 stroke-4"
  />
  <rect x="-90" y="-90" width="180" height="180" stroke="black" stroke-width="2"></rect>

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
