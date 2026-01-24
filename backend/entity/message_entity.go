package entity

import "time"

type Message struct {
	ID        string    `gorm:"column:id;primaryKey"`
	GroupID   string    `gorm:"column:group_id"`
	SenderID  string    `gorm:"column:sender_id"`
	Message   string    `gorm:"column:message"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	Group     *Group    `gorm:"foreignKey:group_id;references:id"`
	User      *User     `gorm:"foreignKey:sender_id;references:id"`
}
