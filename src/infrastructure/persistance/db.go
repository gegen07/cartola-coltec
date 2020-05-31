package persistence

import (
	"database/sql"
	"fmt"
	"github.com/gegen07/cartola-university/domain/repository"
	_ "github.com/lib/pq"
)

type Repositories struct {
	Formation repository.FormationRepository
	Scout repository.ScoutRepository
	Position repository.PositionRepository
	ScoutPosition repository.ScoutPositionRepository
	db *sql.DB
}

func NewRepositories(DbDriver, DbUser, DbPass, DbPort, DbHost, DbName string) (*Repositories, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPass)
	db, err := sql.Open(DbDriver, url)

	if err != nil {
		return nil, err
	}

	return &Repositories{
		Formation: NewFormationRepository(db),
		Scout: NewScoutRepository(db),
		Position: NewPositionRepository(db),
		ScoutPosition: NewScoutPositionRepository(db),
		db: db,
	}, nil
}

func (s *Repositories) Close() error {
	return s.db.Close()
}