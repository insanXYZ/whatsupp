import { ReactNode } from "react";
import { Separator } from "../ui/separator";
import { SidebarInset, SidebarTrigger } from "../ui/sidebar";

export const AppSidebarInset = ({
  header,
  content,
}: {
  header?: ReactNode;
  content?: ReactNode;
}) => {
  return (
    <SidebarInset>
      <header className="bg-background sticky top-0 flex shrink-0 items-center gap-2 border-b p-4">
        <SidebarTrigger className="-ml-1" />
        <Separator
          orientation="vertical"
          className="mr-2 data-[orientation=vertical]:h-4"
        />
        {header}
      </header>
      <div className="flex flex-1 flex-col gap-4 p-4">{content}</div>
    </SidebarInset>
  );
};
