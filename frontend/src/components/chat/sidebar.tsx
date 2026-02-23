import {
  BadgeCheck,
  Camera,
  ChevronsUpDown,
  LoaderCircle,
  LogOut,
  MessageSquarePlus,
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
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  NAV_TITLE_CHAT,
  NAV_TITLE_GROUPS,
  Navigations,
} from "@/navigation/navigation";
import { ReactNode, useEffect, useRef, useState } from "react";
import { useDebounce } from "use-debounce";
import {
  ContentType,
  HttpMethod,
  Mutation,
  useQueryData,
} from "@/utils/tanstack";
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
import {
  CreateGroupConversationDto,
  RowConversationChat,
} from "@/dto/conversation-dto.ts";
import { Tooltip, TooltipContent, TooltipTrigger } from "../ui/tooltip";
import { Controller, useForm } from "react-hook-form";
import z from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Field, FieldError, FieldGroup, FieldLabel } from "../ui/field";
import { Input } from "../ui/input";
import { Button } from "../ui/button";
import { ButtonLoading } from "../ui/button-loading";
import Image from "next/image";
import ReactCrop, { Crop, centerCrop, makeAspectCrop } from "react-image-crop";
import "react-image-crop/dist/ReactCrop.css";
import { DeleteDbIdb } from "@/utils/indexdb";

type AppSidebarProps = {
  contentSidebarDetail?: ReactNode;
  activeItem: string;
  onChangeActiveItem: (v: string) => void;
  onSearch: (v: string) => void;
};

export const AppSidebar = ({
  contentSidebarDetail,
  onChangeActiveItem,
  activeItem,
  onSearch,
}: AppSidebarProps) => {
  const [prevSearch, setPrevSearch] = useState<string>("");
  const [search, setSearch] = useState<string>("");
  const [searchDebounce] = useDebounce(search, 600);

  useEffect(() => {
    if (searchDebounce != prevSearch) {
      setPrevSearch(searchDebounce);
      onSearch(searchDebounce);
    }
  }, [searchDebounce]);

  useEffect(() => {
    onChangeActiveItem(NAV_TITLE_CHAT);
  }, []);

  const handleChangeActiveItem = (v: string) => {
    if (v != activeItem) {
      setSearch("");
      onChangeActiveItem(v);
    }
  };

  return (
    <Sidebar
      collapsible="icon"
      className="overflow-hidden *:data-[sidebar=sidebar]:flex-row"
    >
      <SidebarNavigation
        activeItem={activeItem}
        onActiveItemChange={handleChangeActiveItem}
      />
      <SidebarDetail
        content={contentSidebarDetail}
        title={activeItem}
        search={search}
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
  search: string;
  onSearch: (v: string) => void;
};

const SidebarDetail = ({
  title,
  content,
  search,
  onSearch,
}: SidebarDetailProps) => {
  return (
    <Sidebar collapsible="none" className="hidden flex-1 md:flex">
      <SidebarHeader className="gap-3.5 border-b p-4">
        <div className="flex w-full items-center justify-between">
          <div className="text-foreground text-base font-medium">{title}</div>

          {title === NAV_TITLE_GROUPS && <FormDialogNewGroup />}
        </div>
        <SidebarInput
          value={search}
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

interface CropImageDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  tempImage: string | undefined;
  onCropComplete: (file: File, previewUrl: string) => void;
}

const getCroppedImg = async (
  image: HTMLImageElement,
  crop: Crop,
  fileName: string,
): Promise<File> => {
  const canvas = document.createElement("canvas");

  const scaleX = image.naturalWidth / image.width;
  const scaleY = image.naturalHeight / image.height;

  canvas.width = crop.width!;
  canvas.height = crop.height!;

  const ctx = canvas.getContext("2d");

  ctx?.drawImage(
    image,
    crop.x! * scaleX,
    crop.y! * scaleY,
    crop.width! * scaleX,
    crop.height! * scaleY,
    0,
    0,
    crop.width!,
    crop.height!,
  );

  return new Promise((resolve) => {
    canvas.toBlob((blob) => {
      if (!blob) return;
      resolve(new File([blob], fileName, { type: "image/jpeg" }));
    }, "image/jpeg");
  });
};

const CropImageDialog = ({
  open,
  onOpenChange,
  tempImage,
  onCropComplete,
}: CropImageDialogProps) => {
  const [crop, setCrop] = useState<Crop>();
  const [completedCrop, setCompletedCrop] = useState<Crop>();
  const imgRef = useRef<HTMLImageElement | null>(null);

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>Crop Image</DialogTitle>
        </DialogHeader>

        {tempImage && (
          <ReactCrop
            crop={crop}
            onChange={(c) => setCrop(c)}
            onComplete={(c) => setCompletedCrop(c)}
            aspect={1}
            minWidth={100}
          >
            <img
              ref={imgRef}
              src={tempImage}
              alt="Crop"
              style={{ maxHeight: 400 }}
              onLoad={(e) => {
                const { width, height } = e.currentTarget;
                const crop = makeAspectCrop(
                  { unit: "%", width: 90 },
                  1,
                  width,
                  height,
                );
                setCrop(centerCrop(crop, width, height));
              }}
            />
          </ReactCrop>
        )}

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>

          <Button
            onClick={async () => {
              if (!imgRef.current || !completedCrop) return;

              const croppedFile = await getCroppedImg(
                imgRef.current,
                completedCrop,
                "group.jpg",
              );

              const previewUrl = URL.createObjectURL(croppedFile);
              onCropComplete(croppedFile, previewUrl);
              onOpenChange(false);
            }}
          >
            OK
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

const FormDialogNewGroup = () => {
  const { mutate, isPending, isSuccess } = Mutation(["createGroup"], true);

  const [imageSrc, setImageSrc] = useState<string>();
  const [tempImage, setTempImage] = useState<string>();
  const [cropDialogOpen, setCropDialogOpen] = useState(false);
  const [formDialogOpen, setFormDialogOpen] = useState(false);

  const form = useForm<z.infer<typeof CreateGroupConversationDto>>({
    defaultValues: {
      image: undefined,
      name: "",
      bio: undefined,
    },
    resolver: zodResolver(CreateGroupConversationDto),
  });

  const onSubmit = (data: z.infer<typeof CreateGroupConversationDto>) => {
    const formData = new FormData();
    if (data.image) {
      formData.append("image", data.image!);
    }
    if (data.bio) {
      formData.append("bio", data.bio!);
    }
    formData.append("name", data.name);

    mutate({
      body: formData,
      method: HttpMethod.POST,
      url: "/conversations",
      contentType: ContentType.FORM,
    });
  };

  useEffect(() => {
    if (isSuccess) {
      setFormDialogOpen(false);
    }
  }, [isSuccess]);

  useEffect(() => {
    return () => {
      if (imageSrc) URL.revokeObjectURL(imageSrc);
      if (tempImage) URL.revokeObjectURL(tempImage);
    };
  }, [imageSrc, tempImage]);

  return (
    <>
      <Dialog open={formDialogOpen} onOpenChange={setFormDialogOpen}>
        <Tooltip>
          <TooltipTrigger asChild>
            <DialogTrigger asChild>
              <MessageSquarePlus className="cursor-pointer" />
            </DialogTrigger>
          </TooltipTrigger>
          <TooltipContent side={"bottom"}>
            <p>New Group</p>
          </TooltipContent>
        </Tooltip>

        <DialogContent showCloseButton={false} className="sm:max-w-sm">
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <DialogHeader>
              <DialogTitle>Create Group</DialogTitle>
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
                            <AvatarImage src={imageSrc} alt={imageSrc} />
                            <AvatarFallback>IM</AvatarFallback>
                          </Avatar>
                        ) : (
                          <Camera
                            width={70}
                            className="bg-secondary p-4 rounded-full cursor-pointer text-white"
                            height={70}
                          />
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

              <Controller
                control={form.control}
                name="name"
                render={({ field, fieldState }) => (
                  <Field>
                    <FieldLabel htmlFor="name">Name</FieldLabel>
                    <Input {...field} id="name" />
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
                    <Input {...field} id="bio" />
                    <FieldError errors={[fieldState.error]} />
                  </Field>
                )}
              />

              <DialogFooter>
                <DialogClose asChild>
                  <Button variant="outline" type="button">
                    Cancel
                  </Button>
                </DialogClose>

                <ButtonLoading type="submit" isPending={isPending}>
                  Create
                </ButtonLoading>
              </DialogFooter>
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
      console.log(successLogout);
      const deleteDbAndReload = async () => {
        try {
          await DeleteDbIdb();
        } catch (error) {
          console.log(error);
        }
      };

      deleteDbAndReload();

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
