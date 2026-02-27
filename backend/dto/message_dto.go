package dto

import (
	"time"
)

type Message struct {
	ID             int           `json:"id,omitempty"`
	Message        string        `json:"message,omitempty"`
	CreatedAt      time.Time     `json:"created_at,omitempty"`
	ConversationID int           `json:"conversation_id,omitempty"`
	User           *User         `json:"user,omitempty"`
	Conversation   *Conversation `json:"conversation,omitempty"`
}

type eventMessageWs string
type targetSendMessage string

const (
	EVENT_SEND_MESSAGE              eventMessageWs = "SEND_MESSAGE"
	EVENT_NEW_MESSAGE               eventMessageWs = "NEW_MESSAGE"
	EVENT_NEW_CONVERSATION          eventMessageWs = "NEW_CONVERSATION"
	EVENT_MEMBER_JOIN_CONVERSATION  eventMessageWs = "MEMBER_JOIN_CONVERSATION"
	EVENT_LEAVE_CONVERSATION        eventMessageWs = "LEAVE_CONVERSATION"
	EVENT_MEMBER_LEAVE_CONVERSATION eventMessageWs = "MEMBER_LEAVE_CONVERSATION"
)

const (
	TYPE_TARGET_PRIVATE targetSendMessage = "PRIVATE"
	TYPE_TARGET_GROUP   targetSendMessage = "GROUP"
)

type EventMessageWs struct {
	Event string `json:"event,omitempty"`
	Data  any    `json:"data,omitempty"`
}

type TargetSendMessage struct {
	Type targetSendMessage `json:"type,omitempty"`
	ID   int               `json:"id,omitempty"`
}

type SendMessageRequestWs struct {
	Target            TargetSendMessage `json:"target,omitempty"`
	Message           string            `json:"message,omitempty"`
	TmpConversationID *string           `json:"tmp_conversation_id,omitempty"`
	ConversationID    *int              `json:"conversation_id,omitempty"`
}

type MessageOnNewMessageResponse = ItemGetMessagesResponse

type NewMessageResponse struct {
	ConversationID    int                          `json:"conversation_id,omitempty"`
	TmpConversationID *string                      `json:"tmp_conversation_id,omitempty"`
	Message           *MessageOnNewMessageResponse `json:"message,omitempty"`
}

type ReadPumpClient = EventMessageWs

type GetMessageRequest struct {
	ConversationId int `param:"conversationId"`
}

type ItemGetMessagesResponse struct {
	IsMe bool `json:"is_me"`
	*Message
}

type GetMessagesResponse struct {
	ConversationId int                        `json:"conversation_id,omitempty"`
	Message        []*ItemGetMessagesResponse `json:"messages,omitempty"`
}
