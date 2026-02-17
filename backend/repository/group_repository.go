package repository

import (
	"context"
	"whatsupp-backend/dto"
	"whatsupp-backend/entity"

	"gorm.io/gorm"
)

type GroupRepository struct {
	*repository[*entity.Group]
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{
		repository: &repository[*entity.Group]{
			db: db,
		},
	}
}

func (gr *GroupRepository) WithTx(tx *gorm.DB) *GroupRepository {
	return &GroupRepository{
		repository: &repository[*entity.Group]{
			db: tx,
		},
	}
}

func (gr *GroupRepository) SearchGroupAndUserWithName(
	ctx context.Context,
	userId int,
	name string,
) ([]dto.SearchGroupResponse, error) {

	searchPattern := "%" + name + "%"

	query := `
				SELECT 
				'user' as type,
				u.id,
				u.name,
				u.image,
				u.bio,
				NULL as group_type,
				( SELECT 
					g.id FROM groups g
					INNER JOIN members m1 ON m1.group_id = g.id
					INNER JOIN members m2 ON m2.group_id = g.id
					WHERE g.type = 'PERSONAL'
					AND m1.user_id = ?
					AND m2.user_id = u.id
					LIMIT 1
				) as group_id 
				FROM users u
				WHERE u.name ILIKE ?
				AND u.id != ?

        UNION ALL

        SELECT
            'group' AS type,
            g.id,
            g.name,
            g.image,
            g.bio,
            g.type AS group_type,
            g.id AS group_id
        FROM groups g
        WHERE g.type = 'GROUP'
        AND g.name ILIKE ?
    `

	var results []dto.SearchGroupResponse

	err := gr.db.
		WithContext(ctx).
		Raw(query, userId, searchPattern, userId, searchPattern).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (gr *GroupRepository) TakePersonalGroupBySenderAndReceiverId(ctx context.Context, senderId, receiverId int, dst *entity.Group) error {

	rawQuery := `
                SELECT *
                FROM groups g
                INNER JOIN members m1 ON m1.group_id = g.id
                INNER JOIN members m2 ON m2.group_id = g.id
                WHERE g.type = 'PERSONAL'
                AND m1.user_id = ?
                AND m2.user_id = ?
                LIMIT 1
	`

	tx := gr.db.WithContext(ctx).Raw(rawQuery, senderId, receiverId).Scan(dst)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (gr *GroupRepository) FindGroupTypeWithMemberUserId(ctx context.Context, userId int, dst []*dto.Group) error {
	return gr.db.WithContext(ctx).Preload("Members", "user_id = ?", userId).Find(dst, "type = ?", entity.GROUP).Error
}

func (gr *GroupRepository) TakeGroupWithGroupIdAndUserId(ctx context.Context, groupId, userId int) (*entity.Group, error) {
	group := new(entity.Group)

	gr.db.WithContext(ctx).Joins("JOIN members ON members.group_id = groups.id AND members.user_id = ?", userId).Take(group, "groups.id = ?", groupId)

	return group, nil
}

func (gr *GroupRepository) FindAllGroupWithMemberUserId(
	ctx context.Context,
	userId int,
) ([]dto.LoadRecentGroup, error) {

	query := `
        SELECT
            CASE 
                WHEN g.type = 'PERSONAL' THEN 'user'
                ELSE 'group'
            END AS type,

            COALESCE(u.id, g.id) AS id,
            COALESCE(u.name, g.name) AS name,
            COALESCE(u.image, g.image) AS image,
            COALESCE(u.bio, g.bio) AS bio,

            g.type AS group_type,
            g.id AS group_id

        FROM groups g
        JOIN members me 
            ON me.group_id = g.id

        LEFT JOIN members other
            ON other.group_id = g.id
            AND other.user_id != me.user_id
            AND g.type = 'PERSONAL'

        LEFT JOIN users u
            ON u.id = other.user_id

        WHERE me.user_id = ?
    `

	var result []dto.LoadRecentGroup

	err := gr.db.
		WithContext(ctx).
		Raw(query, userId).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}
