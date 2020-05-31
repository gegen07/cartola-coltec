package repository

import (
	"context"
	"github.com/gegen07/cartola-university/domain/entity/scout"
)

type ScoutRepository interface {
	GetAll(ctx context.Context, page int) ([]scout.Scout, error)
	GetByID(ctx context.Context, id uint64) (*scout.Scout, error)
	Insert(ctx context.Context, scout *scout.Scout) (*scout.Scout, error)
	Update(ctx context.Context, scout *scout.Scout) (*scout.Scout, error)
	Delete(ctx context.Context, id uint64) error
}
