package repository

import (
	"whatsupp-backend/entity"

	"gorm.io/gorm"
)

type GroupMemberRepository struct {
	*repository[*entity.Group]
}

func NewGroupMessageRepository(db *gorm.DB) *GroupMemberRepository {
	return &GroupMemberRepository{
		repository: &repository[*entity.Group]{
			DB: db,
		},
	}
}
