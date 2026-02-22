import { RowConversationChat } from "./conversation-dto.ts";
import { ItemGetMessageResponse } from "./message-dto.js";

type eventWsString = string;
type typeTargetSend = string;

export const EVENT_SEND_MESSAGE: eventWsString = "SEND_MESSAGE";
export const EVENT_NEW_MESSAGE: eventWsString = "NEW_MESSAGE";
export const EVENT_NEW_CONVERSATION: eventWsString = "NEW_CONVERSATION";
export const EVENT_LEAVE_CONVERSATION: eventWsString = "LEAVE_CONVERSATION";

export interface EventWs {
  event: eventWsString;
  data: any;
}

export interface TargetSendMessage {
  type: typeTargetSend;
  id: number;
}

export interface SendMessageRequest {
  target: TargetSendMessage;
  message: string;
  tmp_conversation_id?: string;
  conversation_id?: number;
}

export type NewConversationResponse = RowConversationChat;

export interface NewMessageResponse {
  tmp_conversation_id?: string;
  conversation_id: number;
  message: ItemGetMessageResponse;
}

export interface LeaveConversationResponse {
  conversation_id: number;
}
