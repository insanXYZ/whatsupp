package dto

type Member struct {
	ID           int           `json:"id,omitempty"`
	Role         string        `json:"role,omitempty"`
	Conversation *Conversation `json:"conversation,omitempty"`
	User         *User         `json:"user,omitempty"`
}
