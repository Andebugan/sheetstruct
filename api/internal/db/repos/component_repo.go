package repos

import (
	"context"
	"fmt"

	"github.com/andebugan/sheetstruct/internal/app/dtos"
	"github.com/andebugan/sheetstruct/internal/app/repos"
	"gorm.io/gorm"
)

type GormComponentRepository struct {
	repos.IRepository[dtos.Component]
	DB   *gorm.DB
}

func NewGormComponentRepository(db *gorm.DB) *GormComponentRepository {
	return &GormComponentRepository{
		IRepository: NewGormBaseRepositoty[dtos.Component](db),
		DB: db,
	}
}

func (r *GormComponentRepository) FindForSheet(ctx context.Context, sid uint64) ([]dtos.Component, error) {
	var components []dtos.Component

	err := r.DB.WithContext(ctx). 
				Where("sid = ?", sid).
				Find(&components).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get components: %w", err)
	}

	return components, nil
}
