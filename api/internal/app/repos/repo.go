package repos

import "context"

// Repository interface that defines basic CRUD operations for any entity
type IRepository[T any] interface {
	// Inserts a new record
	Create(ctx context.Context, entity *T) (*T, error)
	
	// Retrieves a record by its ID
	FindByID(ctx context.Context, id uint64) (*T, error)
	
	// Modifies an existing record
	Update(ctx context.Context, entity *T) (*T, error)
	
	// Performs a soft delete on a record
	Delete(ctx context.Context, id uint64) error
	
	// Retrieves all records with pagination
	GetAll(ctx context.Context, offset, limit int) ([]T, error)
}
