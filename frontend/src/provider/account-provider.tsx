"use client";

import { GetMeResponse, UserEntity } from "@/dto/user-dto";
import { HttpMethod, Mutation } from "@/utils/tanstack";
import { createContext, ReactNode, useEffect, useState } from "react";

type AccountContextType = {
  user: UserEntity | null;
  setUser: React.Dispatch<React.SetStateAction<UserEntity | null>>;
};

export const AccountContext = createContext<AccountContextType | undefined>(
  undefined,
);

export const AccountProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<UserEntity | null>(null);

  const { mutate, isSuccess, data } = Mutation(["getMe"]);

  useEffect(() => {
    mutate({
      body: null,
      method: HttpMethod.GET,
      url: "/me",
    });
  }, []);

  useEffect(() => {
    if (isSuccess) {
      setUser(data.data as GetMeResponse);
    }
  }, [isSuccess]);

  return (
    <AccountContext.Provider value={{ user, setUser }}>
      {children}
    </AccountContext.Provider>
  );
};
