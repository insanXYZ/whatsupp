export interface UserEntity {
  bio: string;
  email: string;
  id: number;
  image: string;
  name: string;
}

export type GetMeResponse = UserEntity;

export interface MemberEntity {
  id: number;
  role: string;
  conversation_id: number;
  user: UserEntity;
}

export interface MessageEntity {
  id: number;
  message: string;
  conversation_id: number;
  created_at: string;
}
