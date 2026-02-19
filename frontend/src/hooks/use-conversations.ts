import { RowConversationChat } from "@/dto/conversation-dto.ts";
import { useState } from "react";

export function useConversations() {
  const [conversations, setConversations] = useState<RowConversationChat[]>([]);
  const [searchedConversations, setSearchedConversations] = useState<
    RowConversationChat[]
  >([]);
  const [activeChat, setActiveChat] = useState<RowConversationChat>();

  const addConversation = (conv: RowConversationChat) => {
    setConversations((prev) => [conv, ...prev]);
  };

  const addConversations = (convs: RowConversationChat[]) => {
    setConversations((prev) => [...prev, ...convs]);
  };

  return {
    conversations,
    addConversation,
    addConversations,
    activeChat,
    setActiveChat,
    searchedConversations,
    setSearchedConversations,
  };
}
