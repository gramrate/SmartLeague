package types

type GameStatus int

const (
	GameStatusDraft GameStatus = iota
	GameStatusInProgress
	GameStatusFinished
)
