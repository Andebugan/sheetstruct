package repos

import (
	"context"

	"github.com/andebugan/sheetstruct/internal/app/dtos"
)

// UserRepository defines user-specific operations
type IUserRepository interface {
	IRepository[dtos.User]
	
	// FindByEmail finds a user by email
	FindByEmail(ctx context.Context, email string) (*dtos.User, error)

	// FindByEmail finds a user by email
	FindByName(ctx context.Context, name string) (*dtos.User, error)

	// Returns error if user with same email exists
	CheckDuplicateName(ctx context.Context, name string) error

	// Returns error if user with same name exists
	CheckDuplicateEmail(ctx context.Context, email string) error
}
