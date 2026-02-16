"use client";

import { ChatBannerLoading } from "@/components/chat/banner-loading";
import {
  AppSidebarInset,
  InsetChat,
  InsetHeaderGroupProfile,
} from "@/components/chat/inset";
import { AppSidebar } from "@/components/chat/sidebar";
import { RowGroupChat } from "@/dto/group-dto";
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
  const [activeChat, setActiveChat] = useState<RowGroupChat>();
  const [messagesByChatKey, setMessagesByChatKey] = useState<
    Record<string, ItemGetMessageResponse[]>
  >({});

  const {
    mutate: mutateGetMessages,
    isSuccess: isSuccessGetMessages,
    data: dataGetMessages,
  } = Mutation(["getMessages"]);

  const activeChatKey = activeChat
    ? activeChat.group_id
      ? `group-${activeChat.group_id}`
      : `user-${activeChat.id}`
    : null;

  const handleSendMessage = (v: SendMessageRequest) => {
    wsRef.current?.send(JSON.stringify(v));
  };

  const onClickGroupChat = (v: RowGroupChat) => {
    setActiveChat(v);

    const chatKey = v.group_id ? `group-${v.group_id}` : `user-${v.id}`;

    if (v.group_id && !messagesByChatKey[chatKey]) {
      mutateGetMessages({
        method: HttpMethod.GET,
        url: "/messages/" + v.group_id,
        body: null,
      });
    }
  };

  useEffect(() => {
    if (isSuccessGetMessages && dataGetMessages?.data) {
      const res = dataGetMessages.data as GetMessageResponse;
      const groupId = res.group_id;

      if (!groupId) return;

      const chatKey = `group-${groupId}`;
      const messages = res.messages as ItemGetMessageResponse[];

      setMessagesByChatKey((prev) => ({
        ...prev,
        [chatKey]: messages,
      }));
    }
  }, [isSuccessGetMessages]);

  useEffect(() => {
    ConnectIdb()
      .then((db) => {
        idbRef.current = db;
      })
      .catch(() => {
        ToastError(
          "Error connected IndexedDB",
          "Please refresh this page, or if thats not help, you can send issues to https://github.com/insanXYZ/postgirl/issues",
        );
      });

    wsRef.current = ConnectWS({
      onClose: null,
      onError: null,
      onMessage: (v) => {
        console.log(v.data);
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
            <InsetHeaderGroupProfile
              image={activeChat.image}
              name={activeChat.name}
            />
          )
        }
        content={
          activeChat && (
            <InsetChat
              messages={
                activeChatKey ? (messagesByChatKey[activeChatKey] ?? []) : []
              }
              onSubmit={handleSendMessage}
              groupId={activeChat.group_id}
              receiverId={activeChat.id}
            />
          )
        }
      />
    </>
  );
}
