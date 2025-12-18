package repos

import (
	"context"
	"fmt"

	"github.com/andebugan/sheetstruct/internal/app"
	"github.com/andebugan/sheetstruct/internal/app/dtos"
	"github.com/andebugan/sheetstruct/internal/app/repos"
	"gorm.io/gorm"
)

// GormUserRepository implements user-specific operations
type GormUserRepository struct {
	repos.IRepository[dtos.User]
	DB   *gorm.DB
}

// NewGormUserRepository creates a new UserRepository instance
func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{
		IRepository: NewGormBaseRepositoty[dtos.User](db),
		DB: db,
	}
}

// FindByEmail finds a user with matching email
func (r *GormUserRepository) FindByEmail(ctx context.Context, email string) (*dtos.User, error) {
	var user dtos.User
	err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, app.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	return &user, nil
}

// FindByName finds a user with matching name
func (r *GormUserRepository) FindByName(ctx context.Context, name string) (*dtos.User, error) {
	var user dtos.User
	err := r.DB.WithContext(ctx).Where("name = ?", name).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, app.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user by name: %w", err)
	}
	return &user, nil
}

// Returns true if user with same name exists
func (r *GormUserRepository) CheckDuplicateName(ctx context.Context, name string) error {
	var user dtos.User
	err := r.DB.WithContext(ctx).Where("name = ?", name).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return fmt.Errorf("failed to find user by name: %w", err)
	}
	return app.ErrUserNameAlredyExists
}

// Returns true if user with same email exists
func (r *GormUserRepository) CheckDuplicateEmail(ctx context.Context, email string) error {
	var user dtos.User
	err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return fmt.Errorf("failed to find user by email: %w", err)
	}
	return app.ErrUserEmailAlredyExists
}
