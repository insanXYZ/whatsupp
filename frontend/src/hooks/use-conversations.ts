import { RowConversationChat } from "@/dto/conversation-dto.ts";
import { useEffect, useState } from "react";

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

  const overwriteConversations = (convs: RowConversationChat[]) => {
    setConversations(convs);
  };

  return {
    conversations,
    addConversation,
    addConversations,
    overwriteConversations,
    activeChat,
    setActiveChat,
    searchedConversations,
    setSearchedConversations,
  };
}
