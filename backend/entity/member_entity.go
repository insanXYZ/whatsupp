package entity

import (
	"time"
)

type typeMember string

const (
	MEMBER_TYPE_ADMIN  typeMember = "ADMIN"
	MEMBER_TYPE_MEMBER typeMember = "MEMBER"
)

type Member struct {
	ID             int           `gorm:"column:id;primaryKey,autoIncrement"`
	ConversationID int           `gorm:"column:conversation_id"`
	UserID         int           `gorm:"column:user_id"`
	Role           typeMember    `gorm:"column:role;type:role_member"`
	JoinedAt       time.Time     `gorm:"column:joined_at;autoCreateTime"`
	Conversation   *Conversation `gorm:"foreignKey:ConversationID;references:ID"`
	User           *User         `gorm:"foreignKey:UserID;references:ID"`
}

func (Member) TableName() string {
	return "members"
}
