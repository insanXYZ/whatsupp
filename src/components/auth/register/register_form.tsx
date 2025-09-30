"use client";

import { CardContent } from "@/components/ui/card";
import { AuthButtonSubmit, AuthLabelInput } from "../card_content";
import { FormEvent } from "react";

export default function RegisterForm() {
  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    console.log("test 1");
  };

  return (
    <CardContent>
      <form onSubmit={handleSubmit}>
        <div className="grid gap-6">
          <div className="grid gap-6">
            <AuthLabelInput
              label="Email"
              placeholder="john.doe@example.com"
              type="email"
            />
            <AuthLabelInput label="Password" type="password" />
            <AuthButtonSubmit>Login</AuthButtonSubmit>
          </div>
        </div>
      </form>
    </CardContent>
  );
}
