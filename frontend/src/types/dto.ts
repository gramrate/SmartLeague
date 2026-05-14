import { ClubState, GameStatus, GameType, MafiaRole, Role, SeriesStatus } from "./enums";

export type UUID = string;

export interface HTTPStatus {
  code: number;
  message: string;
}

export interface PaginationInfo {
  total_items: number;
  total_pages: number;
  current_page: number;
  has_next: boolean;
  has_previous: boolean;
}

export interface User {
  id: UUID;
  nickname: string;
  email: string;
  name: string;
  show_name: boolean;
  description?: string | null;
  club_id?: UUID | null;
  club_state: ClubState;
  role: Role;
}

export type Profile = User;

export interface Club {
  id: UUID;
  creator_id: UUID;
  name: string;
  description?: string | null;
}

export interface Series {
  id: UUID;
  club_id: UUID;
  creator_id?: UUID | null;
  name: string;
  scoring_rules: string;
  start_at: string;
  end_at: string;
  description?: string | null;
  price_rub: number;
  is_closed: boolean;
  game_type: GameType;
  status: SeriesStatus;
}

export interface Game {
  id: UUID;
  series_id: UUID;
  name: string;
  number: number;
  description?: string | null;
  host_id?: UUID | null;
  status: GameStatus;
}

export interface GameResultRow {
  profile_id: UUID;
  place?: number | null;
  role?: MafiaRole | null;
  best_move?: string | null;
  first_killed: boolean;
  compensation: number;
  yellow_cards: number;
  removed: number;
  extra_points: number;
  total_points: number;
}

export interface GameFull extends Game {
  participant_ids: UUID[];
  results: GameResultRow[];
}

export interface LeaderboardRow {
  profile_id: UUID;
  points: number;
}
