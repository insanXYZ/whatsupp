package repository

import (
	"context"
	"whatsupp-backend/dto"
	"whatsupp-backend/entity"

	"gorm.io/gorm"
)

type ConversationRepository struct {
	*repository[*entity.Conversation]
}

func NewConversationRepository(db *gorm.DB) *ConversationRepository {
	return &ConversationRepository{
		repository: &repository[*entity.Conversation]{
			db: db,
		},
	}
}

func (cr *ConversationRepository) WithTx(tx *gorm.DB) *ConversationRepository {
	return &ConversationRepository{
		repository: &repository[*entity.Conversation]{
			db: tx,
		},
	}
}

func (cr *ConversationRepository) SearchConversationWithNameAndUserId(
	ctx context.Context,
	userId int,
	name string,
) ([]dto.SearchConversationResponse, error) {

	searchPattern := "%" + name + "%"

	query := `
				SELECT
					'PRIVATE' AS conversation_type,
					u.id,
					u.name,
					u.image,
					u.bio,
					c.id AS conversation_id,
					(c.id IS NOT NULL) AS have_joined
				FROM users u
				LEFT JOIN (
						SELECT
								conv.id,
								other.user_id
						FROM conversations conv
						JOIN members me
								ON me.conversation_id = conv.id
							 AND me.user_id = ?
						JOIN members other
								ON other.conversation_id = conv.id
							 AND other.user_id != ?
						WHERE conv.type = 'PRIVATE'
				) c ON c.user_id = u.id
				WHERE u.id != 1
					AND u.name ILIKE ?
				UNION ALL
				SELECT
					'GROUP' AS conversation_type,
					c.id,
					c.name,
					c.image,
					c.bio,
					c.id AS conversation_id,
					(m.id IS NOT NULL) AS have_joined
				FROM conversations c
				LEFT JOIN members m
					ON m.conversation_id = c.id
					AND m.user_id = ?
				WHERE c.type = 'GROUP'
					AND c.name ILIKE ?;
    `

	var results []dto.SearchConversationResponse

	err := cr.db.
		WithContext(ctx).
		Raw(query, userId, userId, searchPattern, userId, searchPattern).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (cr *ConversationRepository) TakePrivateConversationBySenderAndReceiverId(ctx context.Context, senderId, receiverId int) (*entity.Conversation, error) {

	conversation := new(entity.Conversation)

	rawQuery := `
                SELECT 
								c.id as id,
								c.name as name,
								c.bio as bio,
								c.type as type,
								c.image as image,
								c.created_at as created_at
                FROM conversations c
                INNER JOIN members m1 ON m1.conversation_id = c.id
                INNER JOIN members m2 ON m2.conversation_id = c.id
                WHERE c.type = ?
                AND m1.user_id = ?
                AND m2.user_id = ?
                LIMIT 1
	`

	tx := cr.db.WithContext(ctx).Raw(rawQuery, entity.CONV_TYPE_PRIVATE, senderId, receiverId).Scan(conversation)
	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return conversation, nil
}

func (cr *ConversationRepository) TakeGroupConversationByUserAndConversationId(ctx context.Context, userId, conversationId int) (*entity.Conversation, error) {
	conversation := new(entity.Conversation)

	err := cr.db.WithContext(ctx).Joins("JOIN members ON members.conversation_id = conversations.id AND members.user_id = ?", userId).Preload("Members").Take(conversation, "conversations.id = ?", conversationId).Error

	return conversation, err
}

func (cr *ConversationRepository) TakeGroupConversationLeftJoinMemberByUserAndConversationId(ctx context.Context, userId, conversationId int) (*entity.Conversation, error) {

	conversation := new(entity.Conversation)

	err := cr.db.WithContext(ctx).Preload("Members", "members.user_id = ?", userId).Take(conversation, "conversations.id = ?", conversationId).Error
	if err != nil {
		return nil, err
	}

	return conversation, nil
}

func (cr *ConversationRepository) FindConversationsByUserId(
	ctx context.Context,
	userId int,
) ([]*entity.Conversation, error) {

	var privateConversations []*entity.Conversation

	err := cr.db.Joins("join members on conversations.id = members.conversation_id and members.user_id = ?", userId).Preload("Members").Preload("Members.User").Find(&privateConversations, "conversations.type = 'PRIVATE'").Error
	if err != nil {
		return nil, err
	}

	var groupConversations []*entity.Conversation

	err = cr.db.Joins("join members on conversations.id = members.conversation_id and members.user_id = ?", userId).Preload("Members").Preload("Members.User").Find(&groupConversations, "conversations.type = 'GROUP'").Error
	if err != nil {
		return nil, err
	}

	results := append(privateConversations, groupConversations...)

	return results, nil

}

func (cr *ConversationRepository) TakeConversationByConversationAndUserId(ctx context.Context, conversationId, userId int) (*entity.Conversation, error) {
	conversation := new(entity.Conversation)

	query := `
			SELECT c.*
			FROM conversations c
			JOIN members m ON m.conversation_id = c.id
			JOIN users u ON u.id = m.user_id
			WHERE c.id = ? AND u.id = ?
	`

	tx := cr.db.WithContext(ctx).Raw(query, conversationId, userId).Scan(conversation)

	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return conversation, nil
}
