package model

import "github.com/google/uuid"

type LeaderboardRow struct {
	ProfileID uuid.UUID
	Points    int
}

