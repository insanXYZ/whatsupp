package repository

import (
	"whatsupp-backend/entity"

	"gorm.io/gorm"
)

type MessageAttachmentRepository struct {
	*repository[*entity.Group]
}

func NewMessageAttachmentRepository(db *gorm.DB) *MessageAttachmentRepository {
	return &MessageAttachmentRepository{
		repository: &repository[*entity.Group]{
			DB: db,
		},
	}
}
