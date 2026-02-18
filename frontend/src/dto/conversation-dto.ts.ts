export interface RowConversationChat {
  id: number;
  name: string;
  image: string;
  bio: string;
  conversation_type: string;
  conversation_id?: number;
}

export type SearchConversationResponse = RowConversationChat;

export type RecentConversationsResponse = RowConversationChat;
