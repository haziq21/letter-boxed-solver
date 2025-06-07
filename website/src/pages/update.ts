import type { APIRoute } from "astro";
import { setPuzzle } from "../db";
import { z } from "zod";
import "dotenv/config";

const apiResSchema = z.object({
  date: z.string(),
  sides: z.string().array(),
  solutions: z.string().array().array(),
});

export const GET: APIRoute = async () => {
  const res = await fetch(`${process.env.API_URL!}/today?max-words=2`);
  const puzzle = apiResSchema.parse(await res.json());
  setPuzzle(puzzle.date, { sides: puzzle.sides, solutions: puzzle.solutions });

  return new Response(null, { status: 204 });
};
