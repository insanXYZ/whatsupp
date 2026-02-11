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
			db: db,
		},
	}
}

func (mr *MemberRepository) TakeByUserIdAndGroupId(ctx context.Context, userId, groupId int, dst *entity.Member) error {
	return mr.db.WithContext(ctx).Where("user_id = ? AND group_id = ?", userId, groupId).Take(dst).Error
}

func (mr *MemberRepository) WithTx(tx *gorm.DB) *MemberRepository {
	return &MemberRepository{
		repository: &repository[*entity.Member]{
			db: tx,
		},
	}
}
