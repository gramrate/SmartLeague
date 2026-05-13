package dto

import (
	"SmartLeague/internal/domain/types"

	"github.com/google/uuid"
)

type Game struct {
	ID          uuid.UUID      `json:"id"`
	SeriesID    uuid.UUID      `json:"series_id"`
	Name        string         `json:"name"`
	Number      int            `json:"number"`
	Description *string        `json:"description,omitempty"`
	HostID      *uuid.UUID     `json:"host_id,omitempty"`
	Status      types.GameStatus `json:"status"`
}

type CreateGameRequest struct {
	SeriesID    uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
	Name        *string   `json:"name,omitempty" validate:"omitempty,min=1,max=200"`
	Number      int       `json:"number" validate:"required,min=1,max=100000"`
	Description *string   `json:"description,omitempty" validate:"omitempty,max=5000"`
	HostID      *uuid.UUID `json:"host_id,omitempty" validate:"omitempty,uuid"`
	Status      types.GameStatus `json:"status" validate:"min=0,max=2"`
}

type CreateGameResponse Game

type GetGameRequest struct {
	ID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
}

type GetGameResponse Game

type GetSeriesGamesRequest struct {
	SeriesID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
	Limit    *int      `json:"limit,omitempty" form:"limit" validate:"omitempty,min=1,max=100"`
	Offset   *int      `json:"offset,omitempty" form:"offset" validate:"omitempty,min=0"`
}

type GetSeriesGamesResponse struct {
	Items      []*Game         `json:"items"`
	Pagination PaginationInfo  `json:"pagination"`
}

type UpdateGameRequest struct {
	ID          uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
	Name        *string   `json:"name,omitempty" validate:"omitempty,min=1,max=200"`
	Description *string   `json:"description,omitempty" validate:"omitempty,max=5000"`
	HostID      *uuid.UUID `json:"host_id,omitempty" validate:"omitempty,uuid"`
	Status      *types.GameStatus `json:"status,omitempty" validate:"omitempty,min=0,max=2"`
}

type UpdateGameResponse Game

type SetGameParticipantsRequest struct {
	GameID         uuid.UUID   `json:"-" validate:"required,uuid" swaggerignore:"true"`
	ParticipantIDs []uuid.UUID `json:"participant_ids" validate:"required,min=1,max=100,dive,uuid"`
}

type UpsertGameResultsRequest struct {
	GameID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
	Rows   []GameResultRow `json:"rows" validate:"required,min=1,max=100,dive"`
}

type GameResultRow struct {
	ProfileID    uuid.UUID `json:"profile_id" validate:"required,uuid"`
	Place        *int      `json:"place,omitempty" validate:"omitempty,min=1,max=10"`
	Role         *string   `json:"role,omitempty" validate:"omitempty,min=1,max=100"`
	BestMove     bool      `json:"best_move"`
	FirstKilled  bool      `json:"first_killed"`
	Compensation int       `json:"compensation" validate:"min=0,max=1000000"`
	YellowCards  int       `json:"yellow_cards" validate:"min=0,max=10"`
	Removed      bool      `json:"removed"`
	ExtraPoints  int       `json:"extra_points" validate:"min=-1000000,max=1000000"`
	TotalPoints  int       `json:"total_points" validate:"min=-1000000,max=1000000"`
}

type LeaderboardRow struct {
	ProfileID uuid.UUID `json:"profile_id"`
	Points    int       `json:"points"`
}

type GetSeriesLeaderboardRequest struct {
	SeriesID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
	Limit    *int      `json:"limit,omitempty" form:"limit" validate:"omitempty,min=1,max=100"`
	Offset   *int      `json:"offset,omitempty" form:"offset" validate:"omitempty,min=0"`
}

type GetSeriesLeaderboardResponse struct {
	Items      []*LeaderboardRow `json:"items"`
	Pagination PaginationInfo    `json:"pagination"`
}

