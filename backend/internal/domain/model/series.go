package model

import (
	"SmartLeague/internal/domain/types"
	"time"

	"github.com/google/uuid"
)

type Series struct {
	ID           uuid.UUID
	ClubID       uuid.UUID
	CreatorID    uuid.UUID
	Name         string
	Description  string
	StartAt      time.Time
	EndAt        time.Time
	PriceRub     int
	IsRating     bool
	IsClubOnly   bool
	ShowToAll    bool
	IsClosed     bool
	IsTournament bool
	GameType     types.GameType
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type SeriesUpdatePatch struct {
	Name         *string
	Description  *string
	StartAt      *time.Time
	EndAt        *time.Time
	PriceRub     *int
	IsRating     *bool
	IsClubOnly   *bool
	ShowToAll    *bool
	IsClosed     *bool
	IsTournament *bool
	GameType     *types.GameType
}

type SeriesListItem struct {
	ID           uuid.UUID
	ClubID       uuid.UUID
	ClubName     string
	Name         string
	Description  string
	StartAt      time.Time
	EndAt        time.Time
	PriceRub     int
	IsRating     bool
	IsClubOnly   bool
	ShowToAll    bool
	IsClosed     bool
	IsTournament bool
	GamesCount   int
}

type JudgeRole int16

const (
	JudgeRoleSide JudgeRole = 0
	JudgeRoleMain JudgeRole = 1
)

type SeriesJudge struct {
	SeriesID  uuid.UUID
	ProfileID uuid.UUID
	Role      JudgeRole
	User      *User
}

type PlayerSeriesItem struct {
	Series
	IsJudge   bool
	JudgeRole *JudgeRole
}
