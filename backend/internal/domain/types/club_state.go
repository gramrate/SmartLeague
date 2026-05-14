package types

type ClubState int

const (
	ClubStateNone ClubState = iota
	ClubStateMember
	ClubStateLeader
	ClubStatePresident
)
