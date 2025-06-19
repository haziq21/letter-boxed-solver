import { drizzle } from 'drizzle-orm/neon-http';
import { DATABASE_URL } from '$env/static/private';
import { dictionary, puzzles } from './schema';
import { desc, eq, lte, sql } from 'drizzle-orm';

const db = drizzle(DATABASE_URL);

export async function upsertPuzzle(puzzle: {
  date: Date;
  sides: string[];
  solutions: string[][];
  definitions: Record<string, string>;
}) {
  const { date, sides, solutions, definitions } = puzzle;

  await Promise.all([
    db.insert(puzzles).values({ date, sides, solutions }).onConflictDoUpdate({
      target: puzzles.date,
      set: { sides, solutions }
    }),
    db
      .insert(dictionary)
      .values(Object.entries(definitions).map(([word, definition]) => ({ word, definition })))
      .onConflictDoUpdate({
        target: dictionary.word,
        set: { definition: sql.raw(`excluded.${dictionary.definition.name}`) }
      })
  ]);
}

export async function getPuzzles(options?: {
  maxDate?: Date;
}): Promise<{ date: Date; sides: string[]; solutions: string[][] }[]> {
  if (!options?.maxDate) {
    return await db.select().from(puzzles).orderBy(desc(puzzles.date));
  }
  return await db
    .select()
    .from(puzzles)
    .where(lte(puzzles.date, options.maxDate))
    .orderBy(desc(puzzles.date));
}

export async function getAllDefinitions(options?: {
  maxDate?: Date;
}): Promise<Map<string, string>> {
  let defs: { word: string; definition: string }[];
  if (!options?.maxDate) {
    defs = await db.select().from(dictionary);
  } else {
    defs = await db
      .selectDistinct({
        word: sql<string>`word.word`,
        definition: dictionary.definition
      })
      .from(puzzles)
      .crossJoinLateral(sql`unnest(${puzzles.solutions}) as word`)
      .innerJoin(dictionary, eq(sql`word.word`, dictionary.word))
      .where(lte(puzzles.date, options.maxDate));
  }

  return new Map(defs.map(({ word, definition }) => [word, definition]));
}
