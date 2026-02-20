import {
  Contact,
  Inbox,
  LucideIcon,
  Search,
  Settings,
  Users,
} from "lucide-react";

export interface NavMain {
  title: string;
  isActive: boolean;
  icon: LucideIcon;
}

export const NAV_TITLE_CHAT = "Chats",
  NAV_TITLE_SETTING = "Setting",
  NAV_TITLE_SEARCH = "Search",
  NAV_TITLE_CONTACTS = "Contacts",
  NAV_TITLE_GROUPS = "Groups";

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
    title: NAV_TITLE_GROUPS,
    icon: Users,
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
