import { date, pgTable, text } from 'drizzle-orm/pg-core';

export const puzzles = pgTable('puzzles', {
  date: date('date', { mode: 'date' }).primaryKey(),
  sides: text('sides').array().notNull(),
  solutions: text('solutions').array().array().notNull()
});

export const dictionary = pgTable('dictionary', {
  word: text('word').primaryKey(),
  definition: text('definition').notNull()
});
