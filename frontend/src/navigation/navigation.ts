import { Contact, Inbox, LucideIcon, Search, Settings } from "lucide-react";

export interface NavMain {
  title: string;
  isActive: boolean;
  icon: LucideIcon;
}

export const NAV_TITLE_CHAT = "Chat",
  NAV_TITLE_FRIEND = "Friend",
  NAV_TITLE_Setting = "Setting";

export const Navigations: NavMain[] = [
  {
    title: "Chat",
    icon: Inbox,
    isActive: true,
  },
  {
    title: "Contacts",
    icon: Contact,
    isActive: true,
  },
  {
    title: "Search",
    icon: Search,
    isActive: true,
  },
  {
    title: "Setting",
    icon: Settings,
    isActive: true,
  },
];
