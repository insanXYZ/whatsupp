package repository

import (
	"whatsupp-backend/entity"

	"gorm.io/gorm"
)

type MessageRepository struct {
	*repository[*entity.Group]
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{
		repository: &repository[*entity.Group]{
			DB: db,
		},
	}
}
