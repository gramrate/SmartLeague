package dto

import (
	"SmartLeague/internal/domain/types"

	"github.com/google/uuid"
)

type Game struct {
	ID          uuid.UUID        `json:"id"`
	SeriesID    uuid.UUID        `json:"series_id"`
	Name        string           `json:"name"`
	Number      int              `json:"number"`
	Description *string          `json:"description,omitempty"`
	HostID      *uuid.UUID       `json:"host_id,omitempty"`
	Status      types.GameStatus `json:"status"`
}

type CreateGameRequest struct {
	SeriesID    uuid.UUID        `json:"-" validate:"required,uuid" swaggerignore:"true"`
	Name        *string          `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string          `json:"description,omitempty" validate:"omitempty,max=500"`
	HostID      *uuid.UUID       `json:"host_id,omitempty" validate:"omitempty,uuid"`
	Status      types.GameStatus `json:"status" validate:"min=0,max=2"`
}

type CreateGameDraftRequest struct {
	SeriesID    uuid.UUID         `json:"-" swaggerignore:"true"`
	Name        *string           `json:"name,omitempty"`
	Description *string           `json:"description,omitempty"`
	HostID      *uuid.UUID        `json:"host_id,omitempty"`
	Status      *types.GameStatus `json:"status,omitempty"`
}

type CreateGameResponse Game

type GetGameRequest struct {
	ID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
}

type GetGameResponse Game

type GameFull struct {
	Game
	ParticipantIDs []uuid.UUID     `json:"participant_ids"`
	Results        []GameResultRow `json:"results"`
}

type GetGameFullResponse GameFull

type GetSeriesGamesRequest struct {
	SeriesID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
	Limit    *int      `json:"limit,omitempty" form:"limit" validate:"omitempty,min=1,max=200"`
	Offset   *int      `json:"offset,omitempty" form:"offset" validate:"omitempty,min=0,max=10000"`
}

type GetSeriesGamesResponse struct {
	Items      []*Game        `json:"items"`
	Pagination PaginationInfo `json:"pagination"`
}

type UpdateGameRequest struct {
	ID          uuid.UUID         `json:"-" validate:"required,uuid" swaggerignore:"true"`
	Name        *string           `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string           `json:"description,omitempty" validate:"omitempty,max=500"`
	HostID      *uuid.UUID        `json:"host_id,omitempty" validate:"omitempty,uuid"`
	Status      *types.GameStatus `json:"status,omitempty" validate:"omitempty,min=0,max=2"`
}

type UpdateGameResponse Game

type DeleteGameRequest struct {
	ID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
}

type SetGameParticipantsRequest struct {
	GameID         uuid.UUID   `json:"-" validate:"required,uuid" swaggerignore:"true"`
	ParticipantIDs []uuid.UUID `json:"participant_ids" validate:"required,min=1,max=100,dive,uuid"`
}

type UpsertGameResultsRequest struct {
	GameID uuid.UUID       `json:"-" validate:"required,uuid" swaggerignore:"true"`
	Rows   []GameResultRow `json:"rows" validate:"required,min=1,max=100,dive"`
}

type GameResultRow struct {
	ProfileID     uuid.UUID        `json:"profile_id" validate:"required,uuid"`
	Place         *int             `json:"place,omitempty" validate:"omitempty,min=1,max=10"`
	Role          *types.MafiaRole `json:"role,omitempty"`
	BestMove      *string          `json:"best_move,omitempty" validate:"omitempty,max=50"`
	FirstKilled   bool             `json:"first_killed"`
	Compensation  float64          `json:"compensation" validate:"min=-1000000,max=1000000"`
	YellowCards   float64          `json:"yellow_cards" validate:"min=-1000000,max=1000000"`
	Removed       float64          `json:"removed" validate:"min=-1000000,max=1000000"`
	VictoryPoints float64          `json:"victory_points" validate:"min=-1000000,max=1000000"`
	ExtraPoints   float64          `json:"extra_points" validate:"min=-1000000,max=1000000"`
	TotalPoints   float64          `json:"total_points" validate:"min=-1000000,max=1000000"`
}

type ManageGameRow struct {
	Slot         int              `json:"slot"`
	ProfileID    *uuid.UUID       `json:"profile_id,omitempty"`
	Role         *types.MafiaRole `json:"role,omitempty"`
	BestMove     *string          `json:"best_move,omitempty"`
	Compensation float64          `json:"compensation"`
	YellowCards  float64          `json:"yellow_cards"`
	Removed      float64          `json:"removed"`
	ExtraPoints  float64          `json:"extra_points"`
	TotalPoints  float64          `json:"total_points"`
}

type SaveGameDraftRequest struct {
	GameID uuid.UUID       `json:"-" validate:"required,uuid" swaggerignore:"true"`
	Rows   []ManageGameRow `json:"rows"`
}

type PublishGameRequest struct {
	GameID uuid.UUID       `json:"-" validate:"required,uuid" swaggerignore:"true"`
	Rows   []ManageGameRow `json:"rows" validate:"required,min=10,max=10"`
}

type LeaderboardRow struct {
	ProfileID uuid.UUID `json:"profile_id"`
	Points    float64   `json:"points"`
}

type GetSeriesLeaderboardRequest struct {
	SeriesID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
	Limit    *int      `json:"limit,omitempty" form:"limit" validate:"omitempty,min=1,max=200"`
	Offset   *int      `json:"offset,omitempty" form:"offset" validate:"omitempty,min=0,max=10000"`
}

type GetSeriesLeaderboardResponse struct {
	Items      []*LeaderboardRow `json:"items"`
	Pagination PaginationInfo    `json:"pagination"`
}
