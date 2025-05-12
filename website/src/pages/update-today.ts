import type { APIRoute } from "astro";
import { db, solutions, solutionWords, words } from "../db";
import { z } from "zod";

export const GET: APIRoute = async () => {
  const res = await fetch(`${process.env.API_URL!}/todays-solutions?max-words=2`);
  const sols = z
    .string()
    .array()
    .array()
    .parse(await res.json());

  const today = new Date();
  today.setHours(7, 0, 0, 0); // Letter Boxed updates daily at 7AM UTC

  // Insert all the solutions into the database
  const solIds = await db
    .insert(solutions)
    .values(Array(sols.length).map(() => ({ date: today.toUTCString() })))
    .returning({ id: solutions.id });
  await db
    .insert(words)
    .values(sols.flatMap((s) => s.map((w) => ({ text: w }))))
    .onConflictDoNothing();
  await db.insert(solutionWords).values(
    sols.flatMap((s, i) =>
      s.map((w, j) => ({
        solutionId: solIds[i].id,
        word: w,
        order: j,
      }))
    )
  );

  return new Response(null, { status: 204 });
};
