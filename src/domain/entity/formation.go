package entity

import (
	"time"
)

// Formation is the schema of game team
type Formation struct {
	ID         uint64 `json:"id"`
	Goalkeeper uint64 `json:"goalkeeper"`
	Attackers  uint64 `json:"attackers"`
	Defenders  uint64 `json:"defenders"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type PublicFormation struct {
	ID         uint64 `json:"id"`
	Goalkeeper uint64 `json:"goalkeeper"`
	Attackers  uint64 `json:"attackers"`
	Defenders  uint64 `json:"defenders"`
}

type Formations []Formation

func (formations Formations) PublicFormations() []interface{} {
	result := make([]interface{}, len(formations))

	for index, formation := range formations {
		result[index] = formation.PublicFormation()
	}

	return result
}

func (f *Formation) PublicFormation() interface{} {
	return &PublicFormation {
		ID: f.ID,
		Goalkeeper: f.Goalkeeper,
		Attackers: f.Attackers,
		Defenders: f.Defenders,
	}
}

func (f *Formation) Prepare() {
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()
}


