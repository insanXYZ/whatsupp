package entity

import "time"

type MessageAttachment struct {
	ID        int       `gorm:"column:id;primaryKey,autoIncrement"`
	MessageID int       `gorm:"column:message_id"`
	FileURL   string    `gorm:"column:file_url"`
	FileExt   string    `gorm:"column:file_type"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	Message   *Message  `gorm:"foreignKey:MessageID;references:ID"`
}

func (MessageAttachment) TableName() string {
	return "message_attachments"
}
