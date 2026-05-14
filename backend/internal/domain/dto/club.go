package dto

import "github.com/google/uuid"

type Club struct {
	ID          uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	CreatorID   uuid.UUID `json:"creator_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string    `json:"name" example:"Smart League"`
	Description *string   `json:"description,omitempty" example:"Best club"`
}

type CreateClubRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=200" example:"Smart League"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=2000" example:"Best club"`
}

type CreateClubResponse Club

type GetClubRequest struct {
	ID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
}

type GetClubResponse Club

type GetAllClubsRequest struct {
	Limit  *int `json:"limit,omitempty" form:"limit" validate:"omitempty,min=1,max=100" example:"10"`
	Offset *int `json:"offset,omitempty" form:"offset" validate:"omitempty,min=0" example:"0"`
}

type GetAllClubsResponse struct {
	Items      []*Club        `json:"items"`
	Pagination PaginationInfo `json:"pagination"`
}

type UpdateClubRequest struct {
	ID          uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
	Name        *string   `json:"name,omitempty" validate:"omitempty,min=1,max=200" example:"Smart League"`
	Description *string   `json:"description,omitempty" validate:"omitempty,max=2000" example:"Best club"`
}

type UpdateClubResponse Club

type DeleteClubRequest struct {
	ID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
}

type GetClubMembersRequest struct {
	ClubID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
	Limit  *int      `json:"limit,omitempty" form:"limit" validate:"omitempty,min=1,max=100" example:"10"`
	Offset *int      `json:"offset,omitempty" form:"offset" validate:"omitempty,min=0" example:"0"`
}

type GetClubMembersResponse struct {
	Items      []*Profile     `json:"items"`
	Pagination PaginationInfo `json:"pagination"`
}

type JoinClubRequest struct {
	ProfileID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
	ClubID    uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
}

type LeaveClubRequest struct {
	ProfileID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
}
