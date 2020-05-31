package scout

import (
	"time"
)

// Scout struct represents the stats of each player
type Scout struct {
	ID           	uint64
	Description  	string
	Scout 			string
	Points       	float64
	Positions		[]Position
	CreatedAt    	time.Time
	UpdatedAt 		time.Time
}

type RequestScout struct {
	ID           	uint64 `json:"id"`
	Description  	string `json:"description"`
	Scout 			string `json:"scout"`
	Points       	float64 `json:"points"`
	PositionsID		[]uint64 `json:"positions_id"`
}

type PublicScout struct {
	ID          uint64 `json:"id"`
	Description string `json:"description"`
	Scout 		string `json:"scout"`
	Points      float64 `json:"points"`
	Positions   []Position `json:"positions"`
}

type Scouts []Scout

func (scouts Scouts) PublicScouts() []interface{} {
	result := make([]interface{}, len(scouts))

	for i, scout := range scouts {
		result[i] = scout.PublicScout()
	}

	return result
}

func (s *Scout) PublicScout() *PublicScout {
	return &PublicScout{
		ID:          s.ID,
		Description: s.Description,
		Points:      s.Points,
		Scout: 		 s.Scout,
		Positions:   s.Positions,
	}
}

func (s *RequestScout) ToScout() *Scout {
	return &Scout{
		ID:          s.ID,
		Description: s.Description,
		Scout:       s.Scout,
		Points:      s.Points,
		Positions:   nil,
	}
}

func (s *Scout) Prepare() {
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
}
