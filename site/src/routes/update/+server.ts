import { API_URL } from '$env/static/private';
import type { RequestHandler } from './$types';
import { setPuzzle } from '$lib/server/db';
import { z } from 'zod';

const apiResSchema = z.object({
  date: z.string(),
  sides: z.string().array(),
  solutions: z.string().array().array()
});

export const GET: RequestHandler = async () => {
  const res = await fetch(`${API_URL!}/today?max-words=2`);
  const puzzle = apiResSchema.parse(await res.json());
  setPuzzle(puzzle.date, { sides: puzzle.sides, solutions: puzzle.solutions });

  return new Response(null, { status: 204 });
};
