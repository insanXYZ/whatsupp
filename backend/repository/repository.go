package repository

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

type repository[T any] struct {
	db *gorm.DB
}

func (r *repository[T]) Create(ctx context.Context, model T) error {
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *repository[T]) Creates(ctx context.Context, model []T) error {
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *repository[T]) Update(ctx context.Context, model T) error {
	return r.db.WithContext(ctx).Updates(model).Error
}

func (r *repository[T]) TakeById(ctx context.Context, id any) (T, error) {
	var dst T

	err := r.db.WithContext(ctx).Take(&dst, "id = ?", id).Error

	return dst, err
}

func (r *repository[T]) Transaction(ctx context.Context, f func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return r.db.WithContext(ctx).Transaction(f, opts...)
}

func (r *repository[T]) DeleteById(ctx context.Context, id int) error {
	var model T
	return r.db.WithContext(ctx).Delete(&model, id).Error
}
