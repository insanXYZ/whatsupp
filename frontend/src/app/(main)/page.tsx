"use client";

import { ChatBannerLoading } from "@/components/chat/banner-loading";
import {
  AppSidebarContent,
  RowsGroupChat,
} from "@/components/chat/sidebar-content";
import {
  AppSidebarInset,
  InsetHeaderGroup,
} from "@/components/chat/sidebar-inset";
import { AppSidebarNavigation } from "@/components/chat/sidebar-navigation";
import { Sidebar, useSidebar } from "@/components/ui/sidebar";
import { GroupNavigationContent, SearchGroupResponse } from "@/dto/group-dto";
import { NAV_TITLE_CHAT } from "@/navigation/navigation";
import { HttpMethod, Mutation } from "@/utils/tanstack";
import { ConnectWS } from "@/utils/ws";
import { ReactNode, useEffect, useRef, useState } from "react";
import { useDebounce } from "use-debounce";

export default function Page() {
  // AppSidebarNavigation
  const [activeItem, setActiveItem] = useState<string>(NAV_TITLE_CHAT);
  const { setOpen } = useSidebar();

  // AppSidebarContent
  const [sidebarContent, setSidebarContent] = useState<ReactNode>();
  const [sidebarContentHeader, setSidebarContentHeader] = useState<ReactNode>();
  const [groupId, setGroupId] = useState<number | undefined>(undefined);
  const [search, setSearch] = useState<string>("");
  const [searchDebounce] = useDebounce(search, 600);

  // AppSidebarInset
  const [sidebarInsetContent, setSidebarInsetContent] = useState<ReactNode>();
  const [sidebarInsetHeader, setSidebarInsetHeader] = useState<ReactNode>();

  const wsRef = useRef<WebSocket | null>(null);
  const [connect, setConnect] = useState<boolean>(false);

  const { mutate, isPending, isSuccess, data } = Mutation(["getGroups"]);

  const handleGroupSelected = (group: GroupNavigationContent) => {
    setGroupId(group.id);

    setSidebarInsetHeader(InsetHeaderGroup(group.image, group.name));
  };

  useEffect(() => {
    if (searchDebounce != "") {
      if (activeItem == NAV_TITLE_CHAT) {
        mutate({
          url: "/groups?name=" + search,
          body: null,
          method: HttpMethod.GET,
        });
      }
    }
  }, [searchDebounce]);

  useEffect(() => {
    if (isSuccess) {
      switch (activeItem) {
        case NAV_TITLE_CHAT:
          const groups: SearchGroupResponse[] =
            data.data as SearchGroupResponse[];
          setSidebarContent(() => RowsGroupChat(groups, handleGroupSelected));
      }
    }
  }, [isSuccess]);

  useEffect(() => {
    wsRef.current = ConnectWS();

    const ws = wsRef.current;

    ws.onopen = (ev) => {
      console.log("success open: ", ev.type);
      setConnect(true);
    };

    ws.onmessage = (ev) => {
      console.log("incoming message: ", ev.data);
    };

    return () => {
      console.log("cleanup");
      ws.close();
    };
  }, []);

  return !connect ? (
    <ChatBannerLoading />
  ) : (
    <>
      <Sidebar
        collapsible="icon"
        className="overflow-hidden *:data-[sidebar=sidebar]:flex-row"
      >
        <AppSidebarNavigation
          setOpen={setOpen}
          activeItem={activeItem}
          setActiveItem={setActiveItem}
        />
        <AppSidebarContent
          onSearch={setSearch}
          title={activeItem}
          content={sidebarContent}
          headerComponent={sidebarContentHeader}
        />
      </Sidebar>

      <AppSidebarInset
        header={sidebarInsetHeader}
        content={sidebarInsetContent}
      />
    </>
  );
}
