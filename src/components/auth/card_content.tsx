import { Label } from "@radix-ui/react-label";
import {
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../ui/card";
import { Input } from "../ui/input";
import { Button } from "../ui/button";
import React, { FormEvent, FormEventHandler } from "react";
import Link from "next/link";

interface PropsCardHeader {
  title: string;
  description: string;
}

interface PropsLabelInput {
  label: string;
  placeholder?: string;
  type: string;
}

interface PropsButtonSubmit {
  children: string;
}

interface PropsCardForm {
  children: React.ReactNode;
  onSubmit: () => void;
}

interface PropsFooterForm {
  title: string;
  href: string;
  title_href: string;
}

function AuthCardHeader({ title, description }: PropsCardHeader) {
  return (
    <CardHeader className="text-center">
      <CardTitle className="text-xl">{title}</CardTitle>
      <CardDescription>{description}</CardDescription>
    </CardHeader>
  );
}

function AuthLabelInput(props: PropsLabelInput) {
  return (
    <div className="grid gap-3">
      <Label htmlFor={props.label}>{props.label}</Label>
      <Input
        id={props.label}
        type={props.type}
        placeholder={props.placeholder}
        required
      />
    </div>
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

export { AuthCardHeader, AuthLabelInput, AuthButtonSubmit, AuthFooterForm };
