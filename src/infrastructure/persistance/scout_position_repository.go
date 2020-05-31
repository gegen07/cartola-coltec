package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gegen07/cartola-university/domain/entity/scout"
	"github.com/gegen07/cartola-university/domain/repository"
	"github.com/sirupsen/logrus"
)

type ScoutPositionRepository struct {
	db *sql.DB
}

var _ repository.ScoutPositionRepository = &ScoutPositionRepository{}

func NewScoutPositionRepository(db *sql.DB) *ScoutPositionRepository {
	return &ScoutPositionRepository{
		db: db,
	}
}

func (s ScoutPositionRepository) AddRelation(ctx context.Context, scoutId uint64, positionId uint64) error {
	query := `INSERT INTO scout_position (position_id, scout_id) VALUES ($1, $2)`

	stmt, err := s.db.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	logrus.Info(positionId, scoutId)
	_, err = stmt.ExecContext(ctx, positionId, scoutId)

	if err != nil {
		return err
	}

	return nil
}

func (s ScoutPositionRepository) GetScoutsByPositionID(ctx context.Context, positionId uint64) ([]scout.Scout, error) {
	query := `SELECT s.id, s.scout, s.description, s.points, s.created_at, s.updated_at 
				FROM scouts s INNER JOIN scout_position sp on sp.id = s.id WHERE (sp.position_id = $1)`

	stmt, err := s.db.PrepareContext(ctx, query)

	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx, positionId)

	scouts := make([]scout.Scout, 0)
	for rows.Next() {
		s := scout.Scout{}

		err := rows.Scan(&s.ID,
			&s.Scout,
			&s.Description,
			&s.Points,
			&s.CreatedAt,
			&s.UpdatedAt,
			)

		if err != nil {
			return nil, err
		}

		scouts = append(scouts, s)
	}

	return scouts, nil
}

func (s ScoutPositionRepository) GetPositionsByScoutID(ctx context.Context, scoutId uint64) ([]scout.Position, error) {
	query := `SELECT p.id, p.description, p.created_at, p.updated_at 
				FROM positions p INNER JOIN scout_position sp WHERE sp.scout_id = $1`

	stmt, err := s.db.PrepareContext(ctx, query)

	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx, scoutId)

	if err != nil {
		return nil, err
	}

	positions := make([]scout.Position, 0)
	for rows.Next() {
		p := scout.Position{}

		err = rows.Scan(
			&p.ID,
			&p.Description,
			&p.CreatedAt,
			&p.UpdatedAt,
			)

		if err != nil {
			return nil, err
		}
	}

	return positions, nil
}

func (s ScoutPositionRepository) DeleteRelation(ctx context.Context, position_id uint64, scout_id uint64) error {
	query := `DELETE FROM scout_position sp WHERE sp.position_id=$1 AND sp.scout_id=$2`

	stmt, err := s.db.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, position_id, scout_id)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		err = fmt.Errorf("Weird Behavior %d", rowsAffected)
		return err
	}

	return nil
}
