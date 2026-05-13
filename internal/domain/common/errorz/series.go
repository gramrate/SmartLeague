package errorz

import "errors"

var (
	SeriesNotFound = errors.New("series not found")
	GameNotFound   = errors.New("game not found")
)

