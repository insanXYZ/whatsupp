import {
  EVENT_NEW_CONVERSATION,
  EVENT_NEW_MESSAGE,
  EventWs,
} from "@/dto/ws-dto";
import { ConnectWS } from "@/utils/ws";
import { useEffect, useRef, useState } from "react";

type useChatSocketProps = {
  onNewMessage: (v: EventWs) => void;
  onNewConversation: (v: EventWs) => void;
};

export function useChatSocket({
  onNewMessage,
  onNewConversation,
}: useChatSocketProps) {
  const wsRef = useRef<WebSocket | null>(null);
  const [connected, setConnected] = useState<boolean>(false);

  useEffect(() => {
    wsRef.current = ConnectWS({
      onOpen: () => setConnected(true),
      onClose: () => setConnected(false),
      onError: () => setConnected(false),
      onMessage: (v) => {
        const event = JSON.parse(v.data);

        switch (event.event) {
          case EVENT_NEW_CONVERSATION:
            onNewConversation(event.data);
            break;

          case EVENT_NEW_MESSAGE:
            onNewMessage(event.data);
            break;
        }
      },
    });

    return () => wsRef.current?.close();
  }, []);

  const send = (event: EventWs) => {
    wsRef.current?.send(JSON.stringify(event));
  };

  return { connected, send };
}
