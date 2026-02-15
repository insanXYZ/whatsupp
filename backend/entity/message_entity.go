package entity

import "time"

type Message struct {
	ID        int       `gorm:"column:id;primaryKey,autoIncrement"`
	MemberID  int       `gorm:"column:member_id"`
	Message   string    `gorm:"column:message"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	Member    *Member   `gorm:"foreignKey:MemberID;references:ID"`
}
