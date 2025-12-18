package repos

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// GormBaseRepositoty provides a generic implementation of Repository interface
type GormBaseRepositoty[T any] struct {
	DB *gorm.DB
}

// NewGormBaseRepositoty creates a new GormBaseRepositoty instance
func NewGormBaseRepositoty[T any](db *gorm.DB) *GormBaseRepositoty[T] {
	return &GormBaseRepositoty[T]{DB: db}
}

// Create inserts a new record
func (r *GormBaseRepositoty[T]) Create(ctx context.Context, entity *T) (*T, error) {
	if entity == nil {
		return nil, errors.New("entity cannot be nil")
	}

	err := r.DB.WithContext(ctx).Create(entity).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create entity: %w", err)
	}

	return entity, nil
}

// FindByID retrieves a record by its ID
func (r *GormBaseRepositoty[T]) FindByID(ctx context.Context, id uint64) (*T, error) {
	var entity T
	err := r.DB.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("entity with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to find entity: %w", err)
	}
	return &entity, nil
}

// Update modifies an existing record
func (r *GormBaseRepositoty[T]) Update(ctx context.Context, entity *T) (*T, error) {
	if entity == nil {
		return nil, errors.New("entity cannot be nil")
	}

	err := r.DB.WithContext(ctx).Save(entity).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update entity: %w", err)
	}

	return entity, nil
}

// Delete performs a soft delete on a record
func (r *GormBaseRepositoty[T]) Delete(ctx context.Context, id uint64) error {
	var entity T
	result := r.DB.WithContext(ctx).Delete(&entity, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete entity: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("entity with ID %d not found", id)
	}
	return nil
}

// GetAll retrieves all records with pagination
func (r *GormBaseRepositoty[T]) GetAll(ctx context.Context, offset, limit int) ([]T, error) {
	var entities []T

	err := r.DB.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Find(&entities).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get entities: %w", err)
	}

	return entities, nil
}

// Transaction executes a function within a database transaction
func (r *GormBaseRepositoty[T]) Transaction(ctx context.Context, fn func(txRepo *GormBaseRepositoty[T]) error) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &GormBaseRepositoty[T]{DB: tx}
		return fn(txRepo)
	})
}
