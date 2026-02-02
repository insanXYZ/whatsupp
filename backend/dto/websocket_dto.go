package dto

type event string

const (
	SEND_MESSAGE            event = "SEND_MESSAGE"
	SEND_MESSAGE_ATTACHMENT event = "SEND_MESSAGE_ATTACHMENT"
	SYNC_MESSAGE            event = "SYNC_MESSAGE"
)

type MessageWS struct {
	Event event `json:"event,omitempty"`
	Data  any   `json:"data,omitempty"`
}

type SyncMessageWS struct {
	Success bool
}

type SendMessageWS struct {
	Message string `json:"message,omitempty"`
	GroupID int    `json:"group_id,omitempty"`
}
