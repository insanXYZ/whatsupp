package repository

import (
	"context"

	"gorm.io/gorm"
)

type repository[T any] struct {
	DB *gorm.DB
}

func (r *repository[T]) Create(ctx context.Context, model T) error {
	return r.DB.WithContext(ctx).Create(model).Error
}

func (r *repository[T]) Update(ctx context.Context, model T) error {
	return r.DB.WithContext(ctx).Updates(model).Error
}
