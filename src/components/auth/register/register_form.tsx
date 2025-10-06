"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  AuthButtonSubmit,
  AuthFormContent,
  AuthFormField,
} from "../card_content";
import { RegisterRequest, RegisterRequestSchema } from "@/app/dto";
import { API } from "@/utils/axios";

export default function RegisterForm() {
  const form = useForm<RegisterRequest>({
    resolver: zodResolver(RegisterRequestSchema),
    defaultValues: {
      email: "",
      name: "",
      password: "",
    },
  });

  const onSubmit = async (v: RegisterRequest) => {
    const data: RegisterRequest = {
      name: v.name,
      email: v.email,
      password: v.password,
    };

    try {
      const res = await API.post("/auth/register", data, {
        headers: {
          "Content-Type": "application/json",
        },
      });
      console.log(res);
    } catch (error) {
      console.log("error ", error);
    }
  };

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
      <AuthButtonSubmit>Sign Up</AuthButtonSubmit>
    </AuthFormContent>
  );
}
