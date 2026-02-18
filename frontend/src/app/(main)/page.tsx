"use client";

import { ChatBannerLoading } from "@/components/chat/banner-loading";
import {
  AppSidebarInset,
  InsetChat,
  InsetHeaderConversationProfile,
} from "@/components/chat/inset";
import { AppSidebar } from "@/components/chat/sidebar";
import { RowConversationChat } from "@/dto/conversation-dto.ts";
import {
  EVENT_SEND_MESSAGE,
  EventWs,
  SEND_MESSAGE,
  SendMessageRequest,
} from "@/dto/ws-dto";
import { ConnectIdb } from "@/utils/indexdb";
import { HttpMethod, Mutation } from "@/utils/tanstack";
import { ToastError } from "@/utils/toast";
import { ConnectWS } from "@/utils/ws";
import { IDBPDatabase } from "idb";
import { useEffect, useRef, useState } from "react";

export default function Page() {
  const wsRef = useRef<WebSocket | null>(null);
  const idbRef = useRef<IDBPDatabase | null>(null);

  const [connect, setConnect] = useState<boolean>(false);
  const [activeChat, setActiveChat] = useState<RowConversationChat>();
  const [messagesByChatKey, setMessagesByChatKey] = useState<
    Record<string, ItemGetMessageResponse[]>
  >({});

  const {
    mutate: mutateGetMessages,
    isSuccess: isSuccessGetMessages,
    data: dataGetMessages,
  } = Mutation(["getMessages"]);

  const handleSendMessage = (v: SendMessageRequest) => {
    const req: EventWs = {
      event: EVENT_SEND_MESSAGE,
      data: v,
    };

    console.log("request send message:", req);

    wsRef.current?.send(JSON.stringify(req));
  };

  const onClickGroupChat = (v: RowConversationChat) => {
    setActiveChat(v);
  };

  useEffect(() => {
    // if (isSuccessGetMessages && dataGetMessages?.data) {
    //   const res = dataGetMessages.data as GetMessageResponse;
    //   const conversationId = res.conversation_id;
    //
    //   if (!conversationId) return;
    //
    //   const chatKey = `group-${groupId}`;
    //   const messages = res.messages as ItemGetMessageResponse[];
    //
    //   setMessagesByChatKey((prev) => ({
    //     ...prev,
    //     [chatKey]: messages,
    //   }));
    // }
  }, [isSuccessGetMessages]);

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
        console.log("onMessage: ", v.data);
        // const data = JSON.parse(v.data) as GetMessageResponse;
        // const newMessages = data.messages;
        // const chatKey = `group-${data.group_id}`;
        //
        // setMessagesByChatKey((prev) => ({
        //   ...prev,
        //   [chatKey]: [...(prev[chatKey] ?? []), ...newMessages],
        // }));
      },
      onOpen: () => {
        setConnect(true);
      },
    });

    return () => {
      wsRef.current!.close();
    };
  }, []);

  return !connect ? (
    <ChatBannerLoading />
  ) : (
    <>
      <AppSidebar onClickGroupChat={onClickGroupChat} />
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
              messages={[]}
              onSubmit={handleSendMessage}
              conversationDetail={activeChat}
            />
          )
        }
      />
    </>
  );
}
