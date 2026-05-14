package types

type SeriesStatus int

const (
	SeriesStatusClosed SeriesStatus = iota
	SeriesStatusRegistration
	SeriesStatusClosedRegistration
	SeriesStatusGames
)
