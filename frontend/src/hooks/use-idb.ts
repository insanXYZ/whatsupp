import { RowConversationChat } from "@/dto/conversation-dto.ts";
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

  useEffect(() => {
    ConnectIdb()
      .then((idb) => (idbRef.current = idb))
      .catch(() => {
        // setConnect(false);
        ToastError(
          "Error connected IndexedDB",
          "Please refresh this page, or if thats not help, you can send issues to https://github.com/insanXYZ/whatsupp/issues",
        );
      });
  }, []);

  return {
    AppendConversationsIdb,
    AppendConversationIdb,
    SearchConversationsByNameIdb,
    GetAllConversationsIdb,
    ReplaceConversationsIdb,
  };
};
