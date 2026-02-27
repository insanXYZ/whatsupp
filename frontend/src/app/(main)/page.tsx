"use client";

import { ChatBannerLoading } from "@/components/chat/banner-loading";
import {
  AppSidebarInset,
  InsetChat,
  InsetHeaderConversationProfile,
} from "@/components/chat/inset";
import {
  AppSidebar,
  RenderRowsConversationChat,
} from "@/components/chat/sidebar";
import {
  CONVERSATION_TYPE_GROUP,
  CONVERSATION_TYPE_PRIVATE,
  RowConversationChat,
} from "@/dto/conversation-dto.ts";
import { GetMessageResponse } from "@/dto/message-dto";
import { EVENT_SEND_MESSAGE, EventWs, SendMessageRequest } from "@/dto/ws-dto";
import { useChatSocket } from "@/hooks/use-chat-socket";
import { useConversations } from "@/hooks/use-conversations";
import { useIdb } from "@/hooks/use-idb";
import { useMessages } from "@/hooks/use-message";
import {
  NAV_TITLE_CHAT,
  NAV_TITLE_CONTACTS,
  NAV_TITLE_GROUPS,
  NAV_TITLE_SEARCH,
} from "@/navigation/navigation";
import { HttpMethod, Mutation, useQueryData } from "@/utils/tanstack";
import { useEffect, useState } from "react";

export default function Page() {
  const [activeItem, setActiveItem] = useState<string>(NAV_TITLE_CHAT);
  const [recordHaveMutateMessage, setRecordHaveMutateMessage] = useState<
    Record<number, boolean>
  >({});

  const {
    activeChat,
    setActiveChat,
    conversations,
    addConversation,
    deleteConversationsByConversationId,
    overwriteConversations,
  } = useConversations();

  const { messages, AppendMessage, OverwriteMessages } = useMessages();

  const {
    SearchConversationsByNameIdb,
    GetAllConversationsIdb,
    SearchGroupConversationsByNameIdb,
    SearchPrivateConversationsByNameIdb,
    SearchConversationByIdIdb,
    DeleteConversationIdb,
    AppendConversationIdb,
    AppendConversationsIdb,
    AppendMessagesIdb,
    AppendMessageIdb,
    AppendMembersIdb,
    GetMessagesByConversationIdIdb,
    GetMembersByConversationIdIdb,
  } = useIdb();

  const { send, connected } = useChatSocket({
    onNewMessage: async (data) => {
      const message = data.message;
      if (message) {
        await AppendMessageIdb(message);

        if (activeChat?.conversation_id == data.conversation_id) {
          AppendMessage(message);
        }
      }
    },
    onNewConversation: async (data) => {
      try {
        AppendConversationIdb(data);

        if (activeItem === NAV_TITLE_SEARCH) {
          if (
            activeChat?.id === data.id &&
            activeChat.conversation_type === data.conversation_type
          ) {
            SetActiveChatWithMemberAndMessageConversation(data);
          }
        }

        switch (activeItem) {
          case NAV_TITLE_CHAT:
            addConversation(data);
            break;
          case NAV_TITLE_GROUPS:
            if (data.conversation_type === CONVERSATION_TYPE_GROUP) {
              addConversation(data);
            }
            break;
        }
      } catch (error) {
        console.log(error);
      }
    },
    onLeaveConversation: async (data) => {
      try {
        await DeleteConversationIdb(data.conversation_id);

        if (activeItem === NAV_TITLE_CHAT || activeItem === NAV_TITLE_GROUPS) {
          deleteConversationsByConversationId(data.conversation_id);
        }

        if (activeChat?.conversation_id === data.conversation_id) {
          SetActiveChatWithMemberAndMessageConversation(undefined);
        }
      } catch (err) {
        console.log(err);
      }
    },
    onMemberLeaveConversation: async (data) => { },
    onMemberJoinConversation: async (data) => { },
  });

  const {
    mutate: mutateGetMessages,
    isSuccess: isSuccessGetMessages,
    data: dataGetMessages,
  } = Mutation(["getMessages"]);

  const {
    mutate: mutateSearchConversations,
    isSuccess: isSuccessSearchConversations,
    data: dataSearchConversations,
  } = Mutation(["getConversations"]);

  const {
    isPending: isPendingMembershipGroupConversation,
    mutate: mutateMembershipGroupConversation,
  } = Mutation(["joinGroupConversation"], true);

  const {
    data: dataGetRecentConversations,
    isSuccess: isSuccessGetRecentConversations,
  } = useQueryData(["getRecentConversations"], "/conversations/recent");

  const handleSendMessage = (v: SendMessageRequest) => {
    const req: EventWs = {
      event: EVENT_SEND_MESSAGE,
      data: v,
    };

    send(req);
  };

  const onSubmitMembershipGroupConversation = (v: RowConversationChat) => {
    mutateMembershipGroupConversation({
      body: null,
      method: HttpMethod.PUT,
      url: `/conversations/${v.conversation_id}/members/me_join`,
    });
  };

  const SetActiveChatWithMemberAndMessageConversation = async (
    conversation: RowConversationChat | undefined,
  ) => {
    OverwriteMessages([]);

    if (conversation == undefined) {
      setActiveChat(undefined);
      return;
    }

    if (conversation.conversation_id && !conversation.members) {
      const members = await GetMembersByConversationIdIdb(
        conversation.conversation_id!,
      );
      conversation.members = members ?? members;
    }

    setActiveChat(conversation);

    if (conversation.conversation_id && conversation.have_joined) {
      if (!recordHaveMutateMessage[conversation.conversation_id]) {
        mutateGetMessages({
          body: null,
          method: HttpMethod.GET,
          url: "/messages/" + conversation.conversation_id,
        });
      } else {
        const messages = await GetMessagesByConversationIdIdb(
          conversation.conversation_id!,
        );
        OverwriteMessages(messages ? messages : []);
      }
    }
  };

  const onClickConversationChat = async (v: RowConversationChat) => {
    switch (v.conversation_type) {
      case CONVERSATION_TYPE_PRIVATE:
        if (v.id == activeChat?.id) {
          return;
        }
        break;
      case CONVERSATION_TYPE_GROUP:
        if (v.conversation_id == activeChat?.conversation_id) {
          return;
        }
        break;
    }
    SetActiveChatWithMemberAndMessageConversation(v);
  };

  const onSearch = async (v: string) => {
    try {
      switch (activeItem) {
        case NAV_TITLE_CHAT:
          const conversations = await SearchConversationsByNameIdb(v);
          overwriteConversations(conversations);
          break;
        case NAV_TITLE_SEARCH:
          if (v != "") {
            mutateSearchConversations({
              body: null,
              method: HttpMethod.GET,
              url: "/conversations?name=" + v,
            });
          }
          break;
        case NAV_TITLE_GROUPS:
          const groupConversations = await SearchGroupConversationsByNameIdb(v);
          overwriteConversations(groupConversations);
          break;
        case NAV_TITLE_CONTACTS:
          const privateConversations =
            await SearchPrivateConversationsByNameIdb(v);
          overwriteConversations(privateConversations);
          break;
      }
    } catch (error) {
      console.log(error);
    }
  };

  const onChangeActiveItem = async (v: string) => {
    if (v === activeItem) return;
    setActiveItem(v);
    overwriteConversations([]);
    try {
      switch (v) {
        case NAV_TITLE_SEARCH:
          break;
        case NAV_TITLE_CHAT:
          const allConversations = await GetAllConversationsIdb();
          overwriteConversations(allConversations);
          break;
        case NAV_TITLE_GROUPS:
          const groupConversations =
            await SearchGroupConversationsByNameIdb("");
          overwriteConversations(groupConversations);
          break;
        case NAV_TITLE_CONTACTS:
          const privateConversations =
            await SearchPrivateConversationsByNameIdb("");
          overwriteConversations(privateConversations);
          break;
      }
    } catch (error) {
      console.log(error);
    }
  };

  useEffect(() => {
    if (isSuccessGetMessages && dataGetMessages.data) {
      const data = dataGetMessages.data as GetMessageResponse;
      const messages = data.messages;

      setRecordHaveMutateMessage((prev) => ({
        ...prev,
        [data.conversation_id]: true,
      }));

      AppendMessagesIdb(messages);

      if (activeChat?.conversation_id == data.conversation_id) {
        OverwriteMessages(messages);
      }
    }
  }, [isSuccessGetMessages]);

  useEffect(() => {
    if (isSuccessGetRecentConversations) {
      const conversations = dataGetRecentConversations.data
        ? (dataGetRecentConversations.data as RowConversationChat[])
        : [];

      const conversationsIdbVal: RowConversationChat[] = conversations.map(
        (conversation) => {
          return {
            conversation_type: conversation.conversation_type,
            have_joined: conversation.have_joined,
            bio: conversation.bio,
            id: conversation.id,
            image: conversation.image,
            name: conversation.name,
            conversation_id: conversation.conversation_id,
          };
        },
      );

      AppendConversationsIdb(conversationsIdbVal);

      AppendMembersIdb(
        conversations.flatMap((conversation) => conversation.members!),
      );

      if (activeItem === NAV_TITLE_CHAT) {
        overwriteConversations(conversations);
      }
    }
  }, [isSuccessGetRecentConversations]);

  useEffect(() => {
    if (isSuccessSearchConversations && activeItem == NAV_TITLE_SEARCH) {
      const data = dataSearchConversations.data as RowConversationChat[];
      overwriteConversations(data);
    }
  }, [isSuccessSearchConversations]);

  return !connected && !isSuccessGetRecentConversations ? (
    <ChatBannerLoading />
  ) : (
    <>
      <AppSidebar
        activeItem={activeItem}
        onChangeActiveItem={onChangeActiveItem}
        onSearch={onSearch}
        contentSidebarDetail={
          <RenderRowsConversationChat
            conversations={conversations}
            onClick={onClickConversationChat}
          />
        }
      />
      <AppSidebarInset
        header={
          activeChat && (
            <InsetHeaderConversationProfile conversation={activeChat} />
          )
        }
        content={
          activeChat && (
            <InsetChat
              isPendingJoin={isPendingMembershipGroupConversation}
              onSubmitMembershipGroupConversation={
                onSubmitMembershipGroupConversation
              }
              messages={messages}
              onSubmit={handleSendMessage}
              conversationDetail={activeChat}
            />
          )
        }
      />
    </>
  );
}
