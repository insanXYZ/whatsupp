package entity

import "time"

type Message struct {
	ID        int       `gorm:"column:id;primaryKey,autoIncrement"`
	MemberID  int       `gorm:"column:member_id"`
	Message   string    `gorm:"column:message"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	Member    *Member   `gorm:"foreignKey:member_id;references:id"`
}
