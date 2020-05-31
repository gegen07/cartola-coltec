package application

import (
	"context"
	"github.com/gegen07/cartola-university/domain/entity/scout"
	"github.com/gegen07/cartola-university/domain/repository"
	"golang.org/x/sync/errgroup"
	"time"
)

type ScoutApplication struct {
	repo repository.ScoutRepository
	scoutPositionRepo repository.ScoutPositionRepository
	contextTimeout time.Duration
}

type ScoutApplicationInterface interface {
	GetAll(ctx context.Context, page int) ([]scout.Scout, error)
	GetByID(ctx context.Context, id uint64) (*scout.Scout, error)
	Insert(ctx context.Context, scout *scout.RequestScout) (*scout.Scout, error)
	Update(ctx context.Context, scout *scout.Scout) (*scout.Scout, error)
	Delete(ctx context.Context, id uint64) error
}


func NewScoutApplication(repo repository.ScoutRepository,
						scoutPositionRepo repository.ScoutPositionRepository,
						timeout time.Duration) *ScoutApplication {
	return &ScoutApplication{
		repo:                    repo,
		scoutPositionRepo:       scoutPositionRepo,
		contextTimeout:          timeout,
	}
}

func (sa *ScoutApplication) fillScoutWithPosition(c context.Context, data []scout.Scout) ([]scout.Scout, error) {
	g, ctx := errgroup.WithContext(c)

	for _, scout := range data  {
		g.Go(func() error {
			res, err := sa.scoutPositionRepo.GetPositionsByScoutID(ctx, scout.ID)

			if err != nil {
				return err
			}

			scout.Positions = res

			return nil
		})
	}

	go func() {
		err := g.Wait()

		if err != nil {
			return
		}
	}()

	if err := g.Wait(); err!= nil {
		return nil, err
	}

	return data, nil
}

func (sa ScoutApplication) GetAll(ctx context.Context, page int) ([]scout.Scout, error) {
	ctx, cancel := context.WithTimeout(ctx, sa.contextTimeout)
	defer cancel()

	scouts, err := sa.repo.GetAll(ctx, page)

	if err != nil {
		return nil, err
	}

	return sa.fillScoutWithPosition(ctx, scouts)
}

func (sa ScoutApplication) GetByID(ctx context.Context, id uint64) (*scout.Scout, error) {
	ctx, cancel := context.WithTimeout(ctx, sa.contextTimeout)
	defer cancel()

	scout, err := sa.repo.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	positions, err := sa.scoutPositionRepo.GetPositionsByScoutID(ctx, scout.ID)

	if err != nil {
		return nil, err
	}

	scout.Positions = positions

	return scout, nil
}

func (sa ScoutApplication) Insert(c context.Context, reqScout *scout.RequestScout) (*scout.Scout, error) {
	g, ctx := errgroup.WithContext(c)

	ctx, cancel := context.WithTimeout(ctx, sa.contextTimeout)
	defer cancel()

	var s *scout.Scout
	s = reqScout.ToScout()
	s.Prepare()
	s, err := sa.repo.Insert(ctx, s)

	if err != nil {
		return nil, err
	}

	for _, positionID := range reqScout.PositionsID {
		g.Go(func() error {
			err := sa.scoutPositionRepo.AddRelation(c, s.ID, positionID)

			if err != nil {
				return err
			}

			return nil
		})
	}

	go func() {
		err := g.Wait()

		if err != nil {
			return
		}
	}()

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return s, nil
}

func (sa ScoutApplication) Update(ctx context.Context, scout *scout.Scout) (*scout.Scout, error) {
	ctx, cancel := context.WithTimeout(ctx, sa.contextTimeout)
	defer cancel()

	return sa.repo.Update(ctx, scout)
}

func (sa ScoutApplication) Delete(ctx context.Context, id uint64) error {
	g, ctx := errgroup.WithContext(ctx)

	ctx, cancel := context.WithTimeout(ctx, sa.contextTimeout)
	defer cancel()

	res, err := sa.scoutPositionRepo.GetPositionsByScoutID(ctx, id)

	if err != nil {
		return err
	}

	g.Go(func() error {
		for _, r := range res {
			err := sa.scoutPositionRepo.DeleteRelation(ctx, r.ID, id)

			if err != nil {
				return err
			}
		}

		return nil
	})

	go func() {
		err := g.Wait()

		if err != nil {
			return
		}
	}()

	if err := g.Wait(); err!= nil {
		return err
	}

	return sa.repo.Delete(ctx, id)
}
