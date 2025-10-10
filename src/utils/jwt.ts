import { JWTPayload, SignJWT, jwtVerify } from "jose";

const secret = new TextEncoder().encode(process.env.JWT_SECRET);
const alg = "HS256";

export async function CreateJWT(
  payload: JWTPayload,
  expTime: string
): Promise<string> {
  const jwt = await new SignJWT(payload)
    .setExpirationTime(expTime)
    .setProtectedHeader({
      alg: alg,
    })
    .sign(secret);
  return jwt;
}

export async function VerifyJWT(jwtToken: string): Promise<JWTPayload> {
  const { payload } = await jwtVerify(jwtToken, secret);
  return payload;
}
