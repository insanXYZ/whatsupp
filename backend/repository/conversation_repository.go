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
				'PRIVATE' as conversation_type,
				u.id,
				u.name,
				u.image,
				u.bio,
				( SELECT 
					c.id FROM conversations c
					INNER JOIN members m1 ON m1.conversation_id = c.id
					INNER JOIN members m2 ON m2.conversation_id = c.id
					WHERE c.type = 'PRIVATE'
					AND m1.user_id = ?
					AND m2.user_id = u.id
					LIMIT 1
				) AS conversation_id ,
				true as have_joined
				FROM users u
				WHERE u.name ILIKE ?
				AND u.id != ?

        UNION ALL

        SELECT
            'GROUP' AS conversation_type,
            c.id,
            c.name,
            c.image,
            c.bio,
            c.id AS conversation_id,
						CASE
							WHEN m.id IS NOT NULL THEN true
							ELSE false
						END AS have_joined
        FROM conversations c
				LEFT JOIN members m on m.conversation_id = c.id AND m.user_id = ?
        WHERE c.type = 'GROUP'
        AND c.name ILIKE ?
    `

	var results []dto.SearchConversationResponse

	err := cr.db.
		WithContext(ctx).
		Raw(query, userId, searchPattern, userId, userId, searchPattern).
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

func (cr *ConversationRepository) TakeConversationByUserAndConversationId(ctx context.Context, userId, conversationId int) (*entity.Conversation, error) {
	conversation := new(entity.Conversation)

	err := cr.db.WithContext(ctx).Joins("JOIN users on user").Take(conversation, "conversations.id = ?", conversationId).Error

	return conversation, err
}

func (cr *ConversationRepository) TakeGroupConversationByUserAndConversationId(ctx context.Context, userId, conversationId int) (*entity.Conversation, error) {
	conversation := new(entity.Conversation)

	cr.db.WithContext(ctx).Joins("JOIN members ON members.conversation_id = conversations.id AND members.user_id = ?", userId).Preload("Members").Take(conversation, "conversations.id = ?", conversationId)

	return conversation, nil
}

func (cr *ConversationRepository) FindConversationsByUserId(
	ctx context.Context,
	userId int,
) ([]dto.LoadRecentConversation, error) {

	query := `
        SELECT
            COALESCE(u.id, c.id) AS id,
            COALESCE(u.name, c.name) AS name,
            COALESCE(u.image, c.image) AS image,
            COALESCE(u.bio, c.bio) AS bio,
            c.type AS conversation_type,
            c.id AS conversation_id,
						true AS have_joined

        FROM conversations c
        JOIN members me 
            ON me.conversation_id = c.id

        LEFT JOIN members other
            ON other.conversation_id = c.id
            AND other.user_id != me.user_id
            AND c.type = 'PRIVATE'

        LEFT JOIN users u
            ON u.id = other.user_id

        WHERE me.user_id = ?
    `

	var result []dto.LoadRecentConversation

	err := cr.db.
		WithContext(ctx).
		Raw(query, userId).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}
