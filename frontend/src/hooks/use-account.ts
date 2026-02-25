import { AccountContext } from "@/provider/account-provider";
import { useContext } from "react";

export const useAccount = () => {
  const context = useContext(AccountContext);
  if (!context) {
    throw new Error('useAccount should be in AccountProvider');
  }
  return context;
};