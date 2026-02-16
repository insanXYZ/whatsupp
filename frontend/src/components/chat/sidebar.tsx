import { MessageSquareQuote } from "lucide-react";
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarInput,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "../ui/sidebar";
import { NAV_TITLE_CHAT, Navigations } from "@/navigation/navigation";
import { ReactNode, useEffect, useState } from "react";
import { useDebounce } from "use-debounce";
import { HttpMethod, Mutation, useQueryData } from "@/utils/tanstack";
import {
  RecentGroupsResponse,
  RowGroupChat,
  SearchGroupResponse,
} from "@/dto/group-dto";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";

type AppSidebarProps = {
  onClickGroupChat: (v: RowGroupChat) => void;
};

export const AppSidebar = ({ onClickGroupChat }: AppSidebarProps) => {
  const [activeItem, setActiveItem] = useState<string>(NAV_TITLE_CHAT);
  const [search, setSearch] = useState<string>("");
  const [searchDebounce] = useDebounce(search, 600);

  // sidebar detail
  const [content, setContent] = useState<ReactNode>(null);

  const { data: dataGetMessages, isSuccess: successGetMessages } = useQueryData(
    ["getMessages"],
    "/groups/recent",
  );

  const {
    mutate: mutateGetGroups,
    isSuccess: successGetGroups,
    data: dataGetGroups,
  } = Mutation(["getGroups"]);

  useEffect(() => {
    if (successGetMessages && dataGetMessages.data) {
      const recentGroups = dataGetMessages.data as RecentGroupsResponse[];

      setContent(renderRowsGroupChat(recentGroups, onClickGroupChat));
    }
  }, [successGetMessages]);

  useEffect(() => {
    if (successGetGroups && dataGetGroups.data) {
      const groups = dataGetGroups.data as SearchGroupResponse[];

      setContent(renderRowsGroupChat(groups, onClickGroupChat));
    }
  }, [successGetGroups]);

  useEffect(() => {
    if (searchDebounce == "") {
      // handle render recent groups
    } else {
      mutateGetGroups({
        body: null,
        method: HttpMethod.GET,
        url: "/groups?name=" + searchDebounce,
      });
    }
  }, [searchDebounce]);

  return (
    <Sidebar
      collapsible="icon"
      className="overflow-hidden *:data-[sidebar=sidebar]:flex-row"
    >
      <SidebarNavigation
        activeItem={activeItem}
        onActiveItemChange={setActiveItem}
      />
      <SidebarDetail
        content={content}
        title={activeItem}
        onSearch={setSearch}
      />
    </Sidebar>
  );
};

type SidebarNavigationProps = {
  activeItem: string;
  onActiveItemChange: (v: string) => void;
};

const SidebarNavigation = ({
  activeItem,
  onActiveItemChange,
}: SidebarNavigationProps) => {
  return (
    <Sidebar
      collapsible="none"
      className="w-[calc(var(--sidebar-width-icon)+1px)]! border-r"
    >
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild className="md:h-8 md:p-0">
              <div>
                <div className="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg">
                  <MessageSquareQuote className="size-4" />
                </div>
                <div className="grid flex-1 text-left text-sm leading-tight">
                  <span className="truncate font-medium">Whatsupp</span>
                  <a
                    href="https://github.com/insanXYZ/whatsupp"
                    className="truncate text-xs"
                  >
                    insanXYZ
                  </a>
                </div>
              </div>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupContent className="px-1.5 md:px-0">
            <SidebarMenu>
              {Navigations.map((nav) => {
                return (
                  <SidebarMenuItem key={nav.title}>
                    <SidebarMenuButton
                      tooltip={{
                        children: nav.title,
                        hidden: false,
                      }}
                      onClick={() => {
                        onActiveItemChange(nav.title);
                      }}
                      isActive={nav.title === activeItem}
                      className="px-2.5 md:px-2"
                    >
                      <nav.icon />
                      <span>{nav.title}</span>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                );
              })}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      {/* <SidebarFooter> */}
      {/*   <NavUser /> */}
      {/* </SidebarFooter> */}
    </Sidebar>
  );
};

type SidebarDetailProps = {
  title: string;
  content?: ReactNode;
  onSearch: (v: string) => void;
};

const SidebarDetail = ({ title, content, onSearch }: SidebarDetailProps) => {
  return (
    <Sidebar collapsible="none" className="hidden flex-1 md:flex">
      <SidebarHeader className="gap-3.5 border-b p-4">
        <div className="flex w-full items-center justify-between">
          <div className="text-foreground text-base font-medium">{title}</div>
        </div>
        <SidebarInput
          onChange={(v) => onSearch(v.target.value)}
          placeholder="Type to search..."
        />
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup className="px-0">
          <SidebarGroupContent>{content}</SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
    </Sidebar>
  );
};

export const renderRowsGroupChat = (
  groups: RecentGroupsResponse[],
  onClick: (r: RecentGroupsResponse) => any,
) => {
  return (
    groups &&
    groups.map((g) => (
      <div
        onClick={() => onClick(g)}
        key={g.id}
        className="hover:bg-sidebar-accent hover:text-sidebar-accent-foreground
             flex items-center gap-3 border-b p-4 text-sm leading-tight
             last:border-b-0 min-w-0"
      >
        <Avatar className="h-8 w-8 rounded-lg">
          <AvatarImage src={g.image} alt={g.name} className="bg-gray-profile" />
          <AvatarFallback className="rounded-lg">
            {g.name.toUpperCase().slice(0, 2)}
          </AvatarFallback>
        </Avatar>

        <div className="flex min-w-0 flex-col gap-1">
          <span className="font-medium truncate">{g.name}</span>

          <span className="line-clamp-2 text-xs wrap-break-word">
            {g.bio ?? "~"}
          </span>
        </div>
      </div>
    ))
  );
};
