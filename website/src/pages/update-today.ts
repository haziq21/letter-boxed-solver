import type { APIRoute } from "astro";
import { db, solutions } from "../db";
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
  const todayString = today.toUTCString();

  await db
    .insert(solutions)
    .values(
      sols.map((s) => ({
        date: todayString,
        words: s.join(","),
      }))
    )
    .onConflictDoNothing();

  return new Response(null, { status: 204 });
};
