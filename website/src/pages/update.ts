import type { APIRoute } from "astro";
import { db, solutions, sides } from "../db";
import { z } from "zod";

const apiResSchema = z.object({
  date: z.string(),
  sides: z.string().array(),
  solutions: z.string().array().array(),
});

export const GET: APIRoute = async () => {
  const res = await fetch(`${process.env.API_URL!}/today?max-words=2`);
  const { date, sides: s, solutions: sols } = apiResSchema.parse(await res.json());

  await db
    .insert(sides)
    .values({ date, sides: s.join(",") })
    .onConflictDoNothing();
  await db
    .insert(solutions)
    .values(sols.map((s) => ({ date, words: s.join(",") })))
    .onConflictDoNothing();

  return new Response(null, { status: 204 });
};
