package repository

import (
	"whatsupp-backend/entity"

	"gorm.io/gorm"
)

type MessageAttachmentRepository struct {
	*repository[*entity.MessageAttachment]
}

func NewMessageAttachmentRepository(db *gorm.DB) *MessageAttachmentRepository {
	return &MessageAttachmentRepository{
		repository: &repository[*entity.MessageAttachment]{
			DB: db,
		},
	}
}

func (ma *MessageAttachmentRepository) WithTx(tx *gorm.DB) *MessageAttachmentRepository {
	return &MessageAttachmentRepository{
		repository: &repository[*entity.MessageAttachment]{
			DB: tx,
		},
	}
}
