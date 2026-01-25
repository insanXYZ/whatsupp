import React from "react";
import { Button } from "./button";
import { Icon } from "@iconify/react";

export const ButtonLoading = ({
  children,
  isPending,
}: {
  children: React.ReactNode;
  isPending: boolean;
}) => {
  return (
    <Button disabled={isPending}>
      {isPending && <Icon icon={"line-md:loading-loop"} />}
      {children}
    </Button>
  );
};
