package repository

import (
	"whatsupp-backend/entity"

	"gorm.io/gorm"
)

type MemberRepository struct {
	*repository[*entity.Group]
}

func NewMemberRepository(db *gorm.DB) *MemberRepository {
	return &MemberRepository{
		repository: &repository[*entity.Group]{
			DB: db,
		},
	}
}
