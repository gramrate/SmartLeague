package types

type MafiaRole string

const (
	MafiaRoleCivilian MafiaRole = "civilian"
	MafiaRoleMafia    MafiaRole = "mafia"
	MafiaRoleDon      MafiaRole = "don"
	MafiaRoleSheriff  MafiaRole = "sheriff"
)
