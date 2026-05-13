package dto

import (
	"SmartLeague/internal/domain/types"
	"github.com/google/uuid"
)

type Profile struct {
	ID          uuid.UUID  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Nickname    string     `json:"nickname" example:"mishmish"`
	Name        string     `json:"name" example:"Ivan"`
	ShowName    bool       `json:"show_name" example:"true"`
	Description *string    `json:"description,omitempty" example:"About me"`
	Email       string     `json:"email" example:"user@example.com"`
	ClubID      *uuid.UUID `json:"club_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	ClubState   types.ClubState `json:"club_state" example:"0"`
	Role        types.Role `json:"role" example:"0"`
}

type CreateProfileRequest struct {
	Nickname    *string     `json:"nickname,omitempty" validate:"omitempty,min=1,max=100" example:"mishmish"`
	Name        string      `json:"name" validate:"required,min=1,max=100" example:"Ivan"`
	ShowName    *bool       `json:"show_name,omitempty" validate:"omitempty" example:"true"`
	Description *string     `json:"description,omitempty" validate:"omitempty,max=2000" example:"About me"`
	Email       string      `json:"email" validate:"required,email,min=6,max=254" example:"user@example.com"`
	Password    string      `json:"password" validate:"required,min=8,max=100" example:"SecurePass123!" format:"password"`
	ClubID      *uuid.UUID  `json:"club_id,omitempty" validate:"omitempty,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Role        *types.Role `json:"role,omitempty" validate:"omitempty,role" swaggerignore:"true"`
}

type CreateProfileResponse Profile

type GetProfileRequest struct {
	ID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
}

type GetProfileResponse Profile

type GetAllProfilesRequest struct {
	Limit  *int `json:"limit,omitempty" form:"limit" validate:"omitempty,min=1,max=100" example:"10"`
	Offset *int `json:"offset,omitempty" form:"offset" validate:"omitempty,min=0" example:"0"`
}

type GetAllProfilesResponse struct {
	Items      []*Profile       `json:"items"`
	Pagination PaginationInfo   `json:"pagination"`
}

type UpdateCurrentProfileRequest struct {
	ID          uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
	Nickname    *string   `json:"nickname,omitempty" validate:"omitempty,min=1,max=100" example:"mishmish"`
	Name        *string   `json:"name,omitempty" validate:"omitempty,min=1,max=100" example:"Ivan"`
	ShowName    *bool     `json:"show_name,omitempty" validate:"omitempty" example:"true"`
	Description *string   `json:"description,omitempty" validate:"omitempty,max=2000" example:"About me"`
	ClubID      *uuid.UUID `json:"club_id,omitempty" validate:"omitempty,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type UpdateCurrentProfileResponse Profile

type UpdateEachProfileRequest struct {
	ID            uuid.UUID   `json:"-" validate:"required,uuid" swaggerignore:"true"`
	RequesterRole types.Role  `json:"-" validate:"required,role" swaggerignore:"true"`
	Nickname      *string     `json:"nickname,omitempty" validate:"omitempty,min=1,max=100" example:"mishmish"`
	Name          *string     `json:"name,omitempty" validate:"omitempty,min=1,max=100" example:"Ivan"`
	ShowName      *bool       `json:"show_name,omitempty" validate:"omitempty" example:"true"`
	Description   *string     `json:"description,omitempty" validate:"omitempty,max=2000" example:"About me"`
	ClubID        *uuid.UUID  `json:"club_id,omitempty" validate:"omitempty,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email         *string     `json:"email,omitempty" validate:"omitempty,email,min=6,max=254" example:"user@example.com"`
	Password      *string     `json:"password,omitempty" validate:"omitempty,min=8,max=100" example:"SecurePass123!" format:"password"`
	Role          *types.Role `json:"role,omitempty" validate:"omitempty,role"`
}

type UpdateEachProfileResponse Profile

type DeleteProfileRequest struct {
	ID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
}
