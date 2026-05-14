package model

import "SmartLeague/internal/domain/types"

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Name         string
	Surname      string
	Role         types.Role
}
