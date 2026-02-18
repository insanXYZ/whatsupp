type eventWsString = string;
type typeTargetSend = string;

export const EVENT_SEND_MESSAGE: eventWsString = "SEND_MESSAGE";

export const TARGET_USER: typeTargetSend = "USER";
export const TARGET_GROUP: typeTargetSend = "GROUP";

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
