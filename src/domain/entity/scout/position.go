package scout

import "time"

// Position struct represents the position of the player in a team
type Position struct {
	ID         	uint64
	Description string
	Scouts 		[]Scout
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RequestPosition struct {
	ID         	uint64 `json:"id"`
	Description string `json:"description"`
	ScoutsID  	[]uint64 `json:"scouts_id"`
}

type PublicPosition struct {
	ID          uint64 `json:"id"`
	Description string `json:"description"`
	Scouts 		[]Scout `json:"scouts"`
}

type Positions []Position

func (p Positions) PublicPositions() []interface{} {
	result := make([]interface{}, len(p))

	for i, position := range p {
		result[i] = position.PublicPosition()
	}

	return result
}

func (p *Position) PublicPosition() *PublicPosition {
	return &PublicPosition{
		ID: p.ID,
		Description: p.Description,
		Scouts: p.Scouts,
	}
}

func (p *RequestPosition) ToPosition() *Position {
	return &Position{
		ID:          p.ID,
		Description: p.Description,
	}
}

func (p *Position) Prepare() {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}
