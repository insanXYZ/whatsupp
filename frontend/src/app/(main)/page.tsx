"use client";

import { AppSidebarContent } from "@/components/chat/sidebar-content";
import { AppSidebarInset } from "@/components/chat/sidebar-inset";
import { AppSidebarNavigation } from "@/components/chat/sidebar-navigation";
import { Sidebar, useSidebar } from "@/components/ui/sidebar";
import { NAV_TITLE_CHAT } from "@/navigation/navigation";
import { HttpMethod, Mutation } from "@/utils/tanstack";
import Image from "next/image";
import { ReactNode, useEffect, useState } from "react";
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

  const { mutate, isPending, isSuccess, data } = Mutation(["getGroups"]);

  useEffect(() => {
    if (activeItem == NAV_TITLE_CHAT) {
      mutate({
        url: "/groups?name=" + search,
        body: null,
        method: HttpMethod.GET,
      });
    }
  }, [searchDebounce]);

  useEffect(() => {
    if (isSuccess) {
      switch (activeItem) {
        case NAV_TITLE_CHAT:
          const groups: SearchGroupResponse[] =
            data.data as SearchGroupResponse[];

          setSidebarContent(() =>
            groups.map((g) => (
              <div
                key={g.id}
                className="hover:bg-sidebar-accent hover:text-sidebar-accent-foreground
             flex items-center gap-3 border-b p-4 text-sm leading-tight
             last:border-b-0 min-w-0"
              >
                <Image
                  src={g.image}
                  alt="anu"
                  width={35}
                  height={35}
                  className="shrink-0"
                />

                <div className="flex min-w-0 flex-col gap-1">
                  <span className="font-medium truncate">{g.name}</span>

                  <span className="line-clamp-2 text-xs wrap-break-word">
                    {g.bio ?? "~"}
                  </span>
                </div>
              </div>
            )),
          );
      }
    }
  }, [isSuccess]);

  // const ws = ConnectWS();
  //
  // ws.onopen = (ev) => {
  //   console.log("success open: ", ev.type);
  // };
  //
  // ws.onclose = (ev) => {
  //   console.log("close ws: ", ev.reason);
  // };
  //
  // ws.onerror = (ev) => {
  //   console.log("error ws: ", ev.type);
  // };
  //
  // ws.onmessage = (ev) => {
  //   console.log("incoming message: ", ev.data);
  // };
  //
  return (
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
