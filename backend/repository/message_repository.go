package repository

import (
	"whatsupp-backend/entity"

	"gorm.io/gorm"
)

type MessageRepository struct {
	*repository[*entity.Message]
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{
		repository: &repository[*entity.Message]{
			db: db,
		},
	}
}

func (m *MessageRepository) WithTx(tx *gorm.DB) *MessageRepository {
	return &MessageRepository{
		repository: &repository[*entity.Message]{
			db: tx,
		},
	}

}
