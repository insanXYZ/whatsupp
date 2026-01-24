package entity

import "time"

type MessageAttachment struct {
	ID        string    `gorm:"column:id;primaryKey"`
	MessageID string    `gorm:"column:message_id"`
	FileURL   string    `gorm:"column:file_url"`
	FileType  string    `gorm:"column:file_type"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	Message   *Message  `gorm:"foreignKey:message_id;references:id"`
}

func (MessageAttachment) TableName() string {
	return "message_attachments"
}
