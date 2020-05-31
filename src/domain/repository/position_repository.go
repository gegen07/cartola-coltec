package repository

import (
	"context"
	"github.com/gegen07/cartola-university/domain/entity/scout"
)

type PositionRepository interface {
	GetAll(ctx context.Context, page int) ([]scout.Position, error)
	GetById(ctx context.Context, id uint64) (*scout.Position, error)
	Insert(ctx context.Context, position *scout.Position) (*scout.Position, error)
	Update(ctx context.Context, position *scout.Position) (*scout.Position, error)
	Delete(ctx context.Context, id uint64) error
}