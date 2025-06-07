import "dotenv/config";
import { Redis } from "@upstash/redis";

const redis = Redis.fromEnv();

export async function setPuzzle(date: string, puzzle: { sides: string[]; solutions: string[][] }): Promise<void> {
  await redis.json.set("puzzles", "$", "{}", { nx: true });
  await redis.json.set("puzzles", `$["${date}"]`, JSON.stringify(puzzle));
}

export async function getAllPuzzles(): Promise<Map<string, { sides: string[]; solutions: string[][] }>> {
  const raw = await redis.json.get<{ [date: string]: { sides: string[]; solutions: string[][] } }>("puzzles");
  return new Map(Object.entries(raw ?? {}).sort(([a], [b]) => +new Date(b) - +new Date(a)));
}
