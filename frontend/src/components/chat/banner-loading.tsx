import { LoaderCircle, MessageSquareQuote } from "lucide-react";

export const ChatBannerLoading = () => {
  return (
    <div className="w-full h-screen flex flex-col items-center justify-center bg-primary">
      <MessageSquareQuote className="size-48 text-white" />

      <LoaderCircle className="size-5 text-white animate-spin" />
    </div>
  );
};
