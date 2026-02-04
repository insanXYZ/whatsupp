package dto

const (
	MULTIPART_FORM_NAME = "file"
)

type MessageAttachment struct {
	ID        int      `json:"id,omitempty"`
	MessageID string   `json:"message_id,omitempty"`
	FileURL   string   `json:"file_url,omitempty"`
	FileExt   string   `json:"file_ext,omitempty"`
	Message   *Message `json:"message,omitempty"`
}
