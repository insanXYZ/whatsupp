import { Icon } from "@iconify/react";
import { Button } from "./ui/button";

interface PropsButtonSubmit {
  children: string;
  loading: boolean;
}

export function ButtonLoading({ children, loading }: PropsButtonSubmit)  {
  return (
    <Button type="submit" className="w-full">
      {loading ? <Icon icon={"line-md:loading-loop"} /> : children}
    </Button>
  );
}
