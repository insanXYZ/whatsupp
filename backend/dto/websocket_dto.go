package dto

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
