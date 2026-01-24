package entity

import (
	"database/sql/driver"
	"time"
)

type typeGroupMember string

const (
	ADMIN  typeGroupMember = "ADMIN"
	MEMBER typeGroupMember = "MEMBER"
)

func (tg *typeGroupMember) Scan(value any) error {
	*tg = typeGroupMember(value.([]byte))
	return nil
}

func (tg typeGroupMember) Value() (driver.Value, error) {
	return string(tg), nil
}

type GroupMember struct {
	ID       string          `gorm:"column:id;primaryKey"`
	GroupID  string          `gorm:"column:group_id"`
	UserID   string          `gorm:"column:user_id"`
	Role     typeGroupMember `gorm:"column:role;type:role_group_member"`
	JoinedAt time.Time       `gorm:"column:joined_at;autoCreateTime"`
	Group    *Group          `gorm:"foreignKey:group_id;references:id"`
	User     *User           `gorm:"foreignKey:user_id;references:id"`
}

func (GroupMember) TableName() string {
	return "group_members"
}
