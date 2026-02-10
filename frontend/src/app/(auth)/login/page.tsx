"use client";

import { Banner } from "@/components/auth/banner";
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
import { LoginDto } from "@/dto/auth-dto";
import { HttpMethod, Mutation } from "@/utils/tanstack";
import { zodResolver } from "@hookform/resolvers/zod";
import { MessageSquareQuote } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { Controller, useForm } from "react-hook-form";

import z from "zod";

export default function LoginPage() {
  const router = useRouter();
  const { mutate, isPending, isSuccess } = Mutation(["login"], true);

  const defaultValues: z.infer<typeof LoginDto> = {
    email: "",
    password: "",
  };

  const form = useForm<z.infer<typeof LoginDto>>({
    defaultValues,
    resolver: zodResolver(LoginDto),
  });

  const onSubmit = (data: z.infer<typeof LoginDto>) => {
    mutate({
      body: data,
      method: HttpMethod.POST,
      url: "/login",
    });
  };

  useEffect(() => {
    if (isSuccess) {
      router.push("/");
    }
  }, [isSuccess]);

  return (
    <CardContent className="grid p-0 md:grid-cols-2">
      <Banner />
      <form
        noValidate
        onSubmit={form.handleSubmit(onSubmit)}
        className="p-6 md:p-8"
      >
        <FieldGroup>
          <div className="flex flex-col items-center gap-2 text-center">
            <h1 className="text-2xl font-bold">Welcome back</h1>
            <p className="text-muted-foreground text-balance">
              Login to your Acme Inc account
            </p>
          </div>
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
                  placeholder="***"
                  required
                />
                <FieldError errors={[fieldState.error]} />
              </Field>
            )}
          ></Controller>
          <Field>
            <ButtonLoading isPending={isPending}>Login</ButtonLoading>
          </Field>
          <FieldDescription className="text-center">
            Don&apos;t have an account? <Link href="/register">Sign up</Link>
          </FieldDescription>
        </FieldGroup>
      </form>
    </CardContent>
  );
}
