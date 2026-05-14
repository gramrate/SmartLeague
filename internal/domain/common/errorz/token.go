package errorz

import "errors"

var (
	InvalidToken  = errors.New("invalid token format")
	TokenNotFound = errors.New("token not found")
)
