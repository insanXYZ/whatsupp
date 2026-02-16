export interface RowGroupChat {
  type: string;
  id: number;
  name: string;
  image: string;
  bio?: string;
  group_type?: string;
  group_id?: number;
}

export type SearchGroupResponse = RowGroupChat;

export type RecentGroupsResponse = RowGroupChat;
