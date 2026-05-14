export enum Role {
  User = 0,
  Moderator = 1,
  Admin = 2,
  SuperAdmin = 3
}

export enum ClubState {
  None = 0,
  Member = 1,
  Leader = 2,
  President = 3
}

export enum GameType {
  SportMafia = 0
}

export enum SeriesStatus {
  Closed = 0,
  Registration = 1,
  ClosedRegistration = 2,
  Games = 3
}

export enum GameStatus {
  Draft = 0,
  InProgress = 1,
  Finished = 2
}

export enum MafiaRole {
  Civilian = "civilian",
  Mafia = "mafia",
  Don = "don",
  Sheriff = "sheriff"
}
