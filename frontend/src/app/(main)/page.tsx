"use client";

import { ChatBannerLoading } from "@/components/chat/banner-loading";
import {
  AppSidebarInset,
  InsetChat,
  InsetHeaderConversationProfile,
} from "@/components/chat/inset";
import {
  AppSidebar,
  RenderRowsConversationChat,
} from "@/components/chat/sidebar";
import { RowConversationChat } from "@/dto/conversation-dto.ts";
import { GetMessageResponse } from "@/dto/message-dto";
import {
  EVENT_SEND_MESSAGE,
  EventWs,
  NewConversationResponse,
  NewMessageResponse,
  SendMessageRequest,
} from "@/dto/ws-dto";
import { useChatSocket } from "@/hooks/use-chat-socket";
import { useConversations } from "@/hooks/use-conversations";
import { useMessages } from "@/hooks/use-messages";
import { HttpMethod, Mutation, useQueryData } from "@/utils/tanstack";
import { useEffect } from "react";

export default function Page() {
  const { setMessages, appendMessage, messagesByChatKey } = useMessages();
  const {
    activeChat,
    setActiveChat,
    conversations,
    addConversation,
    addConversations,
    searchedConversations,
    setSearchedConversations,
  } = useConversations();

  const { send, connected } = useChatSocket({
    onNewMessage: (event) => {
      const data = event.data as NewMessageResponse;
      const message = data.message;

      appendMessage(data.conversation_id, message);
    },
    onNewConversation: (event) => {
      const newConversation = event.data as NewConversationResponse;
      addConversation(newConversation);
    },
  });

  const {
    mutate: mutateGetMessages,
    isSuccess: isSuccessGetMessages,
    data: dataGetMessages,
  } = Mutation(["getMessages"]);

  const {
    mutate: mutateSearchConversations,
    isSuccess: isSuccessSearchConversations,
    data: dataSearchConversations,
  } = Mutation(["getConversations"]);

  const {
    data: dataGetRecentConversations,
    isSuccess: isSuccessGetRecentConversations,
  } = useQueryData(["getRecentConversations"], "/conversations/recent");

  const handleSendMessage = (v: SendMessageRequest) => {
    const req: EventWs = {
      event: EVENT_SEND_MESSAGE,
      data: v,
    };

    send(req);
  };

  const onClickConversationChat = (v: RowConversationChat) => {
    setActiveChat(v);

    if (v.conversation_id) {
      mutateGetMessages({
        body: null,
        method: HttpMethod.GET,
        url: "/messages/" + v.conversation_id,
      });
    }
  };

  const onSearch = (v: string) => {
    if (v != "") {
      mutateSearchConversations({
        body: null,
        method: HttpMethod.GET,
        url: "/conversations?name=" + v,
      });
    } else {
      setSearchedConversations([]);
    }
  };

  useEffect(() => {
    if (isSuccessGetMessages && dataGetMessages.data) {
      const data = dataGetMessages.data as GetMessageResponse;
      const messages = data.messages;

      setMessages(data.conversation_id, messages);
    }
  }, [isSuccessGetMessages]);

  useEffect(() => {
    if (isSuccessGetRecentConversations && dataGetRecentConversations.data) {
      const conversations =
        dataGetRecentConversations.data as RowConversationChat[];
      addConversations(conversations);
    }
  }, [isSuccessGetRecentConversations]);

  useEffect(() => {
    if (isSuccessSearchConversations) {
      const conversations =
        dataSearchConversations.data as RowConversationChat[];

      conversations
        ? setSearchedConversations(conversations)
        : setSearchedConversations([]);
    }
  }, [isSuccessSearchConversations]);

  return !connected ? (
    <ChatBannerLoading />
  ) : (
    <>
      <AppSidebar
        onSearch={onSearch}
        contentSidebarDetail={
          <RenderRowsConversationChat
            conversations={
              searchedConversations.length == 0
                ? conversations
                : searchedConversations
            }
            onClick={onClickConversationChat}
          />
        }
      />
      <AppSidebarInset
        header={
          activeChat && (
            <InsetHeaderConversationProfile
              image={activeChat.image}
              name={activeChat.name}
            />
          )
        }
        content={
          activeChat && (
            <InsetChat
              messages={
                activeChat
                  ? activeChat.conversation_id
                    ? messagesByChatKey[
                    `conversation-${activeChat.conversation_id}`
                    ]
                    : messagesByChatKey[
                    `tmp-${activeChat.conversation_type}-${activeChat.id}`
                    ]
                  : []
              }
              onSubmit={handleSendMessage}
              conversationDetail={activeChat}
            />
          )
        }
      />
    </>
  );
}
