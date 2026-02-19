import {
  BadgeCheck,
  ChevronsUpDown,
  LoaderCircle,
  LogOut,
  MessageSquareQuote,
} from "lucide-react";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarInput,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "../ui/sidebar";
import { NAV_TITLE_CHAT, Navigations } from "@/navigation/navigation";
import { ReactNode, useEffect, useState } from "react";
import { useDebounce } from "use-debounce";
import { HttpMethod, Mutation, useQueryData } from "@/utils/tanstack";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu";
import { AlertDialogWithMedia } from "../ui/alert-dialog-media";
import { RowConversationChat } from "@/dto/conversation-dto.ts";

type AppSidebarProps = {
  contentSidebarDetail?: ReactNode;
  onSearch: (v: string) => void;
};

export const AppSidebar = ({
  contentSidebarDetail,
  onSearch,
}: AppSidebarProps) => {
  const [activeItem, setActiveItem] = useState<string>(NAV_TITLE_CHAT);
  const [search, setSearch] = useState<string>("");
  const [searchDebounce] = useDebounce(search, 600);

  useEffect(() => {
    onSearch(searchDebounce);
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
        content={contentSidebarDetail}
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
      <SidebarFooter>
        <NavUser />
      </SidebarFooter>
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

type RenderRowsConversationChatProps = {
  conversations: RowConversationChat[];
  onClick: (r: RowConversationChat) => void;
};

export const RenderRowsConversationChat = ({
  conversations,
  onClick,
}: RenderRowsConversationChatProps) => {
  return (
    conversations &&
    conversations.map((g) => (
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

export const NavUser = () => {
  const { isMobile } = useSidebar();
  const [user, setUser] = useState<GetMeResponse | null>(null);

  const { isPending, isSuccess, data } = useQueryData(["getMe"], "/me");
  const {
    mutate,
    isPending: pendingLogout,
    isSuccess: successLogout,
  } = Mutation(["logout"]);

  useEffect(() => {
    if (isSuccess) {
      setUser(data.data as GetMeResponse);
    }
  }, [isSuccess]);

  useEffect(() => {
    if (successLogout) {
      window.location.reload();
    }
  }, [successLogout]);

  const handleLogout = () => {
    mutate({
      url: "/me/logout",
      body: null,
      method: HttpMethod.GET,
    });
  };

  return isPending ? (
    <LoaderCircle className="animate-spin" />
  ) : (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size="lg"
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground md:h-8 md:p-0"
            >
              <Avatar className="h-8 w-8 rounded-lg">
                <AvatarImage
                  src={user?.image}
                  alt={user?.name}
                  className="bg-gray-profile"
                />
                <AvatarFallback className="rounded-lg">CN</AvatarFallback>
              </Avatar>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-medium">{user?.name}</span>
                <span className="truncate text-xs">{user?.email}</span>
              </div>
              <ChevronsUpDown className="ml-auto size-4" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className="w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"
            side={isMobile ? "bottom" : "right"}
            align="end"
            sideOffset={4}
          >
            <DropdownMenuLabel className="p-0 font-normal">
              <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                <Avatar className="h-8 w-8 rounded-lg">
                  <AvatarImage
                    src={user?.image}
                    alt={user?.name}
                    className="bg-gray-profile"
                  />
                  <AvatarFallback className="rounded-lg">CN</AvatarFallback>
                </Avatar>
                <div className="grid flex-1 text-left text-sm leading-tight">
                  <span className="truncate font-medium">{user?.name}</span>
                  <span className="truncate text-xs">{user?.email}</span>
                </div>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              <DropdownMenuItem>
                <BadgeCheck />
                Account
              </DropdownMenuItem>
              <AlertDialogWithMedia
                isPending={pendingLogout}
                onClick={handleLogout}
                Icon={LogOut}
                title="Log out"
                description="Are you really want logout from whatsupp?"
              >
                <DropdownMenuItem onSelect={(e) => e.preventDefault()}>
                  <LogOut />
                  Log out
                </DropdownMenuItem>
              </AlertDialogWithMedia>
            </DropdownMenuGroup>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  );
};
