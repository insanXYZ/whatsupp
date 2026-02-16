package dto

import "time"

type Message struct {
	ID        int       `json:"id,omitempty"`
	Message   string    `json:"message,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	Member    *Member   `json:"member,omitempty"`
}

type MessageWS struct {
	Message    string `json:"message,omitempty"`
	GroupID    *int   `json:"group_id,omitempty"`
	ReceiverID *int   `json:"receiver_id,omitempty"`
}

type BroadcastMessageWS struct {
	*MessageWS
	ClientID  int
	MessageID int
}

type SyncMessageWS struct {
	Success bool
}

type GetMessageRequest struct {
	GroupID int `param:"groupId"`
}

type GetMessagesResponse struct {
	GroupId  int                        `json:"group_id,omitempty"`
	Messages []*ItemGetMessagesResponse `json:"messages,omitempty"`
}

type ItemGetMessagesResponse struct {
	IsMe bool `json:"is_me"`
	Message
}
