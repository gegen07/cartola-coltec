package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gegen07/cartola-university/domain/entity/scout"
	"github.com/gegen07/cartola-university/domain/repository"
)

type ScoutRepository struct {
	db *sql.DB
}

func NewScoutRepository(db *sql.DB) *ScoutRepository {
	return &ScoutRepository{
		db: db,
	}
}

var _ repository.ScoutRepository = &ScoutRepository{}

func (s ScoutRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]scout.Scout, error) {
	scouts := make([]scout.Scout, 0)

	stmt, err := s.db.PrepareContext(ctx, query)

	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx, args...)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		s := scout.Scout{}

		err := rows.Scan(
			&s.ID,
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

func (s ScoutRepository) GetAll(ctx context.Context, page int) ([]scout.Scout, error) {
	limit := 10
	offset := (page-1)*limit
	query := `SELECT id, scout, description, points, created_at, updated_at FROM scouts`

	scouts, err := s.fetch(ctx, query, limit, offset)

	if err != nil {
		return nil, err
	}

	return scouts, nil
}

func (s ScoutRepository) GetByID(ctx context.Context, id uint64) (*scout.Scout, error) {
	query := `SELECT id, scout, description, points, created_at, updated_at FROM scouts WHERE id=$1`

	scouts, err := s.fetch(ctx, query, id)

	if err != nil {
		return nil, err
	}

	if len(scouts) > 1 {
		// TODO: error
	}

	return &scouts[0], nil
}

func (s ScoutRepository) Insert(ctx context.Context, scout *scout.Scout) (*scout.Scout, error) {
	query := `INSERT INTO scouts (scout, description, points, created_at, updated_at) 
				VALUES ($1, $2, $3, $4, $5) RETURNING id`

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	err = stmt.QueryRowContext(ctx,
				scout.Scout,
				scout.Description,
				scout.Points,
				scout.CreatedAt,
				scout.UpdatedAt).Scan(&scout.ID)

	if err != nil {
		return nil, err
	}

	return scout, nil
}

func (s ScoutRepository) Update(ctx context.Context, scout *scout.Scout) (*scout.Scout, error) {
	query := `UPDATE scouts SET (scout, description, points, updated_at) = ($1, $2, $3, $4) WHERE id=$5`

	stmt, err := s.db.PrepareContext(ctx, query)

	if err != nil {
		return nil, err
	}

	res, err := stmt.ExecContext(ctx, scout.Scout, scout.Description, scout.Points, scout.UpdatedAt)

	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return nil, err
	}

	if rowsAffected != 1 {
		err = fmt.Errorf("Weird Behavior %d", rowsAffected)
		return nil, err
	}

	return scout, nil
}

func (s ScoutRepository) Delete(ctx context.Context, id uint64) error {
	query := `DELETE FROM scouts WHERE id=$1`

	stmt, err := s.db.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, id)

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