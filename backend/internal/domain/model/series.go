package model

import (
	"SmartLeague/internal/domain/types"
	"time"

	"github.com/google/uuid"
)

type Series struct {
	ID          uuid.UUID
	ClubID      uuid.UUID
	CreatorID   uuid.UUID
	Name        string
	Description string
	StartAt     time.Time
	EndAt       time.Time
	PriceRub    int
	IsClosed    bool
	GameType    types.GameType
	Status      types.SeriesStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SeriesUpdatePatch struct {
	Name        *string
	Description *string
	StartAt     *time.Time
	EndAt       *time.Time
	PriceRub    *int
	IsClosed    *bool
	GameType    *types.GameType
	Status      *types.SeriesStatus
}

type SeriesListItem struct {
	ID          uuid.UUID
	ClubID      uuid.UUID
	ClubName    string
	Name        string
	Description string
	StartAt     time.Time
	EndAt       time.Time
	IsClosed    bool
	GamesCount  int
}
