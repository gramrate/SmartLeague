package types

type ClubState int

const (
	ClubStateNone ClubState = iota
	ClubStateMember
	ClubStateResident
	ClubStateLeader
	ClubStatePresident
)
