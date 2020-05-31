package repository

import (
	"context"
	"github.com/gegen07/cartola-university/domain/entity/scout"
)

type ScoutPositionRepository interface {
	AddRelation(ctx context.Context, scoutId uint64, positionId uint64) error
	GetScoutsByPositionID(ctx context.Context, positionId uint64) ([]scout.Scout, error)
	GetPositionsByScoutID(ctx context.Context, scoutId uint64) ([]scout.Position, error)
	DeleteRelation(ctx context.Context, scoutId uint64, positionId uint64) error
}
