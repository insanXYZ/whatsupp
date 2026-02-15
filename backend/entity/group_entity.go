package entity

import (
	"time"
)

type typeGroup = string

const (
	GROUP    typeGroup = "GROUP"
	PERSONAL typeGroup = "PERSONAL"
)

type Group struct {
	ID        int       `gorm:"column:id;primaryKey,autoIncrement"`
	Name      string    `gorm:"column:name"`
	Bio       string    `gorm:"column:bio"`
	GroupType typeGroup `gorm:"column:type;type:type_group"`
	Image     string    `gorm:"column:image"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	Members   []Member  `gorm:"foreignKey:GroupId;references:ID"`
}
