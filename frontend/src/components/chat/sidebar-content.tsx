import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarInput,
} from "../ui/sidebar";
import { GroupNavigationContent } from "@/dto/group-dto";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";

export const AppSidebarContent = ({
  title,
  headerComponent,
  content,
  onSearch,
}: {
  title: string;
  headerComponent?: React.ReactNode;
  content?: React.ReactNode;
  onSearch: (v: string) => void;
}) => {
  return (
    <Sidebar collapsible="none" className="hidden flex-1 md:flex">
      <SidebarHeader className="gap-3.5 border-b p-4">
        <div className="flex w-full items-center justify-between">
          <div className="text-foreground text-base font-medium">{title}</div>
          {headerComponent ?? headerComponent}
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

export const RowsGroupChat = (
  contents: GroupNavigationContent[],
  onClick: (g: GroupNavigationContent) => any,
) => {
  return (
    contents &&
    contents.map((g) => (
      <div
        onClick={() => onClick(g)}
        key={g.id}
        className="hover:bg-sidebar-accent hover:text-sidebar-accent-foreground
             flex items-center gap-3 border-b p-4 text-sm leading-tight
             last:border-b-0 min-w-0"
      >
        <Avatar className="h-8 w-8 rounded-lg">
          <AvatarImage src={g.image} alt={g.name} className="bg-gray-profile" />
          <AvatarFallback className="rounded-lg">CN</AvatarFallback>
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
