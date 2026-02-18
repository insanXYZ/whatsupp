package dto

import "time"

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
	ConversationID   *int   `json:"conversation_id"`
}

type LoadRecentConversation = SearchConversationResponse

type NewConversationResponse = SearchConversationResponse
