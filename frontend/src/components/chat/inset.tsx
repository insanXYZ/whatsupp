import { FormEvent, ReactNode, useState } from "react";
import { SidebarInset, SidebarTrigger } from "../ui/sidebar";
import { Separator } from "../ui/separator";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";
import { Paperclip, Send } from "lucide-react";
import { ButtonLoading } from "../ui/button-loading";

type AppSidebarInsetProps = {
  header?: ReactNode;
  content?: ReactNode;
};

export const AppSidebarInset = ({ header, content }: AppSidebarInsetProps) => {
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

type InsetHeaderGroupProps = {
  image: string;
  name: string;
};

export const InsetHeaderGroupProfile = ({
  image,
  name,
}: InsetHeaderGroupProps) => {
  return (
    <div className="flex gap-5 items-center">
      <Avatar className="h-7 w-7 rounded-lg">
        <AvatarImage src={image} alt={name} className="bg-gray-profile" />
        <AvatarFallback className="rounded-lg">CN</AvatarFallback>
      </Avatar>
      <div>{name}</div>
    </div>
  );
};

type InsetChatProps = {
  groupId?: number;
  receiverId?: number;
  messages: ItemGetMessageResponse[];
  onSubmit: (v: SendMessageRequest) => any;
};

export const InsetChat = ({
  groupId,
  receiverId,
  messages,
  onSubmit,
}: InsetChatProps) => {
  const [message, setMessage] = useState<string>("");

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    onSubmit({
      group_id: groupId,
      receiver_id: receiverId,
      message: message,
    });
  };

  return (
    <>
      <div className="flex-1 overflow-y-auto pb-20 flex flex-col gap-4">
        {messages?.map((v) =>
          v.is_me ? (
            <div key={v.id} className="flex justify-end">
              <div className="bg-blue-200 max-w-2/3 rounded p-2">
                {v.message}
              </div>
            </div>
          ) : (
            <div key={v.id} className="flex justify-start">
              <div className="bg-green-200 max-w-2/3 rounded p-2">
                {v.message}
              </div>
            </div>
          ),
        )}
      </div>

      <form
        onSubmit={handleSubmit}
        className="absolute flex items-center gap-5 bottom-0 left-0 right-0 bg-background border-t p-4"
      >
        <Paperclip />
        <input
          onChange={(v) => setMessage(v.target.value)}
          className="w-full rounded-md border px-3 py-2"
          placeholder="Write message..."
        />
        <ButtonLoading className="p-5 " isPending={false}>
          <Send />
        </ButtonLoading>
      </form>
    </>
  );
};
