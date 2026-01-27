"use client";

import { ButtonLoading } from "@/components/ui/button-loading";
import { CardContent } from "@/components/ui/card";
import {
  Field,
  FieldDescription,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { RegisterDto } from "@/dto/auth-dto";
import { HttpMethod, Mutation } from "@/utils/tanstack";
import { zodResolver } from "@hookform/resolvers/zod";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { Controller, useForm } from "react-hook-form";

import z from "zod";

export default function RegisterPage() {
  const router = useRouter();
  const { mutate, isPending, isSuccess } = Mutation(["register"], true);

  const defaultValues: z.infer<typeof RegisterDto> = {
    name: "",
    email: "",
    password: "",
  };

  const form = useForm<z.infer<typeof RegisterDto>>({
    defaultValues,
    resolver: zodResolver(RegisterDto),
  });

  const onSubmit = (data: z.infer<typeof RegisterDto>) => {
    mutate({
      body: data,
      method: HttpMethod.POST,
      url: "/register",
    });
  };

  useEffect(() => {
    if (isSuccess) {
      router.push("/login");
    }
  }, [isSuccess]);

  return (
    <CardContent className="grid p-0 md:grid-cols-2">
      <form
        noValidate
        onSubmit={form.handleSubmit(onSubmit)}
        className="p-6 md:p-8"
      >
        <FieldGroup>
          <div className="flex flex-col items-center gap-2 text-center">
            <h1 className="text-2xl font-bold">Hola !!!!</h1>
            <p className="text-muted-foreground text-balance">
              Create your account
            </p>
          </div>
          <Controller
            control={form.control}
            name="name"
            render={({ field, fieldState }) => (
              <Field>
                <FieldLabel htmlFor="name">Username</FieldLabel>
                <Input
                  {...field}
                  id="name"
                  type="text"
                  placeholder="john doe"
                  required
                />
                <FieldError errors={[fieldState.error]} />
              </Field>
            )}
          ></Controller>
          <Controller
            control={form.control}
            name="email"
            render={({ field, fieldState }) => (
              <Field>
                <FieldLabel htmlFor="email">Email</FieldLabel>
                <Input
                  {...field}
                  id="email"
                  type="email"
                  placeholder="m@example.com"
                  required
                />
                <FieldError errors={[fieldState.error]} />
              </Field>
            )}
          ></Controller>
          <Controller
            control={form.control}
            name="password"
            render={({ field, fieldState }) => (
              <Field>
                <FieldLabel htmlFor="password">Password</FieldLabel>
                <Input
                  {...field}
                  id="password"
                  type="password"
                  placeholder="m@example.com"
                  required
                />
                <FieldError errors={[fieldState.error]} />
              </Field>
            )}
          ></Controller>
          <Field>
            <ButtonLoading isPending={isPending}>Register</ButtonLoading>
          </Field>
          <FieldDescription className="text-center">
            have account? <Link href="/login">Login</Link>
          </FieldDescription>
        </FieldGroup>
      </form>
      <div className="bg-muted relative hidden md:block">
        <img
          src="/next.svg"
          alt="Image"
          className="absolute inset-0 h-full w-full object-cover dark:brightness-[0.2] dark:grayscale"
        />
      </div>
    </CardContent>
  );
}
