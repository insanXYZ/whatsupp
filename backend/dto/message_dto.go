package dto

import (
	"time"
)

type Message struct {
	ID           int           `json:"id,omitempty"`
	Message      string        `json:"message,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	User         *User         `json:"user,omitempty"`
	Conversation *Conversation `json:"conversation,omitempty"`
}

type eventMessageWs string
type targetSendMessage string

const (
	EVENT_SEND_MESSAGE     eventMessageWs = "SEND_MESSAGE"
	EVENT_NEW_MESSAGE      eventMessageWs = "NEW_MESSAGE"
	EVENT_NEW_CONVERSATION eventMessageWs = "NEW_CONVERSATION"
)

const (
	TYPE_TARGET_USER  targetSendMessage = "USER"
	TYPE_TARGET_GROUP targetSendMessage = "GROUP"
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
	Target         TargetSendMessage `json:"target"`
	Message        string            `json:"message,omitempty"`
	ConversationID *int              `json:"conversation_id,omitempty"`
}

type NewMessageResponse struct {
	IsMe           bool
	ConversationID int
	*Message
}

type ReadPumpClient = EventMessageWs

type GetMessageRequest struct {
	ConversationId int `param:"conversationId"`
}

type ItemGetMessagesResponse struct {
	IsMe    bool     `json:"is_me"`
	Message *Message `json:"message"`
}

type GetMessagesResponse struct {
	ConversationId int                        `json:"conversation_id,omitempty"`
	Message        []*ItemGetMessagesResponse `json:"message,omitempty"`
}
