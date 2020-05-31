package application

import (
	"context"
	"github.com/gegen07/cartola-university/domain/entity/scout"
	"github.com/gegen07/cartola-university/domain/repository"
	"golang.org/x/sync/errgroup"
	"github.com/sirupsen/logrus"
	"time"
)

type PositionApplicationInterface interface {
	GetAll(ctx context.Context, page int) ([]scout.Position, error)
	GetById(ctx context.Context, id uint64) (*scout.Position, error)
	Insert(ctx context.Context, position *scout.RequestPosition) (*scout.Position, error)
	Update(ctx context.Context, position *scout.Position) (*scout.Position, error)
	Delete(ctx context.Context, id uint64) error
}

type PositionApplication struct {
	repo                    repository.PositionRepository
	scoutPositionRepository repository.ScoutPositionRepository
	contextTimeout time.Duration
}

func NewPositionApplication(repo repository.PositionRepository,
							spRepo repository.ScoutPositionRepository,
							timeout time.Duration) *PositionApplication{
	return &PositionApplication{
		repo:                    repo,
		scoutPositionRepository: spRepo,
		contextTimeout:          timeout,
	}
}

func (pa *PositionApplication) fillPositionWithScout(c context.Context, data []scout.Position) ([]scout.Position, error) {
	g, ctx := errgroup.WithContext(c)
	type Result struct {
		scouts []scout.Scout
		positionID uint64
	}

	result := Result{}

	mapScouts := map[uint64][]scout.Scout{}

	for _, position := range data {
		mapScouts[position.ID] = make([]scout.Scout, 0)
	}

	chanScouts := make(chan Result)

	for positionID := range mapScouts  {
		positionID := positionID
		g.Go(func() error {
			res, err := pa.scoutPositionRepository.GetScoutsByPositionID(ctx, positionID)

			if err != nil {
				return err
			}

			result.scouts = res
			result.positionID = positionID

			chanScouts <- result

			return nil
		})
	}

	go func() {
		err := g.Wait()

		if err != nil {
			return
		}

		close(chanScouts)
	}()

	for result := range chanScouts {
		mapScouts[result.positionID] = result.scouts
	}

	if err := g.Wait(); err!= nil {
		return nil, err
	}

	for index, item := range data {
		if s, ok := mapScouts[item.ID]; ok {
			data[index].Scouts = s
		}
	}

	return data, nil
}

func (pa *PositionApplication) GetAll(ctx context.Context, page int) ([]scout.Position, error) {
	ctx, cancel := context.WithTimeout(ctx, pa.contextTimeout)
	defer cancel()

	positions, err := pa.repo.GetAll(ctx, page)

	if err != nil {
		return nil, err
	}

	return pa.fillPositionWithScout(ctx, positions)
}

func (pa *PositionApplication) GetById(ctx context.Context,id uint64) (*scout.Position, error) {
	ctx, cancel := context.WithTimeout(ctx, pa.contextTimeout)
	defer cancel()

	position, err := pa.repo.GetById(ctx, id)

	if err != nil {
		return nil, err
	}

	res, err := pa.scoutPositionRepository.GetScoutsByPositionID(ctx, position.ID)

	if err != nil {
		return nil, err
	}

	position.Scouts = res

	return position, nil
}

func (pa *PositionApplication) Insert(c context.Context, position *scout.RequestPosition) (*scout.Position, error) {
	g, ctx := errgroup.WithContext(c)

	ctx, cancel := context.WithTimeout(ctx, pa.contextTimeout)
	defer cancel()

	p := position.ToPosition()
	res, err := pa.repo.Insert(ctx, p)

	if err != nil {
		return nil, err
	}

	position.ID = res.ID

	for _, scoutID := range position.ScoutsID {
		g.Go(func() error {
			err := pa.scoutPositionRepository.AddRelation(c, scoutID, position.ID)

			if err != nil {
				return err
			}

			return nil
		})
	}

	go func() {
		err := g.Wait()

		if err != nil {
			logrus.Error(err)
			return
		}
	}()

	if err := g.Wait(); err != nil {
		return nil, err
	}

	scouts, err := pa.scoutPositionRepository.GetScoutsByPositionID(ctx, res.ID)

	if err != nil {
		return nil, err
	}

	res.Scouts= scouts

	return res, nil
}

func (pa *PositionApplication) Update(ctx context.Context, position *scout.Position) (*scout.Position, error) {
	ctx, cancel := context.WithTimeout(ctx, pa.contextTimeout)
	defer cancel()

	return pa.repo.Update(ctx, position)
}

func (pa *PositionApplication) Delete(c context.Context, id uint64) error {
	g, ctx := errgroup.WithContext(c)

	ctx, cancel := context.WithTimeout(ctx, pa.contextTimeout)
	defer cancel()

	res, err := pa.scoutPositionRepository.GetScoutsByPositionID(ctx, id)

	if err != nil {
		return err
	}

	g.Go(func() error {
		for _, r := range res {
			err := pa.scoutPositionRepository.DeleteRelation(ctx, r.ID, id)

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

	return pa.repo.Delete(ctx, id)
}