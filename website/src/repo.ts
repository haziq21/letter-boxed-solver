import { db, solutions, solutionWords } from "./db";
import { desc, eq, sql } from "drizzle-orm";

/** Retrieve all solutions from the database. Return a Map of date strings to arrays of solutions. */
export async function getAllSolutions(): Promise<Map<string, string[][]>> {
  const sq = db.select().from(solutionWords).orderBy(solutionWords.solutionId, solutionWords.order).as("sq");
  const sols = await db
    .select({ date: solutions.date, words: sql<string>`group_concat(${sq.word})` })
    .from(solutions)
    .innerJoin(sq, eq(sq.solutionId, solutions.id))
    .groupBy(sq.solutionId)
    .orderBy(desc(solutions.date));

  // Aggregate the solutions by date
  const solsMap = new Map<string, string[][]>();
  for (const sol of sols) {
    const date = new Date(sol.date).toISOString();
    if (!solsMap.has(date)) solsMap.set(date, []);
    solsMap.get(date)!.push(sol.words.split(","));
  }

  return solsMap;
}
