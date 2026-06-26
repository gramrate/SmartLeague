package model

import (
	"github.com/google/uuid"
	"time"
)

type Guest struct {
	ID        uuid.UUID
	SeriesID  uuid.UUID
	Nickname  string
	CreatedAt time.Time
}
