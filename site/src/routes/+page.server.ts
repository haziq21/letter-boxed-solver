import { getAllDefinitions, getPuzzles } from '$lib/server/database';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
  const now = new Date();
  const maxDate = new Date(now.getTime() - (now.getUTCHours() < 7 ? 2 : 1) * 24 * 60 * 60 * 1000);

  const [puzzles, definitions] = await Promise.all([
    getPuzzles({ maxDate }),
    getAllDefinitions({ maxDate })
  ]);
  console.log(
    `Loaded ${puzzles.length} puzzles (${definitions.size} defs) in ${Date.now() - +now}ms`
  );

  return { puzzles, definitions };
};
