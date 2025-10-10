"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  AuthFormContent,
  AuthFormField,
} from "../card_content";
import { RegisterRequest, RegisterRequestSchema } from "@/app/dto";
import { HttpMethod, Mutation } from "@/utils/tanstack";
import { useEffect } from "react";
import { ButtonLoading } from "@/components/button_loading";

export default function RegisterForm() {
  const form = useForm<RegisterRequest>({
    resolver: zodResolver(RegisterRequestSchema),
    defaultValues: {
      email: "",
      name: "",
      password: "",
    },
  });

  const {mutate , isPending , isSuccess, isError , error} = Mutation(["register"]);

  const onSubmit = (v: RegisterRequest) => {
    mutate({
      url: "/auth/register",
      body: v,
      method: HttpMethod.POST
    })
  };

  useEffect(() => {

    console.log("isSuccess",isSuccess)
    console.log("isError",isError)
    console.log("error",error)

  },[isSuccess, isError , error])

  return (
    <AuthFormContent form={form} onSubmit={onSubmit}>
      <AuthFormField
        control={form.control}
        label="Name"
        type="text"
        name="name"
      />
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
      <ButtonLoading loading={isPending}>Register</ButtonLoading>
    </AuthFormContent>
  );
}
