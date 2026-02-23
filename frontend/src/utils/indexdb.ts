import { RowConversationChat } from "@/dto/conversation-dto.ts";
import { openDB, deleteDB, DBSchema } from "idb";

export interface WhatsuppIdbSchema extends DBSchema {
  conversations: {
    key: number;
    value: RowConversationChat;
  };
}

export const ConnectIdb = async () => {
  const db = await openDB<WhatsuppIdbSchema>("whatsupp", 1, {
    upgrade(database) {
      if (!database.objectStoreNames.contains("conversations")) {
        database.createObjectStore("conversations", {
          keyPath: "conversation_id",
        });
      }
    },
  });

  return db;
};

export const DeleteDbIdb = async () => {
  await deleteDB("whatsupp");
};
