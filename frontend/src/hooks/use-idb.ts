import {
  CONVERSATION_TYPE_GROUP,
  CONVERSATION_TYPE_PRIVATE,
  RowConversationChat,
} from "@/dto/conversation-dto.ts";
import { ItemGetMessageResponse } from "@/dto/message-dto";
import { ConnectIdb, WhatsuppIdbSchema } from "@/utils/indexdb";
import { ToastError } from "@/utils/toast";
import { IDBPDatabase } from "idb";
import { useEffect, useRef } from "react";

export const useIdb = () => {
  const idbRef = useRef<IDBPDatabase<WhatsuppIdbSchema> | null>(null);

  // const AddConversation = async (c: RowConversationChat) => { };

  const AppendConversationsIdb = async (c: RowConversationChat[]) => {
    const tx = idbRef.current?.transaction("conversations", "readwrite");
    const store = tx?.objectStore("conversations");

    for (const row of c) {
      await store?.put(row);
    }

    await tx?.done;
  };

  const AppendConversationIdb = async (c: RowConversationChat) => {
    const tx = idbRef.current?.transaction("conversations", "readwrite");
    const store = tx?.objectStore("conversations");

    store?.put(c);

    await tx?.done;
  };

  const SearchConversationsByNameIdb = async (name: string) => {
    const tx = idbRef.current?.transaction("conversations");
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
    const tx = idbRef.current?.transaction("conversations");
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

  const SearchPrivateConversationsByNameIdb = async (name: string) => {
    const tx = idbRef.current?.transaction("conversations");
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

  const SearchConversationByIdIdb = async (id: number) => {
    const conversations = await idbRef.current?.get("conversations", id);
    return conversations;
  };

  const GetAllConversationsIdb = async () => {
    const conversations = await idbRef.current?.getAll("conversations");
    return conversations ? conversations : [];
  };
  const ReplaceConversationsIdb = async (convs: RowConversationChat[]) => {
    const tx = idbRef.current?.transaction("conversations", "readwrite");
    const store = tx?.objectStore("conversations");

    await store?.clear();

    await Promise.all(convs.map((conv) => store?.put(conv)));
  };

  const DeleteConversationIdb = async (conversationId: number) => {
    const tx = idbRef.current?.transaction("conversations", "readwrite");
    const store = tx?.objectStore("conversations");

    await store?.delete(conversationId);
  };

  const AppendMessagesIdb = async (messages: ItemGetMessageResponse[]) => {
    const tx = idbRef.current?.transaction("messages", "readwrite");
    const store = tx?.objectStore("messages");

    for (const message of messages) {
      await store?.put(message);
    }
  };

  const AppendMessageIdb = async (message: ItemGetMessageResponse) => {
    const tx = idbRef.current?.transaction("messages", "readwrite");
    const store = tx?.objectStore("messages");

    await store?.put(message);
  };

  const GetMessagesWithConversationId = async (conversationId: number) => {
    const messages = await idbRef.current?.getFromIndex(
      "messages",
      "conversation_id",
      conversationId,
    );
    return messages;
  };

  useEffect(() => {
    ConnectIdb()
      .then((idb) => (idbRef.current = idb))
      .catch(() => {
        ToastError(
          "Error connected IndexedDB",
          "Please refresh this page, or if thats not help, you can send issues to https://github.com/insanXYZ/whatsupp/issues",
        );
      });
  }, []);

  return {
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
    ReplaceConversationsIdb,
    DeleteConversationIdb,
  };
};
