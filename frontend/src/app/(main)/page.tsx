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
import {
  CONVERSATION_TYPE_GROUP,
  CONVERSATION_TYPE_PRIVATE,
  RowConversationChat,
} from "@/dto/conversation-dto.ts";
import { GetMessageResponse } from "@/dto/message-dto";
import { EVENT_SEND_MESSAGE, EventWs, SendMessageRequest } from "@/dto/ws-dto";
import { useChatSocket } from "@/hooks/use-chat-socket";
import { useConversations } from "@/hooks/use-conversations";
import { useIdb } from "@/hooks/use-idb";
import { useMessages } from "@/hooks/use-messages";
import {
  NAV_TITLE_CHAT,
  NAV_TITLE_CONTACTS,
  NAV_TITLE_GROUPS,
  NAV_TITLE_SEARCH,
} from "@/navigation/navigation";
import { HttpMethod, Mutation, useQueryData } from "@/utils/tanstack";
import { useEffect, useState } from "react";

export default function Page() {
  const [activeItem, setActiveItem] = useState<string>(NAV_TITLE_CHAT);

  const { setMessages, appendMessage, messagesByChatKey } = useMessages();
  const {
    activeChat,
    setActiveChat,
    conversations,
    addConversation,
    addConversations,
    overwriteConversations,
  } = useConversations();

  const {
    AppendConversationsIdb,
    SearchConversationsByNameIdb,
    GetAllConversationsIdb,
    ReplaceConversationsIdb,
    AppendConversationIdb,
  } = useIdb();

  const { send, connected } = useChatSocket({
    onNewMessage: (data) => {
      const message = data.message;
      appendMessage(data.conversation_id, message);
    },
    onNewConversation: async (data) => {
      try {
        AppendConversationIdb(data);

        if (activeItem === NAV_TITLE_SEARCH && !activeChat?.conversation_id) {
          if (
            activeChat?.id === data.id &&
            activeChat.conversation_type === data.conversation_type
          ) {
            setActiveChat(data);
          }
        }

        switch (activeItem) {
          case NAV_TITLE_CHAT:
            addConversation(data);
            break;
          case NAV_TITLE_GROUPS:
            if (data.conversation_type === CONVERSATION_TYPE_GROUP) {
              addConversation(data);
            }
            break;
        }
      } catch (error) {
        console.log(error);
      }
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

  const onSearch = async (v: string) => {
    try {
      switch (activeItem) {
        case NAV_TITLE_CHAT:
          const conversations = await SearchConversationsByNameIdb(v);
          overwriteConversations(conversations);
          break;
        case NAV_TITLE_SEARCH:
          if (v != "") {
            mutateSearchConversations({
              body: null,
              method: HttpMethod.GET,
              url: "/conversations?name=" + v,
            });
          }
          break;
      }
    } catch (error) {
      console.log(error);
    }
  };

  const onChangeActiveItem = async (v: string) => {
    setActiveItem(v);

    try {
      switch (v) {
        case NAV_TITLE_SEARCH:
          overwriteConversations([]);
          break;
        case NAV_TITLE_CHAT:
          const allConversations = await GetAllConversationsIdb();
          overwriteConversations(allConversations);
          break;
      }
    } catch (error) {
      console.log(error);
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
    if (isSuccessGetRecentConversations) {
      const conversations = dataGetRecentConversations.data
        ? (dataGetRecentConversations.data as RowConversationChat[])
        : [];

      ReplaceConversationsIdb(conversations);

      if (activeItem === NAV_TITLE_CHAT) {
        overwriteConversations(conversations);
      }
    }
  }, [isSuccessGetRecentConversations]);

  useEffect(() => {
    if (isSuccessSearchConversations && activeItem == NAV_TITLE_SEARCH) {
      const data = dataSearchConversations.data as RowConversationChat[];

      overwriteConversations(data);
    }
  }, [isSuccessSearchConversations]);

  return !connected && !isSuccessGetRecentConversations ? (
    <ChatBannerLoading />
  ) : (
    <>
      <AppSidebar
        activeItem={activeItem}
        onChangeActiveItem={onChangeActiveItem}
        onSearch={onSearch}
        contentSidebarDetail={
          <RenderRowsConversationChat
            conversations={conversations}
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
                activeChat.conversation_id
                  ? messagesByChatKey[
                  `conversation-${activeChat.conversation_id}`
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
