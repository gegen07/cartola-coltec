package repository

import (
	"context"
	"github.com/gegen07/cartola-university/domain/entity"
)

// Repository represent the round's repository contract
type FormationRepository interface {
	GetAll(ctx context.Context, page int) ([]entity.Formation, error)
	GetByID(ctx context.Context, id uint64) (*entity.Formation, error)
	Insert(ctx context.Context, formation *entity.Formation) (*entity.Formation, error)
	Update(ctx context.Context, formation *entity.Formation) (*entity.Formation, error)
	Delete(ctx context.Context, id uint64) error
}
