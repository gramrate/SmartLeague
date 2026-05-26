package errorz

import "errors"

var (
	Unauthorized     = errors.New("permission denied")
	PermissionDenied = errors.New("permission denied")
	NoCookie         = errors.New("no mandatory cookies")
	InvalidRequest   = errors.New("invalid request")
	SeriesJoinClosed = errors.New("регистрация уже кончилась")
)
