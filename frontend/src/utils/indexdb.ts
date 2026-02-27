import { RowConversationChat } from "@/dto/conversation-dto.ts";
import { ItemGetMessageResponse } from "@/dto/message-dto";
import { MemberEntity } from "@/dto/user-dto";
import { openDB, deleteDB, DBSchema } from "idb";

export interface WhatsuppIdbSchema extends DBSchema {
  conversations: {
    key: number;
    value: RowConversationChat;
  };
  messages: {
    key: number;
    value: ItemGetMessageResponse;
    indexes: {
      conversation_id: number;
      "user.id": number;
    };
  };
  members: {
    key: number;
    value: MemberEntity;
    indexes: {
      conversation_id: number;
      "user.id": number;
    };
  };
}

export const ConnectIdb = async () => {
  const db = await openDB<WhatsuppIdbSchema>("whatsupp", 3, {
    upgrade(database) {
      if (!database.objectStoreNames.contains("conversations")) {
        database.createObjectStore("conversations", {
          keyPath: "conversation_id",
        });
      }

      if (!database.objectStoreNames.contains("messages")) {
        const store = database.createObjectStore("messages", {
          keyPath: "id",
        });

        store.createIndex("conversation_id", "conversation_id");
        store.createIndex("user.id", "user.id");
      }

      if (!database.objectStoreNames.contains("members")) {
        const store = database.createObjectStore("members", {
          keyPath: "id",
        });

        store.createIndex("conversation_id", "conversation_id");
        store.createIndex("user.id", "user.id");
      }
    },
  });

  return db;
};

export const DeleteDbIdb = async () => {
  await deleteDB("whatsupp");
};
