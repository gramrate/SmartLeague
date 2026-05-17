import { apiFetch } from "./client";
import type { Game, GameFull, PaginationInfo, UUID } from "../types/dto";
import type { MafiaRole } from "../types/enums";

export interface CreateGameRequest {
  series_id: UUID;
  name?: string | null;
  description?: string | null;
  host_id?: UUID | null;
  status: number;
}

export function createGame(seriesId: UUID, req: Omit<CreateGameRequest, "series_id">) {
  return apiFetch<Game>(`/api/v1/series/${seriesId}/games`, { method: "POST", body: JSON.stringify({ ...req }) });
}

export interface GetSeriesGamesResponse {
  items: Game[];
  pagination: PaginationInfo;
}

export function listGames(seriesId: UUID, params: { limit?: number; offset?: number }) {
  const q = new URLSearchParams();
  if (params.limit != null) q.set("limit", String(params.limit));
  if (params.offset != null) q.set("offset", String(params.offset));
  return apiFetch<GetSeriesGamesResponse>(`/api/v1/series/${seriesId}/games?${q.toString()}`, { method: "GET" });
}

export function getGame(id: UUID) {
  return apiFetch<Game>(`/api/v1/game/${id}`, { method: "GET" });
}

export function getGameFull(id: UUID) {
  return apiFetch<GameFull>(`/api/v1/game/${id}/full`, { method: "GET" });
}

export function updateGame(id: UUID, patch: Partial<Omit<CreateGameRequest, "series_id">>) {
  return apiFetch<Game>(`/api/v1/game/${id}`, { method: "PATCH", body: JSON.stringify(patch) });
}

export function deleteGame(id: UUID) {
  return apiFetch<void>(`/api/v1/game/${id}`, { method: "DELETE" });
}

export function setGameParticipants(gameId: UUID, participant_ids: UUID[]) {
  return apiFetch<void>(`/api/v1/game/${gameId}/participants`, { method: "POST", body: JSON.stringify({ participant_ids }) });
}

export type UpsertGameResultsRow = {
  profile_id: UUID;
  place?: number | null;
  role?: MafiaRole | null;
  best_move?: string | null;
  first_killed: boolean;
  compensation: number;
  yellow_cards: number;
  removed: number;
  victory_points: number;
  extra_points: number;
  total_points: number;
};

export function upsertGameResults(gameId: UUID, rows: UpsertGameResultsRow[]) {
  return apiFetch<void>(`/api/v1/game/${gameId}/results`, { method: "POST", body: JSON.stringify({ rows }) });
}
