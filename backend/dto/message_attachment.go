package dto

type MessageAttachment struct {
	ID        int      `json:"id,omitempty"`
	MessageID string   `json:"message_id,omitempty"`
	FileURL   string   `json:"file_url,omitempty"`
	FileType  string   `json:"file_type,omitempty"`
	Message   *Message `json:"message,omitempty"`
}
