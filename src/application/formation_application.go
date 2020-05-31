package application

import (
	"context"
	"github.com/gegen07/cartola-university/domain/entity"
	"github.com/gegen07/cartola-university/domain/repository"
)

type formationApplication struct {
	repo repository.FormationRepository
}

var _ FormationApplicationInterface = &formationApplication{}

type FormationApplicationInterface interface {
	GetAll(ctx context.Context, page int) ([]entity.Formation, error)
	GetByID(ctx context.Context, id uint64) (*entity.Formation, error)
	Insert(ctx context.Context, formation *entity.Formation) (*entity.Formation, error)
	Update(ctx context.Context, formation *entity.Formation) (*entity.Formation, error)
	Delete(ctx context.Context, formationId uint64) error
}

func (f formationApplication) GetAll(ctx context.Context, page int) ([]entity.Formation, error) {
	return f.repo.GetAll(ctx, page)
}

func (f formationApplication) GetByID(ctx context.Context, id uint64) (*entity.Formation, error) {
	return f.repo.GetByID(ctx, id)
}

func (f formationApplication) Insert(ctx context.Context, formation *entity.Formation) (*entity.Formation, error) {
	formation.Prepare()
	return f.repo.Insert(ctx, formation)
}

func (f formationApplication) Update(ctx context.Context, formation *entity.Formation) (*entity.Formation, error) {
	return f.repo.Update(ctx, formation)
}

func (f formationApplication) Delete(ctx context.Context, formationId uint64) error {
	return f.repo.Delete(ctx, formationId)
}
