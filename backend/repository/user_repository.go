package repository

import (
	"context"
	"whatsupp-backend/entity"

	"gorm.io/gorm"
)

type UserRepository struct {
	*repository[*entity.User]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		&repository[*entity.User]{
			DB: db,
		},
	}
}

func (u *UserRepository) TakeByEmail(ctx context.Context, email string, dst *entity.User) error {
	return u.DB.WithContext(ctx).Take(dst, "email = ?", email).Error
}

func (u *UserRepository) TakeByName(ctx context.Context, name string, dst *entity.User) error {
	return u.DB.WithContext(ctx).Take(dst, "name = ?", name).Error
}
