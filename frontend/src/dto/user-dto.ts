interface GetMeResponse {
  id: number;
  name: string;
  image: string;
  email: string;
  bio: string;
}

interface UserEntity {
  bio: string;
  email: string;
  id: number;
  image: string;
  name: string;
}

interface MemberEntity {
  id: number;
  group_id: number;
  user_id: number;
  role: string;
  user: UserEntity;
}

interface MessageEntity {
  id: number;
  member_id: number;
  message: string;
  created_at: string;
  member: MemberEntity;
}

interface GetMessageResponse {
  group_id: number;
  messages: ItemGetMessageResponse[];
}

interface ItemGetMessageResponse extends MessageEntity {
  is_me: boolean;
}
