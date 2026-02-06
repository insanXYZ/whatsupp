package repository

import (
	"context"
	"whatsupp-backend/entity"

	"gorm.io/gorm"
)

type MemberRepository struct {
	*repository[*entity.Member]
}

func NewMemberRepository(db *gorm.DB) *MemberRepository {
	return &MemberRepository{
		repository: &repository[*entity.Member]{
			DB: db,
		},
	}
}

func (m *MemberRepository) TakeByUserIdAndGroupId(ctx context.Context, userId, groupId int, dst *entity.Member) error {
	return m.DB.WithContext(ctx).Where("user_id = ? AND group_id = ?", userId, groupId).Take(dst).Error
}
