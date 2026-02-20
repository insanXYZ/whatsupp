import { Contact, Inbox, LucideIcon, Search, Settings } from "lucide-react";

export interface NavMain {
  title: string;
  isActive: boolean;
  icon: LucideIcon;
}

export const NAV_TITLE_CHAT = "Chat",
  NAV_TITLE_SETTING = "Setting",
  NAV_TITLE_SEARCH = "Search",
  NAV_TITLE_CONTACTS = "Contacts";

export const Navigations: NavMain[] = [
  {
    title: NAV_TITLE_CHAT,
    icon: Inbox,
    isActive: true,
  },
  {
    title: NAV_TITLE_CONTACTS,
    icon: Contact,
    isActive: true,
  },
  {
    title: NAV_TITLE_SEARCH,
    icon: Search,
    isActive: true,
  },
  {
    title: NAV_TITLE_SETTING,
    icon: Settings,
    isActive: true,
  },
];
