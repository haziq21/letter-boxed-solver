import { API_URL } from '$env/static/private';
import type { RequestHandler } from './$types';
import { z } from 'zod';
import { upsertPuzzle } from '$lib/server/database';

const apiResSchema = z.object({
  date: z.date({ coerce: true }),
  sides: z.string().array(),
  solutions: z.string().array().array(),
  definitions: z.record(z.string(), z.string())
});

export const GET: RequestHandler = async () => {
  const res = await fetch(`${API_URL!}/today?max-words=2`);
  const puzzle = apiResSchema.parse(await res.json());
  await upsertPuzzle(puzzle);

  return new Response(null, { status: 204 });
};
