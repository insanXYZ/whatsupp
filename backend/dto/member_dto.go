package dto

type Member struct {
	ID             int           `json:"id,omitempty"`
	Role           string        `json:"role,omitempty"`
	ConversationId int           `json:"conversation_id,omitempty"`
	Conversation   *Conversation `json:"conversation,omitempty"`
	User           *User         `json:"user,omitempty"`
}
