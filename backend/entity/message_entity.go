package entity

import "time"

type Message struct {
	ID             int           `gorm:"column:id;primaryKey,autoIncrement"`
	ConversationID int           `gorm:"column:conversation_id"`
	UserID         int           `gorm:"column:user_id"`
	Message        string        `gorm:"column:message"`
	CreatedAt      time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time     `gorm:"column:updated_at"`
	Conversation   *Conversation `gorm:"foreignKey:ConversationID;references:ID"`
	User           *User         `gorm:"foreignKey:UserID;references:ID"`
}
