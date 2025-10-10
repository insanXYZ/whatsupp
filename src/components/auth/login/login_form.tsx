"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { AuthFormContent, AuthFormField } from "../card_content";
import { LoginRequest, LoginRequestSchema } from "@/app/dto";
import { HttpMethod, Mutation } from "@/utils/tanstack";
import { ButtonLoading } from "@/components/button_loading";
import { useEffect } from "react";

export default function LoginForm() {
  const form = useForm<LoginRequest>({
    resolver: zodResolver(LoginRequestSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  const { mutate, isPending, isSuccess, error, data } = Mutation(["login"]);

  const onSubmit = (v: LoginRequest) => {
    mutate({
      method: HttpMethod.POST,
      body: v,
      url: "/auth/login",
    });
  };

  useEffect(() => {
    console.log("isSuccess ", isSuccess);
    console.log("error ", error);
    console.log("data ", data);
  }, [isSuccess, error]);

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
      <ButtonLoading loading={isPending}>Login</ButtonLoading>
    </AuthFormContent>
  );
}
