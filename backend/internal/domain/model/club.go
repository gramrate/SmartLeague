package model

import (
	"time"

	"github.com/google/uuid"
)

type Club struct {
	ID          uuid.UUID
	CreatorID   uuid.UUID
	Name        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ClubUpdatePatch struct {
	Name        *string
	Description *string
}
