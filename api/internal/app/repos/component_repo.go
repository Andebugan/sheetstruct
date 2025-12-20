package repos

import (
	"context"

	"github.com/andebugan/sheetstruct/internal/app/dtos"
)

// IComponentRepository defines component-specific database operations
type IComponentRepository interface {
	IRepository[dtos.Component]

	FindForSheet(ctx context.Context, sid uint64) ([]dtos.Component, error)
}
