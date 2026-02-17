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

func (mr *MemberRepository) WithTx(tx *gorm.DB) *MemberRepository {
	return &MemberRepository{
		repository: &repository[*entity.Member]{
			db: tx,
		},
	}
}

func (mr *MemberRepository) TakeByUserIdAndGroupId(ctx context.Context, userId, groupId int, dst *entity.Member) error {
	err := mr.db.WithContext(ctx).Take(dst, "user_id = ? AND group_id = ?", userId, groupId).Error
	return err
}

func (mr *MemberRepository) FindByUserId(ctx context.Context, userId int, members []*entity.Member) error {
	return mr.db.WithContext(ctx).Where("user_id = ?", userId).Find(members).Error
}

func (mr *MemberRepository) GetUserIdsWithGroupId(ctx context.Context, groupId int) ([]int, error) {
	var members []*entity.Member

	err := mr.db.WithContext(ctx).Select("user_id").Where("group_id = ?", groupId).Find(&members).Error
	if err != nil {
		return nil, err
	}

	if len(members) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	userIds := make([]int, len(members))

	for i, member := range members {
		userIds[i] = member.UserID
	}

	return userIds, nil
}
