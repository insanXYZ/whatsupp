package dto

import (
	"mime/multipart"
	"time"
)

type Conversation struct {
	ID               int       `json:"id,omitempty"`
	Name             string    `json:"name,omitempty"`
	Bio              string    `json:"bio,omitempty"`
	ConversationType string    `json:"conversation_type,omitempty"`
	Image            string    `json:"image,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	Members          []Member  `json:"members,omitempty"`
	Messages         []Message `json:"messages,omitempty"`
}

type SearchConversationRequest struct {
	Name string `query:"name"`
}

type SearchConversationResponse struct {
	ID               int    `json:"id,omitempty"`
	Name             string `json:"name,omitempty"`
	Image            string `json:"image,omitempty"`
	Bio              string `json:"bio,omitempty"`
	ConversationType string `json:"conversation_type,omitempty"`
	ConversationID   *int   `json:"conversation_id,omitempty"`
	HaveJoined       bool   `json:"have_joined"`
}

type LoadRecentConversation struct {
	*SearchConversationResponse
	Members []*Member `json:"members"`
}

type NewConversationResponse = SearchConversationResponse

type CreateGroupConversationRequest struct {
	Name  string                `form:"name" validate:"required,min=3,max=25"`
	Bio   string                `form:"bio" `
	Image *multipart.FileHeader `validate:"-"`
}

type UpdateGroupConversationRequest struct {
	ConversationId int                   `param:"conversationId"`
	Name           string                `form:"name" validate:"required,min=3,max=25"`
	Bio            string                `form:"bio" `
	Image          *multipart.FileHeader `validate:"-"`
}

type ListMembersConversationRequest struct {
	ConversationID int `param:"conversationId"`
}

type JoinGroupConversationRequest struct {
	ConversationID int `param:"conversationId"`
}

type LeaveConversationResponse struct {
	ConversationID int `json:"conversation_id,omitempty"`
}
