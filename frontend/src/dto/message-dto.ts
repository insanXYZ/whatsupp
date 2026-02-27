import { MessageEntity, UserEntity } from "./user-dto";

export interface ItemGetMessageResponse extends MessageEntity {
  is_me: boolean;
  user: UserEntity;
}

export interface GetMessageResponse {
  conversation_id: number;
  messages: ItemGetMessageResponse[];
}
