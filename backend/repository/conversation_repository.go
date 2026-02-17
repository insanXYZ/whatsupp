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
				'user' as conversation_type,
				u.id,
				u.name,
				u.image,
				u.bio,
				( SELECT 
					c.id FROM conversation c
					INNER JOIN members m1 ON m1.group_id = c.id
					INNER JOIN members m2 ON m2.group_id = c.id
					WHERE c.type = 'PERSONAL'
					AND m1.user_id = ?
					AND m2.user_id = u.id
					LIMIT 1
				) as conversation_id 
				FROM users u
				WHERE u.name ILIKE ?
				AND u.id != ?

        UNION ALL

        SELECT
            'group' AS conversation_type,
            g.id,
            g.name,
            g.image,
            g.bio,
            g.id AS conversation_id
        FROM groups g
        WHERE g.type = 'GROUP'
        AND g.name ILIKE ?
    `

	var results []dto.SearchConversationResponse

	err := cr.db.
		WithContext(ctx).
		Raw(query, userId, searchPattern, userId, searchPattern).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (cr *ConversationRepository) TakePersonalConversationBySenderAndReceiverId(ctx context.Context, senderId, receiverId int) (*entity.Conversation, error) {

	conversation := new(entity.Conversation)

	rawQuery := `
                SELECT *
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

func (cr *ConversationRepository) TakeConversationByUserAndConversationId(ctx context.Context, userId, conversationId int) (*entity.Conversation, error) {
	conversation := new(entity.Conversation)

	err := cr.db.WithContext(ctx).Joins("JOIN users on user").Take(conversation, "conversations.id = ?", conversationId).Error

	return conversation, err
}

func (cr *ConversationRepository) TakeGroupConversationByUserAndConversationId(ctx context.Context, userId, conversationId int) (*entity.Conversation, error) {
	conversation := new(entity.Conversation)

	cr.db.WithContext(ctx).Joins("JOIN members ON members.group_id = groups.id AND members.user_id = ?", userId).Take(conversation, "conversations.id = ?", conversationId)

	return conversation, nil
}

func (cr *ConversationRepository) FindConversationsByUserId(
	ctx context.Context,
	userId int,
) ([]dto.LoadRecentGroup, error) {

	query := `
        SELECT
            COALESCE(u.id, c.id) AS id,
            COALESCE(u.name, c.name) AS name,
            COALESCE(u.image, c.image) AS image,
            COALESCE(u.bio, c.bio) AS bio,
            c.type AS conversation_type,
            c.id AS conversation_id

        FROM conversations c
        JOIN members me 
            ON me.group_id = c.id

        LEFT JOIN members other
            ON other.group_id = c.id
            AND other.user_id != me.user_id
            AND c.type = 'PERSONAL'

        LEFT JOIN users u
            ON u.id = other.user_id

        WHERE me.user_id = ?
    `

	var result []dto.LoadRecentGroup

	err := cr.db.
		WithContext(ctx).
		Raw(query, userId).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}
