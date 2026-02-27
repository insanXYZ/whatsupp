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

func (mr *MemberRepository) TakeByUserIdAndConversationId(ctx context.Context, userId, conversationId int) (*entity.Member, error) {
	member := new(entity.Member)
	err := mr.db.WithContext(ctx).Take(member, "user_id = ? AND conversation_id = ?", userId, conversationId).Error
	return member, err
}

func (mr *MemberRepository) FindByUserId(ctx context.Context, userId int, members []*entity.Member) error {
	return mr.db.WithContext(ctx).Where("user_id = ?", userId).Find(members).Error
}

func (mr *MemberRepository) GetUserIdsWithConversationId(ctx context.Context, conversationId int) ([]int, error) {
	var members []*entity.Member

	err := mr.db.WithContext(ctx).Select("user_id").Where("conversation_id = ?", conversationId).Find(&members).Error
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

func (mr *MemberRepository) FindByConversationId(ctx context.Context, conversationId int) ([]*entity.Member, error) {
	var members []*entity.Member

	err := mr.db.WithContext(ctx).Preload("User").Where("conversation_id = ?", conversationId).Find(&members).Error

	return members, err
}

func (mr *MemberRepository) IsAdminConversationByConversationAndUserId(ctx context.Context, conversationId, userId int) (bool, error) {
	member := new(entity.Member)

	err := mr.db.WithContext(ctx).Take(member, "conversation_id = ? AND user_id = ?", conversationId, userId).Error
	if err != nil {
		return false, err
	}

	return member.Role == entity.MEMBER_TYPE_ADMIN, nil
}
