package errorz

import "errors"

var (
	Unauthorized     = errors.New("forbidden: not authorized")
	PermissionDenied = errors.New("forbidden: you are not right enough")
	NoCookie         = errors.New("no mandatory cookies")
	InvalidRequest   = errors.New("invalid request")
	SeriesJoinClosed = errors.New("регистрация уже кончилась")
)
