package dto

import (
	"SmartLeague/internal/domain/types"
	"time"

	"github.com/google/uuid"
)

type Series struct {
	ID           uuid.UUID          `json:"id"`
	ClubID       uuid.UUID          `json:"club_id"`
	CreatorID    *uuid.UUID         `json:"creator_id,omitempty"`
	Name         string             `json:"name"`
	ScoringRules string             `json:"scoring_rules"`
	StartAt      time.Time          `json:"start_at"`
	EndAt        time.Time          `json:"end_at"`
	Description  *string            `json:"description,omitempty"`
	PriceRub     int                `json:"price_rub"`
	IsClosed     bool               `json:"is_closed"`
	GameType     types.GameType     `json:"game_type"`
	Status       types.SeriesStatus `json:"status"`
}

type CreateSeriesRequest struct {
	Name         string             `json:"name" validate:"required,min=1,max=200"`
	ScoringRules string             `json:"scoring_rules" validate:"required,min=1,max=10000"`
	StartAt      time.Time          `json:"start_at" validate:"required"`
	EndAt        time.Time          `json:"end_at" validate:"required"`
	Description  *string            `json:"description,omitempty" validate:"omitempty,max=5000"`
	PriceRub     int                `json:"price_rub" validate:"min=0,max=100000000"`
	IsClosed     bool               `json:"is_closed"`
	GameType     types.GameType     `json:"game_type" validate:"min=0,max=3"`
	Status       types.SeriesStatus `json:"status" validate:"min=0,max=3"`
}

type CreateSeriesResponse Series

type GetSeriesRequest struct {
	ID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
}

type GetSeriesResponse Series

type GetClubSeriesRequest struct {
	ClubID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
	Limit  *int      `json:"limit,omitempty" form:"limit" validate:"omitempty,min=1,max=100"`
	Offset *int      `json:"offset,omitempty" form:"offset" validate:"omitempty,min=0"`
}

type GetClubSeriesResponse struct {
	Items      []*Series       `json:"items"`
	Pagination PaginationInfo  `json:"pagination"`
}

type UpdateSeriesRequest struct {
	ID           uuid.UUID          `json:"-" validate:"required,uuid" swaggerignore:"true"`
	Name         *string            `json:"name,omitempty" validate:"omitempty,min=1,max=200"`
	ScoringRules *string            `json:"scoring_rules,omitempty" validate:"omitempty,min=1,max=10000"`
	StartAt      *time.Time         `json:"start_at,omitempty"`
	EndAt        *time.Time         `json:"end_at,omitempty"`
	Description  *string            `json:"description,omitempty" validate:"omitempty,max=5000"`
	PriceRub     *int               `json:"price_rub,omitempty" validate:"omitempty,min=0,max=100000000"`
	IsClosed     *bool              `json:"is_closed,omitempty"`
	GameType     *types.GameType    `json:"game_type,omitempty" validate:"omitempty,min=0,max=3"`
	Status       *types.SeriesStatus `json:"status,omitempty" validate:"omitempty,min=0,max=3"`
}

type UpdateSeriesResponse Series

type DeleteSeriesRequest struct {
	ID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
}

type GetSeriesParticipantsRequest struct {
	SeriesID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
	Limit    *int      `json:"limit,omitempty" form:"limit" validate:"omitempty,min=1,max=100"`
	Offset   *int      `json:"offset,omitempty" form:"offset" validate:"omitempty,min=0"`
}

type GetSeriesParticipantsResponse struct {
	Items      []*Profile      `json:"items"`
	Pagination PaginationInfo  `json:"pagination"`
}

type JoinSeriesRequest struct {
	SeriesID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
	ProfileID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
}

type LeaveSeriesRequest struct {
	SeriesID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
	ProfileID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
}
