package repository

import (
	"whatsupp-backend/entity"

	"gorm.io/gorm"
)

type GroupRepository struct {
	*repository[*entity.Group]
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{
		repository: &repository[*entity.Group]{
			DB: db,
		},
	}
}
