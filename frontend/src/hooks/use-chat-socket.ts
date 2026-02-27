import {
  EVENT_LEAVE_CONVERSATION,
  EVENT_MEMBER_JOIN_CONVERSATION,
  EVENT_MEMBER_LEAVE_CONVERSATION,
  EVENT_NEW_CONVERSATION,
  EVENT_NEW_MESSAGE,
  EventWs,
  LeaveConversationResponse,
  MemberJoinConversationResponse,
  MemberLeaveConversationResponse,
  NewConversationResponse,
  NewMessageResponse,
} from "@/dto/ws-dto";
import { ConnectWS } from "@/utils/ws";
import { useEffect, useRef, useState } from "react";

type useChatSocketProps = {
  onNewMessage: (data: NewMessageResponse) => void;
  onNewConversation: (data: NewConversationResponse) => void;
  onLeaveConversation: (data: LeaveConversationResponse) => void;
  onMemberLeaveConversation: (data: MemberLeaveConversationResponse) => void;
  onMemberJoinConversation: (data: MemberJoinConversationResponse) => void;
};

export function useChatSocket({
  onNewMessage,
  onNewConversation,
  onLeaveConversation,
  onMemberJoinConversation,
  onMemberLeaveConversation,
}: useChatSocketProps) {
  const wsRef = useRef<WebSocket | null>(null);
  const [connected, setConnected] = useState(false);

  const onNewMessageRef = useRef(onNewMessage);
  const onNewConversationRef = useRef(onNewConversation);
  const onLeaveConversationRef = useRef(onLeaveConversation);
  const onMemberLeaveConversationRef = useRef(onMemberLeaveConversation);
  const onMemberJoinConversationRef = useRef(onMemberJoinConversation);

  useEffect(() => {
    onNewMessageRef.current = onNewMessage;
    onNewConversationRef.current = onNewConversation;
    onLeaveConversationRef.current = onLeaveConversation;
    onMemberLeaveConversationRef.current = onMemberLeaveConversation;
    onMemberJoinConversationRef.current = onMemberJoinConversation;
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
          case EVENT_LEAVE_CONVERSATION:
            onLeaveConversationRef.current(
              event.data as LeaveConversationResponse,
            );
            break;
          case EVENT_MEMBER_LEAVE_CONVERSATION:
            onMemberLeaveConversationRef.current(
              event.data as MemberLeaveConversationResponse,
            );
            break;
          case EVENT_MEMBER_JOIN_CONVERSATION:
            onMemberJoinConversationRef.current(
              event.data as MemberJoinConversationResponse,
            );
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
