import z from "zod";
import { MemberEntity } from "./user-dto";

export interface RowConversationChat {
  id: number;
  name: string;
  image: string;
  bio: string;
  conversation_type: string;
  conversation_id?: number;
  have_joined: boolean;
  members: MemberEntity[];
}

export const CONVERSATION_TYPE_PRIVATE = "PRIVATE";
export const CONVERSATION_TYPE_GROUP = "GROUP";

export type SearchConversationResponse = RowConversationChat;

export type RecentConversationsResponse = RowConversationChat;

export interface CreateGroupConversationRequest {
  name: string;
  bio: string;
}

export const CreateGroupConversationDto = z.object({
  image: z.file().mime(["image/png", "image/jpeg"]).optional(),
  name: z.string().min(3).max(25),
  bio: z.string().optional(),
});

export const EditGroupConversationDto = CreateGroupConversationDto;
