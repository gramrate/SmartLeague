package dto

import (
	"SmartLeague/internal/domain/types"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID       `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Nickname    string          `json:"nickname" example:"mishmish"`
	Name        string          `json:"name" example:"Ivan"`
	ShowName    bool            `json:"show_name" example:"true"`
	Description *string         `json:"description,omitempty" example:"About me"`
	Email       string          `json:"email" example:"user@example.com"`
	ClubID      *uuid.UUID      `json:"club_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	ClubState   types.ClubState `json:"club_state" example:"0"`
	Role        types.Role      `json:"role" example:"0"`
}

type RegisterUserRequest struct {
	Nickname *string     `json:"nickname,omitempty" validate:"omitempty,min=1,max=50" example:"mishmish"`
	Name     string      `json:"name" validate:"required,min=1,max=50" example:"Ivan"`
	ShowName *bool       `json:"show_name,omitempty" validate:"omitempty" example:"true"`
	Email    string      `json:"email" validate:"required,email,min=6,max=254" example:"user@example.com"`
	Password string      `json:"password" validate:"required,min=8,max=100" example:"SecurePass123!" format:"password"`
	Role     *types.Role `json:"role,omitempty" validate:"omitempty,role" swaggerignore:"true"`
}

type RegisterUserResponse struct {
	RefreshToken string `json:"-"`
	User
}

type GetUserRequest struct {
	ID uuid.UUID `json:"-" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000" swaggerignore:"true"`
}

type GetUserResponse User

type GetAllByFilterUsersRequest struct {
	Limit     *int             `json:"limit,omitempty" form:"limit" validate:"omitempty,min=1,max=200" example:"10"`
	Offset    *int             `json:"offset,omitempty" form:"offset" validate:"omitempty,min=0,max=10000" example:"0"`
	Role      *types.Role      `json:"role,omitempty" form:"role" validate:"omitempty,role" example:"0"`
	ClubState *types.ClubState `json:"club_state,omitempty" form:"club_state" validate:"omitempty,min=1,max=3" example:"2"`
	ClubQuery *string          `json:"club,omitempty" form:"club" validate:"omitempty,min=1,max=100" example:"smart"`
	Query     *string          `json:"q,omitempty" form:"q" validate:"omitempty,min=1,max=50" example:"mishmish"`
}

type PaginationInfo struct {
	TotalItems  int  `json:"total_items" example:"1250"`
	TotalPages  int  `json:"total_pages" example:"13"`
	CurrentPage int  `json:"current_page" example:"1"`
	HasNext     bool `json:"has_next" example:"true"`
	HasPrevious bool `json:"has_previous" example:"false"`
}

type GetAllByFilterUsersResponse struct {
	Items      []*User        `json:"items"`
	Pagination PaginationInfo `json:"pagination"`
}

type PlayerGame struct {
	ID         uuid.UUID        `json:"id"`
	SeriesID   uuid.UUID        `json:"series_id"`
	SeriesName string           `json:"series_name"`
	Name       string           `json:"name"`
	Number     int              `json:"number"`
	Status     types.GameStatus `json:"status"`
	CreatedAt  time.Time        `json:"created_at"`
}

type GetUserGamesRequest struct {
	UserID uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
	Limit  *int      `json:"limit,omitempty" form:"limit" validate:"omitempty,min=1,max=200"`
	Offset *int      `json:"offset,omitempty" form:"offset" validate:"omitempty,min=0,max=10000"`
}

type GetUserGamesResponse struct {
	Items      []*PlayerGame  `json:"items"`
	Pagination PaginationInfo `json:"pagination"`
}

type PlayerSeries struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	StartAt  time.Time `json:"start_at"`
	EndAt    time.Time `json:"end_at"`
	PriceRub int       `json:"price_rub"`
	IsRating bool      `json:"is_rating"`
	IsClosed bool      `json:"is_closed"`
}

type GetUserSeriesRequest struct {
	UserID     uuid.UUID `json:"-" validate:"required,uuid" swaggerignore:"true"`
	Limit      *int      `json:"limit,omitempty" form:"limit" validate:"omitempty,min=1,max=200"`
	Offset     *int      `json:"offset,omitempty" form:"offset" validate:"omitempty,min=0,max=10000"`
	Query      *string   `json:"q,omitempty" form:"q" validate:"omitempty,min=1,max=100"`
	From       *string   `json:"from,omitempty" form:"from" validate:"omitempty,len=10"`
	To         *string   `json:"to,omitempty" form:"to" validate:"omitempty,len=10"`
	IsRating   *bool     `json:"is_rating,omitempty" form:"is_rating"`
	ShowPast   *bool     `json:"show_past,omitempty" form:"show_past"`
	ShowClosed *bool     `json:"show_closed,omitempty" form:"show_closed"`
}

type GetUserSeriesResponse struct {
	Items      []*PlayerSeries `json:"items"`
	Pagination PaginationInfo  `json:"pagination"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email,min=6,max=254" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=8,max=100" example:"SecurePass123!" format:"password"`
}

type LoginUserResponse struct {
	RefreshToken string `json:"-"`
	User
}

type LogoutRequest struct {
	ID uuid.UUID `json:"id" validate:"required,uuid" swaggerignore:"true"`
}

type UpdateCurrentUserRequest struct {
	ID          uuid.UUID  `json:"-" validate:"required,uuid" swaggerignore:"true"`
	Nickname    *string    `json:"nickname,omitempty" validate:"omitempty,min=1,max=50" example:"mishmish"`
	Name        *string    `json:"name,omitempty" validate:"omitempty,min=1,max=50" example:"Ivan"`
	ShowName    *bool      `json:"show_name,omitempty" validate:"omitempty" example:"true"`
	Description *string    `json:"description,omitempty" validate:"omitempty,max=500" example:"About me"`
	ClubID      *uuid.UUID `json:"club_id,omitempty" validate:"omitempty,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type UpdateCurrentUserResponse User

type UpdateEachUserRequest struct {
	ID            uuid.UUID        `json:"-" validate:"required,uuid" swaggerignore:"true"`
	RequesterRole types.Role       `json:"-" validate:"required,role" swaggerignore:"true"`
	Nickname      *string          `json:"nickname,omitempty" validate:"omitempty,min=1,max=50" example:"mishmish"`
	Name          *string          `json:"name,omitempty" validate:"omitempty,min=1,max=50" example:"Ivan"`
	ShowName      *bool            `json:"show_name,omitempty" validate:"omitempty" example:"true"`
	Description   *string          `json:"description,omitempty" validate:"omitempty,max=500" example:"About me"`
	ClubID        *uuid.UUID       `json:"club_id,omitempty" validate:"omitempty,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email         *string          `json:"email,omitempty" validate:"omitempty,email,min=6,max=254" example:"user@example.com"`
	Password      *string          `json:"password,omitempty" validate:"omitempty,min=8,max=100" example:"SecurePass123!" format:"password"`
	ClubState     *types.ClubState `json:"club_state,omitempty" validate:"omitempty,min=0,max=4"`
	Role          *types.Role      `json:"role,omitempty" validate:"omitempty,role"`
}

type UpdateEachUserResponse User

type ChangePasswordRequest struct {
	ID          uuid.UUID `json:"id" validate:"required,uuid" swaggerignore:"true"`
	OldPassword string    `json:"old_password" validate:"required,min=8,max=100" example:"OldPass123!" format:"password"`
	NewPassword string    `json:"new_password" validate:"required,min=8,max=100" example:"NewPass456!" format:"password"`
}

type ChangePasswordResponse struct {
	RefreshToken string `json:"-" swaggerignore:"true"`
}

type DeleteUserRequest struct {
	ID uuid.UUID `json:"id" validate:"required,uuid" swaggerignore:"true"`
}
