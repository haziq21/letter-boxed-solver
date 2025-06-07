import { getAllPuzzles } from '$lib/server/db';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
  const startTime = Date.now();
  const allPuzzles = await getAllPuzzles();
  console.log(`Loaded ${allPuzzles.size} puzzles in ${Date.now() - startTime}ms`);

  return {
    solsByDate: new Map(allPuzzles.keys().map((date) => [date, allPuzzles.get(date)!.solutions])),
    sidesByDate: new Map(allPuzzles.keys().map((date) => [date, allPuzzles.get(date)!.sides]))
  };
};
