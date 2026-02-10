import { MessageSquareQuote } from "lucide-react";

export const Banner = () => {
  return (
    <div className="bg-muted relative hidden md:block">
      <div className="absolute bg-primary inset-0 h-full w-full flex items-center justify-center object-cover dark:brightness-[0.2] dark:grayscale">
        <MessageSquareQuote className="size-44 text-white" />
      </div>
    </div>
  );
};
