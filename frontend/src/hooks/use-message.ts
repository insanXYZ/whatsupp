import { ItemGetMessageResponse } from "@/dto/message-dto";
import { setRequestMeta } from "next/dist/server/request-meta";
import { useState } from "react";

export const useMessages = () => {
  const [messages, setMessages] = useState<ItemGetMessageResponse[]>([]);

  const AppendMessage = (message: ItemGetMessageResponse) => {
    setMessages((prev) => [...prev, message]);
  };

  const OverwriteMessages = (messages: ItemGetMessageResponse[]) => {
    setMessages(messages);
  };

  return { messages, AppendMessage, OverwriteMessages };
};
