package entity

import (
	"time"
)

type typeConversation = string

const (
	CONV_TYPE_GROUP   typeConversation = "GROUP"
	CONV_TYPE_PRIVATE typeConversation = "PRIVATE"
)

type Conversation struct {
	ID               int              `gorm:"column:id;primaryKey,autoIncrement"`
	Name             string           `gorm:"column:name"`
	Bio              string           `gorm:"column:bio"`
	ConversationType typeConversation `gorm:"column:type;type:type_conversations"`
	Image            string           `gorm:"column:image"`
	CreatedAt        time.Time        `gorm:"column:created_at;autoCreateTime"`
	Members          []Member         `gorm:"foreignKey:ConversationID;references:ID"`
	Messages         []Message        `gorm:"foreignKey:ConversationID;references:ID"`
}
