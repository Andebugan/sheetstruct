package repos

import (
	"context"

	"github.com/andebugan/sheetstruct/internal/app/dtos"
)

// SheetRepository defines sheet-specific database operations
type ISheetRepository interface {
	IRepository[dtos.Sheet]

	FindByUser(ctx context.Context, uid uint64) ([]dtos.Sheet, error)
}
