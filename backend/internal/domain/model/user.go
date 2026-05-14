package model

import "SmartLeague/internal/domain/types"
import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID
	Nickname     string
	Name         string
	ShowName     bool
	Description  *string
	Email        string
	PasswordHash string
	ClubID       *uuid.UUID
	ClubState    types.ClubState
	Role         types.Role
}
