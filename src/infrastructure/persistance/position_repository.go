package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gegen07/cartola-university/domain/entity/scout"
	"github.com/gegen07/cartola-university/domain/repository"
)

type PositionRepository struct {
	db *sql.DB
}

func NewPositionRepository(db *sql.DB) *PositionRepository {
	return &PositionRepository{
		db: db,
	}
}

var _ repository.PositionRepository = &PositionRepository{}

func (p PositionRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]scout.Position, error) {
	positions := make([]scout.Position, 0)

	stmt, err := p.db.PrepareContext(ctx, query)

	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx, args...)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := scout.Position{}

		err := rows.Scan(
			&p.ID,
			&p.Description,
			&p.CreatedAt,
			&p.UpdatedAt,
			)

		if err != nil {
			return nil, err
		}

		positions = append(positions, p)
	}

	return positions, nil
}

func (p PositionRepository) GetAll(ctx context.Context, page int) ([]scout.Position, error) {
	limit := 10
	offset := (page-1) * limit

	query := `SELECT p.id, p.description, p.created_at, p.updated_at FROM positions p LIMIT $1 OFFSET $2`

	positions, err := p.fetch(ctx, query, limit, offset)

	if err != nil {
		return nil, err
	}

	return positions, nil
}

func (p PositionRepository) GetById(ctx context.Context, id uint64) (*scout.Position, error) {
	query := `SELECT id, description, created_at, updated_at FROM positions WHERE id = $1`

	positions, err := p.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}

	if len(positions) > 1 {
		//TODO error
	}

	return &positions[0], nil
}

func (p PositionRepository) Insert(ctx context.Context, position *scout.Position) (*scout.Position, error) {
	query := `INSERT INTO positions (description, created_at, updated_at) 
				VALUES ($1, $2, $3) RETURNING id`

	stmt, err := p.db.PrepareContext(ctx, query)

	if err != nil {
		return nil, err
	}

	err = stmt.QueryRowContext(ctx,
		position.Description,
		position.CreatedAt,
		position.UpdatedAt).Scan(&position.ID)

	if err != nil {
		return nil, err
	}

	return position, nil
}

func (p PositionRepository) Update(ctx context.Context, position *scout.Position) (*scout.Position, error) {
	query := `UPDATE positions SET (description, updated_at) = ($1, $2) WHERE id=$3`

	stmt, err := p.db.PrepareContext(ctx, query)

	if err != nil {
		return nil, err
	}

	res, err := stmt.ExecContext(ctx,
		position.Description,
		position.UpdatedAt,
		position.ID)

	if err != nil {
		return nil, err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected != 1 {
		err = fmt.Errorf("Weird Behavior %d", rowsAffected)
		return nil, err
	}

	return position, nil
}

func (p PositionRepository) Delete(ctx context.Context, id uint64) error {
	query := `DELETE FROM positions WHERE id=$1`

	stmt, err := p.db.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, id)

	if err != nil {
		return nil
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

