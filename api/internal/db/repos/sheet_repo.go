package repos

import (
	"context"
	"fmt"

	"github.com/andebugan/sheetstruct/internal/app/dtos"
	"github.com/andebugan/sheetstruct/internal/app/repos"
	"gorm.io/gorm"
)

type GormSheetRepository struct {
	repos.IRepository[dtos.Sheet]
	DB *gorm.DB
}

func NewGormSheetRepository(db *gorm.DB) *GormSheetRepository {
	return &GormSheetRepository{
		IRepository: NewGormBaseRepositoty[dtos.Sheet](db),
		DB:          db,
	}
}

// Finds sheets connected to specific user
func (r *GormSheetRepository) FindByUser(ctx context.Context, uid uint64) ([]dtos.Sheet, error) {
	var sheets []dtos.Sheet

	err := r.DB.WithContext(ctx). 
				Where("uid = ?", uid).
				Find(&sheets).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get sheets: %w", err)
	}

	return sheets, nil
}
