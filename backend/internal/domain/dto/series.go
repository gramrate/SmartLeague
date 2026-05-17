package dto

import (
	"SmartLeague/internal/domain/types"
	"time"

	"github.com/google/uuid"
)

type Series struct {
	ID          uuid.UUID          `json:"id"`
	ClubID      uuid.UUID          `json:"club_id"`
	CreatorID   *uuid.UUID         `json:"creator_id,omitempty"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	StartAt     time.Time          `json:"start_at"`
	EndAt       time.Time          `json:"end_at"`
	PriceRub    int                `json:"price_rub"`
	IsClosed    bool               `json:"is_closed"`
	GameType    types.GameType     `json:"game_type"`
	Status      types.SeriesStatus `json:"status"`
}

type CreateSeriesRequest struct {
	Name        string             `json:"name" validate:"required,min=1,max=200"`
	Description string             `json:"description" validate:"required,min=1,max=10000"`
	StartAt     time.Time          `json:"start_at" validate:"required"`
	EndAt       time.Time          `json:"end_at" validate:"required"`
	PriceRub    int                `json:"price_rub" validate:"min=0,max=100000000"`
	IsClosed    bool               `json:"is_closed"`
	GameType    types.GameType     `json:"game_type" validate:"eq=0"`
	Status      types.SeriesStatus `json:"status" validate:"min=0,max=3"`
}

type CreateSeriesResponse Series

type GetSeriesRequest struct {
	ID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
}

type GetSeriesResponse Series

type GetSeriesFullRequest struct {
	ID                 uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
	ParticipantsLimit  *int      `json:"participants_limit,omitempty" form:"participants_limit" validate:"omitempty,min=1,max=100"`
	ParticipantsOffset *int      `json:"participants_offset,omitempty" form:"participants_offset" validate:"omitempty,min=0"`
	GamesLimit         *int      `json:"games_limit,omitempty" form:"games_limit" validate:"omitempty,min=1,max=100"`
	GamesOffset        *int      `json:"games_offset,omitempty" form:"games_offset" validate:"omitempty,min=0"`
	LeaderboardLimit   *int      `json:"leaderboard_limit,omitempty" form:"leaderboard_limit" validate:"omitempty,min=1,max=100"`
	LeaderboardOffset  *int      `json:"leaderboard_offset,omitempty" form:"leaderboard_offset" validate:"omitempty,min=0"`
}

type GetSeriesFullResponse struct {
	Series       *Series                        `json:"series"`
	Participants *GetSeriesParticipantsResponse `json:"participants"`
	Games        *GetSeriesGamesResponse        `json:"games"`
	Leaderboard  *GetSeriesLeaderboardResponse  `json:"leaderboard"`
}

type GetClubSeriesRequest struct {
	ClubID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
	Limit  *int      `json:"limit,omitempty" form:"limit" validate:"omitempty,min=1,max=100"`
	Offset *int      `json:"offset,omitempty" form:"offset" validate:"omitempty,min=0"`
}

type GetClubSeriesResponse struct {
	Items      []*Series      `json:"items"`
	Pagination PaginationInfo `json:"pagination"`
}

type AllSeriesItem struct {
	ID          uuid.UUID `json:"id"`
	ClubID      uuid.UUID `json:"club_id"`
	ClubName    string    `json:"club_name"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	GamesCount  int       `json:"games_count"`
}

type GetAllSeriesRequest struct {
	Limit  *int `json:"limit,omitempty" form:"limit" validate:"omitempty,min=1,max=100"`
	Offset *int `json:"offset,omitempty" form:"offset" validate:"omitempty,min=0"`
}

type GetAllSeriesResponse struct {
	Items      []*AllSeriesItem `json:"items"`
	Pagination PaginationInfo   `json:"pagination"`
}

type UpdateSeriesRequest struct {
	ID          uuid.UUID           `json:"-" validate:"required,uuid" swaggerignore:"true"`
	Name        *string             `json:"name,omitempty" validate:"omitempty,min=1,max=200"`
	Description *string             `json:"description,omitempty" validate:"omitempty,min=1,max=10000"`
	StartAt     *time.Time          `json:"start_at,omitempty"`
	EndAt       *time.Time          `json:"end_at,omitempty"`
	PriceRub    *int                `json:"price_rub,omitempty" validate:"omitempty,min=0,max=100000000"`
	IsClosed    *bool               `json:"is_closed,omitempty"`
	GameType    *types.GameType     `json:"game_type,omitempty" validate:"omitempty,eq=0"`
	Status      *types.SeriesStatus `json:"status,omitempty" validate:"omitempty,min=0,max=3"`
}

type UpdateSeriesResponse Series

type DeleteSeriesRequest struct {
	ID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
}

type GetSeriesParticipantsRequest struct {
	SeriesID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
	Limit    *int      `json:"limit,omitempty" form:"limit" validate:"omitempty,min=1,max=100"`
	Offset   *int      `json:"offset,omitempty" form:"offset" validate:"omitempty,min=0"`
	Query    *string   `json:"q,omitempty" form:"q" validate:"omitempty,min=1,max=100"`
}

type GetSeriesParticipantsResponse struct {
	Items      []*User        `json:"items"`
	Pagination PaginationInfo `json:"pagination"`
}

type JoinSeriesRequest struct {
	SeriesID  uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
	ProfileID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
}

type LeaveSeriesRequest struct {
	SeriesID  uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
	ProfileID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
}
