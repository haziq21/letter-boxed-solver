import { primaryKey, sqliteTable, text } from "drizzle-orm/sqlite-core";

export const solutions = sqliteTable(
  "solutions",
  {
    date: text().notNull(),
    words: text().notNull(),
  },
  (table) => [primaryKey({ columns: [table.date, table.words] })]
);
