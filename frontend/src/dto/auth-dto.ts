import * as z from "zod";

export const LoginDto = z.object({
  email: z.email(),
  password: z.string().min(3),
});

export const RegisterDto = z.object({
  name: z.string().min(3),
  email: z.email(),
  password: z.string().min(8),
});
