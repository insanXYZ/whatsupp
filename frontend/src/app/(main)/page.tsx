"use client";

import { ChatBannerLoading } from "@/components/chat/banner-loading";
import { AppSidebarInset } from "@/components/chat/inset";
import { AppSidebar } from "@/components/chat/sidebar";
import { ConnectIdb } from "@/utils/indexdb";
import { ToastError } from "@/utils/toast";
import { ConnectWS } from "@/utils/ws";
import { IDBPDatabase } from "idb";
import { useEffect, useRef, useState } from "react";

export default function Page() {
  const wsRef = useRef<WebSocket | null>(null);
  const idbRef = useRef<IDBPDatabase | null>(null);
  const [connect, setConnect] = useState<boolean>(false);

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
      onMessage: null,
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
      <AppSidebar />
      <AppSidebarInset />
    </>
  );
}
