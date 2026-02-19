import { ItemGetMessageResponse } from "@/dto/message-dto";
import { useState } from "react";

export function useMessages() {
  const [messagesByChatKey, setMessagesByChatKey] = useState<
    Record<string, ItemGetMessageResponse[]>
  >({});

  const appendMessage = (
    conversationId: number,
    message: ItemGetMessageResponse,
  ) => {
    const key = `conversation-${conversationId}`;

    setMessagesByChatKey((prev) => ({
      ...prev,
      [key]: [...(prev[key] ?? []), message],
    }));
  };

  const setMessages = (
    conversationId: number,
    messages: ItemGetMessageResponse[],
  ) => {
    const key = `conversation-${conversationId}`;

    setMessagesByChatKey((prev) => ({
      ...prev,
      [key]: messages,
    }));
  };

  return {
    messagesByChatKey,
    appendMessage,
    setMessages,
  };
}
