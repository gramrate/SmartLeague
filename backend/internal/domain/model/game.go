package model

import (
	"SmartLeague/internal/domain/types"
	"time"

	"github.com/google/uuid"
)

type Game struct {
	ID          uuid.UUID
	SeriesID    uuid.UUID
	Name        string
	Number      int
	Description *string
	HostID      *uuid.UUID
	Status      types.GameStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type GameUpdatePatch struct {
	Name        *string
	Description *string
	HostID      *uuid.UUID
	Status      *types.GameStatus
}

type GameResultRow struct {
	GameID       uuid.UUID
	ProfileID    uuid.UUID
	Place        *int
	Role         *string
	BestMove     bool
	FirstKilled  bool
	Compensation int
	YellowCards  int
	Removed      bool
	ExtraPoints  int
	TotalPoints  int
}
