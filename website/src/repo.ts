import { db, sides, solutions } from "./db";
import { desc } from "drizzle-orm";

/** Retrieve all solutions from the database. Return a Map of date strings to arrays of solutions. */
export async function getSolutionsByDate(): Promise<Map<string, string[][]>> {
  const sols = await db.select().from(solutions).orderBy(desc(solutions.date), solutions.words);

  // Aggregate the solutions by date
  const solsMap = new Map<string, string[][]>();
  for (const sol of sols) {
    if (!solsMap.has(sol.date)) solsMap.set(sol.date, []);
    solsMap.get(sol.date)!.push(sol.words.split(","));
  }

  return solsMap;
}

/** Retrieve the sides of the puzzles for every date. Return a Map of date strings to arrays of sides. */
export async function getSidesByDate(): Promise<Map<string, string[]>> {
  const allSides = await db.select().from(sides).orderBy(desc(sides.date));

  // Aggregate the sides by date
  return new Map(allSides.map((s) => [s.date, s.sides.split(",")]));
}
