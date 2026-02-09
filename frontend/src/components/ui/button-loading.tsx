import React from "react";
import { Button } from "./button";
import { LoaderCircle } from "lucide-react";

export const ButtonLoading = ({
  children,
  isPending,
  ...props
}: React.ComponentProps<"button"> & {
  children: React.ReactNode;
  isPending: boolean;
}) => {
  return (
    <Button {...props} disabled={isPending}>
      {isPending && <LoaderCircle className="animate-spin" />}
      {children}
    </Button>
  );
};
