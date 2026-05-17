package errorz

import "errors"

var (
	ClubNotFound       = errors.New("club not found")
	AlreadyInThisClub  = errors.New("player already in this club")
	AlreadyInOtherClub = errors.New("already in another club, leave current club first")
	ClubBanned         = errors.New("you are blocked in this club")
	ClubSelfAction     = errors.New("you cannot manage yourself in club")
	ClubRoleRestricted = errors.New("insufficient permissions for target member role")
)
