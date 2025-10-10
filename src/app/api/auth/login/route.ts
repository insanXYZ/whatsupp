import { LoginRequestSchema, ResponseSchema } from "@/app/dto";
import { db } from "@/db/db";
import { usersTable } from "@/db/schema/users";
import { eq } from "drizzle-orm";
import { NextRequest, NextResponse } from "next/server";
import bcrypt from "bcrypt";
import { CreateJWT } from "@/utils/jwt";
import { JWTPayload } from "jose";

export async function POST(req: NextRequest) {
  try {
    const reqJson = await req.json();
    const loginRequest = LoginRequestSchema.parse(reqJson);

    const users = await db
      .select()
      .from(usersTable)
      .where(eq(usersTable.email, loginRequest.email));
    if (users.length == 0) {
      var res: ResponseSchema = {
        message: "email or password incorrect",
      };

      return NextResponse.json(res, {
        status: 400,
      });
    }

    const user = users[0];

    const isCompared = await bcrypt.compare(
      loginRequest.password,
      user.password
    );

    if (!isCompared) {
      var res: ResponseSchema = {
        message: "email or password incorrect",
      };

      return NextResponse.json(res, {
        status: 400,
      });
    }

    const payloadJWT: JWTPayload = {
      sub: user.id,
    };

    const jwtToken = await CreateJWT(payloadJWT, "15m");

    var res: ResponseSchema = {
      message: "login successful",
      data: {
        access_token: jwtToken,
      },
    };

    return NextResponse.json(res, {
      status: 200,
    });
  } catch (error) {
    console.log("error login:", error);

    return NextResponse.json(
      {
        error: error,
      },
      {
        status: 500,
      }
    );
  }
}
