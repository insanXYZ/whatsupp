import { FormEvent, ReactNode, useEffect, useState } from "react";
import { SidebarInset, SidebarTrigger } from "../ui/sidebar";
import { Separator } from "../ui/separator";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";
import { LogOut, Paperclip, Send } from "lucide-react";
import { ButtonLoading } from "../ui/button-loading";
import { SendMessageRequest } from "@/dto/ws-dto";
import {
  CONVERSATION_TYPE_GROUP,
  CONVERSATION_TYPE_PRIVATE,
  EditGroupConversationDto,
  RowConversationChat,
} from "@/dto/conversation-dto.ts";
import { ItemGetMessageResponse } from "@/dto/message-dto";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "../ui/dialog";
import { ContentType, HttpMethod, Mutation } from "@/utils/tanstack";
import z from "zod";
import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Field, FieldError, FieldGroup, FieldLabel } from "../ui/field";
import { Input } from "../ui/input";
import { Button } from "../ui/button";
import { CropImageDialog } from "./sidebar";
import { useAccount } from "@/hooks/use-account";
import { MEMBER_ROLE_ADMIN } from "@/dto/member-dto";

type AppSidebarInsetProps = {
  header?: ReactNode;
  content?: ReactNode;
};

export const AppSidebarInset = ({ header, content }: AppSidebarInsetProps) => {
  return (
    <SidebarInset>
      <header className="bg-background sticky top-0 z-50 flex shrink-0 items-center gap-2 border-b p-4">
        <SidebarTrigger className="-ml-1" />
        <Separator
          orientation="vertical"
          className="mr-2 data-[orientation=vertical]:h-4"
        />
        {header}
      </header>
      {content}
    </SidebarInset>
  );
};

type InsetHeaderConversationProps = {
  conversation: RowConversationChat;
};

export const InsetHeaderConversationProfile = ({
  conversation,
}: InsetHeaderConversationProps) => {
  const { user } = useAccount();

  const { mutate: mutateEditGroup, isPending: isPendingEditGroup } = Mutation(
    ["editGroup"],
    true,
  );

  const { mutate: mutateLeaveGroup, isPending: isPendingLeaveGroup } = Mutation(
    ["memberGroup"],
    true,
  );

  const [imageSrc, setImageSrc] = useState<string>();
  const [tempImage, setTempImage] = useState<string>();
  const [cropDialogOpen, setCropDialogOpen] = useState(false);
  const [formDialogOpen, setFormDialogOpen] = useState(false);
  const [isAdmin, setIsAdmin] = useState<boolean>(false);

  const form = useForm<z.infer<typeof EditGroupConversationDto>>({
    defaultValues: {
      image: undefined,
      name: conversation.name,
      bio: conversation.bio,
    },
    resolver: zodResolver(EditGroupConversationDto),
  });

  const onSubmit = (data: z.infer<typeof EditGroupConversationDto>) => {
    const formData = new FormData();
    if (data.image) {
      formData.append("image", data.image!);
    }
    if (data.bio) {
      formData.append("bio", data.bio!);
    }
    formData.append("name", data.name);

    mutateEditGroup({
      body: formData,
      method: HttpMethod.PUT,
      url: `/conversations/${conversation.id}`,
      contentType: ContentType.FORM,
    });
  };

  const onClickLeaveGroup = () => {
    mutateLeaveGroup({
      body: null,
      method: HttpMethod.PUT,
      url: `/conversations/${conversation.conversation_id}/members/me_join`,
    });
  };

  useEffect(() => {
    if (conversation.members) {
      const filtered = conversation.members.find(
        (member) => member.user.id === user?.id,
      );

      if (filtered) {
        setIsAdmin(filtered.role === MEMBER_ROLE_ADMIN);
      }
    }
  }, [conversation]);

  useEffect(() => {
    return () => {
      if (imageSrc) URL.revokeObjectURL(imageSrc);
      if (tempImage) URL.revokeObjectURL(tempImage);
    };
  }, [imageSrc, tempImage]);

  return (
    <>
      <Dialog open={formDialogOpen} onOpenChange={setFormDialogOpen}>
        <DialogTrigger asChild>
          <div className="w-full cursor-pointer flex items-center ">
            <div className="flex gap-5 items-center">
              <Avatar className="h-7 w-7 rounded-lg">
                <AvatarImage
                  src={conversation.image}
                  alt={conversation.name}
                  className="bg-gray-profile"
                />
                <AvatarFallback className="rounded-lg">
                  {conversation.name.slice(0, 2).toUpperCase()}
                </AvatarFallback>
              </Avatar>
              <div>{conversation.name}</div>
            </div>
          </div>
        </DialogTrigger>

        <DialogContent showCloseButton={false} className="sm:max-w-sm">
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <DialogHeader>
              <DialogTitle>
                Detail{" "}
                {conversation.conversation_type === CONVERSATION_TYPE_PRIVATE
                  ? "user"
                  : "group"}
              </DialogTitle>
            </DialogHeader>

            <FieldGroup>
              <Controller
                control={form.control}
                name="image"
                render={({ field: { ref, name, onBlur }, fieldState }) => (
                  <Field>
                    <FieldLabel htmlFor="image">
                      <div className="w-full flex items-center justify-center">
                        {imageSrc ? (
                          <Avatar className="w-24 h-24">
                            <AvatarImage
                              className="bg-gray-profile"
                              src={imageSrc}
                              alt={imageSrc}
                            />
                            <AvatarFallback>
                              {conversation.name.slice(0, 2).toUpperCase()}
                            </AvatarFallback>
                          </Avatar>
                        ) : (
                          <Avatar className="w-24 h-24">
                            <AvatarImage
                              className="bg-gray-profile"
                              src={conversation.image}
                              alt={conversation.name}
                            />
                            <AvatarFallback>
                              {conversation.name.slice(0, 2).toUpperCase()}
                            </AvatarFallback>
                          </Avatar>
                        )}
                      </div>
                    </FieldLabel>

                    <Input
                      id="image"
                      name={name}
                      ref={ref}
                      type="file"
                      accept="image/*"
                      onBlur={onBlur}
                      hidden
                      disabled={
                        conversation.conversation_type ===
                        CONVERSATION_TYPE_PRIVATE ||
                        (conversation.conversation_type ===
                          CONVERSATION_TYPE_GROUP &&
                          !isAdmin)
                      }
                      onChange={(e) => {
                        const file = e.target.files?.[0];
                        if (!file) return;

                        const url = URL.createObjectURL(file);
                        setTempImage(url);
                        setCropDialogOpen(true);
                      }}
                    />

                    <FieldError errors={[fieldState.error]} />
                  </Field>
                )}
              />

              <div className="flex justify-center">
                <ButtonLoading
                  isPending={isPendingLeaveGroup}
                  onClick={onClickLeaveGroup}
                  className="bg-red-600"
                  type="button"
                >
                  <LogOut color="white" className="text-white" />
                </ButtonLoading>
              </div>

              <Controller
                control={form.control}
                name="name"
                render={({ field, fieldState }) => (
                  <Field>
                    <FieldLabel htmlFor="name">Name</FieldLabel>
                    <Input
                      {...field}
                      id="name"
                      disabled={
                        conversation.conversation_type ===
                        CONVERSATION_TYPE_PRIVATE ||
                        (conversation.conversation_type ===
                          CONVERSATION_TYPE_GROUP &&
                          !isAdmin)
                      }
                    />
                    <FieldError errors={[fieldState.error]} />
                  </Field>
                )}
              />

              <Controller
                control={form.control}
                name="bio"
                render={({ field, fieldState }) => (
                  <Field>
                    <FieldLabel htmlFor="bio">Bio</FieldLabel>
                    <Input
                      {...field}
                      id="bio"
                      disabled={
                        conversation.conversation_type ===
                        CONVERSATION_TYPE_PRIVATE ||
                        (conversation.conversation_type ===
                          CONVERSATION_TYPE_GROUP &&
                          !isAdmin)
                      }
                    />
                    <FieldError errors={[fieldState.error]} />
                  </Field>
                )}
              />

              {conversation.conversation_type != CONVERSATION_TYPE_PRIVATE &&
                isAdmin && (
                  <DialogFooter>
                    <DialogClose asChild>
                      <Button variant="outline" type="button">
                        Cancel
                      </Button>
                    </DialogClose>

                    <ButtonLoading type="submit" isPending={isPendingEditGroup}>
                      Update
                    </ButtonLoading>
                  </DialogFooter>
                )}
            </FieldGroup>
          </form>
        </DialogContent>
      </Dialog>

      <CropImageDialog
        open={cropDialogOpen}
        onOpenChange={setCropDialogOpen}
        tempImage={tempImage}
        onCropComplete={(file, previewUrl) => {
          form.setValue("image", file);
          setImageSrc(previewUrl);
        }}
      />
    </>
  );
};

type InsetChatProps = {
  conversationDetail: RowConversationChat;
  messages: ItemGetMessageResponse[];
  isPendingJoin: boolean;
  onSubmitMembershipGroupConversation: (v: RowConversationChat) => any;
  onSubmit: (v: SendMessageRequest) => any;
};

export const InsetChat = ({
  conversationDetail,
  messages,
  onSubmit,
  isPendingJoin,
  onSubmitMembershipGroupConversation,
}: InsetChatProps) => {
  const [message, setMessage] = useState<string>("");

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const req: SendMessageRequest = {
      message: message,
      conversation_id:
        conversationDetail.conversation_id ??
        conversationDetail.conversation_id,
      tmp_conversation_id: conversationDetail.conversation_id
        ? undefined
        : `tmp-${conversationDetail.conversation_type}-${conversationDetail.id}`,
      target: {
        id: conversationDetail.id,
        type: conversationDetail.conversation_type,
      },
    };

    onSubmit(req);
    setMessage("");
  };

  return (
    <>
      <div className="flex flex-1 flex-col gap-4 p-4 z-0 ">
        <div className="flex-1 overflow-y-auto pb-20 flex flex-col gap-4">
          {messages?.map((v) =>
            v.is_me ? (
              <div key={v.id} className="flex justify-end">
                <div className="bg-blue-200 flex items-end flex-col max-w-2/3 rounded p-2">
                  <div>{v.message}</div>
                </div>
              </div>
            ) : (
              <div key={v.id} className="flex gap-2 max-w-2/3">
                {conversationDetail.conversation_type ===
                  CONVERSATION_TYPE_GROUP && (
                    <Avatar>
                      <AvatarImage
                        src={v.user.image}
                        className="bg-gray-profile"
                      />
                      <AvatarFallback>
                        {v.user.name.slice(0, 2).toUpperCase()}
                      </AvatarFallback>
                    </Avatar>
                  )}
                <div className="flex justify-start">
                  <div className="bg-green-200 flex items-start flex-col rounded p-2">
                    <div className="text-xs text-gray-800">
                      {v.user.name}#{v.user.id}
                    </div>
                    <div>{v.message}</div>
                  </div>
                </div>
              </div>
            ),
          )}
        </div>
      </div>

      {conversationDetail.have_joined ||
        conversationDetail.conversation_type == CONVERSATION_TYPE_PRIVATE ? (
        <form
          onSubmit={handleSubmit}
          className="sticky flex items-center gap-5 bottom-0 z-50 bg-background border-t p-4"
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
      ) : (
        <div className="absolute flex items-center gap-5 bottom-0 left-0 right-0 bg-background border-t p-4">
          <ButtonLoading
            className="w-full"
            isPending={isPendingJoin}
            onClick={() =>
              onSubmitMembershipGroupConversation(conversationDetail)
            }
          >
            Join
          </ButtonLoading>
        </div>
      )}
    </>
  );
};
