import { RegisterRequestSchema } from "@/app/dto";
import { db } from "@/db/db";
import { usersTable } from "@/db/schema/users";
import { eq } from "drizzle-orm";
import { NextRequest, NextResponse } from "next/server";
import { ResponseSchema } from "@/app/dto";
import { v4 as uuidv4 } from "uuid";
import bcrypt from "bcrypt";
import { throws } from "assert";

// register handler
export async function POST(req: NextRequest) {
  try {
    const reqJson = await req.json();
    const parsed = RegisterRequestSchema.parse(reqJson);

    const user = await db
      .select()
      .from(usersTable)
      .where(eq(usersTable.email, parsed.email));

    if (user.length != 0) {
      const res: ResponseSchema = {
        message: "email was used",
      };

      return NextResponse.json(res, {
        status: 400,
      });
    }

    const hash = bcrypt.hashSync(parsed.password, 10);

    const newUser: typeof usersTable.$inferInsert = {
      id: uuidv4(),
      email: parsed.email,
      image: "default",
      name: parsed.name,
      password: hash,
    };

    await db.insert(usersTable).values(newUser);

    const res: ResponseSchema = {
      message: "success register",
    };

    return NextResponse.json(res, {
      status: 200,
    });
  } catch (error) {
    const res: ResponseSchema = {
      message: "failed register",
      error: error,
    };

    return NextResponse.json(res, {
      status: 500,
    });
  }
}
