"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { AuthFormContent, AuthFormField } from "../card_content";
import { LoginRequest, LoginRequestSchema, LoginResponse } from "@/app/dto";
import { HttpMethod, Mutation } from "@/utils/tanstack";
import { ButtonLoading } from "@/components/button_loading";
import { useEffect } from "react";
import { ToastError } from "@/components/toast";
import { useRouter } from "next/navigation";

export default function LoginForm() {
  const router = useRouter();

  const form = useForm<LoginRequest>({
    resolver: zodResolver(LoginRequestSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  const { mutate, isPending, isSuccess, isError, error, data } =
    Mutation<LoginResponse>(["login"]);

  const onSubmit = (v: LoginRequest) => {
    mutate({
      method: HttpMethod.POST,
      body: v,
      url: "/auth/login",
    });
  };

  useEffect(() => {
    if (isSuccess && !isError) {
      const accToken = data.data?.access_token!;
      localStorage.setItem("X-ACC-TOKEN", accToken);

      router.push("/");
    } else if (!isSuccess && isError) {
      ToastError(error.response?.data.message as string);
    }
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
