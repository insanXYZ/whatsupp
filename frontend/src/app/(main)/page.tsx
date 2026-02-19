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
import { GetMessageResponse, ItemGetMessageResponse } from "@/dto/message-dto";
import {
  EVENT_NEW_CONVERSATION,
  EVENT_NEW_MESSAGE,
  EVENT_SEND_MESSAGE,
  EventWs,
  NewConversationResponse,
  NewMessageResponse,
  SendMessageRequest,
} from "@/dto/ws-dto";
import { ConnectIdb } from "@/utils/indexdb";
import { HttpMethod, Mutation, useQueryData } from "@/utils/tanstack";
import { ToastError } from "@/utils/toast";
import { ConnectWS } from "@/utils/ws";
import { IDBPDatabase } from "idb";
import { useEffect, useEffectEvent, useRef, useState } from "react";

export default function Page() {
  const wsRef = useRef<WebSocket | null>(null);
  const idbRef = useRef<IDBPDatabase | null>(null);

  const [connect, setConnect] = useState<boolean>(false);
  const [conversationChat, setConversationChat] = useState<
    RowConversationChat[]
  >([]);
  const [conversationChatSearched, setConversationChatSearched] = useState<
    RowConversationChat[]
  >([]);
  const [activeChat, setActiveChat] = useState<RowConversationChat>();
  const [messagesByChatKey, setMessagesByChatKey] = useState<
    Record<string, ItemGetMessageResponse[]>
  >({});

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

    wsRef.current?.send(JSON.stringify(req));
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
      setConversationChatSearched(() => []);
    }
  };

  useEffect(() => {
    if (isSuccessGetMessages && dataGetMessages.data) {
      const data = dataGetMessages.data as GetMessageResponse;
      const messages = data.messages;

      const chatKey = `conversation-${data.conversation_id}`;

      setMessagesByChatKey((prev) => ({
        ...prev,
        [chatKey]: messages,
      }));
    }
  }, [isSuccessGetMessages]);

  useEffect(() => {
    if (isSuccessGetRecentConversations && dataGetRecentConversations.data) {
      const conversations =
        dataGetRecentConversations.data as RowConversationChat[];

      setConversationChat((prev) => [...prev, ...conversations]);
    }
  }, [isSuccessGetRecentConversations]);

  useEffect(() => {
    if (isSuccessSearchConversations) {
      const conversations =
        dataSearchConversations.data as RowConversationChat[];

      conversations
        ? setConversationChatSearched(() => [...conversations])
        : setConversationChatSearched(() => []);
    }
  }, [isSuccessSearchConversations]);

  useEffect(() => {
    ConnectIdb()
      .then((db) => {
        idbRef.current = db;
      })
      .catch(() => {
        ToastError(
          "Error connected IndexedDB",
          "Please refresh this page, or if thats not help, you can send issues to https://github.com/insanXYZ/whatsupp/issues",
        );
      });

    wsRef.current = ConnectWS({
      onClose: (v) => {
        console.log("close ", v);
        setConnect(false);
      },
      onError: (v) => {
        console.log("error ", v);
        setConnect(false);
      },
      onMessage: (v) => {
        const event = JSON.parse(v.data) as EventWs;
        console.log(event);

        switch (event.event) {
          case EVENT_NEW_CONVERSATION:
            const newConversation = event.data as NewConversationResponse;
            setConversationChat((prev) => [newConversation, ...prev]);
          case EVENT_NEW_MESSAGE:
            const data = event.data as NewMessageResponse;
            const message = data.message;

            const chatKey = `conversation-${data.conversation_id}`;

            setMessagesByChatKey((prev) => ({
              ...prev,
              [chatKey]: [...(prev[chatKey] ?? []), message],
            }));
        }
      },
      onOpen: () => {
        setConnect(true);
      },
    });

    return () => {
      wsRef.current!.close();
    };
  }, []);

  useEffect(() => {
    // console.log(JSON.stringify(messagesByChatKey));
  }, [messagesByChatKey]);

  return !connect ? (
    <ChatBannerLoading />
  ) : (
    <>
      <AppSidebar
        onSearch={onSearch}
        contentSidebarDetail={
          <RenderRowsConversationChat
            conversations={
              conversationChatSearched.length == 0
                ? conversationChat
                : conversationChatSearched
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
