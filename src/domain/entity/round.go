package entity

import (
	"time"
)

// Round struct represents the matches of the week
type Round struct {
	ID              int
	Matches         []Match
	RoundBeginDate  time.Time
	RoundFinishDate time.Time
}
