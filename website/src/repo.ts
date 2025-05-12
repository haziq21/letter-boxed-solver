import { db, solutions } from "./db";
import { desc } from "drizzle-orm";

/** Retrieve all solutions from the database. Return a Map of date strings to arrays of solutions. */
export async function getAllSolutions(): Promise<Map<string, string[][]>> {
  const sols = await db.select().from(solutions).orderBy(desc(solutions.date), solutions.words);

  // Aggregate the solutions by date
  const solsMap = new Map<string, string[][]>();
  for (const sol of sols) {
    const date = new Date(sol.date).toISOString();
    if (!solsMap.has(date)) solsMap.set(date, []);
    solsMap.get(date)!.push(sol.words.split(","));
  }

  return solsMap;
}
