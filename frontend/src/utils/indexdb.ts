import { openDB, deleteDB, wrap, unwrap } from "idb";

export const ConnectIdb = async () => {
  const db = await openDB("whatsupp", 1, {
    upgrade(database, oldVersion, newVersion, transaction, event) {
      if (database.objectStoreNames.contains("groups")) {
        database.deleteObjectStore("groups");
      }

      if (database.objectStoreNames.contains("messages")) {
        database.deleteObjectStore("messages");
      }

      database.createObjectStore("groups");
      database.createObjectStore("messages");
    },
  });

  return db;
};
