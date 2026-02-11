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

func (gr *GroupRepository) SearchGroupAndUserWithName(ctx context.Context, userId int, name string) ([]dto.SearchGroupResponse, error) {
	var userResults []dto.SearchGroupResponse
	var groupResults []dto.SearchGroupResponse

	searchPattern := "%" + name + "%"

	userQuery := `
        SELECT 
            'user' as type,
            u.id,
            u.name,
            u.image,
            u.bio,
            NULL as group_type,
            (
                SELECT g.id
                FROM groups g
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
        LIMIT 10
    `

	err := gr.db.WithContext(ctx).Raw(userQuery, userId, searchPattern, userId).Scan(&userResults).Error
	if err != nil {
		return nil, err
	}

	groupQuery := `
        SELECT 
            'group' as type,
            g.id,
            g.name,
            g.image,
            g.bio,
            g.type as group_type,
            g.id as group_id
        FROM groups g
        WHERE g.name ILIKE ?
        AND g.type = 'GROUP'
        LIMIT 10
    `

	err = gr.db.WithContext(ctx).Raw(groupQuery, searchPattern).Scan(&groupResults).Error
	if err != nil {
		return nil, err
	}

	results := append(userResults, groupResults...)

	return results, nil
}

func (gr *GroupRepository) TakePrivateGroupBySenderAndReceiverId(ctx context.Context, senderId, receiverId int, dst *entity.Group) error {

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

	return gr.db.WithContext(ctx).Raw(rawQuery, senderId, receiverId).Scan(dst).Error
}
