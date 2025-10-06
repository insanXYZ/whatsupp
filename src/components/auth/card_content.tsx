import {
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../ui/card";
import { Input } from "../ui/input";
import { Button } from "../ui/button";
import React from "react";
import Link from "next/link";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "../ui/form";
import { Control, UseFormReturn } from "react-hook-form";

interface PropsCardHeader {
  title: string;
  description: string;
}

interface PropsFormField {
  label: string;
  control: Control<any>;
  name: string;
  placeholder?: string;
  type?: string;
}

interface PropsButtonSubmit {
  children: string;
}

interface PropsFooterForm {
  title: string;
  href: string;
  title_href: string;
}

interface PropsFormContent {
  children: React.ReactNode;
  onSubmit: (v: any) => void;
  form: UseFormReturn<any>;
}

function AuthCardHeader({ title, description }: PropsCardHeader) {
  return (
    <CardHeader className="text-center">
      <CardTitle className="text-xl">{title}</CardTitle>
      <CardDescription>{description}</CardDescription>
    </CardHeader>
  );
}

function AuthFormField(props: PropsFormField) {
  return (
    <FormField
      control={props.control}
      name={props.name}
      render={({ field }) => (
        <FormItem>
          <FormLabel>{props.label}</FormLabel>
          <FormControl>
            <Input
              type={props.type || "text"}
              placeholder={props.placeholder}
              {...field}
            />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  );
}

function AuthButtonSubmit({ children }: PropsButtonSubmit) {
  return (
    <Button type="submit" className="w-full">
      {children}
    </Button>
  );
}

function AuthFooterForm(props: PropsFooterForm) {
  return (
    <div className="text-center text-sm">
      {props.title} <Link href={props.href}>{props.title_href}</Link>
    </div>
  );
}

function AuthFormContent(props: PropsFormContent) {
  return (
    <CardContent>
      <Form {...props.form}>
        <form onSubmit={props.form.handleSubmit(props.onSubmit)} noValidate>
          <div className="grid gap-6">{props.children}</div>
        </form>
      </Form>
    </CardContent>
  );
}

export {
  AuthCardHeader,
  AuthFormField,
  AuthButtonSubmit,
  AuthFooterForm,
  AuthFormContent,
};
