// SmartLeague API DTOs (based on swagger.yaml)

export enum ClubState {
  None = 0,
  Member = 1,
  Resident = 2,
  Leader = 3,
  President = 4,
}

export enum GameStatus {
  Draft = 0,
  InProgress = 1,
  Finished = 2,
}

export enum SeriesStatus {
  Closed = 0,
  Registration = 1,
  ClosedRegistration = 2,
  Games = 3,
}

export enum SystemRole {
  User = 0,
  Moderator = 1,
  Admin = 2,
  SuperAdmin = 3,
}

export type MafiaRole = "civilian" | "mafia" | "don" | "sheriff";

export interface PaginationInfo {
  current_page: number;
  has_next: boolean;
  has_previous: boolean;
  total_items: number;
  total_pages: number;
}

export interface Paged<T> {
  items: T[] | null;
  pagination: PaginationInfo;
}

export interface User {
  id: string;
  email?: string;
  name?: string;
  nickname?: string;
  description?: string;
  show_name?: boolean;
  club_id?: string;
  club_state?: ClubState;
  role?: SystemRole;
}

export interface Club {
  id: string;
  name: string;
  description?: string;
  creator_id?: string;
}

export interface Series {
  id: string;
  club_id: string;
  creator_id?: string;
  name: string;
  description?: string;
  start_at: string;
  end_at: string;
  status?: SeriesStatus;
  game_type?: number;
  is_closed?: boolean;
  price_rub?: number;
}

export interface AllSeriesItem {
  id: string;
  name: string;
  description?: string;
  club_id: string;
  club_name?: string;
  start_at: string;
  end_at: string;
  is_closed: boolean;
  games_count: number;
}

export interface Game {
  id: string;
  series_id: string;
  host_id?: string;
  name: string;
  description?: string;
  number: number;
  status: GameStatus;
}

export interface GameResultRow {
  profile_id: string;
  place?: number;
  role?: MafiaRole;
  best_move?: string;
  victory_points?: number;
  extra_points?: number;
  compensation?: number;
  total_points?: number;
  yellow_cards?: number;
  removed?: number;
  first_killed?: boolean;
}

export interface ManageGameRow {
  slot: number;
  profile_id?: string;
  role?: MafiaRole;
  best_move?: string;
  yellow_cards: number;
  removed: number;
  extra_points: number;
  total_points: number;
}

export interface GameFull extends Game {
  participant_ids: string[] | null;
  results: GameResultRow[] | null;
}

export interface LeaderboardRow {
  profile_id: string;
  points: number;
}

export interface SeriesFull {
  series: Series;
  participants: Paged<User>;
  games: Paged<Game>;
  leaderboard: Paged<LeaderboardRow>;
}

export interface PlayerGame {
  id: string;
  series_id: string;
  series_name: string;
  name: string;
  number: number;
  status: GameStatus;
  created_at?: string;
}

export interface PlayerSeries {
  id: string;
  name: string;
  start_at: string;
  end_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  name: string;
  nickname?: string;
  description?: string;
  show_name?: boolean;
}

export interface CreateClubRequest {
  name: string;
  description?: string;
}

export interface UpdateClubRequest {
  name?: string;
  description?: string;
}

export interface CreateSeriesRequest {
  name: string;
  description: string;
  start_at: string;
  end_at: string;
  price_rub?: number;
  game_type?: number;
  is_closed?: boolean;
  status?: SeriesStatus;
}

export interface UpdateSeriesRequest {
  name?: string;
  description?: string;
  start_at?: string;
  end_at?: string;
  price_rub?: number;
  is_closed?: boolean;
  status?: SeriesStatus;
}

export interface CreateGameRequest {
  name?: string;
  description?: string;
  host_id?: string;
  status?: GameStatus;
}

export interface UpdateCurrentUserRequest {
  name?: string;
  nickname?: string;
  description?: string;
  show_name?: boolean;
  club_id?: string;
}

export interface HTTPStatus {
  code: number;
  message: string;
}
