import { sqliteTable, text, primaryKey } from "drizzle-orm/sqlite-core";

export const sides = sqliteTable("sides", {
  date: text().primaryKey(),
  sides: text().notNull(),
});

export const solutions = sqliteTable(
  "solutions",
  {
    date: text().notNull(),
    words: text().notNull(),
  },
  (table) => [primaryKey({ columns: [table.date, table.words] })]
);

export const dictionary = sqliteTable("dictionary", {
  word: text().primaryKey(),
  definition: text(),
});
