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

func (mr *MessageRepository) GetMessages(ctx context.Context, groupId int) ([]*entity.Message, error) {
	var messages []*entity.Message

	err := mr.db.WithContext(ctx).Joins("Member").Joins("Member.User").Where(`"Member"."group_id" = ?`, groupId).Find(&messages).Error

	return messages, err
}
