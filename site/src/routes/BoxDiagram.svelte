<script lang="ts">
  interface Props {
    sides: string[];
    letterSeq?: string;
    class?: string;
  }

  let { sides, letterSeq = '', class: cls }: Props = $props();
  let points = $derived.by(() => {
    type PointData = {
      cx: number;
      cy: number;
      x: number;
      y: number;
      dominantBaseline: string;
      textAnchor: string;
    };
    const map = new Map<string, PointData>();

    for (let i = 0; i < sides.length; i++) {
      for (let j = 0; j < sides[i].length; j++) {
        const dot = {
          cx: [(j - 1) * 60, 90, (j - 1) * 60, -90][i],
          cy: [-90, (j - 1) * 60, 90, (j - 1) * 60][i]
        };
        map.set(sides[i][j], {
          ...dot,
          x: dot.cx + [0, 20, 0, -20][i],
          y: dot.cy + [-20, 0, 20, 0][i],
          dominantBaseline: ['auto', 'middle', 'hanging', 'middle'][i],
          textAnchor: ['middle', 'start', 'middle', 'end'][i]
        });
      }
    }

    return map;
  });
</script>

<svg
  viewBox="-150 -150 300 300"
  fill="none"
  xmlns="http://www.w3.org/2000/svg"
  class={['w-75 h-75 lg:w-100 lg:h-100', cls]}
>
  <rect x="-90" y="-90" width="180" height="180" fill="white" stroke-width="2" />
  <polyline
    points={letterSeq
      .split('')
      .map((letter) => {
        const { cx, cy } = points.get(letter)!;
        return `${cx},${cy}`;
      })
      .join(' ')}
    class="stroke-4 fill-none stroke-rose-200"
  />
  <rect x="-90" y="-90" width="180" height="180" stroke="black" stroke-width="2" />

  {#each points.entries() as [letter, { cx, cy, x, y, dominantBaseline, textAnchor }]}
    <circle {cx} {cy} r="6" fill="white" stroke="black" stroke-width="2" />
    <text
      {x}
      {y}
      text-anchor={textAnchor}
      dominant-baseline={dominantBaseline}
      class="fill-black text-2xl font-semibold"
    >
      {letter}
    </text>
  {/each}
</svg>
