import { LoginRequest, LoginRequestSchema } from "@/app/dto";
import { cookies } from "next/headers";
import { NextRequest, NextResponse } from "next/server";
import z from "zod";

export async function POST(req: NextRequest) {
  try {
    const reqJson = await req.json();
    const parsed = LoginRequestSchema.parse(reqJson);
  } catch (error) {}
}
