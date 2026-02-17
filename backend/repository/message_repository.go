package repository

import (
	"context"
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

func (mr *MessageRepository) WithTx(tx *gorm.DB) *MessageRepository {
	return &MessageRepository{
		repository: &repository[*entity.Message]{
			db: tx,
		},
	}

}

func (mr *MessageRepository) GetMessages(ctx context.Context, conversationId int) ([]*entity.Message, error) {
	var messages []*entity.Message

	err := mr.db.WithContext(ctx).Preload("User").Find(&messages, "messages.conversation_id = ?", conversationId).Error

	return messages, err
}
