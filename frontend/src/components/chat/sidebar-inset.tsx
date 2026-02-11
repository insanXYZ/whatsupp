import { ReactNode } from "react";
import { Separator } from "../ui/separator";
import { SidebarInset, SidebarTrigger } from "../ui/sidebar";
import { Paperclip } from "lucide-react";
import Image from "next/image";

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
      {/* <div className="flex flex-1 flex-col gap-4 p-4">{content}</div> */}
      <div className="relative flex flex-1 flex-col p-4">
        {/* chat list */}
        <div className="flex-1 overflow-y-auto pb-20 flex flex-col gap-4">
          <div className="flex justify-start">
            <div className="bg-blue-200 max-w-2/3 rounded p-2">
              Lorem ipsum dolor sit amet...
            </div>
          </div>

          <div className="flex justify-end">
            <div className="bg-green-200 max-w-2/3 rounded p-2">
              Lorem ipsum dolor sit amet...
            </div>
          </div>
        </div>

        {/* input chat */}
        <div className="absolute flex items-center gap-5 bottom-0 left-0 right-0 bg-background border-t p-4">
          <Paperclip />
          <input
            className="w-full rounded-md border px-3 py-2"
            placeholder="Ketik pesan..."
          />
        </div>
      </div>
    </SidebarInset>
  );
};

export const InsetHeaderGroup = (image: string, name: string) => {
  return (
    <div className="flex gap-5 items-center">
      <Image src={image} width={28} height={28} alt={name} />
      <div>{name}</div>
    </div>
  );
};
