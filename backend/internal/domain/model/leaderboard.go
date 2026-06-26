package model

import "github.com/google/uuid"

type LeaderboardRow struct {
	ProfileID     *uuid.UUID // nil for guests
	GuestID       *uuid.UUID // nil for registered profiles
	GuestNickname *string
	Points        float64
}
