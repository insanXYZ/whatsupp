import {
  EVENT_NEW_CONVERSATION,
  EVENT_NEW_MESSAGE,
  EventWs,
  NewConversationResponse,
  NewMessageResponse,
} from "@/dto/ws-dto";
import { ConnectWS } from "@/utils/ws";
import { useEffect, useRef, useState } from "react";

type useChatSocketProps = {
  onNewMessage: (data: NewMessageResponse) => void;
  onNewConversation: (data: NewConversationResponse) => void;
};

export function useChatSocket({
  onNewMessage,
  onNewConversation,
}: useChatSocketProps) {
  const wsRef = useRef<WebSocket | null>(null);
  const [connected, setConnected] = useState(false);

  const onNewMessageRef = useRef(onNewMessage);
  const onNewConversationRef = useRef(onNewConversation);

  useEffect(() => {
    onNewMessageRef.current = onNewMessage;
    onNewConversationRef.current = onNewConversation;
  });

  useEffect(() => {
    wsRef.current = ConnectWS({
      onOpen: () => setConnected(true),
      onClose: () => setConnected(false),
      onError: () => setConnected(false),
      onMessage: (v) => {
        const event = JSON.parse(v.data) as EventWs;

        switch (event.event) {
          case EVENT_NEW_CONVERSATION:
            onNewConversationRef.current(event.data as NewConversationResponse);
            break;

          case EVENT_NEW_MESSAGE:
            onNewMessageRef.current(event.data as NewMessageResponse);
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
