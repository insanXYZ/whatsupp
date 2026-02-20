import { RowConversationChat } from "@/dto/conversation-dto.ts";
import { openDB, deleteDB, wrap, unwrap, DBSchema } from "idb";

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
          keyPath: "id",
        });
      }

      // if (!database.objectStoreNames.contains("messages")) {
      //   const store = database.createObjectStore("messages", {
      //     keyPath: "id",
      //   });
      //
      //   store.createIndex("conversation_id", "conversation_id");
      //
      //   store.createIndex("conversation_created", [
      //     "conversation_id",
      //     "created_at",
      //   ]);
      // }
    },
  });

  return db;
};
