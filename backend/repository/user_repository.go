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
			db: db,
		},
	}
}

func (u *UserRepository) WithTx(tx *gorm.DB) *UserRepository {
	return &UserRepository{
		repository: &repository[*entity.User]{
			db: tx,
		},
	}
}

func (u *UserRepository) TakeByEmail(ctx context.Context, email string) (*entity.User, error) {
	dst := new(entity.User)

	err := u.db.WithContext(ctx).Take(dst, "email = ?", email).Error
	return dst, err
}

func (u *UserRepository) FindByName(ctx context.Context, name string) ([]*entity.User, error) {
	var dst []*entity.User

	err := u.db.WithContext(ctx).Where("name LIKE ?", name).Find(dst).Error
	return dst, err
}
