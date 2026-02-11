"use client";

import {
  BadgeCheck,
  Bell,
  ChevronsUpDown,
  CreditCard,
  LoaderCircle,
  LogOut,
  Sparkles,
} from "lucide-react";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar";
import { useEffect, useState } from "react";
import { HttpMethod, Mutation, useQueryData } from "@/utils/tanstack";
import { AlertDialogWithMedia } from "../ui/alert-dialog-media";

export function NavUser() {
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
}
