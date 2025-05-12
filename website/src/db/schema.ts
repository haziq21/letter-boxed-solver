import { sqliteTable, int, text, primaryKey } from "drizzle-orm/sqlite-core";

export const solutions = sqliteTable("solutions", {
  id: int().primaryKey({ autoIncrement: true }),
  date: text().notNull(),
});

export const solutionWords = sqliteTable(
  "solution_words",
  {
    solutionId: int()
      .references(() => solutions.id)
      .notNull(),
    word: text()
      .references(() => words.text)
      .notNull(),
    order: int().notNull(),
  },
  (table) => [primaryKey({ columns: [table.solutionId, table.word] })]
);

export const words = sqliteTable("words", {
  text: text().primaryKey(),
  definition: text(),
});
