import { pgTable, text, timestamp, varchar } from "drizzle-orm/pg-core";

export const usersTable = pgTable("users", {
  id: varchar({ length: 255 }).primaryKey(),
  name: varchar({ length: 100 }).notNull(),
  email: varchar({ length: 100 }).notNull(),
  password: varchar({ length: 255 }).notNull(),
  image: text().notNull(),
  created_at: timestamp().notNull().defaultNow(),
  updated_at: timestamp(),
});
