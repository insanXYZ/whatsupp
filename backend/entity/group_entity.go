package entity

import (
	"database/sql/driver"
	"time"
)

type typeGroup string

const (
	GROUP    typeGroup = "GROUP"
	PERSONAL typeGroup = "PERSONAL"
)

func (tg *typeGroup) Scan(value any) error {
	*tg = typeGroup(value.([]byte))
	return nil
}

func (tg typeGroup) Value() (driver.Value, error) {
	return string(tg), nil
}

type Group struct {
	ID          string    `gorm:"column:id;primaryKey"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
	GroupType   typeGroup `gorm:"column:type;type:type_group"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
}
