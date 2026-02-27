import {
  CONVERSATION_TYPE_GROUP,
  CONVERSATION_TYPE_PRIVATE,
  RowConversationChat,
} from "@/dto/conversation-dto.ts";
import { ItemGetMessageResponse } from "@/dto/message-dto";
import { MemberEntity } from "@/dto/user-dto";
import { ConnectIdb, WhatsuppIdbSchema } from "@/utils/indexdb";
import { ToastError } from "@/utils/toast";
import { IDBPDatabase } from "idb";
import { useEffect, useRef } from "react";

export const useIdb = () => {
  let dbPromise: Promise<IDBPDatabase<WhatsuppIdbSchema>> | null = null;

  const getIdb = async () => {
    if (!dbPromise) {
      dbPromise = ConnectIdb();
    }

    return dbPromise;
  };

  // conversations
  //
  const AppendConversationsIdb = async (
    conversations: RowConversationChat[],
  ) => {
    const db = await getIdb();
    const tx = db.transaction("conversations", "readwrite");
    const store = tx?.objectStore("conversations");

    await Promise.all(
      conversations.map((conversation) => store!.put(conversation)),
    );

    await tx?.done;
  };

  const SearchConversationsByNameIdb = async (name: string) => {
    const db = await getIdb();
    const tx = db.transaction("conversations");
    const store = tx?.objectStore("conversations");
    let cursor = await store?.openCursor();

    const res: RowConversationChat[] = [];

    while (cursor) {
      const conversation = cursor.value;

      if (conversation.name.toLowerCase().includes(name.toLowerCase())) {
        res.push(conversation);
      }

      cursor = await cursor.continue();
    }

    return res;
  };

  const SearchGroupConversationsByNameIdb = async (name: string) => {
    const db = await getIdb();
    const tx = db.transaction("conversations");
    const store = tx?.objectStore("conversations");
    let cursor = await store?.openCursor();

    const res: RowConversationChat[] = [];

    while (cursor) {
      const conversation = cursor.value;
      if (
        conversation.conversation_type === CONVERSATION_TYPE_GROUP &&
        (name ? conversation.name.toLowerCase().includes(name) : true)
      ) {
        res.push(conversation);
      }
      cursor = await cursor.continue();
    }

    return res;
  };

  const AppendConversationIdb = async (c: RowConversationChat) => {
    const db = await getIdb();
    const tx = db.transaction("conversations", "readwrite");
    const store = tx?.objectStore("conversations");

    store?.put(c);

    await tx?.done;
  };

  const SearchConversationByIdIdb = async (id: number) => {
    const db = await getIdb();
    const conversations = await db.get("conversations", id);
    return conversations;
  };

  const GetAllConversationsIdb = async () => {
    const db = await getIdb();
    const conversations = await db.getAll("conversations");
    return conversations ? conversations : [];
  };

  const DeleteConversationIdb = async (conversationId: number) => {
    const db = await getIdb();
    const tx = db.transaction("conversations", "readwrite");
    const store = tx?.objectStore("conversations");

    await store?.delete(conversationId);
  };

  const SearchPrivateConversationsByNameIdb = async (name: string) => {
    const db = await getIdb();
    const tx = db.transaction("conversations");
    const store = tx?.objectStore("conversations");
    let cursor = await store?.openCursor();

    const res: RowConversationChat[] = [];

    while (cursor) {
      const conversation = cursor.value;
      if (
        conversation.conversation_type === CONVERSATION_TYPE_PRIVATE &&
        (name ? conversation.name.toLowerCase().includes(name) : true)
      ) {
        res.push(conversation);
      }
      cursor = await cursor.continue();
    }

    return res;
  };

  // messages

  const GetMessagesByConversationIdIdb = async (conversationId: number) => {
    const db = await getIdb();
    return await db.getAllFromIndex(
      "messages",
      "conversation_id",
      conversationId,
    );
  };

  const AppendMessagesIdb = async (messages: ItemGetMessageResponse[]) => {
    const db = await getIdb();
    const tx = db.transaction("messages", "readwrite");
    const store = tx?.objectStore("messages");

    for (const message of messages) {
      await store?.put(message);
    }
  };

  const AppendMessageIdb = async (message: ItemGetMessageResponse) => {
    const db = await getIdb();
    const tx = db.transaction("messages", "readwrite");
    const store = tx?.objectStore("messages");

    await store?.put(message);
  };

  const GetMessagesWithConversationId = async (conversationId: number) => {
    const db = await getIdb();
    const messages = await db.getFromIndex(
      "messages",
      "conversation_id",
      conversationId,
    );
    return messages;
  };

  // members

  const AppendMembersIdb = async (members: MemberEntity[]) => {
    const db = await getIdb();
    const tx = db.transaction("members", "readwrite");
    const store = tx?.objectStore("members");

    await Promise.all(members.map((member) => store?.put(member)));

    await tx?.done;
  };

  const GetMembersByConversationIdIdb = async (conversationId: number) => {
    const db = await getIdb();
    return await db.getAllFromIndex(
      "members",
      "conversation_id",
      conversationId,
    );
  };

  return {
    AppendMembersIdb,
    AppendConversationsIdb,
    AppendConversationIdb,
    AppendMessagesIdb,
    AppendMessageIdb,
    SearchConversationsByNameIdb,
    SearchGroupConversationsByNameIdb,
    SearchPrivateConversationsByNameIdb,
    SearchConversationByIdIdb,
    GetAllConversationsIdb,
    GetMessagesWithConversationId,
    DeleteConversationIdb,
    GetMessagesByConversationIdIdb,
    GetMembersByConversationIdIdb,
  };
};
