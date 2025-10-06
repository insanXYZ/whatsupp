"use client";

import z from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  AuthButtonSubmit,
  AuthFormContent,
  AuthFormField,
} from "../card_content";
import { API } from "@/utils/axios";
import { LoginRequest, LoginRequestSchema } from "@/app/dto";

export default function LoginForm() {
  const form = useForm<LoginRequest>({
    resolver: zodResolver(LoginRequestSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  const onSubmit = async (v: LoginRequest) => {
    const res = await API.post("/auth/login", {
      email: v.email,
      password: v.password,
    });

    console.log(res.data);
  };

  return (
    <AuthFormContent form={form} onSubmit={onSubmit}>
      <AuthFormField
        control={form.control}
        label="Email"
        type="email"
        name="email"
      />
      <AuthFormField
        control={form.control}
        label="Password"
        name="password"
        type="password"
      />
      <AuthButtonSubmit>Login</AuthButtonSubmit>
    </AuthFormContent>
  );
}
