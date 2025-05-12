import "dotenv/config";
import { defineConfig } from "drizzle-kit";

const dbCredentials =
  process.env.ENVIRONMENT === "production"
    ? { url: process.env.DB_URL!, authToken: process.env.DB_AUTH_TOKEN! }
    : { url: `file:${process.env.DB_FILE_NAME!}` };

export default defineConfig({
  out: "./drizzle",
  schema: "./src/db/schema.ts",
  dialect: process.env.ENVIRONMENT === "production" ? "turso" : "sqlite",
  dbCredentials,
});
