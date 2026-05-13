package model

import (
	"SmartLeague/internal/domain/types"
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	ID           uuid.UUID
	Nickname     string
	Name         string
	ShowName     bool
	Description  *string
	Email        string
	PasswordHash string
	Club         *string
	Role         types.Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ProfileUpdatePatch struct {
	Nickname     *string
	Name         *string
	ShowName     *bool
	Description  *string
	Club         *string
	Email        *string
	PasswordHash *string
	Role         *types.Role
}
