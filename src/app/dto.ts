import z from "zod";

export interface ResponseSchema<T = any> {
  message: string;
  data?: T;
  error?: any;
}

export const LoginRequestSchema = z.object({
  email: z.email(),
  password: z.string().min(6),
});

export type LoginRequest = z.infer<typeof LoginRequestSchema>;

export const RegisterRequestSchema = z.object({
  email: z.email(),
  password: z.string().min(6),
  name: z.string().min(3),
});

export type RegisterRequest = z.infer<typeof RegisterRequestSchema>;
