import "dotenv/config";
import { drizzle } from "drizzle-orm/libsql";

export let db: ReturnType<typeof drizzle>;
export * from "./schema";

if (process.env.ENVIRONMENT === "production") {
  db = drizzle({
    connection: {
      url: process.env.DB_URL!,
      authToken: process.env.DB_AUTH_TOKEN!,
    },
  });
} else {
  db = drizzle(`file:${process.env.DB_FILE_NAME!}`);
}
