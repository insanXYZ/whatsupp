package entity

import (
	"time"
)

type typeMember string

const (
	ADMIN  typeMember = "ADMIN"
	MEMBER typeMember = "MEMBER"
)

type Member struct {
	ID       int        `gorm:"column:id;primaryKey,autoIncrement"`
	GroupID  int        `gorm:"column:group_id"`
	UserID   int        `gorm:"column:user_id"`
	Role     typeMember `gorm:"column:role;type:role_group_member"`
	JoinedAt time.Time  `gorm:"column:joined_at;autoCreateTime"`
	Group    *Group     `gorm:"foreignKey:GroupID;references:ID"`
	User     *User      `gorm:"foreignKey:UserID;references:ID"`
}

func (Member) TableName() string {
	return "members"
}
