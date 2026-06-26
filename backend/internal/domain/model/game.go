package model

import (
	"SmartLeague/internal/domain/types"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Game struct {
	ID                 uuid.UUID
	SeriesID           uuid.UUID
	Name               string
	Number             int
	Description        *string
	HostID             *uuid.UUID
	Status             types.GameStatus
	IsTournament       bool
	JudgeID            *uuid.UUID
	JudgeConfirmed     bool
	JudgeNickname      string
	JudgeName          string
	JudgeShowName      bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type GameUpdatePatch struct {
	Name        *string
	Description *string
	HostID      *uuid.UUID
	Status      *types.GameStatus
}

type GameDraft struct {
	GameID         uuid.UUID
	Rows           json.RawMessage
	JudgeID        *uuid.UUID
	JudgeConfirmed bool
	UpdatedAt      time.Time
}

type GameResultRow struct {
	GameID        uuid.UUID
	ProfileID     *uuid.UUID // nil for guests
	GuestID       *uuid.UUID // nil for registered profiles
	GuestNickname *string    // populated when reading guest rows
	Place         *int
	Role          *types.MafiaRole
	BestMove      *string
	FirstKilled   bool
	Compensation  float64
	YellowCards   float64
	Removed       float64
	VictoryPoints float64
	ExtraPoints   float64
	TotalPoints   float64
}
