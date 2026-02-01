package entity

import (
	"database/sql/driver"
	"time"
)

type typeMember string

const (
	ADMIN  typeMember = "ADMIN"
	MEMBER typeMember = "MEMBER"
)

func (tg *typeMember) Scan(value any) error {
	*tg = typeMember(value.([]byte))
	return nil
}

func (tg typeMember) Value() (driver.Value, error) {
	return string(tg), nil
}

type Member struct {
	ID       int        `gorm:"column:id;primaryKey,autoIncrement"`
	GroupID  string     `gorm:"column:group_id"`
	UserID   string     `gorm:"column:user_id"`
	Role     typeMember `gorm:"column:role;type:role_group_member"`
	JoinedAt time.Time  `gorm:"column:joined_at;autoCreateTime"`
	Group    *Group     `gorm:"foreignKey:group_id;references:id"`
	User     *User      `gorm:"foreignKey:user_id;references:id"`
}

func (Member) TableName() string {
	return "members"
}
