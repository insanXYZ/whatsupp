package entity

import "time"

type User struct {
	ID        int       `gorm:"column:id;primaryKey,autoIncrement"`
	Name      string    `gorm:"column:name"`
	Email     string    `gorm:"column:email"`
	Password  string    `gorm:"column:password"`
	Image     string    `gorm:"column:image"`
	Bio       string    `gorm:"column:bio"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}
